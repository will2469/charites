package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/will2469/charites/internal/rules"
	"github.com/will2469/charites/internal/wiki"
)

// RunWiki mengeksekusi pembuatan ensiklopedia dokumentasi Markdown berbasis rules.Registry.
func RunWiki(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("wiki", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			_, _ = fmt.Fprintf(stdout, "Usage: charites wiki [output_dir]\n\nGenerate ensiklopedia dokumentasi Markdown untuk seluruh rule terdaftar.\n")
			return ExitClean
		}
		_, _ = fmt.Fprintf(stderr, "charites: error: %v\n", err)
		return ExitOperational
	}

	outDir := "wiki"
	if fs.NArg() > 0 {
		outDir = strings.TrimSpace(fs.Arg(0))
	}
	if outDir == "" {
		outDir = "wiki"
	}

	gen := wiki.NewGenerator(rules.DefaultRegistry())
	if err := gen.Generate(outDir); err != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to generate wiki: %v\n", err)
		return ExitOperational
	}

	_, _ = fmt.Fprintf(stdout, "charites: wiki documentation successfully generated at %s\n", outDir)
	return ExitClean
}
