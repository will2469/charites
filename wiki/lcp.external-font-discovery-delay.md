# lcp.external-font-discovery-delay

> **Rule ID:** `lcp.external-font-discovery-delay`
> **Severity:** `WARN`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay), W3C Resource Hints (preconnect and dns-prefetch specification), Web Performance Working Group Connection Handshake Optimization Invariants

---

## 1. Overview & Core Invariant

External font stylesheet loaded without '<link rel="preconnect">' hints, adding 200ms-400ms connection handshake latency to LCP font discovery

### Core Invariant:
> **"External cross-origin font stylesheets must be preceded by '<link rel="preconnect">' hints to eliminate DNS, TCP, and TLS handshake round-trips before font binaries are requested."**

---
## 2. Technical Grounding & Engine Realities

Loading web fonts from third-party CDNs (such as Google Fonts or Adobe Typekit) involves a multi-origin dependency chain: the CSS stylesheet is fetched from one domain (e.g. 'fonts.googleapis.com'), while the font binaries (.woff2) are hosted on a separate storage origin (e.g. 'fonts.gstatic.com').

Without preconnect hints, the browser cannot initiate the DNS resolution, TCP three-way handshake, and TLS negotiation with the font storage origin until the stylesheet is completely downloaded and parsed.

Adding '<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>' allows the browser to perform socket setup in parallel during initial document streaming, shaving 200ms-400ms off LCP font discovery.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Connection Handshake Serialization** | HIGH | Sequential DNS/TCP/TLS handshakes to external font origins add 200ms-400ms round-trip latency to critical text LCP. |
| **Cellular Network Latency Amplification** | MEDIUM | On mobile 3G/4G networks with high RTT (Round Trip Time), serialized connection setup severely degrades user perceptual paint times. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### HTML (External Google Fonts imported without preconnect hints to stylesheet and font binary origins):
```html
<head>
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@700&display=swap" />
</head>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### HTML (Preconnect hints declared early to pre-warm DNS and TLS sockets before stylesheet parsing):
```html
<head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@700&display=swap" />
</head>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: lcp.external-font-discovery-delay"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore lcp.external-font-discovery-delay` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/lcp.external-font-discovery-delay/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for lcp.external-font-discovery-delay"]
        subgraph P ["Positive Corpus (tests/correctness/lcp.external-font-discovery-delay/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/lcp.external-font-discovery-delay/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/lcp.external-font-discovery-delay/adversarial/)"]
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
<!-- charites:ignore lcp.external-font-discovery-delay intentional exception -->
```

```tsx
// charites:ignore lcp.external-font-discovery-delay intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.external-font-discovery-delay:
    severity: warn # error | warn | info | off
```

