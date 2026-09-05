package token_test

import (
	"testing"

	"github.com/will2469/charites/internal/token"
)

func TestCascade_SourceOrder(t *testing.T) {
	// Sesuai CSS Cascade: jika layer, spesifisitas, dan kondisi identik,
	// deklarasi terakhir (source order lebih besar) harus menang.
	input := []byte(`:root {
  --color: #111111;
  --color: #222222;
  --result: var(--color);
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("ParseCSS failed: %v", err)
	}

	resToks := ctx.ByName("--result")
	if len(resToks) != 1 {
		t.Fatalf("expected 1 --result token")
	}

	val, ok, err := ctx.Resolve(resToks[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed to resolve --result: %v", err)
	}
	if val != "#222222" {
		t.Errorf("expected source order to pick later declaration '#222222', got %q", val)
	}
}

func TestCascade_UnlayeredBeatsLayered(t *testing.T) {
	// Sesuai CSS Cascade 5: Deklarasi unlayered (tanpa @layer)
	// selalu mengalahkan deklarasi di dalam @layer terlepas dari source order.
	input := []byte(`:root {
  --brand: #unlayered;
}

@layer utilities {
  :root {
    --brand: #layered;
  }
}

.consumer {
  --active: var(--brand);
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("ParseCSS failed: %v", err)
	}

	consToks := ctx.ByName("--active")
	if len(consToks) != 1 {
		t.Fatalf("expected 1 --active token")
	}

	val, ok, err := ctx.Resolve(consToks[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed to resolve --active: %v", err)
	}
	if val != "#unlayered" {
		t.Errorf("expected unlayered token '#unlayered' to win over layered, got %q", val)
	}
}

func TestCascade_LayerOrdering(t *testing.T) {
	// Sesuai CSS Cascade: @layer reset, framework;
	// mendefinisikan bahwa framework memiliki prioritas lebih tinggi daripada reset,
	// terlepas dari urutan penulisan bloknya di bawah.
	input := []byte(`@layer reset, framework;

@layer framework {
  :root {
    --brand-color: #framework-wins;
  }
}

@layer reset {
  :root {
    --brand-color: #reset-loses;
  }
}

.box {
  --final: var(--brand-color);
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("ParseCSS failed: %v", err)
	}

	finalToks := ctx.ByName("--final")
	if len(finalToks) != 1 {
		t.Fatalf("expected 1 --final token")
	}

	val, ok, err := ctx.Resolve(finalToks[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed to resolve --final: %v", err)
	}
	if val != "#framework-wins" {
		t.Errorf("expected @layer framework '#framework-wins' to beat reset, got %q", val)
	}
}

func TestCascade_ConditionalAtRules(t *testing.T) {
	// Sesuai CSS Cascade: @media (prefers-color-scheme: dark)
	// hanya menang saat target context berada di kondisi dark mode.
	input := []byte(`:root {
  --surface: #ffffff;
  --resolved-light: var(--surface);
}

@media (prefers-color-scheme: dark) {
  :root {
    --surface: #000000;
    --resolved-dark: var(--surface);
  }
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("ParseCSS failed: %v", err)
	}

	// 1. Evaluasi pada scope root biasa (light mode)
	lightToks := ctx.ByName("--resolved-light")
	if len(lightToks) != 1 {
		t.Fatalf("expected 1 --resolved-light token")
	}
	lightVal, ok, err := ctx.Resolve(lightToks[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed resolving light token: %v", err)
	}
	if lightVal != "#ffffff" {
		t.Errorf("expected light token to resolve to '#ffffff', got %q", lightVal)
	}

	// 2. Evaluasi pada scope dark mode
	darkToks := ctx.ByName("--resolved-dark")
	if len(darkToks) != 1 {
		t.Fatalf("expected 1 --resolved-dark token")
	}
	darkVal, ok, err := ctx.Resolve(darkToks[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed resolving dark token: %v", err)
	}
	if darkVal != "#000000" {
		t.Errorf("expected dark token to resolve to '#000000', got %q", darkVal)
	}
}

func TestCascade_ElementDirectOverInheritedRoot(t *testing.T) {
	// Sesuai aturan CSS property inheritance: deklarasi langsung pada elemen (.card)
	// mendahului nilai yang diwarisi dari root (:root) terlepas dari urutan penulisan.
	input := []byte(`:root {
  --pad: 10px;
}

.card {
  --pad: 24px;
  --content-pad: var(--pad);
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("ParseCSS failed: %v", err)
	}

	toks := ctx.ByName("--content-pad")
	if len(toks) != 1 {
		t.Fatalf("expected 1 --content-pad token")
	}

	val, ok, err := ctx.Resolve(toks[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed to resolve: %v", err)
	}
	if val != "24px" {
		t.Errorf("expected direct element declaration '24px', got %q", val)
	}
}

func TestCascade_UnlayeredBeatsConditionalLayered(t *testing.T) {
	// W3C CSS Cascade 5: Deklarasi unlayered selalu mengalahkan deklarasi layered,
	// bahkan jika deklarasi layered dibungkus oleh kondisi @media yang cocok!
	// (Membuktikan tidak ada artificial ConditionScore yang membajak cascade layer).
	input := []byte(`@media (prefers-color-scheme: dark) {
  @layer theme {
    :root {
      --surface: #layered-dark;
    }
  }
}

:root {
  --surface: #unlayered-wins;
}

@media (prefers-color-scheme: dark) {
  .consumer {
    --active: var(--surface);
  }
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("ParseCSS failed: %v", err)
	}

	toks := ctx.ByName("--active")
	if len(toks) != 1 {
		t.Fatalf("expected 1 --active token")
	}

	val, ok, err := ctx.Resolve(toks[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed to resolve: %v", err)
	}
	if val != "#unlayered-wins" {
		t.Errorf("expected unlayered token '#unlayered-wins' to beat layered even under media query, got %q", val)
	}
}

func TestCascade_IsVsWhereSpecificity(t *testing.T) {
	// W3C CSS Selectors 4:
	// :is(#hero) memiliki spesifisitas ID (1, 0, 0).
	// :where(#hero) memiliki spesifisitas NOL (0, 0, 0).
	// Oleh karena itu, aturan dengan :is() harus mengalahkan :where() meskipun :where() ditulis belakangan!
	input := []byte(`:is(#hero) {
  --accent: #is-wins;
}

:where(#hero) {
  --accent: #where-loses;
}

#hero {
  --resolved: var(--accent);
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("ParseCSS failed: %v", err)
	}

	toks := ctx.ByName("--resolved")
	if len(toks) != 1 {
		t.Fatalf("expected 1 --resolved token")
	}

	val, ok, err := ctx.Resolve(toks[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed to resolve: %v", err)
	}
	if val != "#is-wins" {
		t.Errorf("expected :is(#hero) '#is-wins' (1,0,0) to beat :where(#hero) (0,0,0), got %q", val)
	}
}
