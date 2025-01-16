-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  username varchar(100),
  email varchar(50) UNIQUE NOT NULL,
  password varchar(256) NOT NULL,
  is_email_verified boolean NOT NULL DEFAULT false,
  role varchar(10) NOT NULL,
  auth_provider varchar(50),
  is_on_boarded boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
