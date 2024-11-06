-- +goose Up
-- +goose StatementBegin
CREATE TABLE attendances (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  check_in timestamp,
  check_out timestamp,
  status varchar(20) NOT NULL,
  date date NOT NULL,
  location varchar(100) DEFAULT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS attendances;
-- +goose StatementEnd
