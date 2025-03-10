#!/bin/bash
# Create and commit a license file for a Git repository
# Usage: ./licensecreate.sh [license-type] [owner-name]
# Currently supported licenses: MIT
# If owner-name is not provided, uses Git config user.name

source functions.sh

generate_license_file() {
    local license_type=${1:-MIT}
    local year=$(date +%Y)
    # Git user name e.g.: John Doe
    local owner=${2:-$(git config user.name)}
    
    case $license_type in
        "MIT")
            cat > LICENSE << EOF
MIT License

Copyright (c) $year $owner

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF
            ;;
        *)
            error "License type not supported"
            return 1
            ;;
    esac
}

commit_license() {
    execute "git add LICENSE" \
        "Failed to stage LICENSE file" \
        "LICENSE file created and staged successfully" || return $?
    
    execute "git commit -m 'Add $license_type license'" \
        "Failed to commit LICENSE" \
        "LICENSE committed successfully" || return $?
        
    execute "git push" \
        "Failed to push LICENSE" \
        "LICENSE pushed successfully" || return $?
        
    return 0
}

create_license() {
    generate_license_file "$1" "$2" || return $?
    commit_license "$1" || return $?
}

# If the script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    if [ "$#" -gt 2 ]; then
        error "Usage: $0 [license-type] [owner-name]"
        exit 1
    fi
    create_license "$1" "$2"
    exit_code=$?
    successMessages
    exit $exit_code
fi
