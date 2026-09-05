package cli

import (
	"fmt"
	"io"
	"os"
)

// UsageString returns the canonical CLI help message conforming to SPEC-00.
func UsageString() string {
	return `Usage: charites <command> [options] [path]

Commands:
  scan       Scan frontend files for design system, a11y, and performance issues
             Aliases: check, run
  version    Print binary version

Options for 'scan':
  -f, --format string      Output format: inline (default ANSI) or json
  --ext string             Filter by extension: astro, tsx, jsx
  --category string        Filter by category: theme, a11y, perf, layout, seo
  --rule string            Filter by single rule ID: theme.hardcode-opacity-color
  --ignore string          Additional custom ignore pattern
`
}

// Execute is the main entrypoint trampoline for the CLI using standard os streams.
func Execute(args []string) int {
	return ExecuteArgs(args, os.Stdout, os.Stderr)
}

// ExecuteArgs executes the CLI with injected stdout and stderr writers,
// enforcing strict stream routing isolation and exit codes (SPEC-00).
func ExecuteArgs(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(stdout, UsageString())
		return 0
	}

	first := args[0]

	switch first {
	case "-v", "--version", "version":
		_, _ = fmt.Fprint(stdout, VersionString())
		return 0
	case "-h", "--help", "help":
		_, _ = fmt.Fprint(stdout, UsageString())
		return 0
	default:
		_, _ = fmt.Fprintf(stderr, "charites: unknown command or flag \"%s\"\nRun 'charites --help' for usage.\n", first)
		return 2
	}
}
