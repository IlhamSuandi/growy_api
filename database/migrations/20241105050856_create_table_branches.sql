-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS branches (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  company_id bigint NOT NULL,
  name varchar(255) NOT NULL,
  address text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS branches;
-- +goose StatementEnd
