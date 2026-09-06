# inp.large-interaction-layout-scope

> **Rule ID:** `inp.large-interaction-layout-scope`
> **Severity:** `WARN`
> **Category:** `inp`
> **Target Standards:** Google Chrome Core Web Vitals (Interaction to Next Paint Presentation Delay), W3C CSS Containment Module Level 3 (contain: layout / contain: strict), HTML Living Standard HTMLDialogElement Top-Layer Architecture

---

## 1. Overview & Core Invariant

Interactive overlay or drawer element lacks layout containment or native dialog isolation, triggering document-wide reflow on toggle

### Core Invariant:
> **"Large interactive overlays and drawers must establish layout containment ('contain: layout') or utilize the browser top-layer ('<dialog>') to prevent whole-page layout recalculations during interactions."**

---
## 2. Technical Grounding & Engine Realities

When a large overlay, slide-over drawer, or modal toggles its visibility in the standard document flow, the browser's layout engine must invalidate ancestor and sibling boxes, triggering a full document reflow.

For complex interfaces with thousands of elements, this layout recalculation stalls the main thread and inflates Presentation Delay well past 100ms.

Using the HTML5 '<dialog>' element places the modal in the browser's isolated 'top-layer', preventing any layout impact on the document tree. Alternatively, applying CSS layout containment (e.g. 'contain-layout' or '[contain:layout]') constrains layout recalculations strictly inside the overlay container.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Document-Wide Layout Invalidation** | HIGH | Toggling modal/drawer state forces the layout engine to recalculate geometry for every element on the page. |
| **Presentation Delay Frame Drops** | HIGH | Users experience visible stutters and sluggishness when expanding or collapsing sidebars, sheets, or dialogs. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Unconstrained fixed drawer in normal document flow triggering document reflow):
```tsx
<div className={`fixed inset-y-0 right-0 w-96 ${isOpen ? "block" : "hidden"}`}>
  <HeavySidebar />
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Native HTML5 dialog rendered in the browser's isolated top-layer):
```tsx
<dialog ref={dialogRef} className="fixed inset-y-0 right-0 w-96">
  <HeavySidebar />
</dialog>
```
### TSX (Explicit CSS layout containment isolating reflows to the panel):
```tsx
<div className={`fixed inset-y-0 right-0 w-96 contain-layout ${isOpen ? "block" : "hidden"}`}>
  <HeavySidebar />
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: inp.large-interaction-layout-scope"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore inp.large-interaction-layout-scope` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/inp.large-interaction-layout-scope/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for inp.large-interaction-layout-scope"]
        subgraph P ["Positive Corpus (tests/correctness/inp.large-interaction-layout-scope/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/inp.large-interaction-layout-scope/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/inp.large-interaction-layout-scope/adversarial/)"]
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
<!-- charites:ignore inp.large-interaction-layout-scope intentional exception -->
```

```tsx
// charites:ignore inp.large-interaction-layout-scope intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  inp.large-interaction-layout-scope:
    severity: warn # error | warn | info | off
```

