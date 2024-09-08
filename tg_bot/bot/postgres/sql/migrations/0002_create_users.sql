-- +goose Up
CREATE TABLE users (
  tg_user_id BIGINT PRIMARY KEY,
  current_bot_status TEXT,
  current_bot_id BIGSERIAL
);

-- +goose Down
DROP TABLE users;
