-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS salaries (
    id SERIAL PRIMARY KEY, -- Assuming Model includes an `ID` field
    created_at TIMESTAMP NOT NULL DEFAULT now(), -- Assuming Model includes `CreatedAt`
    updated_at TIMESTAMP NOT NULL DEFAULT now(), -- Assuming Model includes `UpdatedAt`
    deleted_at TIMESTAMP, -- Assuming Model includes `DeletedAt` for soft deletes

    employee_email VARCHAR(255) NOT NULL,
    monthly INTEGER NOT NULL DEFAULT 0,
    hourly INTEGER NOT NULL DEFAULT 0,
    current_earnings INTEGER NOT NULL DEFAULT 0,
    
    CONSTRAINT fk_employee FOREIGN KEY (employee_email) 
        REFERENCES employees (employee_email) 
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS salaries;
-- +goose StatementEnd
