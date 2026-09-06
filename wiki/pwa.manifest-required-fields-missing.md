# pwa.manifest-required-fields-missing

> **Rule ID:** `pwa.manifest-required-fields-missing`
> **Severity:** `ERROR`
> **Category:** `pwa`
> **Target Standards:** W3C Web App Manifest Specification Section 5 (Manifest Members), Google Chrome Web App Installability Criteria, W3C Application Lifecycle & Installation Architecture

---

## 1. Overview & Core Invariant

Errors when a Web App Manifest definition is missing required fields (name/short_name, start_url, display, icons)

### Core Invariant:
> **"Web App Manifest declarations (<script type="application/manifest+json">) must declare required installability fields: 'name' (or 'short_name'), 'start_url', 'display', and at least one icon in 'icons'."**

---
## 2. Technical Grounding & Engine Realities

For mobile operating systems (Android, iOS) and modern web engines to recognize a website as an installable Progressive Web App, the Web App Manifest must declare minimum installability metadata.

Omitting 'name' or 'short_name' leaves the OS homescreen launcher with an empty or broken app label. Missing 'start_url' prevents the launcher from knowing which route to boot into. Omitting 'display' forces the app to open inside standard browser tabs rather than immersive standalone mode. Omitting 'icons' results in missing or broken asset icons on the user's home screen.

Declaring all four fundamental fields ensures reliable PWA install prompts and clean native OS integration.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Browser Installation Prompt Suppression** | HIGH | Mobile browsers silently suppress the 'Add to Home Screen' / install banner when manifest required fields are absent. |
| **Broken Application Branding** | MEDIUM | If installed manually, the web application displays placeholder text and fallback generic browser icons. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Manifest missing start_url, display, and icons):
```tsx
<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital"
  })}
</script>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Manifest declares all required installability members):
```tsx
<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital",
    short_name: "Desa",
    start_url: "/",
    display: "standalone",
    icons: [
      { src: "/icon-192.png", sizes: "192x192", type: "image/png" },
      { src: "/icon-512.png", sizes: "512x512", type: "image/png", purpose: "maskable" }
    ]
  })}
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: pwa.manifest-required-fields-missing"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore pwa.manifest-required-fields-missing` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/pwa.manifest-required-fields-missing/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for pwa.manifest-required-fields-missing"]
        subgraph P ["Positive Corpus (tests/correctness/pwa.manifest-required-fields-missing/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/pwa.manifest-required-fields-missing/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/pwa.manifest-required-fields-missing/adversarial/)"]
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
<!-- charites:ignore pwa.manifest-required-fields-missing intentional exception -->
```

```tsx
// charites:ignore pwa.manifest-required-fields-missing intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  pwa.manifest-required-fields-missing:
    severity: error # error | warn | info | off
```

