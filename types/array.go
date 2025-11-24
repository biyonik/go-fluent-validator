// -----------------------------------------------------------------------------
// ArrayType
// -----------------------------------------------------------------------------
// Bu dosya, Go ile geliştirilmiş Fluent Validation yapısının "Array" tipini
// temsil etmektedir. Laravel validation, Symfony Validator veya TypeScript
// Zod tarzındaki güçlü ve akıcı doğrulama yaklaşımlarının Go diline birebir
// aktarılması hedeflenmiştir.
//
// ArrayType, özellikle API gelişiminde sıkça ihtiyaç duyulan:
// - Dizi zorunluluğu (required)
// - Minimum / maksimum eleman sayısı kontrolü
// - Dizinin her bir elemanının belirli bir şemaya göre doğrulanması
// - Her eleman için bağımsız dönüştürme (transform) işlemi
//
// gibi kritik özellikleri destekler.
//
// Bu yapı sayesinde hem dışarıdan gelen JSON verilerini normalize edip hem de
// her bir elemanın kendi tip kurallarına uygunluğunu güvenli ve sistematik
// biçimde kontrol edebilirsiniz.
//
// Ayrıca hatalar, alan yolunu belirterek (`field[0]`, `field[1]`) ilettikleri için,
// özellikle nested dizilerde debug etmek çok kolaylaşır.
//
// Yazar Bilgileri:
// @author   Ahmet ALTUN
// @github   github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email    ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

package types

import (
	"fmt"

	"github.com/biyonik/go-fluent-validator/core"
	"github.com/biyonik/go-fluent-validator/i18n"
)

// ArrayType, bir diziyi (array/slice) doğrulamak için kullanılan gelişmiş bir
// doğrulama tipidir. Uzunluk kontrolü, zorunluluk, element tipi doğrulama gibi
// özellikleri içerir.
type ArrayType struct {
	core.BaseType
	minLength        *int      // Minimum eleman sayısı
	maxLength        *int      // Maksimum eleman sayısı
	elementSchema    core.Type // Her bir elemanın uyacağı şema
	customValidation *core.CustomValidation
	// New validators
	isUnique         bool
	containsValue    *any
	isNotEmpty       bool
}

// Required, alanın zorunlu olduğunu belirtir.
func (a *ArrayType) Required() *ArrayType {
	a.SetRequired()
	return a
}

// Label, kullanıcıya gösterilecek alan adını özelleştirir.
func (a *ArrayType) Label(label string) *ArrayType {
	a.SetLabel(label)
	return a
}

// Min, dizide bulunması gereken minimum eleman sayısını tanımlar.
func (a *ArrayType) Min(length int) *ArrayType {
	a.minLength = &length
	return a
}

// Max, dizide bulunmasına izin verilen maksimum eleman sayısını tanımlar.
func (a *ArrayType) Max(length int) *ArrayType {
	a.maxLength = &length
	return a
}

// Elements, dizinin her elemanının belirli bir doğrulama tipine uymasını sağlar.
// Örneğin: validation.Array().Elements(validation.String().Min(3))
func (a *ArrayType) Elements(schema core.Type) *ArrayType {
	a.elementSchema = schema
	return a
}

// Transform, dizinin kendisini ve varsa elemanlarını dönüştürür.
// Bu aşama, veri normalize etme (ör. trim, type-cast) için kritiktir.
func (a *ArrayType) Transform(value any) (any, error) {
	value, err := a.BaseType.Transform(value)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	slice, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("dizi (array) tipinde olmalıdır")
	}

	if a.elementSchema != nil {
		transformedSlice := make([]any, len(slice))
		for i, item := range slice {
			transformedItem, err := a.elementSchema.Transform(item)
			if err != nil {
				return nil, fmt.Errorf("dizi index %d: %w", i, err)
			}
			transformedSlice[i] = transformedItem
		}
		return transformedSlice, nil
	}
	return slice, nil
}

func (a *ArrayType) Custom(validator func([]any) error) *ArrayType {
	if a.customValidation == nil {
		a.customValidation = core.NewCustomValidation()
	}

	a.customValidation.AddSync(func(value any) error {
		if value == nil {
			return nil
		}

		arr, ok := value.([]any)
		if !ok {
			return fmt.Errorf("value must be array")
		}

		return validator(arr)
	})

	return a
}

func (a *ArrayType) AddRule(rule core.Rule) *ArrayType {
	if a.customValidation == nil {
		a.customValidation = core.NewCustomValidation()
	}
	a.customValidation.AddRule(rule)
	return a
}

// Unique ensures all elements in the array are unique
func (a *ArrayType) Unique() *ArrayType {
	a.isUnique = true
	return a
}

// Contains ensures the array contains a specific value
func (a *ArrayType) Contains(value any) *ArrayType {
	a.containsValue = &value
	return a
}

// NotEmpty ensures the array is not empty
func (a *ArrayType) NotEmpty() *ArrayType {
	a.isNotEmpty = true
	return a
}

// Validate, dizinin uzunluk doğrulamasını ve eleman doğrulamasını yapar.
// Hatalar, `field[0]`, `field[1]` formatında detaylı bir şekilde işlenir.
func (a *ArrayType) Validate(field string, value any, result *core.ValidationResult) {
	a.BaseType.Validate(field, value, result)
	if result.HasErrors() {
		return
	}
	if value == nil {
		return
	}

	slice, ok := value.([]any)
	if !ok {
		result.AddError(field, fmt.Sprintf("%s alanı dizi (array) tipinde olmalıdır", a.GetLabel(field)))
		return
	}

	fieldName := a.GetLabel(field)

	if a.minLength != nil && len(slice) < *a.minLength {
		result.AddError(field, fmt.Sprintf("%s alanında en az %d eleman olmalıdır", fieldName, *a.minLength))
	}
	if a.maxLength != nil && len(slice) > *a.maxLength {
		result.AddError(field, fmt.Sprintf("%s alanında en fazla %d eleman olmalıdır", fieldName, *a.maxLength))
	}

	// New validators
	if a.isNotEmpty && len(slice) == 0 {
		result.AddError(field, i18n.Get(i18n.KeyNotEmpty, fieldName))
	}

	if a.isUnique {
		seen := make(map[string]bool)
		for i, item := range slice {
			key := fmt.Sprintf("%v", item)
			if seen[key] {
				result.AddError(field, i18n.Get(i18n.KeyUnique, fieldName))
				break
			}
			seen[key] = true
			_ = i // prevent unused variable warning
		}
	}

	if a.containsValue != nil {
		found := false
		for _, item := range slice {
			if fmt.Sprintf("%v", item) == fmt.Sprintf("%v", *a.containsValue) {
				found = true
				break
			}
		}
		if !found {
			result.AddError(field, i18n.Get(i18n.KeyArrayContains, fieldName, *a.containsValue))
		}
	}

	if a.customValidation != nil && a.customValidation.HasValidators() {
		a.customValidation.ValidateSync(field, value, result)
	}

	if a.elementSchema != nil {
		for i, item := range slice {
			elementFieldPath := fmt.Sprintf("%s[%d]", field, i)
			a.elementSchema.Validate(elementFieldPath, item, result)
		}
	}
}
