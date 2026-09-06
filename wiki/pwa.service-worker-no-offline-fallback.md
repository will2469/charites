# pwa.service-worker-no-offline-fallback

> **Rule ID:** `pwa.service-worker-no-offline-fallback`
> **Severity:** `WARN`
> **Category:** `pwa`
> **Target Standards:** W3C Service Workers 1 (Offline Resilience Architecture), W3C Cache Storage Specification (Offline Asset Fallback), Google Chrome PWA Reliability Criteria (Offline Support)

---

## 1. Overview & Core Invariant

Warns when a Service Worker intercepts fetch events without providing an offline cache fallback or failure handler

### Core Invariant:
> **"Service Worker fetch event handlers must implement an offline cache fallback (e.g. caches.match) or failure catch handler instead of bare pass-through fetch interception."**

---
## 2. Technical Grounding & Engine Realities

In spotty or rural mobile network conditions (3G/4G signal drops), a Service Worker that intercepts fetch events without an offline cache strategy causes the browser to immediately display a connection-lost screen (the offline dinosaur page).

Providing a resilient cache-first or network-first fallback mechanism guarantees that the application shell and cached pages remain accessible even when completely disconnected from the Internet.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Immediate Offline Blackout** | MEDIUM | Users opening the PWA without network connectivity encounter a browser network failure screen instead of cached application content. |
| **PWA Installability Rejection** | LOW | Mobile browsers may downgrade or reject full PWA installation status due to failing offline resilience audits. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Pass-through fetch interception without offline cache fallback):
```tsx
<script>
  self.addEventListener("fetch", (event) => {
    event.respondWith(fetch(event.request));
  });
</script>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Cache fallback provided via caches.match with network fallback):
```tsx
<script>
  self.addEventListener("fetch", (event) => {
    event.respondWith(
      caches.match(event.request).then((cached) => {
        return cached || fetch(event.request).catch(() => caches.match("/offline.html"));
      })
    );
  });
</script>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: pwa.service-worker-no-offline-fallback"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore pwa.service-worker-no-offline-fallback` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/pwa.service-worker-no-offline-fallback/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for pwa.service-worker-no-offline-fallback"]
        subgraph P ["Positive Corpus (tests/correctness/pwa.service-worker-no-offline-fallback/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/pwa.service-worker-no-offline-fallback/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/pwa.service-worker-no-offline-fallback/adversarial/)"]
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
<!-- charites:ignore pwa.service-worker-no-offline-fallback intentional exception -->
```

```tsx
// charites:ignore pwa.service-worker-no-offline-fallback intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  pwa.service-worker-no-offline-fallback:
    severity: warn # error | warn | info | off
```

