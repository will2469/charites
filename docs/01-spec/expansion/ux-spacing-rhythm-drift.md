# EXPANSION SPECIFICATION: `ux.spacing-rhythm-drift`
> **Kode Dokumen:** `SPEC-EXP-UX-SPACING-RHYTHM-DRIFT`
> **Kategori:** `ux`
> **Pilar:** `01-SPEC` (WHAT - Spesifikasi Perilaku & Kontrak Rule)
> **Status:** Proposed Expansion Specification
> **Standar Rujukan:** W3C DTCG, Gestalt Proximity & Visual Rhythm, Evidence-Based Static Analysis
> **Pilar Terkait:** [01-SPEC: ux.md](ux.md) & [01-SPEC: themes.md](themes.md)

---

## 1. Metadata

| Field                    | Value                                                                                        |
| ------------------------ | -------------------------------------------------------------------------------------------- |
| Rule ID                  | `ux.spacing-rhythm-drift`                                                                    |
| Category                 | `ux.*`                                                                                       |
| Severity                 | `warning`                                                                                    |
| Status                   | Proposed                                                                                     |
| Engine                   | JSX/TSX AST                                                                                  |
| Analysis Type            | Structural + local relational analysis                                                       |
| Primary Goal             | Detect spacing sequences that form an inconsistent local rhythm within the same layout group |
| False Positive Policy    | Conservative; warn only when structural evidence is sufficient                               |
| Auto-fix                 | No                                                                                           |
| Cross-file Analysis      | No                                                                                           |
| Cross-component Analysis | No in v1                                                                                     |

---

## 2. Purpose

`ux.spacing-rhythm-drift` detects **local spacing rhythm inconsistencies** within a single layout group.

The rule does not enforce a universal spacing modulus such as:

* all spacing must be a multiple of `4px`;
* all spacing must be a multiple of `8px`;
* dense interfaces must use multiples of `12px`;
* spacious interfaces must use multiples of `8px`.

Those are design-system conventions rather than universal UI laws.

The rule instead evaluates whether spacing values used for **the same structural relationship** form a coherent local rhythm.

### Core invariant

> Spacing values serving the same spatial role within the same layout group should normally come from a coherent local spacing family.

Example:

```tsx
<div className="flex flex-col gap-3">
  <FieldA />
  <FieldB />
  <FieldC />
  <FieldD />
</div>
```

The repeated relationship is:

```text
FieldA  FieldB
FieldB  FieldC
FieldC  FieldD
```

The spacing model is coherent because all repeated relationships use `gap-3`.

By contrast:

```tsx
<div className="flex flex-col">
  <div className="mb-3"><FieldA /></div>
  <div className="mb-5"><FieldB /></div>
  <div className="mb-3"><FieldC /></div>
</div>
```

contains an inconsistent local rhythm:

```text
3 → 5 → 3
```

provided those margins represent the same sibling-separation relationship.

---

## 3. Non-Goals

This rule MUST NOT attempt to determine whether a spacing value is globally "good", "bad", "modern", "dense", or "spacious" from the numeric value alone.

It MUST NOT:

1. enforce `spacing % 8 == 0`;
2. enforce `spacing % 4 == 0`;
3. enforce `spacing % 12 == 0`;
4. reject Tailwind fractional values solely because they are fractional;
5. replace design-token validation;
6. detect component-internal vs external spacing inversion;
7. infer psychological usability from spacing alone;
8. compare unrelated layout groups;
9. require every sibling in a component to have identical spacing;
10. treat different semantic roles as one spacing sequence.

Ownership is therefore:

```text
theme.hardcode-size
    ↓
token/value validity

ux.spacing-inversion
    ↓
internal spacing vs external group spacing

ux.spacing-rhythm-drift
    ↓
local spacing consistency

ux.spacing-density-mismatch
    ↓
density semantics
```

---

## 4. Design Principle

The rule operates on **relationships**, not isolated scalar values.

A spacing value has limited semantic information by itself.

For example:

```text
gap-4
```

does not reveal whether `4` means:

* label → input,
* field → field,
* card → card,
* section → section,
* page → section.

Therefore the analyzer MUST first determine the **spatial relationship** represented by the value.

