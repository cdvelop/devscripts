#!/bin/bash
# unify-git-author.sh
# Unifica el autor de todos los commits de la rama actual usando la configuración global de git
global_name=$(git config --global user.name)
global_email=$(git config --global user.email)

if [ -z "$global_name" ] || [ -z "$global_email" ]; then
  echo "[ERROR] No se encontró configuración global de nombre o email en git."
  exit 1
fi

echo "Usando nombre: $global_name"
echo "Usando email: $global_email"

echo "Reescribiendo historial de la rama actual..."
git filter-branch --env-filter "
    export GIT_AUTHOR_NAME=\"$global_name\"
    export GIT_AUTHOR_EMAIL=\"$global_email\"
    export GIT_COMMITTER_NAME=\"$global_name\"
    export GIT_COMMITTER_EMAIL=\"$global_email\"
" HEAD

if [ $? -ne 0 ]; then
  echo "[ERROR] Falló la reescritura del historial."
  exit 1
fi

echo "Forzando push al remoto..."
git push --force origin HEAD

echo "Listo. Todos los commits de la rama actual tienen el autor unificado."
ñ