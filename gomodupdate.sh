#!/bin/bash
# Updates Go module name and its references in all project files
# Usage: ./gomodupdate.sh old-module-name new-module-name

source functions.sh
source gomodutils.sh

update_go_module_name() {
    local old_name=$1
    local new_name=$2

    # Validate required arguments
    if [ -z "$old_name" ] || [ -z "$new_name" ]; then
        error "Usage: update_go_module_name <old-name> <new-name>"
        return 1
    fi

    # Get GitHub username from local git config
    gitHubUser=$(git config --get user.name)
    if [ -z "$gitHubUser" ]; then
        error "Unable to get GitHub username from git config"
        return 1
    fi

    local old_module="github.com/$gitHubUser/$old_name"
    local new_module="github.com/$gitHubUser/$new_name"

    # Update go.mod file
    if [ -f "go.mod" ]; then
        execute "sed -i \"s|module $old_module|module $new_module|g\" go.mod" \
            "Failed to update module name in go.mod" \
            "Updated module name in go.mod" || return $?
    else
        warning "go.mod file not found"
        return 1
    fi

    # Find all .go files and update import statements
    success "Updating import statements in Go files..."
    
    # Using find to locate all .go files
    go_files=$(find . -type f -name "*.go")
    
    if [ -z "$go_files" ]; then
        success "No Go files found"
    else
        for file in $go_files; do
            execute "sed -i \"s|$old_module|$new_module|g\" \"$file\"" \
                "Failed to update imports in $file" \
                "Updated imports in $file" || continue
        done
    fi

    # Update go.sum if it exists
    if [ -f "go.sum" ]; then
        execute "sed -i \"s|$old_module|$new_module|g\" go.sum" \
            "Failed to update references in go.sum" \
            "Updated references in go.sum" || return $?
    fi

    # Run go mod tidy to clean up dependencies
    if command -v go >/dev/null 2>&1; then
        execute "go mod tidy" \
            "Failed to tidy up Go modules" \
            "Go modules tidied up successfully" || return $?
    else
        warning "Go command not found. Skip running 'go mod tidy'"
    fi

    return 0
}

# Execute directly if script is not being sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    if [ "$#" -ne 2 ]; then
        error "Usage: $0 <old-name> <new-name>"
        exit 1
    fi

    update_go_module_name "$1" "$2"
    exit_code=$?
    successMessages
    exit $exit_code
fi