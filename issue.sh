#!/bin/bash
# Description: Script to manage GitHub issues
# Usage: source ghissue.sh ej: ./pu.sh "Commit message #123 close" o ./pu.sh "Commit message #123 create"
#        or run directly: ./ghissue.sh create "Issue title" or ./ghissue.sh close 123

# Aseguramos que la CLI de GitHub esté disponible
check_gh_cli() {
    if ! command -v gh &> /dev/null; then
        echo "Error: GitHub CLI (gh) no está instalado."
        return 1
    fi
    return 0
}

# Cierra un issue por su número
close_issue() {
    local issue_number=$1
    if [[ -z "$issue_number" ]]; then
        echo "Error: Número de issue no proporcionado."
        return 1
    fi
    
    check_gh_cli || return 1
    
    echo "Cerrando issue #$issue_number..."
    gh issue close "$issue_number" || {
        echo "Error al cerrar el issue #$issue_number."
        return 1
    }
    echo "Issue #$issue_number cerrado correctamente."
    return 0
}

# Crea un issue con el título proporcionado
create_issue() {
    local title="$1"
    if [[ -z "$title" ]]; then
        echo "Error: Título del issue no proporcionado."
        return 1
    fi
    
    check_gh_cli || return 1
    
    echo "Creando issue: '$title'..."
    local issue_url
    issue_url=$(gh issue create --title "$title" --body "Issue creado automáticamente desde el script pu.sh")
    
    if [[ $? -ne 0 ]]; then
        echo "Error al crear el issue con título: $title"
        return 1
    fi
    
    echo "Issue creado correctamente: $issue_url"
    return 0
}

# Función para extraer información de issues del mensaje de commit
parse_issue_command() {
    local commit_message="$1"
    local issue_pattern='#([0-9]+)[ ]*(close|closes|closed)'
    local create_pattern='#[ ]*(create|new)'
    
    if [[ $commit_message =~ $issue_pattern ]]; then
        local issue_number="${BASH_REMATCH[1]}"
        echo "Detectado comando para cerrar issue #$issue_number"
        close_issue "$issue_number"
    elif [[ $commit_message =~ $create_pattern ]]; then
        # Extraer el título del issue (todo lo que está antes de #)
        local title=$(echo "$commit_message" | sed -n 's/\(.*\)#[ ]*\(create\|new\).*/\1/p' | xargs)
        if [ -n "$title" ]; then
            echo "Detectado comando para crear issue con título: '$title'"
            create_issue "$title"
        else
            echo "Error: No se pudo extraer un título para el nuevo issue."
        fi
    fi
}

# Función para mostrar ayuda
show_help() {
    echo "Uso:"
    echo "  Como script: ./ghissue.sh <comando> [argumentos]"
    echo "  Como fuente: source ghissue.sh"
    echo ""
    echo "Comandos:"
    echo "  create \"Título del issue\"   - Crea un nuevo issue con el título especificado"
    echo "  close NUMERO               - Cierra el issue con el número especificado"
    echo "  parse \"Mensaje de commit\"   - Prueba la función parse_issue_command con un mensaje"
    echo "  help                       - Muestra esta ayuda"
}

# Ejecutar como script independiente si se llama directamente
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    # El script se está ejecutando directamente, no como fuente
    if [ $# -lt 1 ]; then
        show_help
        exit 1
    fi

    comando="$1"
    shift

    case "$comando" in
        create)
            if [ -z "$1" ]; then
                echo "Error: Debe proporcionar un título para el issue."
                show_help
                exit 1
            fi
            create_issue "$1"
            ;;
        close)
            if [ -z "$1" ]; then
                echo "Error: Debe proporcionar un número de issue."
                show_help
                exit 1
            fi
            close_issue "$1"
            ;;
        parse)
            if [ -z "$1" ]; then
                echo "Error: Debe proporcionar un mensaje para analizar."
                show_help
                exit 1
            fi
            parse_issue_command "$1"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            echo "Error: Comando desconocido '$comando'"
            show_help
            exit 1
            ;;
    esac
fi
