package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	logger "example.com/m/v2/internal/pkg"
	"example.com/m/v2/internal/services"
	"github.com/joho/godotenv"
)

// TestApifyService_Integration tests the Apify service with mocked HTTP server
func TestApifyService_Integration(t *testing.T) {
	logger.Init(false)

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	service, err := services.NewApifyService(os.Getenv("APIFY_API_KEY"))
	if err != nil {
		t.Fatalf("Failed to create Apify service: %v", err)
	}

	t.Run("NewApifyService_Success", func(t *testing.T) {
		// service, err := services.NewApifyService("test-token")
		// datasetIDs, err := service.GetAllDatasetIDs(context.Background())
		// if err != nil {
		// 	t.Fatalf("Failed to get dataset IDs: %v", err)
		// }
		// fmt.Println("datasetIDs", datasetIDs)

		fmt.Println(service.GetAllRawDatasetItems(context.Background(), "PS9BhFCDSvRpYvLYu"))
		// require.NoError(t, err)
		// assert.NotNil(t, service)
		// assert.Equal(t, "test-token", service.Config.APIToken)
		// assert.Equal(t, "https://api.apify.com/v2", service.Config.BaseURL)
	})

	// t.Run("NewApifyService_EmptyToken", func(t *testing.T) {
	// 	service, err := services.NewApifyService("")
	// 	assert.Error(t, err)
	// 	assert.Nil(t, service)
	// 	assert.Equal(t, "API token cannot be empty", err.Error())
	// })

	// t.Run("BuildApifyURL_Success", func(t *testing.T) {
	// 	service, err := services.NewApifyService("test-token-123")
	// 	require.NoError(t, err)

	// 	url, err := service.BuildApifyURL("https://api.apify.com/v2/datasets")
	// 	require.NoError(t, err)
	// 	assert.Contains(t, url, "token=test-token-123")
	// })

	// t.Run("BuildApifyURL_InvalidURL", func(t *testing.T) {
	// 	service, err := services.NewApifyService("test-token")
	// 	require.NoError(t, err)

	// 	_, err = service.BuildApifyURL("://invalid-url")
	// 	assert.Error(t, err)
	// })
}

// // TestApifyService_GetAllDatasetIDs tests the GetAllDatasetIDs method with a mock server
// func TestApifyService_GetAllDatasetIDs(t *testing.T) {
// 	t.Run("Success_ReturnsDatasetIDs", func(t *testing.T) {
// 		// Create mock server
// 		mockResponse := services.DatasetResponse{
// 			Data: services.DatasetData{
// 				Items: []services.DatasetItem{
// 					{ID: "dataset-1", Name: "Test Dataset 1"},
// 					{ID: "dataset-2", Name: "Test Dataset 2"},
// 					{ID: "dataset-3", Name: "Test Dataset 3"},
// 				},
// 			},
// 			Count:  3,
// 			Offset: 0,
// 			Limit:  100,
// 			Total:  3,
// 		}

// 		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			assert.Equal(t, "GET", r.Method)
// 			assert.Contains(t, r.URL.Path, "/datasets")
// 			assert.Contains(t, r.URL.RawQuery, "token=")
// 			assert.Contains(t, r.URL.RawQuery, "unnamed=1")

// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusOK)
// 			json.NewEncoder(w).Encode(mockResponse)
// 		}))
// 		defer server.Close()

// 		// Create service with mock server URL
// 		service := &services.ApifyService{
// 			Config: &services.ApifyConfig{
// 				APIToken: "test-token",
// 				BaseURL:  server.URL,
// 			},
// 			HTTPClient: pkg.NewHTTPClient(&http.Client{Timeout: 10 * time.Second}),
// 		}

// 		ctx := context.Background()
// 		datasetIDs, err := service.GetAllDatasetIDs(ctx)

// 		require.NoError(t, err)
// 		assert.Len(t, datasetIDs, 3)
// 		assert.Contains(t, datasetIDs, "dataset-1")
// 		assert.Contains(t, datasetIDs, "dataset-2")
// 		assert.Contains(t, datasetIDs, "dataset-3")
// 	})

// 	t.Run("Success_EmptyDatasets", func(t *testing.T) {
// 		mockResponse := services.DatasetResponse{
// 			Data: services.DatasetData{
// 				Items: []services.DatasetItem{},
// 			},
// 			Count:  0,
// 			Offset: 0,
// 			Limit:  100,
// 			Total:  0,
// 		}

// 		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusOK)
// 			json.NewEncoder(w).Encode(mockResponse)
// 		}))
// 		defer server.Close()

// 		service := &services.ApifyService{
// 			Config: &services.ApifyConfig{
// 				APIToken: "test-token",
// 				BaseURL:  server.URL,
// 			},
// 			HTTPClient: pkg.NewHTTPClient(&http.Client{Timeout: 10 * time.Second}),
// 		}

// 		ctx := context.Background()
// 		datasetIDs, err := service.GetAllDatasetIDs(ctx)

// 		require.NoError(t, err)
// 		assert.Len(t, datasetIDs, 0)
// 	})

// 	t.Run("Error_ServerError", func(t *testing.T) {
// 		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte("Internal Server Error"))
// 		}))
// 		defer server.Close()

// 		service := &services.ApifyService{
// 			Config: &services.ApifyConfig{
// 				APIToken: "test-token",
// 				BaseURL:  server.URL,
// 			},
// 			HTTPClient: pkg.NewHTTPClient(&http.Client{Timeout: 10 * time.Second}),
// 		}

