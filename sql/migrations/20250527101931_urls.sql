-- +goose Up
CREATE TABLE urls (
    id TEXT PRIMARY KEY,
    original TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE urls;

