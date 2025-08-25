package strava

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func NewAuth(clientID, clientSecret string) *auth {
	return &auth{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

type auth struct {
	ClientID     string
	ClientSecret string
}

func (a *auth) GetAuthorizationUrl(redirectURI string) *url.URL {
	baseUrl := "https://www.strava.com/oauth/authorize"
	u, _ := url.Parse(baseUrl)

	params := url.Values{}
	params.Add("client_id", a.ClientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("response_type", "code")
	params.Add("approval_prompt", "auto")
	params.Add("scope", "read,profile:read_all,activity:read_all")

	u.RawQuery = params.Encode()

	return u
}

func (a *auth) ExchangeCodeForAccessToken(code string) (*TokenResponse, error) {
	baseUrl := "https://www.strava.com/oauth/token"
	u, _ := url.Parse(baseUrl)

	params := url.Values{}
	params.Add("client_id", a.ClientID)
	params.Add("client_secret", a.ClientSecret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")

	u.RawQuery = params.Encode()

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange code for access token: %s", resp.Status)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func (a *auth) RefreshToken(refreshToken string) (*TokenResponse, error) {
	baseUrl := "https://www.strava.com/oauth/token"
	u, _ := url.Parse(baseUrl)

	params := url.Values{}
	params.Add("client_id", a.ClientID)
	params.Add("client_secret", a.ClientSecret)
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", refreshToken)

	u.RawQuery = params.Encode()

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to refresh token: %s", resp.Status)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}
