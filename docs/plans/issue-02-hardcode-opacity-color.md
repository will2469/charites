# Fix Rule Gap: theme.hardcode-opacity-color Classification Boundary & Arbitrary Slash Opacity Detection

Close static analysis gap and silent-pass loophole in `theme.hardcode-opacity-color` ([Issue #2](https://github.com/will2469/charites/issues/2)) by establishing an explicit classification-boundary engine pipeline, adding `shadow-` utility class support to the shared prefix registry, rejecting arbitrary uncalibrated slash opacities that lack semantic tokens in `global.css`, providing standardized static diagnostic hints (Option A), and enforcing zero-noise orthogonal rule delegation across the 1-SSOT Tri-Corpus.

---

## User Review Required

> [!IMPORTANT]
> **Golden Snapshot Evolution (`opacity_violations`):**
> In [`tests/fixtures/projects/opacity_violations/src/components/Card.tsx`](../../tests/fixtures/projects/opacity_violations/src/components/Card.tsx#L5), `<div className="bg-primary/10 text-secondary/50 ...">` contains `text-secondary/50`.
> Previously, `secondary/50` silently passed because `FindOpacityReplacement()` returned `false`. Under the new classification contract, `text-secondary/50` is recognized as a semantic candidate lacking a calibrated replacement and is correctly flagged.
> The golden snapshot in [`tests/golden/projects/opacity_violations.golden.json`](../../tests/golden/projects/opacity_violations.golden.json) and [`opacity_violations.golden.txt`](../../tests/golden/projects/opacity_violations.golden.txt) will be updated from 2 to 3 diagnostics via `go test ./tests -run TestPipeline_GoldenSnapshots -update`.

> [!NOTE]
> **SemVer Classification:**
> Per `charites-versioning` (SemVer 2.0.0 Clause 5), this is a backward-compatible false-negative bug fix (PATCH level). Public CLI syntax, MCP server schemas, and `internal/ir` public structures remain 100% backward compatible.

---

## The 4 Hard Acceptance Invariants

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THE 4 NON-NEGOTIABLE ACCEPTANCE INVARIANTS               │
├─────────────────────────────────────────────────────────────────────────────┤
│ 1. Predicate Decoupling Invariant:                                          │
│    FindOpacityReplacement() MUST NEVER be used as the predicate for         │
│    determining whether a class is in-scope. FindOpacityReplacement() == false│
│    MUST NOT imply clean/pass.                                               │
│                                                                             │
│ 2. Bracket-Safe Parsing Invariant:                                          │
│    SplitAlphaModifier() MUST be used as the mandatory parsing mechanism.    │
│    The rule MUST NOT parse '/' inside bracketed arbitrary expressions as an │
│    alpha modifier delimiter (e.g. [color:rgb(0/0/0)] is not alpha).         │
│    Bracket in color-base (bg-[#123456]/10) delegates to theme.hardcode-color,│
│    while bracket in alpha (bg-primary/[0.1]) stays in this rule.            │
│                                                                             │
│ 3. Shadow Disambiguation Invariant:                                         │
│    shadow-* elevation dimensions (shadow-sm, shadow-md, shadow-lg,          │
│    shadow-inner, shadow-none, and even shadow-sm/20) MUST NOT be classified │
│    as colors. Only shadow-* with semantic color bases are analyzed.         │
│                                                                             │
│ 4. Orthogonal Bait Ownership Invariant:                                     │
│    Every fixture in adversarial/ MUST have an explicit owning orthogonal    │
│    rule. Any detection by theme.hardcode-opacity-color indicates an illegal │
│    boundary leakage.                                                        │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Core Classification Contract

Sequential pipeline executed for every class in `node.Classes`:

```text
node.Classes
    │
1. Strip Variants (StripVariantsOnlyBase)
    │  Handles arbitrary variants like [&>svg]:text-primary/20 (0 allocs)
    │
2. Split Alpha Modifier (SplitAlphaModifier)
    ├─ !hasAlpha ────────────────────────► [PASS / CLEAN]
    │  Extracts baseNoAlpha and alpha (e.g. "bg-primary" & "10", "bg-primary" & "[0.1]")
    │
3. Resolve Utility Color Prefix (SplitColorPrefix)
    ├─ !ok ──────────────────────────────► [PASS / CLEAN]
    │  Extracts prefix ("bg-", "text-", "border-", "shadow-") & colorBase
    │
4. Reject Non-Color Utility Keywords
    ├─ prefix == "text-" && IsTailwindFontSize(colorBase) ───────► [IGNORE / CLEAN] (e.g. text-sm/6, text-3xl/9)
    │  (IsTailwindFontSize ONLY classifies known typography steps; MUST NOT classify color names)
    ├─ strings.HasPrefix(prefix, "border") && IsNonColorBorderKeyword(colorBase) ► [IGNORE / CLEAN] (e.g. border-2, border-solid)
    ├─ prefix == "shadow-" && IsShadowSizeKeyword(colorBase) ────► [IGNORE / CLEAN] (e.g. shadow-sm, shadow-sm/20)
    │
5. Delegate Primitive Palette Colors
    ├─ IsTailwindPrimitiveColor(colorBase) ──────────────────────► [DELEGATE: theme.primitive-in-component] (e.g. bg-red-500/10)
    │
6. Delegate Monochrome Colors
    ├─ IsMonochromeColor(colorBase) ─────────────────────────────► [DELEGATE: theme.hardcode-monochrome] (e.g. bg-black/10, text-white/20)
    │
7. Delegate Arbitrary / Raw Colors
    ├─ strings.HasPrefix(colorBase, "[") || IsHexColor(colorBase)► [DELEGATE: theme.hardcode-color] (e.g. bg-[#123456]/10)
    │
8. Treat Remaining Color Bases as Semantic Candidates
    │
    ▼
Query conv.FindOpacityReplacement(colorBase, alpha, tCtx)
    ├─ [FOUND (ok == true)]:
    │    VIOLATION + Specific Semantic Token Hint:
    │    Hint: "Use semantic token \"" + cands[0].Name + "\"."
    │
    └─ [NOT FOUND (ok == false)]:
         VIOLATION + Static Generic Calibrated Token Hint (Option A):
         Hint: "Use an existing semantic token or declare a calibrated semantic token in global.css (e.g. --<base>-<state>) instead of using arbitrary slash opacity modifiers."
         (Static constant; MUST NOT interpolate colorBase, alpha, or state)
```

---

## Adversarial Bait Ownership Matrix

| Bait Pattern | Fixture File | Owning Rule / Reason | Expected Finding |
| :--- | :--- | :--- | :---: |
| `bg-red-500/10`, `text-blue-600/20` | `adversarial/primitive_palette_opacity.astro` | `theme.primitive-in-component` | 0 |
| `bg-black/10`, `text-white/20` | `adversarial/monochrome_opacity.astro` | `theme.hardcode-monochrome` | 0 |
| `bg-[#123456]/10`, `text-[#ff0000]/20` | `adversarial/arbitrary_color_opacity.astro` | `theme.hardcode-color` | 0 |
| `text-sm/6`, `text-3xl/9`, `text-xs/relaxed` | `adversarial/typography_fraction.astro` | Tailwind typography line-height | 0 |
| `border-2`, `border-solid`, `border-0` | `adversarial/border_non_color.astro` | Tailwind border geometry/style | 0 |
| `shadow-sm`, `shadow-md`, `shadow-sm/20` | `adversarial/shadow_elevation.astro` | Tailwind box-shadow elevation | 0 |
| `w-1/2`, `h-1/3`, `aspect-16/9`, `grid-cols-2/3` | `adversarial/layout_fraction.astro` | Tailwind layout dimensions | 0 |

---

## Phased Implementation Steps

### Phase 1: Shared Classifiers & Parsing Helpers (`internal/rules/theme/util.go`)

#### [MODIFY] `internal/rules/theme/util.go`
- **Extend `OrderedColorPrefixes`:**
  Add `"shadow-"` to `OrderedColorPrefixes`.
- **Add `IsShadowSizeKeyword(s string) bool`:**
  ```go
  func IsShadowSizeKeyword(s string) bool {
      switch s {
      case "sm", "md", "lg", "xl", "2xl", "inner", "none":
          return true
      default:
          return false
      }
  }
  ```
- **Add `IsTailwindFontSize(s string) bool`:**
  ```go
  func IsTailwindFontSize(s string) bool {
      switch s {
      case "xs", "sm", "base", "lg", "xl", "2xl", "3xl", "4xl", "5xl", "6xl", "7xl", "8xl", "9xl":
          return true
      default:
          return false
      }
  }
  ```
  Contract: Strictly checks known font size steps; never matches arbitrary semantic colors (`primary`, `warning`, `muted`).

#### [MODIFY] `internal/rules/theme/important_override.go`
- In `isColorUtility(prefix, remainder)`:
  Add:
  ```go
  if prefix == "shadow-" && IsShadowSizeKeyword(remainder) {
      return false
  }
  ```
  Ensures `!shadow-md` is never falsely treated as a color utility.

---

### Phase 2: Rewrite Rule Decision Boundary (`internal/rules/theme/hardcode_opacity_color.go`)

#### [MODIFY] `internal/rules/theme/hardcode_opacity_color.go`
- Define static hint constant:
  ```go
  const uncalibratedOpacityHint = "Use an existing semantic token or declare a calibrated semantic token in global.css (e.g. --<base>-<state>) instead of using arbitrary slash opacity modifiers."
  ```
- Delete private local `stripVariants()` and call `StripVariantsOnlyBase(class)`.
- Use `SplitAlphaModifier(base)` to extract `baseNoAlpha`, `alpha`, and `hasAlpha`.
- Use `SplitColorPrefix(baseNoAlpha)` to extract `prefix` and `colorBase`.
- Apply non-color exclusions (`IsTailwindFontSize`, `IsNonColorBorderKeyword`, `IsShadowSizeKeyword`).
- Apply orthogonal rule delegations (`IsTailwindPrimitiveColor`, `IsMonochromeColor`, bracket/hex check).
- Query `conv.FindOpacityReplacement(colorBase, alpha, tCtx)`.
- Emit specific hint if replacement found; emit static `uncalibratedOpacityHint` if replacement not found.
- Update `Description()` and `Doc()` Bad Examples (`shadow-primary/20`, `hover:border-primary/50`, `bg-muted/20`, `border-warning/40`, `text-warning/90`).

---

### Phase 3: Unit Tests & Regression Suites (`internal/rules/theme/`)

#### [MODIFY] `internal/rules/theme/hardcode_opacity_color_test.go`
- In `TestHardcodeOpacityColorRule_TableDrivenBoundary`:
  - Rename test case from `OutOfScope_unmapped_opacities` to `InScope_unmapped_opacities` and assert 4 violations with static generic hint for `bg-primary/30`, `bg-primary/50`, `bg-primary/100`, and `bg-primary/[0.1]`.
  - Add test case for `shadow-primary/20` (assert mapped hint `Use semantic token "primary-light".`).
  - Add test case for `shadow-sm/20` (assert 0 violations - elevation size keyword).
  - Add test case for `[&>svg]:text-primary/20` (assert variant bracket stripping).
  - Add test case for Issue #2 reproducible snippet:
    `<div className="shadow-primary/20 hover:border-primary/50 bg-muted/20 border-warning/40 text-warning/90">`
    asserting 5 diagnostics with exact expected hints.
  - Maintain clean node benchmark: assert `0 B/op, 0 allocs/op`.
- Add dedicated top-level test:
  ```go
  func TestHardcodeOpacityColorRule_UnmappedSemanticOpacities(t *testing.T)
  ```
  Explicitly asserting that mapped and unmapped semantic opacities produce two distinct hint types.

#### [NEW] `internal/rules/theme/important_override_regression_test.go`
- Add targeted regression test verifying:
  - `!shadow-sm`, `!shadow-md`, `!shadow-lg`, `!shadow-inner`, `!shadow-none` produce 0 diagnostics.
  - `!shadow-primary` produces 1 diagnostic.

---

### Phase 4: 1-SSOT Tri-Corpus Golden Fixtures (`tests/correctness/theme/hardcode-opacity-color/`)

#### [MODIFY] `tests/correctness/theme/hardcode-opacity-color/positive/basic.astro`
- Retain mapped semantic token violations (`bg-primary/10`, `border-destructive/20`, etc.).

#### [NEW] `tests/correctness/theme/hardcode-opacity-color/positive/arbitrary_opacity.astro`
- Add unmapped arbitrary slash opacities (`shadow-primary/20`, `hover:border-primary/50`, `border-warning/40`, `text-warning/90`, `bg-primary/30`, `bg-primary/[0.1]`).

#### [MODIFY] & [NEW] Structured Adversarial Bait Fixtures:
- [MODIFY] `adversarial/unmapped_opacity.astro` $\rightarrow$ split into:
  - [NEW] `adversarial/primitive_palette_opacity.astro` (`bg-red-500/10`, `text-blue-600/20`, `border-gray-300/5`)
  - [NEW] `adversarial/monochrome_opacity.astro` (`bg-black/10`, `text-white/20`)
  - [MODIFY] `adversarial/arbitrary_color_opacity.astro` (`bg-[#123456]/10`, `text-[#ff0000]/20`, `border-[#abcdef]/5`)
  - [MODIFY] `adversarial/typography_fraction.astro` (`text-sm/6`, `text-3xl/9`, `text-xs/relaxed`, `text-base/7`)
  - [NEW] `adversarial/shadow_elevation.astro` (`shadow-sm`, `shadow-md`, `shadow-lg`, `shadow-sm/20`)
  - [MODIFY] `adversarial/slash_layout.astro` (`w-1/2`, `h-1/3`, `aspect-16/9`, `grid-cols-2/3`)
- Delete old monolithic `adversarial/unmapped_opacity.astro`.

#### [MODIFY] `tests/correctness/theme/hardcode-opacity-color/rule_test.go`
- Assert positive violations for both mapped and arbitrary opacities.
- Assert strict 0 diagnostics across all 7 adversarial bait fixtures.

---

### Phase 5: Golden Snapshots Update

#### [MODIFY] `tests/golden/projects/opacity_violations.golden.json` & `tests/golden/projects/opacity_violations.golden.txt`
- Update snapshots via `go test ./tests -run TestPipeline_GoldenSnapshots -update`, recording 3 diagnostics including `text-secondary/50`.

---

### Phase 6: Automated 8-Pillars Wiki Generation

#### [MODIFY] `wiki/theme.hardcode-opacity-color.md`
- Recompile deterministic wiki via `make wiki`.

---

## Verification Plan

### Automated Tests
1. **Unit Tests & Zero-Alloc Benchmarks:**
   ```bash
   go test -v ./internal/rules/theme/...
   go test -bench=BenchmarkEvaluateHardcodeOpacityColor -benchmem ./internal/rules/theme/...
   ```
   *Expected:* All tests pass; `BenchmarkEvaluateHardcodeOpacityColor_Clean` reports `0 B/op, 0 allocs/op`.

2. **Important Override Regression:**
   ```bash
   go test -v -run TestImportantOverride_ShadowRegression ./internal/rules/theme/...
   ```
   *Expected:* Passes; zero false positives on shadow elevation sizes.

3. **Tri-Corpus Golden Correctness Matrix:**
   ```bash
   go test -v ./tests/correctness/theme/hardcode-opacity-color/...
   ```
   *Expected:* Positive violations pass with exact hints; Negative zero noise == 0; Adversarial bait immunity == 0.

4. **Golden Snapshot Integrity:**
   ```bash
   go test -v ./tests -run TestPipeline_GoldenSnapshots
   ```
   *Expected:* Golden snapshots match 100%.

5. **Full Repository Concurrency & Race Detector:**
   ```bash
   go test -race -count=1 ./...
   ```
   *Expected:* 0 race conditions, all packages pass.

6. **Code Hygiene & Linters:**
   ```bash
   golangci-lint run ./...
   make lint
   ```
   *Expected:* Zero lint warnings, clean formatting.

7. **Wiki SSOT Generator Verification:**
   ```bash
   make wiki
   git diff wiki/theme.hardcode-opacity-color.md
   ```
   *Expected:* Byte-for-byte deterministic wiki generation.
