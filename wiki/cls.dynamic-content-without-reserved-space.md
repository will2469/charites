# cls.dynamic-content-without-reserved-space

> **Rule ID:** `cls.dynamic-content-without-reserved-space`
> **Severity:** `WARN`
> **Category:** `cls`
> **Target Standards:** Google Core Web Vitals (Dynamic Content Injection Guidelines), React Suspense & Skeleton Architecture Invariants, W3C Cumulative Layout Shift Mitigation

---

## 1. Overview & Core Invariant

Dynamic widget or banner injected in document flow lacks reserved vertical dimensions (min-h/h), risking content reflow

### Core Invariant:
> **"Dynamic in-flow widgets, promotional banners, or asynchronously injected components must be enclosed in containers with reserved vertical dimensions ('min-h-*') or guarded by Suspense fallback skeletons."**

---
## 2. Technical Grounding & Engine Realities

When asynchronous content (such as personalization widgets, promotional announcements, or dynamic notification bars) loads after initial page paint, injecting it directly into normal document flow without pre-allocated vertical space forces all content below it to shift abruptly downward.

This post-load displacement is one of the single largest real-world contributors to poor Cumulative Layout Shift (CLS) scores, frequently leading to accidental user miss-clicks and navigation disorientation.

Enclosing dynamic elements in a container with an explicit minimum height ('min-h-[120px]') or using a matching skeleton placeholder ensures document flow remains completely stable before and after data resolution.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Post-Load In-Flow Content Snapping** | HIGH | Asynchronous widgets popping into document flow push down articles, forms, or buttons while users are reading or interacting. |
| **Accidental Miss-Clicks** | HIGH | Sudden vertical shifts cause users to accidentally click unintended links or submit buttons. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Unreserved dynamic promotional banner in document flow):
```tsx
<main>
  <h1>Artikel</h1>
  <PromoBanner />
  <Content />
</main>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Dynamic banner enclosed in container with reserved min-height):
```tsx
<main>
  <h1>Artikel</h1>
  <div className="min-h-[120px]">
    <PromoBanner />
  </div>
  <Content />
</main>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore cls.dynamic-content-without-reserved-space intentional exception -->
```

```tsx
// charites:ignore cls.dynamic-content-without-reserved-space intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  cls.dynamic-content-without-reserved-space:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [cls Category Guide](cls).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


