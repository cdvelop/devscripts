package devscripts

import (
	"os"
	"path/filepath"
	"strings"
)

// ScriptInfo represents information about a shell script
type ScriptInfo struct {
	Name        string
	Description string
	Usage       string
}

// ScriptParser handles parsing of shell scripts
type ScriptParser struct {
	scriptsDir string
}

// NewScriptParser creates a new ScriptParser
func NewScriptParser(scriptsDir string) *ScriptParser {
	return &ScriptParser{scriptsDir: scriptsDir}
}

// GetScriptNames obtiene los nombres de los scripts .sh en el directorio
func (sp *ScriptParser) GetScriptNames() ([]string, error) {
	files, err := os.ReadDir(sp.scriptsDir)
	if err != nil {
		return nil, err
	}

	var scripts []string
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".sh" {
			scripts = append(scripts, f.Name())
		}
	}

	return scripts, nil
}

// ParseScripts obtiene las descripciones de los scripts
func (sp *ScriptParser) ParseScripts() ([]ScriptInfo, error) {
	scripts, err := sp.GetScriptNames()
	if err != nil {
		return nil, err
	}

	var scriptsStruct []ScriptInfo

	for _, script := range scripts {
		content, err := os.ReadFile(filepath.Join(sp.scriptsDir, script))
		if err != nil {
			return nil, err
		}

		if len(content) == 0 {
			scriptsStruct = append(scriptsStruct, ScriptInfo{
				Name:        script,
				Description: "Empty script file",
				Usage:       "",
			})
			continue
		}

		lines := strings.Split(string(content), "\n")

		if len(lines) >= 3 {
			// Extract from line 2 (index 1)
			descLine := strings.TrimSpace(lines[1])
			descLine = strings.TrimPrefix(descLine, "# Description:")
			descLine = strings.TrimSpace(descLine)

			// Extract from line 3 (index 2)
			usageLine := strings.TrimSpace(lines[2])
			usageLine = strings.TrimPrefix(usageLine, "# Usage:")
			usageLine = strings.TrimSpace(usageLine)

			// Combine description and usage
			if descLine != "" {
				if usageLine != "" {
					scriptsStruct = append(scriptsStruct, ScriptInfo{
						Name:        script,
						Description: descLine,
						Usage:       usageLine,
					})
				} else {
					scriptsStruct = append(scriptsStruct, ScriptInfo{
						Name:        script,
						Description: descLine,
						Usage:       "",
					})
				}
			} else {
				scriptsStruct = append(scriptsStruct, ScriptInfo{
					Name:        script,
					Description: sp.generateAutoDescription(script, strings.Join(lines, "\n")),
					Usage:       "",
				})
			}
		} else {
			scriptsStruct = append(scriptsStruct, ScriptInfo{
				Name:        script,
				Description: sp.generateAutoDescription(script, strings.Join(lines, "\n")),
				Usage:       "",
			})
		}
	}

	return scriptsStruct, nil
}

func (sp *ScriptParser) generateAutoDescription(name, content string) string {
	keywords := map[string]string{
		"git":    "Git operations",
		"repo":   "Repository management",
		"setup":  "System setup/config",
		"update": "Dependency updates",
		"go":     "Go language utilities",
	}

	var desc []string
	for k, v := range keywords {
		if strings.Contains(strings.ToLower(name), k) {
			desc = append(desc, v)
		}
	}

	if len(desc) == 0 {
		return "Shell script utility"
	}
	return strings.Join(desc, ", ")
}
