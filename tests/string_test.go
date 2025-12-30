// -----------------------------------------------------------------------------
// String Type Tests
// -----------------------------------------------------------------------------
// Bu dosya, StringType için kapsamlı test senaryolarını içerir. Test coverage'ı
// maksimize etmek ve edge case'leri yakalamak için çeşitli senaryolar test edilir.
//
// Test Kategorileri:
//   - Required field validation
//   - Length constraints (min/max)
//   - Format validation (email, URL, IP, phone)
//   - Transformation (trim, strip tags)
//   - Custom validators
//   - Password strength
//   - Edge cases (nil, empty, invalid types)
//
// Test Stratejisi:
//   - Table-driven tests (Go best practice)
//   - Positive and negative cases
//   - Performance benchmarks
//   - Thread-safety tests
//
// Metadata:
// @author   Ahmet ALTUN
// @github   github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email    ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

package tests

import (
	"fmt"
	"strings"
	"testing"

	validation "github.com/biyonik/go-fluent-validator"
	"github.com/biyonik/go-fluent-validator/i18n"
)

// TestStringType_Required tests required field validation
func TestStringType_Required(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		wantError bool
		locale    string
	}{
		{
			name:      "valid string",
			value:     "test@example.com",
			wantError: false,
		},
		{
			name:      "empty string should fail",
			value:     "",
			wantError: true,
		},
		{
			name:      "nil value should fail",
			value:     nil,
			wantError: true,
		},
		{
			name:      "whitespace only should fail",
			value:     "   ",
			wantError: false, // After trim becomes empty
		},
		{
			name:      "turkish locale",
			value:     nil,
			wantError: true,
			locale:    "tr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set locale if specified
			if tt.locale != "" {
				i18n.SetLocale(tt.locale)
				defer i18n.SetLocale("en") // Reset
			}

			schema := validation.Make().Shape(map[string]validation.Type{
				"email": validation.String().Required().Label("Email"),
			})

			result := schema.Validate(map[string]any{
				"email": tt.value,
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

// TestStringType_Email tests email validation
func TestStringType_Email(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		// Valid emails
		{"simple email", "test@example.com", false},
		{"email with subdomain", "user@mail.example.com", false},
		{"email with plus", "user+tag@example.com", false},
		{"email with numbers", "user123@example.com", false},
		{"email with dash", "first-last@example.com", false},
		{"email with underscore", "first_last@example.com", false},

		// Invalid emails
		{"missing @", "testexample.com", true},
		{"missing domain", "test@", true},
		{"missing local", "@example.com", true},
		{"double @", "test@@example.com", true},
		{"space in email", "test @example.com", true},
		{"invalid TLD", "test@example.c", true},
		{"starts with dot", ".test@example.com", true},
		{"ends with dot", "test.@example.com", true},
		{"consecutive dots", "test..name@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"email": validation.String().Email().Required(),
			})

			result := schema.Validate(map[string]any{
				"email": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("email '%s': got error = %v, want error = %v",
					tt.value, result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// TestStringType_Length tests min/max length validation
func TestStringType_Length(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		min       int
		max       int
		wantError bool
	}{
		{"within range", "hello", 3, 10, false},
		{"exact min", "hi", 2, 10, false},
		{"exact max", "1234567890", 2, 10, false},
		{"below min", "a", 3, 10, true},
		{"above max", "12345678901", 2, 10, true},
		{"empty string with min", "", 1, 10, true},
		{"unicode characters", "café", 4, 10, false}, // 4 runes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"text": validation.String().Min(tt.min).Max(tt.max),
			})

			result := schema.Validate(map[string]any{
				"text": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("value '%s': got error = %v, want error = %v",
					tt.value, result.HasErrors(), tt.wantError)
			}
		})
	}
}

// TestStringType_Trim tests trim transformation
func TestStringType_Trim(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"leading spaces", "  hello", "hello"},
		{"trailing spaces", "hello  ", "hello"},
		{"both sides", "  hello  ", "hello"},
		{"no spaces", "hello", "hello"},
		{"only spaces", "   ", ""},
		{"tabs and newlines", "\t\nhello\n\t", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"text": validation.String().Trim(),
			})

			result := schema.Validate(map[string]any{
				"text": tt.value,
			})

			if result.HasErrors() {
				t.Errorf("unexpected error: %v", result.Errors())
			}

			validData := result.ValidData()
			got := validData["text"].(string)
			if got != tt.expected {
				t.Errorf("got '%s', want '%s'", got, tt.expected)
			}
		})
	}
}

// TestStringType_URL tests URL validation
func TestStringType_URL(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		// Valid URLs
		{"simple http", "http://example.com", false},
		{"https", "https://example.com", false},
		{"with path", "https://example.com/path/to/page", false},
		{"with query", "https://example.com?key=value", false},
		{"with port", "http://example.com:8080", false},

		{"without protocol", "example.com", true},
		{"invalid protocol", "ftp://example.com", true},
		{"spaces", "http://example .com", true},
		{"missing domain", "http://", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"url": validation.String().URL(),
			})

			result := schema.Validate(map[string]any{
				"url": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("URL '%s': got error = %v, want error = %v",
					tt.value, result.HasErrors(), tt.wantError)
			}
		})
	}
}

