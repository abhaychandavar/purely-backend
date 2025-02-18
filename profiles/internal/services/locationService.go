package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"profiles/internal/config"
	"profiles/internal/types/locationServiceTypes"
)

type LocationService struct {
}

type LocationResponse struct {
	Results       []interface{} `json:"results"`
	NextPageToken *string       `json:"next_page_token"`
}

func (locationService *LocationService) GetLocations(c *context.Context, data locationServiceTypes.GetLocationsType) (interface{}, error) {
	baseURL := "https://maps.googleapis.com/maps/api/place/textsearch/json"
	apiKey := config.GetConfig().GoogleMapsAPIKey // Replace with your Google Maps API key

	// Build query parameters
	query := url.Values{}
	query.Set("key", apiKey)
	query.Set("types", "locality")

	// Validate and set parameters
	if data.Query != nil {
		query.Set("query", *data.Query)
	} else {
		return nil, fmt.Errorf("either 'place' or both 'lat' and 'lng' must be provided")
	}

	if data.PageToken != nil {
		query.Set("pagetoken", *data.PageToken)
	}

	// Construct the full URL
	fullURL := fmt.Sprintf("%s?%s", baseURL, query.Encode())
	log.Default().Printf("fullURL %v", fullURL)
	// Make the HTTP request
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var locationResponse LocationResponse
	if err := json.Unmarshal(body, &locationResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return locationServiceTypes.GetLocationsResponseType{
		Results:       locationResponse.Results,
		NextPageToken: locationResponse.NextPageToken,
	}, nil
}
