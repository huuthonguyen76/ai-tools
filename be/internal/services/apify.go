package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"time"

	"example.com/m/v2/internal/pkg"
	logger "example.com/m/v2/internal/pkg"
	"go.uber.org/zap"
)

// ApifyConfig holds the Apify API configuration
type ApifyConfig struct {
	APIToken string
	BaseURL  string
}

// ApifyService provides methods for interacting with Apify API
type ApifyService struct {
	Config     *ApifyConfig
	HTTPClient pkg.CustomHTTPClient
}

type ApifyServiceInterface interface {
	RunActor(ctx context.Context, actorID string, input interface{}, options *ActorRunOptions) (*ActorRun, error)
	BuildApifyURL(url string) (string, error)
	GetRawDatasetItems(ctx context.Context, datasetID string, offset, limit int) ([]map[string]interface{}, error)
	GetAllDatasetIDs(ctx context.Context) ([]string, error)
	RunFacebookScraper(ctx context.Context, groupURLs []string, onlyPostsNewerThan string, viewOption string) (*ActorRun, error)
}

// ActorRunInput represents the input configuration for running an actor
type ActorRunInput struct {
	OnlyPostsNewerThan string      `json:"onlyPostsNewerThan,omitempty"`
	StartUrls          []StartURL  `json:"startUrls,omitempty"`
	ViewOption         string      `json:"viewOption,omitempty"`
	MaxPosts           int         `json:"maxPosts,omitempty"`
	CustomInput        interface{} `json:"-"` // For any custom actor-specific input
}

// StartURL represents a URL configuration for the actor
type StartURL struct {
	URL string `json:"url"`
}

// ActorRunOptions contains options for running an actor
type ActorRunOptions struct {
	Build         string        `json:"build,omitempty"`              // Actor build to run (e.g., "latest", "beta")
	Memory        int           `json:"memory,omitempty"`             // Memory in MB (e.g., 128, 256, 512, 1024, etc.)
	Timeout       int           `json:"timeout,omitempty"`            // Timeout in seconds
	WaitForFinish int           `json:"waitForFinish,omitempty"`      // Wait for actor to finish (in seconds)
	Webhooks      []interface{} `json:"webhooks,omitempty"`           // Webhooks to be called on actor events
	MaxItems      int           `json:"maxItems,omitempty"`           // Maximum number of items to scrape
	ProxyGroup    string        `json:"proxyConfiguration,omitempty"` // Proxy configuration
}

// ActorRunResponse represents the response from starting an actor run
type ActorRunResponse struct {
	Data ActorRun `json:"data"`
}

// ActorRun represents an actor run
type ActorRun struct {
	ID                     string                 `json:"id"`
	ActID                  string                 `json:"actId"`
	Status                 string                 `json:"status"` // READY, RUNNING, SUCCEEDED, FAILED, ABORTED
	StartedAt              time.Time              `json:"startedAt"`
	FinishedAt             *time.Time             `json:"finishedAt,omitempty"`
	BuildID                string                 `json:"buildId"`
	ExitCode               int                    `json:"exitCode"`
	DefaultKeyValueStoreID string                 `json:"defaultKeyValueStoreId"`
	DefaultDatasetID       string                 `json:"defaultDatasetId"`
	DefaultRequestQueueID  string                 `json:"defaultRequestQueueId"`
	BuildNumber            string                 `json:"buildNumber"`
	ContainerURL           string                 `json:"containerUrl"`
	Meta                   map[string]interface{} `json:"meta,omitempty"`
	Stats                  ActorRunStats          `json:"stats"`
	Options                map[string]interface{} `json:"options,omitempty"`
	UsageTotalUSD          float64                `json:"usageTotalUsd"`
	UsageUSD               map[string]float64     `json:"usageUsd,omitempty"`
}

