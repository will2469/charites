# performance.astro-unnecessary-client-directive

> **Rule ID:** `performance.astro-unnecessary-client-directive`
> **Severity:** `ERROR`
> **Category:** `performance`
> **Target Standards:** Astro Islands Architecture Specification (Zero-JS Baseline Principle), W3C Web Performance Client-Side Script Minimization Invariants, Astro Official Documentation ('Template Directives: client:*')

---

## 1. Overview & Core Invariant

Menegakkan prinsip Zero-JS Astro dengan melarang penambahan direktif hidrasi (client:*) pada komponen antarmuka yang murni statis.

### Core Invariant:
> **"Static UI components must not include 'client:*' hydration directives; adding hydration directives to non-interactive components forces the framework runtime and component bundle to be downloaded, violating Astro's Zero-JS guarantee."**

---
## 2. Technical Grounding & Engine Realities

Astro by default renders all components to pure, static HTML at build time with zero client-side JavaScript overhead.

When a developer unnecessarily adds a `client:*` directive (`client:load`, `client:idle`, `client:visible`) to a purely presentational component, Astro treats it as an interactive island.

This forces the bundler to extract the component into a separate client bundle and ship the framework runtime (such as React or Vue, weighing 30-50KB+) to the browser, needlessly delaying page interactivity and squandering network bandwidth.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Zero-JS Guarantee Violation** | HIGH | Transmits unnecessary framework runtimes and component code to the client, increasing page weight and parse time. |
| **Main Thread Hydration Lag** | MEDIUM | Wastes browser CPU cycles hydrating static DOM trees that have no event listeners or interactive state. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### ASTRO (Header statis dipaksa terhidrasi ke peramban):
```astro
---
import HeaderStatic from '../components/HeaderStatic.tsx';
---
<HeaderStatic client:load title="Selamat Datang" />
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### ASTRO (Dirender sebagai pure static HTML tanpa JavaScript):
```astro
---
import HeaderStatic from '../components/HeaderStatic.tsx';
---
<HeaderStatic title="Selamat Datang" />
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore performance.astro-unnecessary-client-directive intentional exception -->
```

```tsx
// charites:ignore performance.astro-unnecessary-client-directive intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.astro-unnecessary-client-directive:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [performance Category Guide](performance).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


