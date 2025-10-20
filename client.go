package openaix

import (
	"errors"
	"os"

	"github.com/sashabaranov/go-openai"
)

// ClientFromEnv creates an OpenAI client based on environment variables.
// It's kept here to unify the OpenAI client initialization among different projects.
func ClientFromEnv() (*openai.Client, error) {
	var (
		kind     = os.Getenv("OPENAI_TYPE")
		endpoint = os.Getenv("OPENAI_ENDPOINT")
		version  = os.Getenv("OPENAI_API_VERSION")
		key      = os.Getenv("OPENAI_API_KEY")
	)

	switch kind {
	case "azure":
		config := openai.DefaultAzureConfig(key, endpoint)
		config.APIVersion = version
		return openai.NewClientWithConfig(config), nil
	case "openai":
		return openai.NewClient(key), nil
	}

	return nil, errors.New("openaix: unknown OPENAI_TYPE: " + kind)
}
