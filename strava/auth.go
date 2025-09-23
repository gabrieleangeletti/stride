package strava

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func NewAuth(clientID, clientSecret string) *Auth {
	return &Auth{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

type Auth struct {
	ClientID     string
	ClientSecret string
}

func (a *Auth) GetAuthorizationUrl(redirectURI string) *url.URL {
	u, _ := url.Parse("https://www.strava.com/oauth/authorize")

	params := url.Values{}
	params.Add("client_id", a.ClientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("response_type", "code")
	params.Add("approval_prompt", "auto")
	params.Add("scope", "read,profile:read_all,activity:read_all")

	u.RawQuery = params.Encode()

	return u
}

func (a *Auth) ExchangeCodeForAccessToken(code string) (*TokenResponse, error) {
	u, _ := url.Parse("https://www.strava.com/oauth/token")

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
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to exchange code for access token: %s\n%s", resp.Status, string(bodyBytes))
	}

	var r TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

func (a *Auth) RefreshToken(refreshToken string) (*TokenResponse, error) {
	u, _ := url.Parse("https://www.strava.com/oauth/token")

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
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to refresh token: %s\n%s", resp.Status, string(bodyBytes))
	}

	var r TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

func (a *Auth) RegisterWebhookSubscription(callbackURL, verifyToken string) (*WebhookRegistrationResponse, error) {
	u, _ := url.Parse("https://www.strava.com/api/v3/push_subscriptions")

	params := url.Values{}
	params.Add("client_id", a.ClientID)
	params.Add("client_secret", a.ClientSecret)
	params.Add("callback_url", callbackURL)
	params.Add("verify_token", verifyToken)

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

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to register webhook: %s\n%s", resp.Status, string(bodyBytes))
	}

	var r WebhookRegistrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

func (a *Auth) GetWebhookSubscriptions() ([]WebhookSubscription, error) {
	u, _ := url.Parse("https://www.strava.com/api/v3/push_subscriptions")

	params := url.Values{}
	params.Add("client_id", a.ClientID)
	params.Add("client_secret", a.ClientSecret)

	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
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
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get webhook subscriptions: %s\n%s", resp.Status, string(bodyBytes))
	}

	var r []WebhookSubscription
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, nil
}

func (a *Auth) DeleteWebhookSubscription(subscriptionID int) error {
	u, _ := url.Parse(fmt.Sprintf("https://www.strava.com/api/v3/push_subscriptions/%d", subscriptionID))

	params := url.Values{}
	params.Add("client_id", a.ClientID)
	params.Add("client_secret", a.ClientSecret)

	u.RawQuery = params.Encode()

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete webhook subscription: %s\n%s", resp.Status, string(bodyBytes))
	}

	return nil
}
