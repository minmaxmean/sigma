# Ollama Client Wrapper

A simple wrapper around the official Ollama Go client that provides easy-to-use methods for making requests to Ollama models.

## Features

- Simple client creation and initialization
- Text generation with custom prompts and system messages
- File-based generation (combines file content with prompts)
- Error handling and context support

## Usage

### Basic Setup

```go
import "github.com/minmaxmean/sigma/ollama"

// Create a new client
client, err := ollama.NewClient()
if err != nil {
    log.Fatal(err)
}
```

### Simple Text Generation

```go
ctx := context.Background()
req := ollama.GenerateRequest{
    Model:        "gemma3:12b",
    Prompt:       "What is the capital of France?",
    SystemPrompt: "Provide brief, concise responses",
}

response, err := client.Generate(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Println(response)
```

### Generate with File Content

```go
req := ollama.GenerateRequest{
    Model:        "gemma3:12b",
    Prompt:       "Summarize this data",
    SystemPrompt: "Provide a brief summary",
}

response, err := client.GenerateWithFile(ctx, req, "data.txt")
if err != nil {
    log.Fatal(err)
}

fmt.Println(response)
```



## API Reference

### Types

#### `Client`
The main client wrapper that provides methods for interacting with Ollama.

#### `GenerateRequest`
```go
type GenerateRequest struct {
    Model        string // The model to use (e.g., "gemma3:12b")
    Prompt       string // The main prompt
    SystemPrompt string // System message to guide the model
}
```

### Methods

#### `NewClient() (*Client, error)`
Creates a new Ollama client using environment configuration.

#### `Generate(ctx context.Context, req GenerateRequest) (string, error)`
Generates a text response for the given request.

#### `GenerateWithFile(ctx context.Context, req GenerateRequest, filename string) (string, error)`
Combines file content with the prompt and generates a response.



## Requirements

- Go 1.24.4 or later
- Ollama server running locally or accessible via environment variables
- The `github.com/ollama/ollama` dependency (already included in go.mod)

## Environment Configuration

The client uses the official Ollama client's environment configuration:
- `OLLAMA_HOST`: Ollama server host (default: "http://localhost:11434")
- `OLLAMA_ORIGINS`: Allowed origins for CORS

## Testing

Run the tests with:

```bash
go test ./ollama
```

Tests will be skipped if Ollama is not available or if the required model is not installed. 