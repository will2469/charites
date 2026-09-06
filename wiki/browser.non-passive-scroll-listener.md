# browser.non-passive-scroll-listener

> **Rule ID:** `browser.non-passive-scroll-listener`
> **Severity:** `WARN`
> **Category:** `browser`
> **Target Standards:** W3C DOM Level 4 Events Specification (Passive Event Listeners), Chromium & WebKit Compositor Scrolling Pipeline Guidelines, Google Lighthouse Best Practices (Does not use passive listeners)

---

## 1. Overview & Core Invariant

Enforces { passive: true } option on touch and wheel event listeners to prevent main thread scroll blocking

### Core Invariant:
> **"Event listeners for 'touchstart', 'touchmove', 'wheel', or 'mousewheel' must declare '{ passive: true }' to ensure non-blocking compositor scrolling."**

---
## 2. Technical Grounding & Engine Realities

Browsers execute smooth scrolling on a dedicated compositor thread. Without '{ passive: true }', the compositor must block and wait for JavaScript execution on the main thread to see if 'preventDefault()' is called.

This introduces severe touch response latency and frame rate drops (scroll jank) on mobile devices.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Scroll Jank & Latency** | MEDIUM | Users experience jerky, lagging scrolling and delayed touch gestures on Safari iOS and Android Chrome. |
| **Lighthouse Performance Penalty** | LOW | Fails Lighthouse 'Does not use passive listeners to improve scrolling performance' audit. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### JAVASCRIPT (Adding touchmove listener without passive: true option):
```javascript
window.addEventListener("touchmove", (e) => {
  trackTouchPosition(e);
});
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### JAVASCRIPT (Specifying { passive: true } to unblock the compositor thread):
```javascript
window.addEventListener("touchmove", (e) => {
  trackTouchPosition(e);
}, { passive: true });
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: browser.non-passive-scroll-listener"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore browser.non-passive-scroll-listener` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/browser.non-passive-scroll-listener/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for browser.non-passive-scroll-listener"]
        subgraph P ["Positive Corpus (tests/correctness/browser.non-passive-scroll-listener/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/browser.non-passive-scroll-listener/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/browser.non-passive-scroll-listener/adversarial/)"]
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
<!-- charites:ignore browser.non-passive-scroll-listener intentional exception -->
```

```tsx
// charites:ignore browser.non-passive-scroll-listener intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  browser.non-passive-scroll-listener:
    severity: warn # error | warn | info | off
```

