package a11y_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/a11y"
)

func TestTouchTargetSizeRule(t *testing.T) {
	rule := a11y.NewTouchTargetSizeRule()

	tests := []struct {
		name        string
		node        *ir.Node
		shouldFault bool
	}{
		{
			name: "small_button_24px",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "button",
				Classes: []string{"h-6", "w-6"},
			},
			shouldFault: true,
		},
		{
			name: "small_button_32px",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "button",
				Classes: []string{"size-8"},
			},
			shouldFault: true,
		},
		{
			name: "standard_touch_target_44px",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "button",
				Classes: []string{"h-11", "w-11"},
			},
			shouldFault: false,
		},
		{
			name: "min_size_compensation",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "button",
				Classes: []string{"h-8", "w-8", "min-h-11", "min-w-11"},
			},
			shouldFault: false,
		},
		{
			name: "inline_link_in_paragraph",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "a",
				Classes: []string{"text-blue-600"},
				Parent: &ir.Node{
					Type: ir.NodeElement,
					Tag:  "p",
				},
			},
			shouldFault: false,
		},
		{
			name: "non_interactive_div",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "div",
				Classes: []string{"h-6", "w-6"},
			},
			shouldFault: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := rule.Evaluate(tc.node)
			if tc.shouldFault && len(diags) == 0 {
				t.Errorf("expected violation for %s, got 0", tc.name)
			}
			if !tc.shouldFault && len(diags) > 0 {
				t.Errorf("unexpected violation for %s: %+v", tc.name, diags)
			}
		})
	}
}

func TestTouchTargetSpacingRule(t *testing.T) {
	rule := a11y.NewTouchTargetSpacingRule()

	tests := []struct {
		name        string
		node        *ir.Node
		shouldFault bool
	}{
		{
			name: "flex_buttons_zero_gap",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "div",
				Classes: []string{"flex", "gap-0"},
				Children: []*ir.Node{
					{Type: ir.NodeElement, Tag: "button", Classes: []string{"h-11", "w-11"}},
					{Type: ir.NodeElement, Tag: "button", Classes: []string{"h-11", "w-11"}},
				},
			},
			shouldFault: true,
		},
		{
			name: "flex_buttons_cramped_gap_1",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "div",
				Classes: []string{"flex", "gap-1"},
				Children: []*ir.Node{
					{Type: ir.NodeElement, Tag: "button", Classes: []string{"h-11", "w-11"}},
					{Type: ir.NodeElement, Tag: "button", Classes: []string{"h-11", "w-11"}},
				},
			},
			shouldFault: true,
		},
		{
			name: "flex_buttons_adequate_gap_2",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "div",
				Classes: []string{"flex", "gap-2"},
				Children: []*ir.Node{
					{Type: ir.NodeElement, Tag: "button", Classes: []string{"h-11", "w-11"}},
					{Type: ir.NodeElement, Tag: "button", Classes: []string{"h-11", "w-11"}},
				},
			},
			shouldFault: false,
		},
		{
			name: "flex_single_button",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "div",
				Classes: []string{"flex"},
				Children: []*ir.Node{
					{Type: ir.NodeElement, Tag: "button", Classes: []string{"h-11", "w-11"}},
					{Type: ir.NodeElement, Tag: "span", Classes: []string{"text-sm"}},
				},
			},
			shouldFault: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := rule.Evaluate(tc.node)
			if tc.shouldFault && len(diags) == 0 {
				t.Errorf("expected violation for %s, got 0", tc.name)
			}
			if !tc.shouldFault && len(diags) > 0 {
				t.Errorf("unexpected violation for %s: %+v", tc.name, diags)
			}
		})
	}
}

func TestInputIOSZoomHazardRule(t *testing.T) {
	rule := a11y.NewInputIOSZoomHazardRule()

	tests := []struct {
		name        string
		node        *ir.Node
		shouldFault bool
	}{
		{
			name: "input_text_sm_violation",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "input",
				Classes: []string{"text-sm", "border", "px-3"},
			},
			shouldFault: true,
		},
		{
			name: "input_text_xs_violation",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "input",
				Classes: []string{"text-xs", "border"},
			},
			shouldFault: true,
		},
		{
			name: "select_arbitrary_14px_violation",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "select",
				Classes: []string{"text-[14px]", "border"},
			},
			shouldFault: true,
		},
		{
			name: "input_safe_responsive_override",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "input",
				Classes: []string{"text-base", "sm:text-sm", "border"},
			},
			shouldFault: false,
		},
		{
			name: "input_desktop_only_sm_override",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "input",
				Classes: []string{"sm:text-sm", "border"},
			},
			shouldFault: false,
		},
		{
			name: "input_checkbox_ignored",
			node: &ir.Node{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"type": "checkbox"},
				Classes:    []string{"text-sm"},
			},
			shouldFault: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := rule.Evaluate(tc.node)
			if tc.shouldFault && len(diags) == 0 {
				t.Errorf("expected violation for %s, got 0", tc.name)
			}
			if !tc.shouldFault && len(diags) > 0 {
				t.Errorf("unexpected violation for %s: %+v", tc.name, diags)
			}
		})
	}
}

