#!/bin/bash
# Performs backup operations using FreeFileSync on Windows systems
# Usage: ./autoBackup.sh

source functions.sh

autoBackup() {
  case "$OSTYPE" in
    msys*|mingw*)
      warning "FreeFileSync backup started...."
      (execute '"/c/Program Files/FreeFileSync/FreeFileSync.exe" "/c/Users/$(whoami)/SyncWin/SyncSettings.ffs_batch"' &) >/dev/null 2>&1
      ;;
    linux*)
      warning "FreeFileSync backup started on Linux...."
      freefilesync_path="$(command -v FreeFileSync || command -v freefilesync)"
      if [ -z "$freefilesync_path" ]; then
        error "FreeFileSync is not installed or not in PATH."
        return 1
      fi
      sync_file="$HOME/Own/Sync/SyncSettings.ffs_batch"
      "$freefilesync_path" "$sync_file"
      ;;
    *)
      error "This operating system '$OSTYPE' is not compatible with the backup function."
      ;;
  esac
}

autoBackup