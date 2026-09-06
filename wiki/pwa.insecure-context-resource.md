# pwa.insecure-context-resource

> **Rule ID:** `pwa.insecure-context-resource`
> **Severity:** `ERROR`
> **Category:** `pwa`
> **Target Standards:** W3C Secure Contexts Specification, W3C Mixed Content Level 2 Specification, RFC 7258 Pervasive Monitoring Is an Attack

---

## 1. Overview & Core Invariant

Errors when a resource element loads assets over an insecure HTTP protocol (http://) in violation of W3C Secure Contexts

### Core Invariant:
> **"Resource elements (<script>, <link>, <img>, <iframe>, <video>, <audio>) must not load external assets over insecure 'http://' (except localhost loopback)."**

---
## 2. Technical Grounding & Engine Realities

Progressive Web Apps strictly require a Secure Context (HTTPS) to enable service workers, cache storage, and device hardware APIs.

Loading external assets (scripts, stylesheets, images, media, iframes) via an unencrypted 'http://' connection triggers Mixed Content blocking. Active mixed content (scripts, stylesheets) is blocked immediately by modern mobile browsers, breaking application functionality. Passive mixed content (images, audio) generates security warnings and can be intercepted or tampered with on public Wi-Fi networks.

All asset references must use HTTPS ('https://'), protocol-relative URLs ('//'), or local origin paths. Localhost addresses ('http://localhost' and 'http://127.0.0.1') are excepted for development purposes.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Active Mixed Content Blocking** | HIGH | Browsers completely block insecure scripts and stylesheets, breaking UI styling and interactive application logic. |
| **Man-in-the-Middle Asset Tampering** | HIGH | Unencrypted HTTP traffic can be intercepted, inspected, or modified by malicious actors on untrusted networks. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Insecure HTTP external script and stylesheet links):
```tsx
<div>
  <script src="http://cdn.example.org/tracker.js" />
  <link rel="stylesheet" href="http://assets.desa.id/styles.css" />
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Secure HTTPS asset loading conforming to Secure Contexts):
```tsx
<div>
  <script src="https://cdn.example.org/tracker.js" />
  <link rel="stylesheet" href="https://assets.desa.id/styles.css" />
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: pwa.insecure-context-resource"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore pwa.insecure-context-resource` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/pwa.insecure-context-resource/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for pwa.insecure-context-resource"]
        subgraph P ["Positive Corpus (tests/correctness/pwa.insecure-context-resource/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/pwa.insecure-context-resource/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/pwa.insecure-context-resource/adversarial/)"]
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
<!-- charites:ignore pwa.insecure-context-resource intentional exception -->
```

```tsx
// charites:ignore pwa.insecure-context-resource intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  pwa.insecure-context-resource:
    severity: error # error | warn | info | off
```

