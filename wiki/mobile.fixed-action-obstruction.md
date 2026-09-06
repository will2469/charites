# mobile.fixed-action-obstruction

> **Rule ID:** `mobile.fixed-action-obstruction`
> **Severity:** `WARN`
> **Category:** `mobile`
> **Target Standards:** Apple Human Interface Guidelines (Bottom Toolbars & Screen Clearance), Google Material Design 3 (Bottom App Bars & Safe Content Boundaries), W3C CSS Positioned Layout Module Level 3 (Fixed Positioning)

---

## 1. Overview & Core Invariant

Warns when fixed bottom elements lack compensating bottom padding on parent or content siblings, risking content obstruction

### Core Invariant:
> **"Fixed bottom bars and floating action buttons must be accompanied by compensating bottom padding ('pb-16', 'pb-20', 'pb-24', 'pb-safe') on parent layouts or content siblings to prevent content obstruction."**

---
## 2. Technical Grounding & Engine Realities

Elements anchored with 'fixed bottom-0' float out of normal document flow, permanently covering the lower portion of the viewport.

Without compensating bottom padding (such as 'pb-24' or 'pb-[env(safe-area-inset-bottom)]') on the layout container or content siblings (<main>, <article>, <form>), the final rows of text, interactive inputs, or submit buttons will be permanently hidden behind the fixed bar.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Obstructed Content & Trapped Form Inputs** | MEDIUM | The bottom-most form fields or submit controls are occluded by the fixed bar, blocking user progress. |
| **Accidental Clicks on Fixed Bar** | LOW | Users attempting to tap the bottom of the page accidentally trigger bottom navigation items instead. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Fixed bottom nav without compensating bottom padding on main content):
```tsx
<div className="min-h-screen bg-background">
  <main className="p-4 space-y-4">
    <p>Konten formulir paling bawah...</p>
  </main>
  <nav className="fixed bottom-0 inset-x-0 h-16 bg-card border-t">
    <button type="button">Beranda</button>
  </nav>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Compensating bottom padding (pb-24) ensures full content clearance):
```tsx
<div className="min-h-screen bg-background">
  <main className="p-4 space-y-4 pb-24">
    <p>Konten formulir paling bawah...</p>
  </main>
  <nav className="fixed bottom-0 inset-x-0 h-16 bg-card border-t pb-[env(safe-area-inset-bottom)]">
    <button type="button">Beranda</button>
  </nav>
</div>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: mobile.fixed-action-obstruction"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore mobile.fixed-action-obstruction` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/mobile.fixed-action-obstruction/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for mobile.fixed-action-obstruction"]
        subgraph P ["Positive Corpus (tests/correctness/mobile.fixed-action-obstruction/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/mobile.fixed-action-obstruction/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/mobile.fixed-action-obstruction/adversarial/)"]
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
<!-- charites:ignore mobile.fixed-action-obstruction intentional exception -->
```

```tsx
// charites:ignore mobile.fixed-action-obstruction intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  mobile.fixed-action-obstruction:
    severity: warn # error | warn | info | off
```

