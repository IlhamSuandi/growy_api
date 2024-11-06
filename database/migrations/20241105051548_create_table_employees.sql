-- +goose Up
-- +goose StatementBegin
CREATE TABLE employees (
  branch_id bigint NOT NULL,
  user_id bigint UNIQUE NOT NULL,

  CONSTRAINT fk_employees_branch FOREIGN KEY (
      branch_id
  ) REFERENCES branches (id) ON DELETE CASCADE,
  CONSTRAINT fk_employees_user FOREIGN KEY (user_id) REFERENCES users (
      id
  ) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employees;
-- +goose StatementEnd
