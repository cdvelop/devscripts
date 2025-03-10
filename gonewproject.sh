#!/bin/bash
# Description: Creates a new Go project with standard directory structure and initial files, sets up remote repository
# Usage: ./gonewproject.sh <repo-name> <description> [visibility]

source functions.sh

init_project() {
    local repo_name=$1
    local description=$2
    local visibility=${3:-public}  # Default to public if not specified

    # Validate required arguments
    if [ -z "$repo_name" ] || [ -z "$description" ]; then
        error "Usage: init_project <repo-name> <description> [visibility]"
        return 1
    fi

    # create repository remote 
    if ! repocreate.sh "$repo_name" "$description" "$visibility"; then
        error "Failed to create remote repository"
        return 1
    fi

    # Change to the new repository directory
    cd "$repo_name" || return 1

    # Setup basic go project structure
    if ! gomodinit.sh; then
        error "Failed to initialize go project"
        return 1
    fi

    # Setup repoexistingsetup
    if ! repoexistingsetup.sh; then
        error "Failed to setup existing repo"
        return 1
    fi

    return 0
}

# Execute if script is run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    if [ "$#" -lt 2 ] || [ "$#" -gt 3 ]; then
        error "Usage: $0 <repo-name> <description> [visibility]"
        exit 1
    fi

    init_project "$1" "$2" "$3"
    exit_code=$?
    successMessages
    exit $exit_code
fi