# performance.astro-over-prefetching

> **Rule ID:** `performance.astro-over-prefetching`
> **Severity:** `WARN`
> **Category:** `performance`
> **Target Standards:** Astro Prefetch Configuration Best Practices ('data-astro-prefetch'), W3C Resource Hints & Speculative Parsing Bandwidth Economy, Mobile Web Data Saver & Cellular Network Latency Guidelines

---

## 1. Overview & Core Invariant

Mencegah pemborosan kuota data seluler dengan melarang penempatan strategi prefetch agresif (viewport/load) pada tautan navigasi sekunder atau footer.

### Core Invariant:
> **"Aggressive 'viewport' or 'load' prefetch strategies must not be assigned to secondary or low-conversion navigation links; secondary links should use passive 'hover' or 'tap' prefetching."**

---
## 2. Technical Grounding & Engine Realities

Astro provides link prefetching via `data-astro-prefetch`.

Using aggressive strategies like `data-astro-prefetch="viewport"` causes the browser to immediately fetch all linked documents as soon as their anchors enter the viewport.

When applied to secondary links (such as legal terms, privacy policies, or footer menus), this aggressively consumes user bandwidth and saturates the network connection, starving critical assets such as images and analytical payloads on slow mobile networks.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cellular Data Waste** | MEDIUM | Preemptively downloads full pages that users rarely click, depleting metered mobile data connections. |
| **Network Queue Contention** | MEDIUM | Prefetch network requests crowd the HTTP queue and delay high-priority above-the-fold assets. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Prefetch agresif pada tautan footer sekunder):
```astro
<footer>
  <a href="/terms" data-astro-prefetch="viewport">Syarat & Ketentuan</a>
  <a href="/privacy" data-astro-prefetch="viewport">Kebijakan Privasi</a>
</footer>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Prefetch pasif saat hover untuk tautan footer):
```astro
<footer>
  <a href="/terms" data-astro-prefetch="hover">Syarat & Ketentuan</a>
  <a href="/privacy" data-astro-prefetch="hover">Kebijakan Privasi</a>
</footer>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: performance.astro-over-prefetching"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore performance.astro-over-prefetching` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/performance.astro-over-prefetching/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for performance.astro-over-prefetching"]
        subgraph P ["Positive Corpus (tests/correctness/performance.astro-over-prefetching/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/performance.astro-over-prefetching/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/performance.astro-over-prefetching/adversarial/)"]
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
<!-- charites:ignore performance.astro-over-prefetching intentional exception -->
```

```tsx
// charites:ignore performance.astro-over-prefetching intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.astro-over-prefetching:
    severity: warn # error | warn | info | off
```

