package devscripts

import (
	"fmt"
	"os"

	"github.com/cdvelop/badges"
)

// Badges is the entry function used by the cmd wrapper and the bash shim.
// It follows the ExecuteWithArgs convention: args are the parameters passed from the shell.
func Badges(args ...string) {
	// Create handler for badges package. The package `badges` will handle
	// validation of an injected currentDir argument (from gocurrentdir.sh)
	// including .git checks, so keep this wrapper minimal.
	sectionArgs, err := badges.NewBadgeHandler(args...).BuildBadges()
	if err != nil {
		// Print the error message text so tests that look for specific
		// substrings (like "No badges specified") can see them.
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Forward the returned args directly to SectionUpdate
	// SectionUpdate expects: sectionID, afterLine, content, readmeFile
	SectionUpdate(sectionArgs...)
}
