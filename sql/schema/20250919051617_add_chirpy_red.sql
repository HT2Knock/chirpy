-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    chirps
ADD
    COLUMN is_chirpy_red boolean NOT NULL DEFAULT false;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    chirps DROP COLUMN is_chirpy_red;

-- +goose StatementEnd
