package reporter

import (
	"io"
	"os"
)

// ColorMode menentukan kebijakan pewarnaan ANSI escape codes pada presenter teks.
type ColorMode int

const (
	// ColorAuto mengaktifkan ANSI escape codes hanya jika output diarahkan ke terminal interaktif (TTY).
	ColorAuto ColorMode = iota
	// ColorNever menonaktifkan seluruh kode pewarnaan ANSI (plain text).
	ColorNever
	// ColorAlways memaksa pewarnaan ANSI terlepas dari status TTY (berguna untuk testing/fixture).
	ColorAlways
)

// ResolveColorMode menentukan mode warna berdasarkan flag CLI, variabel lingkungan NO_COLOR, dan tipe writer.
func ResolveColorMode(noColorFlag bool, w io.Writer) ColorMode {
	if noColorFlag {
		return ColorNever
	}
	if os.Getenv("NO_COLOR") != "" {
		return ColorNever
	}
	if !IsTerminal(w) {
		return ColorNever
	}
	return ColorAuto
}

// IsTerminal memeriksa apakah writer terhubung ke perangkat karakter interaktif (TTY).
// Menggunakan Go Standard Library tanpa dependensi eksternal.
func IsTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
