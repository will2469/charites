# Implementation Plan: Number Input Ergonomics & Semantic Integrity (Issue #5)

Resolves **[GitHub Issue #5](https://github.com/will2469/charites/issues/5)** by establishing a domain-driven `internal/rules/form` adapter and implementing three distinct, non-overlapping static analysis rules across the `ergonomy` and `ux` categories:
1. `ergonomy.number-input-wheel-hazard`: Prevents accidental value mutation caused by mouse wheel scrolling when hovering over `<input type="number">`.
2. `ux.number-input-identity-misuse`: Flags the misuse of `<input type="number">` on discrete identifier attributes (e.g. NIK, postal code, phone, credit card, PIN, account number, order ID), enforcing `inputmode="numeric"` with `type="text"`.
3. `ux.number-input-missing-bounds`: Flags unbounded `<input type="number">` elements representing quantitative measures (e.g. quantity, amount, discount, percentage, age) that lack `min` and/or `max` constraints.

---

## 1. Architectural Principles & The Shared Form Adapter

To prevent duplication and semantic drift across rules, all form-related AST extraction is centralized into `internal/rules/form`:

```text
                                AST / IR Node
                                      │
                                      ▼
                      internal/rules/form.ExtractFormFacts()
                                      │
               ┌──────────────────────┼──────────────────────┐
               ▼                      ▼                      ▼
        Input Properties        Event Handlers        Field Semantics
        - Tag & Type            - onWheel / onwheel    - ClassifyIdentifier()
        - Disabled / ReadOnly   - onWheelCapture       - Strong Identity
        - Min / Max Bounds      - blur / preventDefault- Contextual Quantity
               │                      │                      │
               ▼                      ▼                      ▼
    ergonomy.wheel-hazard      ux.identity-misuse     ux.missing-bounds
    (scroll hijacking risk)   (discrete identifiers) (unbounded quantities)
```

### Precedence Hierarchy for Identifier Classification:
1. **Strong Identity** (`KindIdentity`):
   Tokens that unambiguously denote discrete identification numbers (e.g. `nik`, `ktp`, `phone`, `tel`, `mobile`, `postal`, `zip`, `pin`, `otp`, `card`, `account`, `order_id`, `invoice_id`).
2. **Contextual Identity + Quantity Resolution**:
   Compound identifiers like `jumlah_nomor_kk` or `total_account_number` resolve to Identity because the identity anchor is primary.
3. **Quantity** (`KindQuantity`):
   Tokens that unambiguously represent measurable scalar amounts (e.g. `qty`, `quantity`, `count`, `amount`, `nominal`, `discount`, `diskon`, `price`, `harga`, `age`, `umur`, `weight`, `height`).
4. **Unknown** (`KindUnknown`):
   Generic or unclassified inputs.

---

## 2. Rule Specifications & Ownership Boundaries

### Rule 1: `ergonomy.number-input-wheel-hazard`
- **Category:** `ergonomy`
- **Severity:** `WARN`
- **Core Invariant:** Unprotected `<input type="number">` elements must prevent accidental value increments/decrements during vertical page scroll.
- **Suppression Criteria:**
  - Presence of `onWheel`, `onwheel`, or `onWheelCapture` handler invoking `blur()` or `preventDefault()`.
  - Element is `disabled` or `readOnly`.
  - Dynamic type bindings (e.g. `{type}`) are treated with conservative immunity.

### Rule 2: `ux.number-input-identity-misuse`
- **Category:** `ux`
- **Severity:** `ERROR`
- **Core Invariant:** Discrete identifier fields must use `type="text"` (or `type="tel"`) with `inputmode="numeric"` and `pattern="[0-9]*"`, never `type="number"`.
- **Rationale:** `<input type="number">` strips leading zeros (turning `0812` into `812`), alters large integers via floating-point rounding ($> 2^{53}-1$), allows scientific notation (`1e5`), and exposes stepper arrows.

### Rule 3: `ux.number-input-missing-bounds`
- **Category:** `ux`
- **Severity:** `WARN`
- **Core Invariant:** Quantitative `<input type="number">` elements must declare explicit upper (`max`) and/or lower (`min`) bounds.
- **Strict Yield Hierarchy:** If an input is classified as `KindIdentity`, this rule strictly yields to `ux.number-input-identity-misuse` (emits 0 diagnostics), preventing double-penalizing developers for an element that should not be a number input in the first place.

---

## 3. Verification & 1-SSOT Golden Corpus

Each rule is verified through Layer 1 (Structural Presence) and Layer 2 (Metric Gate: `Positive > 0 && Negative == 0 && Adversarial == 0`) across Astro and React TSX fixtures:
- `tests/correctness/ergonomy/number-input-wheel-hazard/`
- `tests/correctness/ux/number-input-identity-misuse/`
- `tests/correctness/ux/number-input-missing-bounds/`
- `tests/correctness/cross_rule_ownership_test.go`: Explicitly asserts cross-rule yielding and precedence contracts.
