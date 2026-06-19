# tt (Terminal Tool)

`tt` is a CLI tool that converts natural language requests into terminal commands using LLM.

### Build From Source

1. Clone the repository, then run the install script. It builds the binary and
   adds a `tt` alias to the right shell startup file for your shell
   (`~/.zshrc`, `~/.bashrc`/`~/.bash_profile`, or fish config):
   ```bash
   ./scripts/install.sh
   ```
   On Windows (PowerShell), run instead:
   ```powershell
   .\scripts\install.ps1
   ```
   Restart your shell (or `source` the startup file) for the alias to take effect.
2. Configure your provider and credentials:
   ```bash
   tt --config
   ```
   This walks you through selecting a provider (Gemini, OpenAI, Azure OpenAI, or
   Anthropic) and entering your API key and an optional model. Credentials are
   stored in `~/.config/ttool/config.json` (mode `0600`). Run `tt --clear-config`
   to delete the saved configuration.

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
