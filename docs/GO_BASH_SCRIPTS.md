# ISSUE: Go Scripts Architecture with Bash Execution

## Problem
Need executable Go scripts that work as terminal commands while maintaining reusable code and organized structure.

## Architecture

```
devscripts/
├── go.mod                    # Module: github.com/cdvelop/devscripts
├── args.go                   # Reusable argument handler
├── gocurrentdir.sh          # Universal script launcher
├── sectionUpdate.go         # Exportable function SectionUpdate(args ...string)
├── sectionUpdate.sh         # Wrapper script for execution
└── cmd/
    └── sectionUpdate.go     # Executable using devscripts.SectionUpdate
```

## Components

### 1. Universal Script Launcher (`gocurrentdir.sh`)
```bash
#!/bin/bash
# Universal launcher that automatically detects calling script and runs corresponding Go command
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CURRENT_DIR="$(pwd)"
CALLING_SCRIPT="$(basename "${BASH_SOURCE[1]}" .sh)"
cd "$SCRIPT_DIR"
go run "cmd/${CALLING_SCRIPT}.go" "$CURRENT_DIR" "$@"
```

**Key Features:**
- Automatically detects the calling script name using `${BASH_SOURCE[1]}`
- Removes `.sh` extension to match Go executable name
- Eliminates need to hardcode script names in wrappers
- Single file pattern for all Go script wrappers

### 2. Argument Handler (`args.go`)
```go
package devscripts

func ExecuteWithArgs(fn func(...string)) {
    if len(os.Args) > 1 {
        os.Chdir(os.Args[1]) // First arg = working directory
    }
    fn(os.Args[2:]...) // Remaining args to function
}
```

### 3. Business Function (Root) EXAMPLE
```go
package devscripts

type MdHandler struct {
    sectionID  string
    afterLine  string
    content    string
    readmeFile string
}

func NewMdHandler(args ...string) *MdHandler {
    h := &MdHandler{readmeFile: "README.md"}
    if len(args) > 0 { h.sectionID = args[0] }
    if len(args) > 1 { h.afterLine = args[1] }
    if len(args) > 2 { h.content = args[2] }
    if len(args) > 3 { h.readmeFile = args[3] }
    return h
}

func (h *MdHandler) SectionUpdate() error {
    // Implementation: update section in markdown file
    fmt.Printf("Updating section %s in %s\n", h.sectionID, h.readmeFile)
    return nil
}

func SectionUpdate(args ...string) {
    handler := NewMdHandler(args...)
    handler.SectionUpdate()
}
```

### 4. Executable (`cmd/`)
```go
//go:build ignore
package main

import "github.com/cdvelop/devscripts"

func main() {
    devscripts.ExecuteWithArgs(devscripts.SectionUpdate)
}
```

### 5. Bash Wrapper
```bash
#!/bin/bash
# Description: Update sections in markdown files dynamically  
# Usage: ./sectionUpdate.sh section_identifier [after_line] new_content [file]
source "$(dirname "$0")/gocurrentdir.sh"
```

**Benefits:**
- Single line of code per wrapper
- Automatic script name detection
- No hardcoded references
- Consistent pattern across all scripts
- Windows compatibility by gitbash
- Testing support with go test

## LLM Instructions

When creating new Go scripts:

1. **Business function** (root): `func Name(args ...string)` with `package devscripts`
2. **Executable** (cmd/): Use `//go:build ignore` and `devscripts.ExecuteWithArgs(devscripts.Name)`
3. **Bash wrapper**: Simply add `source "$(dirname "$0")/gocurrentdir.sh"` after header comments
4. **First argument is always the execution directory**

### Benefits of This Pattern

- **Reduces boilerplate**: Single line wrapper implementation
- **Eliminates errors**: No need to manually update script names
- **Consistent**: Same pattern for all Go script wrappers
- **Maintainable**: Changes to launcher logic only need to be made in one place
- **Testable**: Easy to test Go code using standard Go testing framework instead of bash testing
