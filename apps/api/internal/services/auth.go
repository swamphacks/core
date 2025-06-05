package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/oauth"
)

var (
	ErrProviderUnsupported  = errors.New("this provider is unsupported")
	ErrAuthenticationFailed = errors.New("failed to authenticate user")
	ErrFetchUserFailed      = errors.New("failed to fetch user info")
)

type AuthService struct {
	userRepo *repository.UserRepository
	client   *http.Client
	logger   zerolog.Logger
	authCfg  *config.AuthConfig
}

func NewAuthService(userRepo *repository.UserRepository, client *http.Client, logger zerolog.Logger, authCfg *config.AuthConfig) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		client:   client,
		logger:   logger.With().Str("service", "AuthService").Str("component", "auth").Logger(),
		authCfg:  authCfg,
	}
}

func (s *AuthService) AuthenticateWithOAuth(ctx context.Context, code, provider string) (*sqlc.AuthSession, error) {
	// Authenticate logic here
	switch provider {
	case "discord":
		discordOAuthResp, err := oauth.ExchangeDiscordCode(ctx, s.client, &s.authCfg.Discord, code)
		if err != nil {
			// Log it
			s.logger.Err(err).Msg("Failed to exchange discord code for user authentication")
			return nil, ErrAuthenticationFailed
		}

		discordUser, err := oauth.GetDiscordUserInfo(ctx, s.client, discordOAuthResp.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("%w: provider=%s", ErrFetchUserFailed, provider)
		}

		// Start transactions to create user and session and account

		break
	default:
		return nil, ErrProviderUnsupported
	}

	return nil, nil
}

func (s *AuthService) GetMe(ctx context.Context) (*sqlc.AuthUser, error) {
	// Pull the userId from the context
	return nil, nil
}
