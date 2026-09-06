# ux.unconventional-home-link

> **Rule ID:** `ux.unconventional-home-link`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Jakob's Law of Internet User Experience (Nielsen Norman Group), W3C Web Navigation & Landmark Architecture Guidelines, ISO 9241-110 Ergonomics of Human-System Interaction (Suitability for Learning & Predictability)

---

## 1. Overview & Core Invariant

Enforces Jakob's Law by ensuring header logo/brand identity links to the root home page ('/')

### Core Invariant:
> **"Brand identity and logo elements in the primary header must be enclosed within an anchor ('<a>' or '<Link>') whose destination normalizes to the root homepage ('/')."**

---
## 2. Technical Grounding & Engine Realities

Jakob's Law states that users spend most of their time on sites other than yours. Consequently, they bring deeply ingrained mental models about standard interaction patterns. The most universal web convention is that clicking the brand logo in the top-left header returns to the homepage ('/').

When a logo is unclickable, rendered as a passive image or plain text, or links to an unexpected secondary destination (like /about, /products, or an external portal), users become disoriented and lose their primary visual escape hatch back to the system root.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Mental Model Disorientation** | MEDIUM | Users habitually click the top-left brand mark when seeking the homepage; non-functional or diverted logos induce frustration and cognitive friction. |
| **Accidental Site Exit or Dead End** | LOW | Navigating away from the root application when attempting to reset context forces users to rely on browser history or address bar edits. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Passive brand logo in header without any enclosing link):
```tsx
<header className="flex items-center justify-between px-6 py-4 border-b">
  <img src="/brand-logo.svg" alt="Acme Corporation Logo" className="h-8 w-auto" />
  <nav className="flex gap-4">
    <a href="/features">Features</a>
    <a href="/pricing">Pricing</a>
  </nav>
</header>
```
### ASTRO (Brand logo linking to an internal sub-page instead of root):
```astro
<header class="flex items-center justify-between px-6 py-4">
  <a href="/about" class="brand-logo flex items-center gap-2">
    <img src="/logo.svg" alt="Brand Logo" />
    <span class="font-bold">Portal</span>
  </a>
</header>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Brand logo wrapped in accessible anchor linking directly to root '/'):
```tsx
<header className="flex items-center justify-between px-6 py-4 border-b">
  <a href="/" aria-label="Acme Corporation - Beranda" className="flex items-center gap-2">
    <img src="/brand-logo.svg" alt="Acme Corporation Logo" className="h-8 w-auto" />
    <span className="font-bold text-lg">Acme</span>
  </a>
  <nav className="flex gap-4">
    <a href="/features">Features</a>
  </nav>
</header>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.unconventional-home-link intentional exception -->
```

```tsx
// charites:ignore ux.unconventional-home-link intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.unconventional-home-link:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


