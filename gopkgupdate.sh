#!/bin/bash
# Updates Go packages in go.mod to their latest versions from local repositories
# Usage: ./gopkgupdate.sh

source functions.sh

# Get username
username=$(whoami)

# Go packages directory
go_pkgs="/c/Users/$username/Packages/go"

# Function to get the latest version of a package from Go packages directory
getLatestVersion() {
    local package_name=$1
    local package_dir="$go_pkgs/$package_name"
    local latest_tag=""
    if [ -d "$package_dir" ]; then
        cd "$package_dir"
        latest_tag=$(git describe --tags --abbrev=0)
    fi
    echo "$latest_tag"
}

# Parse go.mod file and compare package versions
while IFS= read -r line; do
   if [[ $line =~ (require[[:space:]]+)?($currentGitHostUserPath/.+)[[:space:]]+([^[:space:]]+) ]]; then
        package_name=$(gawk -v repository="$currentGitHostUserPath" 'match($0, repository"/([^[:space:]]+)", arr) {print arr[1]}' <<< "$line")
        current_tag=$(gawk -v repository="$currentGitHostUserPath" 'match($0, repository"/([^[:space:]]+) v?([0-9]+\\.[0-9]+\\.[0-9]+)", arr) {print arr[2]}' <<< "$line")
        
        # If package_name and current_tag are not empty
        if [[ -n "$package_name" && -n "$current_tag" ]]; then
            # success "package=> [$package_name] current version: [$current_tag]"
            
            # Get the latest version from Go packages directory
            latest_tag=$(getLatestVersion "$package_name")
            # success "latest version: [$latest_tag]"
            
            if [ "$latest_tag" != "$current_tag" ]; then
                # Package has a different version, needs updating
                bash goget.sh $package_name $latest_tag
            fi
        fi
    fi
done < "go.mod"

successMessages
