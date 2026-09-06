# lcp.missing-critical-origin-hint

> **Rule ID:** `lcp.missing-critical-origin-hint`
> **Severity:** `INFO`
> **Category:** `lcp`
> **Target Standards:** Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay), W3C Resource Hints (preconnect and dns-prefetch specification), Web Performance Working Group Network Socket Pre-warming Guidelines

---

## 1. Overview & Core Invariant

Critical LCP visual asset loaded from third-party CDN origin without '<link rel="preconnect">' connection hint in '<head>'

### Core Invariant:
> **"Critical visual LCP assets hosted on external cross-origin domains should have an early '<link rel="preconnect">' hint in '<head>' to eliminate DNS, TCP, and TLS socket handshake round-trips."**

---
## 2. Technical Grounding & Engine Realities

When an above-the-fold hero image is loaded from a third-party domain (e.g. 'images.unsplash.com', 'cdn.shopify.com', or 'res.cloudinary.com'), the browser cannot initiate the network connection until the `<img>` tag is discovered and parsed.

Establishing a new HTTPS connection to an external origin requires three sequential round-trips: DNS resolution, TCP three-way handshake, and TLS cryptographic negotiation, adding 150ms to 400ms of idle latency on cellular connections.

Declaring `<link rel="preconnect" href="https://images.unsplash.com" />` in `<head>` instructs the browser to open the socket in the background during early HTML document streaming.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Connection Handshake Latency Spike** | MEDIUM | Adds 150ms-400ms of socket negotiation delay to LCP Resource Load Delay before image bytes start streaming. |
| **Cellular Round-Trip Time Penalty** | LOW | Multi-RTT connection setup noticeably degrades mobile performance metrics in emerging market networks. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### HTML (Critical hero image loaded from external CDN without preconnect hint in head):
```html
<head>
  <title>Product Showcase</title>
</head>
<body>
  <img src="https://images.unsplash.com/photo-hero" fetchpriority="high" data-perf-role="hero" />
</body>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### HTML (Preconnect hint declared in head to pre-warm CDN connection sockets early):
```html
<head>
  <title>Product Showcase</title>
  <link rel="preconnect" href="https://images.unsplash.com" />
</head>
<body>
  <img src="https://images.unsplash.com/photo-hero" fetchpriority="high" data-perf-role="hero" />
</body>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: lcp.missing-critical-origin-hint"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore lcp.missing-critical-origin-hint` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/lcp.missing-critical-origin-hint/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for lcp.missing-critical-origin-hint"]
        subgraph P ["Positive Corpus (tests/correctness/lcp.missing-critical-origin-hint/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/lcp.missing-critical-origin-hint/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/lcp.missing-critical-origin-hint/adversarial/)"]
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
<!-- charites:ignore lcp.missing-critical-origin-hint intentional exception -->
```

```tsx
// charites:ignore lcp.missing-critical-origin-hint intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  lcp.missing-critical-origin-hint:
    severity: info # error | warn | info | off
```

