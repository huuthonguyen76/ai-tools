package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	logger "example.com/m/v2/internal/pkg"
	"go.uber.org/zap"
)

type DifyServiceInterface interface {
	GetContextualLink(url string) (string, error)
}

type DifyService struct {
	BaseURL string
	APIKey  string
}

type DifyData struct {
	ID          string      `json:"id"`
	WorkflowID  string      `json:"workflow_id"`
	Status      string      `json:"status"`
	Outputs     interface{} `json:"outputs"`
	Error       string      `json:"error"`
	ElapsedTime float64     `json:"elapsed_time"`
}

type DifyResponse struct {
	TaskID        string   `json:"task_id"`
	WorkflowRunID string   `json:"workflow_run_id"`
	Data          DifyData `json:"data"`
}

type ContextualizeLinkResponse struct {
	Link string `json:"link"`
}

func NewDifyService(baseURL, apiKey string) *DifyService {
	return &DifyService{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

func (s *DifyService) callDify(payload map[string]interface{}) (*DifyResponse, error) {
	var response DifyResponse

	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.BaseURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.APIKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *DifyService) GetContextualLink(url string) (string, error) {
	var contextualizeLinkResponse ContextualizeLinkResponse

	payload := map[string]interface{}{
		"inputs": map[string]string{
			"link": url,
		},
		"response_mode": "blocking",
		"user":          "tho-local",
	}

	difyResponse, err := s.callDify(payload)
	if err != nil {
		return "", err
	}

	logger.Info(context.Background(), "response: ", zap.Any("response", difyResponse.Data.Outputs))

	output, err := json.Marshal(difyResponse.Data.Outputs)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(output, &contextualizeLinkResponse)
	if err != nil {
		return "", err
	}

	logger.Info(context.Background(), "contextualizeLinkResponse: ", zap.Any("contextualizeLinkResponse", contextualizeLinkResponse))

	return contextualizeLinkResponse.Link, nil
}
