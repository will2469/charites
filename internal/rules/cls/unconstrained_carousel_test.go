package cls

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
)

func TestClassifySlider_Disambiguation(t *testing.T) {
	tests := []struct {
		name     string
		node     *ir.Node
		expected SliderClassification
	}{
		{
			name: "Cerberus ChallengeSlider (onSolve + solved) -> Control",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "ChallengeSlider",
				Attributes: map[string]string{
					"onSolve":  "cerberus.solveChallenge",
					"onCancel": "onCancel",
					"error":    "cerberus.error",
					"solved":   "cerberus.challengeSolved",
				},
			},
			expected: SliderControl,
		},
		{
			name: "Radix UI Slider (min, max, step, onValueChange) -> Control",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "Slider",
				Attributes: map[string]string{
					"min":           "0",
					"max":           "100",
					"step":          "1",
					"value":         "value",
					"onValueChange": "setValue",
				},
			},
			expected: SliderControl,
		},
		{
			name: "RangeSlider with min and max -> Control",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "RangeSlider",
				Attributes: map[string]string{
					"min": "0",
					"max": "100",
				},
			},
			expected: SliderControl,
		},
		{
			name: "RangeSlider without props -> Unknown (tag alone is not authoritative)",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "RangeSlider",
			},
			expected: SliderUnknown,
		},
		{
			name: "Semantic div role=slider with aria-valuenow -> Control",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "div",
				Attributes: map[string]string{
					"role":          "slider",
					"aria-valuenow": "50",
					"aria-valuemin": "0",
					"aria-valuemax": "100",
				},
			},
			expected: SliderControl,
		},
		{
			name: "Ambiguous standalone <Slider /> without props -> Unknown",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "Slider",
			},
			expected: SliderUnknown,
		},
		{
			name: "Dedicated Carousel with generic min, max, value -> ContentCarousel",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "Carousel",
				Attributes: map[string]string{
					"min":   "0",
					"max":   "10",
					"value": "activeSlide",
				},
			},
			expected: SliderContentCarousel,
		},
		{
			name: "Dedicated Carousel with value -> ContentCarousel",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "Carousel",
				Attributes: map[string]string{
					"value": "activeSlide",
				},
			},
			expected: SliderContentCarousel,
		},
		{
			name: "Dedicated Carousel with spread props -> ContentCarousel",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "Carousel",
				Attributes: map[string]string{
					"{...props}": "",
				},
			},
			expected: SliderContentCarousel,
		},
		{
			name: "BannerSlider with onChange -> ContentCarousel",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "BannerSlider",
				Attributes: map[string]string{
					"onChange": "handleSlideChange",
				},
			},
			expected: SliderContentCarousel,
		},
		{
			name: "Generic Slider with horizontal scroll snap -> ContentCarousel",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "Slider",
				Classes: []string{"overflow-x-auto", "snap-x"},
				Children: []*ir.Node{
					{
						Type: ir.NodeElement,
						Tag:  "Slide",
					},
				},
			},
			expected: SliderContentCarousel,
		},
		{
			name: "Contradictory strong control + structural carousel -> Unknown (conservative suppression)",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "Slider",
				Classes: []string{"overflow-x-auto", "snap-x"},
				Attributes: map[string]string{
					"role":          "slider",
					"aria-valuenow": "50",
				},
				Children: []*ir.Node{
					{
						Type: ir.NodeElement,
						Tag:  "Slide",
					},
				},
			},
			expected: SliderUnknown,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := classifySlider(tc.node)
			if got != tc.expected {
				t.Errorf("classifySlider() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestUnconstrainedCarouselRule_Evaluate(t *testing.T) {
	rule := NewUnconstrainedCarouselRule()

	tests := []struct {
		name          string
		node          *ir.Node
		wantDiagCount int
	}{
		{
			name: "ChallengeSlider -> 0 diags (Control)",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "ChallengeSlider",
				Attributes: map[string]string{
					"onSolve": "cerberus.solveChallenge",
					"solved":  "cerberus.challengeSolved",
				},
			},
			wantDiagCount: 0,
		},
		{
			name: "Radix Slider -> 0 diags (Control)",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "Slider",
				Attributes: map[string]string{
					"min":           "0",
					"max":           "100",
					"step":          "1",
					"onValueChange": "setValue",
				},
			},
			wantDiagCount: 0,
		},
		{
			name: "Standalone <Slider /> -> 0 diags (Unknown)",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "Slider",
			},
			wantDiagCount: 0,
		},
		{
			name: "Unconstrained Carousel with generic range props -> 1 diag",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "Carousel",
				Attributes: map[string]string{
					"min":   "0",
					"max":   "10",
					"value": "activeSlide",
				},
			},
			wantDiagCount: 1,
		},
		{
			name: "Constrained Carousel with generic range props (h-64) -> 0 diags",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "Carousel",
				Classes: []string{"h-64"},
				Attributes: map[string]string{
					"min":   "0",
					"max":   "10",
					"value": "activeSlide",
				},
			},
			wantDiagCount: 0,
		},
		{
			name: "Unconstrained Carousel with spread props -> 1 diag",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "Carousel",
				Attributes: map[string]string{
					"{...props}": "",
				},
			},
			wantDiagCount: 1,
		},
		{
			name: "Unconstrained BannerSlider with onChange -> 1 diag",
			node: &ir.Node{
				Type: ir.NodeElement,
				Tag:  "BannerSlider",
				Attributes: map[string]string{
					"onChange": "handleSlideChange",
				},
			},
			wantDiagCount: 1,
		},
		{
			name: "Generic Slider with snap track and unconstrained child -> 1 diag",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "Slider",
				Classes: []string{"overflow-x-auto", "snap-x"},
				Children: []*ir.Node{
					{
						Type: ir.NodeElement,
						Tag:  "Slide",
					},
				},
			},
			wantDiagCount: 1,
		},
		{
			name: "Generic Slider with snap track and constrained child (aspect-video) -> 0 diags",
			node: &ir.Node{
				Type:    ir.NodeElement,
				Tag:     "Slider",
				Classes: []string{"overflow-x-auto", "snap-x"},
				Children: []*ir.Node{
					{
						Type:    ir.NodeElement,
						Tag:     "Slide",
						Classes: []string{"aspect-video"},
					},
				},
			},
			wantDiagCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := rule.Evaluate(tc.node)
			if len(diags) != tc.wantDiagCount {
				t.Errorf("Evaluate() got %d diagnostics, want %d", len(diags), tc.wantDiagCount)
			}
		})
	}
}
