package types

import (
	"fmt"

	"github.com/biyonik/go-fluent-validator/core"
	"github.com/biyonik/go-fluent-validator/i18n"
)

// ObjectType
//
// Nesne (object) tipindeki alanları doğrulamak için kullanılan tiptir.
// İçerisinde belirli alanlar ve bu alanlara uygulanacak alt şemalar (Type) tanımlanabilir.
//
// Özellikler:
//   - Shape() ile alt alanları ve tiplerini belirleyebilirsiniz
//   - Required() ile boş geçilemez kılınabilir
//   - Transform ve Validate metodları ile alt alanlar derinlemesine doğrulanır
//
// Kullanım alanları:
//   - JSON veya nested veri doğrulamaları
//   - API ve form payload kontrolleri
//
// Yazar Bilgileri:
//   - @author  Ahmet Altun
//   - @github  https://github.com/biyonik
//   - @linkedin https://linkedin.com/in/biyonik
//   - @email   ahmet.altun60@gmail.com
type ObjectType struct {
	core.BaseType
	shape            map[string]core.Type
	customValidation *core.CustomValidation
}

// Required, alanın boş geçilemeyeceğini belirtir.
//
// Döndürür:
//   - *ObjectType
func (o *ObjectType) Required() *ObjectType {
	o.SetRequired()
	return o
}

// Label, doğrulama hatalarında gösterilecek kullanıcı dostu alan adını belirler.
//
// Parametreler:
//   - label (string): kullanıcıya gösterilecek isim
//
// Döndürür:
//   - *ObjectType
func (o *ObjectType) Label(label string) *ObjectType {
	o.SetLabel(label)
	return o
}

// Shape, nesnenin alt alanlarını ve bu alanların tiplerini tanımlar.
//
// Parametreler:
//   - shape (map[string]core.Type): alt alanlar ve tipleri
//
// Döndürür:
//   - *ObjectType
func (o *ObjectType) Shape(shape map[string]core.Type) *ObjectType {
	o.shape = shape
	return o
}

// Custom adds a custom validation function
func (o *ObjectType) Custom(validator func(map[string]any) error) *ObjectType {
	if o.customValidation == nil {
		o.customValidation = core.NewCustomValidation()
	}

	o.customValidation.AddSync(func(value any) error {
		if value == nil {
			return nil
		}

		objVal, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("value must be object (map[string]any)")
		}

		return validator(objVal)
	})

	return o
}

// AddRule adds a custom validation rule
func (o *ObjectType) AddRule(rule core.Rule) *ObjectType {
	if o.customValidation == nil {
		o.customValidation = core.NewCustomValidation()
	}
	o.customValidation.AddRule(rule)
	return o
}

// Transform, alt alanların tip dönüşümlerini uygular.
//
// Parametreler:
//   - value (any): dönüştürülecek değer
//
// Döndürür:
//   - any: dönüştürülmüş değer
//   - error: dönüşüm hatası
func (o *ObjectType) Transform(value any) (any, error) {
	value, err := o.BaseType.Transform(value)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	data, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("nesne (object) tipinde olmalıdır")
	}

	transformedData := make(map[string]any)
	for field, typ := range o.shape {
		subValue := data[field]
		transformedSubValue, err := typ.Transform(subValue)
		if err != nil {
			return nil, fmt.Errorf("alan '%s': %w", field, err)
		}
		transformedData[field] = transformedSubValue
	}
	for k, v := range data {
		if _, ok := transformedData[k]; !ok {
			transformedData[k] = v
		}
	}
	return transformedData, nil
}

// Validate, nesne ve alt alanlarının doğrulamasını gerçekleştirir.
//
// Parametreler:
//   - field (string): alan adı
//   - value (any): doğrulanacak değer
//   - result (*core.ValidationResult): doğrulama sonucu
func (o *ObjectType) Validate(field string, value any, result *core.ValidationResult) {
	o.BaseType.Validate(field, value, result)
	if result.HasErrors() {
		return
	}
	if value == nil {
		return
	}

	data, ok := value.(map[string]any)
	if !ok {
		result.AddError(field, i18n.Get(i18n.KeyObject, o.GetLabel(field)))
		return
	}

	for subField, subSchema := range o.shape {
		subValue := data[subField]
		fullFieldPath := fmt.Sprintf("%s.%s", field, subField)
		subSchema.Validate(fullFieldPath, subValue, result)
	}

	if o.customValidation != nil && o.customValidation.HasValidators() {
		o.customValidation.ValidateSync(field, value, result)
	}
}
