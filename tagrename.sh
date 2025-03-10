#!/bin/bash
# Description: Rename git tags both locally and remotely
# Usage: ./tagrename.sh <old_tag> <new_tag>

if [ "$#" -ne 2 ]; then
    echo "Error: Exactly two parameters are required"
    echo "Usage: $0 <old_tag> <new_tag>"
    exit 1
fi

old_tag=$1
new_tag=$2

git fetch --all
git tag "$new_tag" "$old_tag"
git tag -d "$old_tag"
git push origin :refs/tags/"$old_tag"
git push --tags