The conceptual model is:

```text
Spacing Value
     ↓
Spatial Role
     ↓
Layout Group
     ↓
Peer Relationships
     ↓
Local Rhythm
     ↓
Drift Classification
```

---

## 5. Terminology

### 5.1 Spacing Value

The normalized numeric distance represented by a spacing declaration.

Examples:

```tsx
gap-2
gap-3
gap-4
mt-6
mb-8
```

may normalize to:

```text
2
3
4
6
8
```

The exact normalization unit is implementation-defined but MUST remain consistent within one analysis.

---

### 5.2 Spacing Role

The structural relationship represented by the spacing.

```go
type SpacingRole int

const (
    SpacingRoleUnknown SpacingRole = iota
    SpacingRoleInline
    SpacingRoleComponent
    SpacingRoleGroup
    SpacingRoleSection
)
```

Suggested interpretation:

| Role        | Example           |
| ----------- | ----------------- |
| `Inline`    | icon  text       |
| `Component` | label  input     |
| `Group`     | field  field     |
| `Section`   | section  section |

v1 SHOULD primarily analyze `Component` and `Group` relationships.

`Inline` and `Section` relationships SHOULD require stronger structural evidence.

---

### 5.3 Layout Group

A set of sibling relationships governed by a common layout mechanism.

Examples:

```tsx
<div className="flex flex-col gap-4">
  ...
</div>
```

```tsx
<div className="grid grid-cols-2 gap-6">
  ...
</div>
```

```tsx
<div className="space-y-4">
  ...
</div>
```

A layout group is identified from:

* flex layout;
* grid layout;
* `space-*`;
* repeated margin-based sibling spacing;
* explicit layout primitives supported by the analyzer.

---

### 5.4 Rhythm

A sequence of comparable spacing relationships within the same layout group.

Example:

```text
4, 4, 4, 4
```

is a uniform rhythm.

Example:

```text
3, 4, 6, 8
```

may represent a valid hierarchical scale if the relationships are semantically distinct.

Therefore the rule MUST NOT treat every unequal sequence as drift.

---

### 5.5 Rhythm Drift

A local inconsistency where comparable spatial relationships use spacing values that do not belong to the same inferred spacing pattern.

Example:

```text
4, 4, 7, 4
```

inside repeated field relationships is a strong candidate for drift.

---

## 6. Input Scope

The rule operates on JSX/TSX AST or the repository's normalized JSX IR.

In v1 the analyzer SHOULD support:

### Layout mechanisms

```text
flex
grid
space-x-*
space-y-*
gap-*
gap-x-*
gap-y-*
margin-top
margin-bottom
margin-left
margin-right
```

through the project's existing Tailwind/class extraction infrastructure.

### Supported static values

The rule SHOULD analyze:

```tsx
className="gap-4"
className="space-y-4"
className="mb-6"
className="mt-3"
```

and equivalent statically-resolvable forms:

```tsx
className={"gap-4"}
className={`gap-4`}
```

only when the existing class extraction layer can resolve them deterministically.

Dynamic constructions such as:

```tsx
className={`gap-${size}`}
```

are `Unknown` for this rule unless the existing IR resolves them.

---

## 7. Required Structural Facts

The rule SHOULD consume normalized IR facts rather than repeatedly parsing raw JSX.

Suggested facts:

```go
type SpacingOccurrence struct {
    NodeID        NodeID
    ParentID      NodeID

    Role          SpacingRole
    Axis          Axis
    Value         float64

    Source        SpacingSource
    Confidence    Confidence

    LayoutGroupID NodeID
}
```

Where:

```go
type Axis int

const (
    AxisUnknown Axis = iota
    AxisHorizontal
    AxisVertical
)
```

and:

```go
type SpacingSource int

const (
    SpacingSourceUnknown SpacingSource = iota
    SpacingSourceGap
    SpacingSourceSpace
    SpacingSourceMargin
)
```

The exact names MAY follow existing project conventions.

---

## 8. Classification Pipeline

The analyzer MUST follow a deterministic pipeline.

