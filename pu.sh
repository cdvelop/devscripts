#!/bin/bash
# Description: Script to commit changes, create a new tag, and push to remote
# Usage: ./pu.sh "Commit message"
source functions.sh
source fileIssues.sh
source issue.sh # Incluir script para manejar issues de GitHub

current_folder=$(basename "$(pwd)")

# Accept positional args: first is commit message (quoted), second optional is tag
commit_message="$1"
provided_tag="$2"

# Analiza el mensaje para detectar comandos de issues ANTES de modificarlo
parse_issue_command "$commit_message"

# Obtiene el mensaje de commit desde ISSUE.md si existe
commit_message=$(get_commit_message_from_issue_md "$commit_message")

# Si commit_message está vacío después de la función, asigna un valor predeterminado
if [ -z "$commit_message" ]; then
    commit_message="auto update package"
fi

# Agrega cambios al índice
execute "git add ." "Error al añadir cambios a Git $current_folder." "cambios $current_folder añadidos"

# Realiza el commit solo si hay cambios
if git diff-index --quiet HEAD --; then
    echo "No hay cambios que commitear."
else
    execute "git commit -m '$commit_message'" "Error al crear el nuevo commit $current_folder."
fi

# Si se pasó una etiqueta por parámetro, usarla y saltarse la lógica automática
if [ -n "$provided_tag" ]; then
    new_tag="$provided_tag"
else
    # Obtén la última etiqueta
    latest_tag=$(git describe --abbrev=0 --tags 2>/dev/null)

    if [ -z "$latest_tag" ]; then
        # Si no existe ninguna etiqueta, establece la etiqueta inicial en v0.0.1
        new_tag="v0.0.1"
    else
        # Extrae el número de la etiqueta
        last_number=$(echo "$latest_tag" | grep -oE '[0-9]+$')

        # Incrementa el número en uno
        next_number=$((last_number + 1))

        # Construye la nueva etiqueta
        new_tag=$(echo "$latest_tag" | sed "s/$last_number$/$next_number/")
        
        # Verifica si la etiqueta ya existe
        while git rev-parse "$new_tag" >/dev/null 2>&1; do
            # Si la etiqueta existe, incrementar el número nuevamente
            next_number=$((next_number + 1))
            new_tag=$(echo "$latest_tag" | sed "s/$last_number$/$next_number/")
            echo "La etiqueta ya existe, probando con $new_tag"
        done
    fi
fi

execute "git tag $new_tag" "Error al crear la nueva etiqueta $current_folder." "nueva etiqueta $new_tag"

# Verifica si la rama tiene upstream configurado
branch=$(git symbolic-ref --short HEAD)
upstream=$(git rev-parse --symbolic-full-name --abbrev-ref @{u} 2>/dev/null)

if [ -z "$upstream" ]; then
    # Si no hay upstream, configura el upstream y haz push
    execute "git push --set-upstream origin $branch && git push origin $new_tag" \
    "Error al empujar los cambios y la nueva etiqueta a remoto $current_folder." \
    "Commit y Push $current_folder..."
else
    # Si ya tiene upstream, simplemente haz push
    execute "git push && git push origin $new_tag" \
    "Error al empujar los cambios y la nueva etiqueta a remoto $current_folder." \
    "Commit y Push $current_folder..."
fi

# Imprimir los mensajes acumulados
successMessages
deleteChangesIssueFile
exit 0
