package tests

import (
	"testing"
	"time"

	v "github.com/biyonik/go-fluent-validator"
	"github.com/biyonik/go-fluent-validator/core"
)

// TestUuidValidation tests UUID validation
func TestUuidValidation(t *testing.T) {
	tests := []struct {
		name      string
		uuid      string
		version   int
		shouldErr bool
	}{
		{"Valid UUID v4", "550e8400-e29b-41d4-a716-446655440000", 4, false},
		{"Valid UUID v4 uppercase", "550E8400-E29B-41D4-A716-446655440000", 4, false},
		{"Valid UUID v1", "a0eebc99-9c0b-1ef9-bb6d-00c04fd430c8", 1, false},
		{"Invalid UUID", "not-a-uuid", 4, true},
		{"Empty UUID", "", 4, true},
		{"Valid any version", "550e8400-e29b-41d4-a716-446655440000", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := v.Make().Shape(map[string]v.Type{
				"id": v.Uuid().Version(tt.version).Required(),
			})

			result := schema.Validate(map[string]any{
				"id": tt.uuid,
			})

			if tt.shouldErr && !result.HasErrors() {
				t.Errorf("Expected error for UUID %q but got none", tt.uuid)
			}
			if !tt.shouldErr && result.HasErrors() {
				t.Errorf("Expected no error for UUID %q but got: %v", tt.uuid, result.Errors())
			}
		})
	}
}

// TestUuidCustomValidation tests UUID custom validators
func TestUuidCustomValidation(t *testing.T) {
	schema := v.Make().Shape(map[string]v.Type{
		"id": v.Uuid().Version(4).Custom(func(uuid string) error {
			if uuid == "550e8400-e29b-41d4-a716-446655440000" {
				return nil
			}
			return core.NewValidationError("UUID must be a specific value")
		}),
	})

	// Valid custom check
	result := schema.Validate(map[string]any{
		"id": "550e8400-e29b-41d4-a716-446655440000",
	})
	if result.HasErrors() {
		t.Errorf("Expected no error but got: %v", result.Errors())
	}

	// Invalid custom check
	result = schema.Validate(map[string]any{
		"id": "a0eebc99-9c0b-4ef9-bb6d-00c04fd430c8",
	})
	if !result.HasErrors() {
		t.Error("Expected custom validation error but got none")
	}
}

// TestIbanValidation tests IBAN validation
func TestIbanValidation(t *testing.T) {
	tests := []struct {
		name      string
		iban      string
		country   string
		shouldErr bool
	}{
		{"Valid TR IBAN", "TR330006100519786457841326", "TR", false},
		{"Valid DE IBAN", "DE89370400440532013000", "DE", false},
		{"Invalid IBAN", "INVALID123", "", true},
		{"Empty IBAN", "", "", true},
		{"Wrong country code", "DE89370400440532013000", "TR", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ibanValidator := v.Iban().Required()
			if tt.country != "" {
				ibanValidator = ibanValidator.Country(tt.country)
			}

			schema := v.Make().Shape(map[string]v.Type{
				"account": ibanValidator,
			})

			result := schema.Validate(map[string]any{
				"account": tt.iban,
			})

			if tt.shouldErr && !result.HasErrors() {
				t.Errorf("Expected error for IBAN %q but got none", tt.iban)
			}
			if !tt.shouldErr && result.HasErrors() {
				t.Errorf("Expected no error for IBAN %q but got: %v", tt.iban, result.Errors())
			}
		})
	}
}

// TestIbanCustomValidation tests IBAN custom validators
func TestIbanCustomValidation(t *testing.T) {
	schema := v.Make().Shape(map[string]v.Type{
		"account": v.Iban().Custom(func(iban string) error {
			if len(iban) < 20 {
				return core.NewValidationError("IBAN too short")
			}
			return nil
		}),
	})

	result := schema.Validate(map[string]any{
		"account": "TR330006100519786457841326",
	})
	if result.HasErrors() {
		t.Errorf("Expected no error but got: %v", result.Errors())
	}

	result = schema.Validate(map[string]any{
		"account": "SHORT",
	})
	if !result.HasErrors() {
		t.Error("Expected custom validation error but got none")
	}
}

