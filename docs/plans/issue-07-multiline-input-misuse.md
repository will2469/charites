# Implementation Plan: Multiline Input Misuse Rule (Issue #7)

Resolves **[GitHub Issue #7](https://github.com/will2469/charites/issues/7)** (`ux.multiline-input-misuse`: feat(ux): rule to flag notes/keterangan fields using single-line `<Input>` instead of `<Textarea>`).

---

## 1. Executive Summary & Problem Statement

In web forms, open-ended commentary, notes, descriptions, and rejection reasons frequently misuse single-line `<Input>` / `<input>` controls instead of multiline `<Textarea>`:
```tsx
{/* Specimen from apps/web/src/pages/desa/layanan/surat/create.astro */}
<Label htmlFor="keterangan-input">Keterangan (Opsional)</Label>
<Input id="keterangan-input" placeholder="Catatan tambahan (bila ada)" />
```

This antipattern produces significant usability defects:
1. **Text Review Obstruction (NN/g & GOV.UK Design System):** Single-line inputs clip overflowing text horizontally, preventing users from reviewing their explanations without manual cursor scrubbing.
2. **Loss of Structural Formatting:** Users cannot insert line breaks or structure thoughts into distinct paragraphs.
3. **Control-Type Mismatch:** According to GOV.UK Design System standards, text inputs must only be used for short single-line entries, whereas open-ended commentary requires `<textarea>`.

---

## 2. The 4 Locked Architectural Invariants

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THE 4 LOCKED ARCHITECTURAL INVARIANTS                    │
├─────────────────────────────────────────────────────────────────────────────┤
│ 1. Orthogonal Semantic Pipelines:                                           │
│    IdentifierEvidence (Identity, Quantity, Unknown) and ContentEvidence     │
│    (Multiline, SingleLine, Unknown) are strictly orthogonal.                │
│    Issue #5 (number inputs) and Issue #7 (multiline inputs) never mutate    │
│    or cross-contaminate each other's evidence channels.                     │
│                                                                             │
│ 2. Tri-State Input Type Classification:                                     │
│    InputTypeClass represents:                                               │
│    - InputTypeText: explicit 'text', empty '', or omitted type              │
│    - InputTypeOther: explicit 'number', 'password', 'email', 'tel', etc.    │
│    - InputTypeUnknown: dynamic expressions (e.g. type={customType})         │
│    The rule strictly triggers ONLY when TypeClass == InputTypeText.         │
│                                                                             │
│ 3. Multi-Channel Evaluation with Qualifier Overrides:                       │
│    Evaluates all channels ('name', 'id', 'placeholder', 'aria-label') with: │
│    Strong Single-Line > Strong Multiline > Contextual Single-Line > Unknown. │
│    E.g. <Input name="description" placeholder="Short description" />         │
│    resolves to SingleLine (0 violations) due to the 'short' qualifier.      │
│                                                                             │
│ 4. Channel-Aware Tokenization & Exclusions:                                 │
│    - 'name' & 'id' use TokenizeIdentifier (camelCase, snake_case)           │
│    - 'placeholder' & 'aria-label' use TokenizeTextWords (punctuation split) │
│    - 'desc' is deliberately excluded to prevent collision with sort_desc   │
│    - 'data-testid' is completely excluded from semantic detection in v1     │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Architecture & Shared Adapter Pipelines

```text
                     AST / IR Node
                           │
                           ▼
                  ExtractInputFacts(node)
                           │
         ┌─────────────────┴─────────────────┐
         ▼                                   ▼
ResolveIdentifierEvidence           ResolveContentEvidence
  (name, id, autocomplete,           (name, id, placeholder, aria-label)
   aria-label, placeholder)                  │
         │                                   ▼
         ▼                          Evaluate All Channels
   Identity / Quantity /              Rank with Qualifier Overrides
   Unknown                                   │
                                             ▼
                                    Multiline / SingleLine / Unknown
                                             │
                                             ▼
                               InputFacts.Content.Intent
```

### Content Intent Classification Hierarchy

1. **Strong Single-Line Qualifier (Precedence 1):**
   `code`, `title`, `name`, `id`, `short`, `brief`, `summary`, `count`, `number`, `type`, `status`, `date`, `tanggal`, `judul`, `kode`, `singkat`, `ringkas`, `jumlah`, `nomor`, `no`.
   If matched $\rightarrow$ `ContentIntentSingleLine`.

2. **Strong Multiline Semantic Token (Precedence 2):**
   `keterangan`, `catatan`, `notes`, `note`, `deskripsi`, `description`, `alasan`, `reason`, `komentar`, `comment`, `uraian`, `penjelasan`, `explanation`, `pesan`, `message`, `feedback`, `keluhan`, `saran`.
   If matched $\rightarrow$ `ContentIntentMultiline`.

3. **Contextual Single-Line Tokens (Precedence 3):**
   If a contextual token matches without multiline tokens $\rightarrow$ `ContentIntentSingleLine`.

---

## 4. Verification Results

### Unit Tests
- `internal/rules/form/identifier_test.go`: 100% pass across word tokenization, identifier tokenization, and content intent classification.
- `internal/rules/form/input_test.go`: 100% pass across tri-state input types, multi-channel extraction, and qualifier overrides.
- `internal/rules/ux/multiline_input_misuse_test.go`: 20 unit tests covering Astrades specimens, negative controls, qualifier overrides, and dynamic edge cases.

### 1-SSOT Tri-Corpus Correctness
- `tests/correctness/ux/multiline-input-misuse/positive/`:
  - `astrades_specimen.tsx`
  - `rejection_reason.astro`
  - `customer_notes.tsx`
  - `placeholder_multiline.tsx`
  - `aria_label_multiline.tsx`
- `tests/correctness/ux/multiline-input-misuse/negative/`:
  - `proper_textarea.tsx`
  - `note_title.tsx`
  - `reason_code.astro`
  - `proper_textarea.astro`
  - `short_description.tsx`
- `tests/correctness/ux/multiline-input-misuse/adversarial/`:
  - `password_notes.tsx`
  - `number_notes_count.tsx`
  - `dynamic_type_notes.tsx`
  - `abbreviated_desc.tsx`
  - `data_testid_bait.tsx`

### Cross-Rule Ownership Matrix
- `tests/correctness/cross_rule_ownership_test.go`: 15 scenarios testing interaction between `ux.number-input-identity-misuse`, `ux.number-input-missing-bounds`, `ergonomy.number-input-wheel-hazard`, and `ux.multiline-input-misuse`.
