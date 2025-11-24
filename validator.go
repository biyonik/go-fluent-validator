// pkg/validation/schema.go
package validation

import (
	"fmt"

	"github.com/biyonik/go-fluent-validator/core"
)

//
// -----------------------------------------------------------------------------
// Validation Schema (Şema Tabanlı Doğrulama)
// -----------------------------------------------------------------------------
// Bu dosya, Fluent Validation sistemine ait `ValidationSchema` yapısını içerir.
// Amaç, bir JSON/Map yapısını Type bazlı doğrulamak, dönüştürmek, koşullu rule
// çalıştırmak ve çok alanlı (cross-field) validasyonlar uygulamaktır.
//
// Bu yapı, Laravel'in Validator::make() veya Yup, Joi, Zod gibi JS şema
// doğrulayıcılarının Go karşılığıdır.
//
// Metadata:
// @author    Ahmet ALTUN
// @github    github.com/biyonik
// @linkedin  linkedin.com/in/biyonik
// @email     ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

// Alias'lar (kullanıcı dostu API için)
type ValidationResult = core.ValidationResult
type Type = core.Type
type Schema = core.Schema

// conditionalRule
// -----------------------------------------------------------------------------
// "When" fonksiyonu ile kullanılan koşullu kuralı temsil eder.
// Bir alan belirli bir değere eşitse callback çağrılır ve alt-şema uygulanır.
type conditionalRule struct {
	field         string             // Koşul kontrol edilecek alan
	expectedValue any                // Beklenen değer
	callback      func() core.Schema // Çalıştırılacak alt şema
}

// ValidationSchema
// -----------------------------------------------------------------------------
// Bir validasyon şemasını temsil eder.
//
// Alanlar:
//   - shape: Her field için Type karşılığı
//   - crossValidators: Çok alanlı doğrulama fonksiyonları
//   - conditionalRules: When(...) ile eklenen koşullu doğrulama kuralları
//
// Örnek:
//
//	schema := validation.Make().Shape(map[string]core.Type{
//	    "email": validation.String().Email().Required(),
//	    "age":   validation.Number().Min(18),
//	})
//
// -----------------------------------------------------------------------------
type ValidationSchema struct {
	shape            map[string]core.Type
	crossValidators  []func(data map[string]any) error
	conditionalRules []conditionalRule
}

// Make
// -----------------------------------------------------------------------------
// Yeni bir ValidationSchema oluşturur.
//
// Dönüş:
//   - *ValidationSchema
//
// Örnek:
//
//	schema := validation.Make()
func Make() *ValidationSchema {
	return &ValidationSchema{
		shape:            make(map[string]core.Type),
		conditionalRules: make([]conditionalRule, 0),
	}
}

// Shape
// -----------------------------------------------------------------------------
// Şema için alan–type eşlemesini belirtir.
//
// Parametreler:
//   - shape: map[string]core.Type (örneğin email → StringType)
//
// Dönüş:
//   - core.Schema (chainable)
//
// Örnek:
//
//	schema.Shape(map[string]core.Type{
//	    "email": validation.String().Required().Email(),
//	    "age":   validation.Number().Min(18),
//	})
func (vs *ValidationSchema) Shape(shape map[string]core.Type) core.Schema {
	vs.shape = shape
	return vs
}

// CrossValidate
// -----------------------------------------------------------------------------
// Çok alanlı validasyon ekler. Örneğin password == password_confirm gibi.
//
// Parametreler:
//   - fn: func(data map[string]any) error
//
// Eğer hata dönerse _cross_validation alanına eklenir.
//
// Örnek:
//
//	schema.CrossValidate(func(data map[string]any) error {
//	    if data["password"] != data["password_confirm"] {
//	        return errors.New("Şifreler eşleşmiyor")
//	    }
//	    return nil
//	})
func (vs *ValidationSchema) CrossValidate(fn func(data map[string]any) error) core.Schema {
	vs.crossValidators = append(vs.crossValidators, fn)
	return vs
}

// When
// -----------------------------------------------------------------------------
// Koşullu doğrulama ekler. Belli bir alan belirlenen değere eşitse
// callback çağrılır ve alt-şema çalıştırılır.
//
// Parametreler:
//   - field: Koşul kontrol edilecek alan
//   - expectedValue: Bu değer eşleşirse callback tetiklenir
//   - callback: Alt şema döndüren fonksiyon
//
// Örnek:
//
//	schema.When("type", "corporate", func() core.Schema {
//	    return validation.Make().Shape(map[string]core.Type{
//	        "tax_number": validation.String().Required(),
//	    })
//	})
func (vs *ValidationSchema) When(field string, expectedValue any, callback func() core.Schema) core.Schema {
	vs.conditionalRules = append(vs.conditionalRules, conditionalRule{
		field:         field,
		expectedValue: expectedValue,
		callback:      callback,
	})
	return vs
}

// Validate
// -----------------------------------------------------------------------------
// Verilen veriyi şemaya göre doğrular.
//
// Adımlar:
//  1. Her alan için Transform çalıştırılır (tip dönüşümü).
//  2. Her alan için Validate çalıştırılır.
//  3. When(...) kuralları işlenir.
//  4. CrossValidate fonksiyonları çalıştırılır.
//  5. Hata yoksa ValidData set edilir.
//
// Parametre:
//   - data: map[string]any
//
// Dönüş:
//   - *core.ValidationResult
func (vs *ValidationSchema) Validate(data map[string]any) *core.ValidationResult {
	result := core.NewResult()
	transformedData := make(map[string]any)

	// 1) Transform aşaması
	for field, typ := range vs.shape {
		value := data[field]
		transformedValue, err := typ.Transform(value)
		if err != nil {
			result.AddError(field, fmt.Sprintf("Dönüşüm hatası: %s", err.Error()))
			continue
		}
		transformedData[field] = transformedValue
	}

	// 2) Field-level validation
	for field, typ := range vs.shape {
		typ.Validate(field, transformedData[field], result)
	}

	// 3) Koşullu şemalar
	if len(vs.conditionalRules) > 0 {
		for _, rule := range vs.conditionalRules {
			val, exists := transformedData[rule.field]
			if exists && val == rule.expectedValue {
				subSchema := rule.callback()
				subResult := subSchema.Validate(transformedData)
				if subResult.HasErrors() {
					for f, msgs := range subResult.Errors() {
						for _, msg := range msgs {
							result.AddError(f, msg)
						}
					}
				}
			}
		}
	}

	// 4) Cross-field validation
	if !result.HasErrors() {
		for _, fn := range vs.crossValidators {
			if err := fn(transformedData); err != nil {
				result.AddError("_cross_validation", err.Error())
			}
		}
	}

	// 5) Valid data set
	if !result.HasErrors() {
		result.SetValidData(transformedData)
	}

	return result
}
