# performance.tailwind-dynamic-class-concatenation

> **Rule ID:** `performance.tailwind-dynamic-class-concatenation`
> **Severity:** `ERROR`
> **Category:** `performance`
> **Target Standards:** Tailwind CSS v4 Compiler Scanner Specification (Oxide Static Extraction), Tailwind CSS Official Architecture ('Dynamic Class Names Limitations'), Zero-Runtime CSS Extraction Invariants

---

## 1. Overview & Core Invariant

Mencegah penggabungan string nama kelas dinamis parsial yang merusak deteksi compiler scanner Tailwind CSS v4 (Oxide engine).

### Core Invariant:
> **"Tailwind CSS utility classes must be written as complete, static string literals; dynamic string interpolation on partial class prefixes prevents the static scanner from detecting classes, resulting in missing styles in production."**

---
## 2. Technical Grounding & Engine Realities

Tailwind CSS v4 uses a high-performance static scanner (Oxide engine) that scans source code for complete class tokens without executing JavaScript runtime.

Constructing utility names dynamically via template literals or string concatenation (e.g. `bg-${color}-100` or `'text-' + size`) breaks static extraction completely.

Because the scanner never evaluates runtime variables, it never sees the complete utility string (like `bg-red-100`), causing the required CSS rules to be omitted from the compiled stylesheet.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Missing Production Stylesheet Rules** | HIGH | Utility classes generated through string concatenation are completely missing from the production CSS bundle, causing broken UI visuals. |
| **Silent Runtime Failures** | HIGH | Classes appear functional in local environments if the class was previously cached or generated elsewhere, but fail silently upon clean production builds. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Penggabungan string parsial tidak dapat diekstrak oleh compiler):
```tsx
function Badge({ color }: { color: 'red' | 'blue' }) {
  return <span className={`bg-${color}-100 text-${color}-800`}>Status</span>;
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Menuliskan nama kelas secara utuh dalam kamus statis):
```tsx
const COLOR_MAP = {
  red: 'bg-red-100 text-red-800',
  blue: 'bg-blue-100 text-blue-800',
} as const;

function Badge({ color }: { color: 'red' | 'blue' }) {
  return <span className={COLOR_MAP[color]}>Status</span>;
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore performance.tailwind-dynamic-class-concatenation intentional exception -->
```

```tsx
// charites:ignore performance.tailwind-dynamic-class-concatenation intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.tailwind-dynamic-class-concatenation:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [performance Category Guide](performance).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


