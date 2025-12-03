-- name: GetAllDevices :many
SELECT * FROM devices;

-- name: GetAllDevicesByBrand :many
SELECT * FROM devices 
WHERE brand = $1;

-- name: GetAllDevicesByState :many
SELECT * FROM devices 
WHERE state = $1;

-- name: GetDeviceByID :one
SELECT * FROM devices WHERE id = $1;

-- name: CreateDevice :one
INSERT INTO devices (id, name, brand, state, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateDevice :exec
UPDATE devices
SET name = $1,
    brand = $2,
    state = $3
WHERE id = $4;

-- name: DeleteDevice :exec
DELETE FROM devices WHERE id = $1;

