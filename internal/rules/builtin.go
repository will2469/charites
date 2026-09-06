package rules

import (
	"github.com/will2469/charites/internal/rules/a11y"
	"github.com/will2469/charites/internal/rules/browser"
	"github.com/will2469/charites/internal/rules/cls"
	"github.com/will2469/charites/internal/rules/ergonomy"
	"github.com/will2469/charites/internal/rules/inp"
	"github.com/will2469/charites/internal/rules/mobile"
	"github.com/will2469/charites/internal/rules/pwa"
	"github.com/will2469/charites/internal/rules/responsive"
	"github.com/will2469/charites/internal/rules/theme"
	"github.com/will2469/charites/internal/rules/ux"
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

	// browser Wave 1 (Rendering & Styling Multi-Engine)
	_ Rule = (*browser.AppearanceNativeOverrideRule)(nil)
	_ Rule = (*browser.ScrollbarVendorIncompleteRule)(nil)
	_ Rule = (*browser.ObsoleteVendorPrefixRule)(nil)
	_ Rule = (*browser.HoverOnlyInteractionRule)(nil)

	// browser Wave 2 (Runtime Safety & Event Performance)
	_ Rule = (*browser.ExperimentalAPINoFeaturedetectRule)(nil)
	_ Rule = (*browser.DateInputFormatAssumptionRule)(nil)
	_ Rule = (*browser.NonPassiveScrollListenerRule)(nil)

	// pwa Wave 1 (Web App Manifest & Branding)
	_ Rule = (*pwa.ManifestRequiredFieldsMissingRule)(nil)
	_ Rule = (*pwa.IconMaskableMissingRule)(nil)
	_ Rule = (*pwa.ManifestMissingRule)(nil)
	_ Rule = (*pwa.StartURLInconsistencyRule)(nil)

	// pwa Wave 2 (Apple Standalone & Security)
	_ Rule = (*pwa.AppleMetaMissingRule)(nil)
	_ Rule = (*pwa.InsecureContextResourceRule)(nil)

	// pwa Wave 3 (Service Worker Lifecycle & Offline Cache)
	_ Rule = (*pwa.ServiceWorkerNoOfflineFallbackRule)(nil)
	_ Rule = (*pwa.ServiceWorkerMissingRule)(nil)
	_ Rule = (*pwa.ServiceWorkerRegistrationRule)(nil)
	_ Rule = (*pwa.CacheRuntimeAPIRiskRule)(nil)

	// responsive Wave 1 (Layout Core & Viewport Deterministic)
	_ Rule = (*responsive.MissingBreakpointRule)(nil)
	_ Rule = (*responsive.UnwrappedTableOverflowRule)(nil)
	_ Rule = (*responsive.FixedWidthOverflowRule)(nil)
	_ Rule = (*responsive.ViewportUnitLeakRule)(nil)
	_ Rule = (*responsive.SafeAreaMissingRule)(nil)
	_ Rule = (*responsive.ViewportMetaMissingRule)(nil)

	// responsive Wave 2 (Overflow Integrity & Flex Geometry)
	_ Rule = (*responsive.HorizontalOverflowRule)(nil)
	_ Rule = (*responsive.FlexChildOverflowRule)(nil)
	_ Rule = (*responsive.ImageOverflowRule)(nil)
	_ Rule = (*responsive.MobileTextOverflowRule)(nil)

	// responsive Wave 3 (Mobile-First Content & Viewport Dynamics)
	_ Rule = (*responsive.DesktopOnlyContentRule)(nil)
	_ Rule = (*responsive.MobileDensityOverloadRule)(nil)
	_ Rule = (*responsive.DynamicViewportInconsistencyRule)(nil)

	// responsive Wave 4 (Viewport Dynamics, Keyboard Obstruction & Grid Physics)
	_ Rule = (*responsive.KeyboardObstructionRule)(nil)
	_ Rule = (*responsive.ContainerOverconstraintRule)(nil)
	_ Rule = (*responsive.GridMinColumnRule)(nil)
	_ Rule = (*responsive.AspectRatioOverflowRule)(nil)

	// ux Wave 1 (Spatial Hierarchy, Navigation Chunking & CTA Clarity)
	_ Rule = (*ux.SpacingInversionRule)(nil)
	_ Rule = (*ux.NavOverflowChunkingRule)(nil)
	_ Rule = (*ux.CompetingPrimaryCTARule)(nil)
	_ Rule = (*ux.CamouflagedLinkRule)(nil)

	// ux Wave 2 (Information Architecture, Mental Models & Form Load)
	_ Rule = (*ux.UnconventionalHomeLinkRule)(nil)
	_ Rule = (*ux.RadioOverchoiceRule)(nil)
	_ Rule = (*ux.MonolithicFormBloatRule)(nil)
	_ Rule = (*ux.MissingAutofillRule)(nil)

	// ux Wave 3 (State Transparency, Defensive UI & Feedback Visibility)
	_ Rule = (*ux.EmptyCollectionUnhandledRule)(nil)
	_ Rule = (*ux.DisabledControlNoExplanationRule)(nil)
	_ Rule = (*ux.OrphanedErrorStateRule)(nil)
	_ Rule = (*ux.UnthrottledInputHandlerRule)(nil)

	// ux Wave 4 (Async Control-Flow, Action Safety & Doherty Threshold)
	_ Rule = (*ux.SubmitFeedbackMissingRule)(nil)
	_ Rule = (*ux.UnboundedAsyncFlagRule)(nil)
	_ Rule = (*ux.DestructiveActionUnconfirmedRule)(nil)
	_ Rule = (*ux.SilentCatchSwallowRule)(nil)

	// cls Wave 1 (Rendering Box Reservation, Embed Frames, Ad Slots & Slider Physics)
	_ Rule = (*cls.UnsizedImageRule)(nil)
	_ Rule = (*cls.UnsizedEmbedFrameRule)(nil)
	_ Rule = (*cls.UnreservedAdContainerRule)(nil)
	_ Rule = (*cls.UnconstrainedCarouselRule)(nil)

	// cls Wave 2 (Font Loading & Metric Stability)
	_ Rule = (*cls.FontDisplayMissingRule)(nil)
	_ Rule = (*cls.UnadjustedFontMetricRule)(nil)
	_ Rule = (*cls.FontImportLateDiscoveryRule)(nil)
	_ Rule = (*cls.TextIconLateReflowRule)(nil)

	// cls Wave 3 (CSS Animations, Transitions, Scrollbar Gutter, & Dynamic Table Layouts)
	_ Rule = (*cls.LayoutTriggerAnimationRule)(nil)
	_ Rule = (*cls.LayoutTriggerTransitionRule)(nil)
	_ Rule = (*cls.UnstableScrollbarGutterRule)(nil)
	_ Rule = (*cls.DynamicTableReflowRule)(nil)

	// cls Wave 4 (Lifecycle DOM, Client-Only Hydration, Fixed Headers, Dynamic Injections, & Collapsibles)
	_ Rule = (*cls.ClientOnlyHydrationPopRule)(nil)
	_ Rule = (*cls.UnreservedFixedHeaderRule)(nil)
	_ Rule = (*cls.DynamicContentWithoutReservedSpaceRule)(nil)
	_ Rule = (*cls.CollapsibleHeightJumpRule)(nil)

	// inp Wave 1 (Event Handler Execution & Synchronous Work)
	_ Rule = (*inp.LayoutThrashingRule)(nil)
	_ Rule = (*inp.HeavyEventHandlerRule)(nil)
	_ Rule = (*inp.RepeatedStateUpdateRule)(nil)
	_ Rule = (*inp.UnyieldedLongTaskRule)(nil)
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

		// a11y Wave 3 (Shadcn UI Hierarchy & Astro Image)
		a11y.NewFormLabelMissingControlRule(),
		a11y.NewFormLabelCompositeControlRule(),
		a11y.NewImgMissingAltRule(),

		// a11y Wave 4 (Modal Traps & Semantic Interactions)
		a11y.NewButtonTypeMissingRule(),
		a11y.NewKeyboardTrapMissingEscapeRule(),
		a11y.NewDialogMissingAriaRule(),
		a11y.NewEmptyInteractiveRule(),

		// browser Wave 1 (Rendering & Styling Multi-Engine)
		browser.NewAppearanceNativeOverrideRule(),
		browser.NewScrollbarVendorIncompleteRule(),
		browser.NewObsoleteVendorPrefixRule(),
		browser.NewHoverOnlyInteractionRule(),

		// browser Wave 2 (Runtime Safety & Event Performance)
		browser.NewExperimentalAPINoFeaturedetectRule(),
		browser.NewDateInputFormatAssumptionRule(),
		browser.NewNonPassiveScrollListenerRule(),

		// browser Wave 3 (Browser Capability & Vendor API Isolation)
		browser.NewUserAgentSniffingRule(),
		browser.NewWebKitOnlyAPIRule(),
		browser.NewChromeOnlyAPIRule(),
		browser.NewFirefoxOnlyAPIRule(),
		browser.NewSafariOnlyAPIRule(),

		// ergonomy Wave 1 (Virtual Keypad, Touch Feedback & Gesture Ergonomics)
		ergonomy.NewMissingInputmodeKeyboardRule(),
		ergonomy.NewTapHighlightNotHandledRule(),
		ergonomy.NewGestureWithoutTouchActionRule(),

		// ergonomy Wave 2 (Thumb Zone & Navigation Ergonomics)
		ergonomy.NewBottomNavThumbUnreachableRule(),

		// mobile Wave 3 (Mobile Viewport & Obstruction Physics)
		mobile.NewKeyboardViewportRiskRule(),
		mobile.NewFixedActionObstructionRule(),
		mobile.NewModalViewportLockRule(),
		mobile.NewOrientationLockRiskRule(),
		mobile.NewPointerEventsBlockRule(),

		// pwa Wave 1 (Web App Manifest & Branding)
		pwa.NewManifestRequiredFieldsMissingRule(),
		pwa.NewIconMaskableMissingRule(),
		pwa.NewManifestMissingRule(),
		pwa.NewStartURLInconsistencyRule(),

		// pwa Wave 2 (Apple Standalone & Security)
		pwa.NewAppleMetaMissingRule(),
		pwa.NewInsecureContextResourceRule(),

		// pwa Wave 3 (Service Worker Lifecycle & Offline Cache)
		pwa.NewServiceWorkerNoOfflineFallbackRule(),
		pwa.NewServiceWorkerMissingRule(),
		pwa.NewServiceWorkerRegistrationRule(),
		pwa.NewCacheRuntimeAPIRiskRule(),

		// responsive Wave 1 (Layout Core & Viewport Deterministic)
		responsive.NewMissingBreakpointRule(),
		responsive.NewUnwrappedTableOverflowRule(),
		responsive.NewFixedWidthOverflowRule(),
		responsive.NewViewportUnitLeakRule(),
		responsive.NewSafeAreaMissingRule(),
		responsive.NewViewportMetaMissingRule(),

		// responsive Wave 2 (Overflow Integrity & Flex Geometry)
		responsive.NewHorizontalOverflowRule(),
		responsive.NewFlexChildOverflowRule(),
		responsive.NewImageOverflowRule(),
		responsive.NewMobileTextOverflowRule(),

		// responsive Wave 3 (Mobile-First Content & Viewport Dynamics)
		responsive.NewDesktopOnlyContentRule(),
		responsive.NewMobileDensityOverloadRule(),
		responsive.NewDynamicViewportInconsistencyRule(),

		// responsive Wave 4 (Viewport Dynamics, Keyboard Obstruction & Grid Physics)
		responsive.NewKeyboardObstructionRule(),
		responsive.NewContainerOverconstraintRule(),
		responsive.NewGridMinColumnRule(),
		responsive.NewAspectRatioOverflowRule(),

		// ux Wave 1 (Spatial Hierarchy, Navigation Chunking & CTA Clarity)
		ux.NewSpacingInversionRule(),
		ux.NewNavOverflowChunkingRule(),
		ux.NewCompetingPrimaryCTARule(),
		ux.NewCamouflagedLinkRule(),

		// ux Wave 2 (Information Architecture, Mental Models & Form Load)
		ux.NewUnconventionalHomeLinkRule(),
		ux.NewRadioOverchoiceRule(),
		ux.NewMonolithicFormBloatRule(),
		ux.NewMissingAutofillRule(),

		// ux Wave 3 (State Transparency, Defensive UI & Feedback Visibility)
		ux.NewEmptyCollectionUnhandledRule(),
		ux.NewDisabledControlNoExplanationRule(),
		ux.NewOrphanedErrorStateRule(),
		ux.NewUnthrottledInputHandlerRule(),

		// ux Wave 4 (Async Control-Flow, Action Safety & Doherty Threshold)
		ux.NewSubmitFeedbackMissingRule(),
		ux.NewUnboundedAsyncFlagRule(),
		ux.NewDestructiveActionUnconfirmedRule(),
		ux.NewSilentCatchSwallowRule(),

		// cls Wave 1 (Rendering Box Reservation, Embed Frames, Ad Slots & Slider Physics)
		cls.NewUnsizedImageRule(),
		cls.NewUnsizedEmbedFrameRule(),
		cls.NewUnreservedAdContainerRule(),
		cls.NewUnconstrainedCarouselRule(),

		// cls Wave 2 (Font Loading & Metric Stability)
		cls.NewFontDisplayMissingRule(),
		cls.NewUnadjustedFontMetricRule(),
		cls.NewFontImportLateDiscoveryRule(),
		cls.NewTextIconLateReflowRule(),

		// cls Wave 3 (CSS Animations, Transitions, Scrollbar Gutter, & Dynamic Table Layouts)
		cls.NewLayoutTriggerAnimationRule(),
		cls.NewLayoutTriggerTransitionRule(),
		cls.NewUnstableScrollbarGutterRule(),
		cls.NewDynamicTableReflowRule(),

		// cls Wave 4 (Lifecycle DOM, Client-Only Hydration, Fixed Headers, Dynamic Injections, & Collapsibles)
		cls.NewClientOnlyHydrationPopRule(),
		cls.NewUnreservedFixedHeaderRule(),
		cls.NewDynamicContentWithoutReservedSpaceRule(),
		cls.NewCollapsibleHeightJumpRule(),

		// inp Wave 1 (Event Handler Execution & Synchronous Work)
		inp.NewLayoutThrashingRule(),
		inp.NewHeavyEventHandlerRule(),
		inp.NewRepeatedStateUpdateRule(),
		inp.NewUnyieldedLongTaskRule(),
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
