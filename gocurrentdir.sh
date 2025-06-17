#!/bin/bash
# Description: Generic Go runner that executes cmd/{script_name}.go with current directory context
# Usage: Called from other scripts via: source gocurrentdir.sh

# Get the name of the calling script (without .sh extension)
CALLING_SCRIPT=$(basename "${BASH_SOURCE[1]}" .sh)

# Set up directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[1]}")" && pwd)"
CURRENT_DIR="$(pwd)"

# Change to script directory and run the corresponding Go command
cd "$SCRIPT_DIR"
go run "cmd/${CALLING_SCRIPT}.go" "$CURRENT_DIR" "$@"
