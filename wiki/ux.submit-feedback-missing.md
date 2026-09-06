# ux.submit-feedback-missing

> **Rule ID:** `ux.submit-feedback-missing`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Doherty Threshold (< 400ms Interaction Feedback), Nielsen Heuristic #1: Visibility of System Status, ISO 9241-110 Ergonomics of Human-System Interaction (Self-Descriptiveness)

---

## 1. Overview & Core Invariant

Enforces reentry guard (disabled) and perceivable feedback (aria-busy/spinner) on async mutation triggers

### Core Invariant:
> **"Async mutation triggers must provide both Reentry Guard (R1: 'disabled={isPending}') and Perceivable Feedback (R2: 'aria-busy' or spinner)."**

---
## 2. Technical Grounding & Engine Realities

When users trigger a mutation (such as submitting an order or charging a payment) without immediate feedback, they perceive the system as unresponsive within the 400ms Doherty Threshold.

In the absence of a reentry lock, users repeatedly click the submit button, triggering duplicate requests, double financial charges, and server-side race conditions. Enforcing both R1 (reentry lockout) and R2 (visual feedback) guarantees transactional safety and cognitive assurance.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Duplicate Submissions & Double Charging** | HIGH | Rapid repeated clicks by uncertain users cause duplicate database inserts and multiple payment gateway captures. |
| **Cognitive Disorientation & Interaction Rage** | MEDIUM | Users doubt whether their input registered, leading to rage clicks and workflow cancellation. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Async submit button without disabling control or showing pending feedback):
```tsx
<button
  onClick={async () => {
    await api.post("/orders", orderData);
  }}
  className="bg-primary text-white px-4 py-2"
>
  Bayar Sekarang
</button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Compliant button satisfying R1 (disabled={isPending}) and R2 (aria-busy & dynamic text)):
```tsx
<button
  onClick={handlePayment}
  disabled={isPending}
  aria-busy={isPending}
  className="bg-primary text-white px-4 py-2 disabled:opacity-50"
>
  {isPending ? "Memproses Pembayaran..." : "Bayar Sekarang"}
</button>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: ux.submit-feedback-missing"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore ux.submit-feedback-missing` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/ux.submit-feedback-missing/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for ux.submit-feedback-missing"]
        subgraph P ["Positive Corpus (tests/correctness/ux.submit-feedback-missing/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/ux.submit-feedback-missing/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/ux.submit-feedback-missing/adversarial/)"]
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
<!-- charites:ignore ux.submit-feedback-missing intentional exception -->
```

```tsx
// charites:ignore ux.submit-feedback-missing intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.submit-feedback-missing:
    severity: warn # error | warn | info | off
```

