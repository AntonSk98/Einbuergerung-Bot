package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	TelegramToken     string
	DatabasePath      string
	AuthorizedUserIds []int64
}

// LoadConfig loads the configuration from environment variables and an optional .env file.
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Reading direcetly from environment variables")
	}

	cfg := &Config{
		TelegramToken:     os.Getenv("TELEGRAM_TOKEN"),
		DatabasePath:      os.Getenv("DATABASE_PATH"),
		AuthorizedUserIds: getAuthorizedIdentifiers(),
	}

	return cfg, nil
}

// getAuthorizedIdentifiers parses a comma-separated list of authorized user IDs from the environment.
func getAuthorizedIdentifiers() []int64 {
	authorizedUserIdentifiersString := strings.TrimSpace(os.Getenv("AUTHORIZED_USER_IDS"))
	if authorizedUserIdentifiersString == "" {
		return []int64{}
	}

	var authorizedUserIdentifiers []int64
	for _, identifier := range strings.Split(authorizedUserIdentifiersString, ",") {
		trimmed := strings.TrimSpace(identifier)
		if id, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
			authorizedUserIdentifiers = append(authorizedUserIdentifiers, id)
		} else {
			log.Printf("Warning: failed to parse authorized user ID '%s': %v", trimmed, err)
		}
	}

	return authorizedUserIdentifiers
}
