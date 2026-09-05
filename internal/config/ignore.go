package config

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

var builtinExclusions = []string{
	".git", "node_modules", "dist", ".astro", ".next", ".turbo", "build", "coverage",
}

// IgnoreMatcher mengevaluasi aturan pengabaian berkas sekuensial (.charitesignore)
// dengan penegakan kekebalan direktori builtin (Hard Exclusion Invariant).
type IgnoreMatcher struct {
	patterns []ignorePattern
}

type ignorePattern struct {
	raw      string
	pattern  string
	negation bool
	dirOnly  bool
	anchored bool
}

// BuiltinExclusions mengembalikan salinan daftar direktori builtin yang tidak dapat dinegasi.
func BuiltinExclusions() []string {
	out := make([]string, len(builtinExclusions))
	copy(out, builtinExclusions)
	return out
}

// LoadIgnore membaca berkas .charitesignore jika ada.
// Jika berkas tidak ditemukan, mengembalikan IgnoreMatcher kosong dengan proteksi builtin tetap aktif.
func LoadIgnore(path string) (*IgnoreMatcher, error) {
	if path == "" {
		path = ".charitesignore"
	}

	data, err := os.ReadFile(filepath.Clean(path)) //nolint:gosec // controlled ignore file path
	if err != nil {
		if os.IsNotExist(err) {
			return NewIgnoreMatcher(nil), nil
		}
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return NewIgnoreMatcher(lines), scanner.Err()
}

// NewIgnoreMatcher membuat matcher dari daftar baris aturan pengabaian.
func NewIgnoreMatcher(lines []string) *IgnoreMatcher {
	m := &IgnoreMatcher{
		patterns: make([]ignorePattern, 0, len(lines)),
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		negation := false
		if strings.HasPrefix(trimmed, "!") {
			negation = true
			trimmed = strings.TrimSpace(trimmed[1:])
		}

		dirOnly := false
		if strings.HasSuffix(trimmed, "/") {
			dirOnly = true
			trimmed = strings.TrimSuffix(trimmed, "/")
		}

		anchored := false
		if strings.HasPrefix(trimmed, "/") {
			anchored = true
			trimmed = strings.TrimPrefix(trimmed, "/")
		} else if strings.Contains(trimmed, "/") {
			anchored = true
		}

		// Normalisasi separator ke POSIX '/'
		normPattern := filepath.ToSlash(trimmed)

		m.patterns = append(m.patterns, ignorePattern{
			raw:      line,
			pattern:  normPattern,
			negation: negation,
			dirOnly:  dirOnly,
			anchored: anchored,
		})
	}

	return m
}

// AddPatterns menambahkan pola tambahan secara dinamis (misalnya dari charites.yaml ignore:).
func (m *IgnoreMatcher) AddPatterns(patterns []string) {
	temp := NewIgnoreMatcher(patterns)
	m.patterns = append(m.patterns, temp.patterns...)
}

// HasBuiltinAncestor memeriksa apakah ada segmen path yang cocok dengan builtin exclusion.
// Digunakan untuk proteksi eksplisit target berkas langsung (Explicit Target Safety).
func (m *IgnoreMatcher) HasBuiltinAncestor(path string) bool {
	clean := filepath.Clean(path)
	parts := strings.Split(filepath.ToSlash(clean), "/")
	for _, part := range parts {
		for _, b := range builtinExclusions {
			if part == b {
				return true
			}
		}
	}
	return false
}

// ShouldIgnoreDir mengevaluasi apakah suatu direktori harus dipangkas (early directory pruning).
func (m *IgnoreMatcher) ShouldIgnoreDir(dirName, relPath string) bool {
	// 1. Invarian Hard Exclusion (Builtin tidak bisa dinegasi)
	normRel := filepath.ToSlash(filepath.Clean(relPath))
	parts := strings.Split(normRel, "/")
	for _, part := range parts {
		for _, b := range builtinExclusions {
			if part == b {
				return true
			}
		}
	}
	for _, b := range builtinExclusions {
		if dirName == b {
			return true
		}
	}

	// 2. Evaluasi Sekuensial .charitesignore (Last matching pattern wins)
	ignored := false
	for _, p := range m.patterns {
		if matchPattern(p, normRel, dirName, true) {
			ignored = !p.negation
		}
	}
	return ignored
}

// ShouldIgnoreFile mengevaluasi apakah berkas harus diabaikan dari pemindaian.
func (m *IgnoreMatcher) ShouldIgnoreFile(relPath string) bool {
	// 1. Invarian Hard Exclusion
	if m.HasBuiltinAncestor(relPath) {
		return true
	}

	normRel := filepath.ToSlash(filepath.Clean(relPath))
	baseName := filepath.Base(normRel)

	// 2. Evaluasi Sekuensial .charitesignore (Last matching pattern wins)
	ignored := false
	for _, p := range m.patterns {
		if p.dirOnly {
			continue // Pola dir-only tidak mencocokkan berkas
		}
		if matchPattern(p, normRel, baseName, false) {
			ignored = !p.negation
		}
	}
	return ignored
}

// matchPattern mencocokkan pola ignorePattern terhadap target path.
func matchPattern(p ignorePattern, targetPath, baseName string, isDir bool) bool {
	if p.pattern == "" {
		return false
	}

	// Jika pola menuntut direktori dan target bukan direktori
	if p.dirOnly && !isDir {
		return false
	}

	// Jika pola tidak di-anchor (tanpa slash), cocokkan terhadap baseName
	if !p.anchored {
		if matched, _ := filepath.Match(p.pattern, baseName); matched {
			return true
		}
		if p.pattern == baseName {
			return true
		}
	}

	// Cocokkan terhadap full relative path
	return matchPathGlob(p.pattern, targetPath)
}

// matchPathGlob mendukung wildcard '*' dan recursive '**'.
func matchPathGlob(pattern, path string) bool {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	return matchParts(patternParts, pathParts)
}

func matchParts(patternParts, pathParts []string) bool {
	if len(patternParts) == 0 {
		return len(pathParts) == 0
	}

	if patternParts[0] == "**" {
		// '**' di posisi akhir mencocokkan sisa path apa pun
		if len(patternParts) == 1 {
			return true
		}
		// Coba cocokkan sisa pattern terhadap sub-slice path
		for i := 0; i <= len(pathParts); i++ {
			if matchParts(patternParts[1:], pathParts[i:]) {
				return true
			}
		}
		return false
	}

	if len(pathParts) == 0 {
		return false
	}

	// Cocokkan segmen tunggal
	matched, _ := filepath.Match(patternParts[0], pathParts[0])
	if !matched && patternParts[0] != pathParts[0] {
		return false
	}

	return matchParts(patternParts[1:], pathParts[1:])
}
