package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
}

type DiscordExchangeResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToekn string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    time.Duration `json:"expires_in"`
	Scope        string        `json:"scope"`
}

func ExchangeDiscordCode(client *http.Client, code string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "autorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:8080/auth/callback")
	data.Set("client_id", "ADD FROM ENV")
	data.Set("client_secret", "ADD FROM ENV")

	req, err := http.NewRequest("POST", "https://discord.com/api/v10/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		//TODO: HANDLE THIS!
	}

	var exchangeResp DiscordExchangeResponse

	if err := json.NewDecoder(resp.Body).Decode(&exchangeResp); err != nil {
		return "", err
	}

	return exchangeResp.AccessToken, nil

}

func GetDiscordUserInfo(client *http.Client, accessToken string) (*DiscordUser, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status: %s - %s", resp.Status, string(bodyBytes))
	}

	var user DiscordUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