// 		ctx := context.Background()
// 		_, err := service.GetAllDatasetIDs(ctx)

// 		assert.Error(t, err)
// 	})

// 	t.Run("Error_InvalidJSON", func(t *testing.T) {
// 		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusOK)
// 			w.Write([]byte("invalid json"))
// 		}))
// 		defer server.Close()

// 		service := &services.ApifyService{
// 			Config: &services.ApifyConfig{
// 				APIToken: "test-token",
// 				BaseURL:  server.URL,
// 			},
// 			HTTPClient: pkg.NewHTTPClient(&http.Client{Timeout: 10 * time.Second}),
// 		}

// 		ctx := context.Background()
// 		_, err := service.GetAllDatasetIDs(ctx)

// 		assert.Error(t, err)
// 		assert.Contains(t, err.Error(), "failed to unmarshal")
// 	})
// }

// // TestApifyService_GetRawDatasetItems tests the GetRawDatasetItems method
// func TestApifyService_GetRawDatasetItems(t *testing.T) {
// 	t.Run("Success_ReturnsItems", func(t *testing.T) {
// 		mockItems := []map[string]interface{}{
// 			{"id": "post-1", "text": "Test post 1", "likes": float64(100)},
// 			{"id": "post-2", "text": "Test post 2", "likes": float64(200)},
// 		}

// 		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			assert.Equal(t, "GET", r.Method)
// 			assert.Contains(t, r.URL.Path, "/datasets/")
// 			assert.Contains(t, r.URL.Path, "/items")

// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusOK)
// 			json.NewEncoder(w).Encode(mockItems)
// 		}))
// 		defer server.Close()

// 		service := &services.ApifyService{
// 			Config: &services.ApifyConfig{
// 				APIToken: "test-token",
// 				BaseURL:  server.URL,
// 			},
// 			HTTPClient: pkg.NewHTTPClient(&http.Client{Timeout: 10 * time.Second}),
// 		}

// 		ctx := context.Background()
// 		items, err := service.GetRawDatasetItems(ctx, "test-dataset-id", 0, 100)

// 		require.NoError(t, err)
// 		assert.Len(t, items, 2)
// 		assert.Equal(t, "post-1", items[0]["id"])
// 		assert.Equal(t, "Test post 1", items[0]["text"])
// 	})

// 	t.Run("Error_EmptyDatasetID", func(t *testing.T) {
// 		service, err := services.NewApifyService("test-token")
// 		require.NoError(t, err)

// 		ctx := context.Background()
// 		_, err = service.GetRawDatasetItems(ctx, "", 0, 100)

// 		assert.Error(t, err)
// 		assert.Equal(t, "dataset ID cannot be empty", err.Error())
// 	})

// 	t.Run("DefaultLimit_WhenZeroOrNegative", func(t *testing.T) {
// 		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// Verify the limit parameter is set to default (1000)
// 			assert.Contains(t, r.URL.RawQuery, "limit=1000")

// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusOK)
// 			json.NewEncoder(w).Encode([]map[string]interface{}{})
// 		}))
// 		defer server.Close()

// 		service := &services.ApifyService{
// 			Config: &services.ApifyConfig{
// 				APIToken: "test-token",
// 				BaseURL:  server.URL,
// 			},
// 			HTTPClient: pkg.NewHTTPClient(&http.Client{Timeout: 10 * time.Second}),
// 		}

// 		ctx := context.Background()
// 		_, err := service.GetRawDatasetItems(ctx, "test-dataset", 0, 0)
// 		require.NoError(t, err)
// 	})
// }

// // TestApifyService_RunActor tests the RunActor method
// func TestApifyService_RunActor(t *testing.T) {
// 	t.Run("Error_EmptyActorID", func(t *testing.T) {
// 		service, err := services.NewApifyService("test-token")
// 		require.NoError(t, err)

// 		ctx := context.Background()
// 		_, err = service.RunActor(ctx, "", nil, nil)

// 		assert.Error(t, err)
// 		assert.Equal(t, "actor ID cannot be empty", err.Error())
// 	})
// }

// // TestApifyService_RunFacebookScraper tests the Facebook scraper integration
// func TestApifyService_RunFacebookScraper(t *testing.T) {
// 	t.Run("BuildsCorrectInput", func(t *testing.T) {
// 		// This test verifies that the Facebook scraper correctly builds the input
// 		groupURLs := []string{
// 			"https://facebook.com/groups/group1",
// 			"https://facebook.com/groups/group2",
// 		}

// 		// Verify the input structure is built correctly
// 		startUrls := make([]services.StartURL, len(groupURLs))
// 		for i, url := range groupURLs {
// 			startUrls[i] = services.StartURL{URL: url}
// 		}

// 		input := services.ActorRunInput{
// 			OnlyPostsNewerThan: "2024-01-01",
// 			StartUrls:          startUrls,
// 			ViewOption:         "full",
// 		}

// 		assert.Len(t, input.StartUrls, 2)
// 		assert.Equal(t, "https://facebook.com/groups/group1", input.StartUrls[0].URL)
// 		assert.Equal(t, "2024-01-01", input.OnlyPostsNewerThan)
// 		assert.Equal(t, "full", input.ViewOption)
// 	})
// }
