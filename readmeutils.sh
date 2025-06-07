#!/bin/bash
# Description: README.md utility functions for updating sections dynamically
# Usage: source readmeutils.sh

# Import functions
source functions.sh

# readme_update function - updates or creates a section in README based on identifier
# Parameters:
#   $1: section_identifier (required) - identifier for the section to update
#   $2: new_content (required) - content to insert/update
#   $3: readme_file (optional) - README file path, defaults to README.md
readme_update() {
    local section_identifier="$1"
    local new_content="$2"
    local readme_file="${3:-README.md}"
    
    if [[ -z "$section_identifier" || -z "$new_content" ]]; then
        warning "readme_update requires section_identifier and new_content parameters"
        return 1
    fi
    
    # Handle special case for BADGES -> BADGES_SECTION for backward compatibility
    if [[ "$section_identifier" == "BADGES" ]]; then
        section_identifier="BADGES_SECTION"
    fi
    
    # Check if README exists
    if [[ ! -f "$readme_file" ]]; then
        echo "WARNING: $readme_file not found, creating new file with section" >&2
        local section_start="<!-- START_SECTION:$section_identifier -->"
        local section_end="<!-- END_SECTION:$section_identifier -->"
        local new_section="${section_start}
${new_content}
${section_end}"
        printf "%s\n" "$new_section" > "$readme_file"
        echo "Created with new section"
        return 0
    fi
    
    # Read current README content
    local current_content
    current_content=$(cat "$readme_file")
    
    # Check if section already exists
    local section_start="<!-- START_SECTION:$section_identifier -->"
    local section_end="<!-- END_SECTION:$section_identifier -->"
    
    if echo "$current_content" | grep -q "$section_start"; then
        # Section exists, extract current section content for comparison
        local existing_section
        existing_section=$(echo "$current_content" | sed -n "/$section_start/,/$section_end/p")
        
        # Create new section with identifiers
        local new_section="${section_start}
${new_content}
${section_end}"
        
        # Compare content (ignoring the identifier comments)
        local existing_inner
        existing_inner=$(echo "$existing_section" | sed "1d;\$d")
        
        if [[ "$existing_inner" == "$new_content" ]]; then
            echo "already up to date"
            return 0
        fi
        
        # Update existing section
        local temp_file=$(mktemp)
        awk -v section_start="$section_start" -v section_end="$section_end" -v new_section="$new_section" '
            BEGIN { in_section = 0 }
            $0 == section_start { 
                print new_section
                in_section = 1
                next
            }
            $0 == section_end && in_section {
                in_section = 0
                next
            }
            !in_section { print $0 }
        ' "$readme_file" > "$temp_file"
        mv "$temp_file" "$readme_file"
        echo "Updated existing section"
        return 0
    else
        # Section doesn't exist, add it
        local new_section="${section_start}
${new_content}
${section_end}"
        
        # Add section at the end
        printf "%s\n\n%s\n" "$current_content" "$new_section" > "$readme_file"
        echo "Added new section"
        return 0
    fi
}

# section_update function - generic section updater (wrapper around readme_update)
# Parameters:
#   $1: section_identifier (required)
#   $2: new_content (required)  
#   $3: readme_file (optional, defaults to README.md)
section_update() {
    readme_update "$@"
}

# Main execution - call function based on first argument
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    # Script is being executed directly, not sourced
    if [[ $# -eq 0 ]]; then
        warning "Usage: $0 <function_name> [args...]"
        exit 1
    fi
    
    function_name="$1"
    shift  # Remove function name from arguments
    
    # Call the specified function with remaining arguments
    if [[ "$(type -t "$function_name")" == "function" ]]; then
        "$function_name" "$@"
        exit $?
    else
        warning "Function '$function_name' not found"
        exit 1
    fi
fi
