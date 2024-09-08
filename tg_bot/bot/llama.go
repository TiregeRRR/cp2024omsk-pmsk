package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	postgres "github.com/gulldan/cp2024omsk-pmsk/bot/postgres/generated"
)

const llamaSystemPrompt = `<|user|>
Вы - ИИ секретарь, чья задача конспектировать происходящие на различных рабочих встречах
Задача:
Создание структурированного и краткого протокола на основе JSON с субтитрами разговора с рабочих встреч.

Инструкции:
1. Создание точных транскрипций, без добавления лишней информации, с указанием только ключевых тем обсуждений и задач.
2. Распределение задач с указанием только ответственных лиц и сроков выполнения.
3. Формирование структурированных и кратких протоколов, включающих:
   - Список участников и их роли (только тех, кто активно принимал участие).
   - Повестку дня с кратким описанием обсуждаемых вопросов.
   - Задачи с указанием ответственных и сроков.
   - Дополнительные заметки только по сути.
   - Исключайте любые дублирования информации, ошибки или излишние детали.

Формат ответа:
{
    "name_report": "Совещание в государственной думе",
    "document_type": "pdf",
    "password": "pdf",
    "data": {
      "date": "2024-09-07T21:10:52.564Z",
      "time": "21:10:52.564Z",
      "duration": "P3D",
      "participants": [
        "SPEAKER_00",
        "SPEAKER_01",
      ],
      "agenda": [
        "1. Обсуждение итогов прошедшего дня и ощущений от совещания.",
      ],
      "blocks": [
        {
          "name_block": "Задачи",
          "proposals": [
            {
              "text": "Сделать все возможное для защиты экономики",
              "context": "Не определен",
              "audio_time": {
                "start": 133.98,
                "end": 140.025
              }
            },
            {
              "text": "нужно собраться и нужно организовать работу",
              "context": "Сегодня",
              "audio_time": {
                "start": 171.105,
                "end": 182.45
              }
            }
          ]
        },
        {
          "name_block": "Обратная связь",
          "proposals": [
            {
              "text": "- SPEAKER_01 признал, что нужно нужно собраться и организовать работку как надо",
              "context": "",
              "audio_time": {
                "start": 171.105,
                "end": 182.45
              }
            }
          ]
        }
      ],
      "audio_times": [
        {
          "start": 0,
          "end": 182.45
        }
      ]
    }
  }

На вход подается только JSON с субтитрами в следующем формате:
{
    "result": {
        "segments": [
          {
            "start": 1.088,
            "end": 12.716,
            "text": " Вы знаете, что в Государственную Думу у меня вынесено предложение о назначении вас на должность председателя правительства Российской Федерации.",
            "speaker": "SPEAKER_01"
          },
          {
            "start": 12.756,
            "end": 20.001,
            "text": "Совсем недавно мы встречались с коллегами и оценивали работу правительства за предыдущие годы.",
            "speaker": "SPEAKER_01"
          },
          {
            "start": 20.061,
            "end": 28.807,
            "text": "Сделано в сложных условиях немало, и мне кажется, что было бы правильно, если бы",
            "speaker": "SPEAKER_01"
          },
          {
            "start": 36.504,
            "end": 42.048,
            "text": " Мы с вами говорили и о структуре, говорили о персонале.",
            "speaker": "SPEAKER_01"
          },
          {
            "start": 42.088,
            "end": 44.89,
            "text": "В целом, думаю, мы на правильном пути.",
            "speaker": "SPEAKER_01"
          },
          {
            "start": 44.95,
            "end": 65.183,
            "text": "И очень надеюсь на то, что депутаты Государственной Думы, а вы не так давно были в Госдуме, отчитывались, они знают, что и как правительство, и вами, как председателям правительства, сделано за последние годы, оценят должным образом и поддержат вас в ходе ваших консультаций предстоящих сегодня в Госдуме.",
            "speaker": "SPEAKER_01"
          },
          {
            "start": 94.77,
            "end": 97.152,
            "text": " Спасибо, уважаемый Владимир Владимирович.",
            "speaker": "SPEAKER_00"
          },
          {
            "start": 97.192,
            "end": 113.905,
            "text": "Хочу в первую очередь поблагодарить вас за доверие, которое вы оказали мне, за задачи, которые вы поставили перед Федеральным Собранием в своем послании, и, конечно, те национальные цели развития, которые были указаны в новом майском указе.",
            "speaker": "SPEAKER_00"
          },
          {
            "start": 114.885,
            "end": 118.468,
            "text": " Это ориентир и приоритеты в работе правительства.",
            "speaker": "SPEAKER_00"
          },
          {
            "start": 118.508,
            "end": 121.871,
            "text": "Хочу вас заверить, что никаких пауз в работе правительства не будет.",
            "speaker": "SPEAKER_00"
          },
          {
            "start": 121.911,
            "end": 124.593,
            "text": "Мы будем продолжать текущую работу.",
            "speaker": "SPEAKER_00"
          },
          {
            "start": 124.613,
            "end": 133.92,
            "text": "Также считаю, что мы должны обеспечить преемственность по всем национальным целям, которые были до этого, 204 и 474 указе.",
            "speaker": "SPEAKER_00"
          },
          {
            "start": 133.98,
            "end": 140.025,
            "text": "Сделаем все для развития нашей экономики, чтобы оправдать доверие наших людей.",
            "speaker": "SPEAKER_00"
          },
          {
            "start": 140.065,
            "end": 144.208,
            "text": "И уверен, что под вашим руководством мы все задачи, которые поставлены, решим.",
            "speaker": "SPEAKER_00"
          },
          {
            "start": 145.229,
            "end": 151.012,
            "text": " Мы с вами вместе и с коллегами из правительства формулировали национальные цели развития.",
            "speaker": "SPEAKER_01"
          },
          {
            "start": 151.093,
            "end": 158.578,
            "text": "Это, конечно, главное, к чему мы должны стремиться, к реализации этих целей по всем направлениям.",
            "speaker": "SPEAKER_01"
          },
          {
            "start": 158.658,
            "end": 163.461,
            "text": "И, как показывает практика последних лет, в целом у нас…",
            "speaker": "SPEAKER_01"
          },
          {
            "start": 166.082,
            "end": 171.025,
            "text": " Получается добиваться тех результатов, которые нужны стране.",
            "speaker": "SPEAKER_01"
          },
          {
            "start": 171.105,
            "end": 182.45,
            "text": "А в сегодняшних непростых условиях, конечно, нужно собраться и нужно организовать работу именно так, как мы с вами договорились на последней встрече с правительством, работать без пауз.",
            "speaker": "SPEAKER_01"
          }
        ]
      }
    }
}
На выход должен выдаваться только желаемый JSON.
<|end|>
<|assistant|>`

