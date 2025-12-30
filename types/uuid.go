// -----------------------------------------------------------------------------
// UuidType: UUID Doğrulama Sınıfı
// -----------------------------------------------------------------------------
// Bu sınıf, string olarak gelen UUID değerlerini doğrulamak için kullanılır.
// Laravel/Symfony tarzında, alan bazlı doğrulama zinciri oluşturulabilir.
// Neyi, Nasıl ve Neden:
//   - Neyi: UUID tipindeki verileri doğrulamak
//   - Nasıl: Belirli bir sürüm (version) için regex ve kurallar ile
//   - Neden: UUID değerlerinin geçerli ve tutarlı olmasını sağlamak
//
// Desteklenen başlıca özellikler:
//   - Zorunlu alan, etiket, sürüm bazlı UUID doğrulama
//
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

package types

import (
	"fmt"

	"github.com/biyonik/go-fluent-validator/core"
	"github.com/biyonik/go-fluent-validator/i18n"
	"github.com/biyonik/go-fluent-validator/rules"
)

// UuidType, string tipindeki UUID değerlerini doğrulamak için kullanılır.
type UuidType struct {
	core.BaseType
	version          int
	customValidation *core.CustomValidation
}

// Required, alanın zorunlu olmasını sağlar.
func (u *UuidType) Required() *UuidType {
	u.SetRequired()
	return u
}

// Label, alan için okunabilir bir isim tanımlar.
func (u *UuidType) Label(label string) *UuidType {
	u.SetLabel(label)
	return u
}

// Version, doğrulama için kullanılacak UUID sürümünü belirler (0-5 arası).
func (u *UuidType) Version(v int) *UuidType {
	if v >= 0 && v <= 5 {
		u.version = v
	}
	return u
}

// Custom adds a custom validation function
func (u *UuidType) Custom(validator func(string) error) *UuidType {
	if u.customValidation == nil {
		u.customValidation = core.NewCustomValidation()
	}

	u.customValidation.AddSync(func(value any) error {
		if value == nil {
			return nil
		}

		strVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("value must be string")
		}

		return validator(strVal)
	})

	return u
}

// AddRule adds a custom validation rule
func (u *UuidType) AddRule(rule core.Rule) *UuidType {
	if u.customValidation == nil {
		u.customValidation = core.NewCustomValidation()
	}
	u.customValidation.AddRule(rule)
	return u
}

// Validate, UUID değerini doğrular ve hataları result'a ekler.
func (u *UuidType) Validate(field string, value any, result *core.ValidationResult) {
	u.BaseType.Validate(field, value, result)
	if result.HasErrors() {
		return
	}
	if value == nil {
		return
	}

	str, ok := value.(string)
	if !ok {
		result.AddError(field, i18n.Get(i18n.KeyString, u.GetLabel(field)))
		return
	}

	fieldName := u.GetLabel(field)
	if !rules.IsValidUUID(str, u.version) {
		result.AddError(field, i18n.Get(i18n.KeyUUID, fieldName))
	}

	if u.customValidation != nil && u.customValidation.HasValidators() {
		u.customValidation.ValidateSync(field, value, result)
	}
}
