-- Table: organizations
CREATE TABLE organizations (
    id UUID PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: users
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(64) NOT NULL UNIQUE,
    password text NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: user_organizations
CREATE TABLE user_organizations (
    user_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    PRIMARY KEY (user_id, organization_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (organization_id) REFERENCES organizations (id)
);