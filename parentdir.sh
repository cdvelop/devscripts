#!/bin/bash
# Check if Go packages directory exists in current user's home
# Usage: source parentdir.sh  parentDir=$(get_parent_dir)
# This script verifies if the Go packages directory exists in the user's home directory

# reusing the function in other scripts
get_parent_dir() {
    # Obtener la ruta base retrocediendo un directorio desde donde se ejecuta el script
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local parent_dir="$(dirname "$script_dir")"
    echo "$parent_dir"
}


