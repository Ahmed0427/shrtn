-- +goose Up
CREATE TABLE urls (
    id TEXT PRIMARY KEY,
    original_url TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_accessed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    access_count INT NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE urls;

