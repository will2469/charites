# Charites Static Analysis Rule Catalog

Welcome to the **Charites Static Analysis Rule Catalog**. Charites is an ultra-fast, zero-CGO, zero-Node.js static analysis compiler for Astro, React TSX, and Tailwind CSS design tokens.

---

## Categories

| Category | Domain | Rules Count | Documentation |
| :--- | :--- | :---: | :--- |
| `theme` | Design tokens, color system, OKLCH, and opacity enforcement | 1 | [`theme.md`](theme.md) |

---

## All Registered Rules

| Rule ID | Category | Severity | Description | Documentation |
| :--- | :---: | :---: | :--- | :--- |
| `theme.hardcode-opacity-color` | `theme` | `ERROR` | Detects utility classes with hardcoded slash opacity modifiers that have official semantic token replacements | [`theme.md#themehardcode-opacity-color`](theme.md#themehardcode-opacity-color) / [`theme.hardcode-opacity-color.md`](theme.hardcode-opacity-color.md) |

---

## Architectural Principles

1. **Deterministic Execution:** Pure-function AST visitors without file system or network I/O during evaluation.
2. **1-SSOT Tri-Corpus Assurance:** Every rule is validated against a 3-part golden test corpus (`positive/`, `negative/`, `adversarial/`).
3. **Canonical Semgrep Identifiers:** All rules follow the `<category>.<slug>` standard.
