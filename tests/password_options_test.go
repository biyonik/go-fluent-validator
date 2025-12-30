// -----------------------------------------------------------------------------
// Password Options Comprehensive Tests
// -----------------------------------------------------------------------------
// Enterprise-level password security testing
// -----------------------------------------------------------------------------

package tests

import (
	"fmt"
	"testing"

	validation "github.com/biyonik/go-fluent-validator"
	"github.com/biyonik/go-fluent-validator/types"
)

// TestPassword_AllOptions tests all password options combined
func TestPassword_AllOptions(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		wantError bool
	}{
		{
			name:      "perfect password",
			password:  "MyStr0ng!P@ssw0rd2024",
			wantError: false,
		},
		{
			name:      "too short",
			password:  "Ab1!",
			wantError: true,
		},
		{
			name:      "too long",
			password:  "Ab1!" + string(make([]byte, 130)),
			wantError: true,
		},
		{
			name:      "no uppercase",
			password:  "mystr0ng!p@ssw0rd",
			wantError: true,
		},
		{
			name:      "no lowercase",
			password:  "MYSTR0NG!P@SSW0RD",
			wantError: true,
		},
		{
			name:      "no numbers",
			password:  "MyStrong!P@ssword",
			wantError: true,
		},
		{
			name:      "no special chars",
			password:  "MyStr0ngPassw0rd",
			wantError: true,
		},
		{
			name:      "common password",
			password:  "Password123!",
			wantError: true,
		},
		{
			name:      "keyboard pattern",
			password:  "Qwerty123!",
			wantError: true,
		},
		{
			name:      "too many repeating",
			password:  "Aaaaa123!",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"password": validation.String().Password(
					types.WithMinLength(8),
					types.WithMaxLength(128),
					types.WithRequireUppercase(true),
					types.WithRequireLowercase(true),
					types.WithRequireNumeric(true),
					types.WithRequireSpecial(true),
					types.WithMinUniqueChars(6),
					types.WithRejectCommon(true),
					types.WithCheckKeyboardPatterns(true),
					types.WithMaxRepeatingChars(2),
					types.WithMinEntropy(50.0),
				).Required().Label("Password"),
			})

			result := schema.Validate(map[string]any{
				"password": tt.password,
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

// TestPassword_CommonPasswords tests common password rejection
func TestPassword_CommonPasswords(t *testing.T) {
	commonPasswords := []string{
		"password",
		"123456",
		"qwerty",
		"admin",
		"letmein",
		"welcome",
		"password123",
		"12345678",
	}

	schema := validation.Make().Shape(map[string]validation.Type{
		"password": validation.String().Password(
			types.WithRejectCommon(true),
			types.WithMinLength(6),
		).Required(),
	})

	for _, password := range commonPasswords {
		t.Run("reject "+password, func(t *testing.T) {
			result := schema.Validate(map[string]any{
				"password": password,
			})

			if !result.HasErrors() {
				t.Errorf("common password '%s' should be rejected", password)
			}
		})
	}
}

// TestPassword_KeyboardPatterns tests keyboard pattern detection
func TestPassword_KeyboardPatterns(t *testing.T) {
	keyboardPatterns := []struct {
		password string
		pattern  string
	}{
		{"Qwerty123!", "qwerty"},
		{"Asdfgh123!", "asdfgh"},
		{"123456Abc!", "123456"},
		{"Zxcvbn123!", "zxcvbn"},
	}

	schema := validation.Make().Shape(map[string]validation.Type{
		"password": validation.String().Password(
			types.WithCheckKeyboardPatterns(true),
			types.WithMinLength(6),
			types.WithRequireUppercase(true),
			types.WithRequireLowercase(true),
			types.WithRequireNumeric(true),
			types.WithRequireSpecial(true),
		).Required(),
	})

	for _, tt := range keyboardPatterns {
		t.Run("reject "+tt.pattern, func(t *testing.T) {
			result := schema.Validate(map[string]any{
				"password": tt.password,
			})

			if !result.HasErrors() {
				t.Errorf("password with keyboard pattern '%s' should be rejected", tt.pattern)
			}
		})
	}
}

// TestPassword_EntropyCalculation tests password entropy validation
func TestPassword_EntropyCalculation(t *testing.T) {
	tests := []struct {
		name       string
		password   string
		minEntropy float64
		wantError  bool
	}{
		{
			name:       "high entropy - strong password",
			password:   "Xy9$mK2#pL8@",
			minEntropy: 50.0,
			wantError:  false,
		},
		{
			name:       "low entropy - weak password",
			password:   "aaaaa",
			minEntropy: 50.0,
			wantError:  true,
		},
		{
			name:       "medium entropy - acceptable",
			password:   "HelloWorld123",
			minEntropy: 40.0,
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"password": validation.String().Password(
					types.WithMinEntropy(tt.minEntropy),
				).Required(),
			})

			result := schema.Validate(map[string]any{
				"password": tt.password,
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

// TestPassword_RealWorld_UserRegistration tests realistic user registration
func TestPassword_RealWorld_UserRegistration(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"password": validation.String().Password(
			types.WithMinLength(10),
			types.WithMaxLength(128),
			types.WithRequireUppercase(true),
			types.WithRequireLowercase(true),
			types.WithRequireNumeric(true),
			types.WithRequireSpecial(true),
			types.WithMinUniqueChars(5),
			types.WithRejectCommon(true),
			types.WithCheckKeyboardPatterns(true),
		).Required().Label("Password"),
		"password_confirm": validation.String().Required().Label("Confirm Password"),
	}).CrossValidate(func(data map[string]any) error {
		pass := data["password"].(string)
		confirm := data["password_confirm"].(string)
		if pass != confirm {
			return fmt.Errorf("passwords do not match")
		}
		return nil
	})

	t.Run("valid registration", func(t *testing.T) {
		result := schema.Validate(map[string]any{
			"password":         "MyS3cur3!P@ss",
			"password_confirm": "MyS3cur3!P@ss",
		})

		if result.HasErrors() {
			t.Errorf("valid password should not have errors: %v", result.Errors())
		}
	})

	t.Run("passwords don't match", func(t *testing.T) {
		result := schema.Validate(map[string]any{
			"password":         "MyS3cur3!P@ss",
			"password_confirm": "DifferentP@ss1",
		})

		if !result.HasErrors() {
			t.Error("mismatched passwords should have errors")
		}
	})

	t.Run("weak password", func(t *testing.T) {
		result := schema.Validate(map[string]any{
			"password":         "weak",
			"password_confirm": "weak",
		})

		if !result.HasErrors() {
			t.Error("weak password should have errors")
		}
	})
}

// BenchmarkPassword_FullValidation benchmarks complete password validation
func BenchmarkPassword_FullValidation(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"password": validation.String().Password(
			types.WithMinLength(10),
			types.WithRequireUppercase(true),
			types.WithRequireLowercase(true),
			types.WithRequireNumeric(true),
			types.WithRequireSpecial(true),
			types.WithRejectCommon(true),
			types.WithCheckKeyboardPatterns(true),
		).Required(),
	})

	data := map[string]any{
		"password": "MyS3cur3!P@ssw0rd2024",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}
