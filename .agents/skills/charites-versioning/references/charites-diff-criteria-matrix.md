# Charites Diff Criteria Matrix: Classifying Latest Tag to Current HEAD
**Purpose:** Exact criteria matrix to evaluate `git diff <latest_tag>...HEAD` specifically for the **Charites** repository (`github.com/will2469/charites`) and its Astro, TSX, CSS, and Go compiler ecosystem.

---

## 1. Charites Component Diff Classification Matrix

| Component / Subsystem | **MAJOR (Breaking)** `1.x.x` → `2.0.0` | **MINOR (Additive)** `1.4.0` → `1.5.0` | **PATCH (Fix/Hygiene)** `1.4.1` → `1.4.2` |
| :--- | :--- | :--- | :--- |
| **CLI & Commands** (`cmd/charites/`, `internal/cli/`) | • Removed subcommand (`scan`, `check`, `run`, `wiki`, `mcp`, `update`)<br>• Removed/renamed CLI flag (`--format`, `--ext`, `--category`, `--rule`, `--config`)<br>• Changed default behavior of a CLI flag<br>• Altered JSON stream machine-readable schema (`--format=json`)<br>• Changed process exit codes (0/1/2) | • Added new subcommand or new optional CLI flag<br>• Added new supported output format (e.g. `--format=markdown`)<br>• Added new environment variable override | • Fixed CLI flag parsing crash / nil panic<br>• Improved CLI terminal human-readable ANSI output styling<br>• Fixed `--help` text typos or shell completions |
| **Analyzer Rules** (`internal/rules/`) | • Removed an existing rule analyzer<br>• Changed rule Semgrep ID (e.g. `theme.hardcode-color` renamed without alias)<br>• Incompatible diagnostic position or span reporting break | • Added new rule analyzer (e.g. new `theme.*`, `a11y.*`, `perf.*` rules)<br>• Added new companion AST visitor<br>• Added `// Deprecated:` notice to an existing rule | • Fixed false-positive (FP) in template AST traversal<br>• Fixed false-negative (FN) in edge-case handling<br>• Added missing design token whitelist (e.g. OKLCH standard values) |
| **Configuration** (`charites.yaml`, `.charitesignore`, `internal/config/`) | • Removed existing config key or renamed section<br>• Made previously optional YAML key required<br>• Changed default severity of an existing rule<br>• Made validation reject previously valid `charites.yaml` | • Added new optional configuration key or section<br>• Added support for new directive alias in YAML<br>• Added non-breaking default overrides | • Fixed YAML parsing error message clarity<br>• Fixed unhandled nil pointer in empty config file<br>• Corrected config file search path precedence bug |
| **Leaf IR & Parsers** (`internal/ir/`, `internal/parser/`) | • Removed exported field or method on `ir.Node`<br>• Changed node type discrimination (`NodeType`) incompatibly<br>• Modified traversal signature in `Walk()` | • Added new node property or AST metadata<br>• Added parser support for new file extension (e.g. `.mdx`)<br>• Added `// Deprecated:` doc comment | • Fixed frontmatter line offset calculation bug<br>• Optimized parser memory pool (`sync.Pool`)<br>• Fixed memory leak in CSS/HTML tokenizer |
| **MCP Server** (`internal/mcp/`) | • Removed an MCP tool (`charites_scan`, `charites_explain_rule`, `charites_list_rules`)<br>• Renamed tool or removed required argument<br>• Changed tool JSON schema incompatibly<br>• Changed MCP protocol wire version incompatibly | • Added new MCP tool or resource<br>• Added optional property to existing tool input schema<br>• Added new diagnostic metadata in response | • Fixed MCP session lifecycle or pipe EOF bug<br>• Fixed `_meta` protocol version fallback<br>• Concurrency race condition fix in tool registry |
| **Testing Harness** (`tests/correctness/`, `tests/e2e/`) | • Deleted golden corpus fixture categories | • Added new golden corpus test cases (P1-P5, N1-N5, A1-A7) | • Fixed flaky test or updated golden snapshot assertions |
| **Runtime & Dependencies** (`go.mod`) | • Raised minimum Go version (e.g. Go 1.26 → Go 1.27)<br>• Replaced parser dependencies with breaking breaking API changes | • Added new development or linting dependency<br>• Bumped minor dependency version | • Bumped patch dependency for security vulnerability (CVE) |

---

## 2. Git Diff Evaluation Pipeline

### Step 1: Detect Base Tag
```bash
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || git tag -l --sort=-v:refname | head -n 1)
```

### Step 2: Three-Dot Diff Inspection
```bash
# 1. Overview of changed files
git diff --stat "${LATEST_TAG}...HEAD"

# 2. Check for breaking deletions or signature changes
git diff "${LATEST_TAG}...HEAD" -- 'cmd/' 'internal/rules/' 'internal/ir/' 'internal/mcp/'
```

### Step 3: Conventional Commit Scan
```bash
# Check for breaking change markers
git log "${LATEST_TAG}..HEAD" --format="%s%n%b"
```
- **MAJOR:** `type!:` (e.g. `feat!:`, `fix!:`) or body containing `BREAKING CHANGE:`.
- **MINOR:** `feat:` (without `!`).
- **PATCH:** `fix:`, `perf:`, `refactor:`, `docs:`, `test:`, `chore:`.

---

## 3. Presedensi & Invariant Reset

$$\mathbf{MAJOR > MINOR > PATCH}$$

1. **MAJOR Triggered:**
   `1.4.2` → **`2.0.0`** *(Reset Minor & Patch ke 0)*.
2. **MINOR Triggered (No Major):**
   `1.4.2` → **`1.5.0`** *(Reset Patch ke 0)*.
3. **PATCH Triggered (No Major & No Minor):**
   `1.4.2` → **`1.4.3`** *(Hanya Patch naik)*.
