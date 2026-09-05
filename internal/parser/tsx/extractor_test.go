package tsx_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/parser/tsx"
)

func TestTSX_AttributesAndSelfClosing(t *testing.T) {
	src := `
export function Button() {
  return (
    <button className="btn btn-primary" id="submit-btn" disabled />
  );
}
`
	root, err := tsx.Extract([]byte(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root.Tag != "button" {
		t.Fatalf("expected button tag, got %q", root.Tag)
	}
	if !root.HasClass("btn") || !root.HasClass("btn-primary") {
		t.Errorf("expected classes 'btn btn-primary', got %v", root.Classes)
	}
	if root.Attributes["id"] != `"submit-btn"` {
		t.Errorf("expected id attribute 'submit-btn', got %q", root.Attributes["id"])
	}
	if _, ok := root.Attributes["disabled"]; !ok {
		t.Errorf("expected disabled attribute")
	}
}

func TestTSX_Fragments(t *testing.T) {
	src := `
export default function Layout() {
  return (
    <>
      <header className="header">Header</header>
      <main className="content">Main</main>
    </>
  );
}
`
	root, err := tsx.Extract([]byte(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root.Type != ir.NodeFragment {
		t.Fatalf("expected root fragment, got %v", root.Type)
	}

	var elementTags []string
	for _, child := range root.Children {
		if child.Type == ir.NodeElement {
			elementTags = append(elementTags, child.Tag)
		}
	}

	expected := []string{"header", "main"}
	if len(elementTags) != len(expected) {
		t.Fatalf("expected elements %v, got %v", expected, elementTags)
	}
	for i, tag := range expected {
		if elementTags[i] != tag {
			t.Errorf("at index %d: expected %s, got %s", i, tag, elementTags[i])
		}
	}
}

func TestTSX_StaticTemplateLiterals(t *testing.T) {
	src := `
export const Box = () => (
  <div className={` + "`flex flex-col p-4`}" + `>
    <span className="text-base font-medium">Title</span>
  </div>
);
`
	root, err := tsx.Extract([]byte(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root.Tag != "div" {
		t.Fatalf("expected root div, got %q", root.Tag)
	}
	if root.HasDynamicClasses {
		t.Errorf("pure static template literal should have HasDynamicClasses = false")
	}
	expected := []string{"flex", "flex-col", "p-4"}
	for _, exp := range expected {
		if !root.HasClass(exp) {
			t.Errorf("expected class %q, got %v", exp, root.Classes)
		}
	}
}

func TestTSX_DynamicTemplateLiterals(t *testing.T) {
	src := `
export function Alert({ isOpen, severity }) {
  return (
    <div className={` + "`p-4 ${isOpen ? \"opacity-100\" : \"opacity-0\"} rounded-lg text-sm`}" + `>
      <p>Message</p>
    </div>
  );
}
`
	root, err := tsx.Extract([]byte(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root.Tag != "div" {
		t.Fatalf("expected root div, got %q", root.Tag)
	}
	if !root.HasDynamicClasses {
		t.Errorf("expected HasDynamicClasses to be true")
	}
	if !root.HasClass("p-4") || !root.HasClass("rounded-lg") || !root.HasClass("text-sm") {
		t.Errorf("expected static classes 'p-4 rounded-lg text-sm', got %v", root.Classes)
	}
	if root.HasClass("opacity-100") || root.HasClass("opacity-0") {
		t.Errorf("dynamic opaque region should not leak into Classes, got %v", root.Classes)
	}
}

func TestTSX_LessThanDisambiguation(t *testing.T) {
	src := `
import React from 'react';

// Disambiguation test:
// 1. Comments with '<broken' or '<div'
// 2. String with '<'
// 3. Comparison count < 10
// 4. Comparison in ternary {count < 5 ? <span>Low</span> : <span>High</span>}

/* multi-line comment
   <broken <tag
*/

export function Dashboard({ count }: { count: number }) {
  const isSmall = count < 10;
  const title = "Count < 10";

  return (
    <div className="dashboard-container">
      {count < 5 ? (
        <span className="text-red-500">Low</span>
      ) : (
        <span className="text-green-500">High</span>
      )}
      <input type="text" disabled />
    </div>
  );
}
`
	root, err := tsx.Extract([]byte(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root.Tag != "div" {
		t.Fatalf("expected root tag 'div', got %q", root.Tag)
	}

	var elementTags []string
	for _, child := range root.Children {
		if child.Type == ir.NodeElement {
			elementTags = append(elementTags, child.Tag)
		}
	}

	expected := []string{"span", "span", "input"}
	if len(elementTags) != len(expected) {
		t.Fatalf("expected element tags %v, got %v", expected, elementTags)
	}
	for i, tag := range expected {
		if elementTags[i] != tag {
			t.Errorf("at index %d: expected %s, got %s", i, tag, elementTags[i])
		}
	}
}

func TestTSX_SpreadAttributesAndEdgeCases(t *testing.T) {
	// 1. Spread attributes: {...props}
	srcSpread := `
export const CustomInput = (props) => (
  <input {...props} className="input-field" />
);
`
	root, err := tsx.Extract([]byte(srcSpread))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root.Tag != "input" {
		t.Fatalf("expected input, got %q", root.Tag)
	}
	if !root.HasClass("input-field") {
		t.Errorf("expected class 'input-field', got %v", root.Classes)
	}

	// 2. Empty source
	emptyRoot, err := tsx.Extract([]byte(""))
	if err != nil || emptyRoot == nil {
		t.Fatalf("expected valid root for empty source, got %v", err)
	}

	// 3. Generics and type assertions
	srcGenerics := `
type ValueMap = Map<string, number>;
const identity = <T>(x: T): T => x;
export const Comp = () => <div className="generic-test" />;
`
	genRoot, err := tsx.Extract([]byte(srcGenerics))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if genRoot.Tag != "div" {
		t.Fatalf("expected div, got %q", genRoot.Tag)
	}
}

func TestTSX_RecoveryAndBranchCoverage(t *testing.T) {
	// 1. Broken tag recovery (<broken <button>)
	srcBroken := `
export const App = () => (
  <main>
    <broken <button className='btn-single'>Click</button>
    </
    <img src="pic.jpg" alt='Photo'>
  </main>
);
`
	root, err := tsx.Extract([]byte(srcBroken))
	if err != nil || root == nil {
		t.Fatalf("unexpected error on broken tags: %v", err)
	}
	if root.Tag != "main" {
		t.Errorf("expected root main, got %q", root.Tag)
	}

	// 2. Unclosed string literals and comments
	srcUnclosed := "/* unclosed comment\nconst s = 'unclosed;\nconst t = `unclosed template ${val;"
	root, err = tsx.Extract([]byte(srcUnclosed))
	if err != nil || root == nil {
		t.Fatalf("unexpected error on unclosed structures: %v", err)
	}

	// 3. Unclosed tag and fragment closing
	srcFrag := `<><span>Text</span></>`
	root, err = tsx.Extract([]byte(srcFrag))
	if err != nil || root == nil {
		t.Fatalf("unexpected error on fragment: %v", err)
	}

	// 4. Unquoted attributes and self-closing tags
	srcUnquoted := `// comment
const dummy = "string";
export const Comp = () => (
  <CustomComponent disabled class=simple-btn />
);`
	root, err = tsx.Extract([]byte(srcUnquoted))
	if err != nil || root == nil {
		t.Fatalf("unexpected error on unquoted attribute: %v", err)
	}
	if root.Tag != "CustomComponent" {
		t.Errorf("expected CustomComponent, got %q", root.Tag)
	}
}
