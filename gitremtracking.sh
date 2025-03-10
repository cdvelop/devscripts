#!/bin/bash
# Description: Removes files from git tracking both locally and remotely
# Usage: ./gitremtracking.sh file1.txt file2.txt
# This script will:
# 1. Remove files from local git tracking
# 2. Commit the changes locally
# 3. Remove files from remote tracking
# 4. Push changes to remote repository

# Check if files are provided as arguments
if [ $# -eq 0 ]; then
  echo "Usage: $0 <file1> [file2 ...]"
  exit 1
fi

# Process each file provided as argument
for file in "$@"; do
  # Step 1: Stop tracking the file locally
  git rm --cached "$file"
  
  # Step 2: Update .gitignore file (commented out)
  #echo "$file" >> .gitignore
  
  # Step 3: Commit changes locally
  git commit -m "Stop tracking file $file locally"
  
  # Step 4: Remove file from remote tracking
  git rm --cached "$file"
  
  # Step 5: Commit and push changes to remote repository
  git commit -m "Stop tracking file $file in remote repository"
  git push origin "$(git branch --show-current)"
  
  echo "File $file has been removed from both local and remote tracking."
done