package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/swamphacks/core/apps/api/internal/config"
)

var (
	ErrDiscordExchangeCode = errors.New("error exchanging the discord code")
	ErrDiscordFetchProfile = errors.New("error fetching user discord profile")
)

type DiscordUser struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Avatar        *string `json:"avatar"`
	Discriminator string  `json:"discriminator"`
	Email         string  `json:"email"`
}

type DiscordUserWithAvatarURL struct {
	DiscordUser
	AvatarURL *string
}

// Note: expiresIn is in seconds
type DiscordExchangeResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func ExchangeDiscordCode(ctx context.Context, client *http.Client, oauthCfg *config.OAuthConfig, code string) (*DiscordExchangeResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", oauthCfg.RedirectURI)
	data.Set("client_id", oauthCfg.ClientID)
	data.Set("client_secret", oauthCfg.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://discord.com/api/v10/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req) // Make the request for the code
	if err != nil {
		return nil, err
	}
	defer gracefullyCloseBody(resp)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrDiscordExchangeCode, string(body))
	}

	var exchangeResp DiscordExchangeResponse

	if err := json.NewDecoder(resp.Body).Decode(&exchangeResp); err != nil {
		return nil, err
	}

	return &exchangeResp, nil

}

func GetDiscordUserInfo(ctx context.Context, client *http.Client, accessToken string) (*DiscordUserWithAvatarURL, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer gracefullyCloseBody(resp)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrDiscordFetchProfile, string(bodyBytes))
	}

	var user DiscordUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	// Parse real avatar URL or make nil
	userWithAvatar := DiscordUserWithAvatarURL{
		DiscordUser: user,
		AvatarURL:   user.AvatarURL(),
	}

	return &userWithAvatar, nil
}

func (u *DiscordUser) AvatarURL() *string {
	// Only proceed if Avatar is non-nil *and* not the empty string
	if u.Avatar != nil && *u.Avatar != "" {
		hash := *u.Avatar

		ext := "png"
		if strings.HasPrefix(hash, "a_") {
			ext = "gif"
		}

		url := fmt.Sprintf(
			"https://cdn.discordapp.com/avatars/%s/%s.%s",
			u.ID, hash, ext,
		)
		return &url
	}

	return nil
}

// TODO: Refactor more cleanly
func gracefullyCloseBody(response *http.Response) {
	if err := response.Body.Close(); err != nil {
		log.Warn().Str("component", "oauth").Str("provider", "discord").Msg("Failed to close response body")
	}
}
