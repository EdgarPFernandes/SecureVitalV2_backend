-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_patient (
    id_user BIGINT NOT NULL,
    id_patient BIGINT NOT NULL,
    PRIMARY KEY (id_user, id_patient),
    CONSTRAINT fk_user_patient_user FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_patient_patient FOREIGN KEY (id_patient) REFERENCES patient(id) ON DELETE CASCADE
    );
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_patient;
-- +goose StatementEnd