-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS work_schedules (
    id SERIAL PRIMARY KEY, -- Assuming Model includes an `ID` field
    created_at TIMESTAMP NOT NULL DEFAULT now(), -- Assuming Model includes `CreatedAt`
    updated_at TIMESTAMP NOT NULL DEFAULT now(), -- Assuming Model includes `UpdatedAt`
    deleted_at TIMESTAMP, -- Assuming Model includes `DeletedAt` for soft deletes
    working_day VARCHAR(256) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,

    company_id INTEGER UNIQUE NOT NULL,
    branch_id INTEGER UNIQUE NOT NULL,
    employee_id INTEGER UNIQUE NOT NULL,

    CONSTRAINT fk_company FOREIGN KEY (company_id) 
        REFERENCES companies (id) 
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS work_schedules;
-- +goose StatementEnd
