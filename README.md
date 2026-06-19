# tt (Terminal Tool)

`tt` is a CLI tool that converts natural language requests into terminal commands using LLM.

### Build From Source

1. Clone the repository and build the binary:
   ```bash
   go build -o tt main.go
   ```
2. Move the binary to your PATH (e.g., `/usr/local/bin`) or create an alias:
   ```bash
   alias tt='/path/to/tt'
   ```
3. Set your provider credentials:
   ```bash
   export GEMINI_API_KEY='your-api-key-here'
   export GEMINI_MODEL='gemini-flash-latest'  # optional; defaults to gemini-flash-latest
   ```
   Or use OpenAI:
   ```bash
   export LLM_PROVIDER='openai'
   export OPENAI_API_KEY='your-api-key-here'
   export OPENAI_MODEL='gpt-4o-mini'
   ```
   Or use Azure OpenAI:
   ```bash
   export LLM_PROVIDER='azure'
   export AZURE_OPENAI_ENDPOINT='https://your-resource.openai.azure.com/openai/v1/'
   export AZURE_OPENAI_API_KEY='your-api-key-here'
   export AZURE_OPENAI_DEPLOYMENT_NAME='your-deployment-name'
   ```
   Or use Anthropic:
   ```bash
   export LLM_PROVIDER='anthropic'
   export ANTHROPIC_API_KEY='your-api-key-here'
   export ANTHROPIC_MODEL='claude-opus-4-8'  # optional; defaults to claude-opus-4-8
   ```

Credentials can also be saved interactively with `tt --config`. They are stored
in `~/.config/ttool/config.json` (mode `0600`); environment variables override the
saved file. Run `tt --clear-config` to delete the saved configuration.

## Usage

Simply type `tt` followed by your request:

```bash
tt i want to see all files in this directory
# Output: ls -la

tt find all logs modified in the last 24 hours
# Output: find . -name "*.log" -mmin -1440
```

### Help

```bash
tt --help
```

```
Usage
  tt <request>          Generate a command from natural language
  tt --config           Configure provider and credentials
  tt --clear-config     Clear saved configuration and credentials
  tt --help             Show help
```
