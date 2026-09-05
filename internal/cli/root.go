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
              Aliases: check, run
  version     Cetak versi kompilasi binary, commit git, dan Go runtime
  help        Bantuan penggunaan perintah

Flags:
  -f, --format string      Format output: inline (default ANSI) atau json (default "inline")
  -e, --ext string         Filter ekstensi yang dipindai (default "astro,tsx,jsx")
  -c, --category string    Filter kategori rule (theme, a11y, perf, dll.)
  -r, --rule string        Filter satu Charites Rule ID spesifik (<category>.<slug>)
      --config string      Path kustom berkas konfigurasi (default "charites.yaml")
      --ignore string      Pola glob ignore tambahan (dapat berulang atau koma)
      --no-color           Matikan pewarnaan ANSI di terminal
      --fail-on-warn       Kembalikan exit code 1 jika hanya terdapat warning
  -h, --help               Bantuan penggunaan perintah
  -v, --version            Cetak versi binary
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
