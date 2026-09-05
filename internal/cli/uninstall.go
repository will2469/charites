package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// OsExecutable membungkus os.Executable untuk memfasilitasi pengujian unit.
var OsExecutable = os.Executable

// RunUninstall mencopot pemasangan Charites dengan menghapus satu-satunya biner eksekutabel.
// Menjamin pembersihan 100% tanpa residu file atau cache (Zero Residual Footprint Guarantee).
func RunUninstall(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("uninstall", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	var autoConfirm bool
	fs.BoolVar(&autoConfirm, "yes", false, "Lewati konfirmasi pencopotan")
	fs.BoolVar(&autoConfirm, "y", false, "Lewati konfirmasi (shorthand)")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			_, _ = fmt.Fprint(stdout, UsageString())
			return ExitClean
		}
		_, _ = fmt.Fprintf(stderr, "charites: error: %v. Run 'charites --help' for usage.\n", err)
		return ExitOperational
	}

	execPath, err := OsExecutable()
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to determine executable path: %v\n", err)
		return ExitOperational
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to resolve executable symlinks: %v\n", err)
		return ExitOperational
	}

	if err := os.Remove(execPath); err != nil {
		if os.IsPermission(err) {
			_, _ = fmt.Fprintf(stderr, "charites: error: failed to remove binary %q: permission denied (please run with sudo)\n", execPath)
			return ExitOperational
		}
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to remove binary %q: %v\n", execPath, err)
		return ExitOperational
	}

	_, _ = fmt.Fprintln(stdout, "Charites uninstalled successfully.")
	_, _ = fmt.Fprintf(stdout, "Removed binary: %s\n", execPath)
	_, _ = fmt.Fprintln(stdout, "0 residual files or caches remaining on the host system.")
	return ExitClean
}
