#!/bin/bash
# Description: Utility functions for GitHub repository management and user information retrieval
# Usage: source githubutils.sh && get_github_user

# Add this function to get the current GitHub user
get_github_user() {
    local user
    user=$(gh api user -q '.login' 2>/dev/null)
    if [ -z "$user" ]; then
        error "Failed to determine GitHub username"
        return 1
    fi
    echo "$user"
}