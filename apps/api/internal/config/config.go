package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectUI   string
}

type AuthConfig struct {
	Discord OAuthConfig
	// Feel free to add more as implementations grow
}

type Config struct {
	DatabaseURL string
	Port        string

	Auth AuthConfig
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

	discordClientID := os.Getenv("DISCORD_CLIENT_ID")
	discordClientSecret := os.Getenv("DISCORD_ClIENT_SECRET")
	discordRedirectURI := os.Getenv("DISCORD_REDIRECT_URI")

	return &Config{
		DatabaseURL: dbURL,
		Port:        port,
		Auth: AuthConfig{
			Discord: OAuthConfig{
				ClientID:     discordClientID,
				ClientSecret: discordClientSecret,
				RedirectUI:   discordRedirectURI,
			},
		},
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
