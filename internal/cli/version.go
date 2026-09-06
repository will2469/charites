package cli

import "fmt"

// Canonical version metadata set via build ldflags or defaults.
var (
	Version   = "1.0.0-beta.1"
	Commit    = "none"
	BuildDate = "unknown"
)

// VersionString returns the canonical version string conforming to SPEC-00.
func VersionString() string {
	return fmt.Sprintf("charites version %s (go1.26.x)\n", Version)
}
