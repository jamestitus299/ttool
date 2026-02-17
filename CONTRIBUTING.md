# Contributing

Thanks for your interest in contributing to `tt`! This guide keeps contributions smooth and predictable.

## Quick Start

1. Fork the repo and create a feature branch.
2. Make your changes with tests where appropriate.
3. Run `go fmt` on any touched Go files.
4. Open a PR with a clear description of what changed and why.

## Development

Build:
```bash
go build -o tt main.go
```

Run:
```bash
./tt --help
```

## Coding Guidelines

- Keep changes focused and minimal.
- Prefer small, composable functions.
- Avoid adding heavy dependencies unless necessary.
- Update README or `.env.example` when adding configuration.

## Submitting a PR

Please include:
- A short summary of changes
- Testing notes (what you ran)
- Any follow-ups or known limitations

## Reporting Issues

When filing a bug, include:
- OS and shell
- Steps to reproduce
- Expected vs actual behavior
- Logs or output if available
