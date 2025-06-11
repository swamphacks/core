package config

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type CookieConfig struct {
	Domain string `env:"DOMAIN"`
	Secure bool   `env:"SECURE"`
}

type OAuthConfig struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURI  string `env:"REDIRECT_URI"`
}

type AuthConfig struct {
	Discord OAuthConfig `envPrefix:"DISCORD_"`
	// Feel free to add more as implementations grow
}

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
	Port        string `env:"PORT" envDefault:"8080"`

	Auth      AuthConfig   `envPrefix:"AUTH_"`
	Cookie    CookieConfig `envPrefix:"COOKIE_"`
	ClientUrl string       `env:"CLIENT_URL"`
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
