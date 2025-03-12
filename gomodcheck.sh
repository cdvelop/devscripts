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
    # Run go vet on the whole module
    execute "go vet ./..." "go vet failed in $go_mod_name" "go vet $go_mod_name ok"
    
    # Check if any test files exist in the module
    if [ -n "$(find . -type f -name "*_test.go")" ]; then
      # Run all tests in the module and its subdirectories
      execute "go test ./..." "Test errors found in $go_mod_name" "All tests in $go_mod_name passed"
      
      # Optional: Run tests with race detection
      execute "go test -race ./..." "Race condition tests failed in $go_mod_name" "Race condition tests in $go_mod_name passed"
    else
      info "No test files found in $go_mod_name"
    fi
    
    # If cmd directory exists, handle it specifically
    if [ -d "cmd" ]; then
      if [ -d "cmd/$go_mod_name" ]; then
          go_vet_dir="./cmd/$go_mod_name"
      else
          go_vet_dir="./cmd"
      fi
      execute "go vet $go_vet_dir" "go vet failed in $go_mod_name $go_vet_dir" "go vet $go_mod_name $go_vet_dir ok"
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