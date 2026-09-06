# performance.react-context-domain-coupling

> **Rule ID:** `performance.react-context-domain-coupling`
> **Severity:** `WARN`
> **Category:** `performance`
> **Target Standards:** React Official Architecture (Context Modularity & State Granularity Principles), React Virtual DOM Reconciliation Invariants (Consumer tree re-render pruning), Web Performance Working Group Main-Thread Optimization Guidelines

---

## 1. Overview & Core Invariant

Context.Provider bundles over-coupled multi-domain state (> 5 fields), triggering cascading re-renders across all consumers on any property change

### Core Invariant:
> **"React Context Providers must maintain fine-grained domain boundaries; bundling disparate application states into a monolithic 'God Context' forces all consumer components to re-render whenever any unrelated field updates."**

---
## 2. Technical Grounding & Engine Realities

React Context propagates state updates to all consumers without granular field-level selector filtering. When an application bundles multiple disparate domains (e.g. user authentication, shopping cart, UI modal state, notification badge count, and scroll position) into a single Provider, any update to a high-frequency field (such as a badge increment) triggers a re-render across every component subscribed to the context.

This monolithic coupling bypasses component memoization and saturates the main thread with unnecessary render passes.

Splitting state into focused, domain-specific contexts (e.g. `AuthContext`, `CartContext`, `UIContext`) ensures components only re-render when their specific domain data actually mutates.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Cascading Re-render Blast Radius** | HIGH | A change in a single property (e.g. notification count) forces dozens of unrelated UI components to re-render simultaneously. |
| **Severe Main-Thread Jitter** | MEDIUM | High-frequency state mutations throttle the main thread, resulting in dropped frames and delayed user interaction responses. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Monolithic God Context bundling 7 distinct domain properties in a single provider):
```tsx
<AppContext.Provider value={{
  user,
  cart,
  theme,
  notifications,
  activeModal,
  isSidebarOpen,
  locale,
}}>
  {children}
</AppContext.Provider>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (State decoupled into domain-specific, modular context providers):
```tsx
<AuthProvider value={authValue}>
  <CartProvider value={cartValue}>
    <UIProvider value={uiValue}>
      {children}
    </UIProvider>
  </CartProvider>
</AuthProvider>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore performance.react-context-domain-coupling intentional exception -->
```

```tsx
// charites:ignore performance.react-context-domain-coupling intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  performance.react-context-domain-coupling:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [performance Category Guide](performance).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


