#!/bin/bash
# Mass rename multiple git tags using a file
# Usage: ./tagallrename.sh <filename>
# File format: each line should contain:
# <old_tag_name> <new_tag_name>
#
# The script performs the following operations for each line:
# 1. Creates new tag pointing to the old tag's commit
# 2. Deletes old tag locally
# 3. Deletes old tag from remote repository
# 4. Pushes all tags to remote repository

filename="$1"
git fetch --all

while read -r old_tag new_tag; do
    git tag "$new_tag" "$old_tag"
    git tag -d "$old_tag"
    git push origin :refs/tags/"$old_tag"
    git push --tags
done < "$filename"