// ActorRunStats contains statistics about the actor run
type ActorRunStats struct {
	InputBodyLen    int     `json:"inputBodyLen"`
	RestartCount    int     `json:"restartCount"`
	ResurrectCount  int     `json:"resurrectCount"`
	MemAvgBytes     int64   `json:"memAvgBytes"`
	MemMaxBytes     int64   `json:"memMaxBytes"`
	MemCurrentBytes int64   `json:"memCurrentBytes"`
	CpuAvgUsage     float64 `json:"cpuAvgUsage"`
	CpuMaxUsage     float64 `json:"cpuMaxUsage"`
	CpuCurrentUsage float64 `json:"cpuCurrentUsage"`
	NetRxBytes      int64   `json:"netRxBytes"`
	NetTxBytes      int64   `json:"netTxBytes"`
	DurationMillis  int64   `json:"durationMillis"`
	RunTimeSecs     float64 `json:"runTimeSecs"`
	MetamorphMillis int64   `json:"metamorphMillis"`
	ComputeUnits    float64 `json:"computeUnits"`
}

// DatasetItem represents a single item from the dataset
type DatasetItem struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	UserID         string    `json:"userId"`
	UserName       string    `json:"userName"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	AccessedAt     time.Time `json:"accessedAt"`
	ItemCount      int       `json:"itemCount"`
	CleanItemCount int       `json:"cleanItemCount"`
	ActID          string    `json:"actId"`
	ActRunID       string    `json:"actRunId"`
}

// DatasetResponse represents the response from fetching dataset items
type DatasetResponse struct {
	Data   DatasetData `json:"data"`
	Count  int         `json:"count"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
	Total  int         `json:"total"`
}

// DatasetData holds the actual dataset items
type DatasetData struct {
	Items []DatasetItem `json:"items"`
}

// NewApifyService creates a new Apify service instance
func NewApifyService(apiToken string) (*ApifyService, error) {
	if apiToken == "" {
		return nil, errors.New("API token cannot be empty")
	}

	return &ApifyService{
		Config: &ApifyConfig{
			APIToken: apiToken,
			BaseURL:  "https://api.apify.com/v2",
		},
		HTTPClient: pkg.NewHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
		}),
	}, nil
}

func (s *ApifyService) BuildApifyURL(url string) (string, error) {
	linkParsed, err := urlpkg.Parse(url)
	if err != nil {
		return "", err
	}

	q := linkParsed.Query()
	q.Set("token", s.Config.APIToken)

	linkParsed.RawQuery = q.Encode()
	return linkParsed.String(), nil
}

func (s *ApifyService) RunActor(ctx context.Context, actorID string, input interface{}, options *ActorRunOptions) (*ActorRun, error) {
	var requestBody []byte

	if actorID == "" {
		return nil, errors.New("actor ID cannot be empty")
	}

	// Build the API URL
	url := fmt.Sprintf("%s/acts/%s/runs", s.Config.BaseURL, actorID)
	apifyURL, err := s.BuildApifyURL(url)
	if err != nil {
		logger.Error(ctx, "Failed to build Apify URL", zap.Error(err))
		return nil, fmt.Errorf("failed to build Apify URL: %w", err)
	}

	resp, rawBody, err := s.HTTPClient.Call(ctx, "POST", apifyURL, &requestBody)
	if err != nil {
		logger.Error(ctx, "Error: RunActor", zap.Error(err))
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.Error(ctx, "Actor run request failed",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(*rawBody)))
		return nil, fmt.Errorf("actor run request failed with status %d: %s", resp.StatusCode, string(*rawBody))
	}

	// Parse the response
	var runResponse ActorRunResponse
	if err := json.Unmarshal(*rawBody, &runResponse); err != nil {
		logger.Error(ctx, "Failed to unmarshal actor run response", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal actor run response: %w", err)
	}

	logger.Info(ctx, "Actor run started successfully",
		zap.String("run_id", runResponse.Data.ID),
		zap.String("status", runResponse.Data.Status))

	return &runResponse.Data, nil
}

