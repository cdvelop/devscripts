#!/bin/bash
# Description: Script to manage GitHub issues using functions.sh helpers
# Usage: ./issue.sh <command> [args] eg: []./issue.sh + "My issue" bug] | - 4 "Closed by xxx" | ?

# Source helper functions
source functions.sh

# Aseguramos que la CLI de GitHub esté disponible
check_gh_cli() {
    # Use execute to check command existence silently, add error if fails
    execute "command -v gh &> /dev/null" "GitHub CLI (gh) no está instalado." "" "no_exit"
    # Check the return code of the last execute command
    if [ $? -ne 0 ]; then
        return 1
    fi
    return 0
}

# Cierra un issue por su número
close_issue() {
    local issue_number=$1
    local close_message=$2  # Optional message parameter
    if [[ -z "$issue_number" ]]; then
        addERRORmessage "Número de issue no proporcionado para cerrar."
        return 1
    fi

    check_gh_cli || return 1 # Stop if gh cli check fails

    # Build command based on whether a message was provided
    local command="gh issue close \"$issue_number\""
    if [[ -n "$close_message" ]]; then
        command="$command --comment \"$close_message\""
    fi

    execute "$command" \
            "Error al cerrar el issue #$issue_number." \
            "Issue #$issue_number cerrado correctamente."
    return $? # Return the exit code of the execute command
}

# Crea un issue con el título proporcionado y opcionalmente añade etiquetas
create_issue() {
    local title="$1"
    local labels_str="$2" # Optional second argument for comma-separated labels
    local issue_url=""

    if [[ -z "$title" ]]; then
        addERRORmessage "Título del issue no proporcionado para crear."
        return 1
    fi

    check_gh_cli || return 1 # Stop if gh cli check fails

    # Create the issue first and capture the URL from stdout on success
    # We need to run this slightly differently than the standard 'execute' to capture the URL
    issue_url=$(gh issue create --title "$title" --body "Issue creado automáticamente desde script." 2>&1)
    local exit_code=$?

    if [ $exit_code -ne 0 ]; then
        # If creation failed, add error message using the captured output (which is the error message from gh)
        addERRORmessage "Error al crear el issue con título '$title': $issue_url"
        return $exit_code
    else
        # If creation succeeded, add success message
        addOKmessage "Issue '$title' creado correctamente: $issue_url"

        # Now, if labels were provided, add them
        if [[ -n "$labels_str" ]]; then
            # Save and restore IFS to handle labels with spaces if needed, though gh cli might handle quotes
            local OLD_IFS="$IFS"
            IFS=','
            local labels_array=($labels_str) # Split string into array by comma
            IFS="$OLD_IFS"

            addOKmessage "Añadiendo etiquetas: ${labels_array[*]}" # Show labels being added

            for label in "${labels_array[@]}"; do
                # Trim potential leading/trailing whitespace from label
                label=$(echo "$label" | xargs)
                if [[ -n "$label" ]]; then # Ensure label is not empty after trimming
                    # Use execute to add each label, referencing the issue by URL
                    execute "gh issue edit \"$issue_url\" --add-label \"$label\"" \
                            "Error al añadir etiqueta '$label' al issue $issue_url." \
                            "Etiqueta '$label' añadida."
                    # Optional: could check exit code here and break if one label fails
                fi
            done
        fi
        return 0 # Overall success
    fi
}

# Lista los issues del repositorio actual
list_issues() {
    check_gh_cli || return 1 # Stop if gh cli check fails

    # Execute gh issue list directly to allow output to pass through
    gh issue list
    local exit_code=$? # Capture exit code of gh issue list

    if [ $exit_code -ne 0 ]; then
        addERRORmessage "Error al listar los issues."
    else
        # Add a success message only if the command succeeded,
        # but the main output is the list itself printed above.
        addOKmessage "Listado de issues mostrado."
    fi
    return $exit_code
}

# Muestra un issue por su número
view_issue() {
    local issue_number=$1
    if [[ -z "$issue_number" ]]; then
        addERRORmessage "Número de issue no proporcionado para ver."
        return 1
    fi

    check_gh_cli || return 1 # Stop if gh cli check fails

    # Execute gh issue view directly to allow output to pass through
    gh issue view "$issue_number"
    local exit_code=$? # Capture exit code of gh issue view

    if [ $exit_code -ne 0 ]; then
        addERRORmessage "Error al ver el issue #$issue_number."
    fi
    return $exit_code
}

