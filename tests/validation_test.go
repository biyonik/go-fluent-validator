// -----------------------------------------------------------------------------
// Schema Validation Tests
// -----------------------------------------------------------------------------
// Bu dosya, ValidationSchema için kapsamlı test senaryolarını içerir.
// Cross-field validation, When conditions, ve complex schemas test edilir.
//
// Test Kategorileri:
//   - Basic schema validation
//   - Cross-field validation (password confirmation)
//   - Conditional rules (When)
//   - Nested schemas
//   - Multiple field dependencies
//   - Complex real-world scenarios
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
	"time"

	validation "github.com/biyonik/go-fluent-validator"
)

// -----------------------------------------------------------------------------
// Basic Schema Tests
// -----------------------------------------------------------------------------

// TestSchema_BasicValidation tests simple schema validation
func TestSchema_BasicValidation(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"name":  validation.String().Required().Min(2).Max(50),
		"email": validation.String().Required().Email(),
		"age":   validation.Number().Min(0).Max(150).Integer(),
	})

	tests := []struct {
		name      string
		data      map[string]any
		wantError bool
	}{
		{
			"valid data",
			map[string]any{
				"name":  "John Doe",
				"email": "john@example.com",
				"age":   30,
			},
			false,
		},
		{
			"missing required field",
			map[string]any{
				"email": "john@example.com",
				"age":   30,
			},
			true,
		},
		{
			"invalid email",
			map[string]any{
				"name":  "John Doe",
				"email": "invalid-email",
				"age":   30,
			},
			true,
		},
		{
			"age out of range",
			map[string]any{
				"name":  "John Doe",
				"email": "john@example.com",
				"age":   200,
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(tt.data)

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Cross-Field Validation Tests
// -----------------------------------------------------------------------------

// TestSchema_CrossValidation_PasswordConfirmation tests password confirmation
func TestSchema_CrossValidation_PasswordConfirmation(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"password":         validation.String().Required().Min(8),
		"password_confirm": validation.String().Required(),
	}).CrossValidate(func(data map[string]any) error {
		password, _ := data["password"].(string)
		confirm, _ := data["password_confirm"].(string)

		if password != confirm {
			return fmt.Errorf("passwords do not match")
		}
		return nil
	})

	tests := []struct {
		name      string
		data      map[string]any
		wantError bool
	}{
		{
			"matching passwords",
			map[string]any{
				"password":         "MyPassword123",
				"password_confirm": "MyPassword123",
			},
			false,
		},
		{
			"non-matching passwords",
			map[string]any{
				"password":         "MyPassword123",
				"password_confirm": "DifferentPass",
			},
			true,
		},
		{
			"empty confirmation",
			map[string]any{
				"password":         "MyPassword123",
				"password_confirm": "",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(tt.data)

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// TestSchema_CrossValidation_DateRange tests start/end date validation
func TestSchema_CrossValidation_DateRange(t *testing.T) {
	_ = validation.Make().Shape(map[string]validation.Type{
		"start_date": validation.Date().Required(),
		"end_date":   validation.Date().Required(),
	}).CrossValidate(func(data map[string]any) error {
		startDate, startOk := data["start_date"].(time.Time)
		endDate, endOk := data["end_date"].(time.Time)

		if !startOk || !endOk {
			return fmt.Errorf("invalid date format")
		}

		if endDate.Before(startDate) {
			return fmt.Errorf("end date must be after start date")
		}
		return nil
	})

	// Note: This test requires proper Date type implementation
	// Just demonstrating the pattern here
}

// TestSchema_CrossValidation_MultipleFields tests validation across multiple fields
func TestSchema_CrossValidation_MultipleFields(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"age":       validation.Number().Min(0).Max(150).Integer(),
		"parent_id": validation.Number().Integer(),
	}).CrossValidate(func(data map[string]any) error {
		// Type-safe check
		ageVal, ageOk := data["age"]
		if !ageOk {
			return nil
		}

		// Int veya float64 olabilir
		var age float64
		switch v := ageVal.(type) {
		case int:
			age = float64(v)
		case float64:
			age = v
		default:
			return nil
		}

		parentIDVal, hasParent := data["parent_id"]
		if !hasParent || parentIDVal == nil {
			return nil
		}

		var parentID float64
		switch v := parentIDVal.(type) {
		case int:
			parentID = float64(v)
		case float64:
			parentID = v
		default:
			return nil
		}

		if parentID > 0 && age >= 18 {
			return fmt.Errorf("adults cannot have parent_id set")
		}
		return nil
	})

	tests := []struct {
		name      string
		data      map[string]any
		wantError bool
	}{
		{
			"child with parent",
			map[string]any{
				"age":       10,
				"parent_id": 1,
			},
			false,
		},
		{
			"adult without parent",
			map[string]any{
				"age": 25,
			},
			false,
		},
		{
			"adult with parent - should fail",
			map[string]any{
				"age":       25,
				"parent_id": 1,
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(tt.data)

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Conditional Validation (When) Tests
// -----------------------------------------------------------------------------

// TestSchema_When_PaymentMethod tests conditional validation based on payment method
func TestSchema_When_PaymentMethod(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"payment_method": validation.String().OneOf([]string{"credit_card", "paypal", "bank_transfer"}).Required(),
	}).When("payment_method", "credit_card", func() validation.Schema {
		return validation.Make().Shape(map[string]validation.Type{
			"card_number": validation.CreditCard().Required(),
			"cvv":         validation.String().Min(3).Max(4).Required(),
			"exp_date":    validation.String().Required(),
		})
	}).When("payment_method", "paypal", func() validation.Schema {
		return validation.Make().Shape(map[string]validation.Type{
			"paypal_email": validation.String().Email().Required(),
		})
	}).When("payment_method", "bank_transfer", func() validation.Schema {
		return validation.Make().Shape(map[string]validation.Type{
			"iban": validation.Iban().Required(),
		})
	})

	tests := []struct {
		name      string
		data      map[string]any
		wantError bool
	}{
		{
			"credit card - valid",
			map[string]any{
				"payment_method": "credit_card",
				"card_number":    "4532015112830366",
				"cvv":            "123",
				"exp_date":       "12/25",
			},
			false,
		},
		{
			"credit card - missing cvv",
			map[string]any{
				"payment_method": "credit_card",
				"card_number":    "4532015112830366",
				"exp_date":       "12/25",
			},
			true,
		},
		{
			"paypal - valid",
			map[string]any{
				"payment_method": "paypal",
				"paypal_email":   "user@example.com",
			},
			false,
		},
		{
			"paypal - invalid email",
			map[string]any{
				"payment_method": "paypal",
				"paypal_email":   "invalid-email",
			},
			true,
		},
		{
			"bank transfer - valid",
			map[string]any{
				"payment_method": "bank_transfer",
				"iban":           "TR330006100519786457841326",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(tt.data)

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// TestSchema_When_UserType tests conditional validation based on user type
func TestSchema_When_UserType(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"user_type": validation.String().OneOf([]string{"individual", "corporate"}).Required(),
		"name":      validation.String().Required(),
	}).When("user_type", "corporate", func() validation.Schema {
		return validation.Make().Shape(map[string]validation.Type{
			"company_name": validation.String().Required(),
			"tax_number":   validation.String().Required().Min(10).Max(11),
		})
	}).When("user_type", "individual", func() validation.Schema {
		return validation.Make().Shape(map[string]validation.Type{
			"id_number": validation.String().Required().Min(11).Max(11),
		})
	})

	tests := []struct {
		name      string
		data      map[string]any
		wantError bool
	}{
		{
			"corporate - valid",
			map[string]any{
				"user_type":    "corporate",
				"name":         "John Doe",
				"company_name": "Acme Inc",
				"tax_number":   "1234567890",
			},
			false,
		},
		{
			"corporate - missing tax_number",
			map[string]any{
				"user_type":    "corporate",
				"name":         "John Doe",
				"company_name": "Acme Inc",
			},
			true,
		},
		{
			"individual - valid",
			map[string]any{
				"user_type": "individual",
				"name":      "John Doe",
				"id_number": "12345678901",
			},
			false,
		},
		{
			"individual - invalid id_number length",
			map[string]any{
				"user_type": "individual",
				"name":      "John Doe",
				"id_number": "123456",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(tt.data)

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Complex Real-World Scenarios
// -----------------------------------------------------------------------------

// TestSchema_RealWorld_UserRegistration tests a complete user registration form
func TestSchema_RealWorld_UserRegistration(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"username":         validation.String().Required().Min(3).Max(20),
		"email":            validation.String().Required().Email().Trim(),
		"password":         validation.String().Required().Password(),
		"password_confirm": validation.String().Required(),
		"age":              validation.Number().Min(13).Integer(),
		"terms_accepted":   validation.Boolean().Required(),
	}).CrossValidate(func(data map[string]any) error {
		password, _ := data["password"].(string)
		confirm, _ := data["password_confirm"].(string)
		if password != confirm {
			return fmt.Errorf("passwords must match")
		}

		termsAccepted, _ := data["terms_accepted"].(bool)
		if !termsAccepted {
			return fmt.Errorf("you must accept the terms and conditions")
		}

		return nil
	})

	tests := []struct {
		name      string
		data      map[string]any
		wantError bool
	}{
		{
			"valid registration",
			map[string]any{
				"username":         "john_doe",
				"email":            "  john@example.com  ",
				"password":         "MyP@ssw0rd123",
				"password_confirm": "MyP@ssw0rd123",
				"age":              25,
				"terms_accepted":   true,
			},
			false,
		},
		{
			"weak password",
			map[string]any{
				"username":         "john_doe",
				"email":            "john@example.com",
				"password":         "weak",
				"password_confirm": "weak",
				"age":              25,
				"terms_accepted":   true,
			},
			true,
		},
		{
			"terms not accepted",
			map[string]any{
				"username":         "john_doe",
				"email":            "john@example.com",
				"password":         "MyP@ssw0rd123",
				"password_confirm": "MyP@ssw0rd123",
				"age":              25,
				"terms_accepted":   false,
			},
			true,
		},
		{
			"underage user",
			map[string]any{
				"username":         "kid123",
				"email":            "kid@example.com",
				"password":         "MyP@ssw0rd123",
				"password_confirm": "MyP@ssw0rd123",
				"age":              12,
				"terms_accepted":   true,
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(tt.data)

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}

			// Check that email was trimmed
			if !result.HasErrors() {
				validData := result.ValidData()
				email := validData["email"].(string)
				if email != strings.TrimSpace(email) {
					t.Errorf("email was not trimmed properly")
				}
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Benchmark Tests
// -----------------------------------------------------------------------------

// BenchmarkSchema_Simple benchmarks simple schema validation
func BenchmarkSchema_Simple(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"name":  validation.String().Required(),
		"email": validation.String().Email(),
	})

	data := map[string]any{
		"name":  "John Doe",
		"email": "john@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}

// BenchmarkSchema_Complex benchmarks complex schema with cross-validation
func BenchmarkSchema_Complex(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"username":         validation.String().Required().Min(3).Max(20),
		"email":            validation.String().Required().Email(),
		"password":         validation.String().Required().Password(),
		"password_confirm": validation.String().Required(),
		"age":              validation.Number().Min(18).Integer(),
	}).CrossValidate(func(data map[string]any) error {
		password, _ := data["password"].(string)
		confirm, _ := data["password_confirm"].(string)
		if password != confirm {
			return fmt.Errorf("passwords must match")
		}
		return nil
	})

	data := map[string]any{
		"username":         "john_doe",
		"email":            "john@example.com",
		"password":         "MyP@ssw0rd123",
		"password_confirm": "MyP@ssw0rd123",
		"age":              25,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}

// BenchmarkSchema_ConditionalValidation benchmarks When() conditional validation
func BenchmarkSchema_ConditionalValidation(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"payment_method": validation.String().OneOf([]string{"credit_card", "paypal"}).Required(),
	}).When("payment_method", "credit_card", func() validation.Schema {
		return validation.Make().Shape(map[string]validation.Type{
			"card_number": validation.CreditCard().Required(),
		})
	})

	data := map[string]any{
		"payment_method": "credit_card",
		"card_number":    "4532015112830366",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}
