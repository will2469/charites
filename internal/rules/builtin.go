package rules

import (
	"github.com/will2469/charites/internal/rules/a11y"
	"github.com/will2469/charites/internal/rules/theme"
)

// Invariant static compile-time check: seluruh rule wajib mengimplementasikan interface Rule.
var (
	_ Rule = (*theme.HardcodeOpacityColorRule)(nil)
	_ Rule = (*theme.HardcodeColorRule)(nil)
	_ Rule = (*theme.PrimitiveInComponentRule)(nil)
	_ Rule = (*theme.HardcodeMonochromeRule)(nil)
	_ Rule = (*theme.HardcodeBorderColorRule)(nil)
	_ Rule = (*theme.GradientHardcodeRule)(nil)
	_ Rule = (*theme.InlineStyleHardcodeRule)(nil)
	_ Rule = (*theme.PseudoHardcodeColorRule)(nil)
	_ Rule = (*theme.ImportantOverrideRule)(nil)
	_ Rule = (*theme.HardcodeSizeRule)(nil)
	_ Rule = (*theme.HardcodeBorderRadiusRule)(nil)
	_ Rule = (*theme.HardcodeZIndexRule)(nil)
	_ Rule = (*theme.HardcodeShadowColorRule)(nil)
	_ Rule = (*theme.BackdropBlurHardcodeRule)(nil)
	_ Rule = (*theme.FocusRingHardcodeRule)(nil)
	_ Rule = (*theme.UnpairedDarkVariantRule)(nil)
	_ Rule = (*theme.ShadowWithoutBorderDarkRule)(nil)
	_ Rule = (*theme.NestedOpacityContrastRule)(nil)
	_ Rule = (*theme.ImageThemeHardcodeRule)(nil)
	_ Rule = (*theme.SVGHardcodeFillRule)(nil)
	_ Rule = (*theme.UnlayeredTokenDefinitionRule)(nil)
	_ Rule = (*theme.MissingTokenFallbackRule)(nil)
	_ Rule = (*theme.TokenSourceDriftRule)(nil)
	_ Rule = (*theme.ApplyBloatRule)(nil)

	// a11y Wave 1 (Physical Touch & Mobile Viewport Ergonomics)
	_ Rule = (*a11y.TouchTargetSizeRule)(nil)
	_ Rule = (*a11y.TouchTargetSpacingRule)(nil)
	_ Rule = (*a11y.InputIOSZoomHazardRule)(nil)
	_ Rule = (*a11y.InputCrampedPaddingRule)(nil)
	_ Rule = (*a11y.MissingFocusRingRule)(nil)

	// a11y Wave 2 (Form Relational Integrity & Error Announcements)
	_ Rule = (*a11y.ErrorNotAnnouncedRule)(nil)
	_ Rule = (*a11y.PlaceholderAsLabelRule)(nil)
	_ Rule = (*a11y.LabelMissingControlRule)(nil)
	_ Rule = (*a11y.FormInputMissingNameRule)(nil)
)

func builtinRules() []Rule {
	return []Rule{
		theme.NewHardcodeOpacityColorRule(),
		theme.NewHardcodeColorRule(),
		theme.NewPrimitiveInComponentRule(),
		theme.NewHardcodeMonochromeRule(),
		theme.NewHardcodeBorderColorRule(),
		theme.NewGradientHardcodeRule(),
		theme.NewInlineStyleHardcodeRule(),
		theme.NewPseudoHardcodeColorRule(),
		theme.NewImportantOverrideRule(),
		theme.NewHardcodeSizeRule(),
		theme.NewHardcodeBorderRadiusRule(),
		theme.NewHardcodeZIndexRule(),
		theme.NewHardcodeShadowColorRule(),
		theme.NewBackdropBlurHardcodeRule(),
		theme.NewFocusRingHardcodeRule(),
		theme.NewUnpairedDarkVariantRule(),
		theme.NewShadowWithoutBorderDarkRule(),
		theme.NewNestedOpacityContrastRule(),
		theme.NewImageThemeHardcodeRule(),
		theme.NewSVGHardcodeFillRule(),
		theme.NewUnlayeredTokenDefinitionRule(),
		theme.NewMissingTokenFallbackRule(),
		theme.NewTokenSourceDriftRule(),
		theme.NewApplyBloatRule(),
		theme.NewMissingColorSchemeRule(),
		theme.NewMetaThemeColorMismatchRule(),
		theme.NewNoReducedMotionRule(),
		theme.NewChartColorHardcodeRule(),
		theme.NewDynamicClassRule(),
		theme.NewDualStrategyCollisionRule(),
		theme.NewHydrationThemeMismatchRule(),
		theme.NewSplitThemeStateRule(),

		// a11y Wave 1 (Physical Touch & Mobile Viewport Ergonomics)
		a11y.NewTouchTargetSizeRule(),
		a11y.NewTouchTargetSpacingRule(),
		a11y.NewInputIOSZoomHazardRule(),
		a11y.NewInputCrampedPaddingRule(),
		a11y.NewMissingFocusRingRule(),

		// a11y Wave 2 (Form Relational Integrity & Error Announcements)
		a11y.NewErrorNotAnnouncedRule(),
		a11y.NewPlaceholderAsLabelRule(),
		a11y.NewLabelMissingControlRule(),
		a11y.NewFormInputMissingNameRule(),
	}
}

func init() {
	for _, r := range builtinRules() {
		_ = Register(r)
	}
}

// RegisterBuiltinRules mendaftarkan seluruh built-in rule Charites ke registry yang diberikan.
func RegisterBuiltinRules(reg *Registry) error {
	if reg == nil {
		return ErrNilRule
	}
	for _, r := range builtinRules() {
		if err := reg.Register(r); err != nil {
			return err
		}
	}
	return nil
}
