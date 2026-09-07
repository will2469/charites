# Implementation Plan: Context-Aware Intra-Component State-Gating Analysis for `ux.destructive-action-unconfirmed`

Resolves **[GitHub Issue #11](https://github.com/will2469/charites/issues/11)** (`feat(ux): context-aware state-gating analysis for ux.destructive-action-unconfirmed`).

---

## 1. Executive Summary & Problem Formulation

In modern React development (lists, tables, cards, dynamic forms), destructive operations (e.g. deleting a user, removing an item) are standardly staged through local component state rather than executed immediately:

```tsx
const [confirmIndex, setConfirmIndex] = useState<number | null>(null);

// Inside item row/card:
<Button
  type="button"
  variant="ghost"
  onClick={() => setConfirmIndex(index)}
>
  Hapus
</Button>

// Elsewhere in the same component:
<ActionApprovalDialog
  open={confirmIndex !== null}
  title="Hapus Pengikut"
  onConfirm={() => {
    if (confirmIndex !== null) onHapus(confirmIndex);
    setConfirmIndex(null);
  }}
  onCancel={() => setConfirmIndex(null)}
/>
```

In this pattern:
1. The button **does not** execute the destructive mutation directly.
2. It only activates a staged confirmation state (`setConfirmIndex(index)`).
3. The actual destructive mutation is safely executed downstream inside the `<ActionApprovalDialog>` confirmation handler (`onConfirm`).

Previously, `ux.destructive-action-unconfirmed` only recognized direct inline confirmation (`<AlertDialogTrigger>` ancestor or `window.confirm()` in the handler). It produced false-positive errors on this standard staged pattern.

---

## 2. The 6 Locked Architectural Invariants

1. **State Binding via Declaration, NOT Naming Conventions:**
   State identity is resolved from `useState` declarations in the component AST (`const [state, setter] = useState(...)`). Naming conventions (`setX -> x`) serve only as secondary fallback, never authoritative identity.
2. **Direct Mutation Disqualification (Zero Bypass on Mixed Handlers):**
   Call-level inspection of the handler. If the handler invokes ANY direct destructive mutation callee (`delete*`, `remove*`, `hapus*`, `destroy*`, `purge*`, `revoke*`, `api.delete*`, `mutate*`), state gating is rejected.
3. **Strict Intra-Component Scope Boundary:**
   State analysis is strictly scoped to the nearest enclosing component/function. If state or setters originate from props or outer components, analysis returns `StateGateUnknown` -> warning emitted.
4. **Downstream Confirmation Primitive & Action Requirement:**
   A state setter alone never proves confirmation. The state variable must be consumed by a confirmation dialog (`<Dialog>`, `<AlertDialog>`, `<ActionApprovalDialog>`, `<ConfirmDialog>`, `<Modal>`) AND the dialog must feature a confirmation action (`onConfirm`, `onApprove`, `onDelete`, `onAccept`, `<AlertDialogAction>`). Informational dialogs do not qualify.
5. **Strict State-Consumption Predicate Contract:**
   Supported visibility predicates: `state !== null`, `state != null`, `state === true`, `state`, `Boolean(state)`, `!!state`, and conditional render guards (`{state && <Dialog ...>}`).
6. **Conservative Rule Policy:**
   - `StateGateConfirmed` -> Suppress diagnostic (0 warnings/errors)
   - `StateGateNotConfirmed` -> Emit diagnostic (1 error)
   - `StateGateUnknown` -> Emit diagnostic (1 error)

---

## 3. Implemented Changes

### Component 1: Dedicated State-Gating Analyzer (`internal/rules/ux/destructive_action_state.go`)
- Implemented `StateGateResult` enum (`StateGateUnknown`, `StateGateConfirmed`, `StateGateNotConfirmed`).
- Implemented `AnalyzeStateGating(node *ir.Node) StateGateResult` with:
  - Direct destructive mutation detection at call-level (differentiating state setters like `setDeleteDialogOpen` from actual destructive calls like `deleteUser`).
  - Parameter-aware component scope detection (`findEnclosingComponentScope`, `findComponentDeclarationStart`, `findComponentScopeEnd`) correctly handling parameter destructuring `{ ... }` before function body `{ ... }`.
  - Exact `useState` AST extraction and setter-to-state resolution.
  - Downstream confirmation primitive verification (`Dialog`, `AlertDialog`, `ActionApprovalDialog`, `ConfirmDialog`, `Modal`, `Drawer`).
  - Supported visibility predicate checks and JSX conditional guard recognition.
  - Confirmation action verification (`onConfirm`, `onApprove`, `onDelete`, `confirmText`, `<AlertDialogAction>`, or confirm button children). Decomposed into small helpers to maintain low cognitive complexity (< 10).

### Component 2: Rule Integration (`internal/rules/ux/destructive_action_unconfirmed.go` & `util.go`)
- Kept `Evaluate()` thin: `if AnalyzeStateGating(node) == StateGateConfirmed { return nil }`.
- Expanded `isInteractiveElement` to cover native `button`, `a`, `*Button`, `*Link`, `role=button/link/menuitem`, and `onClick`/`onPress` event handlers.
- Updated `isDestructiveStyledElement` in `util.go` to recognize `variant="destructive"` and `variant="danger"`.
- Excluded event handler attributes from generic string-matching in `hasConfirmationAttribute` in `util.go`.

### Component 3: 17-Case Acceptance Test Suite (`internal/rules/ux/destructive_action_unconfirmed_test.go`)
- Added comprehensive unit tests for all 17 cases:
  1. Direct `deleteUser()` call -> 1 diagnostic.
  2. Wrapped in `AlertDialogTrigger` -> 0 diagnostics.
  3. Inline `window.confirm()` check -> 0 diagnostics.
  4. Staged `setConfirmIndex` -> `open={confirmIndex !== null}` -> `onConfirm` (Issue #11 specimen) -> 0 diagnostics.
  5. `setIsOpen(true)` with informational dialog -> 1 diagnostic.
  6. Unrelated `setSelectedTab()` with direct `deleteUser` in same handler -> 1 diagnostic.
  7. Setter and dialog unrelated to state -> 1 diagnostic.
  8. State declared in another component -> 1 diagnostic.
  9. Setter passed as prop -> 1 diagnostic.
  10. `setDeleteDialogOpen(true)` with valid downstream confirm action -> 0 diagnostics.
  11. `setPendingDeleteId(id)` -> `open={Boolean(pendingDeleteId)}` -> `onConfirm` -> 0 diagnostics.
  12. `pendingDeleteId && <ConfirmDialog onConfirm={...}>` (conditional guard) -> 0 diagnostics.
  13. `pendingDeleteId !== null` but informational `<Dialog>` -> 1 diagnostic.
  14. Direct delete + `setConfirmOpen(true)` in same handler -> 1 diagnostic.
  15. Computed/dynamic state expression -> 1 diagnostic.
  16. Unknown setter naming resolved via `useState` binding -> 0 diagnostics.
  17. Multiple destructive buttons mapping to same confirmation state (list/table) -> 0 diagnostics.

### Component 4: 1-SSOT Tri-Corpus Correctness Fixtures
- **Negative (`negative/`):**
  - `state_gated_approval_dialog.tsx`: Exact specimen from Issue #11.
  - `state_gated_pending_id.tsx`: Table buttons with `setPendingDeleteId` and `<ConfirmDialog>`.
  - `state_gated_custom_setter.tsx`: Custom setter naming resolved via `useState` binding.
- **Positive (`positive/`):**
  - `unrelated_setter_with_delete.tsx`: Direct mutation alongside state setter.
  - `informational_dialog_only.tsx`: Dialog lacking confirmation action.
  - `setter_passed_as_prop.tsx`: Setter passed via props.

---

## 4. Verification Summary
- Unit Tests: `go test -v ./internal/rules/ux -run TestDestructiveActionUnconfirmedRule` (17/17 passed).
- Tri-Corpus Harness: `go test -v ./tests/correctness/ux/destructive-action-unconfirmed/...` (100% pass across Positive, Negative, Adversarial).
- Canonical Contract: `TestUXRules_CanonicalContract` passed.
- Adoption Matrix & Correctness Gate: `go test -v ./tests -run "TestGoldenCorpus_AdoptionMatrix|TestCorrectnessGate"` passed.
- Linter: `make lint` clean (zero warnings, gocognit and gofmt compliant).
- Wiki Generator: `make wiki` executed successfully.
- Race Detector: `go test -race ./...` (zero race conditions).
