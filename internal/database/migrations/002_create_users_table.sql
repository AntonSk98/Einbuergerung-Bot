-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    telegram_id INTEGER PRIMARY KEY,
    selected_federal_state TEXT
);

-- +migrate Down
DROP TABLE IF EXISTS users;