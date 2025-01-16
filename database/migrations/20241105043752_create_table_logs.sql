-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS logs (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  email varchar(50) NOT NULL,
  remote_addr varchar(50),
  action varchar(10),
  method varchar(10),
  path varchar(50),
  status int NOT NULL,
  execution_time varchar(50),
  size int,
  user_agent text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS logs;
-- +goose StatementEnd
