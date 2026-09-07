# design.component-style-drift

> **Rule ID:** `design.component-style-drift`
> **Severity:** `WARN`
> **Category:** `design`
> **Target Standards:** W3C Design Tokens Community Group (DTCG) 3-Tier Architecture, W3C Web Content Accessibility Guidelines (WCAG) 2.2 Criterion 1.4.3 (Contrast Minimum), Concentric Border Radius Geometry (R_inner = R_outer - Padding)

---

## 1. Overview & Core Invariant

Clusters and flags design system style sprawl, rogue geometric outliers, and chromatic color drift per UI component

### Core Invariant:
> **"UI components must maintain uniform geometric tokens (radii, elevation, press feedback) and accessible chromatic contrast across the repository without rogue style drift."**

---
## 2. Technical Grounding & Engine Realities

In large monorepos and evolving design systems, developers frequently introduce Style Drift (Design Entropy).

While static rules catch raw hex codes on a single-file basis, they are blind to cross-file design entropy where the same UI primitive splinters into subtle variations:
- <Button> using rounded-md (98% dominant), but isolated features use rounded-2xl or rounded-none.
- <Button> using active:scale-95, but rogue pages use active:scale-[0.98].
- Lack of first-class variants (e.g. missing 'warning' variant on <Button>) forcing developers to improvise with competing raw Tailwind hues (yellow-400, amber-500, orange-500), creating severe WCAG 1.4.3 contrast failures (bg-yellow-400 with white text = 1.53:1).

Concentric border radius geometry dictates that outer surfaces (<Card> using rounded-xl/2xl) must legitimately differ from inner elements (<Button> using rounded-md). Therefore, Charites enforces Component-Scoped statistical clustering.

Repository analysis groups styles by (ComponentKind, PropertyCategory), applies a minimum sample gate (10 occurrences), determines canonical standards (>= 80%), and flags rogue outliers (<= 10%).

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Visual & Brand Incoherence** | MEDIUM | Fragmented corner radii and micro-interactions degrade perceived app polish and user confidence. |
| **WCAG Accessibility Violations** | HIGH | Ad-hoc yellow and amber button backgrounds paired with white text fail WCAG AA contrast (minimum 4.5:1), rendering labels unreadable for low-vision users. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Severe WCAG contrast hazard with primitive background and white text):
```tsx
<Button className="bg-yellow-400 text-white">Warning Action</Button>
```
### TSX (Rogue outlier corner radius on button diverging from design system canonical rounded-md):
```tsx
<Button className="rounded-2xl">Checkout</Button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Canonical button styling with accessible contrast and semantic tokens):
```tsx
<Button className="rounded-md active:scale-95 bg-primary text-primary-foreground">Submit</Button>
```
### TSX (Card surface container legitimately using larger outer radius):
```tsx
<Card className="rounded-xl border border-border/40 p-6">Card Content</Card>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore design.component-style-drift intentional exception -->
```

```tsx
// charites:ignore design.component-style-drift intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  design.component-style-drift:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [design Category Guide](design).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


