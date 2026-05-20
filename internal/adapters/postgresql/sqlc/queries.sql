-- name: ListAlerts :many
SELECT
    a.id,
    a.date,
    a.iddevice,
    a.idtypealert,
    ta.description AS alert_type_description
FROM alert a
         JOIN device d ON a.iddevice = d.id
         JOIN type_alert ta ON a.idtypealert = ta.id
ORDER BY a.date DESC;

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

-- name: ListAlertsByPatient :many
SELECT
    a.id,
    a.date,
    a.iddevice,
    a.idtypealert,
    ta.description AS alert_type_description
FROM alert a
         JOIN device d ON a.iddevice = d.id
         JOIN type_alert ta ON a.idtypealert = ta.id
WHERE d.id_patient = $1
ORDER BY a.date DESC;



-- name: ListTypeAlerts :many
SELECT * FROM type_alert
ORDER BY id DESC;