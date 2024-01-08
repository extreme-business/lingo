-- Table: account
CREATE TABLE account (
    id UUID PRIMARY KEY,
    email_hash VARCHAR(64) NOT NULL UNIQUE,
    activation_time TIMESTAMP,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: device
CREATE TABLE device (
    id UUID UNIQUE,
    account_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    public_key BYTEA NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, account_id)
);

-- Table: event
CREATE TABLE event (
    id UUID PRIMARY KEY,
    sender_account_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    body BYTEA NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: event_recipient
CREATE TABLE event_recipient (
    event_id UUID NOT NULL REFERENCES event(id) ON DELETE CASCADE,
    account_id UUID NOT NULL,
    device_id UUID NOT NULL,
    decryption_key BYTEA NOT NULL,
    FOREIGN KEY (account_id, device_id) REFERENCES device(account_id, id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, account_id, device_id)
);