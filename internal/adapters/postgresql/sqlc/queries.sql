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

-- name: CreateTypeAlert :one
INSERT INTO type_alert (
    id, description
) VALUES ($1, $2) RETURNING *;

-- name: GetTypeAlertByID :one
SELECT * FROM type_alert
WHERE id = $1;



-- name: ListPatients :many
SELECT * FROM patient
ORDER BY name ASC ;

-- name: CreatePatient :one
INSERT INTO patient (
    name,birth_date, gender, address,emergency_contact)
    VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetPatientByID :one
SELECT  * FROM patient
WHERE id = $1;



-- name: GetDeviceByID :one
SELECT * FROM device
WHERE id = $1;

-- name: ListDevices :many
SELECT
    d.id,
    d.installation_date,
    d.room,
    d.id_patient,
    p.name AS patient_name
FROM device d
            JOIN patient p ON d.id_patient = p.id
ORDER BY installation_date ASC;

-- name: CreateDevice :one
INSERT INTO device (
    installation_date, room, id_patient)
    VALUES ($1, $2, $3) RETURNING *;


-- name: ListUsers :many
SELECT id, name, email, phone_number, role, last_access
FROM users
ORDER BY name ASC;

-- name: CreateUser :one
INSERT INTO users (
    name, email, phone_number, role, last_access)
