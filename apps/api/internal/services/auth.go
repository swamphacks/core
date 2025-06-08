package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

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

func (s *AuthService) AuthenticateWithOAuth(ctx context.Context, code, provider, ipAddress, userAgent string) (*sqlc.AuthSession, error) {
	var session *sqlc.AuthSession
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
		account, err := s.accountRepo.GetByProviderAndAccountID(ctx, sqlc.GetByProviderAndAccountIDParams{
			ProviderID: provider,
			AccountID:  discordUser.ID,
		})

		if err != nil && errors.Is(err, sql.ErrNoRows) {
			err = s.txm.WithTx(ctx, func(tx pgx.Tx) error {
				txUserRepo := s.userRepo.NewTx(tx)
				txAccountRepo := s.accountRepo.NewTx(tx)
				txSessionRepo := s.sessionRepo.NewTx(tx)

				user, err := txUserRepo.Create(ctx, sqlc.CreateUserParams{
					Name:  discordUser.Username,
					Email: &discordUser.Email,
					Image: &discordUser.Avatar,
				})
				if err != nil {
					return err
				}

				// Create account
				_, err = txAccountRepo.Create(ctx, sqlc.CreateAccountParams{
					UserID:               user.ID,
					ProviderID:           provider,
					AccountID:            discordUser.ID,
					AccessToken:          &discordOAuthResp.AccessToken,
					RefreshToken:         &discordOAuthResp.RefreshToken,
					AccessTokenExpiresAt: expiresAt(discordOAuthResp.ExpiresIn),
				})
				if err != nil {
					return err
				}
				// create session
				session, err = txSessionRepo.Create(ctx, sqlc.CreateSessionParams{
					UserID:    user.ID,
					ExpiresAt: time.Now().AddDate(0, 1, 0),
					IpAddress: &ipAddress,
					UserAgent: &userAgent,
				})
				if err != nil {
					return err
				}

				return nil
			})
			if err != nil {
				return nil, err
			}

		} else if err != nil {
			return nil, err
		} else {
			// Handle user already has an account, fetch user and create new session!
			user, err := s.userRepo.GetByID(ctx, account.UserID)
			if err != nil {
				return nil, err
			}

			session, err = s.sessionRepo.Create(ctx, sqlc.CreateSessionParams{
				UserID:    user.ID,
				ExpiresAt: time.Now().AddDate(0, 1, 0),
				IpAddress: &ipAddress,
				UserAgent: &userAgent,
			})
			if err != nil {
				return nil, err
			}

		}
	default:
		return nil, ErrProviderUnsupported
	}

	return session, nil
}

func (s *AuthService) GetMe(ctx context.Context) (*sqlc.AuthUser, error) {
	// Pull the userId from the context
	return nil, nil
}

func expiresAt(duration time.Duration) *time.Time {
	expiredAtTime := time.Now().Add(duration)
	return &expiredAtTime
}
