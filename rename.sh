#!/bin/bash
# Rename a file and update Git tracking
# Usage: ./rename.sh <current_name> <new_name>

# Check if both parameters are provided
if [ $# -ne 2 ]; then
  echo "Usage: $0 <current_name> <new_name>"
  exit 1
fi

# Assign parameters to variables
current_name=$1
new_name=$2

# Rename file locally
mv "$current_name" "$new_name"

# Track the new name in Git
git add "$new_name"

# Commit the changes
git commit -m "Rename file from $current_name to $new_name"

# Print success message
echo "File successfully renamed to $new_name and changes committed to Git."
