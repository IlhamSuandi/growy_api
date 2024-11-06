-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  username varchar(100),
  email varchar(50) UNIQUE NOT NULL,
  password varchar(256) NOT NULL,
  is_email_verified boolean NOT NULL DEFAULT false,
  auth_provider varchar(50)
);

CREATE INDEX idx_users_email ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
