-- +goose Up
-- +goose StatementBegin
CREATE TABLE companies (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  name varchar(255) NOT NULL,
  address text NOT NULL,
  owner_email varchar(50) UNIQUE NOT NULL,

  CONSTRAINT fk_users_company FOREIGN KEY (
      owner_email
  ) REFERENCES users (email) ON DELETE CASCADE
);

CREATE INDEX idx_companies_owner_email ON companies (owner_email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS companies;
-- +goose StatementEnd
