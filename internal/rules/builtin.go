package rules

import (
	"github.com/will2469/charites/internal/rules/theme"
)

// Invariant static compile-time check: HardcodeOpacityColorRule wajib mengimplementasikan interface Rule.
var _ Rule = (*theme.HardcodeOpacityColorRule)(nil)

func init() {
	_ = Register(theme.NewHardcodeOpacityColorRule())
}

// RegisterBuiltinRules mendaftarkan seluruh built-in rule Charites ke registry yang diberikan.
func RegisterBuiltinRules(reg *Registry) error {
	if reg == nil {
		return ErrNilRule
	}
	return reg.Register(theme.NewHardcodeOpacityColorRule())
}
