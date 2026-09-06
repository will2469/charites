# ergonomy.missing-inputmode-keyboard

> **Rule ID:** `ergonomy.missing-inputmode-keyboard`
> **Severity:** `INFO`
> **Category:** `ergonomy`
> **Target Standards:** HTML Living Standard Section 4.10.5.3 (The inputmode attribute), Tesler's Law (Conservation of Complexity in Virtual Keyboards), Apple iOS & Android Mobile Virtual Keyboard Guidelines

---

## 1. Overview & Core Invariant

Enforces contextual virtual keyboard inputmode and type attributes on mobile form inputs (Tesler's Law)

### Core Invariant:
> **"Form text inputs collecting numeric, phone, or email values must declare contextual 'inputmode' or specialized 'type' attributes to directly open the optimized mobile virtual keypad."**

---
## 2. Technical Grounding & Engine Realities

On mobile devices, focusing an input without specialized type or inputmode opens the full standard QWERTY keyboard.

For numeric, telephone, or OTP fields, this forces the user to repeatedly toggle keyboard layers to find digits. According to Tesler's Law, complexity must be absorbed by software rather than offloaded to user manual effort. Declaring 'inputmode="numeric"' or 'type="tel"' instantly summons large, thumb-friendly numeric keypads.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mobile Cognitive Friction** | LOW | Users are forced to manually switch keyboard layers on small touchscreens to enter digits or phone numbers. |
| **High Form Abandonment & Typing Errors** | LOW | Entering OTP or financial amounts on dense QWERTY keys leads to frequent miss-taps and delayed checkout flows. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Phone input missing type or inputMode):
```tsx
<input
  name="nomor_hp"
  placeholder="08123456789"
  className="h-11 px-3.5 py-2.5 border rounded-xl"
/>
```
### ASTRO (OTP field defaulting to QWERTY keyboard):
```astro
<input
  id="otp_code"
  placeholder="123456"
  class="h-11 px-3.5 py-2.5 border rounded-xl"
/>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Explicit telephone keypad and autocomplete):
```tsx
<input
  name="nomor_hp"
  type="tel"
  inputMode="tel"
  autoComplete="tel"
  placeholder="08123456789"
  className="h-11 px-3.5 py-2.5 border rounded-xl"
/>
```
### ASTRO (Numeric keypad for OTP verification):
```astro
<input
  id="otp_code"
  type="text"
  inputmode="numeric"
  pattern="[0-9]*"
  placeholder="123456"
  class="h-11 px-3.5 py-2.5 border rounded-xl"
/>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: ergonomy.missing-inputmode-keyboard"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore ergonomy.missing-inputmode-keyboard` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/ergonomy.missing-inputmode-keyboard/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for ergonomy.missing-inputmode-keyboard"]
        subgraph P ["Positive Corpus (tests/correctness/ergonomy.missing-inputmode-keyboard/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/ergonomy.missing-inputmode-keyboard/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/ergonomy.missing-inputmode-keyboard/adversarial/)"]
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
<!-- charites:ignore ergonomy.missing-inputmode-keyboard intentional exception -->
```

```tsx
// charites:ignore ergonomy.missing-inputmode-keyboard intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ergonomy.missing-inputmode-keyboard:
    severity: info # error | warn | info | off
```