type ReportedRequest struct {
	NameReport   string `json:"name_report"`
	DocumentType string `json:"document_type"`
	Password     string `json:"password"`
	Data         struct {
		Date         time.Time `json:"date"`
		Time         string    `json:"time"`
		Duration     string    `json:"duration"`
		Participants []string  `json:"participants"`
		Agenda       []string  `json:"agenda"`
		Blocks       []struct {
			NameBlock string `json:"name_block"`
			Proposals []struct {
				Text      string `json:"text"`
				Context   string `json:"context"`
				AudioTime struct {
					Start string `json:"start"`
					End   string `json:"end"`
				} `json:"audio_time"`
			} `json:"proposals"`
		} `json:"blocks"`
		AudioTimes []struct {
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"audio_times"`
	} `json:"data"`
}

type CompletionResp struct {
	Content            string `json:"content"`
	GenerationSettings struct {
		FrequencyPenalty float64       `json:"frequency_penalty"`
		IgnoreEos        bool          `json:"ignore_eos"`
		LogitBias        []interface{} `json:"logit_bias"`
		Mirostat         int           `json:"mirostat"`
		MirostatEta      float64       `json:"mirostat_eta"`
		MirostatTau      float64       `json:"mirostat_tau"`
		Model            string        `json:"model"`
		NCtx             int           `json:"n_ctx"`
		NKeep            int           `json:"n_keep"`
		NPredict         int           `json:"n_predict"`
		NProbs           int           `json:"n_probs"`
		PenalizeNl       bool          `json:"penalize_nl"`
		PresencePenalty  float64       `json:"presence_penalty"`
		RepeatLastN      int           `json:"repeat_last_n"`
		RepeatPenalty    float64       `json:"repeat_penalty"`
		Seed             int64         `json:"seed"`
		Stop             []interface{} `json:"stop"`
		Stream           bool          `json:"stream"`
		Temp             float64       `json:"temp"`
		TfsZ             float64       `json:"tfs_z"`
		TopK             int           `json:"top_k"`
		TopP             float64       `json:"top_p"`
		TypicalP         float64       `json:"typical_p"`
	} `json:"generation_settings"`
	Model        string `json:"model"`
	Prompt       string `json:"prompt"`
	Stop         bool   `json:"stop"`
	StoppedEos   bool   `json:"stopped_eos"`
	StoppedLimit bool   `json:"stopped_limit"`
	StoppedWord  bool   `json:"stopped_word"`
	StoppingWord string `json:"stopping_word"`
	Timings      struct {
		PredictedMs         float64     `json:"predicted_ms"`
		PredictedN          int         `json:"predicted_n"`
		PredictedPerSecond  float64     `json:"predicted_per_second"`
		PredictedPerTokenMs float64     `json:"predicted_per_token_ms"`
		PromptMs            float64     `json:"prompt_ms"`
		PromptN             int         `json:"prompt_n"`
		PromptPerSecond     interface{} `json:"prompt_per_second"`
		PromptPerTokenMs    float64     `json:"prompt_per_token_ms"`
	} `json:"timings"`
	TokensCached    int  `json:"tokens_cached"`
	TokensEvaluated int  `json:"tokens_evaluated"`
	TokensPredicted int  `json:"tokens_predicted"`
	Truncated       bool `json:"truncated"`
}

const TestResponse = `{
    "name_report": "Протокол совещания",
    "document_type": "docx",
    "password": "pdf",
    "data": {
      "date": "2024-09-07T21:10:52.564Z",
      "time": "21:10:52.564Z",
      "duration": "P3D",
      "participants": [
        "SPEAKER_00 (ведущий)",
        "SPEAKER_01 (Яна)",
        "SPEAKER_02 (Аня)"
      ],
      "agenda": [
        "1. Обсуждение итогов прошедшего дня и ощущений от совещания.",
        "2. Оценка текущих задач и распределение обязанностей на завтра."
      ],
      "blocks": [
        {
          "name_block": "Задачи",
          "proposals": [
            {
              "text": "1. Завершить работу над кейсом для Hackathon до конца сегодняшнего дня (SPEAKER_01).",
              "context": "- Срок: До вечера.",
              "audio_time": {
                "start": "21:10:52.564Z",
                "end": "21:10:52.564Z"
              }
            },
            {
              "text": "2. Подготовить презентацию по речи, используя шаблон от SPEAKER_00 (SPEAKER_02).",
              "context": "- Срок: Завтра.",
              "audio_time": {
                "start": "21:10:52.564Z",
                "end": "21:10:52.564Z"
              }
            },
            {
              "text": "3. Окончательная доработка и проверка транскрипции совещания для публикации.",
              "context": "- Ответственный: SPEAKER_01\\n- Срок: До завтрашнего дня.",
              "audio_time": {
                "start": "21:10:52.564Z",
                "end": "21:10:52.564Z"
              }
            },
            {
              "text": "4. Создание таблицы по перспективам ИИ (на основе экселептической таблицы).",
              "context": "   - Ответственные: SPEAKER_00 и SPEAKER_02\\n- Срок: Завтра, утро.",
              "audio_time": {
                "start": "21:10:52.564Z",
                "end": "21:10:52.564Z"
              }
            },
            {
              "text": "5. Подготовка к завтрашним встречам:",
              "context": "- Продолжение работы над проектом для Хакатона (SPEAKER_01).\\n- Анализ информации и подготовка данных (SPEAKER_02).",
              "audio_time": {
                "start": "21:10:52.564Z",
                "end": "21:10:52.564Z"
              }
            }
          ]
        },
        {
          "name_block": "Дополнительные заметки",
          "proposals": [
            {
              "text": "- SPEAKER_00 отметил важность задачи по созданию комплексного набора материалов, включающего аудиозапись, расшифровку и протокол встречи.",
              "context": "",
              "audio_time": {
                "start": "21:10:52.564Z",
                "end": "21:10:52.564Z"
              }
            },
            {
              "text": "- SPEAKER_01 поделилась ощущениями от завершения кейса на фотон и упомянула необходимость дополнительной работы для полноценного завершения задачи.",
              "context": "",
              "audio_time": {
                "start": "21:10:52.564Z",
                "end": "21:10:52.564Z"
              }
            }
          ]
        },
        {
          "name_block": "Обратная связь",
          "proposals": [
            {
              "text": "- SPEAKER_02 признала, что в работе над проектами чувствуется недостаток пространства для самовыражения сотрудников. Важность наладить более тесный контакт и взаимодействие была подчеркнута.",
              "context": "",
              "audio_time": {
                "start": "21:10:52.564Z",
                "end": "21:10:52.564Z"
              }
            }
          ]
        }
      ],
      "audio_times": [
        {
          "start": "21:10:52.564Z",
          "end": "21:15:52.564Z"
        }
      ]
    }
  }`

type CompletionReq struct {
	Temperature      float64  `json:"temperature,omitempty"`
	TopK             int      `json:"top_k,omitempty"`
	TopP             float64  `json:"top_p,omitempty"`
	NPredict         int      `json:"n_predict,omitempty"`
	NKeep            int      `json:"n_keep,omitempty"`
	Stream           bool     `json:"stream,omitempty"`
	Prompt           string   `json:"prompt,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	TfsZ             float64  `json:"tfs_z,omitempty"`
	TypicalP         float64  `json:"typical_p,omitempty"`
	RepeatPenalty    float64  `json:"repeat_penalty,omitempty"`
	RepeatLastN      int      `json:"repeat_last_n,omitempty"`
	PenalizeNl       bool     `json:"penalize_nl,omitempty"`
	PrecensePenalty  float64  `json:"precence_penalty,omitempty"`
	FrequencyPenalty float64  `json:"frequency_penalty,omitempty"`
	Mirostat         int      `json:"mirostat,omitempty"`
	MirostatTAU      float64  `json:"mirostat_tau,omitempty"`
	MirostatETA      float64  `json:"mirostat_eta,omitempty"`
	Seed             int      `json:"seed,omitempty"`
	IgnoreEOS        bool     `json:"ignore_eos,omitempty"`
}

func (bw *BotWrapper) llamaComplete(ctx context.Context, text string, pgID, chatID, messageID int64) {
	bw.updateStatus(ctx, StatusNers, pgID, chatID, messageID)

	body, err := bw.llamaRequest(context.Background(), CompletionReq{
		Prompt: llamaSystemPrompt + " - " + text,
	}, "/completion")
	if err != nil {
		bw.log.Error().Err(err).Msg("failed to make completeion req")

		return
	}

	fmt.Println(string(body))

	if err := bw.psql.UpdateLlamaOutput(ctx, postgres.UpdateLlamaOutputParams{
		LlamaOutput: pgtype.Text{
			String: string(body),
			Valid:  true,
		},
		ID: pgID,
	}); err != nil {
		bw.log.Error().Err(err).Msg("failed to update llama output")

		return
	}
}

func (bw *BotWrapper) llamaRequest(ctx context.Context, req interface{}, handler string) ([]byte, error) {
	bt, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(bt)

	r, err := http.NewRequestWithContext(ctx, "POST", bw.cfg.LlamaAddr+handler, b)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resp.Body.Close()

	return body, nil
}

func (bw *BotWrapper) officialReport(ctx context.Context, pgID, chatID int64, reportType string) ([]byte, error) {
	tr, err := bw.psql.GetTranscribition(ctx, pgID)
	if err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("failed to get transcribition")

		return nil, err
	}

	var compResp CompletionResp

	if err := json.Unmarshal([]byte(tr.LlamaOutput.String), &compResp); err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("unmarshal json response failed")

		return nil, err
	}

	var reportedReq ReportedRequest
	if err := json.Unmarshal([]byte(compResp.Content), &reportedReq); err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("unmarshal json response failed")

		return nil, err
	}

	bytesResp, err := json.Marshal(reportedReq)
	if err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("unmarshal json marshal failed")

		return nil, err
	}

	buf := bytes.NewBuffer(bytesResp)

	llama, _ := url.Parse(bw.cfg.ReporterAddr + "/reports/official")
	newReq, err := http.NewRequest(http.MethodPost, llama.String(), buf)
	if err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("get report failed")

		return nil, err
	}

	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("request failed")

		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("read body failed")

		return nil, err
	}
	resp.Body.Close()

	return b, nil
}

func (bw *BotWrapper) unofficialReport(ctx context.Context, pgID, chatID int64, reportType string) ([]byte, error) {
	tr, err := bw.psql.GetTranscribition(ctx, pgID)
	if err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("failed to get transcribition")

		return nil, err
	}

	var llamaResp CompletionResp

	if err := json.Unmarshal([]byte(tr.LlamaOutput.String), &llamaResp); err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("unmarshal json response failed")

		return nil, err
	}

	bytesResp, err := json.Marshal(llamaResp)
	if err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("unmarshal json marshal failed")

		return nil, err
	}

	buf := bytes.NewBuffer(bytesResp)

	llama, _ := url.Parse(bw.cfg.ReporterAddr + "/reports/unofficial")
	newReq, err := http.NewRequest(http.MethodPost, llama.String(), buf)
	if err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("get report failed")

		return nil, err
	}

	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("request failed")

		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		bw.log.Error().Err(err).Int64("chatID", chatID).Msg("read body failed")

		return nil, err
	}
	resp.Body.Close()

	return b, nil
}
