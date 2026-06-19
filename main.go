package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jamestitus/ttool/internal/config"
	"github.com/jamestitus/ttool/internal/env"
	"github.com/jamestitus/ttool/internal/llm"
	"github.com/jamestitus/ttool/internal/term"
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
)

func label(text string) string {
	return colorBold + colorCyan + text + colorReset
}

func prompt(text string) string {
	return colorBold + text + colorReset
}

func info(text string) string {
	return colorDim + text + colorReset
}

func spinner(stop <-chan struct{}) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-stop:
			fmt.Print("\r\033[K")
			return
		default:
			fmt.Printf("\r%s", colorYellow+frames[i%len(frames)]+colorReset)
			time.Sleep(90 * time.Millisecond)
			i++
		}
	}
}

func box(title string, lines []string) {
	maxLen := len(title)
	if maxLen != 0 {
		for _, line := range lines {
			if len(line) > maxLen {
				maxLen = len(line)
			}
		}
		width := maxLen + 2
		top := "┌" + strings.Repeat("─", width) + "┐"
		mid := "│ " + title + strings.Repeat(" ", width-len(title)-1) + "│"
		bot := "└" + strings.Repeat("─", width) + "┘"
		fmt.Println(top)
		fmt.Println(mid)
		fmt.Println(bot)
	}

	for _, line := range lines {
		fmt.Println("  " + line)
	}
}

func runConfigWizard() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println(label("Configuration"))
	fmt.Println(info("Select the provider and enter credentials."))
	fmt.Println("  1) Gemini")
	fmt.Println("  2) Azure OpenAI")
	fmt.Println("  3) OpenAI")
	fmt.Println("  4) Anthropic")
	fmt.Print(prompt("Enter choice (1/2/3/4): "))
	var choice string
	fmt.Scanln(&choice)
	choice = strings.TrimSpace(choice)

	var provider string
	switch choice {
	case "2":
		provider = "azure"
	case "3":
		provider = "openai"
	case "4":
		provider = "anthropic"
	default:
		provider = "gemini"
	}

	switch provider {
	case "gemini":
		fmt.Print(prompt("Enter your Gemini API Key: "))
		var key string
		fmt.Scanln(&key)
		fmt.Print(prompt("Enter Gemini Model (optional, default gemini-flash-latest): "))
		var model string
		fmt.Scanln(&model)
		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("API key cannot be empty")
		}
		cfg.GeminiAPIKey = key
		cfg.GeminiModel = strings.TrimSpace(model)
	case "azure":
		fmt.Print(prompt("Enter Azure OpenAI Endpoint (e.g. https://your-resource.services.ai.azure.com/openai/v1/): "))
		var endpoint string
		fmt.Scanln(&endpoint)
		fmt.Print(prompt("Enter Azure OpenAI API Key: "))
		var key string
		fmt.Scanln(&key)
		fmt.Print(prompt("Enter Azure OpenAI Deployment Name: "))
		var deployment string
		fmt.Scanln(&deployment)

		cfg.AzureOpenAIEndpoint = strings.TrimSpace(endpoint)
		cfg.AzureOpenAIAPIKey = strings.TrimSpace(key)
		cfg.AzureOpenAIDeploymentName = strings.TrimSpace(deployment)
	case "openai":
		fmt.Print(prompt("Enter OpenAI API Key: "))
		var key string
		fmt.Scanln(&key)
		fmt.Print(prompt("Enter OpenAI Model (optional): "))
		var model string
		fmt.Scanln(&model)
		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("API key cannot be empty")
		}
		cfg.OpenAIAPIKey = key
		cfg.OpenAIModel = strings.TrimSpace(model)
	case "anthropic":
		fmt.Print(prompt("Enter Anthropic API Key: "))
		var key string
		fmt.Scanln(&key)
		fmt.Print(prompt("Enter Anthropic Model (optional, default claude-opus-4-8): "))
		var model string
		fmt.Scanln(&model)
		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("API key cannot be empty")
		}
		cfg.AnthropicAPIKey = key
		cfg.AnthropicModel = strings.TrimSpace(model)
	}

	cfg.Provider = provider
	return config.SaveConfig(cfg)
}

