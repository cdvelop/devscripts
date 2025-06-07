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

# Initialize error tracking
has_errors=0
error_messages=""

# Function to add error message
add_error() {
    local message="$1"
    has_errors=1
    if [ -z "$error_messages" ]; then
        error_messages="$message"
    else
        error_messages="$error_messages\n$message"
    fi
}

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

# Initialize results
test_status="Failed"
coverage_percent="0"
race_status="Detected"
vet_status="Issues"

# Run go vet
echo "Running go vet..."
if go vet ./... 2>/dev/null; then
    vet_status="OK"
    info "go vet $go_mod_name ok"
else
    add_error "go vet failed in $go_mod_name"
fi

# Check if test files exist
if [ -n "$(find . -type f -name "*_test.go")" ]; then
    # Run tests
    echo "Running tests..."
    if go test ./... 2>/dev/null; then
        test_status="Passing"
        info "All tests in $go_mod_name passed"
    else
        add_error "Test errors found in $go_mod_name"
    fi
    
    # Run race detection tests
    echo "Running race detection tests..."
    if go test -race ./... 2>/dev/null; then
        race_status="Clean"
        info "Race condition tests in $go_mod_name passed"
    else
        add_error "Race condition tests failed in $go_mod_name"
    fi
    
    # Calculate coverage
    echo "Calculating test coverage..."
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
            fi
        fi
    else
        add_error "Failed to calculate coverage in $go_mod_name"
    fi
else
    info "No test files found in $go_mod_name"
    coverage_percent="0"
fi

# Get license type (optional parameter)
license_type=$(source license.sh && get_license_type)

# Call gobadge.sh to update README
echo "Updating badges in README.md..."
bash gobadge.sh "$go_mod_name" "$test_status" "$coverage_percent" "$race_status" "$vet_status" "$license_type"

# Print accumulated error messages if any
if [ $has_errors -eq 1 ]; then
    echo ""
    echo "ERRORS FOUND:"
    echo -e "$error_messages"
    exit 1
fi

echo "All checks completed successfully for $go_mod_name"
exit 0
