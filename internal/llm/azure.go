package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type AzureOpenAIClient struct {
	client     openai.Client
	deployment string
}

func NewAzureOpenAIClient(endpoint, apiKey, deployment string) (*AzureOpenAIClient, error) {
	endpoint = strings.TrimSpace(endpoint)
	apiKey = strings.TrimSpace(apiKey)
	deployment = strings.TrimSpace(deployment)

	switch {
	case endpoint == "":
		return nil, fmt.Errorf("Azure OpenAI endpoint is empty")
	case apiKey == "":
		return nil, fmt.Errorf("Azure OpenAI API key is empty")
	case deployment == "":
		return nil, fmt.Errorf("Azure OpenAI deployment name is empty")
	}

	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}

	return &AzureOpenAIClient{
		client: openai.NewClient(
			option.WithBaseURL(endpoint),
			option.WithAPIKey(apiKey),
		),
		deployment: deployment,
	}, nil
}

func (c *AzureOpenAIClient) GetCommand(request, osSystem, shell string) (CommandResponse, error) {
	ctx := context.Background()

	prompt := fmt.Sprintf("%s\n%s", buildSystemPrompt(), buildUserPrompt(request, osSystem, shell))
	resp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(prompt)},
		Model: c.deployment,
	})
	if err != nil {
		return CommandResponse{}, fmt.Errorf("Azure OpenAI request failed: %v", err)
	}

	content := strings.TrimSpace(resp.OutputText())
	if content == "" {
		return CommandResponse{}, fmt.Errorf("Azure OpenAI returned empty output")
	}
	content = extractJSON(content)

	var cmdResp CommandResponse
	if err := json.Unmarshal([]byte(content), &cmdResp); err != nil {
		return CommandResponse{}, fmt.Errorf("failed to parse Azure OpenAI structured output: %v", err)
	}

	return cmdResp, nil
}
