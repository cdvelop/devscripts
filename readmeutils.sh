#!/bin/bash
# Description: README.md utility functions for updating sections dynamically
# Usage: source readmeutils.sh
# Main function: section_update section_identifier [after_line] new_content [readme_file]

# Import functions
source functions.sh

# section_update function - updates or creates a section in README based on identifier
# Parameters:
#   $1: section_identifier (required) - identifier for the section to update
#   $2: after_line (optional) - line number after which to place the section (empty for default behavior)
#   $3: new_content (required) - content to insert/update
#   $4: readme_file (optional) - README file path, defaults to README.md
section_update() {
    local section_identifier="$1"
    local after_line="$2"
    local new_content="$3"
    local readme_file="${4:-README.md}"
    
    if [[ -z "$section_identifier" || -z "$new_content" ]]; then
        warning "section_update requires section_identifier and new_content parameters"
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
        
        if [[ -n "$after_line" ]]; then
            # Create file with section at specified position
            printf "%s\n" "$new_section" > "$readme_file"
        else
            # Default behavior - just create the section
            printf "%s\n" "$new_section" > "$readme_file"
        fi
        echo "Created with new section"
        return 0
    fi
    
    # Read current README content into array
    local lines=()
    while IFS= read -r line; do
        lines+=("$line")
    done < "$readme_file"
    
    local section_start="<!-- START_SECTION:$section_identifier -->"
    local section_end="<!-- END_SECTION:$section_identifier -->"
    
    # Find existing section
    local existing_start_line=-1
    local existing_end_line=-1
    
    for i in "${!lines[@]}"; do
        if [[ "${lines[$i]}" == "$section_start" ]]; then
            existing_start_line=$i
        elif [[ "${lines[$i]}" == "$section_end" ]] && [[ $existing_start_line -ge 0 ]]; then
            existing_end_line=$i
            break
        fi
    done
    
    # Extract existing content if section exists
    local existing_content=""
    if [[ $existing_start_line -ge 0 ]] && [[ $existing_end_line -ge 0 ]]; then
        for ((i=existing_start_line+1; i<existing_end_line; i++)); do
            if [[ -n "$existing_content" ]]; then
                existing_content+=$'\n'
            fi
            existing_content+="${lines[$i]}"
        done
    fi
    
    # Check if content is the same
    if [[ "$existing_content" == "$new_content" ]]; then
        # Content is the same, check if position needs to change
        if [[ -n "$after_line" ]] && [[ $existing_start_line -ge 0 ]]; then
            local target_line=$after_line
            # Check if section is already in the correct position
            if [[ $existing_start_line -eq $target_line ]]; then
                echo "already up to date"
                return 0
            fi
            # Position is different, need to move section
        else
            echo "already up to date"
            return 0
        fi
    fi
    
    # Remove existing section if it exists
    if [[ $existing_start_line -ge 0 ]] && [[ $existing_end_line -ge 0 ]]; then
        # Remove from end to start to maintain line numbers
        for ((i=existing_end_line; i>=existing_start_line; i--)); do
            unset 'lines[$i]'
        done
        # Reindex array
        local temp_lines=()
        for line in "${lines[@]}"; do
            temp_lines+=("$line")
        done
        lines=("${temp_lines[@]}")
    fi
    
    # Determine where to insert the new section
    local insert_position
    if [[ -n "$after_line" ]]; then
        # Insert after specified line
        insert_position=$after_line
        if [[ $insert_position -gt ${#lines[@]} ]]; then
            insert_position=${#lines[@]}
        fi
    else
        # Default behavior - append at end
        insert_position=${#lines[@]}
    fi
    
    # Create new section
    local new_section_lines=(
        "$section_start"
        "$new_content"
        "$section_end"
    )
    
    # Insert new section
    local final_lines=()
    
    # Add lines before insertion point
    for ((i=0; i<insert_position; i++)); do
        if [[ $i -lt ${#lines[@]} ]]; then
            final_lines+=("${lines[$i]}")
        fi
    done
    
    # Add new section
    for section_line in "${new_section_lines[@]}"; do
        final_lines+=("$section_line")
    done
    
    # Add remaining lines after insertion point
    for ((i=insert_position; i<${#lines[@]}; i++)); do
        final_lines+=("${lines[$i]}")
    done
    
    # Write updated content to file
    printf "%s\n" "${final_lines[@]}" > "$readme_file"
    
    if [[ $existing_start_line -ge 0 ]]; then
        echo "Updated existing section"
    else
        echo "Added new section"
    fi
    
    return 0
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
