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
    go_vet_dir=""
    # success "directory $go_mod_name"
    if [ -d "cmd" ]; then
      # success "$go_mod_name contains cmd folder"
      if [ -d "cmd/$go_mod_name" ]; then
          go_vet_dir="./cmd/$go_mod_name"
      else
          go_vet_dir="./cmd"
      fi
      execute "go vet $go_vet_dir" "go vet failed in $go_mod_name $go_vet_dir" "go vet $go_mod_name $go_vet_dir ok"
      
      if [ -n "$(find . -type f -name "*_test.go")" ]; then
          execute "go test" "Test errors found in $go_mod_name" "root test $go_mod_name ok"
      fi
      
      # Find folders containing 'test' in name and run tests
      test_folders=$(find -type d -name "*test*")
      for folder in $test_folders; do
          if [ -n "$(find $folder -type f -name "*_test.go")" ]; then
              execute "go test $folder" "Test errors found in folder $folder in $go_mod_name" "$go_mod_name folder: $folder test ok"
          fi
      done
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
