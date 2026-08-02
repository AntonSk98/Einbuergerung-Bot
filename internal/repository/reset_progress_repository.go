package repository

import (
	"github.com/jmoiron/sqlx"
)

// ResetProgressRepository manages database queries related to resetting a user's data and progress.
type ResetProgressRepository struct {
	db *sqlx.DB
}

// NewResetProgressRepository creates a new instance of ResetProgressRepository with the given database connection.
func NewResetProgressRepository(db *sqlx.DB) *ResetProgressRepository {
	return &ResetProgressRepository{db: db}
}

// ResetProgress completely wipes out a user's learning progress and user account data within a single transaction.
func (repo *ResetProgressRepository) ResetProgress(userId int64) error {
	// Start database transaction
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	// Ensure rollback if anything fails, commit if successful
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Delete all learning progress for the user
	_, err = tx.Exec("DELETE FROM user_progress WHERE user_id = ?", userId)
	if err != nil {
		return err
	}

	// Delete the user record itself from the users table
	_, err = tx.Exec("DELETE FROM users WHERE telegram_id = ?", userId)
	if err != nil {
		return err
	}

	return nil
}
