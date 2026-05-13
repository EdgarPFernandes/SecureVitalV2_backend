-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS type_alert (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS type_alert;
-- +goose StatementEnd