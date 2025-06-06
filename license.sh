#!/bin/bash
# Description: Detect license type from LICENSE files
# Usage: license.sh
# Returns: License type (MIT, Apache, GNU, etc.) or "MIT" as default

# Function to get license type from LICENSE files
get_license_type() {
    local license_files=("LICENSE.txt" "LICENSE" "LICENSE.md")
    
    for license_file in "${license_files[@]}"; do
        if [ -f "$license_file" ]; then
            # Read first line and extract license type
            first_line=$(head -n 1 "$license_file")
            
            # Remove "License" word and extract first word
            license_type=$(echo "$first_line" | sed 's/[Ll]icense//g' | awk '{print $1}')
            
            # Clean up common patterns
            case "$license_type" in
                "MIT"|"mit")
                    echo "MIT"
                    return 0
                    ;;
                "Apache"|"apache")
                    echo "Apache"
                    return 0
                    ;;
                "GNU"|"gnu")
                    echo "GNU"
                    return 0
                    ;;
                "BSD"|"bsd")
                    echo "BSD"
                    return 0
                    ;;
                *)
                    if [ -n "$license_type" ]; then
                        echo "$license_type"
                        return 0
                    fi
                    ;;
            esac
        fi
    done
    
    # Default to MIT if no license file found or no type detected
    echo "MIT"
    return 0
}

# If script is executed directly, call the function
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    get_license_type
fi
