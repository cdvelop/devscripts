#!/bin/bash
# Setup additional files and tags for an existing Git repository
# Usage: ./repoexistingsetup.sh
# This script will:
# 1. Check if repository exists
# 2. Create changes file
# 3. Create initial tag
source functions.sh
source gitutils.sh
source doingmdfile.sh

setup_existing_project() {
    local current_folder=$(basename "$(pwd)")
    
    # Check if git repository already exists
    if [ ! -d ".git" ]; then
        error "Not a git repository. Please initialize git first."
        return 1
    fi
    create_doing_md_file || return $?
    create_git_tag || return $?
    return 0
}

# Execute directly if script is not being sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    setup_existing_project
    exit_code=$?
    successMessages
    exit $exit_code
fi