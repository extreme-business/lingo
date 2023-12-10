-- Table: users
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email_hash VARCHAR(64) NOT NULL UNIQUE,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: users_devices
CREATE TABLE users_devices (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id UUID NOT NULL,
    name VARCHAR(64) NOT NULL,
    PRIMARY KEY (user_id, device_id)
);

-- Table: events
CREATE TABLE events (
    id UUID PRIMARY KEY,
    sender_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body BYTEA NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: event_recipients
CREATE TABLE events_recipients (
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    user_device_id UUID NOT NULL,
    decryption_key VARCHAR(64) NOT NULL,
    PRIMARY KEY (event_id, user_id, user_device_id),
    FOREIGN KEY (user_id, user_device_id) REFERENCES users_devices(user_id, device_id) ON DELETE CASCADE
);