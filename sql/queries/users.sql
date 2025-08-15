-- name: CreateUser :one
INSERT INTO
    users (email, hashed_password, created_at, updated_at)
VALUES
    ($1, $2, $3, $4)
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

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1;
