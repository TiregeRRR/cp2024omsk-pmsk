package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	postgres "github.com/gulldan/cp2024omsk-pmsk/bot/postgres/generated"
)

type TranscriptionResp struct {
	ID      string `json:"identifier"`
	Message string `json:"message"`
}

func (bw *BotWrapper) startTranscription(ctx context.Context, pgID, chatID int64, file string, messageID int64) {
	go func() {
		bw.log.Info().Int64("chatID", chatID).Str("file", file).Msg("start transcription")

		v, err := bw.runTranscription(ctx, file)
		if err != nil {
			bw.log.Error().Err(err).Int64("chatID", chatID).Str("file", file).Msg("run transcription failed")

			return
		}

		b, err := json.Marshal(TaskResponseMarshal{
			Result: v.Result,
		})
		if err != nil {
			bw.log.Error().Err(err).Int64("chatID", chatID).Str("file", file).Msg("json marshal failed")

			return
		}

		if err := bw.psql.UpdateTranscription(ctx, postgres.UpdateTranscriptionParams{
			Transcription: pgtype.Text{
				String: string(b),
				Valid:  true,
			},
			ID: pgID,
		}); err != nil {
			bw.log.Error().Err(err).Int64("chatID", chatID).Str("file", file).Msg("update transcription failed")

			return
		}

		bw.log.Info().Any("resp", v).Msg("hello")

		bw.llamaComplete(ctx, string(b), pgID, chatID, messageID)

		bw.updateStatus(ctx, StatusDone, pgID, chatID, messageID)
	}()
}

type TaskResponseMarshal struct {
	Result struct {
		Segments []struct {
			Start   float64 `json:"start"`
			End     float64 `json:"end"`
			Text    string  `json:"text"`
			Speaker string  `json:"speaker"`
		} `json:"segments"`
	}
}

type TaskResponse struct {
	Status string `json:"status"`
	Result struct {
		Segments []struct {
			Start   float64 `json:"start"`
			End     float64 `json:"end"`
			Text    string  `json:"text"`
			Speaker string  `json:"speaker"`
		} `json:"segments"`
	} `json:"result"`
	Metadata struct {
		TaskType   string `json:"task_type"`
		TaskParams struct {
			Language             string `json:"language"`
			Task                 string `json:"task"`
			Model                string `json:"model"`
			Device               string `json:"device"`
			DeviceIndex          int    `json:"device_index"`
			Threads              int    `json:"threads"`
			BatchSize            int    `json:"batch_size"`
			ComputeType          string `json:"compute_type"`
			AlignModel           any    `json:"align_model"`
			InterpolateMethod    string `json:"interpolate_method"`
			ReturnCharAlignments bool   `json:"return_char_alignments"`
			AsrOptions           struct {
				BeamSize                  int     `json:"beam_size"`
				Patience                  float64 `json:"patience"`
				LengthPenalty             float64 `json:"length_penalty"`
				Temperatures              float64 `json:"temperatures"`
				CompressionRatioThreshold float64 `json:"compression_ratio_threshold"`
				LogProbThreshold          float64 `json:"log_prob_threshold"`
				NoSpeechThreshold         float64 `json:"no_speech_threshold"`
				InitialPrompt             any     `json:"initial_prompt"`
				SuppressTokens            []int   `json:"suppress_tokens"`
				SuppressNumerals          bool    `json:"suppress_numerals"`
			} `json:"asr_options"`
			VadOptions struct {
				VadOnset  float64 `json:"vad_onset"`
				VadOffset float64 `json:"vad_offset"`
			} `json:"vad_options"`
			MinSpeakers any `json:"min_speakers"`
			MaxSpeakers any `json:"max_speakers"`
		} `json:"task_params"`
		Language      string  `json:"language"`
		FileName      string  `json:"file_name"`
		URL           any     `json:"url"`
		Duration      float64 `json:"duration"`
		AudioDuration any     `json:"audio_duration"`
	} `json:"metadata"`
	Error any `json:"error"`
}

func (bw *BotWrapper) runTranscription(ctx context.Context, file string) (TaskResponse, error) {
	b, err := bw.min.DownloadFile(ctx, file, bw.min.GetAudioBucket())
	if err != nil {
		return TaskResponse{}, fmt.Errorf("download file failed: %w", err)
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	var fw io.Writer
	if fw, err = w.CreateFormFile("file", filepath.Base(file)); err != nil {
		return TaskResponse{}, err
	}

	if _, err = io.Copy(fw, b); err != nil {
		return TaskResponse{}, err
	}

	w.Close()

	whisper, err := url.Parse(bw.cfg.WhisperAddr + "/speech-to-text")
	if err != nil {
		return TaskResponse{}, fmt.Errorf("url parser failed: %w", err)
	}
	values := whisper.Query()
	values.Add("model", "large-v3")
	values.Add("language", "ru")
	whisper.RawQuery = values.Encode()

	newReq, err := http.NewRequest("POST", whisper.String(), &buf)
	if err != nil {
		return TaskResponse{}, err
	}

	newReq.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		return TaskResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TaskResponse{}, err
	}

	fmt.Println(string(body))

	defer resp.Body.Close()

	var respCreated TranscriptionResp
	if err := json.Unmarshal(body, &respCreated); err != nil {
		return TaskResponse{}, err
	}

loop:
	for {
		time.Sleep(time.Second)
		whisper, err = url.Parse(bw.cfg.WhisperAddr + "/task/" + respCreated.ID)
		if err != nil {
			return TaskResponse{}, err
		}
		newReq, err = http.NewRequest("GET", whisper.String(), http.NoBody)
		if err != nil {
			return TaskResponse{}, err
		}

		resp, err = http.DefaultClient.Do(newReq)
		if err != nil {
			return TaskResponse{}, err
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return TaskResponse{}, err
		}
		resp.Body.Close()

		fmt.Println(string(body))

		var r TaskResponse
		if err := json.Unmarshal(body, &r); err != nil {
			return TaskResponse{}, err
		}

		bw.log.Debug().Any("resp", r).Msg("waiting")

		switch r.Status {
		case "completed":
			return r, nil
		case "failed":
			return TaskResponse{}, errors.New("task failed")
		default:
			continue loop
		}
	}
}
