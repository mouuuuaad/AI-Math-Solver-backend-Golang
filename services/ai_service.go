package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"maths-solution-backend/config"
	"maths-solution-backend/models"
)

type AIService struct {
	config *config.Config
	client *http.Client
}

func NewAIService(cfg *config.Config) *AIService {
	return &AIService{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.AI.Timeout,
		},
	}
}

type AIRequest struct {
	Expression string `json:"expression"`
}

type AIResponse struct {
	Steps []models.SolutionStep `json:"steps"`
	Final string                `json:"final"`
}

func (s *AIService) SolveMath(expression string) (*AIResponse, error) {
	reqBody := AIRequest{Expression: expression}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", s.config.AI.ServiceURL+"/solve", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	var aiResp AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, fmt.Errorf("failed to decode AI service response: %w", err)
	}

	return &aiResp, nil
}
