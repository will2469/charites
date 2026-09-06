# performance.react-index-as-key

> **Rule ID:** `performance.react-index-as-key`
> **Severity:** `WARN`
> **Category:** `performance`
> **Target Standards:** React Official Documentation (Lists and Keys Invariants), React Reconciliation Algorithm Specification (Diffing with stable keys), Robin Pokorny Guidelines ('Index as a key is an anti-pattern')

---

## 1. Overview & Core Invariant

Using array index as 'key' in dynamic collection mapping breaks VDOM reconciliation when items reorder or mutate

### Core Invariant:
> **"Dynamic collections mapped with '.map()' must use stable, unique item identifiers (e.g. 'item.id') as the 'key' attribute rather than numeric array indexes."**

---
## 2. Technical Grounding & Engine Realities

React relies on the `key` attribute to identify which items in a list have changed, been added, or been removed during reconciliation.

When an array index is used (`key={index}`), rearranging, filtering, prepending, or deleting items shifts the indexes of subsequent elements.

This index drift causes React to confuse element identities, erroneously preserving local uncontrolled state (such as form inputs, focus, and CSS transitions) on the wrong items and forcing redundant re-renders of the entire subtree.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Component State Desynchronization** | HIGH | Internal component state (e.g. input values, selection states) remains bound to the array position rather than the underlying data entity. |
| **Unnecessary Subtree DOM Repaints** | MEDIUM | React fails to recognize moved nodes and completely remounts DOM subtrees instead of performing lightweight repositioning. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Dynamic transactions list using array index as key):
```tsx
{transactions.map((tx, index) => (
  <TransactionRow key={index} data={tx} />
))}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Persistent unique entity identifier used as reconciliation key):
```tsx
{transactions.map((tx) => (
  <TransactionRow key={tx.id} data={tx} />
))}
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: performance.react-index-as-key"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore performance.react-index-as-key` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/performance.react-index-as-key/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for performance.react-index-as-key"]
        subgraph P ["Positive Corpus (tests/correctness/performance.react-index-as-key/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/performance.react-index-as-key/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/performance.react-index-as-key/adversarial/)"]
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
<!-- charites:ignore performance.react-index-as-key intentional exception -->
```

```tsx
// charites:ignore performance.react-index-as-key intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.react-index-as-key:
    severity: warn # error | warn | info | off
```

