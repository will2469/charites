# performance.tailwind-untracked-package-source

> **Rule ID:** `performance.tailwind-untracked-package-source`
> **Severity:** `ERROR`
> **Category:** `performance`
> **Target Standards:** Tailwind CSS v4 Configuration Architecture (@source Directive Specification), Monorepo Multi-Package Style Discovery Standards, Oxide Engine Workspace Scanning Invariants

---

## 1. Overview & Core Invariant

Mewajibkan pendaftaran direktif @source pada berkas CSS root Tailwind v4 ketika mengimpor paket workspace monorepo eksternal.

### Core Invariant:
> **"Tailwind CSS v4 root stylesheets importing external monorepo packages must declare '@source' path directives; without '@source', the Oxide scanner skips external package directories, silently dropping all utility styles from compiled builds."**

---
## 2. Technical Grounding & Engine Realities

In Tailwind CSS v4, the legacy `tailwind.config.js` `content` array is replaced by CSS-first `@source` directives in the main stylesheet.

By default, Tailwind v4 only scans files in the immediate project directory. If the project imports components from external workspace packages (e.g. `@repo/ui` or `../../packages/...`), those package directories are ignored by default.

Failing to add `@source "../../packages/ui";` causes all Tailwind utility classes used inside those shared packages to be completely absent from the final CSS bundle.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Missing Monorepo Component Styles** | HIGH | Shared monorepo UI components render completely unstyled in production because utility classes inside them were never scanned. |
| **Silent Build Failures** | HIGH | No build errors are thrown; stylesheets simply compile without the required utility declarations. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### CSS (Berkas CSS root tidak menyertakan direktif @source untuk paket monorepo):
```css
/* Pelanggaran: Mengimpor tailwindcss tanpa @source untuk monorepo packages */
@import "tailwindcss";
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### CSS (Mendaftarkan path paket eksternal via @source):
```css
/* Patuh: Menyertakan direktif @source untuk paket monorepo */
@import "tailwindcss";
@source "../../packages/ui";
```

---

## 6. Detection & Verification Pipeline (How The Rule Evaluates Code)
This rule evaluates source code through the standard AST inspection pipeline:

```mermaid
flowchart TD
    Node["AST Node (Astro / TSX element)"] --> Inspect["1. Inspect Element & Attributes"]
    Inspect --> Invariant{"2. Evaluate Rule Invariant"}
    Invariant -- "Compliant" --> Safe["Pass"]
    Invariant -- "Non-Compliant" --> IgnoreCheck{"3. Check charites:ignore directive"}
    IgnoreCheck -- "Ignored" --> Safe
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: performance.tailwind-untracked-package-source"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore performance.tailwind-untracked-package-source` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/performance.tailwind-untracked-package-source/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for performance.tailwind-untracked-package-source"]
        subgraph P ["Positive Corpus (tests/correctness/performance.tailwind-untracked-package-source/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/performance.tailwind-untracked-package-source/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/performance.tailwind-untracked-package-source/adversarial/)"]
            A1["A1: Template Literal Interpolations"]
            A2["A2: Ternary Conditional Expressions"]
            A3["A3: Spread Properties & Dynamic Overrides"]
            A4["A4: Dynamic Object Class Syntax"]
            A5["A5: Shadowed Variable Identifiers"]
            A6["A6: Nested Closures & HOC Wrappers"]
            A7["A7: Obfuscated Classes & Cyclic Tokens"]
        end
    end

    P --> TestRunner["Automated Runner (rule_test.go)"]
    N --> TestRunner
    A --> TestRunner
    TestRunner --> Gates["Quality Gates: Zero Panic, Zero False-Positive, Zero Bypass"]
```

- **Positive Fixtures (P1-P5):** Verified to trigger diagnostics at exact lines and column spans.
- **Negative Fixtures (N1-N5):** Verified to produce zero diagnostics on valid tokens and legitimate exemptions.
- **Adversarial Fixtures (A1-A7):** Verified to prevent evasion across dynamic expressions, string interpolations, and cyclic references.

---

## 8. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore performance.tailwind-untracked-package-source intentional exception -->
```

```tsx
// charites:ignore performance.tailwind-untracked-package-source intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.tailwind-untracked-package-source:
    severity: error # error | warn | info | off
```

