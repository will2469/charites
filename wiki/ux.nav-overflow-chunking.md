# ux.nav-overflow-chunking

> **Rule ID:** `ux.nav-overflow-chunking`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Miller's Law (Information Processing Capacity: 7 ± 2 Chunks), Information Architecture Chunking & Category Hierarchy (Rosenfeld & Morville), W3C WAI-ARIA Authoring Practices Guide 1.2 (Navigation Menubars)

---

## 1. Overview & Core Invariant

Warns when a navigation landmark contains more than 7 direct navigation links without chunking mechanisms

### Core Invariant:
> **"Navigation landmarks ('<nav>' or 'role="navigation"') must not present more than 7 flat direct links without grouping into disclosures, dropdown menus, or category drawers."**

---
## 2. Technical Grounding & Engine Realities

Miller's Law dictates that human working memory can reliably retain only 7 ± 2 distinct chunks of information at any single time.

When a main navigation bar presents 8 or more flat links in a single row or list without visual or hierarchical chunking, users experience choice paralysis and elevated cognitive scan latency.

To maintain optimal information architecture, high-density menus should group secondary destinations into nested dropdowns, accordions, or an overflow 'More...' disclosure container.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cognitive Overload & Choice Paralysis** | MEDIUM | Users take significantly longer to locate key navigation targets and frequently miss secondary features. |
| **Visual Clutter on Narrow Viewports** | MEDIUM | Flat multi-link navigation rows wrap awkwardly or cause accidental taps on mobile touchscreens. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Ten flat navigation links inside <nav> without grouping or overflow chunking):
```tsx
<nav className="flex gap-4">
  <a href="/home">Beranda</a>
  <a href="/profil">Profil</a>
  <a href="/layanan">Layanan</a>
  <a href="/berita">Berita</a>
  <a href="/transparansi">Transparansi</a>
  <a href="/anggaran">Anggaran</a>
  <a href="/regulasi">Regulasi</a>
  <a href="/galeri">Galeri</a>
  <a href="/kontak">Kontak</a>
  <a href="/bantuan">Bantuan</a>
</nav>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Primary destinations kept to 4 links, with remaining links grouped into a DropdownMenu):
```tsx
<nav className="flex gap-4 items-center">
  <a href="/home">Beranda</a>
  <a href="/profil">Profil</a>
  <a href="/layanan">Layanan</a>
  <a href="/berita">Berita</a>
  <DropdownMenu>
    <button type="button" aria-expanded="false">Lainnya</button>
  </DropdownMenu>
</nav>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.nav-overflow-chunking intentional exception -->
```

```tsx
// charites:ignore ux.nav-overflow-chunking intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.nav-overflow-chunking:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


