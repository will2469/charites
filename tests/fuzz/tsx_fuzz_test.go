package fuzz

import (
	"os"
	"testing"

	"github.com/will2469/charites/internal/parser/tsx"
)

func FuzzTSXParser(f *testing.F) {
	// Seed corpus dengan fixture dan variasi malformed
	if sample, err := os.ReadFile("../fixtures/sample.tsx"); err == nil {
		f.Add(sample)
	}

	seeds := [][]byte{
		[]byte(""),
		[]byte("export const Box = () => <div className='box'>Hello</div>;"),
		[]byte("<button {...props} className={`p-4 ${dyn} text-sm`}>Click</button>"),
		[]byte("const isSmall = count < 10; const type = Map<string, number>;"),
		[]byte("<broken <tag attr='unclosed>"),
		[]byte("/* unclosed comment\nconst a = `unclosed template"),
		[]byte("<><span>Fragment</span></>"),
		[]byte("export function App() { return (<div>{count < 5 ? <span>Low</span> : <span>High</span>}</div>); }"),
		[]byte("<T>(x: T) => x;"),
		[]byte("<CustomComponent disabled attr=val />"),
		[]byte("<<<<<///// >>>>>"),
		[]byte("return <input type=\"text\""),
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		// Zero-panic invariant: parser dilarang keras memicu runtime panic
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic triggered in tsx.Extract with input %q: %v", string(data), r)
			}
		}()

		root, _ := tsx.Extract(data)
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
