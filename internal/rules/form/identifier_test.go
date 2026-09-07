package form_test

import (
	"reflect"
	"testing"

	"github.com/will2469/charites/internal/rules/form"
)

func TestTokenizeIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"postal_code", []string{"postal", "code"}},
		{"postalCode", []string{"postal", "code"}},
		{"postal-code", []string{"postal", "code"}},
		{"noRT", []string{"no", "rt"}},
		{"orderID", []string{"order", "id"}},
		{"XMLHttp", []string{"xml", "http"}},
		{"start_date", []string{"start", "date"}},
		{"alert_count", []string{"alert", "count"}},
		{"php_version", []string{"php", "version"}},
		{"short_code", []string{"short", "code"}},
		{"total_account_number", []string{"total", "account", "number"}},
		{"jumlah_nomor_kk", []string{"jumlah", "nomor", "kk"}},
		{"user.profile.phone_number", []string{"user", "profile", "phone", "number"}},
	}

	for _, tt := range tests {
		got := form.TokenizeIdentifier(tt.input)
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("TokenizeIdentifier(%q) = %v, expected %v", tt.input, got, tt.expected)
		}
	}
}

func TestClassifyIdentifier(t *testing.T) {
	tests := []struct {
		name          string
		identifier    string
		expectedClass form.IdentifierClass
		expectedMatch string
	}{
		// Strong Exact Identity
		{"Exact NIK", "nik", form.IdentifierIdentity, "nik"},
		{"Exact KK", "kk", form.IdentifierIdentity, "kk"},
		{"Exact KTP", "ktp", form.IdentifierIdentity, "ktp"},
		{"Exact NPWP", "npwp", form.IdentifierIdentity, "npwp"},
		{"Exact BPJS", "bpjs", form.IdentifierIdentity, "bpjs"},
		{"Exact PIN", "pin", form.IdentifierIdentity, "pin"},
		{"Exact OTP", "otp", form.IdentifierIdentity, "otp"},
		{"Exact Passport", "passport", form.IdentifierIdentity, "passport"},
		{"Exact SSN", "ssn", form.IdentifierIdentity, "ssn"},
		{"Exact Serial", "serial", form.IdentifierIdentity, "serial"},

		// Strong Qualified Identity Sequences
		{"Postal Code snake", "postal_code", form.IdentifierIdentity, "postal_code"},
		{"Postal Code camel", "postalCode", form.IdentifierIdentity, "postal_code"},
		{"Zip Code", "zip_code", form.IdentifierIdentity, "zip_code"},
		{"Account Number", "account_number", form.IdentifierIdentity, "account_number"},
		{"Order ID", "order_id", form.IdentifierIdentity, "order_id"},
		{"Order Number", "order_number", form.IdentifierIdentity, "order_number"},
		{"Tracking Number", "tracking_number", form.IdentifierIdentity, "tracking_number"},
		{"Invoice Number", "invoice_number", form.IdentifierIdentity, "invoice_number"},
		{"Phone Number", "phone_number", form.IdentifierIdentity, "phone_number"},
		{"Security PIN", "security_pin", form.IdentifierIdentity, "pin"}, // "pin" exact or "security_pin"
		{"Credit Card", "credit_card", form.IdentifierIdentity, "credit_card"},
		{"Nomor KK", "nomor_kk", form.IdentifierIdentity, "kk"}, // "kk" exact or "nomor_kk"

		// Precedence: Strong Identity over Quantity prefix
		{"Total Account Number", "total_account_number", form.IdentifierIdentity, "account_number"},
		{"Batch Order ID", "batch_order_id", form.IdentifierIdentity, "order_id"},
		{"Jumlah Nomor KK", "jumlah_nomor_kk", form.IdentifierIdentity, "kk"},

		// Broad nouns alone are NOT identity
		{"Order alone", "order", form.IdentifierUnknown, ""},
		{"Account alone", "account", form.IdentifierUnknown, ""},
		{"Invoice alone", "invoice", form.IdentifierUnknown, ""},
		{"Order Quantity", "order_quantity", form.IdentifierQuantity, "quantity"},
		{"Order Count", "order_count", form.IdentifierQuantity, "count"},
		{"Account Quantity", "account_quantity", form.IdentifierQuantity, "quantity"},
		{"Account Count", "account_count", form.IdentifierQuantity, "count"},
		{"Invoice Total", "invoice_total", form.IdentifierQuantity, "total"},
		{"Phone Count", "phone_count", form.IdentifierQuantity, "count"},
		{"Postal Count", "postal_count", form.IdentifierQuantity, "count"},

		// Contextual Identity Resolution
		{"RT alone", "rt", form.IdentifierIdentity, "rt"},
		{"No RT", "no_rt", form.IdentifierIdentity, "rt"},
		{"RT Code", "rt_code", form.IdentifierIdentity, "rt"},
		{"RW alone", "rw", form.IdentifierIdentity, "rw"},
		{"HP alone", "hp", form.IdentifierIdentity, "hp"},
		{"WA alone", "wa", form.IdentifierIdentity, "wa"},
		{"RT Count", "rt_count", form.IdentifierQuantity, "count"},
		{"Total RT", "total_rt", form.IdentifierQuantity, "total"},
		{"Jumlah RT", "jumlah_rt", form.IdentifierQuantity, "jumlah"},
		{"WA Count", "wa_count", form.IdentifierQuantity, "count"},

		// Substring Safeguards (Immunity)
		{"Start Date (not rt)", "start_date", form.IdentifierUnknown, ""},
		{"Alert Count (not rt)", "alert_count", form.IdentifierQuantity, "count"},
		{"PHP Version (not hp)", "php_version", form.IdentifierUnknown, ""},
		{"Short Code (not rt)", "short_code", form.IdentifierUnknown, ""},
		{"Heart Rate (not rt)", "heart_rate", form.IdentifierQuantity, "rate"},

		// General Quantities
		{"Quantity", "quantity", form.IdentifierQuantity, "quantity"},
		{"Item Count", "item_count", form.IdentifierQuantity, "count"},
		{"Discount Percent", "discount_percent", form.IdentifierQuantity, "percent"},
		{"Temperature", "temperature", form.IdentifierUnknown, ""}, // Unknown domain, but numeric input
		{"Weight", "package_weight", form.IdentifierQuantity, "weight"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := form.TokenizeIdentifier(tt.identifier)
			gotClass, gotMatch := form.ClassifyIdentifier(tokens)
			if gotClass != tt.expectedClass {
				t.Errorf("ClassifyIdentifier(%v) class = %v, expected %v", tokens, gotClass, tt.expectedClass)
			}
			if tt.expectedMatch != "" && gotMatch != tt.expectedMatch {
				t.Errorf("ClassifyIdentifier(%v) matched = %q, expected %q", tokens, gotMatch, tt.expectedMatch)
			}
		})
	}
}
