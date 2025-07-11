package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/oauth"
)

var (
	ErrProviderUnsupported       = errors.New("this provider is unsupported")
	ErrAuthenticationFailed      = errors.New("failed to authenticate user")
	ErrFetchUserFailed           = errors.New("failed to fetch user info")
	ErrFetchSessionContextFailed = errors.New("failed to fetch session context")
	ErrInvalidateSessionFailed   = errors.New("failed to invalidate the session")
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

func (s *AuthService) AuthenticateWithOAuth(ctx context.Context, code, provider string, ipAddress, userAgent *string) (*sqlc.AuthSession, error) {
	switch provider {
	case "discord":
		return s.authenticateWithDiscord(ctx, code, ipAddress, userAgent)
	default:
		return nil, ErrProviderUnsupported
	}
}

func (s *AuthService) GetMe(ctx context.Context) (*middleware.UserContext, error) {
	userContext, ok := ctx.Value(middleware.UserContextKey).(*middleware.UserContext)
	if !ok || userContext == nil {
		return nil, ErrFetchUserFailed
	}

	return userContext, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	sessionContext, ok := ctx.Value(middleware.SessionContextKey).(*middleware.SessionContext)
	if !ok || sessionContext == nil {
		return ErrFetchSessionContextFailed
	}

	err := s.sessionRepo.Invalidate(ctx, sessionContext.SessionID)
	if err != nil {
		return ErrInvalidateSessionFailed
	}

	return nil
}

func (s *AuthService) authenticateWithDiscord(ctx context.Context, code string, ipAddress, userAgent *string) (*sqlc.AuthSession, error) {
	discordOAuthResp, err := oauth.ExchangeDiscordCode(ctx, s.client, &s.authCfg.Discord, code)
	if err != nil {
		// Log it
		s.logger.Err(err).Msg("Failed to exchange discord code for user authentication")
		return nil, ErrAuthenticationFailed
	}

	discordUser, err := oauth.GetDiscordUserInfo(ctx, s.client, discordOAuthResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("%w: provider=discord", ErrFetchUserFailed)
	}

	// Check if account already exists!
	account, err := s.accountRepo.GetByProviderAndAccountID(ctx, sqlc.GetByProviderAndAccountIDParams{
		ProviderID: "discord",
		AccountID:  discordUser.ID,
	})

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return s.registerNewDiscordUser(ctx, discordUser, discordOAuthResp, ipAddress, userAgent)
	} else if err != nil {
		return nil, err
	}

	return s.createSessionForExistingUser(ctx, account.UserID, ipAddress, userAgent)
}

func (s *AuthService) registerNewDiscordUser(ctx context.Context, userInfo *oauth.DiscordUserWithAvatarURL, oauthResp *oauth.DiscordExchangeResponse, ipAddress, userAgent *string) (*sqlc.AuthSession, error) {
	var session *sqlc.AuthSession

	err := s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		txUserRepo := s.userRepo.NewTx(tx)
		txAccountRepo := s.accountRepo.NewTx(tx)
		txSessionRepo := s.sessionRepo.NewTx(tx)

		// Default avatar if no discord avatar
		avatar := userInfo.AvatarURL
		if avatar == nil {
			custom := fmt.Sprintf("https://api.dicebear.com/9.x/initials/png?seed=%s", url.QueryEscape(userInfo.Username))
			avatar = &custom
		}

		user, err := txUserRepo.Create(ctx, sqlc.CreateUserParams{
			Name:  userInfo.Username,
			Email: &userInfo.Email,
			Image: avatar,
		})
		if err != nil {
			return err
		}

		// Create account
		_, err = txAccountRepo.Create(ctx, sqlc.CreateAccountParams{
			UserID:               user.ID,
			ProviderID:           "discord",
			AccountID:            userInfo.ID,
			AccessToken:          &oauthResp.AccessToken,
			RefreshToken:         &oauthResp.RefreshToken,
			AccessTokenExpiresAt: expiresAt(time.Duration(oauthResp.ExpiresIn) * time.Second), // Must convert from seconds to time.Duration
		})
		if err != nil {
			return err
		}
		// create session
		session, err = txSessionRepo.Create(ctx, sqlc.CreateSessionParams{
			UserID:    user.ID,
			ExpiresAt: time.Now().AddDate(0, 1, 0),
			IpAddress: ipAddress,
			UserAgent: userAgent,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) createSessionForExistingUser(ctx context.Context, userID uuid.UUID, ipAddress, userAgent *string) (*sqlc.AuthSession, error) {
	return s.sessionRepo.Create(ctx, sqlc.CreateSessionParams{
		UserID:    userID,
		ExpiresAt: time.Now().AddDate(0, 1, 0),
		IpAddress: ipAddress,
		UserAgent: userAgent,
	})
}

/* HELPER FUNCTIONS BELOW THIS LINE */
func expiresAt(duration time.Duration) *time.Time {
	expiredAtTime := time.Now().Add(duration)
	return &expiredAtTime
}
