// -----------------------------------------------------------------------------
// ValidationResult Tests
// -----------------------------------------------------------------------------
// Tests for AllErrors() and error handling
// -----------------------------------------------------------------------------

package tests

import (
	"testing"

	validation "github.com/biyonik/go-fluent-validator"
)

// TestResult_AllErrors tests the AllErrors() method
func TestResult_AllErrors(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"name":  validation.String().Required().Min(3).Label("Name"),
		"email": validation.String().Required().Email().Label("Email"),
		"age":   validation.Number().Required().Min(18).Label("Age"),
	})

	t.Run("no errors", func(t *testing.T) {
		result := schema.Validate(map[string]any{
			"name":  "John Doe",
			"email": "john@example.com",
			"age":   25,
		})

		allErrors := result.AllErrors()
		if len(allErrors) != 0 {
			t.Errorf("expected no errors, got %d: %v", len(allErrors), allErrors)
		}
	})

	t.Run("single field error", func(t *testing.T) {
		result := schema.Validate(map[string]any{
			"name":  "Jo",
			"email": "john@example.com",
			"age":   25,
		})

		allErrors := result.AllErrors()
		if len(allErrors) != 1 {
			t.Errorf("expected 1 error, got %d: %v", len(allErrors), allErrors)
		}
	})

	t.Run("multiple field errors", func(t *testing.T) {
		result := schema.Validate(map[string]any{
			"name":  "Jo",
			"email": "invalid-email",
			"age":   15,
		})

		allErrors := result.AllErrors()
		if len(allErrors) != 3 {
			t.Errorf("expected 3 errors, got %d: %v", len(allErrors), allErrors)
		}

		// Verify it's a flat array
		for _, err := range allErrors {
			if err == "" {
				t.Error("error message should not be empty")
			}
		}
	})

	t.Run("multiple errors same field", func(t *testing.T) {
		complexSchema := validation.Make().Shape(map[string]validation.Type{
			"username": validation.String().
				Required().
				Min(3).
				Max(10).
				Alpha().
				Label("Username"),
		})

		result := complexSchema.Validate(map[string]any{
			"username": "ab123", // too short + contains numbers
		})

		allErrors := result.AllErrors()
		if len(allErrors) < 2 {
			t.Errorf("expected at least 2 errors, got %d: %v", len(allErrors), allErrors)
		}
	})
}

// TestResult_ErrorsVsAllErrors compares Errors() and AllErrors()
func TestResult_ErrorsVsAllErrors(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"field1": validation.String().Required().Label("Field 1"),
		"field2": validation.String().Required().Label("Field 2"),
		"field3": validation.String().Required().Label("Field 3"),
	})

	result := schema.Validate(map[string]any{})

	// Errors() returns map[string][]string
	errorsMap := result.Errors()
	if len(errorsMap) != 3 {
		t.Errorf("expected 3 fields with errors, got %d", len(errorsMap))
	}

	// AllErrors() returns []string
	allErrors := result.AllErrors()
	if len(allErrors) != 3 {
		t.Errorf("expected 3 total errors, got %d", len(allErrors))
	}

	// Count should match
	totalFromMap := 0
	for _, errs := range errorsMap {
		totalFromMap += len(errs)
	}

	if totalFromMap != len(allErrors) {
		t.Errorf("error count mismatch: map has %d, AllErrors has %d", totalFromMap, len(allErrors))
	}
}

// TestResult_RealWorld_APIResponse tests API response formatting
func TestResult_RealWorld_APIResponse(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"email":    validation.String().Required().Email().Label("Email"),
		"password": validation.String().Required().Min(8).Label("Password"),
	})

	t.Run("API error response format", func(t *testing.T) {
		result := schema.Validate(map[string]any{
			"email":    "invalid",
			"password": "short",
		})

		if !result.HasErrors() {
			t.Fatal("expected validation errors")
		}

		// Simulated API response with AllErrors()
		type APIErrorResponse struct {
			Success  bool                `json:"success"`
			Errors   map[string][]string `json:"errors"`   // Detailed
			Messages []string            `json:"messages"` // Simple list
		}

		apiResp := APIErrorResponse{
			Success:  false,
			Errors:   result.Errors(),
			Messages: result.AllErrors(),
		}

		if apiResp.Success {
			t.Error("success should be false")
		}

		if len(apiResp.Errors) == 0 {
			t.Error("errors map should not be empty")
		}

		if len(apiResp.Messages) == 0 {
			t.Error("messages array should not be empty")
		}

		t.Logf("API Response: %+v", apiResp)
	})
}
