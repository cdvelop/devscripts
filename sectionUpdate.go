package devscripts

import (
	"fmt"
	"os"

	"github.com/cdvelop/mdgo"
)

// SectionUpdate is a small CLI-style entry calls mdgo directly.
//
// Args order: sectionID (required), afterLine (optional), content (required),
// readmeFile (optional, default: `README.md`).
func SectionUpdate(args ...string) {
	if len(args) == 0 {
		fmt.Println("Error: sectionID required")
		os.Exit(1)
	}

	sectionID := args[0]
	afterLine := ""
	content := ""
	readmeFile := "README.md"

	if len(args) > 1 {
		afterLine = args[1]
	}
	if len(args) > 2 {
		content = args[2]
	}
	if len(args) > 3 {
		readmeFile = args[3]
	}

	m := mdgo.New(".", ".", func(name string, data []byte) error {
		return os.WriteFile(name, data, 0644)
	})
	m.InputPath(readmeFile, func(name string) ([]byte, error) {
		return os.ReadFile(name)
	})

	err := m.UpdateSection(sectionID, content, afterLine)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
