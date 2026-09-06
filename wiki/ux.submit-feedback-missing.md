# ux.submit-feedback-missing

> **Rule ID:** `ux.submit-feedback-missing`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Doherty Threshold (< 400ms Interaction Feedback), Nielsen Heuristic #1: Visibility of System Status, ISO 9241-110 Ergonomics of Human-System Interaction (Self-Descriptiveness)

---

## 1. Overview & Core Invariant

Enforces reentry guard (disabled) and perceivable feedback (aria-busy/spinner) on async mutation triggers

### Core Invariant:
> **"Async mutation triggers must provide both Reentry Guard (R1: 'disabled={isPending}') and Perceivable Feedback (R2: 'aria-busy' or spinner)."**

---
## 2. Technical Grounding & Engine Realities

When users trigger a mutation (such as submitting an order or charging a payment) without immediate feedback, they perceive the system as unresponsive within the 400ms Doherty Threshold.

In the absence of a reentry lock, users repeatedly click the submit button, triggering duplicate requests, double financial charges, and server-side race conditions. Enforcing both R1 (reentry lockout) and R2 (visual feedback) guarantees transactional safety and cognitive assurance.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Duplicate Submissions & Double Charging** | HIGH | Rapid repeated clicks by uncertain users cause duplicate database inserts and multiple payment gateway captures. |
| **Cognitive Disorientation & Interaction Rage** | MEDIUM | Users doubt whether their input registered, leading to rage clicks and workflow cancellation. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Async submit button without disabling control or showing pending feedback):
```tsx
<button
  onClick={async () => {
    await api.post("/orders", orderData);
  }}
  className="bg-primary text-white px-4 py-2"
>
  Bayar Sekarang
</button>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Compliant button satisfying R1 (disabled={isPending}) and R2 (aria-busy & dynamic text)):
```tsx
<button
  onClick={handlePayment}
  disabled={isPending}
  aria-busy={isPending}
  className="bg-primary text-white px-4 py-2 disabled:opacity-50"
>
  {isPending ? "Memproses Pembayaran..." : "Bayar Sekarang"}
</button>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.submit-feedback-missing intentional exception -->
```

```tsx
// charites:ignore ux.submit-feedback-missing intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.submit-feedback-missing:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


