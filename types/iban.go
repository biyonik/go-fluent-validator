package types

import (
	"fmt"

	"github.com/biyonik/go-fluent-validator/core"
	"github.com/biyonik/go-fluent-validator/i18n"
	"github.com/biyonik/go-fluent-validator/rules"
)

// IbanType
//
// Uluslararası Banka Hesap Numarası (IBAN) doğrulamasını sağlayan tiptir.
// Bu tip, gelen string veriyi IBAN formatına göre kontrol eder ve opsiyonel
// olarak ülke kodu kısıtlaması getirilebilir.
//
// Özellikler:
//   - string tipinde değer alır
//   - Opsiyonel olarak belirli ülke kodu ile doğrulama yapılabilir
//   - Required() ile boş geçilemez kılınabilir
//
// Kullanım alanları:
//   - Banka hesap bilgisi doğrulaması
//   - API ve form validasyonları
//
// Yazar Bilgileri:
//   - @author  Ahmet Altun
//   - @github  https://github.com/biyonik
//   - @linkedin https://linkedin.com/in/biyonik
//   - @email   ahmet.altun60@gmail.com
type IbanType struct {
	core.BaseType
	countryCode      string
	customValidation *core.CustomValidation
}

// Required, alanın boş geçilemeyeceğini belirtir.
//
// Döndürür:
//   - *IbanType: zincirleme kullanım için aynı örnek
func (i *IbanType) Required() *IbanType {
	i.SetRequired()
	return i
}

// Label, doğrulama hatalarında gösterilecek kullanıcı dostu alan adını belirler.
//
// Parametreler:
//   - label (string): kullanıcıya gösterilecek isim
//
// Döndürür:
//   - *IbanType
func (i *IbanType) Label(label string) *IbanType {
	i.SetLabel(label)
	return i
}

// Country, IBAN doğrulamasında belirli bir ülke kodu zorunluluğu ekler.
//
// Parametreler:
//   - code (string): örn: "TR", "DE"
//
// Döndürür:
//   - *IbanType
func (i *IbanType) Country(code string) *IbanType {
	i.countryCode = code
	return i
}

// Custom adds a custom validation function
func (i *IbanType) Custom(validator func(string) error) *IbanType {
	if i.customValidation == nil {
		i.customValidation = core.NewCustomValidation()
	}

	i.customValidation.AddSync(func(value any) error {
		if value == nil {
			return nil
		}

		strVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("value must be string")
		}

		return validator(strVal)
	})

	return i
}

// AddRule adds a custom validation rule
func (i *IbanType) AddRule(rule core.Rule) *IbanType {
	if i.customValidation == nil {
		i.customValidation = core.NewCustomValidation()
	}
	i.customValidation.AddRule(rule)
	return i
}

// Validate, IBAN alanının geçerliliğini kontrol eder.
//
// İşlem sırası:
//  1. BaseType doğrulamaları (required, nullable)
//  2. Value'nin string olup olmadığı
//  3. IBAN formatı kontrolü (opsiyonel ülke kodu ile)
//
// Parametreler:
//   - field (string): alan adı
//   - value (any): doğrulanacak değer
//   - result (*core.ValidationResult): doğrulama sonucu
func (i *IbanType) Validate(field string, value any, result *core.ValidationResult) {
	i.BaseType.Validate(field, value, result)
	if result.HasErrors() || value == nil {
		return
	}

	str, ok := value.(string)
	if !ok {
		result.AddError(field, i18n.Get(i18n.KeyString, i.GetLabel(field)))
		return
	}

	if !rules.IsValidIBAN(str, i.countryCode) {
		result.AddError(field, i18n.Get(i18n.KeyIBAN, i.GetLabel(field)))
	}

	if i.customValidation != nil && i.customValidation.HasValidators() {
		i.customValidation.ValidateSync(field, value, result)
	}
}
