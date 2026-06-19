package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

type Client interface {
	GetCommand(request, osSystem, shell string) (CommandResponse, error)
}

const defaultGeminiModel = "gemini-flash-latest"

type GeminiClient struct {
	client *genai.Client
	model  string
}

type CommandResponse struct {
	Command string `json:"command"`
	Purpose string `json:"purpose"`
}

func buildSystemPrompt() string {
	return "You are an expert at terminal commands. Respond with a JSON object containing the \"command\" and a brief \"purpose\" of the command."
}

func buildUserPrompt(request, osSystem, shell string) string {
	return fmt.Sprintf(`The user wants to: %q
Operating System: %s
Shell: %s`, request, osSystem, shell)
}

func buildGeminiPrompt(request, osSystem, shell string) string {
	return fmt.Sprintf(`%s
%s`, buildSystemPrompt(), buildUserPrompt(request, osSystem, shell))
}

func NewClient(apiKey, model string) (*GeminiClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key is empty")
	}
	model = strings.TrimSpace(model)
	if model == "" {
		model = defaultGeminiModel
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}

	return &GeminiClient{client: client, model: model}, nil
}

func (c *GeminiClient) GetCommand(request, osSystem, shell string) (CommandResponse, error) {
	prompt := buildGeminiPrompt(request, osSystem, shell)

	ctx := context.Background()
	result, err := c.client.Models.GenerateContent(
		ctx,
		c.model,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			Temperature:      genai.Ptr(float32(0.0)),
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"command": {Type: genai.TypeString},
					"purpose": {Type: genai.TypeString},
				},
				Required: []string{"command", "purpose"},
			},
		},
	)
	if err != nil {
		return CommandResponse{}, fmt.Errorf("failed to generate command: %v", err)
	}

	if len(result.Candidates) == 0 || result.Candidates[0].Content == nil || len(result.Candidates[0].Content.Parts) == 0 {
		return CommandResponse{}, fmt.Errorf("no response candidates from Gemini")
	}

	var cmdResp CommandResponse
	err = json.Unmarshal([]byte(result.Text()), &cmdResp)
	if err != nil {
		return CommandResponse{}, fmt.Errorf("failed to parse structured output: %v", err)
	}

	return cmdResp, nil
}
