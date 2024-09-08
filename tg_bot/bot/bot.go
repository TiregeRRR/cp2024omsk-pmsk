package bot

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/gulldan/cp2024omsk-pmsk/bot/minio"
	"github.com/gulldan/cp2024omsk-pmsk/config"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"

	postgres "github.com/gulldan/cp2024omsk-pmsk/bot/postgres/generated"
	pgxstdlib "github.com/jackc/pgx/v5/stdlib"
)

const (
	REPORT_DOCX_OFF   = "report_docx_off"
	REPORT_PDF_OFF    = "report_pdf_off"
	REPORT_DOCX_UNOFF = "report_docx_unoff"
	REPORT_PDF_UNOFF  = "report_pdf_unoff"

	StatusMessageWait = "Пожалуйста, ожидайте.\nТекущий статус задачи: %s"
	StatusMessageDone = "Задача выполнена."
)

type BotWrapper struct {
	log  *zerolog.Logger
	min  *minio.MinioClient
	psql *postgres.Queries
	cfg  *config.Config
	b    *bot.Bot
}

//go:embed postgres/sql/migrations/*.sql
var embedMigrations embed.FS

func New(cfg *config.Config) error {
	min, err := minio.NewMinioClient(cfg)
	if err != nil {
		return fmt.Errorf("minio connect failed: %w", err)
	}

	connStr := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(cfg.PostgresUsername, cfg.PostgresPassword),
		Host:   cfg.PostgresAddress,
		Path:   cfg.PostgresDatabase,
	}

	pg, err := pgxpool.New(context.Background(), connStr.String())
	if err != nil {
		return fmt.Errorf("postgres connect failed: %w", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(pgxstdlib.OpenDBFromPool(pg), "postgres/sql/migrations"); err != nil {
		panic(err)
	}

	var bw BotWrapper

	log := zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stdout})
	bw.log = &log
	bw.min = min
	bw.psql = postgres.New(pg)
	bw.cfg = cfg

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDebug(),
		bot.WithCheckInitTimeout(time.Minute),
		bot.WithMessageTextHandler(START, bot.MatchTypeExact, bw.startHandler),
		bot.WithDefaultHandler(bw.downloadHandler),
		bot.WithCallbackQueryDataHandler("report", bot.MatchTypePrefix, bw.reportCallbackQuery),
	}

	b, err := bot.New("6813542343:AAHfbZx-TjvJ3qf5B9L95X0kRkhO9emnWbU", opts...)
	if err != nil {
		panic(err)
	}
	bw.b = b

	bw.serveApi(ctx)
	b.Start(ctx)

	return nil
}

func (bw *BotWrapper) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   START_TEXT,
	})
	if err != nil {
		bw.log.Error().Err(err).Msg("send start message failed")

		return
	}
}

