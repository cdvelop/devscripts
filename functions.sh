#!/bin/bash
# Helper functions for git and script execution management
# Usage: source functions.sh

# currentGitHostUserPath expected eg: github.com/your-user
currentGitHostUserPath=$(git config --get remote.origin.url | sed -E 's#(git@|https://)([^:/]+)[/:]([^/]+)/.*#\2/\3#')


# Variable to store success messages
message=""

# user name expected eg: Juanin
username=$(whoami)

# Look for and update packages in "\Packages\go" directory
go_pkgs="/c/Users/$username/Packages/go"

# Function to display a success message
success() {
  echo -e "\033[0;32m$1\033[0m" # green color
}

# Function to display a warning message
warning() {
  echo -e "\033[0;33m$1\033[0m" # yellow color
}

# Function to display an error message
error() {
  echo -e "\033[0;31mError: $1 $2\033[0m" # red color
}

# Function to perform an action and show error message on failure
execute() {
 output=$(eval "$1" 2>&1)
 local exit_code=$?
  if [ $exit_code -ne 0 ]; then
    error "$2" "$output"
    if [ -z "$4" ]; then
      # warning "fourth parameter [no exist] not sent."
      exit 1
    fi
  else
    # Concatenate success message to message variable if provided
    if [ -n "$3" ]; then
      addOKmessage "$3"
    fi
  fi
  return $exit_code
}

addOKmessage(){
  if [ -n "$1" ]; then
      symbol="\033[0;33m=>OK\033[0m"  # Orange symbol
      text="\033[0;32m$1\033[0m"      # Green text
      message+="\n$symbol $text"       # Concatenate success message with symbol and text
  fi
}

addERRORmessage(){
  if [ -n "$1" ]; then
      symbol="\033[0;31m!ERROR!\033[0m"  # Red symbol
      text="\033[0;31m$1\033[0m"         # Red text
      message+="\n$symbol $text"
  fi
}

# Print accumulated messages
successMessages(){
  echo -e "$message"
  message=""
}

# Function to check if "changes.txt" exists and clear its content
function deleteChangesFileContent() {
    if [ -f "changes.txt" ] && [ -s "changes.txt" ]; then
        # Clear file content
        echo "" > changes.txt
    fi
}
