CREATE TABLE IF NOT EXISTS records (
    id SERIAL PRIMARY KEY,
    data TEXT NOT NULL,
    deleted BOOLEAN,
    updated_at TIMESTAMP NOT NULL,
    user_id BIGINT NOT NULL
);