```text
AST / IR
   ↓
Layout Node Detection
   ↓
Spacing Extraction
   ↓
Layout Group Construction
   ↓
Spatial Role Classification
   ↓
Peer Relationship Grouping
   ↓
Rhythm Model Construction
   ↓
Drift Classification
   ↓
Diagnostic
```

No diagnostic logic should directly inspect arbitrary class strings without going through normalized spacing facts.

---

## 9. Layout Node Detection

A node is a layout candidate when its styling establishes a spatial relationship among children.

Examples:

```tsx
<div className="flex flex-col gap-4">
```

```tsx
<div className="grid gap-6">
```

```tsx
<div className="space-y-4">
```

The layout classifier SHOULD identify:

```go
type LayoutKind int

const (
    LayoutUnknown LayoutKind = iota
    LayoutFlex
    LayoutGrid
    LayoutStack
    LayoutOther
)
```

The analyzer MUST preserve axis information.

Example:

```text
flex-col + gap-4
    → vertical sibling spacing

flex-row + gap-4
    → horizontal sibling spacing
```

---

## 10. Spacing Extraction

Spacing extraction converts supported syntax into normalized spacing occurrences.

Examples:

```text
gap-4
→ value=4
→ source=Gap
→ axis=Unknown until layout direction is resolved
```

```text
gap-y-4
→ value=4
→ axis=Vertical
```

```text
space-y-4
→ value=4
→ source=Space
→ axis=Vertical
```

```text
mb-4
→ value=4
→ source=Margin
→ axis=Vertical
```

---

## 11. Relationship Classification

Spacing values SHOULD only be compared when they represent comparable relationships.

Example:

```tsx
<div className="flex flex-col gap-3">
  <Field />
  <Field />
  <Field />
</div>
```

The relationships are equivalent:

```text
Field  Field
Field  Field
```

Therefore they form one rhythm group.

However:

```tsx
<Card className="p-6">
  <Field />
</Card>
```

and:

```tsx
<Card className="mb-8">
  ...
</Card>
```

do not represent the same relationship.

They MUST NOT be merged into one rhythm sequence.

---

## 12. Evidence Hierarchy

The analyzer SHOULD prefer structural evidence over lexical assumptions.

Recommended ranking:

```text
VERY_STRONG
    repeated same-layout relationship
    repeated sibling spacing primitive

STRONG
    common parent layout
    same axis
    same spacing source
    same semantic role

MEDIUM
    same component subtree
    repeated structural pattern

WEAK
    class-name similarity
    arbitrary lexical interpretation

UNKNOWN
    dynamic class
    unresolved conditional
    ambiguous relationship
```

Diagnostics MUST require at least `STRONG` evidence.

Weak evidence MUST NOT independently trigger a warning.

---

## 13. Rhythm Model

A rhythm model represents comparable spacing occurrences.

Suggested structure:

```go
type RhythmGroup struct {
    LayoutGroupID NodeID
    Role          SpacingRole
    Axis          Axis

    Occurrences   []SpacingOccurrence
}
```

For example:

```tsx
<div className="flex flex-col gap-4">
    ...
</div>
```

produces:

```text
RhythmGroup
  layout = flex
  axis   = vertical
  role   = group
  values = [4]
```

Repeated explicit margins may produce:

```text
RhythmGroup
  role   = group
  axis   = vertical
  values = [4, 4, 6, 4]
```

---

## 14. Drift Detection

The rule MUST distinguish between:

1. uniform repetition;
2. coherent scale progression;
3. isolated outlier;
4. insufficient evidence.

### 14.1 Uniform Rhythm

```text
4, 4, 4, 4
```

Result:

```text
No diagnostic
```

---

### 14.2 Coherent Token Family

```text
2, 4, 6, 8
```

does not automatically indicate drift.

The analyzer MUST avoid treating ascending spacing as an error when each value corresponds to a distinguishable structural relationship.

---

### 14.3 Isolated Outlier

Example:

```text
4, 4, 7, 4
```

If the `7` occurrence belongs to exactly the same peer relationship as the surrounding `4`s:

```text
4
4
7  ← outlier
4
```

then:

```text
Diagnostic: yes
```

---

### 14.4 Insufficient Evidence

Example:

```text
gap-4
```

or:

```text
gap-4 + one unrelated mb-7
```

