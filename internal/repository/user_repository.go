package repository

import (
	"database/sql"
	"einbuergerung-bot/internal/models"

	"github.com/jmoiron/sqlx"
)

// UserRepository manages database queries related to users.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new instance of UserRepository with the given database connection.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FederalStateSelected checks whether a user has already selected a non-empty federal state.
func (repo *UserRepository) FederalStateSelected(userId int64) bool {
	var selectedState sql.NullString
	query := "SELECT selected_federal_state FROM users WHERE telegram_id = ?"

	err := repo.db.Get(&selectedState, query, userId)
	if err != nil || !selectedState.Valid {
		return false
	}

	return selectedState.String != ""
}

// PersistFederalState saves or updates the user's selected federal state in the database.
func (repo *UserRepository) PersistFederalState(userId int64, federalState string) error {
	query := `
        INSERT INTO users (telegram_id, selected_federal_state) 
        VALUES (?, ?)
        ON CONFLICT(telegram_id) DO UPDATE SET 
        selected_federal_state = excluded.selected_federal_state
    `

	_, err := repo.db.Exec(query, userId, federalState)
	return err
}

// FindUserById retrieves a user from the database by their Telegram ID.
func (repo *UserRepository) FindUserById(userId int64) (*models.User, error) {
	query := `
        SELECT telegram_id, selected_federal_state 
        FROM users 
        WHERE telegram_id = ?
    `

	var user models.User
	err := repo.db.Get(&user, query, userId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
