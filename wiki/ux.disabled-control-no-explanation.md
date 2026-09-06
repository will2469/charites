# ux.disabled-control-no-explanation

> **Rule ID:** `ux.disabled-control-no-explanation`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Nielsen Heuristic #1: Visibility of System Status, Feedforward Principle & Gulf of Evaluation (Don Norman), WCAG 2.2 Success Criterion 3.3.2 (Labels or Instructions)

---

## 1. Overview & Core Invariant

Enforces feedforward explanation for disabled interactive controls to prevent user dead ends

### Core Invariant:
> **"Disabled interactive controls ('disabled', 'aria-disabled="true"') must provide a feedforward explanation via 'aria-describedby', tooltip, or visible helper text."**

---
## 2. Technical Grounding & Engine Realities

When users encounter a disabled button or locked form control without any visible reason, they reach a cognitive dead end: they cannot proceed and have no information on how to unlock the action.

Providing contextual feedforward (e.g. 'Minimum belanja Rp 50.000 untuk checkout' or an explanatory tooltip) clarifies system constraints and transforms frustration into actionable user guidance.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cognitive Dead Ends & Workflow Abandonment** | MEDIUM | Users are unable to diagnose why an action is blocked and abandon conversion or submission workflows. |
| **Accessibility Barrier (Screen Reader Confusion)** | MEDIUM | Assisted technology users hear 'button disabled' without auditory explanation of required missing inputs. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Disabled checkout button without explanation of minimum order constraint):
```tsx
<div className="mt-4">
  <button disabled={cartTotal < 50000} className="bg-primary text-white px-4 py-2 rounded">
    Checkout
  </button>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Disabled button accompanied by aria-describedby linking to explanatory hint text):
```tsx
<div className="mt-4">
  <button
    disabled={cartTotal < 50000}
    aria-describedby="min-order-hint"
    className="bg-primary text-white px-4 py-2 rounded disabled:opacity-50"
  >
    Checkout
  </button>
  {cartTotal < 50000 && (
    <p id="min-order-hint" className="text-xs text-muted-foreground mt-1">
      Minimum belanja Rp 50.000 untuk melanjutkan checkout.
    </p>
  )}
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: ux.disabled-control-no-explanation"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore ux.disabled-control-no-explanation` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/ux.disabled-control-no-explanation/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for ux.disabled-control-no-explanation"]
        subgraph P ["Positive Corpus (tests/correctness/ux.disabled-control-no-explanation/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/ux.disabled-control-no-explanation/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/ux.disabled-control-no-explanation/adversarial/)"]
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
<!-- charites:ignore ux.disabled-control-no-explanation intentional exception -->
```

```tsx
// charites:ignore ux.disabled-control-no-explanation intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.disabled-control-no-explanation:
    severity: warn # error | warn | info | off
```

