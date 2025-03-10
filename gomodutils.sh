#!/bin/bash
# Utility functions for managing Go modules and version updates
# Usage: source gomodutils.sh && update_single_go_module "mymodule" "v1.2.3"

# Function to get current module version from go.mod
get_go_module_version() {
    local module_name=$1
    local go_mod_file=$2
    gawk -v package="$module_name" -v common="$currentGitHostUserPath" \
        'match($0, "^require[[:space:]]+" common "/" package "[[:space:]]+([^[:space:]]+)", tag) {print tag[1]; exit} \
         $1==common "/" package {print $2}' "$go_mod_file"
}

# Function to run go mod tidy and verify tests
update_and_verify_go_module() {
    local pkg_name=$1
    
    execute "go mod tidy" "go mod tidy failed in $pkg_name" || return 1
    
    bash gomodcheck.sh
    return $?
}

# Function to update a specific module
update_single_go_module() {
    local module_name=$1
    local new_version=$2
    execute "go get $currentGitHostUserPath/$module_name@$new_version" \
        "failed to update package $module_name" || return 1
    execute "go get -u all" "failed to update packages" || return 1
}
