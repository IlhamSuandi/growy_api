-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  refresh_token varchar(512) NOT NULL,
  is_revoked boolean NOT NULL DEFAULT false,
  user_agent text NOT NULL,
  ip_address text NOT NULL,
  expires_at timestamp NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
