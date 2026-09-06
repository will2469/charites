# ux.destructive-action-unconfirmed

> **Rule ID:** `ux.destructive-action-unconfirmed`
> **Severity:** `ERROR`
> **Category:** `ux`
> **Target Standards:** Nielsen Heuristic #5: Error Prevention (Slips and Lapses), ISO 9241-110 Ergonomics of Human-System Interaction (Suitability for the Task), Material Design & WCAG Defensive Action Guidelines

---

## 1. Overview & Core Invariant

Enforces confirmation gating for destructive actions to prevent accidental data loss from slips

### Core Invariant:
> **"Destructive user operations ('delete', 'remove', 'destroy', 'purge', 'revoke') must be gated by a confirmation dialog or 2-step verification."**

---
## 2. Technical Grounding & Engine Realities

Destructive actions such as deleting user accounts, clearing billing databases, or revoking credentials cause permanent, often irreversible data loss.

Executing these operations on a single click without confirmation exposes users to motor slips, touchscreen taps during scrolling, and mistaken identity clicks. Gating destructive actions behind an explicit confirmation dialog (e.g. '<AlertDialog>' or 'window.confirm') provides a cognitive pause and protects against catastrophic slips.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Irreversible Data Destruction from Motor Slips** | CRITICAL | Accidental single-click actions immediately wipe critical business data or terminate accounts without user consent. |
| **User Anxiety & Hesitation** | MEDIUM | Users fear interacting with danger-styled buttons when no confirmation boundary protects them from permanent loss. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Destructive button triggering account deletion directly on single click):
```tsx
<button
  onClick={() => deleteAccount(user.id)}
  className="bg-destructive text-destructive-foreground px-4 py-2 rounded"
>
  Hapus Akun Permanen
</button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Destructive button safely wrapped in AlertDialogTrigger confirmation modal):
```tsx
<AlertDialogTrigger asChild>
  <button className="bg-destructive text-destructive-foreground px-4 py-2 rounded">
    Hapus Akun Permanen
  </button>
</AlertDialogTrigger>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.destructive-action-unconfirmed intentional exception -->
```

```tsx
// charites:ignore ux.destructive-action-unconfirmed intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.destructive-action-unconfirmed:
    severity: error # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


