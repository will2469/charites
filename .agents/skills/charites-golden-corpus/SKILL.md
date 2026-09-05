---
name: charites-golden-corpus
description: "MANDATORY 1-SSOT ADVERSARIAL HARNESS & RESILIENCE CORPUS FOR CHARITES: Comprehensive testing framework and Single Source of Truth (SSOT) for Charites AST rules (Astro, TSX, CSS). Enforces the 17-pattern matrix: Positive (P1-P5 obvious, indirect, helper, nested, alias violations), Negative (N1-N5 valid token, explicit ignore, third-party lib, semantic HTML, config exception), and Adversarial (A1-A7 template string interp, conditional/ternary, spread props, dynamic classes, shadowed vars, nested closures/HOC, obfuscated classes). Tracks rule-by-rule adoption progress via TestGoldenCorpus_AdoptionMatrix ('golden corpus', 'test corpus', 'adversarial corpus', 'P1-P5', 'N1-N5', 'A1-A7', 'ssot test', 'adoption matrix')."
compatibility: "Go 1.26+, Leaf IR Parser, bash"
metadata:
  version: "1.0.0"
  author: "Charites Team / Will"
  license: "MIT"
  citations:
    - "IETF RFC 2119: Key words for use in RFCs to Indicate Requirement Levels"
    - "W3C CSS Color Module Level 4 (OKLCH, Display-P3)"
    - "W3C Web Content Accessibility Guidelines (WCAG 2.2)"
    - "Core Web Vitals & INP/LCP Optimization Standards"
---

# Charites Golden Adversarial Corpus Standard (`charites-golden-corpus`)

> **Core Thesis:** Linter atau static analyzer yang hanya menguji *"apakah kode yang jelas-jelas salah terdeteksi?"* memiliki ilusi keamanan (*false sense of security*). Integritas compiler-linter sejati menuntut pembuktian atas dua pertanyaan kritis:
> 1. **"Apakah kode yang benar dan aman terjamin lolos tanpa false positive?"** (Mencegah erosi kepercayaan developer akibat peringatan palsu).
> 2. **"Apakah kode yang sengaja disamarkan atau dievaluasi secara dinamis dapat lolos (*evade*)?"** (Mencegah false negative pada pola template literal, conditional classes, computed properties, dan closure components).
>
> **The 1-SSOT Mandate:** Seluruh aturan Charites berpusat secara tunggal pada arsitektur **1-SSOT Tri-Corpus** di bawah `tests/correctness/<category>.<slug>/` sebagai **Satu-Satunya Sumber Kebenaran (Single Source of Truth)**. CLI scanner, MCP server, dan automated test suite semuanya mengevaluasi sumber kanonikal yang sama.

---

## 0. The 4-Layer Quality Pyramid

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THE 4-LAYER COMPILER-GRADE QUALITY GATES                 │
├─────────────────────────────────────────────────────────────────────────────┤
│  Layer 1: Golden Corpus Adoption (Structural Presence Gate)                 │
│  ├─► Memverifikasi keberadaan struktur 1-SSOT (P/N/A fixtures & runner)     │
│  └─► Gate: TestGoldenCorpus_AdoptionMatrix (Status: 100% File Presence)     │
│                                                                             │
│  Layer 2: Golden Corpus Correctness (Semantic Execution Gate)               │
│  ├─► Positive Gate: P1-P5 memicu diagnostik persis pada line & span target  │
│  ├─► Negative Gate: N1-N5 nol false-positive pada token valid & legitimate  │
│  ├─► Adversarial Gate: A1-A7 tertangkap (template strings, ternaries, etc.) │
│  └─► Format Parity Gate: Inline ANSI == JSON Stream == Markdown Table       │
│                                                                             │
│  Layer 3: Mutation & Evasion Testing (Resilience Gate)                      │
│  ├─► Expression Inversion: boolean ternaries, nullish coalescing, fallback  │
│  ├─► Scope Evasion: Shadowed CSS vars, imported token alias, computed keys  │
│  ├─► Lexical Spoofing: Comments containing false directives, string quotes  │
│  ├─► AST Evasion: Deep slots, fragments, HOC wrappers, nested JSX elements  │
│  ├─► Fail-Closed AST Invariant: Template unparseable dilaporkan, no bypass  │
│  └─► Target: Mutation Kill Rate = 100% (Zero Surviving Mutants)             │
│                                                                             │
│  Layer 4: Cross-Rule Regression & Interaction Matrix (Isolation Gate)       │
│  ├─► Multi-rule concurrent execution (seluruh rule aktif simultan)          │
│  ├─► Shared Leaf IR isolation (Zero cache poisoning di AST cache)           │
│  ├─► Directive scoping isolation (/* charites:ignore */ scoped per rule)    │
│  └─► Whole-Program Golden Corpus (`tests/e2e/fixtures/`)                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 1. The 17-Pattern Tri-Corpus Taxonomy

