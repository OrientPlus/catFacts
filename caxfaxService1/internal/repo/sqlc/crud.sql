-- query.sql

-- name: GetFactByID :one
SELECT id, message, length
FROM facts
WHERE id = $1;

-- name: GetFactByMessage :one
SELECT id, message, length
FROM facts
WHERE message = $1;


-- name: AddFact :one
INSERT INTO facts (message, length)
VALUES ($1, $2) RETURNING id;

-- name: UpdateFactByID :exec
UPDATE facts
SET message = $1,
    length = $2
WHERE id = $3;

-- name: DeleteFactByID :exec
DELETE FROM facts
WHERE id = $1;


-- name: DeleteFactByMessage :exec
DELETE FROM facts
WHERE message = $1;