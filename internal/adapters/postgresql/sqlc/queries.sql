-- name: ListAlerts :many
SELECT * FROM alert
ORDER BY date DESC;

-- name: CreateAlert :one
INSERT INTO alert (
    iddevice, idtypealert
) VALUES ($1, $2) RETURNING *;

-- name: GetDeviceByID :one
SELECT * FROM device
WHERE id = $1;

-- name: GetTypeAlertByID :one
SELECT * FROM type_alert
WHERE id = $1;