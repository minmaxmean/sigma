package ollama

import (
	"context"
	"fmt"
	"log"
	"os"
)

// Example demonstrates basic usage of the Ollama client wrapper
func Example() {
	// Create a new client
	client, err := NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Simple text generation
	req := GenerateRequest{
		Model:        "gemma3:12b",
		Prompt:       "What is the capital of France?",
		SystemPrompt: "Provide brief, concise responses",
	}

	response, err := client.Generate(ctx, req)
	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}

	fmt.Printf("Response: %s\n", response)

	// Example 2: Generate with file content
	if len(os.Args) > 1 {
		filename := os.Args[1]
		fileReq := GenerateRequest{
			Model:        "gemma3:12b",
			Prompt:       "Summarize this data",
			SystemPrompt: "Provide a brief summary",
		}

		fileResponse, err := client.GenerateWithFile(ctx, fileReq, filename)
		if err != nil {
			log.Fatalf("Failed to generate response with file: %v", err)
		}

		fmt.Printf("File Response: %s\n", fileResponse)
	}

}