func TestInputCrampedPaddingRule(t *testing.T) {
	rule := a11y.NewInputCrampedPaddingRule()

	tests := []struct {
		name        string
		node        *ir.Node
		shouldFault bool
	}{
		{
			name: "input_cramped_h8",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "input",
				Classes: []string{"h-8", "px-3", "border"},
			},
			shouldFault: true,
		},
		{
			name: "input_cramped_py1",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "input",
				Classes: []string{"py-1", "px-3", "border"},
			},
			shouldFault: true,
		},
		{
			name: "input_ergonomic_h11",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "input",
				Classes: []string{"h-11", "px-3.5", "border"},
			},
			shouldFault: false,
		},
		{
			name: "input_ergonomic_py2_5",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "input",
				Classes: []string{"py-2.5", "px-3.5", "border"},
			},
			shouldFault: false,
		},
		{
			name: "input_min_h_compensation",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "input",
				Classes: []string{"h-8", "min-h-11", "px-3"},
			},
			shouldFault: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := rule.Evaluate(tc.node)
			if tc.shouldFault && len(diags) == 0 {
				t.Errorf("expected violation for %s, got 0", tc.name)
			}
			if !tc.shouldFault && len(diags) > 0 {
				t.Errorf("unexpected violation for %s: %+v", tc.name, diags)
			}
		})
	}
}

func TestMissingFocusRingRule(t *testing.T) {
	rule := a11y.NewMissingFocusRingRule()

	tests := []struct {
		name        string
		node        *ir.Node
		shouldFault bool
	}{
		{
			name: "button_outline_none_violation",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "button",
				Classes: []string{"outline-none", "bg-primary", "text-white"},
			},
			shouldFault: true,
		},
		{
			name: "link_focus_outline_none_violation",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "a",
				Classes: []string{"focus:outline-none", "text-blue-600"},
			},
			shouldFault: true,
		},
		{
			name: "button_with_focus_visible_ring_safe",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "button",
				Classes: []string{"outline-none", "focus-visible:ring-2", "focus-visible:ring-primary"},
			},
			shouldFault: false,
		},
		{
			name: "button_with_focus_ring_safe",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "button",
				Classes: []string{"focus:outline-none", "focus:ring-2"},
			},
			shouldFault: false,
		},
		{
			name: "non_interactive_div_outline_none_safe",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "div",
				Classes: []string{"outline-none"},
			},
			shouldFault: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := rule.Evaluate(tc.node)
			if tc.shouldFault && len(diags) == 0 {
				t.Errorf("expected violation for %s, got 0", tc.name)
			}
			if !tc.shouldFault && len(diags) > 0 {
				t.Errorf("unexpected violation for %s: %+v", tc.name, diags)
			}
		})
	}
}

func TestErrorNotAnnouncedRule(t *testing.T) {
	rule := a11y.NewErrorNotAnnouncedRule()

	tests := []struct {
		name        string
		node        *ir.Node
		shouldFault bool
	}{
		{
			name: "aria_invalid_without_describedby_violation",
			node: &ir.Node{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"aria-invalid": "true"},
			},
			shouldFault: true,
		},
		{
			name: "aria_invalid_with_describedby_safe",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "input",
				Attributes: map[string]string{
					"aria-invalid":     "true",
					"aria-describedby": "email-error",
				},
			},
			shouldFault: false,
		},
		{
			name: "aria_invalid_false_safe",
			node: &ir.Node{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"aria-invalid": "false"},
			},
			shouldFault: false,
		},
		{
			name: "no_aria_invalid_safe",
			node: &ir.Node{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"id": "username"},
			},
			shouldFault: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := rule.Evaluate(tc.node)
			if tc.shouldFault && len(diags) == 0 {
				t.Errorf("expected violation for %s, got 0", tc.name)
			}
			if !tc.shouldFault && len(diags) > 0 {
				t.Errorf("unexpected violation for %s: %+v", tc.name, diags)
			}
		})
	}
}