must produce:

```text
No diagnostic
```

because there is no sufficiently large comparable rhythm.

---

## 15. Minimum Sample Requirement

The rule SHOULD avoid warnings from tiny samples.

Recommended v1 threshold:

```text
minimum comparable occurrences = 3
```

For example:

```text
4, 4, 7
```

contains enough evidence for one possible outlier.

But:

```text
4, 7
```

does not.

Configurable form:

```go
type RhythmConfig struct {
    MinOccurrences int
    OutlierRatio   float64
}
```

Suggested defaults:

```text
MinOccurrences = 3
OutlierRatio   = 0.50
```

The exact threshold SHOULD remain centralized rather than embedded in individual detectors.

---

## 16. Outlier Classification

An occurrence is a candidate outlier when:

1. it belongs to a rhythm group with sufficient observations;
2. its spatial role matches the majority group;
3. its axis matches the majority group;
4. its value differs from the dominant local value;
5. no structural evidence explains the difference;
6. the difference exceeds the configured outlier threshold.

Example:

```text
4, 4, 4, 7, 4
```

Dominant value:

```text
4
```

Candidate outlier:

```text
7
```

Result:

```text
warning
```

---

## 17. Majority-Based Detection

The initial implementation SHOULD prefer a simple deterministic majority model over statistical clustering.

For:

```text
[4, 4, 4, 7, 4]
```

frequency:

```text
4 → 4
7 → 1
```

Dominant rhythm:

```text
4
```

Outlier:

```text
7
```

A more advanced clustering algorithm SHOULD NOT be introduced until repository evidence demonstrates that simple majority detection creates material false positives.

KISS is the preferred v1 implementation.

---

## 18. Semantic Exceptions

The analyzer MUST avoid warnings when different spacing values are structurally intentional.

Examples:

```tsx
<div className="flex flex-col gap-4">
  <ProfileHeader className="mb-8" />
  <Field />
  <Field />
  <Footer className="mt-12" />
</div>
```

Here:

```text
header → content
field → field
content → footer
```

are not necessarily the same spatial relationship.

The analyzer MUST classify them separately rather than flattening:

```text
8, 4, 4, 12
```

into one rhythm.

---

## 19. Dominant Relationship Rule

When multiple spacing roles exist, the analyzer SHOULD build independent rhythm groups.

Example:

```text
Component spacing:
    3, 3, 3

Group spacing:
    6, 6, 6

Section spacing:
    10, 10
```

This is coherent.

The fact that:

```text
3 ≠ 6 ≠ 10
```

is not itself a problem.

The rule is about **drift inside a relationship class**, not global numerical uniformity.

---

## 20. Explicit Token Awareness

The rule MAY reuse normalized spacing tokens from the theme engine.

Example:

```text
space-3
space-4
space-6
space-8
```

should remain distinguishable from arbitrary values:

```text
space-[13px]
space-[17px]
```

However, token validity belongs to:

```text
theme.hardcode-size
```

Therefore `ux.spacing-rhythm-drift` MUST NOT emit a duplicate diagnostic merely because a value is non-standard.

Example:

```tsx
<div className="gap-[13px]">
```

should be owned by:

```text
theme.hardcode-size
```

not automatically by:

```text
ux.spacing-rhythm-drift
```

---

## 21. Interaction With `ux.spacing-inversion`

The two rules have different ownership boundaries.

### `ux.spacing-inversion`

Detects:

```text
internal child spacing
    >=
external parent sibling spacing
```

Example:

```tsx
<section className="space-y-3">
  <div className="flex flex-col gap-8">
```

Primary concern:

```text
hierarchical spacing relationship
```

### `ux.spacing-rhythm-drift`

Detects:

```text
same relationship
    +
inconsistent local spacing
```

Example:

```tsx
<div className="space-y-4">
  <Field className="mb-4" />
  <Field className="mb-7" />
  <Field className="mb-4" />
</div>
```

Primary concern:

```text
local rhythm consistency
```

A node MUST NOT receive two diagnostics for the same underlying invariant unless the violations are independently meaningful.

---

## 22. Density Independence

`ux.spacing-rhythm-drift` MUST NOT infer:

