-- +goose Up
-- +goose StatementBegin
CREATE TABLE branch_options (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  branch_id bigint UNIQUE NOT NULL,
  use_checkout boolean NOT NULL DEFAULT true
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS branch_options;
-- +goose StatementEnd
