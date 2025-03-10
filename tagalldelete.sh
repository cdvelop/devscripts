#!/bin/bash
# Bulk delete git tags listed in a text file
# Usage: ./tagalldelete.sh <filename>
# The file should contain one tag name per line

# Read the text file containing tags to delete
while read tag; do
  # Delete tag locally
  git tag -d "$tag"
  # Delete tag from remote repository
  git push origin --delete "$tag"
done < "$1"
