package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db"
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
	userRepo    *repository.UserRepository
	accountRepo *repository.AccountRepository
	sessionRepo *repository.SessionRepository
	txm         *db.TransactionManager
	client      *http.Client
	logger      zerolog.Logger
	authCfg     *config.AuthConfig
}

func NewAuthService(userRepo *repository.UserRepository, accountRepo *repository.AccountRepository, sessionRepo *repository.SessionRepository, txm *db.TransactionManager, client *http.Client, logger zerolog.Logger, authCfg *config.AuthConfig) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		accountRepo: accountRepo,
		sessionRepo: sessionRepo,
		txm:         txm,
		client:      client,
		logger:      logger.With().Str("service", "AuthService").Str("component", "auth").Logger(),
		authCfg:     authCfg,
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

		// Check if account already exists!
		_, err = s.accountRepo.GetByProviderAndAccountID(ctx, sqlc.GetByProviderAndAccountIDParams{
			ProviderID: provider,
			AccountID:  discordUser.ID,
		})

		if err != nil && errors.Is(err, sql.ErrNoRows) {
			s.txm.WithTx(ctx, func(tx pgx.Tx) error {
				txUserRepo := s.userRepo.NewTx(tx)
				txAccountRepo := s.accountRepo.NewTx(tx)
				txSessionRepo := s.sessionRepo.NewTx(tx)

				user, err := txUserRepo.Create(ctx, sqlc.CreateUserParams{
					Name:  discordUser.Username,
					Email: discordUser.Email,
					Image: &discordUser.Avatar,
				})
				if err != nil {
					return err
				}

				// Create account
				// create session
			})
		} else if err != nil {
			return nil, err
		} else {
			// Handle user already has an account, fetch user and create new session!
		}

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
