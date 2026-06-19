#!/usr/bin/env sh
# Adds a `tt` alias pointing at the built binary to your shell startup file.
# Supports zsh, bash, and fish on macOS and Linux.
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPO_ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
BIN="$REPO_ROOT/tt"

# Build the binary if it isn't there yet.
if [ ! -x "$BIN" ]; then
	echo "Building tt..."
	(cd "$REPO_ROOT" && go build -o tt main.go)
fi

# Pick the startup file for the current shell.
shell_name=$(basename "${SHELL:-sh}")
case "$shell_name" in
	zsh)
		RC="$HOME/.zshrc"
		;;
	bash)
		if [ "$(uname)" = "Darwin" ]; then
			RC="$HOME/.bash_profile"
		else
			RC="$HOME/.bashrc"
		fi
		;;
	fish)
		RC="$HOME/.config/fish/config.fish"
		;;
	*)
		RC="$HOME/.profile"
		;;
esac

ALIAS_LINE="alias tt='$BIN'"

mkdir -p "$(dirname "$RC")"
touch "$RC"

# Idempotent: refresh an existing tt alias instead of appending a duplicate.
if grep -q "alias tt=" "$RC" 2>/dev/null; then
	tmp=$(mktemp)
	grep -v "alias tt=" "$RC" >"$tmp"
	mv "$tmp" "$RC"
	echo "Updated existing tt alias in $RC"
else
	printf '\n# Added by ttool install script\n' >>"$RC"
	echo "Added tt alias to $RC"
fi
printf '%s\n' "$ALIAS_LINE" >>"$RC"

echo "Done. Restart your shell or run: source \"$RC\""