```text
gap-4 = spacious
gap-3 = comfortable
gap-2 = dense
```

as absolute truth.

Density classification belongs to:

```text
ux.spacing-density-mismatch
```

This preserves separation:

```text
rhythm
≠
density
```

A dense table can have a perfectly consistent rhythm:

```text
2, 2, 2, 2
```

and a spacious landing page can also have a perfectly consistent rhythm:

```text
8, 8, 8, 8
```

Both are valid.

---

## 23. AST Constraints

v1 MUST remain local.

The analyzer MAY inspect:

```text
current node
parent node
direct children
normalized class/style facts
```

The analyzer MUST NOT perform unrestricted recursive tree walking from a rule detector.

Recommended architecture:

```text
AST
 ↓
IR extraction
 ↓
SpacingOccurrence
 ↓
LayoutGroup
 ↓
RhythmGroup
 ↓
Detector
```

This keeps rule logic deterministic and testable.

---

## 24. Dynamic Values

Dynamic values MUST be conservative.

Example:

```tsx
<div className={`gap-${spacing}`}>
```

If `spacing` cannot be resolved statically:

```text
SpacingValue = Unknown
```

The occurrence SHOULD be excluded from numeric rhythm comparison.

Example:

```tsx
<div className="flex flex-col gap-4">
  <A />
</div>

<div className={`flex flex-col gap-${spacing}`}>
  <B />
</div>
```

The analyzer MUST NOT infer:

```text
4 vs unknown → drift
```

because there is insufficient evidence.

---

## 25. Conditional Spacing

Example:

```tsx
<div className={large ? "gap-8" : "gap-4"}>
```

If both values are statically recoverable, the analyzer MAY represent:

```text
candidate values = {4, 8}
```

but MUST NOT treat the conditional branch itself as drift.

The rule should instead ask whether another independent peer occurrence breaks the same rhythm.

Example:

```tsx
<div className="flex flex-col">
  <A className={large ? "mb-4" : "mb-8"} />
  <B className="mb-7" />
  <C className="mb-4" />
</div>
```

The dynamic branch creates uncertainty and SHOULD lower confidence rather than force a warning.

---

## 26. Responsive Behavior

Responsive spacing MUST be analyzed per effective breakpoint state.

Example:

```tsx
<div className="gap-4 md:gap-6">
```

represents:

```text
base:
    gap = 4

md:
    gap = 6
```

The analyzer MUST NOT flatten:

```text
4, 6
```

into one rhythm and report drift.

Instead:

```text
base rhythm
md rhythm
```

are separate states.

A responsive cascade MUST be resolved by stylesheet semantics, not by the textual order of classes.

---

## 27. Diagnostic Ownership

Diagnostic responsibility belongs to the rhythm group that contains the outlier.

Example:

```tsx
<div className="flex flex-col">
  <Field className="mb-4" />
  <Field className="mb-4" />
  <Field className="mb-7" />
  <Field className="mb-4" />
</div>
```

Diagnostic target:

```text
Field with mb-7
```

not:

```text
parent div
```

unless the spacing declaration physically belongs to the parent.

---

## 28. Diagnostic Message

Recommended message:

```text
Spacing rhythm drifts from the surrounding sibling group: `mb-7` is inconsistent with the local `mb-4` pattern.
```

Optional explanation:

```text
The surrounding peer relationships consistently use `4`, while this relationship uses `7`. Consider using the established local spacing token when no semantic distinction is intended.
```

The message MUST NOT claim:

```text
7px is wrong
```

or:

```text
spacing must be a multiple of 8
```

because the rule does not establish either claim.

---

## 29. Suggested Diagnostic Structure

```go
type RhythmDriftDiagnostic struct {
    NodeID              NodeID
    ActualValue         float64
    DominantValue       float64

    Role                SpacingRole
    Axis                Axis

    GroupSize           int
    DominantOccurrences int

    Source              SpacingSource
}
```

This allows diagnostics and tests to validate classification without parsing message strings.

---

## 30. Decision Algorithm

Recommended v1 algorithm:

