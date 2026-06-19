package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const defaultAnthropicModel = "claude-opus-4-8"

type AnthropicClient struct {
	client anthropic.Client
	model  string
}

func NewAnthropicClient(apiKey, model string) (*AnthropicClient, error) {
	apiKey = strings.TrimSpace(apiKey)
	model = strings.TrimSpace(model)
	if apiKey == "" {
		return nil, fmt.Errorf("Anthropic API key is empty")
	}
	if model == "" {
		model = defaultAnthropicModel
	}

	return &AnthropicClient{
		client: anthropic.NewClient(option.WithAPIKey(apiKey)),
		model:  model,
	}, nil
}

func (c *AnthropicClient) GetCommand(request, osSystem, shell string) (CommandResponse, error) {
	ctx := context.Background()

	resp, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: buildSystemPrompt()},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(buildUserPrompt(request, osSystem, shell))),
		},
	})
	if err != nil {
		return CommandResponse{}, fmt.Errorf("Anthropic request failed: %v", err)
	}

	var content strings.Builder
	for _, block := range resp.Content {
		if text, ok := block.AsAny().(anthropic.TextBlock); ok {
			content.WriteString(text.Text)
		}
	}

	text := extractJSON(strings.TrimSpace(content.String()))
	if text == "" {
		return CommandResponse{}, fmt.Errorf("Anthropic returned empty output")
	}

	var cmdResp CommandResponse
	if err := json.Unmarshal([]byte(text), &cmdResp); err != nil {
		return CommandResponse{}, fmt.Errorf("failed to parse Anthropic structured output: %v", err)
	}

	return cmdResp, nil
}
