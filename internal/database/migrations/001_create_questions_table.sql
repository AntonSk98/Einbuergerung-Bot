-- +migrate Up
CREATE TABLE IF NOT EXISTS questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    federal_state TEXT,
    question_text TEXT,
    image_data BLOB,
    option_a TEXT,
    option_b TEXT,
    option_c TEXT,
    option_d TEXT,
    correct_answer TEXT
);

-- +migrate Down
DROP TABLE IF EXISTS questions;