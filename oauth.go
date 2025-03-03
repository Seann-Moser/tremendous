package tremendous

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type OauthConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
}

func (c *Client) SendOauthRequest(ctx context.Context, data *AccessTokenRequest) (*TokenResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access token request: %w", err)
	}
	e := strings.ReplaceAll(c.endpoint, "/api/v2/", "")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, joinURL(e, "/oauth/token"), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed sending request: %s", resp.Status)
	}
	at := &TokenResponse{}
	if err := json.NewDecoder(resp.Body).Decode(at); err != nil {
		return nil, fmt.Errorf("failed to decode access token response: %w", err)
	}
	return at, nil
}
