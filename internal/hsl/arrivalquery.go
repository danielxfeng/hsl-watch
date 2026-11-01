package hsl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielxfeng/hsl-watch/internal/model"
)

const endpoint = "https://api.digitransit.fi/routing/v2/hsl/gtfs/v1"

const arrivalQuery = `
query Arrivals($gtfsIds: [String!]!) {
  stops(ids: $gtfsIds) {
    gtfsId
    name
    stoptimesWithoutPatterns(numberOfDepartures: 3) {
      serviceDay
      realtimeArrival
      realtime
      headsign
    }
  }
}`

type HslResponseDTO struct {
	model.Arrival
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *HslClient) GetArrivals(ctx context.Context, ids []string) (*model.Arrival, error) {
	payload := struct {
		Query     string `json:"query"`
		Variables struct {
			GtfsIDs []string `json:"gtfsIds"`
		} `json:"variables"`
	}{
		Query: arrivalQuery,
	}

	payload.Variables.GtfsIDs = ids

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("digitransit-subscription-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HSL API returned status %s", resp.Status)
	}

	var result HslResponseDTO

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode HSL response: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf(
			"HSL GraphQL error: %s",
			result.Errors[0].Message,
		)
	}
	return &result.Arrival, nil
}
