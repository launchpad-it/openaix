# OpenAIX

OpenAIX is a powerful, configuration-driven Go library designed to simplify interactions with the OpenAI API. It provides a robust and convenient interface for chat completions, including support for structured JSON responses, multi-modal inputs, and seamless client configuration for both OpenAI and Azure OpenAI services.

## Features

-   **Unified Client Factory**: Automatically configure an OpenAI client for standard OpenAI or Azure OpenAI services using environment variables.
-   **Configuration-Driven Completions**: Define your chat completion settings (model, temperature, prompts) in external configuration files (YAML, JSON, etc.) for easy management.
-   **Typed JSON Responses**: Leverage Go generics to automatically unmarshal structured JSON responses from the API directly into your Go structs.
-   **Automatic JSON Schema Generation**: Generate JSON schemas from your Go types on the fly to ensure reliable and correctly formatted JSON output from the model.
-   **Multi-modal Chat**: Easily include images in your chat completions for more context-aware interactions.
-   **Prompt Templating**: Use Go's built-in template engine to dynamically insert data into your system and user prompts.
-   **Conversation History**: Automatically maintains conversation history within a chat context.

## Installation

```sh
go get github.com/launchpad-it/openaix
```

## Configuration

### Client Configuration (Environment Variables)

The `openaix.ClientFromEnv()` function creates a client based on the following environment variables:

-   `OPENAI_API_KEY`: Your API key.
-   `OPENAI_TYPE`: The type of service. Use `openai` (default) or `azure`.
-   `OPENAI_ENDPOINT`: (Azure only) The endpoint URL for your Azure OpenAI resource.
-   `OPENAI_API_VERSION`: (Azure only) The API version to use.

### Completion Configuration (File)

OpenAIX uses [Viper](https://github.com/spf13/viper) to manage completion settings. You can define them in a YAML, JSON, or TOML file.

**`ai.yaml` Example:**

```yaml
my-chat-task:
  model: gpt-4o
  temperature: 0.7
  max_tokens: 1024
  prompts:
    system: "You are a helpful assistant. The user is a {{.Role}}. Respond in a friendly tone."
    user: "Extract the user's name and age from this text: {{.Input}}"
  json:
    name: "user_info"
    description: "Extracts user information"
    strict: true # Use response_format for guaranteed JSON
    reflect: true # Reflect the schema from the Go struct
```

## Usage

### 1. Basic Chat Completion

```go
package main

import (
	"fmt"
	"log"

	"github.com/launchpad-it/openaix"
)

func main() {
	// Create a client using environment variables
	client, err := openaix.ClientFromEnv()
	if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

	// Load completion settings from a config file
	if err := openaix.Read("config.yaml"); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Create a new chat context for string responses
	chat := openaix.Chat[string](client, "my-chat-task")

	// Define variables for prompt templating
	vars := struct {
		Role  string
		Input string
	}{
		Role:  "Developer",
		Input: "My name is John Doe and I am 30 years old.",
	}

	// Get the completion
	response, err := chat.Completion(vars)
	if err != nil {
		log.Fatalf("Chat completion failed: %v", err)
	}

	fmt.Println(response)
}
```

### 2. Structured JSON Response

Define a Go struct for the desired JSON output. OpenAIX will handle the schema generation and unmarshaling.

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/launchpad-it/openaix"
)

// Define the target struct for the JSON output
type UserInfo struct {
	Name string `json:"name" jsonschema:"description=The user's full name"`
	Age  int    `json:"age" jsonschema:"description=The user's age"`
}

func main() {
	client, err := openaix.ClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	if err := openaix.Read("config.yaml"); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Use the UserInfo struct as the generic type parameter
	chat := openaix.Chat[UserInfo](client, "my-chat-task")

	vars := struct {
		Role  string
		Input string
	}{
		Role:  "Developer",
		Input: "My name is Jane Doe and I am 25 years old.",
	}

	// The response will be an instance of UserInfo
	userInfo, err := chat.Completion(vars)
	if err != nil {
		log.Fatalf("Chat completion failed: %v", err)
	}

	// Marshal to pretty-print the JSON
	prettyJSON, _ := json.MarshalIndent(userInfo, "", "  ")
	fmt.Println(string(prettyJSON))
	// Output:
	// {
	//   "name": "Jane Doe",
	//   "age": 25
	// }
}
```

### 3. Multi-modal Chat (with Image)

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/launchpad-it/openaix"
)

func main() {
	client, err := openaix.ClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	if err := openaix.Read("config.yaml"); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Create a chat context for a simple string response
	chat := openaix.Chat[string](client, "my-chat-task")

	// Open the image file
	imageFile, err := os.Open("image.png")
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer imageFile.Close()

	// Pass the image reader to the completion method
	response, err := chat.CompletionWithImage(imageFile, openaix.Map{
		"Input": "What is in this image?",
	})
	if err != nil {
		log.Fatalf("Chat completion with image failed: %v", err)
	}

	fmt.Println(response)
}
```
