CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username TEXT NOT NULL UNIQUE,
                                     password BYTEA NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                     email TEXT NOT NULL UNIQUE,
                                     is_activated BOOLEAN NOT NULL DEFAULT FALSE,
                                     version INT NOT NULL DEFAULT 1
);