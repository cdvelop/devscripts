#!/bin/bash
# Lists git tags with their commit messages, sorted by date
# Usage: ./tags.sh

git for-each-ref --format='%(refname:short) %(subject)' --sort=-taggerdate refs/tags --sort=committerdate
