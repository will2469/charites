package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/will2469/charites/internal/mcp"
	"github.com/will2469/charites/internal/rules"
)

// RunMCP menjalankan Model Context Protocol (MCP) server berbasis Stdio JSON-RPC 2.0.
func RunMCP(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("mcp", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	var workspace string
	fs.StringVar(&workspace, "workspace", ".", "Direktori akar ruang kerja (workspace root)")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			_, _ = fmt.Fprintf(stdout, "Usage: charites mcp [--workspace=<path>]\n\nMulai server Model Context Protocol (MCP) berbasis Stdio JSON-RPC 2.0.\n")
			return ExitClean
		}
		_, _ = fmt.Fprintf(stderr, "charites: error: %v\n", err)
		return ExitOperational
	}

	if stdin == nil {
		stdin = os.Stdin
	}

	srv := mcp.NewServer(workspace, stdin, stdout, stderr, rules.DefaultRegistry(), Version)
	if err := srv.Run(); err != nil {
		_, _ = fmt.Fprintf(stderr, "charites mcp: error: %v\n", err)
		return ExitOperational
	}

	return ExitClean
}
