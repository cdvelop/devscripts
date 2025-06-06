#!/bin/bash
# Description: Automated workflow for Go projects: checks modules, updates dependencies, creates tags, backs up and pushes to remote
# Usage: gopu.sh "Commit message"

# Function to show usage instructions
show_usage() {
    echo "================================================================================"
    echo "gopu.sh - Automated Go Development Workflow"
    echo "================================================================================"
    echo ""
    echo "DESCRIPTION:"
    echo "    This script automates the following sequence of actions:"
    echo "    1. Adds all changed files to the Git staging area (similar to 'git add .')"
    echo "    2. Runs Go tests (equivalent to 'go test ./...')"
    echo "    3. Checks for race conditions in the Go code (equivalent to 'go test -race ./...')"
    echo "    4. Commits the staged changes using the message provided as an argument"
    echo "    5. Pushes the commit to the remote repository"
    echo "    6. Creates and pushes a tag"
    echo ""
    echo "USAGE:"
    echo "    gopu.sh \"Your detailed commit message here\""
    echo ""
    echo "EXAMPLES:"
    echo "    gopu.sh \"feat: implement user authentication module\""
    echo "    gopu.sh \"fix: resolve memory leak in string processing\""
    echo "    gopu.sh \"docs: update README with installation instructions\""
    echo ""
    echo "ERROR: Commit message is required!"
    echo "================================================================================"
}

source functions.sh

# Check if commit message is provided
if [ $# -eq 0 ] || [ -z "$*" ]; then
    show_usage
    exit 1
fi

# Concatena los parámetros en una sola cadena
commit_message="$*"

bash gomodcheck.sh
if [ $? -eq 0 ]; then # Verificar si es 0

  bash pu.sh "$commit_message"
  if [ $? -eq 0 ]; then # Verificar el código de salida

    # actualizar los otros módulos donde este paquete es utilizado
    latest_tag=$(git describe --abbrev=0 --tags) # Obtén la última etiqueta
    
    #obtenemos el nombre del modulo go
    go_mod_name=$(gawk -v pattern=$currentGitHostUserPath/ 'NR==1 && match($0, pattern "([^/]+)", arr) { print arr[1] }' go.mod)
  
    bash gomodtagupdate.sh "$go_mod_name" "$latest_tag"
    if [ $? -eq 0 ]; then # si es 0 realizamos backup
      source backupwindows.sh
    fi
    
  fi

fi