```text
1. Extract static spacing occurrences.
2. Resolve effective axis.
3. Resolve layout group.
4. Classify spatial role.
5. Partition occurrences by:
      layout group
      axis
      spatial role
6. Reject groups with fewer than MinOccurrences.
7. Count spacing values.
8. Identify dominant local value.
9. Identify minority values.
10. Reject values with insufficient confidence.
11. Reject values explained by distinct semantic relationships.
12. Flag isolated outlier(s).
```

Pseudo-code:

```go
func DetectRhythmDrift(group RhythmGroup, cfg RhythmConfig) []Diagnostic {
    if len(group.Occurrences) < cfg.MinOccurrences {
        return nil
    }

    dominant, count := dominantValue(group.Occurrences)

    if count < 2 {
        return nil
    }

    diagnostics := make([]Diagnostic, 0)

    for _, occurrence := range group.Occurrences {
        if occurrence.Value == dominant {
            continue
        }

        if !isComparable(occurrence, group) {
            continue
        }

        if !exceedsOutlierThreshold(
            occurrence.Value,
            dominant,
            cfg.OutlierRatio,
        ) {
            continue
        }

        diagnostics = append(diagnostics, buildDiagnostic(
            occurrence,
            dominant,
            count,
        ))
    }

    return diagnostics
}
```

The implementation MAY evolve beyond this algorithm, but the semantic ownership MUST remain unchanged.

---

## 31. False Positive Guards

The rule MUST NOT report when:

### Different semantic roles

```text
label → input = 3
field → field = 6
section → field-group = 8
```

### Insufficient sample size

```text
4, 7
```

### Dynamic values

```text
gap-${spacing}
```

### Responsive states

```text
base: 4
md: 6
```

### Different layout groups

```tsx
<Card className="gap-4" />
<Modal className="gap-8" />
```

### Intentional structural boundaries

```text
header → content
content → footer
```

### Existing component variants

When a component explicitly declares separate variants:

```tsx
<Stack gap="sm" />
<Stack gap="lg" />
```

and the variant semantics are known to the analyzer, those occurrences SHOULD remain separate.

---

## 32. Positive Detection Cases

### Case A: Repeated margin outlier

```tsx
<div>
  <Field className="mb-4" />
  <Field className="mb-4" />
  <Field className="mb-7" />
  <Field className="mb-4" />
</div>
```

Expected:

```text
1 diagnostic
```

---

### Case B: Repeated group spacing

```tsx
<div className="flex flex-col">
  <Group className="mb-4" />
  <Group className="mb-4" />
  <Group className="mb-6" />
  <Group className="mb-4" />
</div>
```

Expected:

```text
1 diagnostic
```

assuming all four margins represent the same sibling relationship.

---

## 33. Negative Detection Cases

### Case A: Uniform rhythm

```tsx
<div>
  <Field className="mb-4" />
  <Field className="mb-4" />
  <Field className="mb-4" />
</div>
```

Expected:

```text
0 diagnostics
```

---

### Case B: Hierarchical spacing

```tsx
<section className="space-y-8">
  <div className="flex flex-col gap-3">
    <Label />
    <Input />
  </div>

  <div className="flex flex-col gap-3">
    <Label />
    <Input />
  </div>
</section>
```

Expected:

```text
0 diagnostics
```

The `3` values and `8` values have different roles.

---

### Case C: Valid token progression

```tsx
<div className="flex flex-col">
  <Header className="mb-8" />
  <Content />
  <Footer className="mt-12" />
</div>
```

Expected:

```text
0 diagnostics
```

These represent different structural boundaries.

---

### Case D: Dynamic spacing

```tsx
<div className={`flex flex-col gap-${spacing}`}>
```

Expected:

```text
0 diagnostics
```

unless `spacing` is statically resolved elsewhere in the IR.

---

### Case E: Responsive spacing

```tsx
<div className="gap-4 md:gap-6">
```

Expected:

```text
0 diagnostics
```

because the values belong to different responsive states.

---

## 34. Tri-Corpus Fixture Requirements

The rule MUST use a three-corpus fixture model:

```text
correct/
wrong/
edge/
```

### Correct corpus

Must include:

```text
uniform spacing
hierarchical spacing
different semantic roles
responsive variants
different layout groups
dense layouts
spacious layouts
valid token progression
```