// TestCreditCardValidation tests credit card validation
func TestCreditCardValidation(t *testing.T) {
	tests := []struct {
		name      string
		card      string
		cardType  string
		shouldErr bool
	}{
		{"Valid Visa", "4532015112830366", "visa", false},
		{"Valid MasterCard", "5425233430109903", "mastercard", false},
		{"Valid Amex", "374245455400126", "amex", false},
		{"Invalid card", "1234567890123456", "", true},
		{"Empty card", "", "", true},
		{"Wrong type", "4532015112830366", "mastercard", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardValidator := v.CreditCard().Required()
			if tt.cardType != "" {
				cardValidator = cardValidator.Type(tt.cardType)
			}

			schema := v.Make().Shape(map[string]v.Type{
				"card": cardValidator,
			})

			result := schema.Validate(map[string]any{
				"card": tt.card,
			})

			if tt.shouldErr && !result.HasErrors() {
				t.Errorf("Expected error for card %q but got none", tt.card)
			}
			if !tt.shouldErr && result.HasErrors() {
				t.Errorf("Expected no error for card %q but got: %v", tt.card, result.Errors())
			}
		})
	}
}

// TestCreditCardCustomValidation tests credit card custom validators
func TestCreditCardCustomValidation(t *testing.T) {
	schema := v.Make().Shape(map[string]v.Type{
		"card": v.CreditCard().Custom(func(card string) error {
			if card[0] != '4' && card[0] != '5' {
				return core.NewValidationError("Only Visa and MasterCard accepted")
			}
			return nil
		}),
	})

	result := schema.Validate(map[string]any{
		"card": "4532015112830366",
	})
	if result.HasErrors() {
		t.Errorf("Expected no error but got: %v", result.Errors())
	}

	result = schema.Validate(map[string]any{
		"card": "374245455400126", // Amex starting with 3
	})
	if !result.HasErrors() {
		t.Error("Expected custom validation error but got none")
	}
}

// TestDateValidation tests date validation
func TestDateValidation(t *testing.T) {
	tests := []struct {
		name      string
		date      string
		format    string
		minDate   string
		maxDate   string
		shouldErr bool
	}{
		{"Valid date", "2024-01-15", "2006-01-02", "", "", false},
		{"Invalid format", "15-01-2024", "2006-01-02", "", "", true},
		{"Date too early", "2020-01-01", "2006-01-02", "2021-01-01", "", true},
		{"Date too late", "2025-01-01", "2006-01-02", "", "2024-12-31", true},
		{"Date in range", "2023-06-15", "2006-01-02", "2023-01-01", "2024-01-01", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dateValidator := v.Date().Required()
			if tt.format != "" {
				dateValidator = dateValidator.Format(tt.format)
			}
			if tt.minDate != "" {
				dateValidator = dateValidator.Min(tt.minDate)
			}
			if tt.maxDate != "" {
				dateValidator = dateValidator.Max(tt.maxDate)
			}

			schema := v.Make().Shape(map[string]v.Type{
				"date": dateValidator,
			})

			result := schema.Validate(map[string]any{
				"date": tt.date,
			})

			if tt.shouldErr && !result.HasErrors() {
				t.Errorf("Expected error for date %q but got none", tt.date)
			}
			if !tt.shouldErr && result.HasErrors() {
				t.Errorf("Expected no error for date %q but got: %v", tt.date, result.Errors())
			}
		})
	}
}

// TestDateCustomValidation tests date custom validators
func TestDateCustomValidation(t *testing.T) {
	schema := v.Make().Shape(map[string]v.Type{
		"date": v.Date().Format("2006-01-02").Custom(func(date time.Time) error {
			if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
				return core.NewValidationError("Date cannot be on weekend")
			}
			return nil
		}),
	})

	// Monday
	result := schema.Validate(map[string]any{
		"date": "2024-01-15", // This is a Monday
	})
	if result.HasErrors() {
		t.Errorf("Expected no error for weekday but got: %v", result.Errors())
	}

	// Saturday
	result = schema.Validate(map[string]any{
		"date": "2024-01-13", // This is a Saturday
	})
	if !result.HasErrors() {
		t.Error("Expected custom validation error for weekend but got none")
	}
}

