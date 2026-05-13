-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS device (
    id BIGSERIAL PRIMARY KEY,
    installation_date TIMESTAMP NOT NULL,
    room VARCHAR(100),
    id_patient BIGINT NOT NULL,
    CONSTRAINT fk_device_patient FOREIGN KEY (id_patient) REFERENCES patient(id) ON DELETE CASCADE
    );
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS device;
-- +goose StatementEnd