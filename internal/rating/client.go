package rating

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AccountClient struct {
	baseURL string
	client  *http.Client
}

func NewAccountClient(baseURL string) *AccountClient {
	return &AccountClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type updateRatingBody struct {
	Mu    float64 `json:"mu"`
	Sigma float64 `json:"sigma"`
}

func (c *AccountClient) UpdateRating(playerID string, mu, sigma float64) error {
	body, err := json.Marshal(updateRatingBody{Mu: mu, Sigma: sigma})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/players/%s/rating", c.baseURL, playerID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("account-service returned status %d", resp.StatusCode)
	}
	return nil
}
