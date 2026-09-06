# ux.orphaned-error-state

> **Rule ID:** `ux.orphaned-error-state`
> **Severity:** `WARN`
> **Category:** `ux`
> **Target Standards:** Nielsen Heuristic #9: Help Users Recognize, Diagnose, and Recover from Errors, ISO 9241-110 Ergonomics of Human-System Interaction (Error Tolerance), WCAG 2.2 Success Criterion 3.3.1 (Error Identification)

---

## 1. Overview & Core Invariant

Flags error state updates in event handlers that lack corresponding UI error presentation elements

### Core Invariant:
> **"Validation error setters invoked in component handlers must have a corresponding error presentation element in the UI."**

---
## 2. Technical Grounding & Engine Realities

When client-side validation logic flags invalid input and updates internal component state (e.g. 'setEmailError("Format salah")'), that state must be surfaced to the user.

If the error state is never rendered in JSX or communicated via accessible error indicators, the form silently blocks submission while the user remains completely unaware of what went wrong.

---
## 3. Vulnerability & Risk Taxonomy

| Risk Vector | Severity | Impact |
| :--- | :---: | :--- |
| **Silent Submission Failure & Ghost Validation** | HIGH | Users submit forms that silently fail validation without displaying any error messages, causing confusion and frustration. |
| **Inaccessible Error Notification** | MEDIUM | Screen reader users and keyboard navigators receive no feedback when form constraints are violated. |

---
## 4. Non-Compliant Code Patterns (Bad Examples)
### TSX (Event handler updates error state but no error display exists in the JSX):
```tsx
export function LoginForm() {
  const [email, setEmail] = useState("");
  const [emailError, setEmailError] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!email.includes("@")) {
      setEmailError("Format email tidak valid");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <input value={email} onChange={e => setEmail(e.target.value)} />
      <button type="submit">Masuk</button>
    </form>
  );
}
```

---
## 5. Compliant Implementation Patterns (Good Examples)
### TSX (Error state rendered visibly and accessibly with role='alert' and destructive color):
```tsx
export function LoginForm() {
  const [email, setEmail] = useState("");
  const [emailError, setEmailError] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!email.includes("@")) {
      setEmailError("Format email tidak valid");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <input value={email} onChange={e => setEmail(e.target.value)} />
      {emailError && (
        <p role="alert" className="text-sm text-destructive font-medium">
          {emailError}
        </p>
      )}
      <button type="submit">Masuk</button>
    </form>
  );
}
```

---

## 6. How to Suppress (Ignore Directives)

If this pattern is required for an intentional exception, suppress the diagnostic using the canonical Charites Rule ID:

```astro
<!-- charites:ignore ux.orphaned-error-state intentional exception -->
```

```tsx
// charites:ignore ux.orphaned-error-state intentional exception
```

---

## 7. Configuration Reference (`charites.yaml`)

```yaml
rules:
  ux.orphaned-error-state:
    severity: warn # error | warn | info | off
```

---

## 8. Architectural Domain & Verification Reference

- **Domain Architecture & Analysis Pipeline:** See the [ux Category Guide](ux).
- **1-SSOT Golden Tri-Corpus Test Matrix:** See the [Verification Harness](Home#how-testing-works-across-charites-the-4-layer-verification-harness) on the Master Home page.
- **Rule Catalog Index:** See [Master Wiki Home](Home).