Setiap rule di Charites wajib mengimplementasikan matriks kanonikal 17-pola:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                   CHARITES 17-PATTERN CORPUS TAXONOMY                       │
├─────────────────────────────────────────────────────────────────────────────┤
│  POSITIVE (Violations - MUST Trigger 100% with exact span)                  │
│  • P1: Direct Violation      - Hardcoded inline style, raw hex/rgb color    │
│  • P2: Indirect Violation    - Class string concatenation, template literal │
│  • P3: Helper Violation      - Utility function / clsx / cn wrapper call    │
│  • P4: Nested Violation      - Nested slot, deep JSX child, CSS inside tag  │
│  • P5: Alias Violation       - Aliased import of style object or constant   │
├─────────────────────────────────────────────────────────────────────────────┤
│  NEGATIVE (Safe - MUST SURVIVE with 0 False Positives)                      │
│  • N1: Valid Design Token    - Standard var(--color-*) or tailwind token    │
│  • N2: Explicit Exemption    - charites:ignore directive on line or block   │
│  • N3: Third-Party Lib       - Component inside node_modules or vendor dir  │
│  • N4: Semantic Static HTML  - Native HTML attributes without styling issue │
│  • N5: Config Exception      - Path or rule disabled in charites.yaml       │
├─────────────────────────────────────────────────────────────────────────────┤
│  ADVERSARIAL (Edge Cases & Subtle Evaders)                                  │
│  • A1: Template String Interp- `bg-[${dynamic}]`, computed class names      │
│  • A2: Conditional/Ternary   - isActive ? '#ff0000' : 'var(--primary)'      │
│  • A3: Spread Props          - {...restProps} overriding class / style      │
│  • A4: Dynamic Classes       - object syntax { 'bg-red-500': hasError }     │
│  • A5: Shadowed Variables    - local const color hides global token name    │
│  • A6: Nested Closures / HOC - higher-order component wrapping UI element   │
│  • A7: Obfuscated Classes    - arbitrary values with comments / escapes     │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Directory Hierarchy & 1-SSOT Architecture

Semua fixture pengujian berpusat secara eksklusif di direktori `tests/`:

```text
tests/
├── corpus_status_test.go                # Automated Adoption Matrix Checker (100% ADOPTED Gate)
└── correctness/                         # Rule Test Suites
    ├── theme/
    │   └── hardcode-opacity-color/
    │       ├── positive/
    │       │   ├── direct.astro         # P1: Obvious inline violation
    │       │   └── indirect.tsx         # P2: Template string / variant concatenation
    │       ├── negative/
    │       │   ├── tokens.astro         # N1: Standard design tokens
    │       │   └── clean.tsx            # N2: Clean standard classes
    │       ├── adversarial/
    │       │   ├── slash_layout.astro   # A1: Dimension / layout fractions
    │       │   ├── line_height.tsx      # A2: Typography line-height modifiers
    │       │   ├── arbitrary_color.tsx  # A3: Arbitrary hex colors
    │       │   └── unmapped_opacity.astro # A4: Unmapped opacities & raw palettes
    │       └── rule_test.go             # Automated runner for this rule
    └── ... (other categories and rules)
```

---

## 3. Anotasi Harapan Uji (*Want Directives*)

Pada berkas `positive/`, baris pelanggaran ditandai dengan komentar khusus sesuai bahasa:

```astro
---
// Astro frontmatter
---
<!-- want "theme.hardcode-color: avoid hardcoded hex color #ff0000" -->
<div style="background-color: #ff0000;">Bad Color</div>
```

```tsx
// TSX file
// want "theme.hardcode-color"
<div className="bg-[#123456]">Hex Tailwind</div>
```

```css
/* CSS file */
.card {
  /* want "theme.hardcode-color" */
  color: rgb(255, 0, 0);
}
```

---

## 4. Checklist Verifikasi Adopsi Golden Corpus

- [ ] **1. Structural Presence:** Fixture `positive/`, `negative/`, dan `adversarial/` lengkap ada di disk.
- [ ] **2. Semantic Execution:** Seluruh kasus P1-P5 memicu diagnostik dengan line number tepat.
- [ ] **3. Zero False-Positives:** N1-N5 menghasilkan tepat 0 diagnostik.
- [ ] **4. Adversarial Resilience:** A1-A7 menangkap upaya evasion tanpa crash atau bypass.
- [ ] **5. Invariant Adoption Test:** `go test ./tests -run TestGoldenCorpus_AdoptionMatrix` lulus 100%.
