package types

import (
	"fmt"

	"github.com/biyonik/go-fluent-validator/core"
	"github.com/biyonik/go-fluent-validator/rules"
)

// CreditCardType
//
// Kredi kartı numaralarının doğrulanmasını sağlayan gelişmiş tip sınıfıdır.
// Bu tip, BaseType’ın tüm temel doğrulama özelliklerini devralır ve bunlara
// ek olarak kredi kartı numarası formatlarını, uzunluk yapısını ve Luhn
// algoritmasını kontrol eder.
//
// Ayrıca opsiyonel olarak yalnızca belirli bir kart tipine (visa, mastercard,
// amex vb.) ait numaraların doğrulanması sağlanabilir.
//
// Sınıf Özellikleri:
//   - cardType : (string) Yalnızca belirli bir kart markasını doğrulamak için kullanılır.
//     Örneğin: "visa", "mastercard", "amex" vb.
//
// Kullanım Örneği:
//
//	v := validator.New()
//	v.Field("card", types.CreditCard().Required().Type("visa"))
//
// Yazar Bilgileri:
//   - @author  Ahmet Altun
//   - @github  https://github.com/biyonik
//   - @company Biyonik Software
//   - @email   admin@biyonik.dev
type CreditCardType struct {
	core.BaseType
	cardType string // Örn: "visa", "mastercard", "amex" vb.
}

// Required işareti, alanın boş bırakılamayacağını belirtir.
//
// Döndürür:
//   - *CreditCardType : Zincirleme yapı için aynı nesneyi döndürür.
func (c *CreditCardType) Required() *CreditCardType {
	c.SetRequired()
	return c
}

// Label, doğrulama hatalarında görünecek kullanıcı dostu adı belirler.
//
// Parametreler:
//   - label (string): Hatalarda görüntülenecek alan adı.
//
// Döndürür:
//   - *CreditCardType : Zincirleme yapı için aynı örneği döndürür.
func (c *CreditCardType) Label(label string) *CreditCardType {
	c.SetLabel(label)
	return c
}

// Type, yalnızca belirli bir kart markasına ait kredi kartı numarasının
// kabul edilmesini sağlar.
//
// Örnek: .Type("visa") → sadece Visa kartlarına izin verilir.
//
// Parametreler:
//   - typ (string): Kart tipi. Örneğin: "visa", "mastercard", "amex".
//
// Döndürür:
//   - *CreditCardType : Zincirleme yapı için aynı örneği döndürür.
func (c *CreditCardType) Type(typ string) *CreditCardType {
	c.cardType = typ
	return c
}

// Validate, kredi kartı numarasının geçerliliğini kontrol eder.
//
// Gerçekleştirilen kontroller:
//  1. BaseType doğrulama (required, nullable kontrolü, dönüşüm vb.)
//  2. Değerin string olup olmadığı
//  3. Kart numarasının format olarak geçerli olup olmadığı
//  4. Luhn algoritması kontrolü
//  5. Eğer kart tipi tanımlanmışsa, yalnızca o markaya uygunluğunun doğrulanması
//
// Parametreler:
//   - field (string)                : Alan adı (path)
//   - value (any)                   : Doğrulanacak değer
//   - result (*core.ValidationResult): Doğrulama sonuç nesnesi
func (c *CreditCardType) Validate(field string, value any, result *core.ValidationResult) {

	// BaseType doğrulamalarını uygula
	c.BaseType.Validate(field, value, result)

	if result.HasErrors() {
		return
	}
	if value == nil {
		return
	}

	// Değerin string olması gerekir
	str, ok := value.(string)
	if !ok {
		result.AddError(field,
			fmt.Sprintf("%s alanı metin tipinde olmalıdır", c.GetLabel(field)))
		return
	}

	// Kredi kartı doğrulaması (format + Luhn + kart markası kontrolü)
	if !rules.IsValidCreditCard(str, c.cardType) {
		result.AddError(field,
			fmt.Sprintf("%s alanı geçerli bir kredi kartı numarası olmalıdır", c.GetLabel(field)))
	}
}
