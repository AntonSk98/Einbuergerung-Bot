-- +migrate Up
CREATE TABLE IF NOT EXISTS user_progress (
    user_id INTEGER,
    question_id INTEGER,
    correct_answer_weight INTEGER DEFAULT 0,
    last_answered_question_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, question_id),
    FOREIGN KEY (user_id) REFERENCES users(telegram_id),
    FOREIGN KEY (question_id) REFERENCES questions(id)
);

-- +migrate Down
DROP TABLE IF EXISTS user_progress;