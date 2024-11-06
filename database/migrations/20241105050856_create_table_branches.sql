-- +goose Up
-- +goose StatementBegin
CREATE TABLE branches (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  company_id bigint NOT NULL,
  name varchar(255) NOT NULL,
  address text NOT NULL
);

ALTER TABLE ONLY roles
    ADD CONSTRAINT fk_branches_roles FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS branches;
-- +goose StatementEnd