### Wrong corpus

Must include:

```text
isolated local outlier
repeated margin drift
repeated gap drift
mixed local rhythm
same role + same axis + one anomalous value
```

### Edge corpus

Must include:

```text
dynamic className
conditional className
responsive className
insufficient sample size
nested groups
mixed spacing roles
mixed axes
multiple outliers
unknown layout
unknown spacing value
```

---

## 35. Required Test Matrix

| Scenario                                    | Expected |
| ------------------------------------------- | :------: |
| `4,4,4,4`                                   |    0     |
| `4,4,7,4`                                   |    1     |
| `3,3,3,6`                                   | 1 when all roles are equivalent |
| `3,6`                                       |    0     |
| `3,3,6,8` with distinct semantic roles      |    0     |
| `gap-4` vs `gap-6` at different breakpoints |    0     |
| dynamic spacing                             |    0     |
| different parents                           |    0     |
| different axes                              |    0     |
| different semantic roles                    |    0     |
| one known token outlier in repeated group   |    1     |
| two equally common values                   |    0     |
| unresolved layout relationship              |    0     |

---

## 36. Multiple Outliers

The rule MAY report multiple outliers when there is a clearly dominant pattern.

Example:

```text
4, 4, 7, 9, 4
```

Expected:

```text
dominant = 4
outliers = 7, 9
```

However:

```text
4, 4, 6, 6
```

must not automatically report `6`.

There is no sufficiently strong dominant rhythm.

This prevents the detector from forcing an arbitrary winner.

---

## 37. Tie Handling

For:

```text
4, 4, 6, 6
```

the detector MUST return:

```text
No diagnostic
```

unless additional structural evidence establishes that one value belongs to the same relationship class and the other does not.

The rule MUST NOT resolve ties arbitrarily.

---

## 38. Repository-Level Design Tokens

Repository token knowledge MAY improve confidence.

Example:

```css
:root {
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
}
```

A repeated local pattern:

```text
8, 8, 8, 12, 8
```

can still be flagged even though both `8` and `12` are valid design tokens.

Reason:

```text
token validity
≠
local rhythm consistency
```

This distinction is fundamental to the rule.

---

## 39. Configuration

Recommended v1:

```go
type RhythmConfig struct {
    MinOccurrences int
    MinDominance   float64
}
```

Suggested defaults:

```text
MinOccurrences = 3
MinDominance   = 0.60
```

Meaning:

```text
dominant occurrences / total occurrences >= 0.60
```

before an outlier may be considered.

Configuration MUST NOT contain:

```text
BaseSpacing = 8
DenseSpacing = 12
```

because that would convert a repository-independent structural rule into an arbitrary design-system rule.

---

## 40. Performance Constraints

The rule MUST remain linear with respect to inspected local occurrences.

Target complexity:

```text
O(n)
```

for a layout group.

The rule MUST NOT:

* perform repository-wide AST scans;
* repeatedly reparse class strings;
* recursively inspect unrelated subtrees;
* perform inter-file analysis;
* invoke network operations;
* depend on runtime rendering.

---

## 41. Ownership Boundary With `theme.hardcode-size`

Example:

```tsx
<div className="gap-[17px]">
```

Primary owner:

```text
theme.hardcode-size
```

Example:

```tsx
<div>
  <Field className="mb-4" />
  <Field className="mb-4" />
  <Field className="mb-7" />
</div>
```

Primary owner:

```text
ux.spacing-rhythm-drift
```

Even when `7` is also a non-standard spacing value, the rhythm rule SHOULD avoid duplicate messaging when `theme.hardcode-size` already owns the scalar validity violation.

A shared diagnostic deduplication layer SHOULD be used if the engine supports one.

---

## 42. Severity

Recommended severity:

```text
warning
```

Rationale:

Spacing rhythm inconsistencies are design-quality defects, not generally correctness, security, or accessibility failures.

The rule should therefore remain advisory.

---

## 43. Auto-Fix Policy

No auto-fix in v1.

The analyzer does not know with sufficient certainty whether:

```text
4 → 6
```

or:

```text
6 → 4
```

is the intended design.

