# performance.tailwind-duplicate-utility-definition

> **Rule ID:** `performance.tailwind-duplicate-utility-definition`
> **Severity:** `WARN`
> **Category:** `performance`
> **Target Standards:** Tailwind CSS v4 '@utility' Directive Specification, Compiled CSS Output Deduplication & Bundle Hygiene, Atomic CSS Design Invariants

---

## 1. Overview & Core Invariant

Mencegah duplikasi deklarasi utilitas CSS kustom (@utility) yang properti dan nilainya sudah disediakan oleh utilitas core bawaan Tailwind CSS v4.

### Core Invariant:
> **"Custom '@utility' declarations must not duplicate built-in Tailwind CSS core utilities; redundant definitions generate unnecessary stylesheet bytes and break atomic CSS composability."**

---
## 2. Technical Grounding & Engine Realities

The `@utility` directive in Tailwind CSS v4 is designed to register brand-new utilities for modern or proprietary CSS features not yet included in core Tailwind.

Defining custom `@utility` blocks for combinations already covered by core utilities (such as `@utility center-flex { display: flex; align-items: center; }`) produces duplicate CSS rules in the compiled stylesheet.

Composing canonical core utilities (`flex items-center`) directly in markup preserves atomic stylesheet economy and avoids redundant selector bloat.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Compiled Stylesheet Bloat** | MEDIUM | Adds unnecessary custom CSS rules to the production build that duplicate pre-existing atomic utilities. |
| **Bypassed Utility Composability** | LOW | Custom wrapper utilities fracture atomic consistency and make component classes harder to refactor. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### CSS (Mendefinisikan @utility yang menduplikasi utilitas core flexbox):
```css
@utility center-flex {
  display: flex;
  align-items: center;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Menggunakan kombinasi utilitas core native langsung di markup):
```tsx
<div className="flex items-center">Konten</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore performance.tailwind-duplicate-utility-definition intentional exception -->
```

```tsx
// charites:ignore performance.tailwind-duplicate-utility-definition intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.tailwind-duplicate-utility-definition:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [performance Category Guide](performance).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


