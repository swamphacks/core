package config

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type OAuthConfig struct {
	ClientID     string `env:"CLIENTID"`
	ClientSecret string `env:"CLIENTSECRET"`
	RedirectUI   string `env:"REDIRECTUI"`
}

type AuthConfig struct {
	Discord OAuthConfig `envPrefix:"DISCORD_"`
	// Feel free to add more as implementations grow
}

type Config struct {
	DatabaseURL string `env:"DATABASEURL"`
	Port        string `env:"PORT" envDefault:"8080"`

	Auth AuthConfig `envPrefix:"AUTH_"`
}

func Load() *Config {
	loadEnv()

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		log.Fatal().Msgf("Failed to parse env: %v", err)
	}

	return &cfg
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
