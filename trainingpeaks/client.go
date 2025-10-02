package trainingpeaks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func NewClient(accessToken string) *client {
	return &client{
		BaseUrl:     "https://tpapi.trainingpeaks.com/fitness/v6",
		AccessToken: accessToken,
	}
}

type client struct {
	BaseUrl     string
	AccessToken string
}

func (p *client) GetWorkoutSummaries(athleteID int, startTime, endTime time.Time) ([]TrainingPeaksWorkoutSummary, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/athletes/%d/workouts/%s/%s", p.BaseUrl, athleteID, startTime.Format("2006-01-02"), endTime.Format("2006-01-02")))

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get workout summaries: %s", resp.Status)
	}

	var summaries []TrainingPeaksWorkoutSummary
	err = json.NewDecoder(resp.Body).Decode(&summaries)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

func (p *client) GetWorkoutDetails(athleteID int, workoutID string) (TrainingPeaksWorkoutDetail, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/athletes/%d/workouts/%s/detaildata", p.BaseUrl, athleteID, workoutID))

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return TrainingPeaksWorkoutDetail{}, err
	}

	req.Header.Set("Authorization", "Bearer "+p.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TrainingPeaksWorkoutDetail{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TrainingPeaksWorkoutDetail{}, fmt.Errorf("failed to get workout details: %s", resp.Status)
	}

	var detail TrainingPeaksWorkoutDetail
	err = json.NewDecoder(resp.Body).Decode(&detail)
	if err != nil {
		return TrainingPeaksWorkoutDetail{}, err
	}

	return detail, nil
}
