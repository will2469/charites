# ux.competing-primary-cta

> **Rule ID:** `ux.competing-primary-cta`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Von Restorff Effect (The Isolation Effect / Visual Dominance), Hick-Hyman Law (Logarithmic Decision Latency), Nielsen Norman Group (Visual Hierarchy for Action Buttons)

---

## 1. Overview & Core Invariant

Warns when an action group or interactive container contains more than one primary call-to-action button

### Core Invariant:
> **"An action container or button group must contain at most one primary call-to-action button, ensuring a clear visual focal point and zero decision ambiguity."**

---
## 2. Technical Grounding & Engine Realities

The Von Restorff Effect (Isolation Effect) predicts that when multiple similar items are presented, the one that differs from the rest is most likely to be remembered and acted upon.

When an interface presents two or more buttons styled identically as primary actions (e.g., two 'bg-primary' or 'variant="primary"' buttons side by side in a modal footer or form actions), visual hierarchy collapses.

This competing prominence causes choice paralysis (Hick-Hyman Law), forces users to pause and re-read labels carefully, and drastically increases the probability of accidental mis-clicks. Every decision context must have exactly one visually distinct primary action, while supporting actions should be styled with 'outline', 'secondary', or 'ghost' variants.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Choice Paralysis & Decision Latency** | HIGH | Users hesitate when confronted with equal-weight visual cues, increasing conversion drop-off rates. |
| **Accidental Action Slips** | MEDIUM | Users mistake secondary or cancel triggers for primary confirmation due to identical color and elevation styling. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Dialog footer with two competing primary buttons creating visual ambiguity):
```tsx
<div className="flex justify-end gap-3 p-4">
  <button type="button" className="bg-primary text-white px-4 py-2 rounded-md">Simpan Draf</button>
  <button type="submit" className="bg-primary text-white px-4 py-2 rounded-md">Publikasikan</button>
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Clear hierarchy: one primary button paired with a secondary outline button):
```tsx
<div className="flex justify-end gap-3 p-4">
  <button type="button" className="border border-input bg-transparent px-4 py-2 rounded-md">Simpan Draf</button>
  <button type="submit" className="bg-primary text-white px-4 py-2 rounded-md">Publikasikan</button>
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.competing-primary-cta intentional exception -->
```

```tsx
// charites:ignore ux.competing-primary-cta intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.competing-primary-cta:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


