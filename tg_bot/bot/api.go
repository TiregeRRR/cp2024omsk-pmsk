package bot

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (bw *BotWrapper) serveApi(ctx context.Context) {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Location"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.MaxMultipartMemory = 32 << 20
	router.GET("/get_transcriptions", bw.getTranscriptions)
	router.GET("/audio/:id", bw.getMinioLink)
	router.GET("/send_report/docx/unofficial/:id", bw.sendDocxUnofficialHandler)
	router.GET("/send_report/docx/official/:id", bw.sendDocxOfficialHandler)
	router.GET("/send_report/pdf/unofficial/:id", bw.sendPdfUnofficialHandler)
	router.GET("/send_report/pdf/official/:id", bw.sendPdfOfficialHandler)

	go router.Run("0.0.0.0:8888")
}

type getTranscriptionsResponse struct {
	ID        int64     `json:"id"`
	Status    int       `json:"status"`
	Name      string    `json:"name"`
	AudioLink string    `json:"audio_link"`
	CreatedAt time.Time `json:"created_at"`
}

func (bw *BotWrapper) getTranscriptions(c *gin.Context) {
	respPG, err := bw.psql.GetTranscribitions(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "can't get GetTranscribitions: " + err.Error(),
		})
		return
	}

	resp := make([]getTranscriptionsResponse, len(respPG))
	for i := range resp {
		resp[i] = getTranscriptionsResponse{
			ID:        respPG[i].ID,
			Status:    int(respPG[i].Status.Int32),
			Name:      "Совещание",
			AudioLink: "/audio/" + strconv.Itoa(int(respPG[i].ID)),
			CreatedAt: respPG[i].CreatedAt.Time,
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (bw *BotWrapper) getMinioLink(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id is not number" + err.Error(),
		})
		return
	}

	tr, err := bw.psql.GetTranscribition(c.Request.Context(), int64(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "get transcribition failed" + err.Error(),
		})
		return
	}

	f, err := bw.min.DownloadFile(c.Request.Context(), tr.AudioNameMinio.String, tr.AudioBucketMinio.String)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "download file failed: " + err.Error(),
		})
		return
	}

	b, err := io.ReadAll(f)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "read body failed: " + err.Error(),
		})
		return
	}

	switch {
	case strings.HasSuffix(tr.AudioNameMinio.String, ".mp3"):
		c.Data(http.StatusOK, "audio/mpeg", b)
	case strings.HasSuffix(tr.AudioNameMinio.String, ".ogg"):
		c.Data(http.StatusOK, "audio/ogg", b)
	case strings.HasSuffix(tr.AudioNameMinio.String, ".wav"):
		c.Data(http.StatusOK, "audio/vnd.wav", b)
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "unknown type: " + tr.AudioNameMinio.String,
		})
		return
	}
}

func (bw *BotWrapper) sendPdfUnofficialHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id is not number" + err.Error(),
		})
		return
	}

	ts, err := bw.psql.GetTranscribition(c.Request.Context(), int64(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "" + err.Error(),
		})
		return
	}

	bw.sendPdfUnofficial(c.Request.Context(), ts.TgUserID)
}

func (bw *BotWrapper) sendPdfOfficialHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id is not number" + err.Error(),
		})
		return
	}

	ts, err := bw.psql.GetTranscribition(c.Request.Context(), int64(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "" + err.Error(),
		})
		return
	}

	bw.sendPdfOfficial(c.Request.Context(), ts.TgUserID)
}

func (bw *BotWrapper) sendDocxUnofficialHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id is not number" + err.Error(),
		})
		return
	}

	ts, err := bw.psql.GetTranscribition(c.Request.Context(), int64(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "" + err.Error(),
		})
		return
	}

	bw.sendDocxUnofficial(c.Request.Context(), ts.TgUserID)
}

func (bw *BotWrapper) sendDocxOfficialHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id is not number" + err.Error(),
		})
		return
	}

	ts, err := bw.psql.GetTranscribition(c.Request.Context(), int64(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "" + err.Error(),
		})
		return
	}

	bw.sendDocxOfficial(c.Request.Context(), ts.TgUserID)
}
