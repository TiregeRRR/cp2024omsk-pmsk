-- name: CreateTranscribition :one
INSERT INTO transcribitions (
  tg_user_id,
  message_to_edit
) VALUES (
  $1, $2
)
RETURNING id;

-- name: GetUser :one
SELECT * FROM users
WHERE tg_user_id = $1 LIMIT 1;

-- name: GetTranscribition :one
SELECT * FROM transcribitions
WHERE id = $1 LIMIT 1;

-- name: GetTranscribitions :many
SELECT * FROM transcribitions;

-- name: UpdateMinioLink :exec
UPDATE transcribitions
SET audio_name_minio = $1,
    audio_bucket_minio = $2
WHERE id = $3;

-- name: UpdateTranscription :exec
UPDATE transcribitions
SET transcription = $1
WHERE id = $2;

-- name: UpdateStatus :exec
UPDATE transcribitions
SET status = $1
WHERE id = $2;

-- name: UpdateLlamaOutput :exec
UPDATE transcribitions
SET llama_output = $1
WHERE id = $2;

-- name: CreateUser :exec
INSERT INTO users (
  tg_user_id 
) VALUES (
  $1
)
ON CONFLICT(tg_user_id) 
DO NOTHING;

-- name: UpdateCurrentBotID :exec
UPDATE users
SET current_bot_id = $1
WHERE tg_user_id = $2;

-- name: UpdateCurrentBotStatus :exec
UPDATE users
SET current_bot_status = $1
WHERE tg_user_id = $2;
