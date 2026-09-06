# performance.tailwind-duplicate-arbitrary-rules

> **Rule ID:** `performance.tailwind-duplicate-arbitrary-rules`
> **Severity:** `WARN`
> **Category:** `performance`
> **Target Standards:** Tailwind CSS v4 Design Tokens & Spacing Scale Standards, Compiled CSS Output Deduplication & Payload Economy, W3C Stylesheet Declarative Optimization Guidelines

---

## 1. Overview & Core Invariant

Menganjurkan penggunaan utilitas skala inti bawaan Tailwind v4 alih-alih nilai arbitrary sembarang yang menghasilkan deklarasi CSS duplikat.

### Core Invariant:
> **"Arbitrary value utilities (e.g. 'p-[16px]', 'mt-[1rem]') that match standard Tailwind core scale tokens should use the canonical core utility (e.g. 'p-4', 'mt-4') to avoid duplicate CSS rule generation in compiled bundles."**

---
## 2. Technical Grounding & Engine Realities

Tailwind CSS includes a refined, consistent default spacing and sizing scale.

When developers write ad-hoc arbitrary values like `p-[16px]` alongside `p-4` (which also resolves to `padding: 1rem / 16px`), Tailwind generates separate unique CSS selector rules for both.

Consolidating arbitrary values to their core scale equivalents eliminates redundant rule definitions, shrinks the production CSS footprint, and ensures consistent visual rhythm across the application.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Compiled CSS Bloat** | MEDIUM | Inflates stylesheet size with duplicate CSS selector blocks that declare identical CSS properties and values. |
| **Visual Rhythm Inconsistency** | LOW | Ad-hoc arbitrary values drift away from the design system's cohesive 4px/8px modular grid. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Menggunakan nilai arbitrary yang menduplikasi utilitas core p-4 dan mt-4):
```tsx
<div className="p-[16px] mt-[1rem]">Konten</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Menggunakan utilitas skala core standar):
```tsx
<div className="p-4 mt-4">Konten</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore performance.tailwind-duplicate-arbitrary-rules intentional exception -->
```

```tsx
// charites:ignore performance.tailwind-duplicate-arbitrary-rules intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.tailwind-duplicate-arbitrary-rules:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [performance Category Guide](performance).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


