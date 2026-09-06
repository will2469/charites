# ux.competing-primary-cta

> **Rule ID:** `ux.competing-primary-cta`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Von Restorff Effect (The Isolation Effect / Visual Dominance), Hick-Hyman Law (Logarithmic Decision Latency), Nielsen Norman Group (Visual Hierarchy for Action Buttons)

---

## 1. Overview & Core Invariant

Warns when an action group or interactive container contains more than one primary call-to-action button

### Core Invariant:
> **"An action container or button group must contain at most one primary call-to-action button, ensuring a clear visual focal point and zero decision ambiguity."**

---
## 2. Technical Grounding & Engine Realities

The Von Restorff Effect (Isolation Effect) predicts that when multiple similar items are presented, the one that differs from the rest is most likely to be remembered and acted upon.

When an interface presents two or more buttons styled identically as primary actions (e.g., two 'bg-primary' or 'variant="primary"' buttons side by side in a modal footer or form actions), visual hierarchy collapses.

This competing prominence causes choice paralysis (Hick-Hyman Law), forces users to pause and re-read labels carefully, and drastically increases the probability of accidental mis-clicks. Every decision context must have exactly one visually distinct primary action, while supporting actions should be styled with 'outline', 'secondary', or 'ghost' variants.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Choice Paralysis & Decision Latency** | HIGH | Users hesitate when confronted with equal-weight visual cues, increasing conversion drop-off rates. |
| **Accidental Action Slips** | MEDIUM | Users mistake secondary or cancel triggers for primary confirmation due to identical color and elevation styling. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Dialog footer with two competing primary buttons creating visual ambiguity):
```tsx
<div className="flex justify-end gap-3 p-4">
  <button type="button" className="bg-primary text-white px-4 py-2 rounded-md">Simpan Draf</button>
  <button type="submit" className="bg-primary text-white px-4 py-2 rounded-md">Publikasikan</button>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Clear hierarchy: one primary button paired with a secondary outline button):
```tsx
<div className="flex justify-end gap-3 p-4">
  <button type="button" className="border border-input bg-transparent px-4 py-2 rounded-md">Simpan Draf</button>
  <button type="submit" className="bg-primary text-white px-4 py-2 rounded-md">Publikasikan</button>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: ux.competing-primary-cta"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore ux.competing-primary-cta` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/ux.competing-primary-cta/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for ux.competing-primary-cta"]
        subgraph P ["Positive Corpus (tests/correctness/ux.competing-primary-cta/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/ux.competing-primary-cta/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/ux.competing-primary-cta/adversarial/)"]
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
<!-- charites:ignore ux.competing-primary-cta intentional exception -->
```

```tsx
// charites:ignore ux.competing-primary-cta intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.competing-primary-cta:
    severity: warn # error | warn | info | off
```

