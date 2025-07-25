package ollama

import (
	"context"
	"fmt"
	"os"

	"github.com/ollama/ollama/api"
)

// Client wraps the Ollama API client with simplified methods
type Client struct {
	apiClient *api.Client
}

// NewClient creates a new Ollama client using environment configuration
func NewClient() (*Client, error) {
	apiClient, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ollama client: %w", err)
	}
	return &Client{apiClient: apiClient}, nil
}

// GenerateRequest represents a request to generate text
type GenerateRequest struct {
	Model        string
	Prompt       string
	SystemPrompt string
}

// GenerateResponse represents the response from text generation
type GenerateResponse struct {
	Response string
	Done     bool
}

// Generate sends a text generation request to Ollama
func (c *Client) Generate(ctx context.Context, req GenerateRequest) (string, error) {
	apiReq := &api.GenerateRequest{
		Model:  req.Model,
		Prompt: req.Prompt,
		System: req.SystemPrompt,
	}

	var response string
	respFunc := func(resp api.GenerateResponse) error {
		response += resp.Response
		return nil
	}

	if err := c.apiClient.Generate(ctx, apiReq, respFunc); err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	return response, nil
}

// GenerateWithFile combines file content with a prompt and generates a response
func (c *Client) GenerateWithFile(ctx context.Context, req GenerateRequest, filename string) (string, error) {
	// Read file content
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Combine file content with prompt
	combinedPrompt := fmt.Sprintf("data: %s\nprompt: %s\n", string(fileContent), req.Prompt)

	// Create new request with combined prompt
	newReq := req
	newReq.Prompt = combinedPrompt

	return c.Generate(ctx, newReq)
}
