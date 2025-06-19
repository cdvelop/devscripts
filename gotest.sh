#!/bin/bash
# Description: Run Go tests, vet, race detection and coverage analysis
# Usage: gotest.sh [project_directory]
# Returns: 0 if all tests pass, 1 if any issues found

source functions.sh
source gomodutils.sh

# Change to project directory if provided
if [ $# -gt 0 ]; then
    project_dir="$1"
    if [ -d "$project_dir" ]; then
        cd "$project_dir" || exit 1
        info "Changed to project directory: $project_dir"
    else
        error "Directory does not exist: $project_dir"
        exit 1
    fi
fi

# Get Go module name
if [ -f "go.mod" ]; then
    go_mod_name=$(gawk -v pattern=$currentGitHostUserPath/ 'NR==1 && match($0, pattern "([^/]+)", arr) { print arr[1] }' go.mod)
    if [ -z "$go_mod_name" ]; then
        go_mod_name=$(awk 'NR==1 {print $2}' go.mod | sed 's/.*\///')
    fi
else
    echo "Error: go.mod file not found"
    exit 1
fi

# Show what will be executed
warning "Running tests, race detection, vet and coverage analysis for $go_mod_name..."

# Initialize results
test_status="Failed"
coverage_percent="0"
race_status="Detected"
vet_status="Issues"

# Run go vet
execute "go vet ./..." "go vet failed in $go_mod_name" "vet passed" "no_exit"
if [ $? -eq 0 ]; then
    vet_status="OK"
fi

# Check if test files exist
if [ -n "$(find . -type f -name "*_test.go")" ]; then
    # Run tests
    execute "go test ./..." "Test errors found in $go_mod_name" "tests passed" "no_exit"
    if [ $? -eq 0 ]; then
        test_status="Passing"
    fi
    
    # Run race detection tests
    execute "go test -race ./..." "Race condition tests failed in $go_mod_name" "race detection passed" "no_exit"
    if [ $? -eq 0 ]; then
        race_status="Clean"
    fi
      # Calculate coverage
    coverage_output=$(go test -cover ./... 2>/dev/null | grep "coverage:")
    if [ -n "$coverage_output" ]; then
        # Extract coverage percentages, excluding 0.0% (directories without tests)
        percentages=$(echo "$coverage_output" | grep -o '[0-9]\+\.[0-9]\+%' | sed 's/%//' | grep -v "^0\.0$")
        if [ -n "$percentages" ]; then
            total=0
            count=0
            while IFS= read -r percent; do
                if [ -n "$percent" ] && [ "$percent" != "0.0" ]; then
                    total=$(awk "BEGIN {print $total + $percent}")
                    count=$((count + 1))
                fi
            done <<< "$percentages"
            
            if [ $count -gt 0 ]; then
                average=$(awk "BEGIN {printf \"%.0f\", $total / $count}")
                coverage_percent="$average"
                addOKmessage "coverage calculated"
            fi
        fi
    else
        addERRORmessage "Failed to calculate coverage in $go_mod_name"
    fi
else
    addOKmessage "no test files found in $go_mod_name"
    coverage_percent="0"
fi

# Get license type (optional parameter)
license_type=$(source license.sh && get_license_type)

# Call gobadge.sh to update README
gobadge.sh "$go_mod_name" "$test_status" "$coverage_percent" "$race_status" "$vet_status" "$license_type"
badge_exit_code=$?

# Print accumulated messages
successMessages
exit $badge_exit_code
