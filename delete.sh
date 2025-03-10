#!/bin/bash
# Description: Script to delete a file locally and track the deletion in Git
# Usage: ./delete.sh filename.txt

# Check if parameter is provided
if [ $# -ne 1 ]; then
  echo "Usage: $0 <filename>"
  exit 1
fi

# Assign filename
filename=$1

# Delete file locally
rm "$filename"

# Track the deletion in Git
git rm "$filename"

# Commit changes
git commit -m "Delete file $filename"

# Print success message
echo "File $filename successfully deleted and changes committed to Git."