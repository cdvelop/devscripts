#!/bin/bash
# Check if a Go package uses syscall/js imports
# Usage: ./syscall.sh <package_name>

syscall(){
    pkg_name=$1
    # warning "Error found in package $pkg_name"
    found=0  # Variable to track if match was found
    for file in $(find . -type f -name "*.go"); do
        # echo "File found: $file of type syscall/js in package $pkg_name"
        if grep -q "syscall/js" "$file"; then
            found=1
            break
        fi
    done
    # success "package $pkg_name is of type 'syscall/js'."
    return $found
}

# Example usage:
# execute "go vet" "go vet failed in $go_mod_name" "go vet $go_mod_name ok" "no exist"
# syscall $? $go_mod_name