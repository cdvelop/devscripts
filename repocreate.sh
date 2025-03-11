#!/bin/bash
# Creates a new GitHub repository with initial README and license files
# Usage: ./repocreate.sh my-repo "My description" [public|private]

###########################################
# Usage Examples:                        #
###########################################
#
# 1. Create private repository (default):
#    ./repocreate.sh my-repo "My description"
#
# 2. Create public repository:
#    ./repocreate.sh my-repo "My description" public
#
# 3. When importing in another script:
#    source ./repocreate.sh
#    create_repository "my-repo" "My description" "public"
###########################################

# Import functions from functions.sh
source functions.sh
source githubutils.sh
source licensecreate.sh

create_repository() {
    local repo_name=$1
    local description=$2
    local visibility=${3:-public}  # Default to public if not specified

    # Validate required arguments
    if [ -z "$repo_name" ] || [ -z "$description" ]; then
        error "Usage: create_repository <repo-name> <description> [visibility]"
        return 1
    fi

    # Validate visibility parameter
    if [ "$visibility" != "private" ] && [ "$visibility" != "public" ]; then
        error "Visibility must be either 'public' or 'private'"
        return 1
    fi

    # Create repository with specified visibility
    execute "gh repo create $repo_name --$visibility --description \"$description\"" \
        "Failed to create repository" \
        "Repository $repo_name created successfully as $visibility" || return $?

    execute "git clone https://github.com/$gitHubOwner/$repo_name.git" \
        "Failed to clone repository" \
        "Repository cloned successfully" || return $?

    cd $repo_name || return 1

    echo "# $repo_name" > README.md
    echo -e "\n$description" >> README.md

    generate_license_file "MIT" $gitUserName

    execute "git add README.md" \
        "Failed to stage README.md" || return $?

    execute "git commit -m 'Initial commit with README'" \
        "Failed to commit changes" \
        "Initial commit created" || return $?

    execute "git push -u origin main" \
        "Failed to push to remote" \
        "Successfully pushed to remote" || return $?

    cd ..
    return 0
}

# If the script is executed directly (not imported)
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    if [ "$#" -lt 2 ] || [ "$#" -gt 3 ]; then
        error "Usage: $0 <repo-name> <description> [visibility]"
        exit 1
    fi

    create_repository "$1" "$2" "$3"
    exit_code=$?
    successMessages
    exit $exit_code
fi