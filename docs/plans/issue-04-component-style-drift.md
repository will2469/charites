# Implementation Plan: Component-Scoped Style Drift Analyzer (Issue #4)

## 1. Executive Summary & Problem Statement
In growing frontend codebases and design systems, developers introduce **Design System Drift (Style Sprawl / Entropi Desain)**:
- `<Button>` using `rounded-md` (98% canonical across the repo), but isolated features use `rounded-2xl`, `rounded-lg`, or `rounded-none`.
- `<Button>` micro-interaction press feedback using `active:scale-95` (88% dominant), but rogue pages use `active:scale-[0.98]`.
- Missing first-class semantic variants (e.g. lack of `warning` variant on `<Button>`) causing developers to improvise with competing raw Tailwind hues (`yellow-400`, `amber-500`, `orange-500`), introducing visual discord and severe WCAG 1.4.3 contrast failures (`bg-yellow-400` with white text = $1.54:1$).

Concentric border radius geometry ($R_{\text{inner}} = R_{\text{outer}} - \text{padding}$) dictates that Cards (`rounded-xl`/`rounded-2xl`) must legitimately differ from inner Buttons/Inputs (`rounded-md`). Therefore, grouping must be **Component-Scoped**.

Furthermore, `rounded-2xl` is **not intrinsically wrong**. If a repository has 150 instances of `rounded-2xl` on `<Button>` and 3 of `rounded-md`, `rounded-2xl` is the canonical standard. Hence, **empirical drift is a repository-level statistical determination, never a single-node guess**.

---

## 2. The 4 Locked Architectural Invariants

### Invariant 1: Race-Free by Ownership (No Shared Collector Mutex)
Workers in `Pool.Run()` extract `[]drift.StyleOccurrence` locally per file. The worker returns `FileAnalysisResult{Diagnostics, Occurrences}`. The single-threaded coordinator/aggregator collects all occurrences into a unified slice.
Zero `sync.Mutex`. Zero shared mutable state across threads. 100% race-free by linear ownership transfer.

### Invariant 2: Repository Analyzer as a First-Class Pipeline
`rules.Rule` remains strictly a per-node evaluator (`Evaluate(node *ir.Node) []ir.Diagnostic`) for local deterministic invariants.
Repository analysis is governed by a distinct extension point:
```go
type RepositoryAnalyzer interface {
    ID() string
    Analyze(occurrences []drift.StyleOccurrence, cfg config.DriftConfig) drift.DriftReport
}
```
In `charites scan`: files are processed in parallel, occurrences are aggregated, and if `design.component-style-drift` is enabled, `drift.Analyze()` executes on the aggregated population to synthesize diagnostics.
In `charites drift [path] [--report]`: runs repository drift analysis directly with rich terminal visualizations.

### Invariant 3: Confidence-Aware Scope Clustering (`ScopeKey`)
Elements are classified with explicit confidence levels:
```go
type ComponentKind int

const (
    KindButton ComponentKind = iota // <Button>, <button>, [role="button"], <TabsTrigger>
    KindInput                       // <Input>, <textarea>, <select>, <SelectTrigger>, [role="textbox"]
    KindCard                        // <Card>
    KindDialog                      // <DialogContent>, [role="dialog"]
    KindSheet                       // <SheetContent>
    KindBadge                       // <Badge>, [role="status"]
    KindMenu                        // PopoverContent, DropdownMenuContent, TooltipContent
    KindAlert                       // <Alert>, [role="alert"]
)

type ScopeConfidence int

const (
    ConfidenceExact    ScopeConfidence = iota // <Button>, <Card>, <Input>
    ConfidenceNative                          // <button>, <textarea>, <select>
    ConfidenceSemantic                        // role="button", role="dialog"
)

type ScopeKey struct {
    Kind       ComponentKind
    Confidence ScopeConfidence
}
```
By default, `Exact` components (`<Button>`) are not merged with `Native` (`<button>`) or `Semantic` (`[role="button"]`) into the same statistical cluster, preventing false comparisons across different abstraction layers.

### Invariant 4: Context-Aware Exceptions & Hand-Written Config Parser
Exceptions support contextual checking (e.g. `rounded-full` on `Button` is permitted when `size="icon"` or explicitly declared in `charites.yaml`):
```yaml
rules:
  design.component-style-drift: warn

drift:
  dominance_threshold: 80.0
  outlier_threshold: 10.0
  min_cluster_occurrences: 10

  exceptions:
    Button:
      rounded:
        - token: rounded-full
          props:
            size: icon
```
The hand-written parser in `internal/config/config.go` is cleanly extended to parse the `drift:` block.

---

## 3. Package & Pipeline Architecture

