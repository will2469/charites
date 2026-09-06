# browser.hover-only-interaction

> **Rule ID:** `browser.hover-only-interaction`
> **Severity:** `ERROR`
> **Category:** `browser`
> **Target Standards:** W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 2.1.1 (Keyboard), WICG / WHATWG Touch Events & Pointer Events Level 3 (Touch vs Hover Ergonomics), Apple Human Interface Guidelines for iOS Touch Interactions

---

## 1. Overview & Core Invariant

Ensures interactive actions and state reveals have keyboard and touch counterparts instead of relying solely on hover

### Core Invariant:
> **"Interactive controls and revealed elements must not rely exclusively on ':hover' or 'group-hover:' without keyboard/touch counterparts ('focus-visible:', 'group-focus-within:')."**

---
## 2. Technical Grounding & Engine Realities

Touchscreen devices (the majority of web traffic on Safari iOS and Chrome Android) have no physical cursor and cannot perform genuine mouse hover.

When critical action buttons (e.g. delete, edit, copy) are hidden by default with 'opacity-0' and only revealed via 'group-hover:opacity-100':
1. Total Mobile Inaccessibility: Touchscreen users cannot see or activate the controls because hovering does not exist on mobile.
2. iOS Sticky Hover Bug: Tapping an element on Safari iOS triggers an inconsistent 'sticky hover' state, requiring multiple confusing taps.
3. Keyboard Navigation Failure: Users navigating with the 'Tab' key cannot discover or focus hidden controls unless focus-within or focus-visible is bound.

Charites enforces that any hover-revealed element provides accessible keyboard and touch parity.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Touch Exclusion** | HIGH | Critical action controls are completely invisible and unreachable on smartphones and tablets. |
| **Keyboard Accessibility Barrier** | MEDIUM | Fails WCAG 2.2 Level A keyboard navigation audits when hidden controls cannot be focused with the Tab key. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Delete button hidden by default and only revealed on group-hover (invisible on touchscreens)):
```tsx
<div className="group flex items-center justify-between p-3 border rounded-xl">
  <span>Berkas_KTP.pdf</span>
  <button
    type="button"
    onClick={handleDelete}
    className="opacity-0 group-hover:opacity-100 text-destructive text-sm"
  >
    Hapus
  </button>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Button accessible via hover, keyboard Tab navigation, and touch focus):
```tsx
<div className="group flex items-center justify-between p-3 border rounded-xl">
  <span>Berkas_KTP.pdf</span>
  <button
    type="button"
    onClick={handleDelete}
    className="opacity-0 group-hover:opacity-100 group-focus-within:opacity-100 focus-visible:opacity-100 text-destructive text-sm"
  >
    Hapus
  </button>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: browser.hover-only-interaction"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore browser.hover-only-interaction` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/browser.hover-only-interaction/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for browser.hover-only-interaction"]
        subgraph P ["Positive Corpus (tests/correctness/browser.hover-only-interaction/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/browser.hover-only-interaction/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/browser.hover-only-interaction/adversarial/)"]
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
<!-- charites:ignore browser.hover-only-interaction intentional exception -->
```

```tsx
// charites:ignore browser.hover-only-interaction intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  browser.hover-only-interaction:
    severity: error # error | warn | info | off
```

