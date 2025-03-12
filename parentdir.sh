#!/bin/bash
# Gets the parent directory of the script's location
# Usage: source parentdir.sh  parentDir=$(get_parent_dir)

# reusing the function in other scripts
get_parent_dir() {
    # Obtener la ruta base retrocediendo un directorio desde donde se ejecuta el script
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local parent_dir="$(dirname "$script_dir")"
    echo "$parent_dir"
}


