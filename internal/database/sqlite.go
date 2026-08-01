package database

import (
	"embed"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"einbuergerung-bot/internal/config"
	"einbuergerung-bot/internal/models"

	"github.com/jmoiron/sqlx"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed assets/questions.json
var assetsFS embed.FS

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Initialize opens the SQLite database connection, executes migrations, and synchronizes questions.
func Initialize(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	if err := synchronizeQuestionsOnchange(db); err != nil {
		log.Fatalf("Failed to load question from JSON: %v", err)
	}

	return db, nil
}

// createTables executes embedded database migrations up to the latest version.
func createTables(db *sqlx.DB) error {
	// Tell the migrator to read from your embedded folder
	assetsMigrationSource := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationFiles,
		Root:       "migrations",
	}

	// Run the migrations up to the latest version
	_, err := migrate.Exec(db.DB, "sqlite3", assetsMigrationSource, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}

// synchronizeQuestionsOnchange checks if embedded JSON questions differ from the database and updates them if needed.
func synchronizeQuestionsOnchange(db *sqlx.DB) error {
	var questionCountInDatabase int
	if err := db.Get(&questionCountInDatabase, "SELECT COUNT(*) FROM questions"); err != nil {
		return err
	}

	questionsFromJson, err := parseQuestionsFromJson()
	if err != nil {
		return err
	}

	// If counts match, no need to update
	if questionCountInDatabase == len(questionsFromJson) {
		return nil
	}

	return synchronizeQuestions(db, questionsFromJson)
}

// synchronizeQuestions clears existing questions and bulk-inserts the latest ones from JSON within a transaction.
func synchronizeQuestions(db *sqlx.DB, questionsFromJson []models.Question) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM questions"); err != nil {
		return err
	}

	query := `
        INSERT INTO questions (federal_state, question_text, image_data, option_a, option_b, option_c, option_d, correct_answer)
        VALUES (:federal_state, :question_text, :image_data, :option_a, :option_b, :option_c, :option_d, :correct_answer)
    `

	for _, q := range questionsFromJson {
		if _, err := tx.NamedExec(query, q); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// parseQuestionsFromJson reads and unmarshals the embedded JSON question dataset.
func parseQuestionsFromJson() ([]models.Question, error) {
	data, err := assetsFS.ReadFile("assets/questions.json")
	if err != nil {
		return nil, err
	}

	var questions []models.Question
	if err := json.Unmarshal(data, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}
