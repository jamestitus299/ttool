package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	GeminiAPIKey              string `json:"gemini_api_key"`
	GeminiModel               string `json:"gemini_model"`
	Provider                  string `json:"provider"`
	OpenAIAPIKey              string `json:"openai_api_key"`
	OpenAIModel               string `json:"openai_model"`
	AnthropicAPIKey           string `json:"anthropic_api_key"`
	AnthropicModel            string `json:"anthropic_model"`
	AzureOpenAIEndpoint       string `json:"azure_openai_endpoint"`
	AzureOpenAIAPIKey         string `json:"azure_openai_api_key"`
	AzureOpenAIDeploymentName string `json:"azure_openai_deployment_name"`
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "ttool", "config.json"), nil
}

func Load() (*Config, error) {
	cfg := Config{}

	// Load config file first
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
	} else {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	}

	// Environment variables override config file
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		cfg.GeminiAPIKey = key
	}
	if model := os.Getenv("GEMINI_MODEL"); model != "" {
		cfg.GeminiModel = model
	}
	if provider := os.Getenv("LLM_PROVIDER"); provider != "" {
		cfg.Provider = strings.ToLower(strings.TrimSpace(provider))
	}
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		cfg.OpenAIAPIKey = key
	}
	if model := os.Getenv("OPENAI_MODEL"); model != "" {
		cfg.OpenAIModel = model
	}
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		cfg.AnthropicAPIKey = key
	}
	if model := os.Getenv("ANTHROPIC_MODEL"); model != "" {
		cfg.AnthropicModel = model
	}
	if endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT"); endpoint != "" {
		cfg.AzureOpenAIEndpoint = endpoint
	}
	if key := os.Getenv("AZURE_OPENAI_API_KEY"); key != "" {
		cfg.AzureOpenAIAPIKey = key
	}
	if deployment := os.Getenv("AZURE_OPENAI_DEPLOYMENT_NAME"); deployment != "" {
		cfg.AzureOpenAIDeploymentName = deployment
	}

	return &cfg, nil
}

func Save(key string) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	cfg := Config{}
	if data, err := ioutil.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &cfg)
	}
	cfg.GeminiAPIKey = key

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0600)
}

func SaveConfig(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0600)
}

// Clear removes the saved config file. It is not an error if the file does
// not exist.
func Clear() error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
