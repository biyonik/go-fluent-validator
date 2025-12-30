package types

import (
	"fmt"
	"math"

	"github.com/biyonik/go-fluent-validator/core"
	"github.com/biyonik/go-fluent-validator/i18n"
)

// NumberType
//
// Sayısal değerleri doğrulamak için kullanılan tiptir.
// Float veya integer değerleri kabul eder ve opsiyonel olarak
// minimum, maksimum veya tamsayı (integer) kısıtlamaları eklenebilir.
//
// Özellikler:
//   - float64 veya int tipinde değer kabul eder
//   - Min() ve Max() ile değer aralığı kısıtlanabilir
//   - Integer() ile yalnızca tamsayı değerler kabul edilir
//   - Required() ile boş geçilemez kılınabilir
//
// Kullanım alanları:
//   - Form ve API doğrulamaları
//   - Finansal veya sayısal verilerin kontrolü
//
// Yazar Bilgileri:
//   - @author  Ahmet Altun
//   - @github  https://github.com/biyonik
//   - @linkedin https://linkedin.com/in/biyonik
//   - @email   ahmet.altun60@gmail.com
type NumberType struct {
	core.BaseType
	min              *float64
	max              *float64
	isInteger        bool
	customValidation *core.CustomValidation
	// New validators
	isPositive  bool
	isNegative  bool
	multipleOf  *float64
	betweenMin  *float64
	betweenMax  *float64
}

// Required, alanın boş geçilemeyeceğini belirtir.
//
// Döndürür:
//   - *NumberType
func (n *NumberType) Required() *NumberType {
	n.SetRequired()
	return n
}

// Label, doğrulama hatalarında gösterilecek kullanıcı dostu alan adını belirler.
//
// Parametreler:
//   - label (string): kullanıcıya gösterilecek isim
//
// Döndürür:
//   - *NumberType
func (n *NumberType) Label(label string) *NumberType {
	n.SetLabel(label)
	return n
}

// Default, alanın varsayılan değerini belirler.
//
// Parametreler:
//   - value (any): int, float32 veya float64
//
// Döndürür:
//   - *NumberType
func (n *NumberType) Default(value any) *NumberType {
	switch v := value.(type) {
	case int:
		n.SetDefault(float64(v))
	case float64:
		n.SetDefault(v)
	case float32:
		n.SetDefault(float64(v))
	default:
		n.SetDefault(value)
	}
	return n
}

// Min, alanın alabileceği minimum değeri belirler.
//
// Parametreler:
//   - val (float64): minimum değer
//
// Döndürür:
//   - *NumberType
func (n *NumberType) Min(val float64) *NumberType {
	n.min = &val
	return n
}

// Max, alanın alabileceği maksimum değeri belirler.
//
// Parametreler:
//   - val (float64): maksimum değer
//
// Döndürür:
//   - *NumberType
func (n *NumberType) Max(val float64) *NumberType {
	n.max = &val
	return n
}

// Integer, alanın yalnızca tamsayı değer almasını zorunlu kılar.
//
// Döndürür:
//   - *NumberType
func (n *NumberType) Integer() *NumberType {
	n.isInteger = true
	return n
}

func (n *NumberType) Custom(validator func(float64) error) *NumberType {
	if n.customValidation == nil {
		n.customValidation = core.NewCustomValidation()
	}

	n.customValidation.AddSync(func(value any) error {
		if value == nil {
			return nil
		}

		var num float64
		switch v := value.(type) {
		case int:
			num = float64(v)
		case int8:
			num = float64(v)
		case int16:
			num = float64(v)
		case int32:
			num = float64(v)
		case int64:
			num = float64(v)
		case float32:
			num = float64(v)
		case float64:
			num = v
		default:
			return fmt.Errorf("value must be number")
		}

		return validator(num)
	})

	return n
}

func (n *NumberType) AddRule(rule core.Rule) *NumberType {
	if n.customValidation == nil {
		n.customValidation = core.NewCustomValidation()
	}
	n.customValidation.AddRule(rule)
	return n
}

// Positive ensures the number is greater than zero
func (n *NumberType) Positive() *NumberType {
	n.isPositive = true
	return n
}

// Negative ensures the number is less than zero
func (n *NumberType) Negative() *NumberType {
	n.isNegative = true
	return n
}

// MultipleOf ensures the number is a multiple of the given value
func (n *NumberType) MultipleOf(value float64) *NumberType {
	n.multipleOf = &value
	return n
}

// Between ensures the number is between min and max (inclusive)
func (n *NumberType) Between(min, max float64) *NumberType {
	n.betweenMin = &min
	n.betweenMax = &max
	return n
}

// Validate, alanın sayısal geçerliliğini kontrol eder.
//
// İşlem sırası:
//  1. BaseType doğrulamaları (required, nullable)
//  2. Tip kontrolü (int, float32, float64)
//  3. Integer kısıtlaması (opsiyonel)
//  4. Min ve Max kontrolleri (opsiyonel)
//
// Parametreler:
//   - field (string): alan adı
//   - value (any): doğrulanacak değer
//   - result (*core.ValidationResult): doğrulama sonucu
func (n *NumberType) Validate(field string, value any, result *core.ValidationResult) {
	n.BaseType.Validate(field, value, result)
	if result.HasErrors() {
		return
	}
	if value == nil {
		return
	}

	var num float64
	var ok bool

	switch v := value.(type) {
	case int:
		num = float64(v)
		ok = true
	case int8:
		num = float64(v)
		ok = true
	case int16:
		num = float64(v)
		ok = true
	case int32:
		num = float64(v)
		ok = true
	case int64:
		num = float64(v)
		ok = true
	case float64:
		num = v
		ok = true
	case float32:
		num = float64(v)
		ok = true
	default:
		ok = false
	}

	fieldName := n.GetLabel(field)

	if !ok {
		result.AddError(field, i18n.Get(i18n.KeyNumeric, fieldName))
		return
	}

	if n.isInteger && num != float64(int64(num)) {
		result.AddError(field, i18n.Get(i18n.KeyInteger, fieldName))
	}
	if n.min != nil && num < *n.min {
		result.AddError(field, i18n.Get(i18n.KeyMin, fieldName, *n.min))
	}
	if n.max != nil && num > *n.max {
		result.AddError(field, i18n.Get(i18n.KeyMax, fieldName, *n.max))
	}

	// New validators
	if n.isPositive && num <= 0 {
		result.AddError(field, i18n.Get(i18n.KeyPositive, fieldName))
	}

	if n.isNegative && num >= 0 {
		result.AddError(field, i18n.Get(i18n.KeyNegative, fieldName))
	}

	if n.multipleOf != nil {
		// Check if num is a multiple of multipleOf using modulo with floating point precision
		remainder := math.Mod(num, *n.multipleOf)
		if math.Abs(remainder) > 1e-9 { // Use small epsilon for floating point comparison
			result.AddError(field, i18n.Get(i18n.KeyMultipleOf, fieldName, *n.multipleOf))
		}
	}

	if n.betweenMin != nil && n.betweenMax != nil {
		if num < *n.betweenMin || num > *n.betweenMax {
			result.AddError(field, i18n.Get(i18n.KeyBetween, fieldName, *n.betweenMin, *n.betweenMax))
		}
	}

	if n.customValidation != nil && n.customValidation.HasValidators() {
		n.customValidation.ValidateSync(field, value, result)
	}
}
