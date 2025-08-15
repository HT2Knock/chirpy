-- name: CreateUser :one
INSERT INTO
    users (email, created_at, updated_at)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: DeleteUsers :exec
DELETE FROM
    users;

-- name: GetUser :one
SELECT
    id,
    email
FROM
    users
WHERE
    id = $1
LIMIT
    1;