# Edita interactivamente el cuerpo de un issue usando Notepad
edit_issue_interactive() {
    local issue_number=$1
    local temp_file=""
    local original_body=""
    local modified_body=""

    if [[ -z "$issue_number" ]]; then
        addERRORmessage "Número de issue no proporcionado para editar."
        return 1
    fi

    check_gh_cli || return 1

    # Crear archivo temporal seguro y preparar limpieza automática
    temp_file=$(mktemp --suffix=.md)
    if [ $? -ne 0 ] || [ -z "$temp_file" ]; then
         addERRORmessage "No se pudo crear archivo temporal para la edición."
         return 1
    fi
    # Asegurar que el archivo temporal se borre al salir de la función (éxito, error o interrupción)
    trap 'rm -f "$temp_file"' RETURN

    addOKmessage "Obteniendo cuerpo actual del issue #$issue_number..."
    # Usar gh issue view para obtener solo el cuerpo
    original_body=$(gh issue view "$issue_number" --json body --jq .body)
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        addERRORmessage "Error al obtener el cuerpo del issue #$issue_number. ¿Existe el issue?"
        # El trap limpiará el archivo temporal
        return $exit_code
    fi

    # Guardar el cuerpo original en el archivo temporal
    printf "%s" "$original_body" > "$temp_file"

    addOKmessage "Abriendo issue #$issue_number en VS Code. Guarda los cambios y cierra la pestaña para continuar." # Mensaje actualizado
    # Usar 'code -r -w' para abrir en la ventana existente y esperar
    code -r -w "$temp_file"
    # Usar 'code --wait' para abrir en VS Code y esperar a que se cierre la pestaña
    #code --wait "$temp_file"
    local editor_exit_code=$? # Capturar código de salida del editor (opcional, puede no ser fiable)

    # Verificar si el editor se cerró correctamente (código 0)
    # Aunque --wait funciona, el código de salida puede no siempre indicar si se guardó.
    # La comparación de contenido sigue siendo la forma más fiable.
    if [ $editor_exit_code -ne 0 ]; then
        warning "VS Code se cerró con un código de error ($editor_exit_code). Verificando cambios de todos modos."
    fi

    # Leer el contenido modificado del archivo temporal
    # Usar mapfile (readarray) es más seguro para leer archivos con saltos de línea
    mapfile -t body_lines < "$temp_file"
    modified_body=$(printf "%s\n" "${body_lines[@]}")
    # Quitar el último salto de línea que printf añade si el archivo no terminaba en uno
    if [[ $(wc -c < "$temp_file") -gt 0 && $(tail -c1 "$temp_file" | wc -l) -eq 0 ]]; then
         modified_body=${modified_body%?}
    fi


    # Comparar si hubo cambios
    if [[ "$original_body" == "$modified_body" ]]; then
        addOKmessage "El contenido del issue #$issue_number no ha cambiado. No se requiere actualización."
        # El trap limpiará el archivo temporal
        return 0
    fi

    addOKmessage "Actualizando el cuerpo del issue #$issue_number en GitHub..."
    # Usar --body-file es más robusto que --body "$variable" para contenido largo o complejo
    execute "gh issue edit \"$issue_number\" --body-file \"$temp_file\"" \
            "Error al actualizar el issue #$issue_number." \
            "Issue #$issue_number actualizado correctamente."

    # El trap limpiará el archivo temporal al retornar
    return $?
}

# Función para extraer información de issues del mensaje de commit
parse_issue_command() {
    local commit_message="$1"
    # addOKmessage "Analizando mensaje para comandos de issue: '$commit_message'" # Optional: uncomment for debugging

    # Pattern to capture close command and issue number
    local issue_close_pattern='(close|closes|closed).*#([0-9]+)|#([0-9]+).*(close|closes|closed)'
    # Pattern to capture title, create command, and optional labels (comma-separated, no spaces around commas assumed for simplicity)
    local create_pattern='(.*)#(create|new)[ ]*([^#]*)' # Group 1: Title, Group 2: Command, Group 3: Labels

    local issue_number=""
    local action=""
    local title=""
    local labels=""
    local title=""

    if [[ $commit_message =~ $issue_close_pattern ]]; then
        if [[ -n "${BASH_REMATCH[2]}" ]]; then issue_number="${BASH_REMATCH[2]}"; fi
        if [[ -n "${BASH_REMATCH[3]}" ]]; then issue_number="${BASH_REMATCH[3]}"; fi
        action="close"
        # addOKmessage "Patrón de cierre coincide. Issue: #$issue_number" # Optional debug
        close_issue "$issue_number"
    elif [[ $commit_message =~ $create_pattern ]]; then
        title=$(echo "${BASH_REMATCH[1]}" | xargs) # Trim whitespace title
        action="${BASH_REMATCH[2]}"
        labels=$(echo "${BASH_REMATCH[3]}" | xargs) # Trim whitespace labels string

        if [[ -n "$title" ]]; then
            # addOKmessage "Patrón de creación coincide. Título: '$title', Etiquetas: '$labels'" # Optional debug
            create_issue "$title" "$labels" # Pass title and labels string
        else
            addERRORmessage "No se pudo extraer un título para el nuevo issue desde '$commit_message'."
        fi
    else
        # Only show message if called directly with 'parse' command, not when sourced
         if [[ "${BASH_SOURCE[0]}" == "${0}" ]] && [[ "$current_command" == "parse" ]]; then
             addERRORmessage "Ningún patrón (cerrar o crear) coincidió para el mensaje: '$commit_message'"
         fi
    fi
}

