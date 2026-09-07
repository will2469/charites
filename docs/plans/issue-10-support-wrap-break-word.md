# Implementation Plan: Support Tailwind v4 `wrap-break-word` in Overflow Rules (`responsive.mobile-text-overflow`)

Resolves **[GitHub Issue #10](https://github.com/will2469/charites/issues/10)** (`[enhancement] feat(tailwind-v4): support 'wrap-break-word' in overflow rules and suggest modern replacement for 'break-words'`).

---

## 1. Executive Summary & Problem Formulation

In Tailwind CSS v4 and modern CSS specifications (CSS Text Module Level 4 / CSS Overflow Level 3):
- The CSS property `overflow-wrap: break-word` is canonically represented in Tailwind v4 as **`wrap-break-word`** (harmonized with `wrap-normal`, `wrap-nowrap`, `wrap-balance`, `wrap-pretty`), replacing the legacy Tailwind v3 class **`break-words`**.
- In Charites, `responsive.mobile-text-overflow` checks for `break-words` and `break-all` on inline `<code>` and `whitespace-nowrap` text containers.
- Previously, `internal/rules/responsive/util.go` only recognized `"break-words"` and `"break-all"`. Developers using modern Tailwind v4 `wrap-break-word` triggered false-positive warnings.

---

## 2. Implemented Changes

### Component 1: Engine Utility in `internal/rules/responsive/util.go`
- Added `"wrap-break-word"` to `hasTextOverflowProtection()` alongside `"break-words"`, `"break-all"`, `"truncate"`, `"overflow-hidden"`, `"overflow-x-auto"`.
- Added `"wrap-break-word"` to `hasCodeWrapOrScroll()` alongside `"break-all"`, `"break-words"`, `"whitespace-normal"`, `"whitespace-pre-wrap"`.

### Component 2: Rule Documentation & Educational Hints in `internal/rules/responsive/mobile_text_overflow.go`
- Updated `Doc()` 8-pillars documentation (CoreInvariant, Grounding, GoodExamples) to reference `wrap-break-word` (Tailwind v4 standard).
- Updated diagnostic message & educational hint:
  - `Message`: `"Inline <code> element lacks word breaking ('break-all', 'wrap-break-word', 'break-words') or horizontal scroll container wrapper. Long snippets will blow out mobile container boundaries."`
  - `Hint`: `"Add 'break-all' or modern 'wrap-break-word' (Tailwind CSS v4) to the <code> element, or wrap it inside an 'overflow-x-auto' container."`

### Component 3: 1-SSOT Tri-Corpus Correctness Fixtures
- **Negative (`negative/`):**
  - `modern_wrap_break_word.tsx`: `<code className="wrap-break-word font-mono text-xs">...</code>` and `<div className="whitespace-nowrap wrap-break-word text-sm">...</div>` (asserts 0 violations for modern Tailwind v4 class).
- **Adversarial (`adversarial/`):**
  - `modern_wrap_break_word.astro`: Astro component utilizing `wrap-break-word` on dynamic code and text tags (0 violations).

### Component 4: Documentation & Wikis
- Regenerated `wiki/responsive.mobile-text-overflow.md` via `make wiki`.

---

## 3. Verification Summary
- Tri-Corpus Harness: `TestRule_ResponsiveMobileTextOverflow_TriCorpus` (100% pass across Positive, Negative, Adversarial).
- Responsive Rules: all unit tests in `internal/rules/responsive/...` pass.
- Linter: `make lint` clean (0 errors).
- Adoption Matrix & Correctness Gate: 100% pass.
- Race Detector: `go test -race ./...` (0 race conditions).
