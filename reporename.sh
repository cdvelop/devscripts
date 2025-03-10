#!/bin/bash
# Description: Renames a repository both locally and on remote GitHub, updates Git remotes and module names if applicable
# Usage: ./reporename.sh <old-name> <new-name>

source functions.sh

check_rename_permissions() {
    # First check if gh CLI is authenticated at all
    if ! gh auth status >/dev/null 2>&1; then
        warning "GitHub CLI not authenticated. Please run: gh auth login"
        return 1
    fi
    
    return 0
}

rename_repository() {
    local old_name=$1
    local new_name=$2

    # Validate required arguments
    if [ -z "$old_name" ] || [ -z "$new_name" ]; then
        error "Usage: rename_repository <old-name> <new-name>"
        return 1
    fi

    # Check if the repository directory exists in current directory
    if [ ! -d "$old_name" ]; then
        error "Directory '$old_name' not found in current directory. Please run this script from the parent directory of the repository."
        return 1
    fi

    # Check permissions first
    check_rename_permissions || return $?

    # Get GitHub username from local git config
    gitHubUser=$(git config --get user.name)
    if [ -z "$gitHubUser" ]; then
        error "Unable to get GitHub username from git config"
        return 1
    fi
    
    # Get directory structure - analyze the target path
    local target_dir=$(basename "$(realpath "$old_name")")
    local parent_dir=$(basename "$(pwd)")
    local grandparent_dir=$(basename "$(dirname "$(pwd)")")
    
    # Confirm rename
    read -p "Are you sure you want to rename repository '$old_name' to '$new_name'? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        warning "Operation cancelled by user"
        return 1
    fi

    # Rename repository using GitHub CLI (from outside the directory)
    execute "gh api -X PATCH repos/$gitHubUser/$old_name -f name=$new_name" \
        "Failed to rename repository" \
        "Repository renamed from $old_name to $new_name successfully" || return $?
    
    # Check if we need to rename local folder based on directory structure
    local should_rename=false
    if [ "$parent_dir" = "Packages" ] || [ "$parent_dir" = "Internal" -a "$grandparent_dir" = "Packages" ]; then
        should_rename=true
    fi
    
    if [ "$should_rename" = true ]; then
        # Rename the folder
        execute "mv '$old_name' '$new_name'" \
            "Failed to rename local directory" \
            "Local directory renamed from $old_name to $new_name" || return $?
        
        # Update Git remotes - need to cd into the new directory
        pushd "$new_name" > /dev/null
        execute "git remote set-url origin https://github.com/$gitHubUser/$new_name.git" \
            "Failed to update Git remote URL" \
            "Git remote URL updated successfully" || { popd > /dev/null; return $?; }
        
        # If it's a Go module, update module name
        if [ -f "go.mod" ]; then
            success "Go module detected. Updating module name..."
            if command -v repomoduleupdate.sh >/dev/null 2>&1; then
                repomoduleupdate.sh "$old_name" "$new_name"
            else
                warning "repomoduleupdate.sh script not found. Module name not updated."
            fi
        fi
        popd > /dev/null
    else
        success "Directory structure doesn't match standard pattern. Local folder not renamed."
    fi

    return 0
}

# Execute directly if script is not being sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    if [ "$#" -ne 2 ]; then
        error "Usage: $0 <old-name> <new-name>"
        exit 1
    fi

    rename_repository "$1" "$2"
    exit_code=$?
    successMessages
    exit $exit_code
fi