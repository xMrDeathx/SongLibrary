CREATE TABLE song
(
    id uuid NOT NULL PRIMARY KEY,
    band VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    release_date VARCHAR(255),
    text TEXT,
    link VARCHAR(1000),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);