package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// UsageString mengembalikan pesan bantuan penggunaan CLI kanonikal sesuai SPEC-05-CLI.
func UsageString() string {
	return `Usage: charites [command] [flags] [path]

Available Commands:
  scan        Pindai berkas frontend untuk audit kualitas dan token semantik (Default)
  check       Alias identik untuk 'scan'
  run         Alias identik untuk 'scan'
  version     Cetak versi kompilasi binary, commit git, dan Go runtime
  help        Bantuan penggunaan perintah
`
}

// Execute adalah titik masuk trampoline utama untuk CLI menggunakan stream standar sistem operasi.
func Execute(args []string) int {
	return ExecuteArgs(args, os.Stdout, os.Stderr)
}

// ExecuteArgs mengeksekusi CLI dengan injeksi writer stdout dan stderr terisolasi.
func ExecuteArgs(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return RunScan([]string{"."}, stdout, stderr)
	}

	first := args[0]

	switch first {
	case "scan", "check", "run":
		return RunScan(args[1:], stdout, stderr)
	case "version", "-v", "--version":
		_, _ = fmt.Fprint(stdout, VersionString())
		return ExitClean
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(stdout, UsageString())
		return ExitClean
	default:
		// Jika argumen diawali '-' (flag langsung) atau berupa path langsung
		if strings.HasPrefix(first, "-") || isPath(first) {
			return RunScan(args, stdout, stderr)
		}
		_, _ = fmt.Fprintf(stderr, "charites: error: unknown command \"%s\". Run 'charites --help' for usage.\n", first)
		return ExitOperational
	}
}

func isPath(arg string) bool {
	if arg == "." || arg == ".." {
		return true
	}
	if strings.HasPrefix(arg, "./") || strings.HasPrefix(arg, "../") || strings.HasPrefix(arg, "/") {
		return true
	}
	if strings.Contains(arg, "/") || strings.Contains(arg, "\\") {
		return true
	}
	if _, err := os.Stat(arg); err == nil {
		return true
	}
	ext := strings.ToLower(filepath.Ext(arg))
	switch ext {
	case ".astro", ".tsx", ".jsx", ".ts", ".js", ".json", ".yaml", ".yml":
		return true
	}
	return false
}
