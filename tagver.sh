#!/bin/bash
# Description: Compare local and remote git tag versions
# Usage: ./tagver.sh

source functions.sh

# Check for pending changes
if [ -n "$(git status --porcelain)" ]; then
   warning "There are pending changes, performing commit..."
   
   bash pu.sh
fi

# Get local tag version
local_version=$(git describe --tags --abbrev=0)
success "local_version: $local_version"

# Get remote tag version
remote_version=$(git ls-remote --tags origin | awk '{print $2}' | cut -d '/' -f 3 | sort -V | tail -n 1)
success "remote_version: $remote_version"

# Compare versions
if [ "$local_version" != "$remote_version" ]; then
 warning "Versions are different. Push is needed."
    # bash pu.sh 
else
 success "Versions are equal. No action needed."
fi

successMessages
exit 0