#!/bin/bash

# Script to install Git hooks from the .husky directory to .git/hooks

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Define source and destination directories
HOOKS_SOURCE_DIR="$PROJECT_ROOT/.husky"
HOOKS_DEST_DIR="$PROJECT_ROOT/.git/hooks"

# Check if the .husky directory exists
if [ ! -d "$HOOKS_SOURCE_DIR" ]; then
    echo "Error: .husky directory not found at $HOOKS_SOURCE_DIR"
    exit 1
fi

# Check if the .git/hooks directory exists
if [ ! -d "$HOOKS_DEST_DIR" ]; then
    echo "Error: .git/hooks directory not found at $HOOKS_DEST_DIR"
    echo "Make sure you're running this script from within a Git repository"
    exit 1
fi

# Copy all hooks from the source to the destination
echo "Installing Git hooks..."
for hook in "$HOOKS_SOURCE_DIR"/*; do
    if [ -f "$hook" ]; then
        hook_name=$(basename "$hook")
        cp "$hook" "$HOOKS_DEST_DIR/$hook_name"
        chmod +x "$HOOKS_DEST_DIR/$hook_name"
        echo "Installed hook: $hook_name"
    fi
done

echo "Git hooks installation completed!"
echo "The following hooks are now active:"
echo "- commit-msg: Validates commit message format"
echo "- pre-commit: Automatically runs 'go mod tidy' before commits"