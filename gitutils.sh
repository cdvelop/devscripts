#!/bin/bash
# Git utilities for repository initialization and management
# Usage: source gitutils.sh && init_new_repo "my-project" "github.com/username"

source functions.sh
source fileIssues.sh




# Create README.md file
create_readme() {
    local current_folder=$1

    execute "echo '# $current_folder' > README.md" \
        "Failed to create README.md" \
        "README.md created" || return $?

    return 0
}



# Initialize files for new repository
init_base_files() {
    local current_folder=$1

    create_readme "$current_folder" || return $?
    create_issue_md_file || return $?

    return 0
}

# Initialize new git repository
init_new_repo() {
    local current_folder=$1
    local remote_url=$2

    if [ -d ".git" ]; then
        warning "Directory already initialized with Git: $current_folder"
        return 1
    fi

    execute "git init" \
        "Failed to initialize git" \
        "Git repository initialized" || return $?

    execute "git branch -M main" \
        "Failed to rename branch" || return $?

    return 0
}

# Configure remote for existing or new repository
setup_git_remote() {
    local current_folder=$1
    local remote_url=$2

    execute "git remote add origin https://$remote_url/$current_folder.git" \
        "Failed to add remote" \
        "Remote added successfully" || return $?

    return 0
}

# Push changes to remote
push_to_remote() {
    execute "git push -u origin main" \
        "Failed to push to remote" \
        "Pushed to remote" || return $?

    return 0
}

# Create and push a tag
create_git_tag() {
    local tag_name=${1:-"v0.0.1"}

    execute "git tag $tag_name" \
        "Failed to create tag" \
        "Tag created: $tag_name" || return $?

    execute "git push origin $tag_name" \
        "Failed to push tag" \
        "Tag pushed to remote" || return $?

    return 0
}

# Create initial commit
create_initial_commit() {
    execute "git add ." \
        "Failed to stage files" || return $?

    execute "git commit -m 'Initial commit'" \
        "Failed to commit" \
        "Files committed" || return $?

    return 0
}