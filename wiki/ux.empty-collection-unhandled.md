# ux.empty-collection-unhandled

> **Rule ID:** `ux.empty-collection-unhandled`
> **Severity:** `INFO`
> **Category:** `ux`
> **Target Standards:** Zero-State Usability & Mental Model Continuity (Nielsen Norman Group), Feedforward Principle & Gulf of Evaluation (Don Norman), ISO 9241-110 Ergonomics of Human-System Interaction (Suitability for Learning)

---

## 1. Overview & Core Invariant

Advises handling empty collection state when mapping dynamic items to avoid zero-state blindness

### Core Invariant:
> **"Dynamic collection rendering expressions must handle empty collection states ('collection.length === 0') with informative fallback zero-state UI."**

---
## 2. Technical Grounding & Engine Realities

When dynamic lists, tables, or feed collections contain 0 records and render nothing, users are stranded in an ambiguous vacuum: did the request fail, is it still loading, or are there genuinely zero records?

Zero-state blindness forces users to refresh repeatedly or assume the application is broken. A dedicated empty state component (e.g. '<EmptyState />' with an illustration, clarifying text, and a call-to-action button) confirms system status and proactively guides user next steps.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Zero-State Blindness & System Status Ambiguity** | LOW | Users perceive blank empty containers as silent application crashes or perpetual loading freezes. |
| **Workflow Dead Ends** | LOW | Without an actionable empty state CTA (e.g. 'Create First Invoice'), users cannot self-discover how to populate the collection. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Dynamic list rendering items via .map() without handling empty array state):
```tsx
<div className="space-y-3">
  <h2 className="text-lg font-bold">Daftar Tagihan</h2>
  <List items={invoices.map(inv => <InvoiceRow key={inv.id} data={inv} />)} />
</div>
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Explicit empty state fallback branch when collection has 0 items):
```tsx
<div className="space-y-3">
  <h2 className="text-lg font-bold">Daftar Tagihan</h2>
  {invoices.length === 0 ? (
    <EmptyState
      title="Belum Ada Tagihan"
      description="Buat tagihan pertama Anda untuk mulai menerima pembayaran."
      actionText="Buat Tagihan"
    />
  ) : (
    <List items={invoices.map(inv => <InvoiceRow key={inv.id} data={inv} />)} />
  )}
</div>
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.empty-collection-unhandled intentional exception -->
```

```tsx
// charites:ignore ux.empty-collection-unhandled intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.empty-collection-unhandled:
    severity: info # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


