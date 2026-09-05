package ir_test

import (
	"fmt"
	"testing"

	"github.com/will2469/charites/internal/ir"
)

func TestBuilder_NormalHierarchy(t *testing.T) {
	b := ir.NewBuilder()

	// <div><span>text</span><button /></div>
	div := b.OpenElement("div", ir.Span{Line: 1, Column: 1}, "container", []string{"container"}, nil, false)
	span := b.OpenElement("span", ir.Span{Line: 2, Column: 3}, "", nil, nil, false)
	text := b.AddText("hello", ir.Span{Line: 2, Column: 9})
	b.CloseElement("span")
	btn := b.AddSelfClosingElement("button", ir.Span{Line: 3, Column: 3}, "btn", []string{"btn"}, nil, false)
	b.CloseElement("div")

	root := b.Root()
	if root != div {
		t.Fatalf("expected root to be div, got %v", root)
	}
	if b.StackDepth() != 0 {
		t.Errorf("expected empty stack after closing root, got depth %d", b.StackDepth())
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children under root, got %d", len(root.Children))
	}
	if root.Children[0] != span || root.Children[1] != btn {
		t.Errorf("children hierarchy mismatch")
	}
	if span.Parent != div || btn.Parent != div {
		t.Errorf("parent pointers not set correctly")
	}
	if len(span.Children) != 1 || span.Children[0] != text {
		t.Errorf("text child not attached to span")
	}
	if text.Parent != span {
		t.Errorf("text parent pointer mismatch")
	}
}

func TestBuilder_UnmatchedClosingTag(t *testing.T) {
	b := ir.NewBuilder()

	b.OpenElement("div", ir.Span{Line: 1, Column: 1}, "", nil, nil, false)
	b.OpenElement("p", ir.Span{Line: 2, Column: 1}, "", nil, nil, false)

	initialDepth := b.StackDepth()
	if initialDepth != 2 {
		t.Fatalf("expected depth 2, got %d", initialDepth)
	}

	// Tag penutup yang tidak ada di stack: </article>
	b.CloseElement("article")

	// Stack harus TIDAK BERUBAH
	if b.StackDepth() != initialDepth {
		t.Errorf("expected stack depth to remain %d, got %d", initialDepth, b.StackDepth())
	}

	// Tag penutup valid
	b.CloseElement("p")
	if b.StackDepth() != 1 {
		t.Errorf("expected depth 1 after closing p, got %d", b.StackDepth())
	}

	b.CloseElement("div")
	if b.StackDepth() != 0 {
		t.Errorf("expected depth 0 after closing div, got %d", b.StackDepth())
	}

	// Menutup tag saat stack kosong tidak boleh panic
	b.CloseElement("anything")
	if b.StackDepth() != 0 {
		t.Errorf("expected depth 0, got %d", b.StackDepth())
	}
}

func TestBuilder_StackUnwinding(t *testing.T) {
	b := ir.NewBuilder()

	// <div><section><span><button>
	div := b.OpenElement("div", ir.Span{Line: 1, Column: 1}, "", nil, nil, false)
	sec := b.OpenElement("section", ir.Span{Line: 2, Column: 1}, "", nil, nil, false)
	span := b.OpenElement("span", ir.Span{Line: 3, Column: 1}, "", nil, nil, false)
	btn := b.OpenElement("button", ir.Span{Line: 4, Column: 1}, "", nil, nil, false)

	if b.StackDepth() != 4 {
		t.Fatalf("expected depth 4, got %d", b.StackDepth())
	}

	// </div> dipanggil langsung, melompati </button>, </span>, </section>
	b.CloseElement("div")

	if b.StackDepth() != 0 {
		t.Errorf("expected stack depth 0 after unwinding to div, got %d", b.StackDepth())
	}

	// Elemen perantara yang tidak ditutup tetap merupakan anak sah dari parent masing-masing
	if btn.Parent != span {
		t.Errorf("btn.Parent should be span")
	}
	if span.Parent != sec {
		t.Errorf("span.Parent should be sec")
	}
	if sec.Parent != div {
		t.Errorf("sec.Parent should be div")
	}
	if len(span.Children) != 1 || span.Children[0] != btn {
		t.Errorf("span should contain btn as child")
	}
}

