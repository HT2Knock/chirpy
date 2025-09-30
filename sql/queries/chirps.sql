-- name: CreateChirp :one
INSERT INTO
    chirps (body, user_id, created_at, updated_at)
VALUES
    ($1, $2, $3, $4)
RETURNING
    *;

-- name: DeleteChirp :exec
DELETE FROM
    chirps
WHERE
    user_id = $1;

-- name: GetChirps :many
SELECT
    *
FROM
    chirps
ORDER BY
    created_at;

-- name: GetChirpsByAuthor :many
SELECT
    *
FROM
    chirps
WHERE
    user_id = $1
ORDER BY
    created_at;

-- name: GetChirp :one
SELECT
    *
FROM
    chirps
WHERE
    id = $1;

-- name: DeleteChirpByID :exec
DELETE FROM
    chirps
WHERE
    id = $1
    AND user_id = $2;
