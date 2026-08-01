package repository

import (
	"github.com/jmoiron/sqlx"
)

// UserProgressRepository manages database queries related to user progress and statistics.
type UserProgressRepository struct {
	db *sqlx.DB
}

// NewUserProgressRepository creates a new instance of UserProgressRepository with the given database connection.
func NewUserProgressRepository(db *sqlx.DB) *UserProgressRepository {
	return &UserProgressRepository{db: db}
}

// HandleCorrectAnswer increments or initializes the correct answer weight when a user answers correctly.
func (repo *UserProgressRepository) HandleCorrectAnswer(userId int64, questionId int) error {
	query := `
        INSERT INTO user_progress (user_id, question_id, correct_answer_weight, last_answered_question_at) 
        VALUES (?, ?, 1, CURRENT_TIMESTAMP)
        ON CONFLICT(user_id, question_id) DO UPDATE SET 
        correct_answer_weight = correct_answer_weight + 1,
        last_answered_question_at = CURRENT_TIMESTAMP
    `

	_, err := repo.db.Exec(query, userId, questionId)
	return err
}

// HandleWrongAnswer decrements or initializes the correct answer weight when a user answers incorrectly.
func (repo *UserProgressRepository) HandleWrongAnswer(userId int64, questionId int) error {
	query := `
        INSERT INTO user_progress (user_id, question_id, correct_answer_weight, last_answered_question_at) 
        VALUES (?, ?, -1, CURRENT_TIMESTAMP)
        ON CONFLICT(user_id, question_id) DO UPDATE SET 
        correct_answer_weight = correct_answer_weight - 1,
        last_answered_question_at = CURRENT_TIMESTAMP
    `

	_, err := repo.db.Exec(query, userId, questionId)
	return err
}
