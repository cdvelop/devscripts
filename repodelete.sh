#!/bin/bash
# Description: Deletes a remote GitHub repository after confirmation and permission checks
# Usage: ./reporemotedelete.sh <repo-name> [force_delete] [owner]

source functions.sh
source githubutils.sh

check_delete_permissions() {
    # First check if gh CLI is authenticated at all
    if ! gh auth status >/dev/null 2>&1; then
        warning "GitHub CLI not authenticated. Please run: gh auth login"
        return 1
    fi
    
    # Check if we have delete_repo scope
    if ! gh auth status 2>&1 | grep -q "repo"; then
        warning "Requesting repository permissions..."
        
    # Request full repo scope for private repository access
    if ! gh auth refresh -h github.com -s repo; then
            error "Failed to obtain repo permission"
            return 1
        fi
        
        # Verify we got the permission
        if ! gh auth status 2>&1 | grep -q "repo"; then
            error "Still missing repo permission after refresh"
            return 1
        fi
        
        success "Delete permission granted successfully"
    fi
    return 0
}

delete_repository() {
    local repo_name=$1
    local force_delete=${2:-false}  # Default to non-force delete
    local repo_owner=${3:-"$(get_github_user)"}  # Use provided owner or default to current user

    # Validate required arguments
    if [ -z "$repo_name" ]; then
        error "Usage: delete_repository <repo-name> [true] [owner]"
        return 1
    fi

    # Check permissions first
    check_delete_permissions || return $?

    # Confirm deletion unless force flag is set
    if [ "$force_delete" != "true" ]; then
        read -p "Are you sure you want to delete repository '$repo_owner/$repo_name'? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            warning "Operation cancelled by user"
            return 1
        fi
    fi

    # Delete repository using --yes instead of --confirm
    execute "gh repo delete $repo_owner/$repo_name --yes" \
        "Failed to delete repository" \
        "Repository $repo_owner/$repo_name deleted successfully" || return $?

    return 0
}

# Execute directly if script is not being sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    if [ "$#" -lt 1 ] || [ "$#" -gt 3 ]; then
        error "Usage: $0 <repo-name> [true] [owner]"
        exit 1
    fi

    delete_repository "$1" "$2" "$3"
    exit_code=$?
    successMessages
    exit $exit_code
fi
