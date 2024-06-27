-- Create enum type for user status
CREATE TYPE user_status AS ENUM ('active', 'deleted', 'inactive');

-- Create organizations table
CREATE TABLE organizations (
    id UUID PRIMARY KEY,
    legal_name VARCHAR(255) NOT NULL UNIQUE,
    Slug VARCHAR(100) NOT NULL UNIQUE,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    CONSTRAINT organizations_chk_slug CHECK (slug ~ '^[a-z0-9-]+$')
);

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
    deleted_time TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations (id)
);

-- Create unique index with a partial condition on status
CREATE UNIQUE INDEX users_unique_active_email ON users (email) WHERE status = 'active';

-- Create organizations_users table
CREATE TABLE organizations_users (
    organization_id UUID NOT NULL,
    user_id UUID NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, organization_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (organization_id) REFERENCES organizations (id)
);