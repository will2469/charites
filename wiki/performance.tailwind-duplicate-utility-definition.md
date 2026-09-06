# performance.tailwind-duplicate-utility-definition

> **Rule ID:** `performance.tailwind-duplicate-utility-definition`
> **Severity:** `WARN`
> **Category:** `performance`
> **Target Standards:** Tailwind CSS v4 '@utility' Directive Specification, Compiled CSS Output Deduplication & Bundle Hygiene, Atomic CSS Design Invariants

---

## 1. Overview & Core Invariant

Mencegah duplikasi deklarasi utilitas CSS kustom (@utility) yang properti dan nilainya sudah disediakan oleh utilitas core bawaan Tailwind CSS v4.

### Core Invariant:
> **"Custom '@utility' declarations must not duplicate built-in Tailwind CSS core utilities; redundant definitions generate unnecessary stylesheet bytes and break atomic CSS composability."**

---
## 2. Technical Grounding & Engine Realities

The `@utility` directive in Tailwind CSS v4 is designed to register brand-new utilities for modern or proprietary CSS features not yet included in core Tailwind.

Defining custom `@utility` blocks for combinations already covered by core utilities (such as `@utility center-flex { display: flex; align-items: center; }`) produces duplicate CSS rules in the compiled stylesheet.

Composing canonical core utilities (`flex items-center`) directly in markup preserves atomic stylesheet economy and avoids redundant selector bloat.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Compiled Stylesheet Bloat** | MEDIUM | Adds unnecessary custom CSS rules to the production build that duplicate pre-existing atomic utilities. |
| **Bypassed Utility Composability** | LOW | Custom wrapper utilities fracture atomic consistency and make component classes harder to refactor. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### CSS (Mendefinisikan @utility yang menduplikasi utilitas core flexbox):
```css
@utility center-flex {
  display: flex;
  align-items: center;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Menggunakan kombinasi utilitas core native langsung di markup):
```tsx
<div className="flex items-center">Konten</div>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: performance.tailwind-duplicate-utility-definition"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore performance.tailwind-duplicate-utility-definition` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/performance.tailwind-duplicate-utility-definition/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for performance.tailwind-duplicate-utility-definition"]
        subgraph P ["Positive Corpus (tests/correctness/performance.tailwind-duplicate-utility-definition/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/performance.tailwind-duplicate-utility-definition/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/performance.tailwind-duplicate-utility-definition/adversarial/)"]
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
<!-- charites:ignore performance.tailwind-duplicate-utility-definition intentional exception -->
```

```tsx
// charites:ignore performance.tailwind-duplicate-utility-definition intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.tailwind-duplicate-utility-definition:
    severity: warn # error | warn | info | off
```

