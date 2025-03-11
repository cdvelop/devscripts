#!/bin/bash
# Description: Utility functions for GitHub repository management and user information retrieval
# Usage: source githubutils.sh

# expected eg: juanin654
gitHubOwner=$(gh api user --jq .login)


