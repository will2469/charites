package astro_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/parser/astro"
)

func TestAstro_LineOffsetPreservation(t *testing.T) {
	// 10 baris frontmatter, elemen template berada tepat di baris 11
	src := `---
// Line 2
// Line 3
// Line 4
// Line 5
// Line 6
// Line 7
// Line 8
// Line 9
---
<div class="container mx-auto">
  <span class="text-sm">Hello</span>
</div>`

	root, err := astro.Parse([]byte(src))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if root.Tag != "div" {
		t.Fatalf("expected root tag 'div', got %q", root.Tag)
	}
	if root.Span.Line != 11 {
		t.Errorf("expected div to be on line 11 after 10 lines of frontmatter, got line %d", root.Span.Line)
	}

	// Periksa child span
	if len(root.Children) == 0 {
		t.Fatalf("expected children under root")
	}

	var spanNode *ir.Node
	for _, child := range root.Children {
		if child.Tag == "span" {
			spanNode = child
			break
		}
	}
	if spanNode == nil {
		t.Fatalf("span child not found")
	}
	if spanNode.Span.Line != 12 {
		t.Errorf("expected span to be on line 12, got line %d", spanNode.Span.Line)
	}
}

func TestAstro_CustomComponents(t *testing.T) {
	src := `
<Card class="p-6 shadow-md">
  <slot name="header" />
  <p class="text-base">Body content</p>
  <Footer client:load />
</Card>
`
	root, err := astro.Parse([]byte(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root.Tag != "Card" {
		t.Fatalf("expected root 'Card', got %q", root.Tag)
	}
	if !root.HasClass("p-6") || !root.HasClass("shadow-md") {
		t.Errorf("expected classes 'p-6 shadow-md', got %v", root.Classes)
	}

	var slotNode, footerNode *ir.Node
	for _, child := range root.Children {
		if child.Tag == "slot" {
			slotNode = child
		}
		if child.Tag == "Footer" {
			footerNode = child
		}
	}

	if slotNode == nil {
		t.Errorf("slot node not found")
	} else if slotNode.Attributes["name"] != `"header"` {
		t.Errorf("expected slot name 'header', got %q", slotNode.Attributes["name"])
	}

	if footerNode == nil {
		t.Errorf("Footer node not found")
	} else if _, ok := footerNode.Attributes["client:load"]; !ok {
		t.Errorf("expected client:load attribute on Footer")
	}
}

func TestAstro_VoidElements(t *testing.T) {
	src := `
<div class="card">
  <img src="avatar.png" alt="Avatar" class="w-10 h-10" />
  <input type="text" disabled />
  <p class="desc">Text</p>
</div>
`
	root, err := astro.Parse([]byte(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// img dan input adalah void elements. Mereka tidak boleh memerangkap <p> sebagai child!
	if root.Tag != "div" {
		t.Fatalf("expected root div, got %q", root.Tag)
	}

	var childTags []string
	for _, child := range root.Children {
		if child.Type == ir.NodeElement {
			childTags = append(childTags, child.Tag)
		}
	}

	expected := []string{"img", "input", "p"}
	if len(childTags) != len(expected) {
		t.Fatalf("expected children tags %v, got %v", expected, childTags)
	}
	for i, tag := range expected {
		if childTags[i] != tag {
			t.Errorf("at index %d: expected %s, got %s", i, tag, childTags[i])
		}
	}
}

func TestAstro_Recovery(t *testing.T) {
	src := `
<div class="wrapper">
  <broken <button class="btn btn-primary">Submit</button>
  </unmatchedTag>
  <span class="badge">Valid</span>
</div>
`
	root, err := astro.Parse([]byte(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root.Tag != "div" {
		t.Fatalf("expected root div, got %q", root.Tag)
	}

	var elementTags []string
	for _, child := range root.Children {
		if child.Type == ir.NodeElement {
			elementTags = append(elementTags, child.Tag)
		}
	}

	expected := []string{"button", "span"}
	if len(elementTags) != len(expected) {
		t.Fatalf("expected tags %v, got %v", expected, elementTags)
	}
	for i, tag := range expected {
		if elementTags[i] != tag {
			t.Errorf("at index %d: expected %s, got %s", i, tag, elementTags[i])
		}
	}
}

func TestAstro_TemplateLiterals(t *testing.T) {
	src := "<div class={`p-4 ${isActive ? 'bg-primary' : 'bg-muted'} rounded-lg`}></div>"

	root, err := astro.Parse([]byte(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !root.HasDynamicClasses {
		t.Errorf("expected HasDynamicClasses to be true")
	}
	if !root.HasClass("p-4") || !root.HasClass("rounded-lg") {
		t.Errorf("expected static classes 'p-4' and 'rounded-lg', got %v", root.Classes)
	}
	if root.HasClass("bg-primary") || root.HasClass("bg-muted") {
		t.Errorf("opaque dynamic region should not leak into Classes, got %v", root.Classes)
	}
}

func TestAstro_CommentsAndFragments(t *testing.T) {
	src := `
<>
  <!-- HTML Comment -->
  {/* JSX Comment */}
  <section class="hero">Content</section>
</>
`
	root, err := astro.Parse([]byte(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root.Type != ir.NodeFragment {
		t.Fatalf("expected root fragment, got %v", root.Type)
	}

	hasComment := false
	hasSection := false
	for _, child := range root.Children {
		if child.Type == ir.NodeComment {
			hasComment = true
		}
		if child.Type == ir.NodeElement && child.Tag == "section" {
			hasSection = true
		}
	}

	if !hasComment {
		t.Errorf("expected comment child")
	}
	if !hasSection {
		t.Errorf("expected section child")
	}
}

func TestAstro_EdgeCases(t *testing.T) {
	// 1. Empty source
	root, err := astro.Parse([]byte(""))
	if err != nil || root == nil {
		t.Fatalf("expected valid root for empty source, got err: %v", err)
	}

	// 2. Unclosed frontmatter
	srcUnclosedFM := "---\nconst x = 1;\n<div class='test'></div>"
	root, err = astro.Parse([]byte(srcUnclosedFM))
	if err != nil || root == nil {
		t.Fatalf("unexpected error on unclosed frontmatter: %v", err)
	}

	// 3. DOCTYPE and single quoted attributes
	srcDoctype := "<!DOCTYPE html><html lang='en'><head><meta charset='utf-8'></head><body>hello</body></html>"
	root, err = astro.Parse([]byte(srcDoctype))
	if err != nil || root == nil {
		t.Fatalf("unexpected error on doctype: %v", err)
	}

	// 4. Unclosed HTML comment and unclosed JSX comment
	srcUnclosedComment := "<!-- unclosed comment\n<div class=\"ok\"></div>"
	root, err = astro.Parse([]byte(srcUnclosedComment))
	if err != nil || root == nil {
		t.Fatalf("unexpected error on unclosed comment: %v", err)
	}

	srcUnclosedJSXComment := "{/* unclosed jsx comment\n<span class=\"ok\"></span>"
	root, err = astro.Parse([]byte(srcUnclosedJSXComment))
	if err != nil || root == nil {
		t.Fatalf("unexpected error on unclosed jsx comment: %v", err)
	}

	// 5. Unclosed quotes in attributes
	srcUnclosedQuotes := `<div class="unclosed>`
	root, err = astro.Parse([]byte(srcUnclosedQuotes))
	if err != nil || root == nil {
		t.Fatalf("unexpected error on unclosed quotes: %v", err)
	}
}
