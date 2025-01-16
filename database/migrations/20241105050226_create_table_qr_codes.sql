-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS qr_codes (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  user_id bigint UNIQUE NOT NULL,
  code text NOT NULL,
  is_used boolean NOT NULL DEFAULT false,
  expires_at timestamp NULL,

  CONSTRAINT fk_users_qr_code FOREIGN KEY (user_id) REFERENCES users (
      id
  ) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_qr_codes_user_id ON qr_codes (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS qr_codes;
-- +goose StatementEnd
