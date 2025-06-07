#!/bin/bash
# Description: Update README.md with project status badges
# Usage: gobadge.sh <module_name> <test_status> <coverage_percent> <race_status> <vet_status> [license_type]

source functions.sh
source gomodutils.sh

# Validate parameters
if [ $# -lt 5 ]; then
    echo "Usage: gobadge.sh <module_name> <test_status> <coverage_percent> <race_status> <vet_status> [license_type]"
    exit 1
fi

module_name="$1"
test_status="$2"
coverage_percent="$3"
race_status="$4"
vet_status="$5"
license_type="${6:-$(source license.sh && get_license_type)}"

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

# Generate badge HTML with CSS classes
generate_badge_html() {
    local license_class=$(get_css_class "license" "$license_type")
    local go_class=$(get_css_class "go" "$go_version")
    local test_class=$(get_css_class "tests" "$test_status")
    local coverage_class=$(get_css_class "coverage" "$coverage_percent")
    local race_class=$(get_css_class "race" "$race_status")
    local vet_class=$(get_css_class "vet" "$vet_status")
    
    cat << EOF
<!-- Generated dynamically by gotest.sh from github.com/cdvelop/devscripts -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/cdvelop/devscripts@main/badges.css">
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
}

# Update README.md
update_readme() {
    local readme_file="README.md"
    
    # Check if README.md exists
    if [ ! -f "$readme_file" ]; then
        warning "README.md not found, skipping badge update"
        return 0
    fi
    
    # Generate new badge HTML
    local new_badge_html=$(generate_badge_html)
    
    # Check if file has a title (first line starting with #)
    local first_title_line=$(grep -n "^#[^#]" "$readme_file" | head -n 1 | cut -d: -f1)
    
    if [ -z "$first_title_line" ]; then
        # No title found, add one at the beginning
        echo "Adding title to README.md"
        temp_file=$(mktemp)
        echo "# $module_name" > "$temp_file"
        echo "$new_badge_html" >> "$temp_file"
        echo "" >> "$temp_file"
        cat "$readme_file" >> "$temp_file"
        mv "$temp_file" "$readme_file"
    else
        # Title exists, check if badges already exist
        local comment_line=$(grep -n "Generated dynamically by gotest.sh" "$readme_file" | cut -d: -f1)
        
        if [ -n "$comment_line" ]; then
            # Badges exist, replace them
            echo "Updating existing badges in README.md"
            
            # Find the closing </div> for project-badges div
            # We need to count div open/close after we find project-badges
            local end_line=""
            local line_num=$comment_line
            local found_project_div=false
            local div_depth=0
            
            while IFS= read -r line; do
                # Check if this line contains project-badges div opening
                if [[ "$line" == *"project-badges"* ]] && [[ "$line" == *"<div"* ]]; then
                    found_project_div=true
                    div_depth=1  # We found the opening project-badges div
                elif [[ "$found_project_div" == true ]]; then
                    # Count div openings and closings after finding project-badges
                    local open_count=$(echo "$line" | grep -o '<div' | wc -l)
                    local close_count=$(echo "$line" | grep -o '</div>' | wc -l)
                    
                    div_depth=$((div_depth + open_count - close_count))
                    
                    # When div_depth reaches 0, we've closed the project-badges div
                    if [ $div_depth -eq 0 ]; then
                        end_line=$line_num
                        break
                    fi
                fi
                
                line_num=$((line_num + 1))
            done < <(tail -n +$comment_line "$readme_file")
            
            if [ -n "$end_line" ]; then
                # Create temp file with content before badges, new badges, and content after badges
                temp_file=$(mktemp)
                head -n $((comment_line - 1)) "$readme_file" > "$temp_file"
                echo "$new_badge_html" >> "$temp_file"
                
                # Add empty line before continuing content if the next line isn't empty
                local next_line_num=$((end_line + 1))
                local next_line=$(sed -n "${next_line_num}p" "$readme_file")
                if [ -n "$next_line" ] && [ "$next_line" != " " ]; then
                    echo "" >> "$temp_file"
                fi
                
                tail -n +$next_line_num "$readme_file" >> "$temp_file"
                mv "$temp_file" "$readme_file"
            else
                warning "Could not find end of badge section, adding new badges instead"
                # Fallback: replace everything from comment line onwards until we find meaningful content
                temp_file=$(mktemp)
                head -n $((comment_line - 1)) "$readme_file" > "$temp_file"
                echo "$new_badge_html" >> "$temp_file"
                echo "" >> "$temp_file"
                
                # Skip lines until we find content that doesn't look like badges/HTML
                local skip_line=$comment_line
                while IFS= read -r line; do
                    skip_line=$((skip_line + 1))
                    # If line is empty or doesn't contain HTML tags, start including content
                    if [[ "$line" != *"<"* ]] && [[ "$line" != *">"* ]] && [ -n "$(echo "$line" | xargs)" ]; then
                        tail -n +$skip_line "$readme_file" >> "$temp_file"
                        break
                    fi
                done < <(tail -n +$((comment_line + 1)) "$readme_file")
                mv "$temp_file" "$readme_file"            
            fi
        else
            # No badges exist, add them after the title
            echo "Adding new badges to README.md"
            temp_file=$(mktemp)
            head -n $first_title_line "$readme_file" > "$temp_file"
            echo "$new_badge_html" >> "$temp_file"
            tail -n +$((first_title_line + 1)) "$readme_file" >> "$temp_file"
            mv "$temp_file" "$readme_file"
        fi
    fi
    
    info "README.md badges updated successfully"
    return 0
}

# Main execution
update_readme
exit $?
