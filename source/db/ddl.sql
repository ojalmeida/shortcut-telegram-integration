CREATE TABLE IF NOT EXISTS users
(

    id CHAR(36) NOT NULL PRIMARY KEY,
    telegram_chat_id BIGINT UNIQUE,
    shortcut_mention_name TEXT UNIQUE NOT NULL,
    notification_rate INT NOT NULL -- in seconds

);


CREATE TABLE IF NOT EXISTS reviews
(
    id              CHAR(36) NOT NULL PRIMARY KEY,
    url             TEXT NOT NULL,
    number          BIGINT NOT NULL,
    requester_id    CHAR(36) NOT NULL REFERENCES users(id),
    requested_id    CHAR(36) NOT NULL REFERENCES users(id),
    suppressed_by_requester BOOL NOT NULL,
    suppressed_by_requested BOOL NOT NULL,
    notify_at       BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS chats
(

    id BIGINT NOT NULL PRIMARY KEY,
    authorized BOOL NOT NULL

)
