# performance.tailwind-untracked-package-source

> **Rule ID:** `performance.tailwind-untracked-package-source`
> **Severity:** `ERROR`
> **Category:** `performance`
> **Target Standards:** Tailwind CSS v4 Configuration Architecture (@source Directive Specification), Monorepo Multi-Package Style Discovery Standards, Oxide Engine Workspace Scanning Invariants

---

## 1. Overview & Core Invariant

Mewajibkan pendaftaran direktif @source pada berkas CSS root Tailwind v4 ketika mengimpor paket workspace monorepo eksternal.

### Core Invariant:
> **"Tailwind CSS v4 root stylesheets importing external monorepo packages must declare '@source' path directives; without '@source', the Oxide scanner skips external package directories, silently dropping all utility styles from compiled builds."**

---
## 2. Technical Grounding & Engine Realities

In Tailwind CSS v4, the legacy `tailwind.config.js` `content` array is replaced by CSS-first `@source` directives in the main stylesheet.

By default, Tailwind v4 only scans files in the immediate project directory. If the project imports components from external workspace packages (e.g. `@repo/ui` or `../../packages/...`), those package directories are ignored by default.

Failing to add `@source "../../packages/ui";` causes all Tailwind utility classes used inside those shared packages to be completely absent from the final CSS bundle.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Missing Monorepo Component Styles** | HIGH | Shared monorepo UI components render completely unstyled in production because utility classes inside them were never scanned. |
| **Silent Build Failures** | HIGH | No build errors are thrown; stylesheets simply compile without the required utility declarations. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### CSS (Berkas CSS root tidak menyertakan direktif @source untuk paket monorepo):
```css
/* Pelanggaran: Mengimpor tailwindcss tanpa @source untuk monorepo packages */
@import "tailwindcss";
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### CSS (Mendaftarkan path paket eksternal via @source):
```css
/* Patuh: Menyertakan direktif @source untuk paket monorepo */
@import "tailwindcss";
@source "../../packages/ui";
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore performance.tailwind-untracked-package-source intentional exception -->
```

```tsx
// charites:ignore performance.tailwind-untracked-package-source intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.tailwind-untracked-package-source:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [performance Category Guide](performance).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


