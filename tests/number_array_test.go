// -----------------------------------------------------------------------------
// Number & Array Type Tests
// -----------------------------------------------------------------------------
// Bu dosya, NumberType ve ArrayType için kapsamlı test senaryolarını içerir.
//
// Test Kategorileri:
//   Number:
//     - Integer validation
//     - Min/max constraints
//     - Type conversion
//     - Default values
//   Array:
//     - Element validation
//     - Length constraints
//     - Nested arrays
//     - Element schema validation
//
// Metadata:
// @author   Ahmet ALTUN
// @github   github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email    ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

package tests

import (
	"testing"

	validation "github.com/biyonik/go-fluent-validator"
)

// -----------------------------------------------------------------------------
// Number Type Tests
// -----------------------------------------------------------------------------

// TestNumberType_Required tests required number field
func TestNumberType_Required(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		wantError bool
	}{
		{"valid integer", 42, false},
		{"valid float", 42.5, false},
		{"zero value", 0, false},
		{"nil value", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"age": validation.Number().Required(),
			})

			result := schema.Validate(map[string]any{
				"age": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
			}
		})
	}
}

// TestNumberType_MinMax tests min/max constraints
func TestNumberType_MinMax(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		min       float64
		max       float64
		wantError bool
	}{
		{"within range", 25, 18, 100, false},
		{"exact min", 18, 18, 100, false},
		{"exact max", 100, 18, 100, false},
		{"below min", 17, 18, 100, true},
		{"above max", 101, 18, 100, true},
		{"float within range", 25.5, 18, 100, false},
		{"float below min", 17.9, 18, 100, true},
		{"negative within range", -5, -10, 0, false},
		{"negative below min", -11, -10, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"value": validation.Number().Min(tt.min).Max(tt.max),
			})

			result := schema.Validate(map[string]any{
				"value": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("value %v: got error = %v, want error = %v",
					tt.value, result.HasErrors(), tt.wantError)
			}
		})
	}
}

// TestNumberType_Integer tests integer constraint
func TestNumberType_Integer(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		wantError bool
	}{
		{"integer", 42, false},
		{"zero", 0, false},
		{"negative integer", -42, false},
		{"float with .0", 42.0, false}, // Should be considered integer
		{"float with decimals", 42.5, true},
		{"small decimal", 42.01, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"count": validation.Number().Integer(),
			})

			result := schema.Validate(map[string]any{
				"count": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("value %v: got error = %v, want error = %v",
					tt.value, result.HasErrors(), tt.wantError)
			}
		})
	}
}

// TestNumberType_TypeConversion tests different number type handling
func TestNumberType_TypeConversion(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		wantError bool
	}{
		{"int", int(42), false},
		{"int8", int8(42), false},
		{"int16", int16(42), false},
		{"int32", int32(42), false},
		{"int64", int64(42), false},
		{"float32", float32(42.5), false},
		{"float64", float64(42.5), false},
		{"string", "42", true}, // Should not auto-convert strings
		{"bool", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"number": validation.Number(),
			})

			result := schema.Validate(map[string]any{
				"number": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("type %T: got error = %v, want error = %v",
					tt.value, result.HasErrors(), tt.wantError)
			}
		})
	}
}

// TestNumberType_Default tests default value
func TestNumberType_Default(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"count": validation.Number().Default(0),
	})

	tests := []struct {
		name     string
		value    any
		expected float64
	}{
		{"provided value", 42, 42.0},
		{"nil uses default", nil, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(map[string]any{
				"count": tt.value,
			})

			if result.HasErrors() {
				t.Errorf("unexpected error: %v", result.Errors())
			}

			validData := result.ValidData()
			var got float64
			switch v := validData["count"].(type) {
			case int:
				got = float64(v)
			case float64:
				got = v
			}

			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Array Type Tests
// -----------------------------------------------------------------------------

// TestArrayType_Required tests required array field
func TestArrayType_Required(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		wantError bool
	}{
		{"valid array", []any{"a", "b", "c"}, false},
		{"empty array", []any{}, false},
		{"nil value", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"tags": validation.Array().Required(),
			})

			result := schema.Validate(map[string]any{
				"tags": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
			}
		})
	}
}

