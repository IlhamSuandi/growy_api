-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  branch_id bigint NOT NULL,
  user_id bigint UNIQUE,
  name varchar(100) NOT NULL,

  CONSTRAINT fk_users_role FOREIGN KEY (user_id) REFERENCES users (
      id
  ) ON DELETE CASCADE
);

CREATE INDEX idx_roles_branch_id ON roles (branch_id);
CREATE INDEX idx_roles_user_id ON roles (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
