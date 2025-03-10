#!/bin/bash
# Description: Script to automatically increment the last number in a Git tag
# Usage: ./tag.sh (will get the latest tag and suggest the next one)

# Get the latest tag
latest_tag=$(git describe --abbrev=0 --tags)

# Extract the number from the tag
last_number=$(echo "$latest_tag" | grep -oE '[0-9]+$')

# Increment the number by one
next_number=$((last_number + 1))

# Build the new tag
new_tag=$(echo "$latest_tag" | sed "s/$last_number$/$next_number/")

echo "Current tag is: $latest_tag"
echo "Next tag will be: $new_tag"