// TestObjectValidation tests object (nested) validation
func TestObjectValidation(t *testing.T) {
	schema := v.Make().Shape(map[string]v.Type{
		"user": v.Object().Shape(map[string]v.Type{
			"name":  v.String().Required().Min(3),
			"email": v.String().Email().Required(),
			"age":   v.Number().Min(18).Required(),
		}),
	})

	// Valid object
	result := schema.Validate(map[string]any{
		"user": map[string]any{
			"name":  "John Doe",
			"email": "john@example.com",
			"age":   25,
		},
	})
	if result.HasErrors() {
		t.Errorf("Expected no error but got: %v", result.Errors())
	}

	// Invalid nested field
	result = schema.Validate(map[string]any{
		"user": map[string]any{
			"name":  "Jo", // Too short
			"email": "john@example.com",
			"age":   25,
		},
	})
	if !result.HasErrors() {
		t.Error("Expected validation error for short name but got none")
	}
}

// TestObjectCustomValidation tests object custom validators
func TestObjectCustomValidation(t *testing.T) {
	schema := v.Make().Shape(map[string]v.Type{
		"user": v.Object().Shape(map[string]v.Type{
			"firstName": v.String().Required(),
			"lastName":  v.String().Required(),
		}).Custom(func(obj map[string]any) error {
			firstName, _ := obj["firstName"].(string)
			lastName, _ := obj["lastName"].(string)
			if firstName == lastName {
				return core.NewValidationError("First and last name cannot be the same")
			}
			return nil
		}),
	})

	result := schema.Validate(map[string]any{
		"user": map[string]any{
			"firstName": "John",
			"lastName":  "Doe",
		},
	})
	if result.HasErrors() {
		t.Errorf("Expected no error but got: %v", result.Errors())
	}

	result = schema.Validate(map[string]any{
		"user": map[string]any{
			"firstName": "John",
			"lastName":  "John",
		},
	})
	if !result.HasErrors() {
		t.Error("Expected custom validation error but got none")
	}
}

// TestCrossValidationTiming tests that cross-validation runs even when field validation fails
func TestCrossValidationTiming(t *testing.T) {
	schema := v.Make().Shape(map[string]v.Type{
		"password":        v.String().Required().Min(8),
		"passwordConfirm": v.String().Required(),
	}).CrossValidate(func(data map[string]any) error {
		pass, _ := data["password"].(string)
		confirm, _ := data["passwordConfirm"].(string)

		if pass != confirm {
			return core.NewValidationError("Passwords do not match")
		}
		return nil
	})

	// Test with valid password but not matching
	result := schema.Validate(map[string]any{
		"password":        "validpass123",
		"passwordConfirm": "different123",
	})

	// Should have cross-validation error
	if !result.HasErrors() {
		t.Error("Expected cross-validation error but got none")
	}
	if _, ok := result.Errors()["_cross_validation"]; !ok {
		t.Error("Expected _cross_validation error key but not found")
	}

	// Test with invalid password (too short) AND not matching
	result = schema.Validate(map[string]any{
		"password":        "short", // Less than 8 characters
		"passwordConfirm": "different",
	})

	// Should have both field error AND cross-validation error
	if !result.HasErrors() {
		t.Error("Expected errors but got none")
	}
	if _, ok := result.Errors()["password"]; !ok {
		t.Error("Expected password field error but not found")
	}
	if _, ok := result.Errors()["_cross_validation"]; !ok {
		t.Error("Expected _cross_validation error even with field errors, but not found")
	}
}

// TestDomainValidation tests RFC 1035 compliant domain validation
func TestDomainValidation(t *testing.T) {
	tests := []struct {
		name      string
		domain    string
		shouldErr bool
	}{
		{"Valid domain", "example.com", false},
		{"Valid subdomain", "www.example.com", false},
		{"Valid with hyphen", "my-site.com", false},
		{"Invalid - starts with hyphen", "-example.com", true},
		{"Invalid - ends with hyphen", "example-.com", true},
		{"Invalid - single char TLD", "example.c", true},
		{"Valid - two char TLD", "example.co", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := v.Make().Shape(map[string]v.Type{
				"domain": v.AdvancedString().Domain(true).Required(),
			})

			result := schema.Validate(map[string]any{
				"domain": tt.domain,
			})

			if tt.shouldErr && !result.HasErrors() {
				t.Errorf("Expected error for domain %q but got none", tt.domain)
			}
			if !tt.shouldErr && result.HasErrors() {
				t.Errorf("Expected no error for domain %q but got: %v", tt.domain, result.Errors())
			}
		})
	}
}
