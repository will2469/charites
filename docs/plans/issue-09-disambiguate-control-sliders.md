# Implementation Plan: Disambiguate Control/Verification Sliders from Content Carousels (`cls.unconstrained-carousel`)

Resolves **[GitHub Issue #9](https://github.com/will2469/charites/issues/9)** (`fix(cls): disambiguate control/verification sliders from content carousels in cls.unconstrained-carousel`).

---

## 1. Executive Summary & Problem Formulation

In the static analyzer rule `cls.unconstrained-carousel`, the matcher in `internal/rules/cls/util.go` previously classified any JSX element whose tag name ended with `Slider` (or matched `/[Ss]lider/`) as a **content carousel track**:
```go
switch node.Tag {
case "Carousel", "Slider", "Swiper", "EmblaCarousel":
    return true
}
if strings.HasSuffix(node.Tag, "Carousel") || strings.HasSuffix(node.Tag, "Slider") {
    return true
}
```

This naive lexical match produced false-positive warnings on interactive control and security verification sliders:
```tsx
// Specimen from apps/kosmos/src/features/auth/relaxed-auth/components/QRFallbackLogin.tsx
<ChallengeSlider
    onSolve={cerberus.solveChallenge}
    onCancel={onCancel}
    error={cerberus.error}
    solved={cerberus.challengeSolved}
/>
```

The scanner produced false alarm:
> `Carousel/slider container <ChallengeSlider> lacks bounded vertical dimensions ('h-*', 'min-h-*', or 'aspect-*'). Dynamically loaded slides or slide transitions may cause vertical layout jumps.`

Furthermore, the existing `Evaluate()` routine contained a blanket bypass on spread attributes:
```go
if hasSpreadProps(node.Attributes) {
    return nil
}
```
This caused real unconstrained carousels like `<Carousel {...props} />` to silently escape inspection.

---

## 2. The 6 Locked Architectural Invariants & Evidence Taxonomy

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                    THE 6 LOCKED ARCHITECTURAL INVARIANTS                    │
├─────────────────────────────────────────────────────────────────────────────┤
│ 1. Generic Range Props Never Hijack Dedicated Carousels:                    │
│    'min', 'max', 'value', 'defaultValue' are generic numeric state props    │
│    frequently used on carousels (e.g. min={0} max={10} value={activeSlide}).│
│    They are NEVER sufficient alone or in raw numeric combination to trigger │
│    SliderControl without control-specific semantics or control tag hints.   │
│                                                                             │
│ 2. Evidence-Based Tri-State Classification (SliderClassification):          │
│    - SliderControl: skip CLS check (0 warnings).                            │
│    - SliderContentCarousel: run CLS bounded-height analysis.                │
│    - SliderUnknown: conservative suppression (0 warnings).                  │
│                                                                             │
│ 3. Structural vs Lexical Carousel Decoupling:                               │
│    - CarouselStructural: concrete DOM/CSS evidence only                     │
│      (horizontal scroll-snap 'overflow-x-*' + 'snap-*', <Slide> children). │
│      Never claims IR '.map' expression facts (not modeled in IR).           │
│    - CarouselLexical: naming hints (BannerSlider, ImageSlider).             │
│    - CarouselDedicated: canonical carousel tags (Carousel, Swiper, Embla).  │
│                                                                             │
│ 4. Control Tags are Lexical Hints, Not Authoritative Proof:                 │
│    Tag names (RangeSlider, VolumeSlider) require supporting props.          │
│    Standalone <RangeSlider /> with no props -> SliderUnknown (0 warnings).  │
│    <RangeSlider min={0} max={100} /> -> SliderControl (0 warnings).         │
│                                                                             │
│ 5. Conservative Contradiction & Ambiguity Resolution:                       │
│    - Very Strong Control + Structural Carousel -> SliderUnknown (0 warnings)│
│    - Ambiguous standalone <Slider /> without props -> SliderUnknown (0).    │
│    - Absence of evidence != evidence of carousel.                           │
│                                                                             │
│ 6. Spread Props are Unknown Attributes, Not Suppressions:                   │
│    Remove blanket 'hasSpreadProps() => nil'. Unconstrained <Carousel {...props}>│
│    remains a Carousel and is properly inspected for bounded height.         │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Evidence Taxonomy Matrix

```text
Control Evidence
├── VERY STRONG (Semantic / Security Invariants)
│   ├── role="slider"
│   ├── aria-valuenow / aria-valuemin / aria-valuemax
│   └── onSolve / onVerify (security challenge step-up)
├── STRONG (Control Events & Paired Hints)
│   ├── onValueChange / onValueCommit (Radix / modern slider events)
│   ├── solved (verification challenge state)
│   └── ControlTagHint + RangeCluster (min + max)
├── MEDIUM (Numeric Range Clusters & Thresholds)
│   ├── min + max + step
│   ├── min + max + (value || defaultValue)
│   └── threshold
└── WEAK (Generic State Props - NEVER Conclusive Alone)
    ├── min alone, max alone
    ├── value alone
    └── onChange alone

Carousel Evidence
├── STRUCTURAL (Concrete DOM / CSS Evidence)
│   ├── Horizontal scroll snap: (overflow-x-auto || overflow-x-scroll) && snap-*
│   └── Slide child elements: <Slide>, <CarouselItem>, <SwiperSlide>, <CarouselSlide>
├── DEDICATED (Canonical Carousel Tag)
│   └── Carousel, Swiper, EmblaCarousel
└── LEXICAL (Content Slider Naming Hints)
    └── Suffix Carousel, BannerSlider, ImageSlider, MediaSlider, HeroSlider, ProductSlider
```

### Precedence Resolution Algorithm

```text
1. Discard non-candidates:
   If node has no slider/carousel tags, no role="slider", and no scroll-snap classes -> SliderUnknown.

2. Extract ControlEvidence & CarouselEvidence:
   - VeryStrongControl: role="slider", aria-valuenow, aria-valuemin, aria-valuemax, onSolve, onVerify
   - StrongControl: onValueChange, onValueCommit, solved, (ControlTagHint && hasMin && hasMax)
   - StructuralCarousel: (hasOverflowX && hasSnapX) || hasSlideChildElement
   - DedicatedCarousel: tag in {"Carousel", "Swiper", "EmblaCarousel"}
   - LexicalContentSlider: tag ends with "Carousel" || tag in {"BannerSlider", "ImageSlider", ...}

3. Precedence:
   a. VeryStrongControl && StructuralCarousel:
      -> SliderUnknown (Conservative suppression, zero false alarms)
   b. VeryStrongControl || StrongControl:
      -> SliderControl
   c. StructuralCarousel:
      -> SliderContentCarousel
   d. DedicatedCarousel:
      -> SliderContentCarousel (Generic min/max/value will NOT demote it!)
   e. LexicalContentSlider:
      -> SliderContentCarousel
   f. Otherwise:
      -> SliderUnknown (Absence of evidence != Carousel; standalone <Slider /> returns Unknown)
```
