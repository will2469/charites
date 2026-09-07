# Implementation Plan: Flag Arbitrary Inline `scale-[...]` in `theme.hardcode-size` (Issue #3)

Resolves **[GitHub Issue #3](https://github.com/will2469/charites/issues/3)** by establishing an explicit, typed classification taxonomy in `theme.hardcode-size` to detect and flag arbitrary inline transform scale modifiers (`scale-[...]`, `-scale-[...]`, `scale-x-[...]`, `-scale-x-[...]`, `scale-y-[...]`, `-scale-y-[...]`, `scale-z-[...]`, `-scale-z-[...]`, `[scale:...]`), enforcing token-backed CSS variables (`scale-[var(--scale-press)]`, `[scale:var(--scale-press)]`) and standard Tailwind scale steps (`scale-95`, `scale-105`), while enforcing strict non-ownership boundaries (rejecting `rotate`, `duration`, `opacity`, `skew`, `translate`), robust numeric validation, and zero heap allocations on clean nodes (`QUAL-03`).

---

## 1. Ownership & Acceptance Invariants

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THE 5 NON-NEGOTIABLE ACCEPTANCE INVARIANTS               │
├─────────────────────────────────────────────────────────────────────────────┤
│ 1. Explicit Taxonomy Invariant:                                             │
│    Classification MUST determine ownership explicitly via SizeClassKind    │
│    (ArbitraryScale, NonStandardFraction, ArbitraryScalar, None).            │
│    ClassifySizeUtility() is the SSOT; IsHardcodedSizeUtility() is an adapter.│
│                                                                             │
│ 2. Precedence Invariant:                                                    │
│    Transform-scale classification has precedence over generic scalar        │
│    classification, ensuring scale diagnostics are never shadowed.           │
│                                                                             │
│ 3. Strict Numeric Validation Invariant:                                     │
│    Only syntactically valid numbers (e.g. 0.99, 1.02, .98, -0.95, +1.05)   │
│    and percentages (e.g. 98%, 105%) trigger scale violations. Corrupt forms │
│    (scale-[abc], scale-[calc(1+2)]) MUST NOT trigger numeric scale finding. │
│                                                                             │
│ 4. Token-Backed Scale Exemption Invariant:                                  │
│    Any scale value containing a CSS custom-property reference (`var(--...)`)│
│    is classified as token-backed and is not reported by this rule.          │
│                                                                             │
│ 5. Strict Non-Ownership Boundary & QUAL-03 Zero Alloc Invariant:            │
│    theme.hardcode-size DOES NOT OWN rotate, duration, opacity, or skew.     │
│    All clean nodes, standard scale steps (scale-95), CSS-var scale, and     │
│    unowned arbitrary classes MUST report 0 B/op, 0 allocs/op.               │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Ownership Acceptance Matrix

| Input | `theme.hardcode-size` Status | Diagnostic Type |
| :--- | :---: | :--- |
| `scale-95`, `scale-105` | CLEAN | None |
| `scale-[0.99]`, `-scale-[0.99]` | VIOLATION | ArbitraryScale |
| `scale-x-[0.95]`, `-scale-x-[0.95]` | VIOLATION | ArbitraryScale |
| `scale-y-[1.05]`, `-scale-y-[1.05]` | VIOLATION | ArbitraryScale |
| `scale-z-[1.05]`, `-scale-z-[1.05]` | VIOLATION | ArbitraryScale |
| `[scale:0.98]`, `[scale-x:0.98]` | VIOLATION | ArbitraryScale |
| `scale-[98%]` | VIOLATION | ArbitraryScale |
| `scale-[var(--scale-press)]` | CLEAN | Token-backed custom property |
| `[scale:var(--scale-press)]` | CLEAN | Token-backed custom property |
| `scale-[calc(1+2)]` | CLEAN / NOT OWNED | Non-numeric scale (not owned) |
| `rotate-[45deg]` | CLEAN / NOT OWNED | Non-scale transform (not owned) |
| `skew-x-[10deg]` | CLEAN / NOT OWNED | Non-scale transform (not owned) |
| `duration-[300ms]` | CLEAN / NOT OWNED | Transition duration (not owned) |
| `opacity-[0.85]` | CLEAN / NOT OWNED | Color/opacity (not owned) |
| `p-[19px]` | VIOLATION | ArbitraryScalar |
| `p-3.25` | VIOLATION | NonStandardFraction |

---

## 3. Classification Pipeline

```text
raw class
   ↓
StripVariantsOnlyBase
   ↓
ClassifySizeUtility
   ↓
SizeClassKind
   ├── ArbitraryScale
   │    ├── contains "var(--"  => CLEAN
   │    ├── isValidScaleValue  => VIOLATION (Arbitrary scale diagnostic)
   │    └── otherwise          => NOT OWNED / None
   ├── NonStandardFraction    => VIOLATION (Fractional scale diagnostic)
   ├── ArbitraryScalar        => VIOLATION (Spatial scalar diagnostic)
   └── None                   => CLEAN
```

Diagnostic specification for `SizeClassArbitraryScale`:
- **Message:** `"Arbitrary inline scale modifier: \"" + class + "\""`
- **Hint:** `"Avoid arbitrary inline scale modifiers. Define a standardized scale token/variable in global.css (e.g. --scale-press) or use standard Tailwind scale steps (scale-95, scale-105)."`

---

## 4. Verification & Hardening

1. **Table-Driven Unit Tests:** Full directional and signed matrix in `internal/rules/theme/hardcode_size_test.go`.
2. **Hardened QUAL-03 Benchmark:** Clean node containing standard scale, CSS-var scale, and unowned non-spatial arbitrary properties verifying `0 B/op, 0 allocs/op`.
3. **Tri-Corpus Golden Correctness Matrix:** `tests/correctness/theme/hardcode-size/` with positive scale violations, negative standard scale, and adversarial token-backed variables / non-spatial utilities.
4. **Ergonomy Follow-Up Compatibility Cleanup:** Align `internal/rules/ergonomy/` recommendation from `active:scale-[0.99]` to `active:scale-95`.
5. **Full Repository Gate:** `make all` clean.
