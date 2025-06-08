#!/bin/bash
# Description: Determine Go project badges and generate them using badges.sh
# Usage: gobadge.sh <module_name> <test_status> <coverage_percent> <race_status> <vet_status> [license_type]

source functions.sh
source gomodutils.sh
source githubutils.sh

# Validate parameters and set defaults
readme_file="README.md"

# Check if first parameter is a readme file (ends with .md)
if [[ $# -eq 1 && "$1" =~ \.md$ ]]; then
    readme_file="$1"
    # Set default values for testing
    module_name="testmodule"
    test_status="Passing"
    coverage_percent="85"
    race_status="Clean"
    vet_status="OK"
    license_type="MIT"
elif [[ $# -eq 0 ]]; then
    # No parameters - use defaults
    module_name="testmodule"
    test_status="Passing"
    coverage_percent="85"
    race_status="Clean"
    vet_status="OK"
    license_type="MIT"
elif [[ $# -lt 5 ]]; then
    echo "Usage: gobadge.sh [readme_file.md] OR gobadge.sh <module_name> <test_status> <coverage_percent> <race_status> <vet_status> [license_type]"
    exit 1
else
    module_name="$1"
    test_status="$2"
    coverage_percent="$3"
    race_status="$4"
    vet_status="$5"
    license_type="${6:-$(source license.sh && get_license_type)}"
fi

# Get Go version
go_version=$(get_go_version)
if [ -z "$go_version" ]; then
    go_version="1.22"
fi

# Function to get color for badge values
get_badge_color() {
    local badge_type="$1"
    local value="$2"
    
    case "$badge_type" in
        "license")
            echo "#007acc"
            ;;
        "go")
            echo "#00ADD8"
            ;;
        "tests")
            if [ "$value" = "Passing" ]; then
                echo "#4c1"
            else
                echo "#e05d44"
            fi
            ;;
        "coverage")
            local num=$(echo "$value" | sed 's/%//' | tr -d ' ')
            if [ "$num" -ge 80 ]; then
                echo "#4c1"
            elif [ "$num" -ge 60 ]; then
                echo "#dfb317"
            elif [ "$num" -gt 0 ]; then
                echo "#fe7d37"
            else
                echo "#e05d44"
            fi
            ;;
        "race")
            if [ "$value" = "Clean" ]; then
                echo "#4c1"
            else
                echo "#e05d44"
            fi
            ;;
        "vet")
            if [ "$value" = "OK" ]; then
                echo "#4c1"
            else
                echo "#e05d44"
            fi
            ;;
        *)
            echo "#007acc"
            ;;
    esac
}

# Generate badges by calling badges.sh with determined parameters
generate_badges() {
    # Get colors for each badge
    local license_color=$(get_badge_color "license" "$license_type")
    local go_color=$(get_badge_color "go" "$go_version")
    local test_color=$(get_badge_color "tests" "$test_status")
    local coverage_color=$(get_badge_color "coverage" "$coverage_percent")
    local race_color=$(get_badge_color "race" "$race_status")
    local vet_color=$(get_badge_color "vet" "$vet_status")
    
    # Call badges.sh with all badge parameters and custom readme file
    ./badges.sh \
        "readmefile:$readme_file" \
        "License:$license_type:$license_color" \
        "Go:$go_version:$go_color" \
        "Tests:$test_status:$test_color" \
        "Coverage:${coverage_percent}%:$coverage_color" \
        "Race:$race_status:$race_color" \
        "Vet:$vet_status:$vet_color"
      local exit_code=$?
    if [[ $exit_code -ne 0 ]]; then
        error "Failed to generate badges"
        return 1
    fi
    
    return 0
}

# Main execution
generate_badges
exit $?