func main() {
	configureWizard := flag.Bool("config", false, "Configure provider and credentials")
	clearConfig := flag.Bool("clear-config", false, "Clear saved configuration and credentials")
	showHelp := flag.Bool("help", false, "Show help")
	flag.BoolVar(showHelp, "h", false, "Show help")
	flag.Parse()

	if *showHelp {
		fmt.Println(label("Usage"))
		fmt.Println("  tt <request>          Generate a command from natural language")
		fmt.Println("  tt --config           Configure provider and credentials")
		fmt.Println("  tt --clear-config     Clear saved configuration and credentials")
		fmt.Println("  tt --help             Show help")
		os.Exit(0)
	}

	if *clearConfig {
		if err := config.Clear(); err != nil {
			fmt.Fprintf(os.Stderr, "Error clearing config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(colorGreen + "Configuration cleared." + colorReset)
		os.Exit(0)
	}

	if *configureWizard {
		if err := runConfigWizard(); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(colorGreen + "Configuration saved." + colorReset)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println(label("Usage"))
		fmt.Println("  tt <request>          Generate a command from natural language")
		fmt.Println("  tt --config           Configure provider and credentials")
		fmt.Println("  tt --clear-config     Clear saved configuration and credentials")
		fmt.Println("  tt --help             Show help")
		os.Exit(1)
	}

	request := strings.Join(args, " ")

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))
	if provider == "" && cfg.GeminiAPIKey == "" && cfg.AzureOpenAIAPIKey == "" && cfg.OpenAIAPIKey == "" && cfg.AnthropicAPIKey == "" {
		fmt.Println(label("First-time setup"))
		if err := runConfigWizard(); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(colorGreen + "Configuration saved." + colorReset)
		cfg, err = config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		provider = strings.ToLower(strings.TrimSpace(cfg.Provider))
	}

	if provider == "" {
		switch {
		case cfg.AzureOpenAIAPIKey != "" && cfg.AzureOpenAIEndpoint != "" && cfg.AzureOpenAIDeploymentName != "":
			provider = "azure"
		case cfg.OpenAIAPIKey != "":
			provider = "openai"
		case cfg.AnthropicAPIKey != "":
			provider = "anthropic"
		default:
			provider = "gemini"
		}
	}

	e := env.Detect()
	var client llm.Client
	switch provider {
	case "azure":
		client, err = llm.NewAzureOpenAIClient(
			cfg.AzureOpenAIEndpoint,
			cfg.AzureOpenAIAPIKey,
			cfg.AzureOpenAIDeploymentName,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "openai":
		if cfg.OpenAIAPIKey == "" {
			fmt.Println(colorYellow + "OpenAI API key not found." + colorReset)
			fmt.Print(prompt("Please enter your OpenAI API Key: "))
			var key string
			fmt.Scanln(&key)
			key = strings.TrimSpace(key)
			if key == "" {
				fmt.Fprintln(os.Stderr, "Error: API key cannot be empty.")
				os.Exit(1)
			}
			cfg.OpenAIAPIKey = key
			if err := config.SaveConfig(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}
		}
		client, err = llm.NewOpenAIClient(
			cfg.OpenAIAPIKey,
			cfg.OpenAIModel,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "gemini":
		if cfg.GeminiAPIKey == "" {
			fmt.Println(colorYellow + "Gemini API key not found." + colorReset)
			fmt.Print(prompt("Please enter your Gemini API Key: "))
			var key string
			fmt.Scanln(&key)
			key = strings.TrimSpace(key)
			if key == "" {
				fmt.Fprintln(os.Stderr, "Error: API key cannot be empty.")
				os.Exit(1)
			}
			if err := config.Save(key); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("API Key saved successfully.")
			cfg.GeminiAPIKey = key
		}
		client, err = llm.NewClient(cfg.GeminiAPIKey, cfg.GeminiModel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "anthropic":
		if cfg.AnthropicAPIKey == "" {
			fmt.Println(colorYellow + "Anthropic API key not found." + colorReset)
			fmt.Print(prompt("Please enter your Anthropic API Key: "))
			var key string
			fmt.Scanln(&key)
			key = strings.TrimSpace(key)
			if key == "" {
				fmt.Fprintln(os.Stderr, "Error: API key cannot be empty.")
				os.Exit(1)
			}
			cfg.AnthropicAPIKey = key
			if err := config.SaveConfig(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}
		}
		client, err = llm.NewAnthropicClient(
			cfg.AnthropicAPIKey,
			cfg.AnthropicModel,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown LLM provider %q (use \"gemini\", \"azure\", \"openai\", or \"anthropic\")\n", provider)
		os.Exit(1)
	}

	stop := make(chan struct{})
	go spinner(stop)
	resp, err := client.GetCommand(request, e.OS, e.Shell)
	close(stop)
	fmt.Print("\r\033[K")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()
	box("", []string{
		// fmt.Sprintf("%s %s", label("Provider:"), provider),
		fmt.Sprintf("%s %s", label("Command:"), colorYellow+resp.Command+colorReset),
		fmt.Sprintf("%s %s", label("Purpose:"), colorBlue+resp.Purpose+colorReset),
	})

	// Inject the command into the terminal input buffer
	if err := term.Inject(resp.Command); err != nil {
		fmt.Print("An error occured. Could not set command.")
	}
}