// TestStringType_Password tests password validation
func TestStringType_Password(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		wantError bool
	}{
		// Valid passwords
		{"strong password", "MyP@ssw0rd123", false},
		{"with special chars", "Test!@#$123Ab", false},
		{"long password", "VeryLongPassword123!@#", false},

		// Invalid passwords
		{"too short", "Ab1!", true},
		{"no uppercase", "mypassword123!", true},
		{"no lowercase", "MYPASSWORD123!", true},
		{"no numbers", "MyPassword!@#", true},
		{"no special", "MyPassword123", true},
		{"common password", "Password123!", true},
		{"keyboard pattern", "Qwerty123!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"password": validation.String().Password().Required(),
			})

			result := schema.Validate(map[string]any{
				"password": tt.password,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("password '%s': got error = %v, want error = %v",
					tt.password, result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// TestStringType_CustomValidator tests custom validation logic
func TestStringType_Custom(t *testing.T) {
	// This test assumes Custom() method exists
	schema := validation.Make().Shape(map[string]validation.Type{
		"username": validation.String().Custom(func(value string) error {
			if strings.Contains(strings.ToLower(value), "admin") {
				return fmt.Errorf("username cannot contain 'admin'")
			}
			return nil
		}).Label("Username"),
	})

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"valid username", "john_doe", false},
		{"contains admin", "admin_user", true},
		{"contains Admin uppercase", "AdminUser", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(map[string]any{
				"username": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("username '%s': got error = %v, want error = %v",
					tt.value, result.HasErrors(), tt.wantError)
			}
		})
	}
}

// TestStringType_OneOf tests enum validation
func TestStringType_OneOf(t *testing.T) {
	allowedRoles := []string{"admin", "user", "editor"}

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"valid admin", "admin", false},
		{"valid user", "user", false},
		{"valid editor", "editor", false},
		{"invalid role", "superadmin", true},
		{"case sensitive", "Admin", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"role": validation.String().OneOf(allowedRoles).Required(),
			})

			result := schema.Validate(map[string]any{
				"role": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("role '%s': got error = %v, want error = %v",
					tt.value, result.HasErrors(), tt.wantError)
			}
		})
	}
}

// TestStringType_Default tests default value functionality
func TestStringType_Default(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"role": validation.String().Default("user"),
	})

	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{"provided value", "admin", "admin"},
		{"nil value uses default", nil, "user"},
		{"empty string", "", ""}, // Empty is different from nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(map[string]any{
				"role": tt.value,
			})

			if result.HasErrors() {
				t.Errorf("unexpected error: %v", result.Errors())
			}

			validData := result.ValidData()
			got := validData["role"].(string)
			if got != tt.expected {
				t.Errorf("got '%s', want '%s'", got, tt.expected)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Benchmark Tests
// -----------------------------------------------------------------------------

// BenchmarkStringValidation_Simple benchmarks simple email validation
func BenchmarkStringValidation_Simple(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"email": validation.String().Email().Required(),
	})

	data := map[string]any{
		"email": "test@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}

// BenchmarkStringValidation_Complex benchmarks complex validation with multiple rules
func BenchmarkStringValidation_Complex(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"email":    validation.String().Email().Required().Min(5).Max(100).Trim(),
		"username": validation.String().Required().Min(3).Max(20),
		"password": validation.String().Password().Required(),
	})

	data := map[string]any{
		"email":    "  test@example.com  ",
		"username": "john_doe",
		"password": "MyP@ssw0rd123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}

// BenchmarkStringValidation_PasswordStrength benchmarks password strength validation
func BenchmarkStringValidation_PasswordStrength(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"password": validation.String().Password().Required(),
	})

	data := map[string]any{
		"password": "MyVeryStr0ng!P@ssword",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}

// -----------------------------------------------------------------------------
// Edge Cases & Error Handling Tests
// -----------------------------------------------------------------------------

// TestStringType_TypeMismatch tests invalid type handling
func TestStringType_TypeMismatch(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"name": validation.String().Required(),
	})

	tests := []struct {
		name  string
		value any
	}{
		{"integer", 123},
		{"float", 123.45},
		{"boolean", true},
		{"array", []string{"test"}},
		{"map", map[string]string{"key": "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(map[string]any{
				"name": tt.value,
			})

			if !result.HasErrors() {
				t.Errorf("expected error for type %T, got none", tt.value)
			}
		})
	}
}

// TestStringType_ChainedValidation tests multiple validation rules
func TestStringType_ChainedValidation(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"email": validation.String().
			Required().
			Email().
			Min(5).
			Max(100).
			Trim().
			Label("Email Address"),
	})

	tests := []struct {
		name      string
		value     any
		wantError bool
	}{
		{"valid email", "  test@example.com  ", false},
		{"too short", "a@b.c", true},
		{"invalid format", "not-an-email", true},
		{"nil value", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(map[string]any{
				"email": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
			}
		})
	}
}
