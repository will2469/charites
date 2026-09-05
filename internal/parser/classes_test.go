package parser_test

import (
	"reflect"
	"testing"

	"github.com/will2469/charites/internal/parser"
)

func TestExtractClasses(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantClasses []string
		wantDynamic bool
	}{
		{
			name:        "empty string",
			input:       "",
			wantClasses: nil,
			wantDynamic: false,
		},
		{
			name:        "whitespace only",
			input:       "   \t\n",
			wantClasses: nil,
			wantDynamic: false,
		},
		{
			name:        "double quoted static string",
			input:       `"p-4 bg-primary text-sm"`,
			wantClasses: []string{"p-4", "bg-primary", "text-sm"},
			wantDynamic: false,
		},
		{
			name:        "single quoted static string",
			input:       `'flex flex-col items-center'`,
			wantClasses: []string{"flex", "flex-col", "items-center"},
			wantDynamic: false,
		},
		{
			name:        "unquoted string",
			input:       `btn btn-primary`,
			wantClasses: []string{"btn", "btn-primary"},
			wantDynamic: false,
		},
		{
			name:        "jsx double quoted in braces",
			input:       `{"w-full h-full"}`,
			wantClasses: []string{"w-full", "h-full"},
			wantDynamic: false,
		},
		{
			name:        "jsx single quoted in braces",
			input:       `{'w-full h-full'}`,
			wantClasses: []string{"w-full", "h-full"},
			wantDynamic: false,
		},
		{
			name:        "pure static template literal in braces",
			input:       "{`grid grid-cols-2 gap-4`}",
			wantClasses: []string{"grid", "grid-cols-2", "gap-4"},
			wantDynamic: false,
		},
		{
			name:        "direct template literal",
			input:       "`grid grid-cols-2 gap-4`",
			wantClasses: []string{"grid", "grid-cols-2", "gap-4"},
			wantDynamic: false,
		},
		{
			name:        "dynamic template literal with static prefixes and suffixes",
			input:       "`p-4 ${foo ? \"bg-red\" : \"bg-blue\"} text-sm`",
			wantClasses: []string{"p-4", "text-sm"},
			wantDynamic: true,
		},
		{
			name:        "jsx dynamic template literal in braces",
			input:       "{`p-4 ${foo ? 'bg-red' : 'bg-blue'} text-sm`}",
			wantClasses: []string{"p-4", "text-sm"},
			wantDynamic: true,
		},
		{
			name:        "template literal with multiple dynamic expressions",
			input:       "`flex ${active ? 'opacity-100' : 'opacity-50'} gap-2 ${size} font-bold`",
			wantClasses: []string{"flex", "gap-2", "font-bold"},
			wantDynamic: true,
		},
		{
			name:        "only dynamic expression in template literal",
			input:       "`${customClass}`",
			wantClasses: nil,
			wantDynamic: true,
		},
		{
			name:        "dynamic variable in jsx braces",
			input:       "{computedClassName}",
			wantClasses: nil,
			wantDynamic: true,
		},
		{
			name:        "nested braces inside template expression",
			input:       "`card ${obj[{key: 1}.key]} p-6`",
			wantClasses: []string{"card", "p-6"},
			wantDynamic: true,
		},
		{
			name:        "unclosed dynamic expression",
			input:       "`card ${unclosed p-6`",
			wantClasses: []string{"card"},
			wantDynamic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClasses, gotDynamic := parser.ExtractClasses(tt.input)
			if gotDynamic != tt.wantDynamic {
				t.Errorf("ExtractClasses(%q) gotDynamic = %v, want %v", tt.input, gotDynamic, tt.wantDynamic)
			}
			if len(gotClasses) == 0 && len(tt.wantClasses) == 0 {
				return
			}
			if !reflect.DeepEqual(gotClasses, tt.wantClasses) {
				t.Errorf("ExtractClasses(%q) gotClasses = %v, want %v", tt.input, gotClasses, tt.wantClasses)
			}
		})
	}
}