func TestPlaceholderAsLabelRule(t *testing.T) {
	rule := a11y.NewPlaceholderAsLabelRule()

	rootWithLabel := &ir.Node{
		Type: ir.NodeElement,
		Tag:  "form",
		Children: []*ir.Node{
			{
				Type:       ir.NodeElement,
				Tag:        "label",
				Attributes: map[string]string{"htmlFor": "email"},
			},
			{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"id": "email", "placeholder": "nama@domain.com"},
			},
		},
	}
	rootWithLabel.Children[0].Parent = rootWithLabel
	rootWithLabel.Children[1].Parent = rootWithLabel

	tests := []struct {
		name        string
		node        *ir.Node
		shouldFault bool
	}{
		{
			name: "bare_placeholder_without_label_violation",
			node: &ir.Node{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"placeholder": "Enter email"},
			},
			shouldFault: true,
		},
		{
			name: "input_with_aria_label_safe",
			node: &ir.Node{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"placeholder": "Search...", "aria-label": "Search"},
			},
			shouldFault: false,
		},
		{
			name:        "input_with_associated_label_in_tree_safe",
			node:        rootWithLabel.Children[1],
			shouldFault: false,
		},
		{
			name: "input_enclosed_in_label_safe",
			node: func() *ir.Node {
				parentLabel := &ir.Node{Type: ir.NodeElement, Tag: "label"}
				childInput := &ir.Node{
					Type:       ir.NodeElement,
					Tag:        "input",
					Attributes: map[string]string{"placeholder": "Search"},
					Parent:     parentLabel,
				}
				parentLabel.Children = []*ir.Node{childInput}
				return childInput
			}(),
			shouldFault: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := rule.Evaluate(tc.node)
			if tc.shouldFault && len(diags) == 0 {
				t.Errorf("expected violation for %s, got 0", tc.name)
			}
			if !tc.shouldFault && len(diags) > 0 {
				t.Errorf("unexpected violation for %s: %+v", tc.name, diags)
			}
		})
	}
}

func TestLabelMissingControlRule(t *testing.T) {
	rule := a11y.NewLabelMissingControlRule()

	docWithInput := &ir.Node{
		Type: ir.NodeElement,
		Tag:  "form",
		Children: []*ir.Node{
			{
				Type:       ir.NodeElement,
				Tag:        "label",
				Attributes: map[string]string{"htmlFor": "user_id"},
			},
			{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"id": "userId"}, // Typo! user_id != userId
			},
			{
				Type:       ir.NodeElement,
				Tag:        "label",
				Attributes: map[string]string{"htmlFor": "userId"}, // Matching!
			},
		},
	}
	for _, ch := range docWithInput.Children {
		ch.Parent = docWithInput
	}

	tests := []struct {
		name        string
		node        *ir.Node
		shouldFault bool
	}{
		{
			name:        "label_htmlfor_mismatch_violation",
			node:        docWithInput.Children[0],
			shouldFault: true,
		},
		{
			name:        "label_htmlfor_matching_safe",
			node:        docWithInput.Children[2],
			shouldFault: false,
		},
		{
			name: "label_wrapping_input_without_htmlfor_safe",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "label",
			},
			shouldFault: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := rule.Evaluate(tc.node)
			if tc.shouldFault && len(diags) == 0 {
				t.Errorf("expected violation for %s, got 0", tc.name)
			}
			if !tc.shouldFault && len(diags) > 0 {
				t.Errorf("unexpected violation for %s: %+v", tc.name, diags)
			}
		})
	}
}

func TestFormInputMissingNameRule(t *testing.T) {
	rule := a11y.NewFormInputMissingNameRule()

	tests := []struct {
		name        string
		node        *ir.Node
		shouldFault bool
	}{
		{
			name: "input_without_name_or_id_violation",
			node: &ir.Node{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"type": "text"},
			},
			shouldFault: true,
		},
		{
			name: "input_with_name_safe",
			node: &ir.Node{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"name": "email"},
			},
			shouldFault: false,
		},
		{
			name: "input_with_id_safe",
			node: &ir.Node{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"id": "email"},
			},
			shouldFault: false,
		},
		{
			name: "submit_button_input_safe",
			node: &ir.Node{
				Type:       ir.NodeElement,
				Tag:        "input",
				Attributes: map[string]string{"type": "submit"},
			},
			shouldFault: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := rule.Evaluate(tc.node)
			if tc.shouldFault && len(diags) == 0 {
				t.Errorf("expected violation for %s, got 0", tc.name)
			}
			if !tc.shouldFault && len(diags) > 0 {
				t.Errorf("unexpected violation for %s: %+v", tc.name, diags)
			}
		})
	}
}
