#!/bin/bash
# Script to change the remote URL of a Git repository
# Usage: ./changeremote.sh https://github.com/username/repository.git

# Get current remote URL
remote_info=$(git remote -v)

# Print current remote URL
echo "Current remote URL:"
echo "$remote_info"

# Check if new origin was provided as argument
if [ -z "$1" ]; then
  echo "ERROR: No new origin provided as argument."
  echo "e.g.: change-remote.sh https://github.com/your-username/your-repository.git"
  exit 1
fi

# Get new origin from argument
new_origin=$1

# Change remote URL
git remote set-url origin "$new_origin"

# Print updated remote URL
echo "Updated remote URL:"
git remote -v
