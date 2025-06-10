package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"carbon_intensity/models"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewHTTPClient() *Client {
	return &Client{
		BaseURL: "https://api.carbonintensity.org.uk",
		HTTPClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) GetCarbonIntensityForecast(ctx context.Context, fromTime time.Time) ([]models.CarbonIntensityPeriod, error) {
	response := models.ExternalAPIResponse{}
	fromTimeStr := fromTime.Format("2006-01-02T15:04")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(c.BaseURL+"/intensity/%s/fw24h", fromTimeStr), nil)
	if err != nil {
		return response.Data, err
	}

	fmt.Println("Request URL:", req.URL.String())

	req.Header.Set("Content-Type", "application/json")

	if err = c.do(req, &response); err != nil {
		return response.Data, fmt.Errorf("failed to get carbon intensity forecast: %w", err)
	}

	return response.Data, nil

}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error response: %s", resp.Status)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}
