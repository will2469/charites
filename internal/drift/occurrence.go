package drift

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ComponentKind mengidentifikasi jenis komponen UI yang dikenali oleh drift analyzer.
type ComponentKind int

// Daftar varian ComponentKind yang didukung.
const (
	// KindUnknown menandai komponen yang tidak teridentifikasi.
	KindUnknown ComponentKind = iota
	KindButton                // <Button>, <button>, [role="button"], <TabsTrigger>
	KindInput                 // <Input>, <textarea>, <select>, <SelectTrigger>, [role="textbox"]
	KindCard                  // <Card>
	KindDialog                // <DialogContent>, [role="dialog"]
	KindSheet                 // <SheetContent>
	KindBadge                 // <Badge>, [role="status"]
	KindMenu                  // PopoverContent, DropdownMenuContent, TooltipContent
	KindAlert                 // <Alert>, [role="alert"]
)

func (k ComponentKind) String() string {
	switch k {
	case KindButton:
		return "Button"
	case KindInput:
		return "Input"
	case KindCard:
		return "Card"
	case KindDialog:
		return "Dialog"
	case KindSheet:
		return "Sheet"
	case KindBadge:
		return "Badge"
	case KindMenu:
		return "Menu"
	case KindAlert:
		return "Alert"
	default:
		return "Unknown"
	}
}

// ScopeConfidence menyatakan tingkat kepastian semantik dari pengenalan elemen AST.
type ScopeConfidence int

// Daftar tingkat kepastian ScopeConfidence.
const (
	// ConfidenceExact menandai komponen bernama persis (misal: <Button>, <Card>).
	ConfidenceExact    ScopeConfidence = iota
	ConfidenceNative                   // Tag HTML natif (misal: <button>, <textarea>)
	ConfidenceSemantic                 // Atribut ARIA semantik (misal: role="button", role="dialog")
)

func (c ScopeConfidence) String() string {
	switch c {
	case ConfidenceExact:
		return "Exact"
	case ConfidenceNative:
		return "Native"
	case ConfidenceSemantic:
		return "Semantic"
	default:
		return "Unknown"
	}
}

// ScopeKey merepresentasikan identitas gabungan jenis komponen dan tingkat kepastiannya.
type ScopeKey struct {
	Kind       ComponentKind
	Confidence ScopeConfidence
}

func (s ScopeKey) String() string {
	return fmt.Sprintf("%s (%s)", s.Kind.String(), s.Confidence.String())
}

// PropertyCategory mendefinisikan kategori properti mikro dan geometris yang dinormalisasi.
type PropertyCategory int

// Daftar kategori properti styling PropertyCategory.
const (
	// CatUnknown menandai kategori properti yang tidak dikenali.
	CatUnknown PropertyCategory = iota
	CatRounded
	CatActiveScale
	CatDisabledOpacity
	CatDisabledCursor
	CatFocusRingWidth
	CatFocusRingColor
	CatFocusRingOffset
	CatShadow
	CatBorderWidth
	CatBorderColor
	CatBackground
	CatHeight
	CatPaddingX
	CatPaddingY
)

func (p PropertyCategory) String() string {
	switch p {
	case CatRounded:
		return "rounded"
	case CatActiveScale:
		return "active:scale"
	case CatDisabledOpacity:
		return "disabled:opacity"
	case CatDisabledCursor:
		return "disabled:cursor"
	case CatFocusRingWidth:
		return "focus-visible:ring-width"
	case CatFocusRingColor:
		return "focus-visible:ring-color"
	case CatFocusRingOffset:
		return "focus-visible:ring-offset"
	case CatShadow:
		return "shadow"
	case CatBorderWidth:
		return "border-width"
	case CatBorderColor:
		return "border-color"
	case CatBackground:
		return "background"
	case CatHeight:
		return "height"
	case CatPaddingX:
		return "padding-x"
	case CatPaddingY:
		return "padding-y"
	default:
		return "unknown"
	}
}

// ClusterKey mendefinisikan kunci unik pengelompokan populasi gaya di seluruh repositori.
type ClusterKey struct {
	Scope    ScopeKey
	Category PropertyCategory
}

func (ck ClusterKey) String() string {
	return fmt.Sprintf("%s -> %s", ck.Scope.String(), ck.Category.String())
}

// StyleOccurrence menyimpan kemunculan token gaya pada elemen dan file tertentu.
type StyleOccurrence struct {
	FilePath   string
	Span       ir.Span
	Scope      ScopeKey
	Category   PropertyCategory
	Token      string            // Token yang dinormalisasi (misal: "rounded-md", "ring-2", "scale-95")
	Variant    string            // Varian awalan jika ada (misal: "focus-visible", "active", "disabled")
	RawClasses string            // Seluruh string class pada elemen untuk konteks diagnostik
	Tag        string            // Tag AST asli (misal: "Button", "button", "div")
	Props      map[string]string // Atribut statis elemen (misal: size="icon", variant="outline")
}
