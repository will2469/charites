# ux.empty-collection-unhandled

> **Rule ID:** `ux.empty-collection-unhandled`
> **Severity:** `INFO`
> **Category:** `ux`
> **Target Standards:** Zero-State Usability & Mental Model Continuity (Nielsen Norman Group), Feedforward Principle & Gulf of Evaluation (Don Norman), ISO 9241-110 Ergonomics of Human-System Interaction (Suitability for Learning)

---

## 1. Overview & Core Invariant

Advises handling empty collection state when mapping dynamic items to avoid zero-state blindness

### Core Invariant:
> **"Dynamic collection rendering expressions must handle empty collection states ('collection.length === 0') with informative fallback zero-state UI."**

---
## 2. Technical Grounding & Engine Realities

When dynamic lists, tables, or feed collections contain 0 records and render nothing, users are stranded in an ambiguous vacuum: did the request fail, is it still loading, or are there genuinely zero records?

Zero-state blindness forces users to refresh repeatedly or assume the application is broken. A dedicated empty state component (e.g. '<EmptyState />' with an illustration, clarifying text, and a call-to-action button) confirms system status and proactively guides user next steps.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Zero-State Blindness & System Status Ambiguity** | LOW | Users perceive blank empty containers as silent application crashes or perpetual loading freezes. |
| **Workflow Dead Ends** | LOW | Without an actionable empty state CTA (e.g. 'Create First Invoice'), users cannot self-discover how to populate the collection. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Dynamic list rendering items via .map() without handling empty array state):
```tsx
<div className="space-y-3">
  <h2 className="text-lg font-bold">Daftar Tagihan</h2>
  <List items={invoices.map(inv => <InvoiceRow key={inv.id} data={inv} />)} />
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Explicit empty state fallback branch when collection has 0 items):
```tsx
<div className="space-y-3">
  <h2 className="text-lg font-bold">Daftar Tagihan</h2>
  {invoices.length === 0 ? (
    <EmptyState
      title="Belum Ada Tagihan"
      description="Buat tagihan pertama Anda untuk mulai menerima pembayaran."
      actionText="Buat Tagihan"
    />
  ) : (
    <List items={invoices.map(inv => <InvoiceRow key={inv.id} data={inv} />)} />
  )}
</div>
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
    IgnoreCheck -- "Not Ignored" --> Diag["4. Emit Diagnostic: ux.empty-collection-unhandled"]
```

### Step-by-Step Evaluation:
1. **AST Traversal:** Visits normalized `ir.Node` structures during AST streaming.
2. **Invariant Check:** Compares node attributes against rule invariants.
3. **Directive Suppression Check:** Honors inline `charites:ignore ux.empty-collection-unhandled` directives.
4. **Diagnostic Emission:** Emits compiler-grade diagnostic on invariant violation.

---

## 7. Verification & Test Harness (How The Test Works: 1-SSOT Tri-Corpus)

This rule is rigorously tested and validated across the canonical **1-SSOT Tri-Corpus** in `tests/correctness/ux.empty-collection-unhandled/`:

```mermaid
flowchart TD
    subgraph GoldenCorpus ["1-SSOT Tri-Corpus Test Matrix for ux.empty-collection-unhandled"]
        subgraph P ["Positive Corpus (tests/correctness/ux.empty-collection-unhandled/positive/)"]
            P1["P1: Obvious Direct Violation"]
            P2["P2: Indirect / Variant Concatenation"]
            P3["P3: Helper / clsx / cn Wrapper"]
            P4["P4: Deeply Nested Elements"]
            P5["P5: Aliased Imports / Re-exports"]
        end
        subgraph N ["Negative Corpus (tests/correctness/ux.empty-collection-unhandled/negative/)"]
            N1["N1: Valid Design Tokens"]
            N2["N2: Explicit charites:ignore Directive"]
            N3["N3: Third-Party / Vendor Components"]
            N4["N4: Clean Semantic HTML"]
            N5["N5: Untokenized Custom Values (Banana Test)"]
        end
        subgraph A ["Adversarial Corpus (tests/correctness/ux.empty-collection-unhandled/adversarial/)"]
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
<!-- charites:ignore ux.empty-collection-unhandled intentional exception -->
```

```tsx
// charites:ignore ux.empty-collection-unhandled intentional exception
```

---

## 9. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.empty-collection-unhandled:
    severity: info # error | warn | info | off
```

