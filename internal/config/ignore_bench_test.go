package config_test

import (
	"fmt"
	"testing"

	"github.com/will2469/charites/internal/config"
)

func BenchmarkIgnore_Matcher(b *testing.B) {
	patterns := []string{
		"node_modules/**",
		"dist/**",
		"coverage/**",
		"tests/fixtures/**/*.tmp",
		"!tests/fixtures/keep.tmp",
		"legacy_vendor/**",
		"*.log",
		"*.lock",
	}
	matcher := config.NewIgnoreMatcher(patterns)

	// Buat 1.000 path sintetis
	paths := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		switch {
		case i%5 == 0:
			paths[i] = fmt.Sprintf("tests/fixtures/sub%d/file_%d.tmp", i%10, i)
		case i%7 == 0:
			paths[i] = fmt.Sprintf("legacy_vendor/pkg%d/index.tsx", i%10)
		default:
			paths[i] = fmt.Sprintf("src/components/feature%d/Button%d.tsx", i%20, i)
		}
	}

	b.ReportAllocs()

	for b.Loop() {
		for _, p := range paths {
			_ = matcher.ShouldIgnoreFile(p)
		}
	}
}
