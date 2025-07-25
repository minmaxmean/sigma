package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ollama/ollama/api"
)

func main() {
	dataFile := os.Args[1]
	prompt := os.Args[2]

	ctx := context.Background()
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("Error initilizing ollama client: %v", err)
	}
	systemPrompt := "Provide very brief, concise responses"

	promptWithData, err := promptWithFile(dataFile, prompt)
	log.Printf("prompt loaded\n")

	req := &api.GenerateRequest{
		Model: "gemma3:12b",
		// Stream: new(bool),
		Prompt: promptWithData,
		System: systemPrompt,
	}

	respFunc := func(resp api.GenerateResponse) error {
		// log.Printf("Model: %+v\n", resp)
		fmt.Printf("%s", resp.Response)
		return nil
	}

	log.Printf("Generting response:\n")
	if err := client.Generate(ctx, req, respFunc); err != nil {
		log.Fatalf("Error chatting: %v", err)
	}
	log.Printf("Done!\n")
}

func promptWithFile(filename string, prompt string) (string, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	// file := "<data>"
	return fmt.Sprintf("data: %s\nprompt: %s\n", string(file), prompt), nil
}
