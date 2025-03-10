#!/bin/bash
# Performs backup operations using FreeFileSync on Windows systems
# Usage: ./backupwindows.sh

source functions.sh

backup() {
  case "$OSTYPE" in
    msys*|mingw*)
      warning "FreeFileSync backup started...."
      (execute '"/c/Program Files/FreeFileSync/FreeFileSync.exe" "/c/Users/$(whoami)/SyncWin/SyncSettings.ffs_batch"' &) >/dev/null 2>&1
      # Adding & at the end of the command will run it in background and free up the terminal
      ;;
    *)
      error "This operating system '$OSTYPE' is not compatible with the backup function."
      ;;
  esac
}

backup