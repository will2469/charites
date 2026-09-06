# responsive.keyboard-obstruction

> **Rule ID:** `responsive.keyboard-obstruction`
> **Severity:** `WARN`
> **Category:** `responsive`
> **Target Standards:** WCAG 2.2 Guideline 2.1 (Keyboard Accessible), Material Design 3 Mobile Form Guidelines, iOS Human Interface Guidelines (Managing the Virtual Keyboard)

---

## 1. Overview & Core Invariant

Warns against fixed bottom action bars in forms lacking vertical scroll containers, which can be obstructed by the mobile virtual keyboard

### Core Invariant:
> **"Forms containing text inputs and fixed/sticky bottom action bars must provide a scrollable container ('overflow-y-auto') so inputs and actions are never obscured when the virtual keyboard expands."**

---
## 2. Technical Grounding & Engine Realities

When a user taps an input on a smartphone, the virtual software keyboard slides up from the bottom of the screen, consuming 40% to 50% of the visible viewport.

Elements styled with 'fixed bottom-0' or 'sticky bottom-0' remain pinned above the viewport bottom or above the keyboard. If the parent form is not wrapped in a vertical scroll container ('overflow-y-auto'), the active input field gets pushed behind the keyboard or under the fixed button, leaving the user unable to view their input or complete submission.

Charites enforces scrollable viewport resilience for mobile form layouts.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Virtual Keyboard Input Obstruction** | HIGH | Users cannot see text being typed into lower form inputs because fixed bottom bars pin directly over them. |
| **Form Abandonment & Submission Blockers** | MEDIUM | When keyboard expansion pushes inputs offscreen without scroll capabilities, conversion rates drop. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Fixed bottom submit bar in a rigid form lacking scrollable region):
```tsx
<form className="h-screen flex flex-col justify-between">
  <div className="p-4 space-y-4">
    <input type="text" placeholder="Nama Lengkap" />
    <input type="email" placeholder="Alamat Surel" />
    <textarea placeholder="Pesan Anda" />
  </div>
  <div className="fixed bottom-0 inset-x-0 p-4 bg-surface border-t">
    <button type="submit" className="w-full bg-primary text-white py-3 rounded">Kirim</button>
  </div>
</form>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Scrollable body with fixed bottom bar allowing smooth mobile keyboard reflow):
```tsx
<form className="h-screen flex flex-col">
  <div className="flex-1 overflow-y-auto p-4 space-y-4">
    <input type="text" placeholder="Nama Lengkap" />
    <input type="email" placeholder="Alamat Surel" />
    <textarea placeholder="Pesan Anda" />
  </div>
  <div className="p-4 bg-surface border-t">
    <button type="submit" className="w-full bg-primary text-white py-3 rounded">Kirim</button>
  </div>
</form>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: responsive.keyboard-obstruction"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore responsive.keyboard-obstruction` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/responsive.keyboard-obstruction/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for responsive.keyboard-obstruction"]
        subgraph P ["Positive Corpus (tests/correctness/responsive.keyboard-obstruction/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/responsive.keyboard-obstruction/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/responsive.keyboard-obstruction/adversarial/)"]
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
<!-- charites:ignore responsive.keyboard-obstruction intentional exception -->
```

```tsx
// charites:ignore responsive.keyboard-obstruction intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  responsive.keyboard-obstruction:
    severity: warn # error | warn | info | off
```

