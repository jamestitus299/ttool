package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultOpenAIModel = "gpt-4o-mini"

type OpenAIClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

type openAIInputMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	OutputText string `json:"output_text"`
	Output     []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func NewOpenAIClient(apiKey, model string) (*OpenAIClient, error) {
	apiKey = strings.TrimSpace(apiKey)
	model = strings.TrimSpace(model)
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is empty")
	}
	if model == "" {
		model = defaultOpenAIModel
	}

	return &OpenAIClient{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *OpenAIClient) GetCommand(request, osSystem, shell string) (CommandResponse, error) {
	ctx := context.Background()

	reqBody := map[string]any{
		"model":        c.model,
		"instructions": buildSystemPrompt(),
		"input": []openAIInputMessage{
			{
				Role:    "user",
				Content: buildUserPrompt(request, osSystem, shell),
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return CommandResponse{}, fmt.Errorf("failed to marshal OpenAI request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(bodyBytes))
	if err != nil {
		return CommandResponse{}, fmt.Errorf("failed to create OpenAI request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return CommandResponse{}, fmt.Errorf("OpenAI request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return CommandResponse{}, fmt.Errorf("failed to read OpenAI response: %v", err)
	}
	if resp.StatusCode >= 400 {
		return CommandResponse{}, fmt.Errorf("OpenAI request failed with status %s", resp.Status)
	}

	var out openAIResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return CommandResponse{}, fmt.Errorf("failed to parse OpenAI response: %v", err)
	}
	if out.Error != nil && out.Error.Message != "" {
		return CommandResponse{}, fmt.Errorf("OpenAI error: %s", out.Error.Message)
	}

	content := strings.TrimSpace(out.OutputText)
	if content == "" {
		content = strings.TrimSpace(firstOutputText(out.Output))
	}
	if content == "" {
		return CommandResponse{}, fmt.Errorf("OpenAI returned empty output")
	}

	content = extractJSON(content)

	var cmdResp CommandResponse
	if err := json.Unmarshal([]byte(content), &cmdResp); err != nil {
		return CommandResponse{}, fmt.Errorf("failed to parse OpenAI structured output: %v", err)
	}

	return cmdResp, nil
}

func firstOutputText(outputs []struct {
	Type    string `json:"type"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}) string {
	for _, out := range outputs {
		if out.Type != "message" {
			continue
		}
		for _, content := range out.Content {
			if content.Type == "output_text" && strings.TrimSpace(content.Text) != "" {
				return content.Text
			}
		}
	}
	return ""
}

func extractJSON(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```") {
		if idx := strings.Index(content, "\n"); idx != -1 {
			content = content[idx+1:]
		}
		if idx := strings.LastIndex(content, "```"); idx != -1 {
			content = content[:idx]
		}
	}
	return strings.TrimSpace(content)
}
