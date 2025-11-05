#!/bin/bash
# Description: Removes files/directories from git tracking both locally and remotely
# Usage: ./gitremtracking.sh file1.txt dir1/ file2.txt
# This script will:
# 1. Remove files/directories from local git tracking
# 2. Commit the changes locally
# 3. Remove files/directories from remote tracking
# 4. Push changes to remote repository

# Check if files/directories are provided as arguments
if [ $# -eq 0 ]; then
  echo "Usage: $0 <file1|dir1> [file2|dir2 ...]"
  exit 1
fi

# Process each file/directory provided as argument
for item in "$@"; do
  # Check if item exists
  if [ ! -e "$item" ] && ! git ls-files --error-unmatch "$item" &>/dev/null; then
    echo "Warning: '$item' does not exist or is not tracked by git. Skipping..."
    continue
  fi
  
  # Determine if it's a file or directory
  if [ -d "$item" ]; then
    echo "Processing directory: $item"
    TYPE="directory"
    # Remove directory recursively from git tracking
    git rm -r --cached "$item"
  else
    echo "Processing file: $item"
    TYPE="file"
    # Remove file from git tracking
    git rm --cached "$item"
  fi
  
  # Check if git rm was successful
  if [ $? -ne 0 ]; then
    echo "Error: Failed to remove '$item' from git tracking. Skipping..."
    continue
  fi
  
  # Step 2: Update .gitignore file (commented out)
  #echo "$item" >> .gitignore
  
  # Step 3: Commit and push changes
  git commit -m "Stop tracking $TYPE '$item' in repository"
  
  if [ $? -eq 0 ]; then
    git push origin "$(git branch --show-current)"
    
    if [ $? -eq 0 ]; then
      echo "✓ $TYPE '$item' has been removed from both local and remote tracking."
    else
      echo "✗ Error: Failed to push changes for '$item' to remote repository."
    fi
  else
    echo "✗ Error: Failed to commit changes for '$item'."
  fi
  
  echo "---"
done

echo "Process completed!"