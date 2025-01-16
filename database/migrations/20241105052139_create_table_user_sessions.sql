-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_sessions (
  session_id bigint NOT NULL,
  user_id bigint NOT NULL,

  CONSTRAINT fk_user_sessions_session FOREIGN KEY (
      session_id
  ) REFERENCES sessions (id) ON DELETE CASCADE,
  CONSTRAINT fk_user_sessions_user FOREIGN KEY (
      user_id
  ) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_sessions;
-- +goose StatementEnd
