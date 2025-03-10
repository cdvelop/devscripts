#!/bin/bash
# Updates the version tag of a Go module in go.mod file
# Usage: ./taggo.sh <package_name>

pkg_updated=$1
new_tag="v10.0.0"
go_mod_file="go.mod"
# Common package address prefix
pkg_common="github.com/cdvelop"
old_tag=$(gawk -v package="$pkg_updated" -v common="$pkg_common" 'match($0, "^require[[:space:]]+" common "/" package "[[:space:]]+([^[:space:]]+)", tag) {print tag[1]; exit} $1==common "/" package {print $2}' "$go_mod_file")
echo "=> Updating package: test module: $pkg_updated previous version: $old_tag new version: $new_tag"