func (bw *BotWrapper) downloadHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	mimeType := ""
	fileID := ""

	if update.Message.Audio != nil {
		if update.Message.Audio.MimeType != "audio/mpeg" && update.Message.Audio.MimeType != "audio/ogg" && update.Message.Audio.MimeType != "audio/vnd.wav" {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   NOT_SUPPORTED_TYPE,
			})
			if err != nil {
				bw.log.Error().Err(err).Msg("send start message failed")
			}

			return
		}

		mimeType = update.Message.Audio.MimeType
		fileID = update.Message.Audio.FileID
	}

	if update.Message.Voice != nil {
		if update.Message.Voice.MimeType != "audio/mpeg" && update.Message.Voice.MimeType != "audio/ogg" && update.Message.Voice.MimeType != "audio/vnd.wav" {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   NOT_SUPPORTED_TYPE,
			})
			if err != nil {
				bw.log.Error().Err(err).Msg("send start message failed")
			}

			return
		}

		mimeType = update.Message.Voice.MimeType
		fileID = update.Message.Voice.FileID
	}

	if fileID == "" {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   NO_AUDIO_ATTACHED,
		})
		if err != nil {
			bw.log.Error().Err(err).Msg("send start message failed")
		}

		return
	}

	mf, err := b.GetFile(ctx, &bot.GetFileParams{
		FileID: fileID,
	})
	if err != nil {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   FAILED_TO_DOWNLOAD_FILE,
		})
		if err != nil {
			bw.log.Error().Err(err).Msg("send start message failed")
		}

		return
	}

	err = bw.psql.CreateUser(ctx, update.Message.Chat.ID)
	if err != nil {
		bw.log.Error().Err(err).Msg("CreateUser failed")
		return
	}

	file, err := downloadFileFromLink(ctx, b.FileDownloadLink(mf), mimeType)
	if err != nil {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   FAILED_TO_DOWNLOAD_FILE,
		})
		if err != nil {
			bw.log.Error().Err(err).Msg("send start message failed")
		}

		return
	}

	m, _ := bw.b.SendMessage(ctx, &bot.SendMessageParams{
		Text:   fmt.Sprintf(StatusMessageWait, "загружено."),
		ChatID: update.Message.Chat.ID,
	})

	trID, err := bw.psql.CreateTranscribition(ctx, postgres.CreateTranscribitionParams{
		TgUserID: update.Message.Chat.ID,
		MessageToEdit: pgtype.Int8{
			Int64: int64(m.ID),
			Valid: true,
		},
	})
	if err != nil {
		bw.log.Error().Err(err).Msg("CreateTranscribition failed")
		return
	}

	bw.updateStatus(ctx, StatusUploaded, trID, update.Message.Chat.ID, int64(m.ID))

	if err := bw.psql.UpdateCurrentBotID(ctx, postgres.UpdateCurrentBotIDParams{
		CurrentBotID: pgtype.Int8{
			Int64: trID,
			Valid: true,
		},
		TgUserID: update.Message.Chat.ID,
	}); err != nil {
		bw.log.Error().Err(err).Msg("UpdateCurrentBotID failed")
		return
	}

	fileName, bucket, err := bw.uploadFileToMinio(ctx, file)
	if err != nil {
		bw.log.Error().Err(err).Msg("UploadToMinio failed")
		return
	}

	if err := bw.psql.UpdateMinioLink(ctx, postgres.UpdateMinioLinkParams{
		AudioNameMinio: pgtype.Text{
			String: fileName,
			Valid:  true,
		},
		AudioBucketMinio: pgtype.Text{
			String: bucket,
			Valid:  true,
		},
		ID: trID,
	}); err != nil {
		bw.log.Error().Err(err).Msg("UpdateCurrentBotID failed")
		return
	}

	bw.updateStatus(ctx, StatusTranscription, trID, update.Message.Chat.ID, int64(m.ID))
	bw.startTranscription(ctx, trID, update.Message.Chat.ID, file, int64(m.ID))
}

func (bw *BotWrapper) uploadFileToMinio(ctx context.Context, file string) (fileName, bucket string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return "", "", fmt.Errorf("failed to open file: %w", err)
	}

	fs, err := f.Stat()
	if err != nil {
		return "", "", fmt.Errorf("failed to get file stat: %w", err)
	}

	err = bw.min.UploadFile(ctx, f, fs.Size(), fs.Name(), bw.min.GetAudioBucket())
	if err != nil {
		return "", "", fmt.Errorf("failed to upload file: %w", err)
	}

	return fs.Name(), bw.min.GetAudioBucket(), nil
}

func (bw *BotWrapper) updateStatus(ctx context.Context, status int, pgID, chatID, messageID int64) {
	if err := bw.psql.UpdateStatus(ctx, postgres.UpdateStatusParams{
		Status: pgtype.Int4{
			Int32: int32(status),
			Valid: true,
		},
		ID: int64(pgID),
	}); err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to update status")
	}

	var err error
	switch status {
	case StatusUploaded:
		_, err = bw.b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: int(messageID),
			Text:      fmt.Sprintf(StatusMessageWait, "загружено."),
		})
	case StatusTranscription:
		_, err = bw.b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: int(messageID),
			Text:      fmt.Sprintf(StatusMessageWait, "транскрибация."),
		})
	case StatusNers:
		_, err = bw.b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: int(messageID),
			Text:      fmt.Sprintf(StatusMessageWait, "выделение информации для отчета."),
		})
	case StatusReport:
		_, err = bw.b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: int(messageID),
			Text:      fmt.Sprintf(StatusMessageWait, "генерация отчета."),
		})
	case StatusDone:
		_, err = bw.b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: int(messageID),
			Text:      "Задача завершена",
			ReplyMarkup: &models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{
						{Text: "Официальный DOCX", CallbackData: REPORT_DOCX_OFF},
						{Text: "Официальный PDF", CallbackData: REPORT_PDF_OFF},
					}, {
						{Text: "Неофициальный DOCX", CallbackData: REPORT_DOCX_UNOFF},
						{Text: "Неофициальный PDF", CallbackData: REPORT_PDF_UNOFF},
					},
				},
			},
		})
	}

	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to edit message")
	}
}

