-- +goose Up
CREATE TABLE transcribitions (
  id                    BIGSERIAL PRIMARY KEY,
  tg_user_id            BIGINT NOT NULL,
  audio_name_minio      TEXT,
  audio_bucket_minio    TEXT,
  formal_report_minio   TEXT,
  informal_report_minio TEXT,
  transcription         TEXT,
  status                INT,
  created_at            timestamp default current_timestamp,
  llama_output          TEXT,
  message_to_edit       BIGINT
);

-- +goose Down
DROP TABLE transcribitions;
