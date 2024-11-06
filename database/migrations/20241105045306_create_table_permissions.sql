-- +goose Up
-- +goose StatementBegin
CREATE TABLE permissions (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  role_id bigint NOT NULL,
  resource varchar(255) NOT NULL,
  action varchar(50) NOT NULL
);

CREATE INDEX idx_permissions_role_id ON permissions (role_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS permissions;
-- +goose StatementEnd
