#!/bin/bash
# Delete git tags locally and remotely
# Usage: tagdelete.sh tag1 tag2 tag3

# Receive tags as space-separated arguments
tags="$@"

for tag in $tags; do
  # Delete tag locally
  git tag -d "$tag"
  
  # Delete tag in remote repository
  git push origin --delete "$tag"
done