package form_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/form"
)

func TestExtractInputFacts(t *testing.T) {
	tests := []struct {
		name          string
		node          *ir.Node
		wantInput     bool
		wantNumber    bool
		wantDisabled  bool
		wantReadOnly  bool
		wantWheel     bool
		wantMin       bool
		wantIdClass   form.IdentifierClass
		wantIdSource  form.IdentifierSource
		wantIdVal     string
		wantIdMatched string
	}{
		{
			name: "Native input type=number postal_code",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "input",
				Attributes: map[string]string{
					"type": `"number"`,
					"name": `"postal_code"`,
				},
			},
			wantInput:     true,
			wantNumber:    true,
			wantDisabled:  false,
			wantReadOnly:  false,
			wantWheel:     false,
			wantMin:       false,
			wantIdClass:   form.IdentifierIdentity,
			wantIdSource:  form.IdentifierSourceName,
			wantIdVal:     "postal_code",
			wantIdMatched: "postal_code",
		},
		{
			name: "JSX Input component with braces and wheel handler",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "Input",
				Attributes: map[string]string{
					"type":    `{"number"}`,
					"name":    `"quantity"`,
					"min":     `"0"`,
					"onWheel": `{(e) => e.currentTarget.blur()}`,
				},
			},
			wantInput:     true,
			wantNumber:    true,
			wantDisabled:  false,
			wantReadOnly:  false,
			wantWheel:     true,
			wantMin:       true,
			wantIdClass:   form.IdentifierQuantity,
			wantIdSource:  form.IdentifierSourceName,
			wantIdVal:     "quantity",
			wantIdMatched: "quantity",
		},
		{
			name: "Dynamic type expression is ignored",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "input",
				Attributes: map[string]string{
					"type": `{dynamicType}`,
					"name": `"quantity"`,
				},
			},
			wantInput:    true,
			wantNumber:   false,
			wantDisabled: false,
			wantReadOnly: false,
			wantWheel:    false,
			wantMin:      false,
			wantIdClass:  form.IdentifierQuantity,
			wantIdSource: form.IdentifierSourceName,
			wantIdVal:    "quantity",
		},
		{
			name: "Disabled input with empty attribute",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "input",
				Attributes: map[string]string{
					"type":     `"number"`,
					"disabled": ``,
					"name":     `"fixed_amount"`,
				},
			},
			wantInput:    true,
			wantNumber:   true,
			wantDisabled: true,
			wantReadOnly: false,
			wantWheel:    false,
			wantMin:      false,
			wantIdClass:  form.IdentifierQuantity,
		},
		{
			name: "ReadOnly input with readOnly attribute",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "input",
				Attributes: map[string]string{
					"type":     `"number"`,
					"readOnly": `"true"`,
					"name":     `"calculated_total"`,
				},
			},
			wantInput:    true,
			wantNumber:   true,
			wantDisabled: false,
			wantReadOnly: true,
			wantWheel:    false,
			wantMin:      false,
			wantIdClass:  form.IdentifierQuantity,
		},
		{
			name: "Dynamic min bound min={domainMin}",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "input",
				Attributes: map[string]string{
					"type": `"number"`,
					"min":  `{domainMin}`,
					"name": `"batch_size"`,
				},
			},
			wantInput:   true,
			wantNumber:  true,
			wantMin:     true,
			wantIdClass: form.IdentifierQuantity,
		},
		{
			name: "Empty min string min=\"\" is missing bound",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "input",
				Attributes: map[string]string{
					"type": `"number"`,
					"min":  `""`,
					"name": `"batch_size"`,
				},
			},
			wantInput:   true,
			wantNumber:  true,
			wantMin:     false,
			wantIdClass: form.IdentifierQuantity,
		},
		{
			name: "Undefined min min={undefined} is missing bound",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "input",
				Attributes: map[string]string{
					"type": `"number"`,
					"min":  `{undefined}`,
					"name": `"batch_size"`,
				},
			},
			wantInput:   true,
			wantNumber:  true,
			wantMin:     false,
			wantIdClass: form.IdentifierQuantity,
		},
		{
			name: "Negative bound min=\"-50\"",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "input",
				Attributes: map[string]string{
					"type": `"number"`,
					"min":  `"-50"`,
					"name": `"temperature"`,
				},
			},
			wantInput:  true,
			wantNumber: true,
			wantMin:    true,
		},
		{
			name: "Ranked evidence: autocomplete over placeholder",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "input",
				Attributes: map[string]string{
					"type":         `"number"`,
					"autocomplete": `"tel"`,
					"placeholder":  `"Enter phone"`,
				},
			},
			wantInput:     true,
			wantNumber:    true,
			wantIdClass:   form.IdentifierIdentity,
			wantIdSource:  form.IdentifierSourceAutocomplete,
			wantIdVal:     "tel",
			wantIdMatched: "tel",
		},
		{
			name: "Ranked evidence: placeholder fallback when no name/id",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "input",
				Attributes: map[string]string{
					"type":        `"number"`,
					"placeholder": `"Enter your NIK"`,
				},
			},
			wantInput:     true,
			wantNumber:    true,
			wantIdClass:   form.IdentifierIdentity,
			wantIdSource:  form.IdentifierSourcePlaceholder,
			wantIdVal:     "Enter your NIK",
			wantIdMatched: "nik",
		},
		{
			name: "Non-input tag is ignored",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "button",
				Attributes: map[string]string{
					"type": `"number"`,
				},
			},
			wantInput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facts := form.ExtractInputFacts(tt.node)
			if facts.IsInputTag != tt.wantInput {
				t.Errorf("IsInputTag = %v, want %v", facts.IsInputTag, tt.wantInput)
			}
			if facts.IsNumberType != tt.wantNumber {
				t.Errorf("IsNumberType = %v, want %v", facts.IsNumberType, tt.wantNumber)
			}
			if facts.IsDisabled != tt.wantDisabled {
				t.Errorf("IsDisabled = %v, want %v", facts.IsDisabled, tt.wantDisabled)
			}
			if facts.IsReadOnly != tt.wantReadOnly {
				t.Errorf("IsReadOnly = %v, want %v", facts.IsReadOnly, tt.wantReadOnly)
			}
			if facts.HasWheelHandler != tt.wantWheel {
				t.Errorf("HasWheelHandler = %v, want %v", facts.HasWheelHandler, tt.wantWheel)
			}
			if facts.HasDeclaredMin != tt.wantMin {
				t.Errorf("HasDeclaredMin = %v, want %v", facts.HasDeclaredMin, tt.wantMin)
			}
			if tt.wantIdClass != 0 && facts.Identifier.Class != tt.wantIdClass {
				t.Errorf("Identifier.Class = %v, want %v", facts.Identifier.Class, tt.wantIdClass)
			}
			if tt.wantIdSource != 0 && facts.Identifier.Source != tt.wantIdSource {
				t.Errorf("Identifier.Source = %v, want %v", facts.Identifier.Source, tt.wantIdSource)
			}
			if tt.wantIdVal != "" && facts.Identifier.Value != tt.wantIdVal {
				t.Errorf("Identifier.Value = %q, want %q", facts.Identifier.Value, tt.wantIdVal)
			}
			if tt.wantIdMatched != "" && facts.Identifier.Matched != tt.wantIdMatched {
				t.Errorf("Identifier.Matched = %q, want %q", facts.Identifier.Matched, tt.wantIdMatched)
			}
		})
	}
}
