#!/bin/bash
# Updates a Go package to its latest tagged version
# Usage: ./goget.sh package-name

source functions.sh

pkg_name=$1

pkg_path=$go_pkgs/$pkg_name
# warning "Package directory to update: $pkg_path"

current_dir=$(pwd)

# Go to package directory to update
cd "$pkg_path"

tag_version=$(git describe --tags --abbrev=0)
# warning "$pkg_name local tag version: $tag_version"

# Return to original path
cd "$current_dir"

execute "go get $currentGitHostUserPath/$1@$tag_version" "failed to get package $pkg_name" "$pkg_name updated successfully"

execute "go mod tidy" "go mod tidy failed" "go mod tidy ok"

successMessages
exit 0