package repository

import (
	"einbuergerung-bot/internal/models"

	"github.com/jmoiron/sqlx"
)

// QuestionRepository manages database queries related to questions.
type QuestionRepository struct {
	db *sqlx.DB
}

// NewQuestionRepository creates a new instance of QuestionRepository with the given database connection.
func NewQuestionRepository(db *sqlx.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// FindQuestionById retrieves a specific question by its unique identifier.
func (repo *QuestionRepository) FindQuestionById(questionId int) (*models.Question, error) {
	var question models.Question
	query := "SELECT id, question_text, option_a, option_b, option_c, option_d, correct_answer FROM questions WHERE id = ?"

	err := repo.db.Get(&question, query, questionId)
	if err != nil {
		return nil, err
	}

	return &question, nil
}

// SelectQuestion selects the next appropriate question for a user based on weight, last answered time, and federal state.
func (repo *QuestionRepository) SelectQuestion(userId int64, federalState string) (*models.Question, error) {
	query := `
        SELECT q.id, q.federal_state, q.question_text, q.image_data, q.option_a, q.option_b, q.option_c, q.option_d, q.correct_answer
        FROM questions q
        LEFT JOIN user_progress up ON q.id = up.question_id AND up.user_id = ?
        WHERE (q.federal_state = ? OR q.federal_state = 'general' OR q.federal_state = '')
          AND (up.last_answered_question_at IS NULL OR datetime(up.last_answered_question_at) < datetime('now', '-5 minutes'))
        ORDER BY 
            COALESCE(up.correct_answer_weight, 0) ASC,
            COALESCE(up.last_answered_question_at, '1970-01-01 00:00:00') ASC,
            RANDOM()
        LIMIT 1;
    `

	var question models.Question
	if err := repo.db.Get(&question, query, userId, federalState); err != nil {
		return nil, err
	}

	return &question, nil
}
