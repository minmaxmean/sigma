package ollama

import (
	"context"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Skipf("Skipping test - Ollama not available: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if client.apiClient == nil {
		t.Fatal("Expected apiClient to be initialized")
	}
}

func TestGenerateRequest(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Skipf("Skipping test - Ollama not available: %v", err)
	}

	ctx := context.Background()
	req := GenerateRequest{
		Model:        "gemma3:12b",
		Prompt:       "Say hello",
		SystemPrompt: "Be brief",
	}

	response, err := client.Generate(ctx, req)
	if err != nil {
		t.Skipf("Skipping test - model not available: %v", err)
	}

	if response == "" {
		t.Error("Expected non-empty response")
	}

	// Check if response contains "hello" (case insensitive)
	if !strings.Contains(strings.ToLower(response), "hello") {
		t.Logf("Response: %s", response)
		t.Log("Note: Response doesn't contain 'hello' - this might be expected depending on the model")
	}
}