The correct replacement requires design intent that cannot reliably be inferred from AST alone.

Therefore:

```text
detect = yes
suggest = optional
autofix = no
```

---

## 44. Suggested Future Extension

A later repository-aware analyzer could construct a canonical rhythm model:

```text
Component
   ↓
Layout Families
   ↓
Observed Spacing Relations
   ↓
Dominant Local Tokens
   ↓
Outlier Detection
```

For example:

```text
FormField
  label → input = 12
  input → help  = 4
  field → field  = 16
```

Then another `FormField` with:

```text
label → input = 18
```

could be detected as a component-family drift.

This MUST remain outside the v1 local AST rule.

---

## 45. Future Density Integration

`ux.spacing-density-mismatch` MAY later consume the same normalized facts:

```text
SpacingOccurrence
        ↓
Density Classifier
        ↓
Dense / Comfortable / Spacious
```

Possible future evidence:

```text
table
data-grid
form
toolbar
card
hero
section
navigation
```

But `ux.spacing-rhythm-drift` MUST remain independent.

Shared infrastructure is preferred:

```text
internal/rules/layout
    ├── spacing extraction
    ├── layout classification
    ├── axis resolution
    ├── relationship classification
    └── rhythm grouping
```

Rule-specific policy remains under:

```text
internal/rules/ux
```

---

## 46. Proposed Package Structure

```text
internal/rules/
├── layout/
│   ├── spacing.go
│   ├── layout.go
│   ├── relationship.go
│   └── rhythm.go
│
└── ux/
    └── spacing_rhythm_drift.go
```

Possible shared types:

```go
package layout

type SpacingOccurrence struct {
    ...
}

type RhythmGroup struct {
    ...
}

func ExtractSpacing(...)
func ClassifySpacingRole(...)
func BuildRhythmGroups(...)
```

Rule package:

```go
package ux

func EvaluateSpacingRhythmDrift(...)
```

The rule SHOULD consume normalized facts and MUST NOT duplicate spacing parsing implemented elsewhere.

---

## 47. Acceptance Criteria

Implementation is considered complete when all conditions below are satisfied.

### Functional

* detects repeated local spacing outliers;
* distinguishes spatial roles;
* distinguishes axes;
* distinguishes layout groups;
* handles static Tailwind spacing;
* handles supported margin/space/gap forms;
* does not enforce a universal `4/8/12` modulus;
* handles responsive spacing per breakpoint;
* rejects insufficient evidence;
* rejects ties;
* rejects unresolved dynamic values;
* avoids duplicate ownership with `theme.hardcode-size`;
* avoids overlap with `ux.spacing-inversion`.

### Architectural

* no hidden recursive tree walk inside the rule;
* no raw class-string classification in the diagnostic layer;
* normalized spacing facts are reusable;
* classification occurs before diagnostic generation;
* unknown values are never treated as valid evidence;
* rule remains deterministic.

### Testing

* tri-corpus fixtures exist;
* all required positive cases have one diagnostic;
* all required negative/edge cases have zero diagnostics;
* diagnostic target points to the actual outlier;
* regression coverage exists for responsive and dynamic classes.

---

## 48. Reference Invariant

The complete rule can be summarized as:

```text
Same layout group
+
Same spatial role
+
Same axis
+
Sufficient comparable observations
+
Dominant local rhythm
+
Isolated unexplained deviation
=
ux.spacing-rhythm-drift
```

Conversely:

```text
Different role
OR
different layout group
OR
different responsive state
OR
insufficient evidence
OR
ambiguous classification
=
no diagnostic
```

---

## 49. Final Rule Contract

`ux.spacing-rhythm-drift` is a **local relational consistency rule**.

It does not answer:

> "Is `8px` better than `12px`?"

It answers:

> "Given this specific structural relationship, does one spacing declaration unexpectedly break the rhythm established by its comparable peers?"

That distinction is the core contract of the rule.

The rule therefore treats:

```text
8
12
16
```

as values without inherent semantic correctness.

The semantics emerge from:

```text
node relationship
+
layout structure
+
peer consistency
+
local context
```

This keeps the analyzer AST-deterministic while allowing Charites to detect a class of design defects that ordinary token linters cannot see.
