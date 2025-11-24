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
	"github.com/biyonik/go-fluent-validator/rules"
)

// UuidType, string tipindeki UUID değerlerini doğrulamak için kullanılır.
type UuidType struct {
	core.BaseType
	version int
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
		result.AddError(field, fmt.Sprintf("%s alanı metin tipinde olmalıdır", u.GetLabel(field)))
		return
	}

	fieldName := u.GetLabel(field)
	if !rules.IsValidUUID(str, u.version) {
		result.AddError(field, fmt.Sprintf("%s alanı geçerli bir UUID olmalıdır", fieldName))
	}
}
