package types

import (
	"fmt"

	"github.com/biyonik/go-fluent-validator/core"
)

// -----------------------------------------------------------------------------
// BooleanType
// -----------------------------------------------------------------------------
// Bu sınıf, doğrulama katmanında boolean (true/false) türündeki alanların
// kontrol edilmesi için özel olarak tasarlanmış bir veri tipi doğrulayıcısıdır.
// Özellikle form işlemleri, API veri alımları ve model doğrulamalarında boolean
// değerlerin beklenen formatta iletilip iletilmediğini güvenli biçimde tespit
// eder.
//
// Neyi yapar?
// - Bir alanın zorunlu olup olmadığını kontrol eder.
// - Varsayılan değer atayabilir.
// - Gelen değerin gerçek bir boolean olup olmadığını doğrular.
// - Hatalı tip ile karşılaşıldığında okunaklı, Türkçe hata mesajı üretir.
//
// Neden vardır?
// Yazılım sistemlerinde boolean değerler sıklıkla kullanıcı girişlerinden veya
// JSON payload'larından geldiği için metinsel veya sayısal olarak hatalı
// aktarılabilmektedir. BooleanType, bu tür olası hataları merkezî bir yapı
// içinde güvenle ele almak için geliştirilmiştir.
//
// Nasıl çalışır?
// BaseType üzerinden aldığı altyapı ile önce "gereklilik", "label", "default"
// gibi temel validasyon kurallarını işletir. Ardından gelen değerin gerçek bir
// boolean olup olmadığını Go'nun type assertion yöntemi ile kontrol eder. Eğer
// tip uyuşmazlığı varsa ValidationResult içerisine anlamlı bir hata mesajı
// ekler.
//
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------
type BooleanType struct {
	core.BaseType
}

// Required, ilgili boolean alanın zorunlu olduğunu işaretler.
// Eğer değer sağlanmazsa BaseType içindeki zorunluluk kontrolü hata üretir.
func (b *BooleanType) Required() *BooleanType {
	b.SetRequired()
	return b
}

// Label, alan için okunabilir ve anlamlı bir isim belirler.
// Bu label, hata mesajlarında kullanıcıya daha anlaşılır geri bildirim vermek
// için kullanılır.
func (b *BooleanType) Label(label string) *BooleanType {
	b.SetLabel(label)
	return b
}

// Default, bu boolean alan için bir varsayılan değer tanımlar.
// Veri gelmediğinde veya boş olduğunda bu değer otomatik olarak atanır.
func (b *BooleanType) Default(value bool) *BooleanType {
	b.SetDefault(value)
	return b
}

// Validate, ilgili alanın doğrulama sürecini yürütür.
// 1. BaseType doğrulama kurallarını çalıştırır (required, default, label...).
// 2. Gelen değer nil ise (zorunlu değilse) işlem durdurulur.
// 3. Gelen değerin gerçek bir boolean olup olmadığı kontrol edilir.
// 4. Tip uyumsuzluğunda ValidationResult içine kullanıcı dostu bir hata eklenir.
func (b *BooleanType) Validate(field string, value any, result *core.ValidationResult) {
	b.BaseType.Validate(field, value, result)
	if result.HasErrors() {
		return
	}

	if value == nil {
		return
	}

	_, ok := value.(bool)
	if !ok {
		result.AddError(field, fmt.Sprintf("%s alanı boolean tipinde olmalıdır", b.GetLabel(field)))
	}
}
