package strava

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gabrieleangeletti/stride"
)

func NewClient(accessToken string) *client {
	return &client{
		BaseUrl:     "https://www.strava.com/api/v3",
		AccessToken: accessToken,
	}
}

type client struct {
	BaseUrl     string
	AccessToken string
}

func (c *client) GetAthlete() (*Athlete, error) {
	u := fmt.Sprintf("%s/athlete", c.BaseUrl)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, stride.ErrRateLimitExceeded
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get athlete: %s\n%s", resp.Status, string(bodyBytes))
	}

	var athlete Athlete
	if err := json.NewDecoder(resp.Body).Decode(&athlete); err != nil {
		return nil, err
	}

	return &athlete, nil
}

func (c *client) GetActivitySummaries(startTime, endTime time.Time, page int) ([]ActivitySummary, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/athlete/activities", c.BaseUrl))

	params := url.Values{}
	params.Add("after", fmt.Sprintf("%d", startTime.Unix()))
	params.Add("before", fmt.Sprintf("%d", endTime.Unix()))
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("per_page", "200")

	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, stride.ErrRateLimitExceeded
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get activities: %s\n%s", resp.Status, string(bodyBytes))
	}

	var activities []ActivitySummary
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return nil, err
	}

	return activities, nil
}

func (c *client) GetActivitySummary(id int, includeAllEfforts bool) (*ActivitySummary, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/activities/%d", c.BaseUrl, id))

	params := url.Values{}
	params.Add("include_all_efforts", strconv.FormatBool(includeAllEfforts))

	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, stride.ErrRateLimitExceeded
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get activities: %s\n%s", resp.Status, string(bodyBytes))
	}

	var activity ActivitySummary
	if err := json.NewDecoder(resp.Body).Decode(&activity); err != nil {
		return nil, err
	}

	return &activity, nil
}

func (c *client) GetActivityStreams(id int) (*ActivityStream, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/activities/%d/streams", c.BaseUrl, id))

	params := url.Values{}
	params.Add("keys", "velocity_smooth,cadence,distance,altitude,heartrate,time")
	params.Add("key_by_type", "true")

	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, stride.ErrRateLimitExceeded
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get activity streams: %s\n%s", resp.Status, string(bodyBytes))
	}

	var streams ActivityStream
	if err := json.NewDecoder(resp.Body).Decode(&streams); err != nil {
		return nil, err
	}

	return &streams, nil
}
