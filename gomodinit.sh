#!/bin/bash
# Script to initialize a Go module and create basic project structure
# Usage: ./gomodinit.sh

# Import common functions
source functions.sh

init_go_module() {
    local current_folder=$1
    local git_remote=$2

    if [ -z "$current_folder" ] || [ -z "$git_remote" ]; then
        error "Usage: init_go_module <folder_name> <git_remote>"
        return 1
    fi

    if [ ! -f "go.mod" ]; then
        execute "go mod init $git_remote/$current_folder" \
            "Failed to initialize go mod" \
            "Go module initialized successfully" || return $?
    fi
    return 0
}

create_handler_file() {
    local current_folder=$1
    local struct="${current_folder^}" # Capitalize first letter
    local file="$current_folder.go"
    # Get first letter of folder name
    local handler="${current_folder:0:1}"

    local model="type $struct struct{}"
    local func="func New() *$struct {\n\n    $handler := &$struct{}\n\n    return $handler\n}"
    
    execute "echo -e 'package $current_folder\n\n$model\n\n$func' > $file" \
        "Failed to create file $file.go" \
        "file $file created successfully" || return $?
    
    return 0
}

setup_go_project() {
    local current_folder=$(basename "$(pwd)")
    
    # Initialize go module
    init_go_module "$current_folder" "$currentGitHostUserPath" || return $?
    
    # Create handler file
    create_handler_file "$current_folder" || return $?
    
    return 0
}

# Execute directly if script is not being sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    setup_go_project
    exit_code=$?
    successMessages
    exit $exit_code
fi