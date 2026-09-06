# ux.monolithic-form-bloat

> **Rule ID:** `ux.monolithic-form-bloat`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Sweller's Cognitive Load Theory (Intrinsic & Germane Cognitive Load Management), Progressive Task Completion & Form Usability Principles (Wroblewski & Nielsen Norman Group), W3C WAI-ARIA Authoring Practices Guide (Form Landmark & Fieldset Segmentation)

---

## 1. Overview & Core Invariant

Warns when a monolithic form contains excessive unchunked inputs (> 9 total or > 7 per chunk), violating Cognitive Load Theory

### Core Invariant:
> **"Forms containing more than 9 total interactive fields must segment fields into chunks ('<fieldset>', Stepper, or Tabs), with no single chunk exceeding 7 fields."**

---
## 2. Technical Grounding & Engine Realities

Long, monolithic forms containing 10 or more unorganized inputs overload user working memory and generate psychological intimidation. According to Cognitive Load Theory, breaking complex information into manageable, cohesive chunks reduces task completion friction and error rates.

Forms must segment large input groups into semantic '<fieldset>' elements with explanatory '<legend>' titles, or utilize multi-step wizards ('<Stepper>', progressive tabs). Furthermore, each individual chunk must not exceed 7 active inputs to preserve perceptual focus.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Severe Form Abandonment** | MEDIUM | Users faced with tall, endless walls of inputs drop off significantly before completing registration or checkout flows. |
| **Field Omission & Data Entry Fatigue** | MEDIUM | Dense unchunked inputs cause users to overlook required fields, leading to repeated validation errors upon submission. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Monolithic form with 10 flat interactive inputs without fieldset or step chunking):
```tsx
<form onSubmit={handleSubmit} className="space-y-4 max-w-md">
  <input name="f1" placeholder="First Name" />
  <input name="f2" placeholder="Last Name" />
  <input name="f3" placeholder="Email" />
  <input name="f4" placeholder="Phone" />
  <input name="f5" placeholder="Address" />
  <input name="f6" placeholder="City" />
  <input name="f7" placeholder="State" />
  <input name="f8" placeholder="Zip Code" />
  <input name="f9" placeholder="Company" />
  <input name="f10" placeholder="Job Title" />
  <button type="submit">Kirim</button>
</form>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Segmented into two semantic fieldsets with legends, keeping fields per chunk under 7):
```tsx
<form onSubmit={handleSubmit} className="space-y-6 max-w-md">
  <fieldset className="space-y-4 border p-4 rounded-lg">
    <legend className="font-semibold text-sm">Informasi Pribadi</legend>
    <input name="f1" placeholder="First Name" />
    <input name="f2" placeholder="Last Name" />
    <input name="f3" placeholder="Email" />
    <input name="f4" placeholder="Phone" />
  </fieldset>

  <fieldset className="space-y-4 border p-4 rounded-lg">
    <legend className="font-semibold text-sm">Alamat Pengiriman</legend>
    <input name="f5" placeholder="Address" />
    <input name="f6" placeholder="City" />
    <input name="f7" placeholder="State" />
    <input name="f8" placeholder="Zip Code" />
  </fieldset>

  <button type="submit">Kirim</button>
</form>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: ux.monolithic-form-bloat"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore ux.monolithic-form-bloat` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/ux.monolithic-form-bloat/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for ux.monolithic-form-bloat"]
        subgraph P ["Positive Corpus (tests/correctness/ux.monolithic-form-bloat/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/ux.monolithic-form-bloat/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/ux.monolithic-form-bloat/adversarial/)"]
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
<!-- charites:ignore ux.monolithic-form-bloat intentional exception -->
```

```tsx
// charites:ignore ux.monolithic-form-bloat intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.monolithic-form-bloat:
    severity: warn # error | warn | info | off
```

