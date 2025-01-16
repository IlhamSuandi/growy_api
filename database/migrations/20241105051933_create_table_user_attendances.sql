-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_attendances (
  attendance_id bigint NOT NULL,
  user_id bigint NOT NULL,

  CONSTRAINT fk_user_attendances_attendance FOREIGN KEY (
      attendance_id
  ) REFERENCES attendances (id) ON DELETE CASCADE,
  CONSTRAINT fk_user_attendances_user FOREIGN KEY (
      user_id
  ) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_attendances;
-- +goose StatementEnd
