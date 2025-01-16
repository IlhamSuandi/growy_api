-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS permissions (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  user_id bigint NOT NULL,
  resource varchar(255) NOT NULL,
  action varchar(50) NOT NULL
);

CREATE INDEX idx_permissions_user_id ON permissions (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS permissions;
-- +goose StatementEnd
