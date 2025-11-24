package validation

import "github.com/biyonik/go-fluent-validator/types"

//
// -----------------------------------------------------------------------------
// Veri Tipi Yardımcı Fonksiyonları
// -----------------------------------------------------------------------------
// Bu dosya, Fluent Validator yapısında kullanılacak veri tipi fonksiyonlarını içerir.
// Her fonksiyon, ilgili Type struct'ını oluşturur ve başlatır.
//
// Laravel/Symfony Validator tarzında, okunabilir ve zincirleme doğrulama
// (fluent validation) mantığı sunar.
//
// Metadata:
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

// String
// -----------------------------------------------------------------------------
// Yeni bir StringType nesnesi oluşturur.
//
// Dönüş:
//   - *types.StringType → string tipli doğrulama nesnesi
func String() *types.StringType {
	return &types.StringType{}
}

// Number
// -----------------------------------------------------------------------------
// Yeni bir NumberType nesnesi oluşturur.
//
// Dönüş:
//   - *types.NumberType → sayısal doğrulama nesnesi
func Number() *types.NumberType {
	return &types.NumberType{}
}

// Boolean
// -----------------------------------------------------------------------------
// Yeni bir BooleanType nesnesi oluşturur.
//
// Dönüş:
//   - *types.BooleanType → boolean doğrulama nesnesi
func Boolean() *types.BooleanType {
	return &types.BooleanType{}
}

// Array
// -----------------------------------------------------------------------------
// Yeni bir ArrayType nesnesi oluşturur.
//
// Dönüş:
//   - *types.ArrayType → dizi doğrulama nesnesi
func Array() *types.ArrayType {
	return &types.ArrayType{}
}

// Object
// -----------------------------------------------------------------------------
// Yeni bir ObjectType nesnesi oluşturur.
//
// Dönüş:
//   - *types.ObjectType → nesne doğrulama nesnesi
func Object() *types.ObjectType {
	return &types.ObjectType{}
}

// Date
// -----------------------------------------------------------------------------
// Yeni bir DateType nesnesi oluşturur.
//
// Dönüş:
//   - *types.DateType → tarih doğrulama nesnesi
func Date() *types.DateType {
	return &types.DateType{}
}

// Uuid
// -----------------------------------------------------------------------------
// Yeni bir UuidType nesnesi oluşturur.
//
// Dönüş:
//   - *types.UuidType → UUID doğrulama nesnesi
func Uuid() *types.UuidType {
	return &types.UuidType{}
}

// Iban
// -----------------------------------------------------------------------------
// Yeni bir IbanType nesnesi oluşturur.
//
// Dönüş:
//   - *types.IbanType → IBAN doğrulama nesnesi
func Iban() *types.IbanType {
	return &types.IbanType{}
}

// CreditCard
// -----------------------------------------------------------------------------
// Yeni bir CreditCardType nesnesi oluşturur.
//
// Dönüş:
//   - *types.CreditCardType → kredi kartı doğrulama nesnesi
func CreditCard() *types.CreditCardType {
	return &types.CreditCardType{}
}

// AdvancedString
// -----------------------------------------------------------------------------
// Yeni bir AdvancedStringType nesnesi oluşturur. Daha gelişmiş string doğrulama
// ve manipülasyon özellikleri sunar.
//
// Dönüş:
//   - *types.AdvancedStringType → gelişmiş string doğrulama nesnesi
func AdvancedString() *types.AdvancedStringType {
	return &types.AdvancedStringType{}
}