func TestBuilder_NestingGuard256(t *testing.T) {
	b := ir.NewBuilder()

	// Push 300 elemen bersarang
	var nodes []*ir.Node
	for i := 1; i <= 300; i++ {
		tag := fmt.Sprintf("node-%d", i)
		n := b.OpenElement(tag, ir.Span{Line: i, Column: 1}, "", nil, nil, false)
		nodes = append(nodes, n)

		depth := b.StackDepth()
		if depth > ir.MaxStackDepth {
			t.Fatalf("at node %d: stack depth %d exceeded MaxStackDepth %d", i, depth, ir.MaxStackDepth)
		}
	}

	// Stack depth harus persis 256
	if b.StackDepth() != ir.MaxStackDepth {
		t.Fatalf("expected stack depth to be capped at %d, got %d", ir.MaxStackDepth, b.StackDepth())
	}

	// Node ke-256 (indeks 255) adalah parent di puncak stack
	parent256 := nodes[255]

	// Node ke-257 hingga ke-300 harus disematkan sebagai anak flat dari node-256
	expectedChildrenCount := 300 - 256
	if len(parent256.Children) != expectedChildrenCount {
		t.Fatalf("expected %d flat children under depth-256 parent, got %d", expectedChildrenCount, len(parent256.Children))
	}

	for i := 256; i < 300; i++ {
		child := nodes[i]
		if child.Parent != parent256 {
			t.Errorf("node %d parent expected to be parent256, got %v", i+1, child.Parent)
		}
	}

	// Tutup node flat (tidak ada di stack, harus discarded secara hening)
	b.CloseElement("node-290")
	if b.StackDepth() != ir.MaxStackDepth {
		t.Errorf("closing flat child should not change stack depth")
	}

	// Tutup node-256: harus sukses mem-pop level 256
	b.CloseElement("node-256")
	if b.StackDepth() != 255 {
		t.Errorf("expected stack depth 255, got %d", b.StackDepth())
	}
}

func TestBuilder_Fragments(t *testing.T) {
	b := ir.NewBuilder()

	frag := b.OpenFragment(ir.Span{Line: 1, Column: 1})
	b.AddSelfClosingElement("span", ir.Span{Line: 2, Column: 3}, "", nil, nil, false)
	b.CloseFragment()

	if b.StackDepth() != 0 {
		t.Errorf("expected depth 0, got %d", b.StackDepth())
	}
	root := b.Root()
	if root != frag || root.Type != ir.NodeFragment {
		t.Errorf("expected root to be fragment")
	}
}

func TestBuilder_MultipleRoots(t *testing.T) {
	b := ir.NewBuilder()

	// Dua elemen level atas terpisah
	d1 := b.OpenElement("div", ir.Span{Line: 1, Column: 1}, "", nil, nil, false)
	b.CloseElement("div")

	d2 := b.OpenElement("span", ir.Span{Line: 3, Column: 1}, "", nil, nil, false)
	b.CloseElement("span")

	root := b.Root()
	if root.Type != ir.NodeFragment {
		t.Fatalf("expected root to become NodeFragment for multiple roots, got %v", root.Type)
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children under root fragment, got %d", len(root.Children))
	}
	if root.Children[0] != d1 || root.Children[1] != d2 {
		t.Errorf("children under root fragment mismatch")
	}
}

func TestBuilder_CommentsAndReset(t *testing.T) {
	b := ir.NewBuilder()

	div := b.OpenElement("div", ir.Span{Line: 1, Column: 1}, "", nil, nil, false)
	c := b.AddComment("/* comment */", ir.Span{Line: 1, Column: 6})
	b.CloseElement("div")

	if len(div.Children) != 1 || div.Children[0] != c {
		t.Errorf("comment child not attached")
	}
	if c.Type != ir.NodeComment {
		t.Errorf("expected NodeComment type")
	}

	b.Reset()
	if b.StackDepth() != 0 {
		t.Errorf("expected stack depth 0 after reset")
	}
	emptyRoot := b.Root()
	if emptyRoot.Type != ir.NodeFragment {
		t.Errorf("expected empty root fragment after reset")
	}
}
