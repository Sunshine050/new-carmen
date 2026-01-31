package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/new-carmen/backend/internal/config"
)

type Client struct {
	BaseURL string
	Model   string
	client  *http.Client
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type AnalyzeRequest struct {
	Question string `json:"question"`
}

type AnalyzeResponse struct {
	IsAmbiguous bool     `json:"is_ambiguous"`
	Reason      string   `json:"reason,omitempty"`
	Candidates  []string `json:"candidates,omitempty"`
}

func NewClient() *Client {
	cfg := config.AppConfig.Ollama
	return &Client{
		BaseURL: cfg.URL,
		Model:   cfg.Model,
		client:  &http.Client{},
	}
}

func (c *Client) Generate(prompt string) (string, error) {
	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.client.Post(
		fmt.Sprintf("%s/api/generate", c.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var generateResp GenerateResponse
	if err := json.Unmarshal(body, &generateResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return generateResp.Response, nil
}

func (c *Client) AnalyzeQuestionClarity(question string) (*AnalyzeResponse, error) {
	prompt := fmt.Sprintf(
		"Analyze if this question is ambiguous or unclear: '%s'. "+
			"Respond with JSON: {\"is_ambiguous\": true/false, \"reason\": \"...\", \"candidates\": [...]}",
		question,
	)

	response, err := c.Generate(prompt)
	if err != nil {
		return nil, err
	}

	var analyzeResp AnalyzeResponse
	if err := json.Unmarshal([]byte(response), &analyzeResp); err != nil {
		// If JSON parsing fails, assume not ambiguous
		return &AnalyzeResponse{IsAmbiguous: false}, nil
	}

	return &analyzeResp, nil
}

func (c *Client) GenerateClarifyingQuestion(question string, candidates []string) (string, error) {
	candidatesStr := ""
	for i, c := range candidates {
		if i > 0 {
			candidatesStr += ", "
		}
		candidatesStr += fmt.Sprintf("%d. %s", i+1, c)
	}

	prompt := fmt.Sprintf(
		"Generate a clarifying question for this ambiguous question: '%s'. "+
			"Possible meanings: %s. "+
			"Create a clear, concise question that helps the user specify what they mean.",
		question, candidatesStr,
	)

	return c.Generate(prompt)
}

func (c *Client) GenerateAnswer(context string, question string) (string, error) {
	prompt := fmt.Sprintf(
		"Based on the following context, answer the question. "+
			"Context: %s\n\nQuestion: %s\n\nAnswer:",
		context, question,
	)

	return c.Generate(prompt)
}
