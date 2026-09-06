# cls.client-only-hydration-pop

> **Rule ID:** `cls.client-only-hydration-pop`
> **Severity:** `WARN`
> **Category:** `cls`
> **Target Standards:** Astro Islands Architecture (client:only directives & fallback slots), W3C Core Web Vitals (Cumulative Layout Shift Prevention), Progressive Enhancement & Skeleton Shell Invariants

---

## 1. Overview & Core Invariant

Astro client:only island lacks a slot='fallback' shell or reserved min-height container, causing hydration layout shift

### Core Invariant:
> **"Astro components utilizing 'client:only' must define an official fallback shell (<div slot='fallback'>) or be enclosed within a container with reserved min-height."**

---
## 2. Technical Grounding & Engine Realities

In Astro's island architecture, the 'client:only' directive explicitly opts out of server-side rendering (SSR), omitting initial HTML markup for the component during build time.

Without a server-rendered placeholder or designated fallback shell, the browser initially renders an empty 0-height space. When the client-side JavaScript bundle finishes downloading, parsing, and executing, the rendered component abruptly expands and pushes all subsequent document content downward.

Providing a dedicated fallback shell via '<div slot="fallback" class="min-h-[...]">...</div>' ensures that the space is permanently reserved in initial server HTML, completely neutralizing Cumulative Layout Shift upon client hydration.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Post-Hydration Content Displacement** | HIGH | Delayed hydration of client-only islands causes sudden vertical document jumping when interactive components finish booting. |
| **Blank Hole Flash** | MEDIUM | Users experience an empty white space where interactive widgets or charts belong prior to JavaScript execution. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (client:only island without fallback slot or reserved height):
```astro
<main class="space-y-4">
  <h1>Dashboard</h1>
  <AnalyticsChart client:only="react" />
  <p>Live stats</p>
</main>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (client:only island with dedicated fallback slot shell):
```astro
<main class="space-y-4">
  <h1>Dashboard</h1>
  <AnalyticsChart client:only="react">
    <div slot="fallback" class="w-full min-h-[350px] bg-muted/20 animate-pulse rounded-lg flex items-center justify-center">
      <span>Memuat grafik...</span>
    </div>
  </AnalyticsChart>
  <p>Live stats</p>
</main>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: cls.client-only-hydration-pop"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore cls.client-only-hydration-pop` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/cls.client-only-hydration-pop/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for cls.client-only-hydration-pop"]
        subgraph P ["Positive Corpus (tests/correctness/cls.client-only-hydration-pop/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/cls.client-only-hydration-pop/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/cls.client-only-hydration-pop/adversarial/)"]
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
<!-- charites:ignore cls.client-only-hydration-pop intentional exception -->
```

```tsx
// charites:ignore cls.client-only-hydration-pop intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.client-only-hydration-pop:
    severity: warn # error | warn | info | off
```

