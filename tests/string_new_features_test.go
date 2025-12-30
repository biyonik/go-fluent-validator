// -----------------------------------------------------------------------------
// String New Features Tests
// -----------------------------------------------------------------------------
// Tests for Length(), NotOneOf(), and IP() string version
// -----------------------------------------------------------------------------

package tests

import (
	"testing"

	validation "github.com/biyonik/go-fluent-validator"
)

// TestString_Length tests exact length validation
func TestString_Length(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		length    int
		wantError bool
	}{
		{
			name:      "exact length - valid",
			value:     "12345",
			length:    5,
			wantError: false,
		},
		{
			name:      "too short",
			value:     "1234",
			length:    5,
			wantError: true,
		},
		{
			name:      "too long",
			value:     "123456",
			length:    5,
			wantError: true,
		},
		{
			name:      "ZIP code - 5 digits",
			value:     "90210",
			length:    5,
			wantError: false,
		},
		{
			name:      "postal code - invalid",
			value:     "9021",
			length:    5,
			wantError: true,
		},
		{
			name:      "empty string",
			value:     "",
			length:    5,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"code": validation.String().
					Length(tt.length).
					Label("Code"),
			})

			result := schema.Validate(map[string]any{
				"code": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// TestString_NotOneOf tests blacklist validation
func TestString_NotOneOf(t *testing.T) {
	blockedUsernames := []string{"admin", "root", "system", "administrator"}

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{
			name:      "allowed username",
			value:     "john_doe",
			wantError: false,
		},
		{
			name:      "blocked - admin",
			value:     "admin",
			wantError: true,
		},
		{
			name:      "blocked - root",
			value:     "root",
			wantError: true,
		},
		{
			name:      "blocked - system",
			value:     "system",
			wantError: true,
		},
		{
			name:      "allowed - admin123 (not exact match)",
			value:     "admin123",
			wantError: false,
		},
		{
			name:      "case sensitive - Admin (allowed)",
			value:     "Admin",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"username": validation.String().
					Required().
					NotOneOf(blockedUsernames).
					Label("Username"),
			})

			result := schema.Validate(map[string]any{
				"username": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// TestString_IP_StringVersion tests IP validation with string parameters
func TestString_IP_StringVersion(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		version   string
		wantError bool
	}{
		// IPv4 tests
		{
			name:      "valid IPv4 with 'v4'",
			value:     "192.168.1.1",
			version:   "v4",
			wantError: false,
		},
		{
			name:      "valid IPv4 with '4'",
			value:     "192.168.1.1",
			version:   "4",
			wantError: false,
		},
		{
			name:      "invalid IPv4 format",
			value:     "256.1.1.1",
			version:   "v4",
			wantError: true,
		},
		{
			name:      "IPv6 when expecting IPv4",
			value:     "2001:0db8::1",
			version:   "v4",
			wantError: true,
		},

		// IPv6 tests
		{
			name:      "valid IPv6 with 'v6'",
			value:     "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			version:   "v6",
			wantError: false,
		},
		{
			name:      "valid IPv6 compressed with '6'",
			value:     "2001:db8::1",
			version:   "6",
			wantError: false,
		},
		{
			name:      "IPv4 when expecting IPv6",
			value:     "192.168.1.1",
			version:   "v6",
			wantError: true,
		},

		// Any version tests
		{
			name:      "IPv4 with empty version",
			value:     "192.168.1.1",
			version:   "",
			wantError: false,
		},
		{
			name:      "IPv6 with empty version",
			value:     "2001:db8::1",
			version:   "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"ip": validation.String().
					IP(tt.version).
					Required().
					Label("IP Address"),
			})

			result := schema.Validate(map[string]any{
				"ip": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// TestString_RealWorld_ZIPCode tests ZIP code validation
func TestString_RealWorld_ZIPCode(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"zip_code": validation.String().
			Required().
			Length(5).
			Numeric().
			Label("ZIP Code"),
	})

	tests := []struct {
		name      string
		zipCode   string
		wantError bool
	}{
		{"valid US ZIP", "90210", false},
		{"valid NYC ZIP", "10001", false},
		{"too short", "9021", true},
		{"too long", "902101", true},
		{"contains letters", "9021A", true},
		{"contains spaces", "90 210", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(map[string]any{
				"zip_code": tt.zipCode,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
			}
		})
	}
}
