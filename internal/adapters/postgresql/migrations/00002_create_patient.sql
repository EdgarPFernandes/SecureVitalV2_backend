-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS patient (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    gender CHAR(1) NOT NULL,
    address TEXT,
    emergency_contact VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS patient;
-- +goose StatementEnd