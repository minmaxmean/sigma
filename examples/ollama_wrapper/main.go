package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minmaxmean/sigma/ollama"
)

func main() {
	// Create a new client
	client, err := ollama.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Simple text generation
	fmt.Println("=== Simple Text Generation ===")
	req := ollama.GenerateRequest{
		Model:        "gemma3:12b",
		Prompt:       "What is the capital of France?",
		SystemPrompt: "Provide brief, concise responses",
	}

	response, err := client.Generate(ctx, req)
	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}

	fmt.Printf("Response: %s\n\n", response)

	// Example 2: Generate with file content (if file provided)
	if len(os.Args) > 1 {
		fmt.Println("=== File-based Generation ===")
		filename := os.Args[1]
		fileReq := ollama.GenerateRequest{
			Model:        "gemma3:12b",
			Prompt:       "Summarize this data",
			SystemPrompt: "Provide a brief summary",
		}

		fileResponse, err := client.GenerateWithFile(ctx, fileReq, filename)
		if err != nil {
			log.Fatalf("Failed to generate response with file: %v", err)
		}

		fmt.Printf("File Response: %s\n\n", fileResponse)
	}

	fmt.Println("\nDone!")
}
