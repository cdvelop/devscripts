#!/bin/bash
# A test script to demonstrate gorunscript functionality
# Usage: ./testScript.sh [error]

echo "Script executed successfully"
echo "Number of arguments: $#"
echo "Received arguments: $@"

# If "error" argument is received, return error code
if [[ "$1" == "error" ]]; then
  echo "ERROR: Requested error! Received 'error' argument"
  echo "Requested to finish with error"
  exit 1
fi

# Simple command execution
echo "Running command: echo 'This is an execute test'"
if echo 'This is an execute test'; then
  echo "Command executed successfully"
else
  echo "Command execution failed"
fi

exit 0