// TestArrayType_MinMax tests array length constraints
func TestArrayType_MinMax(t *testing.T) {
	tests := []struct {
		name      string
		value     []any
		min       int
		max       int
		wantError bool
	}{
		{"within range", []any{1, 2, 3}, 2, 5, false},
		{"exact min", []any{1, 2}, 2, 5, false},
		{"exact max", []any{1, 2, 3, 4, 5}, 2, 5, false},
		{"below min", []any{1}, 2, 5, true},
		{"above max", []any{1, 2, 3, 4, 5, 6}, 2, 5, true},
		{"empty below min", []any{}, 1, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"items": validation.Array().Min(tt.min).Max(tt.max),
			})

			result := schema.Validate(map[string]any{
				"items": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("array length %d: got error = %v, want error = %v",
					len(tt.value), result.HasErrors(), tt.wantError)
			}
		})
	}
}

// TestArrayType_Elements tests element validation
func TestArrayType_Elements(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"emails": validation.Array().Elements(
			validation.String().Email(),
		),
	})

	tests := []struct {
		name      string
		value     []any
		wantError bool
	}{
		{
			"all valid emails",
			[]any{"test@example.com", "user@domain.com"},
			false,
		},
		{
			"one invalid email",
			[]any{"test@example.com", "invalid-email"},
			true,
		},
		{
			"all invalid emails",
			[]any{"invalid1", "invalid2"},
			true,
		},
		{
			"empty array",
			[]any{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(map[string]any{
				"emails": tt.value,
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

// TestArrayType_NestedObjects tests array of objects validation
func TestArrayType_NestedObjects(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"users": validation.Array().Elements(
			validation.Object().Shape(map[string]validation.Type{
				"name":  validation.String().Required().Min(2),
				"age":   validation.Number().Min(0).Max(150),
				"email": validation.String().Email(),
			}),
		),
	})

	tests := []struct {
		name      string
		value     []any
		wantError bool
	}{
		{
			"all valid users",
			[]any{
				map[string]any{"name": "John", "age": 30, "email": "john@example.com"},
				map[string]any{"name": "Jane", "age": 25, "email": "jane@example.com"},
			},
			false,
		},
		{
			"one invalid user - missing name",
			[]any{
				map[string]any{"name": "John", "age": 30, "email": "john@example.com"},
				map[string]any{"age": 25, "email": "jane@example.com"},
			},
			true,
		},
		{
			"one invalid user - invalid email",
			[]any{
				map[string]any{"name": "John", "age": 30, "email": "invalid-email"},
			},
			true,
		},
		{
			"one invalid user - age out of range",
			[]any{
				map[string]any{"name": "John", "age": 200, "email": "john@example.com"},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(map[string]any{
				"users": tt.value,
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

// TestArrayType_TypeMismatch tests invalid type handling
func TestArrayType_TypeMismatch(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"items": validation.Array().Required(),
	})

	tests := []struct {
		name  string
		value any
	}{
		{"string", "not an array"},
		{"integer", 123},
		{"boolean", true},
		{"map", map[string]string{"key": "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.Validate(map[string]any{
				"items": tt.value,
			})

			if !result.HasErrors() {
				t.Errorf("expected error for type %T, got none", tt.value)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Benchmark Tests
// -----------------------------------------------------------------------------

// BenchmarkNumberValidation benchmarks number validation
func BenchmarkNumberValidation(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"age":   validation.Number().Min(0).Max(150).Integer().Required(),
		"price": validation.Number().Min(0),
	})

	data := map[string]any{
		"age":   25,
		"price": 99.99,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}

// BenchmarkArrayValidation_Simple benchmarks simple array validation
func BenchmarkArrayValidation_Simple(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"tags": validation.Array().Min(1).Max(10),
	})

	data := map[string]any{
		"tags": []any{"go", "validation", "testing"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}

// BenchmarkArrayValidation_WithElements benchmarks array with element validation
func BenchmarkArrayValidation_WithElements(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"emails": validation.Array().Elements(
			validation.String().Email(),
		),
	})

	data := map[string]any{
		"emails": []any{
			"user1@example.com",
			"user2@example.com",
			"user3@example.com",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}

// BenchmarkArrayValidation_NestedObjects benchmarks array of objects
func BenchmarkArrayValidation_NestedObjects(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"users": validation.Array().Elements(
			validation.Object().Shape(map[string]validation.Type{
				"name":  validation.String().Required(),
				"email": validation.String().Email(),
			}),
		),
	})

	data := map[string]any{
		"users": []any{
			map[string]any{"name": "John", "email": "john@example.com"},
			map[string]any{"name": "Jane", "email": "jane@example.com"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}
