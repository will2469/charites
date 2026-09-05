package fuzz

import (
	"os"
	"testing"

	"github.com/will2469/charites/internal/parser/astro"
)

func FuzzAstroParser(f *testing.F) {
	// Seed corpus dengan fixture dan variasi malformed
	if sample, err := os.ReadFile("../fixtures/sample.astro"); err == nil {
		f.Add(sample)
	}

	seeds := [][]byte{
		[]byte(""),
		[]byte("---\nconst a = 1;\n---\n<div></div>"),
		[]byte("<broken <button class='btn'>Click</button>"),
		[]byte("<!-- unclosed comment"),
		[]byte("<div class={`p-4 ${dynamic} text-sm`}></div>"),
		[]byte("<slot name='header' /><input type='text' disabled>"),
		[]byte("<><span>Text</span></>"),
		[]byte("<!DOCTYPE html><html><body><h1>Title</h1></body></html>"),
		[]byte("---"),
		[]byte("<<<<<<>>>>>>"),
		[]byte("<tag attr=val unquoted="),
		[]byte("{/* unclosed comment"),
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		// Zero-panic invariant: parser dilarang keras memicu runtime panic
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic triggered in astro.Parse with input %q: %v", string(data), r)
			}
		}()

		root, _ := astro.Parse(data)
		if root != nil {
			// Pastikan tree traversal tidak panic
			count := 0
			for range root.Walk() {
				count++
				if count > 10000 {
					break
				}
			}
		}
	})
}
