#!/bin/bash
# Description: Checks and updates Go module dependencies, runs tests and performs data race detection
# Usage: ./gomodcheck.sh

source functions.sh
source syscall.sh

if [ -f "go.mod" ]; then
  # Get Go module name
  go_mod_name=$(gawk -v pattern=$currentGitHostUserPath/ 'NR==1 && match($0, pattern "([^/]+)", arr) { print arr[1] }' go.mod)
  
  # Update current packages
  execute "go get -u all" "failed to update packages in $go_mod_name" 
  
  # Verify code
  execute "go mod tidy" "go mod tidy failed in $go_mod_name" "go mod tidy $go_mod_name ok"
    syscall $go_mod_name
  syscall_found=$?
    if [ $syscall_found -eq 0 ]; then
    # Run gotest.sh for comprehensive testing and badge updates
    gotest.sh
    gotest_exit_code=$?
    
    # Return gotest exit code to preserve error status
    if [ $gotest_exit_code -ne 0 ]; then
      exit $gotest_exit_code
    fi
  fi
  
  # Run data race detection with go run
  race_output=$(go run -race "$go_mod_name" 2>&1)
  if [[ $race_output == *"WARNING: DATA RACE"* ]]; then
    warning "$race_output"
  fi
  
  successMessages
  exit 0
fi

exit 1