func (s *ApifyService) GetAllDatasetIDs(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/datasets?unnamed=1", s.Config.BaseURL)
	apifyURL, err := s.BuildApifyURL(url)
	if err != nil {
		logger.Error(ctx, "Failed to build Apify URL", zap.Error(err))
		return nil, fmt.Errorf("failed to build Apify URL: %w", err)
	}

	_, rawBody, err := s.HTTPClient.Call(ctx, "GET", apifyURL, nil)
	if err != nil {
		logger.Error(ctx, "Error: GetAllDatasetIDs", zap.Error(err))
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var datasets DatasetResponse
	if err := json.Unmarshal(*rawBody, &datasets); err != nil {
		logger.Error(ctx, "Error: GetAllDatasetIDs", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal datasets: %w", err)
	}

	datasetIDs := make([]string, len(datasets.Data.Items))
	for i, dataset := range datasets.Data.Items {
		datasetIDs[i] = dataset.ID
	}

	return datasetIDs, nil
}

// GetRawDatasetItems retrieves raw items from a dataset as unstructured data
func (s *ApifyService) GetRawDatasetItems(ctx context.Context, datasetID string, offset, limit int) ([]map[string]interface{}, error) {
	if datasetID == "" {
		return nil, errors.New("dataset ID cannot be empty")
	}

	if limit <= 0 {
		limit = 1000
	}

	url := fmt.Sprintf("%s/datasets/%s/items?offset=%d&limit=%d",
		s.Config.BaseURL, datasetID, s.Config.APIToken, offset, limit)

	apifyURL, err := s.BuildApifyURL(url)
	if err != nil {
		logger.Error(ctx, "Failed to build Apify URL", zap.Error(err))
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	_, rawBody, err := s.HTTPClient.Call(ctx, "GET", apifyURL, nil)
	if err != nil {
		logger.Error(ctx, "Error: GetRawDatasetItems", zap.Error(err))
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// The dataset items endpoint returns an array directly
	var items []map[string]interface{}
	if err := json.Unmarshal(*rawBody, &items); err != nil {
		logger.Error(ctx, "Error: GetRawDatasetItems", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal raw dataset items: %w", err)
	}

	logger.Info(ctx, "Retrieved raw dataset items",
		zap.String("dataset_id", datasetID),
		zap.Int("count", len(items)))

	return items, nil
}

// GetAllRawDatasetItems retrieves all items from a dataset with pagination
func (s *ApifyService) GetAllRawDatasetItems(ctx context.Context, datasetID string) ([]map[string]interface{}, error) {
	if datasetID == "" {
		return nil, errors.New("dataset ID cannot be empty")
	}

	var allItems []map[string]interface{}
	offset := 0
	limit := 1000

	for {
		items, err := s.GetRawDatasetItems(ctx, datasetID, offset, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to get items at offset %d: %w", offset, err)
		}

		if len(items) == 0 {
			break
		}

		allItems = append(allItems, items...)
		offset += len(items)

		// If we got fewer items than the limit, we've reached the end
		if len(items) < limit {
			break
		}
	}

	logger.Info(ctx, "Retrieved all raw dataset items",
		zap.String("dataset_id", datasetID),
		zap.Int("total_count", len(allItems)))

	return allItems, nil
}

// Example function for running Facebook scraper actor
func (s *ApifyService) RunFacebookScraper(ctx context.Context, groupURLs []string, onlyPostsNewerThan string, viewOption string) (*ActorRun, error) {
	startUrls := make([]StartURL, len(groupURLs))
	for i, url := range groupURLs {
		startUrls[i] = StartURL{URL: url}
	}

	input := ActorRunInput{
		OnlyPostsNewerThan: onlyPostsNewerThan,
		StartUrls:          startUrls,
		ViewOption:         viewOption,
	}

	// Replace with your actual Facebook scraper actor ID
	actorID := "apify~facebook-groups-scraper" // Update this with the correct actor ID

	return s.RunActor(ctx, actorID, input, &ActorRunOptions{
		Memory:  4096,
		Timeout: 3600,
	})
}
