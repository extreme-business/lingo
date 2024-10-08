-- Create organizations table
CREATE TABLE organizations (
    id UUID PRIMARY KEY,
    legal_name VARCHAR(255) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT organizations_chk_slug CHECK (slug ~ '^[a-z0-9-]+$')
);

-- Create enum type for user status
CREATE TYPE user_status AS ENUM ('active', 'inactive');

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    display_name VARCHAR(16) NOT NULL,
    email VARCHAR(128) NOT NULL,
    hashed_password TEXT NOT NULL,
    status user_status NOT NULL DEFAULT 'active',
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    delete_time TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations (id)
);

-- Create unique index with a partial condition on status
CREATE UNIQUE INDEX users_unique_active_email ON users (email) WHERE status = 'active' AND delete_time IS NULL;
