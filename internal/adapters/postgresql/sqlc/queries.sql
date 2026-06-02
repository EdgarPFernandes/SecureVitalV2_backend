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
    name,
    email,
    phone_number,
    password_hash,
    role,
    photo
)
VALUES (
           $1,
           $2,
           $3,
           $4,
           $5,
           $6
       )
    RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateLastAccess :exec
UPDATE users
SET last_access = NOW()
WHERE id = $1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: ListPatientsByUser :many
SELECT p.*
FROM patient p
JOIN user_patient up ON up.id_patient = p.id
WHERE up.id_user = $1;

-- name: CountAlertsByPatient :one
SELECT COUNT(*) 
FROM alert a
JOIN device d ON d.id = a.idDevice
WHERE d.id_patient = $1;

-- name: CountAlerts24h :one
SELECT COUNT(*) 
FROM alert a
JOIN device d ON d.id = a.idDevice
WHERE d.id_patient = $1
AND a.date >= NOW() - INTERVAL '24 hours';

-- name: CountAlerts7d :one
SELECT COUNT(*) 
FROM alert a
JOIN device d ON d.id = a.idDevice
WHERE d.id_patient = $1
AND a.date >= NOW() - INTERVAL '7 days';

-- name: ListAlertTypesByPatient :many
SELECT ta.description, COUNT(*) as count
FROM alert a
JOIN type_alert ta ON ta.id = a.idTypeAlert
JOIN device d ON d.id = a.idDevice
WHERE d.id_patient = $1
GROUP BY ta.description;

-- name: CountAlertsByMonth :many
SELECT 
  EXTRACT(MONTH FROM a.date)::int AS month,
  COUNT(*)::bigint AS total
FROM alert a
JOIN device d ON d.id = a.idDevice
WHERE d.id_patient = $1
GROUP BY month
ORDER BY month;

-- name: CountAlertsByYear :many
SELECT 
  EXTRACT(YEAR FROM a.date)::int AS year,
  COUNT(*)::bigint AS total
FROM alert a
JOIN device d ON d.id = a.idDevice
WHERE d.id_patient = $1
GROUP BY year
ORDER BY year;

-- name: ListLastAlertsByPatient :many
SELECT
    a.id,
    a.date,
    ta.description AS type
FROM alert a
JOIN device d ON d.id = a.idDevice
JOIN type_alert ta ON ta.id = a.idTypeAlert
WHERE d.id_patient = $1
ORDER BY a.date DESC
LIMIT 5;

-- name: ListAllPatients :many
SELECT * FROM patient;

-- name: ListAllAlertTypes :many
SELECT ta.description, COUNT(a.id) as count
FROM type_alert ta
LEFT JOIN alert a ON a.idtypealert = ta.id
GROUP BY ta.description;

-- name: CountPatients :one
SELECT COUNT(*) FROM patient;

-- name: CountAlerts24hGlobal :one
SELECT COUNT(*) FROM alert
WHERE date >= NOW() - INTERVAL '24 hours';

-- name: CountAlerts7dGlobal :one
SELECT COUNT(*) FROM alert
WHERE date >= NOW() - INTERVAL '7 days';

-- name: CountAlertsLast7Days :many
SELECT
  DATE(date) as day,
  COUNT(*) as total
FROM alert
WHERE date >= NOW() - INTERVAL '7 days'
GROUP BY day
ORDER BY day;

-- name: CountAlertsByTypeGlobal :many
SELECT
  ta.description,
  COUNT(*) as total
FROM alert a
JOIN type_alert ta ON ta.id = a.idTypeAlert
GROUP BY ta.description;

-- name: TopPatientsByAlerts24h :many
SELECT
  p.id,
  p.name,
  COUNT(a.id) as alerts24h
FROM patient p
JOIN device d ON d.id_patient = p.id
JOIN alert a ON a.idDevice = d.id
WHERE a.date >= NOW() - INTERVAL '24 hours'
GROUP BY p.id
ORDER BY alerts24h DESC;

-- name: TopPatientsByAlertsTotal :many
SELECT
  p.id,
  p.name,
  COUNT(a.id) as totalAlerts
FROM patient p
JOIN device d ON d.id_patient = p.id
JOIN alert a ON a.idDevice = d.id
GROUP BY p.id
ORDER BY totalAlerts DESC;

-- name: CountManyAlertsByPatient :many
SELECT
  p.id,
  COUNT(a.id) as total
FROM patient p
JOIN device d ON d.id_patient = p.id
JOIN alert a ON a.idDevice = d.id
GROUP BY p.id;

-- name: CountAlertsLast12Months :many
SELECT
  EXTRACT(MONTH FROM date) AS month,
  COUNT(*) as total
FROM alert
WHERE date >= NOW() - INTERVAL '12 months'
GROUP BY month
ORDER BY month;

-- name: CountAlerts7To14Days :many
SELECT
  DATE(date) as day,
  COUNT(*) as total
FROM alert
WHERE date >= NOW() - INTERVAL '14 days'
  AND date < NOW() - INTERVAL '7 days'
GROUP BY day
ORDER BY day;

-- name: TopPatientsByAlerts7d :many
SELECT
  p.id,
  p.name,
  COUNT(a.id) as alerts7d
FROM patient p
JOIN device d ON d.id_patient = p.id
JOIN alert a ON a.idDevice = d.id
WHERE a.date >= NOW() - INTERVAL '7 days'
GROUP BY p.id
ORDER BY alerts7d DESC;

-- name: CountAlertsLast12MonthsGlobal :many
SELECT
  EXTRACT(MONTH FROM date)::int AS month,
  COUNT(*)::bigint AS total
FROM alert
WHERE date >= NOW() - INTERVAL '12 months'
GROUP BY month
ORDER BY month;

-- name: GetPatientDashboardInfo :one
SELECT
    id,
    name,
    birth_date,
    gender,
    address,
    emergency_contact,
    created_at
FROM patient
WHERE id = $1;

-- name: ListCaregiversByPatient :many
SELECT
    u.id,
    u.name
FROM users u
JOIN user_patient up ON up.id_user = u.id
WHERE up.id_patient = $1;