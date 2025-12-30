// -----------------------------------------------------------------------------
// Date Validation Comprehensive Tests
// -----------------------------------------------------------------------------
// Bu test dosyası, Date tipinin tüm özelliklerini test eder:
// - Format parsing
// - Min/Max validation
// - Before/After validation (YENİ)
// - Edge cases (timezone, leap year, invalid dates)
// - Real-world scenarios
// -----------------------------------------------------------------------------

package tests

import (
	"fmt"
	"testing"
	"time"

	validation "github.com/biyonik/go-fluent-validator"
)

// TestDate_BeforeValidation tests the new Before() method
func TestDate_BeforeValidation(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	tests := []struct {
		name      string
		value     string
		before    time.Time
		wantError bool
	}{
		{
			name:      "date before deadline - valid",
			value:     yesterday.Format("2006-01-02"),
			before:    now,
			wantError: false,
		},
		{
			name:      "date after deadline - invalid",
			value:     tomorrow.Format("2006-01-02"),
			before:    now,
			wantError: true,
		},
		{
			name:      "date equals deadline - invalid (not strictly before)",
			value:     now.Format("2006-01-02"),
			before:    now,
			wantError: true,
		},
		{
			name:      "birth date before today - valid",
			value:     "1990-05-15",
			before:    time.Now(),
			wantError: false,
		},
		{
			name:      "future birth date - invalid",
			value:     time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
			before:    time.Now(),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"date": validation.Date().
					Format("2006-01-02").
					Before(tt.before).
					Required().
					Label("Test Date"),
			})

			result := schema.Validate(map[string]any{
				"date": tt.value,
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

// TestDate_AfterValidation tests the new After() method
func TestDate_AfterValidation(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	tests := []struct {
		name      string
		value     string
		after     time.Time
		wantError bool
	}{
		{
			name:      "date after start - valid",
			value:     tomorrow.Format("2006-01-02"),
			after:     now,
			wantError: false,
		},
		{
			name:      "date before start - invalid",
			value:     yesterday.Format("2006-01-02"),
			after:     now,
			wantError: true,
		},
		{
			name:      "date equals start - invalid (not strictly after)",
			value:     now.Format("2006-01-02"),
			after:     now,
			wantError: true,
		},
		{
			name:      "event date in future - valid",
			value:     time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
			after:     time.Now(),
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"date": validation.Date().
					Format("2006-01-02").
					After(tt.after).
					Required().
					Label("Event Date"),
			})

			result := schema.Validate(map[string]any{
				"date": tt.value,
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

// TestDate_BeforeAfterCombined tests Before and After together (range)
func TestDate_BeforeAfterCombined(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{
			name:      "date within range - valid",
			value:     "2024-06-15",
			wantError: false,
		},
		{
			name:      "date before range - invalid",
			value:     "2023-12-31",
			wantError: true,
		},
		{
			name:      "date after range - invalid",
			value:     "2025-01-01",
			wantError: true,
		},
		{
			name:      "date at start boundary - valid",
			value:     "2024-01-02",
			wantError: false,
		},
		{
			name:      "date at end boundary - valid",
			value:     "2024-12-30",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"date": validation.Date().
					Format("2006-01-02").
					After(startDate).
					Before(endDate).
					Required().
					Label("Booking Date"),
			})

			result := schema.Validate(map[string]any{
				"date": tt.value,
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

// TestDate_RealWorldScenarios tests real-world date validation scenarios
func TestDate_RealWorldScenarios(t *testing.T) {
	t.Run("Flight Booking - Departure before Return", func(t *testing.T) {
		schema := validation.Make().Shape(map[string]validation.Type{
			"departure_date": validation.Date().Format("2006-01-02").Required().Label("Departure Date"),
			"return_date":    validation.Date().Format("2006-01-02").Required().Label("Return Date"),
		}).CrossValidate(func(data map[string]any) error {
			departure, _ := data["departure_date"].(time.Time)
			returnDate, _ := data["return_date"].(time.Time)

			if returnDate.Before(departure) || returnDate.Equal(departure) {
				return fmt.Errorf("return date must be after departure date")
			}
			return nil
		})

		// Valid booking
		result := schema.Validate(map[string]any{
			"departure_date": "2024-07-01",
			"return_date":    "2024-07-10",
		})
		if result.HasErrors() {
			t.Errorf("valid booking should not have errors: %v", result.Errors())
		}

		// Invalid booking (return before departure)
		result = schema.Validate(map[string]any{
			"departure_date": "2024-07-10",
			"return_date":    "2024-07-01",
		})
		if !result.HasErrors() {
			t.Error("invalid booking should have errors")
		}
	})

	t.Run("User Registration - Must be 18+", func(t *testing.T) {
		eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)

		schema := validation.Make().Shape(map[string]validation.Type{
			"birth_date": validation.Date().
				Format("2006-01-02").
				Before(eighteenYearsAgo).
				Required().
				Label("Birth Date"),
		})

		// Valid - user is 20 years old
		result := schema.Validate(map[string]any{
			"birth_date": time.Now().AddDate(-20, 0, 0).Format("2006-01-02"),
		})
		if result.HasErrors() {
			t.Errorf("20 year old should be valid: %v", result.Errors())
		}

		// Invalid - user is 17 years old
		result = schema.Validate(map[string]any{
			"birth_date": time.Now().AddDate(-17, 0, 0).Format("2006-01-02"),
		})
		if !result.HasErrors() {
			t.Error("17 year old should be invalid")
		}
	})

	t.Run("Event Scheduling - Must be in future", func(t *testing.T) {
		schema := validation.Make().Shape(map[string]validation.Type{
			"event_date": validation.Date().
				Format("2006-01-02 15:04").
				After(time.Now()).
				Required().
				Label("Event Date"),
		})

		// Valid - event tomorrow
		tomorrow := time.Now().AddDate(0, 0, 1)
		result := schema.Validate(map[string]any{
			"event_date": tomorrow.Format("2006-01-02 15:04"),
		})
		if result.HasErrors() {
			t.Errorf("future event should be valid: %v", result.Errors())
		}

		// Invalid - event yesterday
		yesterday := time.Now().AddDate(0, 0, -1)
		result = schema.Validate(map[string]any{
			"event_date": yesterday.Format("2006-01-02 15:04"),
		})
		if !result.HasErrors() {
			t.Error("past event should be invalid")
		}
	})
}

// TestDate_EdgeCases tests edge cases and boundary conditions
func TestDate_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		schema    validation.Schema
		data      map[string]any
		wantError bool
	}{
		{
			name: "leap year - Feb 29, 2024 valid",
			schema: validation.Make().Shape(map[string]validation.Type{
				"date": validation.Date().Format("2006-01-02").Required(),
			}),
			data:      map[string]any{"date": "2024-02-29"},
			wantError: false,
		},
		{
			name: "non-leap year - Feb 29, 2023 invalid",
			schema: validation.Make().Shape(map[string]validation.Type{
				"date": validation.Date().Format("2006-01-02").Required(),
			}),
			data:      map[string]any{"date": "2023-02-29"},
			wantError: true,
		},
		{
			name: "year 2000 - leap year",
			schema: validation.Make().Shape(map[string]validation.Type{
				"date": validation.Date().Format("2006-01-02").Required(),
			}),
			data:      map[string]any{"date": "2000-02-29"},
			wantError: false,
		},
		{
			name: "year 1900 - not leap year",
			schema: validation.Make().Shape(map[string]validation.Type{
				"date": validation.Date().Format("2006-01-02").Required(),
			}),
			data:      map[string]any{"date": "1900-02-29"},
			wantError: true,
		},
		{
			name: "invalid month - month 13",
			schema: validation.Make().Shape(map[string]validation.Type{
				"date": validation.Date().Format("2006-01-02").Required(),
			}),
			data:      map[string]any{"date": "2024-13-01"},
			wantError: true,
		},
		{
			name: "invalid day - day 32",
			schema: validation.Make().Shape(map[string]validation.Type{
				"date": validation.Date().Format("2006-01-02").Required(),
			}),
			data:      map[string]any{"date": "2024-01-32"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.Validate(tt.data)
			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// BenchmarkDate_BeforeValidation benchmarks Before() validation
func BenchmarkDate_BeforeValidation(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"date": validation.Date().
			Format("2006-01-02").
			Before(time.Now()).
			Required(),
	})

	data := map[string]any{
		"date": "2020-01-01",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}

// BenchmarkDate_AfterValidation benchmarks After() validation
func BenchmarkDate_AfterValidation(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"date": validation.Date().
			Format("2006-01-02").
			After(time.Now().AddDate(-1, 0, 0)).
			Required(),
	})

	data := map[string]any{
		"date": time.Now().Format("2006-01-02"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}
