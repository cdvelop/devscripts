#!/bin/bash
# Functions to work with issues.md file
# usage: source fileIssues.sh

# Function to get commit message from issues.md file
# Only returns completed tasks marked with [x]
# Usage: get_commit_message_from_doing_md "Initial message"
# Returns: Updated commit message with completed tasks from issues.md
function get_commit_message_from_doing_md() {
    local initial_message="$1"
    local result="$initial_message"
    local completed_tasks=""
    
    # Check if issues.md exists and has content
    if [ -f "issues.md" ] && [ -s "issues.md" ]; then
        # Extract only completed tasks from the file
        while IFS= read -r line; do
            if [[ $line == *"[x]"* ]]; then
                # Extract the task text without the [x] marker
                local task_text=$(echo "$line" | sed -e 's/\[x\] *//')
                
                # Add to completed tasks list
                if [ -n "$completed_tasks" ]; then
                    completed_tasks="$completed_tasks, $task_text"
                else
                    completed_tasks="$task_text"
                fi
            fi
        done < "issues.md"
        
        # If we found completed tasks, add them to the result
        if [ -n "$completed_tasks" ]; then
            if [ -n "$initial_message" ]; then
                # Concatenate the initial message with the completed tasks
                result="$initial_message: $completed_tasks"
            else
                result="$completed_tasks"
            fi
        fi
    fi
    
    echo "$result"
}

# Create issues.md file with initial template
create_doing_md_file() {
    local content="[x] init code\n[ ] task 1"
    execute "echo -e \"$content\" > issues.md" \
        "Failed to create issues.md" \
        "issues.md created" || return $?

    return 0
}

# Function to remove completed tasks from issues.md but keep incomplete ones
function deleteChangesIssueFile() {
    if [ -f "issues.md" ] && [ -s "issues.md" ]; then
        # Create a temporary file
        local temp_file=$(mktemp)
        
        # Keep only incomplete tasks
        while IFS= read -r line; do
            # If line is not a completed task or is empty/other content, keep it
            if [[ ! $line =~ ^\[x\] ]]; then
                echo "$line" >> "$temp_file"
            fi
        done < "issues.md"
        
        # Replace original file with filtered version
        mv "$temp_file" "issues.md"
    fi
}