CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    long_link TEXT NOT NULL,
    short_link TEXT NOT NULL,
    expiration_time TIMESTAMP NOT NULL
);