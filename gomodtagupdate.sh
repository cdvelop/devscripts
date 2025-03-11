#!/bin/bash
# Updates Go module versions across all projects that use them
# Usage: ./gomodtagupdate.sh <package-name> <new-version>
# This script updates a specific Go module version in all projects
# that depend on it, running tests to verify the update

source functions.sh
source gomodutils.sh

# Update module version in all projects that use it
update_module_version() {
    local pkg_name=$1
    local new_version=$2
    if [ -d "$go_pkgs" ]; then
        for observed_pkg in "$go_pkgs"/*; do
            if [ "$(basename "$observed_pkg")" != "$pkg_name" ]; then
                go_mod_file="$observed_pkg/go.mod"
                
                if [ -f "$go_mod_file" ]; then
                    old_version=$(get_go_module_version "$pkg_name" "$go_mod_file")
                    
                    if [ -n "$old_version" ] && [ "$old_version" != "$new_version" ]; then
                        current_dir=$(pwd)
                        cd "$observed_pkg"
                        
                        update_single_go_module "$pkg_name" "$new_version"
                        
                        if update_and_verify_go_module "$(basename "$observed_pkg")"; then
                            bash pu.sh "update module: $pkg_name"
                            addOKmessage "package $pkg_name updated in $(basename "$observed_pkg")"
                        else
                            addERRORmessage "tests failed in: $(basename "$observed_pkg") while updating: [$pkg_name]"
                        fi
                        
                        cd "$current_dir"
                    fi
                fi
            fi
        done
    fi
}

if [ "$#" -ne 2 ] || [ -z "$1" ] || [ -z "$2" ]; then
    error "Usage: $0 <package-name> <new-version>"
    exit 1
fi

update_module_version "$1" "$2"
successMessages
exit 0
