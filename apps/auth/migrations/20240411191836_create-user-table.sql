-- Table: account
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(64) NOT NULL UNIQUE,
    password_hash VARCHAR(64) NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);