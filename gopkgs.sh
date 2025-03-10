#!/bin/bash
# Check if Go packages directory exists in current user's home
# Usage: ./gopkgs.sh
# This script verifies if the Go packages directory exists in the user's home directory

username=$(whoami)
go_pkgs="/c/Users/$username/Packages/go"
echo "Go packages directory: $go_pkgs"

if [ -d $go_pkgs ]; then
   echo "exists: $go_pkgs"
fi
