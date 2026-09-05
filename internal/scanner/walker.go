package scanner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/will2469/charites/internal/config"
)

// MaxScanFileSize membatasi ukuran berkas sumber frontend maksimal 10 Megabytes.
const MaxScanFileSize = 10 * 1024 * 1024

// DefaultExtensions mendefinisikan ekstensi berkas yang didukung secara default.
var DefaultExtensions = []string{".astro", ".tsx", ".jsx"}

// Walker melakukan traversal direktori berkecepatan tinggi dengan proteksi symlink,
// batas ukuran berkas, serta pengecekan keamanan direct-target safety.
type Walker struct {
	matcher *config.IgnoreMatcher
	extMap  map[string]bool
}

// NewWalker membuat instans Walker baru.
// Jika extensions kosong, menggunakan DefaultExtensions (.astro, .tsx, .jsx).
func NewWalker(matcher *config.IgnoreMatcher, extensions []string) *Walker {
	if matcher == nil {
		matcher = config.NewIgnoreMatcher(nil)
	}

	exts := extensions
	if len(exts) == 0 {
		exts = DefaultExtensions
	}

	extMap := make(map[string]bool, len(exts))
	for _, ext := range exts {
		if ext != "" {
			if ext[0] != '.' {
				ext = "." + ext
			}
			extMap[ext] = true
		}
	}

	return &Walker{
		matcher: matcher,
		extMap:  extMap,
	}
}

// Walk menelusuri target path dan mengantrekan berkas yang valid ke channel jobs.
// Walker mematuhi kontrak direct-target safety: jika target berada di dalam
// direktori terlarang (builtin hard exclusions), eksekusi langsung digagalkan
// dan 0 pekerjaan diantrekan ke channel jobs.
func (w *Walker) Walk(ctx context.Context, target string, jobs chan<- string) error {
	cleanTarget := filepath.Clean(target)

	// 1. Direct-Target Safety Check: Tolak target jika memiliki leluhur terlarang
	if w.matcher.HasBuiltinAncestor(cleanTarget) {
		return fmt.Errorf("scan target %q is within excluded directory (builtin hard exclusion)", cleanTarget)
	}

	fi, err := os.Stat(cleanTarget)
	if err != nil {
		return err
	}

	// 2. Penanganan Target Berkas Tunggal Secara Langsung (Single File Target)
	if !fi.IsDir() {
		return w.handleSingleFile(ctx, cleanTarget, fi, jobs)
	}

	// 3. Traversal Direktori Rekursif
	return filepath.WalkDir(cleanTarget, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Lanjutkan traversal jika file individual tidak dapat dibaca
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		return w.walkDirEntry(ctx, cleanTarget, path, d, jobs)
	})
}

func (w *Walker) handleSingleFile(ctx context.Context, cleanTarget string, fi os.FileInfo, jobs chan<- string) error {
	ext := filepath.Ext(cleanTarget)
	if w.extMap[ext] && fi.Size() <= MaxScanFileSize && !w.matcher.ShouldIgnoreFile(cleanTarget) {
		select {
		case jobs <- cleanTarget:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (w *Walker) walkDirEntry(ctx context.Context, cleanTarget, path string, d os.DirEntry, jobs chan<- string) error {
	rel, _ := filepath.Rel(cleanTarget, path)

	// Proteksi Symlink: Jangan ikuti direktori symlink dan lewati berkas symlink
	if d.Type()&os.ModeSymlink != 0 {
		if d.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	// Direktori: Evaluasi Early Pruning
	if d.IsDir() {
		if rel != "." && w.matcher.ShouldIgnoreDir(d.Name(), rel) {
			return filepath.SkipDir
		}
		return nil
	}

	// Batas Ekstensi Berkas
	if !w.extMap[filepath.Ext(path)] {
		return nil
	}

	// Batas Ukuran Berkas Maksimal 10 MB
	info, err := d.Info()
	if err != nil || info.Size() > MaxScanFileSize {
		return nil
	}

	// Filter Ignore Berkas
	if w.matcher.ShouldIgnoreFile(rel) {
		return nil
	}

	select {
	case jobs <- path:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
