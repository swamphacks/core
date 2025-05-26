package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() *Config {
	loadEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal().Msg("DATABASE_URL is required in env")
	}

	return &Config{
		DatabaseURL: dbURL,
		Port:        port,
	}
}

func loadEnv() {
	files := []string{".env.local", ".env.development", ".env"}
	for _, f := range files {
		if _, err := os.Stat(f); err == nil {
			log.Info().Str("file", f).Msg("Loading environment file.")
			if err := godotenv.Load(f); err != nil {
				log.Fatal().Err(err).Msgf("Failed to load %s", f)
			}
			break
		}
	}
}
