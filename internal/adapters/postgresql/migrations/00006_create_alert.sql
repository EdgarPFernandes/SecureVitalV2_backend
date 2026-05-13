-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS alert (
    id SERIAL PRIMARY KEY,
    date TIMESTAMPTZ DEFAULT NOW(),
    idDevice INT NOT NULL,
    idTypeAlert INT NOT NULL,
    CONSTRAINT fk_alert_device
    FOREIGN KEY (idDevice) REFERENCES device(id) ON DELETE CASCADE,
    CONSTRAINT fk_alert_type_alert
    FOREIGN KEY (idTypeAlert) REFERENCES type_alert(id) ON DELETE CASCADE
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS alert;
-- +goose StatementEnd