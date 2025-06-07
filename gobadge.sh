#!/bin/bash
# Description: Update README.md with project status badges
# Usage: gobadge.sh <module_name> <test_status> <coverage_percent> <race_status> <vet_status> [license_type]

source functions.sh
source gomodutils.sh
source githubutils.sh
source readmeutils.sh

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

# Function to get CSS class for badge values
get_css_class() {
    local badge_type="$1"
    local value="$2"
    
    case "$badge_type" in
        "license")
            echo "license"
            ;;
        "go")
            echo "go-version"
            ;;
        "tests")
            if [ "$value" = "Passing" ]; then
                echo "tests-passing"
            else
                echo "tests-failing"
            fi
            ;;
        "coverage")
            local num=$(echo "$value" | sed 's/%//' | tr -d ' ')
            if [ "$num" -ge 80 ]; then
                echo "coverage-high"
            elif [ "$num" -ge 60 ]; then
                echo "coverage-medium"
            elif [ "$num" -gt 0 ]; then
                echo "coverage-low"
            else
                echo "coverage-none"
            fi
            ;;
        "race")
            if [ "$value" = "Clean" ]; then
                echo "race-clean"
            else
                echo "race-detected"
            fi
            ;;
        "vet")
            if [ "$value" = "OK" ]; then
                echo "vet-ok"
            else
                echo "vet-issues"
            fi
            ;;
        *)
            echo "default"
            ;;
    esac
}

# Function to create badges.css file in .github directory
create_badges_css() {
    local github_dir=$(ensure_github_directory)
    local css_file="$github_dir/badges.css"
    local source_css="badges.css"
    
    # Check if source badges.css exists in current directory
    if [ -f "$source_css" ]; then
        # Always copy the source badges.css to .github directory to ensure it's up to date
        cp "$source_css" "$css_file"
        return 0
    else
        error "badges.css file not found in current directory"
        return 1
    fi
}

# Generate badge HTML with CSS classes
generate_badge_html() {
    # Ensure .github directory exists and create badges.css
    create_badges_css
    
    local license_class=$(get_css_class "license" "$license_type")
    local go_class=$(get_css_class "go" "$go_version")
    local test_class=$(get_css_class "tests" "$test_status")
    local coverage_class=$(get_css_class "coverage" "$coverage_percent")
    local race_class=$(get_css_class "race" "$race_status")
    local vet_class=$(get_css_class "vet" "$vet_status")
    
    local badge_html
    badge_html=$(cat << EOF
<link rel="stylesheet" href=".github/badges.css">
<div class="project-badges">
    <div class="badge-group">
        <span class="badge-label">License</span><span class="badge-value $license_class">$license_type</span>
    </div>
    <div class="badge-group">
        <span class="badge-label">Go</span><span class="badge-value $go_class">$go_version</span>
    </div>
    <div class="badge-group">
        <span class="badge-label">Tests</span><span class="badge-value $test_class">$test_status</span>
    </div>
    <div class="badge-group">
        <span class="badge-label">Coverage</span><span class="badge-value $coverage_class">${coverage_percent}%</span>
    </div>
    <div class="badge-group">
        <span class="badge-label">Race</span><span class="badge-value $race_class">$race_status</span>
    </div>
    <div class="badge-group">
        <span class="badge-label">Vet</span><span class="badge-value $vet_class">$vet_status</span>
    </div>
</div>
EOF
)    # Use section_update from readmeutils.sh to handle README update
    local output
    output=$(section_update "BADGES_SECTION" "$badge_html" "$readme_file" 2>&1)
    local exit_code=$?
    
    if [[ $exit_code -ne 0 ]]; then
        error "Failed to update badges"
        return 1
    fi
    
    # Parse the output to determine what happened and generate appropriate message
    if echo "$output" | grep -q "WARNING:"; then
        # Show the warning about missing file
        echo "$output" | grep "WARNING:"
        echo "Adding new badges"
    elif echo "$output" | grep -q "already up to date"; then
        echo "Badges are already up to date"
    elif echo "$output" | grep -q "Added new section"; then
        echo "Adding new badges"
    elif echo "$output" | grep -q "Updated existing section"; then
        echo "Updating existing badges"
    elif echo "$output" | grep -q "Created.*with new section"; then
        echo "Adding new badges"
    else
        # Default case
        echo "Adding new badges"
    fi
    
    return 0
}

# Main execution
generate_badge_html
exit $?