func (bw *BotWrapper) sendPdfUnofficial(ctx context.Context, chatID int64) {
	user, err := bw.psql.GetUser(ctx, chatID)
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get user")
		return
	}

	tr, err := bw.psql.GetTranscribition(ctx, user.CurrentBotID.Int64)
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get transcibition")
		return
	}

	b, err := bw.unofficialReport(ctx, tr.ID, chatID, "pdf")
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get report")
		return
	}

	if _, err := bw.b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID:   chatID,
		Document: &models.InputFileUpload{Filename: "unofficial.pdf", Data: bytes.NewReader(b)},
		Caption:  "Document",
	}); err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to send message")
		return
	}
}

func (bw *BotWrapper) sendPdfOfficial(ctx context.Context, chatID int64) {
	user, err := bw.psql.GetUser(ctx, chatID)
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get user")
		return
	}

	tr, err := bw.psql.GetTranscribition(ctx, user.CurrentBotID.Int64)
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get transcibition")
		return
	}

	b, err := bw.officialReport(ctx, tr.ID, chatID, "pdf")
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get report")
		return
	}

	if _, err := bw.b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID:   chatID,
		Document: &models.InputFileUpload{Filename: "official.pdf", Data: bytes.NewReader(b)},
		Caption:  "Document",
	}); err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to send message")
		return
	}
}

func (bw *BotWrapper) sendDocxUnofficial(ctx context.Context, chatID int64) {
	user, err := bw.psql.GetUser(ctx, chatID)
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get user")
		return
	}

	tr, err := bw.psql.GetTranscribition(ctx, user.CurrentBotID.Int64)
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get transcibition")
		return
	}

	b, err := bw.unofficialReport(ctx, tr.ID, chatID, "docx")
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get report")
		return
	}

	if _, err := bw.b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID:   chatID,
		Document: &models.InputFileUpload{Filename: "unofficial.docx", Data: bytes.NewReader(b)},
		Caption:  "Document",
	}); err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to send message")
		return
	}
}

func (bw *BotWrapper) sendDocxOfficial(ctx context.Context, chatID int64) {
	user, err := bw.psql.GetUser(ctx, chatID)
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get user")
		return
	}

	tr, err := bw.psql.GetTranscribition(ctx, user.CurrentBotID.Int64)
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get transcibition")
		return
	}

	b, err := bw.officialReport(ctx, tr.ID, chatID, "docx")
	if err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to get report")
		return
	}

	if _, err := bw.b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID:   chatID,
		Document: &models.InputFileUpload{Filename: "official.docx", Data: bytes.NewReader(b)},
		Caption:  "Document",
	}); err != nil {
		bw.log.Error().Int64("id", chatID).Err(err).Msg("failed to send message")
		return
	}
}

func (bw *BotWrapper) reportCallbackQuery(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	switch update.CallbackQuery.Data {
	case REPORT_DOCX_OFF:
		bw.sendDocxOfficial(ctx, update.CallbackQuery.From.ID)
	case REPORT_PDF_OFF:
		bw.sendPdfOfficial(ctx, update.CallbackQuery.From.ID)
	case REPORT_DOCX_UNOFF:
		bw.sendDocxUnofficial(ctx, update.CallbackQuery.From.ID)
	case REPORT_PDF_UNOFF:
		bw.sendPdfUnofficial(ctx, update.CallbackQuery.From.ID)
	}
}
