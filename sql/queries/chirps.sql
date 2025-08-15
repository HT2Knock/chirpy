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
