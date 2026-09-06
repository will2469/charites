# browser.experimental-api-no-featuredetect

> **Rule ID:** `browser.experimental-api-no-featuredetect`
> **Severity:** `ERROR`
> **Category:** `browser`
> **Target Standards:** WICG Web Share API Specification (navigator.share), WICG File System Access API (showOpenFilePicker), W3C CSS View Transitions Module Level 1 (startViewTransition), ECMA-262 Feature Detection & Defensive JavaScript Guidelines

---

## 1. Overview & Core Invariant

Detects invocation of experimental Web APIs without runtime feature detection guards

### Core Invariant:
> **"Experimental or non-universal Web APIs must be guarded with runtime capability checks ('prop' in obj, if (obj.prop), optional chaining, or try/catch) before invocation."**

---
## 2. Technical Grounding & Engine Realities

Modern Web APIs are adopted unevenly across browser vendors. For instance, 'showOpenFilePicker' is exclusive to Chromium and crashes instantly on Firefox and Safari.

'navigator.share' throws an uncaught TypeError on desktop Firefox or non-secure contexts.

Directly calling these APIs without feature guards results in severe runtime exceptions that crash SPAs.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Runtime Exception Crash** | HIGH | Uncaught TypeError: undefined is not a function immediately terminates JavaScript execution in non-supporting browsers. |
| **Broken Core User Action** | HIGH | Primary user actions like document sharing or file importing fail completely without informative fallback UX. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Direct invocation of navigator.share without runtime feature guard (crashes on desktop Firefox)):
```tsx
<button onClick={() => {
  navigator.share({ title: "Surat", url: window.location.href });
}}>
  Bagikan
</button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Defensive feature detection with fallback to clipboard copy):
```tsx
<button onClick={async () => {
  if (typeof navigator !== "undefined" && navigator.share) {
    await navigator.share({ title: "Surat", url: window.location.href });
  } else {
    await navigator.clipboard?.writeText(window.location.href);
  }
}}>
  Bagikan
</button>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: browser.experimental-api-no-featuredetect"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore browser.experimental-api-no-featuredetect` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/browser.experimental-api-no-featuredetect/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for browser.experimental-api-no-featuredetect"]
        subgraph P ["Positive Corpus (tests/correctness/browser.experimental-api-no-featuredetect/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/browser.experimental-api-no-featuredetect/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/browser.experimental-api-no-featuredetect/adversarial/)"]
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
<!-- charites:ignore browser.experimental-api-no-featuredetect intentional exception -->
```

```tsx
// charites:ignore browser.experimental-api-no-featuredetect intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  browser.experimental-api-no-featuredetect:
    severity: error # error | warn | info | off
```

