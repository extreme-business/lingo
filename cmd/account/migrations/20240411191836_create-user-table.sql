CREATE TABLE organizations (
    id UUID PRIMARY KEY,
    legal_name VARCHAR(255) NOT NULL UNIQUE,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations (id)
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    display_name VARCHAR(16) NOT NULL,
    email VARCHAR(128) NOT NULL UNIQUE,
    password text NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations (id)
);

CREATE TABLE organizations_users (
    organization_id UUID NOT NULL,
    user_id UUID NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, organization_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (organization_id) REFERENCES organizations (id)
);

CREATE TABLE accounts_users (
    account_id UUID NOT NULL,
    user_id UUID NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, account_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (account_id) REFERENCES accounts (id)
);

-- Trigger Function to Check User Organization
CREATE OR REPLACE FUNCTION check_user_organization()
RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM organizations_users ou
        JOIN accounts a ON ou.organization_id = a.organization_id
        WHERE ou.user_id = NEW.user_id AND a.id = NEW.account_id
    ) THEN
        RAISE EXCEPTION 'User must belong to the same organization as the account';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to Check User Organization
CREATE TRIGGER check_user_organization
BEFORE INSERT ON accounts_users
FOR EACH ROW
EXECUTE FUNCTION check_user_organization();

