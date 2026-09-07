# Implementation Plan: Flexbox Fractional Width & Gap Spacing Geometry (`responsive.fractional-width-gap-drift`)

Resolves **[GitHub Issue #8](https://github.com/will2469/charites/issues/8)** (`responsive.fractional-width-gap-drift`: `[enhancement] feat(responsive): detect mobile baseline fractional widths and flexbox gap drift`).

---

## 1. Executive Summary & Problem Formulation

In web layout development with Tailwind CSS, percentage and fractional width utilities (`w-1/2`, `w-2/3`, `w-3/4`, `w-[60%]`, etc.) interact with CSS Flexbox spacing geometry:

```tsx
// Anti-Pattern 1: Unexpected Line Wrapping
<div className="flex flex-wrap gap-4">
  <div className="w-1/2">A (50%)</div>
  <div className="w-1/2">B (50%) - Patah ke baris baru karena gap-4 (16px)</div>
</div>

// Anti-Pattern 2: Rigid Shrink-Disabled Overflow Blowout
<div className="flex gap-4">
  <div className="w-1/2 shrink-0">A (50% rigid)</div>
  <div className="w-1/2 shrink-0">B (50% rigid) - Overflow container 16px</div>
</div>
```

### The CSS Specification & Mathematical Proof:
Under the W3C CSS Flexible Box Layout specification:
1. **Negative Free Space & Native Shrinking:** When a flex container has `gap > 0` and fractional children totaling 100%, the total row width is $100\% + \text{gap} > 100\%$. However, under **default `flex-wrap: nowrap` and default `flex-shrink: 1`**, Flexbox automatically shrinks each item by its proportional share of the negative free space. Therefore:
   $$\text{fractional sum} \ge 100\% + \text{gap} \ne \text{defect by itself}$$
   The browser can absorb the gap surplus cleanly through native shrinking.
2. **When it becomes a Structurally Non-Resolvable Defect:**
   A true structural defect occurs **only** when the horizontal free-space conflict cannot be resolved through normal shrinking:
   - **Case A (`flex-wrap` / `flex-wrap-reverse` active):** Under Flexbox line-breaking (Section 9.3), $50\% + 16\text{px} + 50\% > 100\%$, forcing the second item to wrap to a new line, destroying the intended multi-column layout.
   - **Case B (Shrink disabled via `shrink-0` / `flex-none` on fractional items):** Items refuse to shrink, forcing the row to overflow the container by the gap width.
3. **Standalone `w-*` is NOT a Defect:** Standalone `<div className="w-1/2">` or `<div className="w-2/3">` is not an automatic violation. AST analysis only warns when container constraints prove that width allocations conflict with spacing and wrap/shrink geometry.

---

## 2. The 4 Locked Architectural Invariants

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THE 4 LOCKED ARCHITECTURAL INVARIANTS                    │
├─────────────────────────────────────────────────────────────────────────────┤
│ 1. Non-Resolvable Structural Proof Invariant:                               │
│    Drift is flagged ONLY when:                                              │
│    - Container is Flex in horizontal row direction at that tier             │
│    - Container has horizontal gap > 0                                       │
│    - Children's fractional widths sum >= 100% (count >= 2)                  │
│    - AND horizontal free-space conflict is structurally non-resolvable:     │
│        (a) wrap is enabled (flex-wrap or flex-wrap-reverse), OR             │
│        (b) at least one fractional child has shrink disabled (shrink-0/     │
│            flex-none).                                                      │
│    Default nowrap with default shrink-1 is SAFE (zero false alarms).        │
│                                                                             │
│ 2. Per-Tier DisplayMode & Breakpoint Cascade:                               │
│    DisplayMode (Flex, Block, Grid, None) is tier-aware. Overrides like      │
│    'md:block' or 'md:grid' deactivate Flex state at 'md:'. Container        │
│    direction, wrap, gap, and child shrink/width cascade independently.      │
│                                                                             │
│ 3. Intent-Driven Carousel Heuristic (No Blanket Suppression):               │
│    Carousel intent requires affirmative evidence:                           │
│    overflow-x-auto/scroll + horizontal flex + >=2 repeated fractional items │
│    + shrink-0. Does NOT provide blanket immunity to arbitrary containers.   │
│                                                                             │
│ 4. Pure Computed Geometry Architecture:                                     │
│    Clean separation between AST fact extraction (ResolvedGeometry) and      │
│    rule classification. Container owns exactly 1 diagnostic. Zero mutable   │
│    global state. Orthogonal to responsive.flex-child-overflow.              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Verification Results

### Unit Tests
- `internal/rules/responsive/geometry_test.go`: 100% pass across fraction parsing, arbitrary percentages, gap precedence (`gap-4 gap-x-0` order-insensitive), display resets (`flex md:block`), and classifier proofs.
- `internal/rules/responsive/fractional_width_gap_drift_test.go`: 25 test cases covering all 17 acceptance matrix scenarios and edge cases.
- `internal/rules/responsive/contract_test.go`: Passes `TestResponsiveRules_CanonicalContract/responsive.fractional-width-gap-drift`.

### 1-SSOT Tri-Corpus Correctness
- `tests/correctness/responsive/fractional-width-gap-drift/positive/`:
  - `flex_wrap_half_half.tsx`
  - `flex_shrink0_half_half.tsx`
  - `flex_wrap_two_thirds_one_third.astro`
  - `flex_wrap_breakpoint_drift.tsx`
  - `flex_wrap_arbitrary_percentage.tsx`
- `tests/correctness/responsive/fractional-width-gap-drift/negative/`:
  - `default_nowrap_half_half.tsx` (safe native shrinking)
  - `standalone_fractional_width.tsx` (no container conflict)
  - `css_grid_alternative.tsx`
  - `flex_one_fluid.tsx`
  - `carousel_slider_peek.tsx`
- `tests/correctness/responsive/fractional-width-gap-drift/adversarial/`:
  - `gap_zero_no_drift.tsx`
  - `gap_y_only_no_drift.tsx`
  - `gap_x_zero_override.tsx`
  - `fractional_sum_under_100.tsx`
  - `non_fractional_shrink0.tsx`
