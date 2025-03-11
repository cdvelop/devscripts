# DevScripts Package
scripts commonly used by a developer in his daily workflow

<!-- SCRIPTS_SECTION_START -->
## Available Scripts
<small>This section is automatically generated.</small>

| Script Name            | Description                                                                                                           | Usage                          |
| ---------------------- | --------------------------------------------------------------------------------------------------------------------- | ------------------------------ |
| `backupwindows.sh`     | Performs backup operations using FreeFileSync on Windows systems                                                      | `./backupwindows.sh`           |
| `changeremote.sh`      | Script to change the remote URL of a Git repository                                                                   | `./changeremote.sh https://github.com/username/repository.git` |
| `delete.sh`            | Script to delete a file locally and track the deletion in Git                                                         | `./delete.sh filename.txt`     |
| `functions.sh`         | Helper functions for git and script execution management                                                              | `source functions.sh`          |
| `githubutils.sh`       | Utility functions for GitHub repository management and user information retrieval                                     | `source githubutils.sh`        |
| `gitremtracking.sh`    | Removes files from git tracking both locally and remotely                                                             | `./gitremtracking.sh file1.txt file2.txt` |
| `gitutils.sh`          | Git utilities for repository initialization and management                                                            | `source gitutils.sh && init_new_repo "my-project" "github.com/username"` |
| `goaddtest.sh`         | Script to generate Go test files with unit test and benchmark templates                                               | `./goaddtest.sh CreateFile create` |
| `goget.sh`             | Updates a Go package to its latest tagged version                                                                     | `./goget.sh package-name`      |
| `gomodcheck.sh`        | Checks and updates Go module dependencies, runs tests and performs data race detection                                | `./gomodcheck.sh`              |
| `gomodinit.sh`         | Script to initialize a Go module and create basic project structure                                                   | `./gomodinit.sh`               |
| `gomodrename.sh`       | Rename a Go module and update all its references                                                                      | `./gomodrename.sh old-module-name new-module-name` |
| `gomodulesupdate.sh`   | Updates Go module versions across all projects that use them                                                          | `./gomodulesupdate.sh <package-name> <new-version>` |
| `gomodutils.sh`        | Utility functions for managing Go modules and version updates                                                         | `source gomodutils.sh && update_single_go_module "mymodule" "v1.2.3"` |
| `gonewproject.sh`      | Creates a new Go project with standard directory structure and initial files, sets up remote repository               | `./gonewproject.sh <repo-name> <description> [visibility]` |
| `gopkgs.sh`            | Check if Go packages directory exists in current user's home                                                          | `./gopkgs.sh`                  |
| `gopkgupdate.sh`       | Updates Go packages in go.mod to their latest versions from local repositories                                        | `./gopkgupdate.sh`             |
| `gopu.sh`              | Automated workflow for Go projects: checks modules, updates dependencies, creates tags, backs up and pushes to remote | `./gopu.sh "Commit message"`   |
| `gorenameproject.sh`   | Script to rename a Go project and update its module references                                                        | `./gorenameproject.sh old-project-name new-project-name` |
| `goupgrade.sh`         | Updates Go packages and tidies up module dependencies                                                                 | `./goupgrade.sh`               |
| `licensecreate.sh`     | Create and commit a license file for a Git repository                                                                 | `./licensecreate.sh [license-type] [owner-name]` |
| `pu.sh`                | Script to commit changes, create a new tag, and push to remote                                                        | `./pu.sh "Commit message"`     |
| `rename.sh`            | Rename a file and update Git tracking                                                                                 | `./rename.sh <current_name> <new_name>` |
| `repocreate.sh`        | Creates a new GitHub repository with initial README and license files                                                 | `./repocreate.sh my-repo "My description" [public|private]` |
| `repodelete.sh`        | Deletes a remote GitHub repository after confirmation and permission checks                                           | `./repodelete.sh <repo-name> [force_delete] [owner]` |
| `repoexistingsetup.sh` | Setup additional files and tags for an existing Git repository                                                        | `./repoexistingsetup.sh`       |
| `repolocalinit.sh`     | Initialize a new local Git repository with basic files and remote setup                                               | `./repolocalinit.sh`           |
| `reporename.sh`        | Renames a repository both locally and on remote GitHub, updates Git remotes                                           | `./reporename.sh <old-name> <new-name>` |
| `syscall.sh`           | Check if a Go package uses syscall/js imports                                                                         | `./syscall.sh <package_name>`  |
| `tag.sh`               | Script to automatically increment the last number in a Git tag                                                        | `./tag.sh (will get the latest tag and suggest the next one)` |
| `tagalldelete.sh`      | Bulk delete git tags listed in a text file                                                                            | `./tagalldelete.sh <filename>` |
| `tagallrename.sh`      | Mass rename multiple git tags using a file                                                                            | `./tagallrename.sh <filename>` |
| `tagdelete.sh`         | Delete git tags locally and remotely                                                                                  | `tagdelete.sh tag1 tag2 tag3`  |
| `taggo.sh`             | Updates the version tag of a Go module in go.mod file                                                                 | `./taggo.sh <package_name>`    |
| `tagrename.sh`         | Rename git tags both locally and remotely                                                                             | `./tagrename.sh <old_tag> <new_tag>` |
| `tags.sh`              | Lists git tags with their commit messages, sorted by date                                                             | `./tags.sh`                    |
| `tagver.sh`            | Compare local and remote git tag versions                                                                             | `./tagver.sh`                  |
| `testScript.sh`        | A test script to demonstrate gorunscript functionality                                                                | `./testScript.sh [error]`      |
| `vpssetupbase.sh`      | Base VPS setup for Debian-based Linux servers                                                                         | `sudo ./vpssetupbase.sh <username> <ssh_key>` |
| `vpssetupsecurity.sh`  | VPS security setup script for Debian-based Linux servers                                                              | `sudo ./vpssetupsecurity.sh <username> <new_ssh_port>` |

