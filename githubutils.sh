#!/bin/bash
# Description: Utility functions for GitHub repository management and user information retrieval
# Usage: source githubutils.sh

# expected eg: juanin654
gitHubOwner=$(gh api user --jq .login)

# Function to ensure .github directory exists and is hidden on Windows
ensure_github_directory() {
    local github_dir=".github"
    
    # Create .github directory if it doesn't exist  
    if [ ! -d "$github_dir" ]; then
        mkdir -p "$github_dir"
        
        # On Windows, make the directory hidden using attrib command
        if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || -n "$MSYSTEM" ]]; then
            attrib +h "$github_dir" 2>/dev/null || true
        fi
    else
        # Ensure it's hidden on Windows if it already exists
        if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || -n "$MSYSTEM" ]]; then
            attrib +h "$github_dir" 2>/dev/null || true
        fi
    fi
    
    echo "$github_dir"
}


