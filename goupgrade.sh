#!/bin/bash
# Description: Updates Go packages and tidies up module dependencies
# Usage: ./goupgrade.sh

# Update packages
go get -u -t ./...

# Run "go mod tidy" to update Go modules
go mod tidy