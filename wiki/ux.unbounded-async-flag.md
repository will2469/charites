# ux.unbounded-async-flag

> **Rule ID:** `ux.unbounded-async-flag`
> **Severity:** `ERROR`
> **Category:** `ux`
> **Target Standards:** Mental Model Continuity & Deadlock Prevention, Nielsen Heuristic #1: Visibility of System Status, ISO 9241-110 Ergonomics of Human-System Interaction (Error Tolerance)

---

## 1. Overview & Core Invariant

Detects async handlers setting loading flags without guaranteed reset in finally/catch exit paths

### Core Invariant:
> **"Async operations setting loading state before 'await' must guarantee state reset in all exit paths or a 'finally' block."**

---
## 2. Technical Grounding & Engine Realities

When asynchronous functions activate loading state (e.g. 'setLoading(true)') before awaiting a network operation and fail to ensure that the flag is reset in a 'finally' block or error handler, any unexpected rejection (500 internal server error, timeout, network dropout) leaves the UI permanently frozen in a loading spinner.

Users can neither re-try the action nor interact with adjacent controls, creating an unrecoverable dead end.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Permanent UI Lockup & Infinite Spinner Deadlock** | HIGH | Failed API requests leave spinners active indefinitely, blocking subsequent interactions and forcing users to hard-refresh. |
| **Silent Failure Masking** | HIGH | Users believe an operation is still in flight even after the underlying request failed and aborted. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Loading state activated before await without reset in catch/finally block):
```tsx
<button
  onClick={async () => {
    setLoading(true);
    try {
      await api.fetchUsers();
    } catch (err) {
      console.error(err);
      // setLoading(false) terlupakan! Spinner berputar selamanya saat API gagal.
    }
  }}
>
  Muat Data
</button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Guaranteed loading reset in finally block ensuring UI unlock under all outcomes):
```tsx
<button
  onClick={async () => {
    setLoading(true);
    try {
      await api.fetchUsers();
    } finally {
      setLoading(false);
    }
  }}
>
  Muat Data
</button>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.unbounded-async-flag intentional exception -->
```

```tsx
// charites:ignore ux.unbounded-async-flag intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.unbounded-async-flag:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


