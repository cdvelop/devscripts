#!/bin/bash
# Description: Automated workflow for Go projects: checks modules, updates dependencies, creates tags, backs up and pushes to remote
# Usage: ./gopu.sh "Commit message"

source functions.sh

# Concatena los parámetros en una sola cadena
commit_message="$*"

bash gomodcheck.sh
if [ $? -eq 0 ]; then # Verificar si es 0

  bash pu.sh "$commit_message"
  if [ $? -eq 0 ]; then # Verificar el código de salida

    # actualizar los otros módulos donde este paquete es utilizado
    latest_tag=$(git describe --abbrev=0 --tags) # Obtén la última etiqueta
    
    #obtenemos el nombre del modulo go
    go_mod_name=$(gawk -v pattern=$currentGitHostUserPath/ 'NR==1 && match($0, pattern "([^/]+)", arr) { print arr[1] }' go.mod)
  
    bash gomodrename.sh "$go_mod_name" "$latest_tag"
    if [ $? -eq 0 ]; then # si es 0 realizamos backup
      source backupwindows.sh
    fi
    
  fi

fi


