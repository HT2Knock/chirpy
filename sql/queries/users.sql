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

-- name: UpdateUser :exec
UPDATE
    users
SET
    email = $2,
    hashed_password = $3,
    updated_at = $4
WHERE
    id = $1;

-- name: UpdateChirpyRed :exec
UPDATE
    users
SET
    is_chirpy_red = $2
WHERE
    id = $1;