<!-- SCRIPTS_SECTION_END -->

## Use with Go

This package provides functionality for executing different types of scripts with Go in a consistent way across different operating systems.

## ScriptRunner

The `ScriptRunner` handles the execution of scripts with appropriate interpreters based on file extensions.

### Basic Usage

```go
import "path/to/devscripts"

// Create a new script runner with default directory (current working directory)
runner := devscripts.NewScriptRunner()

// Or with a specific scripts directory
runner := devscripts.NewScriptRunner("/path/to/scripts")

// Execute a script with arguments
exitCode, output, err := runner.ExecScript("myscript.sh", "arg1", "arg2")
if err != nil {
    // Handle error
    fmt.Printf("Script execution failed: %v\n", err)
    return
}

fmt.Printf("Script executed with exit code %d\n", exitCode)
fmt.Printf("Output:\n%s\n", output)
```

### Chained Script Execution

The package supports chained script execution, where scripts run sequentially and execution stops if any script fails:

```go
// Create a script chain
exitCode, output, err := runner.Chain().
    Then("script1.sh", "arg1").
    Then("script2.py", "arg2", "arg3").
    Then("script3.sh").
    Execute()

if err != nil {
    fmt.Printf("Chain execution failed: %v\n", err)
    return
}

fmt.Printf("All scripts executed successfully with combined output:\n%s\n", output)
```

### Accessing Individual Results

You can also access the results of the last executed script in a chain:

```go
chain := runner.Chain().
    Then("script1.sh").
    Then("script2.sh")

// Execute the chain
chain.Execute()

// Access results
exitCode := chain.ExitCode()
output := chain.Output()
err := chain.Error()

fmt.Printf("Last script exit code: %d\n", exitCode)
fmt.Printf("Last script output: %s\n", output)
if err != nil {
    fmt.Printf("Last script error: %v\n", err)
}
```

## Supported Script Types

By default, the following script types are supported:
- Shell scripts (.sh) - executed with bash
- Python scripts (.py) - executed with python

On Windows, Git Bash is used for executing shell scripts.