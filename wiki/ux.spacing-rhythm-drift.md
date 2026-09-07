# ux.spacing-rhythm-drift

> **Rule ID:** `ux.spacing-rhythm-drift`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Gestalt Law of Proximity & Visual Rhythm, W3C Design Tokens Community Group (DTCG v2025.10 - Spatial Cadence), Evidence-Based Static Analysis for Design Systems

---

## 1. Overview & Core Invariant

Detects spacing sequences that form an inconsistent local rhythm within the same layout group

### Core Invariant:
> **"Spacing values serving the same spatial role within the same layout group should normally come from a coherent local spacing family. When a dominant local rhythm is established (>= 3 occurrences, dominance ratio >= 60%), isolated unexplained outliers are reported as warning."**

---
## 2. Technical Grounding & Engine Realities

In modern user interface engineering, consistent spacing intervals between peer components establish visual cadence, cognitive grouping, and predictable layout rhythm.

When a layout container contains repeated sibling relationships with homogeneous spacing, but one relationship unexpectedly deviates (e.g. mb-4, mb-4, mb-7, mb-4), users perceive visual jarring and broken alignment.

Rather than enforcing arbitrary universal moduli (e.g. strict 4px or 8px grid), ux.spacing-rhythm-drift evaluates local relational consistency within verified layout groups, reporting isolated unexplained outliers while remaining immune to intentional responsive variants, distinct semantic roles, and small sample sizes.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Visual Cadence Disruption** | MEDIUM | Irregular spacing among peer items erodes user trust and creates visual jarring without intentional hierarchy. |
| **Cognitive Grouping Confusion** | LOW | Users misperceive arbitrary spacing outliers as deliberate semantic distinctions or broken layout states. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Repeated form fields use uniform mb-4 spacing, but one field introduces an unexplained mb-7 outlier):
```tsx
<div>
  <Field className="mb-4" />
  <Field className="mb-4" />
  <Field className="mb-7" />
  <Field className="mb-4" />
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Repeated form fields adhere strictly to a coherent local cadence with homogeneous mb-4 spacing):
```tsx
<div>
  <Field className="mb-4" />
  <Field className="mb-4" />
  <Field className="mb-4" />
  <Field className="mb-4" />
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.spacing-rhythm-drift intentional exception -->
```

```tsx
// charites:ignore ux.spacing-rhythm-drift intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.spacing-rhythm-drift:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