# Función para mostrar ayuda (uses standard echo)
show_help() {
    echo "Uso:"
    echo "  Como script: ./issue.sh <comando> [argumentos]"
    echo "  Como fuente: source issue.sh (usado por pu.sh)"
    echo ""
    echo "Comandos:"
    echo "  (sin argumentos)                - Lista los issues del repositorio actual"
    echo "                                 Ej: ./issue.sh"
    echo "  list | l                        - Lista los issues del repositorio actual"
    echo "                                 Ej: ./issue.sh l"
    echo "  create | + \"Título\" [etiquetas] - Crea issue. Añade etiquetas opcionales."
    echo "                                 Ej: ./issue.sh + \"Mi issue\" bug"
    echo "  close | - NUMERO [mensaje]      - Cierra el issue con el número especificado"
    echo "                                 y un mensaje opcional."
    echo "                                 Ej: ./issue.sh - 4 \"Cerrado por xxx\""
    echo "  NUMERO                          - Muestra detalles del issue"
    echo "                                 Ej: ./issue.sh 123"
    echo "  parse \"Mensaje\"                 - Prueba la función parse_issue_command con un mensaje"
    echo "  ? | help | h                    - Muestra esta ayuda"
    echo ""
    echo "Nota: Cuando se usa como fuente (source), solo se exponen las funciones."
    echo "      Los mensajes de éxito/error se manejan a través de functions.sh."
}

# Ejecutar como script independiente si se llama directamente
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    current_command="" # Variable to track the command for context

    # Check for #<number> pattern first
    # Check for number first
    if [[ "$1" =~ ^([0-9]+)$ ]]; then
        view_issue "$1"
        # Print accumulated success/error messages
        successMessages
        # Exit with 1 if there were errors
        if [[ "$message" == *"ERROR"* ]]; then
            exit 1
        else
            exit 0
        fi
    fi

    # If no arguments provided, list issues instead of showing help
    if [ $# -lt 1 ]; then
        current_command="list"
        list_issues
        successMessages
        if [[ "$message" == *"ERROR"* ]]; then
            exit 1
        else
            exit 0
        fi
    fi

    comando="$1"
    current_command="$comando" # Store command
    shift

    case "$comando" in
        create|+)
            if [ -z "$1" ]; then
                addERRORmessage "Debe proporcionar un título para el issue."
                show_help
                successMessages # Print accumulated messages before exiting
                exit 1
            fi
            title="$1"
            shift
            labels="$1"  # Optional label parameter
            create_issue "$title" "$labels"
            ;;
        close|-)
            if [ -z "$1" ]; then
                addERRORmessage "Debe proporcionar un número de issue."
                show_help
                successMessages # Print accumulated messages before exiting
                exit 1
            fi
            issue_number="$1"
            shift
            close_message="$1"  # Optional closing message
            close_issue "$issue_number" "$close_message"
            ;;
        edit|e) # Añadir caso para editar interactivamente
            if [ -z "$1" ]; then
                addERRORmessage "Debe proporcionar un número de issue para editar."
                show_help
                successMessages
                exit 1
            fi
            # Verificar que no se pasen argumentos extra (este modo solo edita el cuerpo)
            if [ $# -gt 1 ]; then
                 addERRORmessage "El comando 'edit' interactivo solo acepta el número de issue."
                 addOKmessage "Para editar otros campos (título, etiquetas, etc.), usa 'gh issue edit <num> [flags...]' directamente."
                 successMessages
                 exit 1
            fi
            issue_number="$1"
            edit_issue_interactive "$issue_number"
            ;;
        parse)
            if [ -z "$1" ]; then
                addERRORmessage "Debe proporcionar un mensaje para analizar."
                show_help
                successMessages # Print accumulated messages before exiting
                exit 1
            fi
            parse_issue_command "$1"
            ;;
        list|l)
            list_issues
            ;;
        \?|help|h)
            show_help
            # No messages to print for help
            exit 0
            ;;
        *)
            addERRORmessage "Comando desconocido '$comando'"
            show_help
            successMessages # Print accumulated messages before exiting
            exit 1
            ;;
    esac

    # Print accumulated success/error messages at the end of direct execution
    successMessages
    # Exit with 1 if there were errors (message variable contains ERROR symbol)
    if [[ "$message" == *"ERROR"* ]]; then
        exit 1
    else
        exit 0
    fi
fi