```text
internal/drift/
├── occurrence.go  # ComponentKind, ScopeConfidence, PropertyCategory, StyleOccurrence
├── classifier.go  # ResolveScope(), ExtractOccurrences(), normalized token/variant parsing
├── cluster.go     # Statistical clustering, MinClusterOccurrences gate, Dominance & Outlier math
├── chromatic.go   # ChromaticFamily (WarmAmber, GreenEmerald, BlueSky, RedRose), CandidateIntent
├── contrast.go    # Static Tailwind sRGB dictionary, RelativeLuminance, ContrastRatio (WCAG 1.4.3)
├── diagnostic.go  # SynthesizeDiagnostics() for scan integration, inline ignore handling
├── analyzer.go    # DriftAnalyzer, Analyze(occurrences, cfg) -> DriftReport, CLI report rendering
└── drift_test.go  # Unit & benchmark test suite
```

### 3.1. Fine-Grained Property Categories
To prevent "prefix soup", properties are strictly disaggregated:
- `CatRounded`: `rounded`, `rounded-md`, `rounded-full`
- `CatActiveScale`: `active:scale-95`
- `CatDisabledOpacity`: `disabled:opacity-50`
- `CatDisabledCursor`: `disabled:cursor-not-allowed`
- `CatFocusRingWidth`: `focus-visible:ring-2`
- `CatFocusRingColor`: `focus-visible:ring-ring`
- `CatFocusRingOffset`: `focus-visible:ring-offset-2`
- `CatShadow`: `shadow-sm`, `shadow-md`
- `CatBorderWidth`: `border`, `border-2`
- `CatBorderColor`: `border-border`
- `CatBackground`: `bg-background`
- `CatHeight`: `h-9`, `h-10` (partitioned by variant/size prop to prevent false positives)

### 3.2. Statistical Decision Engine
```text
Total = Sum(occurrences in ClusterKey{ScopeKey, Category})

if Total < MinClusterOccurrences (default: 10):
    ClusterStatus = InsufficientData (NO DRIFT DECISION)
else:
    Canonical = Token with Max(count)
    Dominance = (Canonical.Count / Total) * 100
    if Dominance >= DominanceThreshold (default: 80%):
        for each Token != Canonical:
            Freq = (Token.Count / Total) * 100
            if Freq <= OutlierThreshold (default: 10%):
                if not IsException(Token, Context):
                    Emit RogueOutlierDiagnostic
```

### 3.3. Chromatic Family & Decoupled WCAG Contrast
- **Chromatic Family:** Maps `yellow-*`, `amber-*`, `orange-*` to `FamilyWarmAmber` with `CandidateIntent = Warning`.
  Detects hue fragmentation across files and emits recommendations:
  1. Define `--warning` and `--warning-foreground` in `global.css`.
  2. Add `variant="warning"` in component definition.
- **WCAG Contrast:** Computed mathematically from static sRGB values ($L = 0.2126R + 0.7152G + 0.0722B$).
  Identifies absolute accessibility hazards independently of cluster distribution (e.g., `bg-yellow-400 text-white`).

---

## 4. Verification Plan

1. **Unit & Benchmark Tests (`internal/drift/`):**
   - Scope resolution with exact, native, and semantic confidence.
   - Fine-grained category normalization (preserving variants).
   - Statistical clustering math with `MinClusterOccurrences` gate.
   - Chromatic family competition detection and recommendations.
   - Contrast ratio calculation against known sRGB reference values.
   - Allocation and benchmark test (`go test -bench=. -benchmem`).
2. **Scanner & Engine Integration:**
   - Worker pool aggregation without mutex.
   - Seamless diagnostic reporting in `charites scan` (ANSI inline, JSON, Markdown).
   - Dedicated CLI command `charites drift [path] [--report]`.
3. **1-SSOT Tri-Corpus (`tests/correctness/design/component-style-drift/`):**
   - `positive/`: 10x canonical `rounded-md` + 1x rogue `rounded-2xl` (population denominator = 11, rogue = 9.1%), plus chromatic drift and contrast hazard.
   - `negative/`: 12x canonical `rounded-md`, explicit `// charites:ignore`, and `size="icon"` full-radius exception.
   - `adversarial/`: Concentric Card (`rounded-2xl` on `<Card>`), generic `<div>`, and small sample cluster (< 10) passing `InsufficientData` gate.
4. **Gates & Quality:**
   - `go test -v ./tests -run TestGoldenCorpus_AdoptionMatrix`
   - `go test -v ./tests -run TestCorrectnessGate`
   - `make lint` (golangci-lint clean)
   - `make wiki` (generates `wiki/design.md` and `wiki/design.component-style-drift.md`)
   - `go test -race ./...` (0 data races)
