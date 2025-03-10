#!/bin/bash
# Script to rename a Go project and update its module references
# Usage: ./gorenameproject.sh old-project-name new-project-name

source functions.sh

# Check if the required scripts exist
check_required_scripts() {
    local missing_scripts=0
    
    if ! command -v reporename.sh >/dev/null 2>&1; then
        error "Required script 'reporename.sh' not found in PATH"
        missing_scripts=1
    fi
    
    if ! command -v gomodupdate.sh >/dev/null 2>&1; then
        error "Required script 'gomodupdate.sh' not found in PATH"
        missing_scripts=1
    fi
    
    return $missing_scripts
}

# Main function to rename a Go project
rename_go_project() {
    local old_name=$1
    local new_name=$2

    # Validate required arguments
    if [ -z "$old_name" ] || [ -z "$new_name" ]; then
        error "Usage: rename_go_project <old-name> <new-name>"
        return 1
    }
    
    # Check if the repository directory exists in current directory
    if [ ! -d "$old_name" ]; then
        error "Directory '$old_name' not found in current directory. Please run this script from the parent directory of the repository."
        return 1
    }
    
    # Check for required scripts
    check_required_scripts || return $?
    
    # Step 1: First check if it's a Go project by looking for go.mod in the target directory
    if [ ! -f "$old_name/go.mod" ]; then
        warning "No go.mod file found in $old_name. This doesn't appear to be a Go project."
        read -p "Continue anyway? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            warning "Operation cancelled by user"
            return 1
        fi
    fi
    
    # Step 2: Rename the repository
    success "Step 1/2: Renaming repository from '$old_name' to '$new_name'..."
    if ! reporename.sh "$old_name" "$new_name"; then
        error "Repository rename failed"
        return 1
    fi
    
    # Step 3: Update Go module references if go.mod exists
    if [ -f "$new_name/go.mod" ]; then
        success "Step 2/2: Updating Go module references..."
        
        # Use pushd/popd to safely enter and exit the new directory
        pushd "$new_name" > /dev/null
        if ! gomodupdate.sh "$old_name" "$new_name"; then
            error "Module update failed"
            warning "Repository was renamed but module references may not be fully updated"
            popd > /dev/null
            return 1
        fi
        
        # Make sure we're displaying the correct, updated module name by reading it after updates
        sync
        sleep 1
        module_name=$(grep "^module" go.mod | cut -d ' ' -f 2)
        success "New module name: $module_name"
        
        popd > /dev/null
    else
        success "Skipping module update as this is not a Go project"
    fi
    
    success "Go project '$old_name' has been successfully renamed to '$new_name'"
    
    # Additional information
    gitHubUser=$(git config --get user.name)
    success "New repository URL: https://github.com/$gitHubUser/$new_name"
    
    return 0
}

# Execute directly if script is not being sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    if [ "$#" -ne 2 ]; then
        error "Usage: $0 <old-name> <new-name>"
        exit 1
    fi

    rename_go_project "$1" "$2"
    exit_code=$?
    successMessages
    exit $exit_code
fi