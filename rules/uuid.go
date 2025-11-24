package rules

import "regexp"

//
// -----------------------------------------------------------------------------
// UUID Doğrulama Kuralları
// -----------------------------------------------------------------------------
// Bu dosya, farklı versiyon UUID’leri doğrulamak için kullanılır.
// PHP'deki UuidValidationTrait portudur.
//
// Özellikler:
//   - UUID v1, v3, v4, v5 doğrulama
//   - Genel UUID formatı kontrolü
//   - Performans için regex global değişken olarak derlenir
//
// Metadata:
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

var (
	// UUID v1 regex
	uuidV1Regex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-1[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	// UUID v3 regex
	uuidV3Regex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	// UUID v4 regex
	uuidV4Regex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	// UUID v5 regex
	uuidV5Regex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	// Genel UUID format regex (tüm versiyonlar)
	uuidGenRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

// IsValidUUID
// -----------------------------------------------------------------------------
// Verilen string’in geçerli bir UUID olup olmadığını kontrol eder.
//
// Parametreler:
//   - uuid: doğrulanacak UUID stringi
//   - version: UUID versiyonu (1,3,4,5) veya 0 = tüm versiyonlar
//
// Dönüş:
//   - bool → geçerliyse true, değilse false
//
// Açıklama:
//   - Performans için regexler global olarak derlenmiştir
//   - Versiyon 0 seçilirse, herhangi bir UUID formatı kabul edilir
//   - PHP portudur, Go'da case-insensitive flag yerine regex küçük harf varsayıyor
func IsValidUUID(uuid string, version int) bool {
	switch version {
	case 1:
		return uuidV1Regex.MatchString(uuid)
	case 3:
		return uuidV3Regex.MatchString(uuid)
	case 4:
		return uuidV4Regex.MatchString(uuid)
	case 5:
		return uuidV5Regex.MatchString(uuid)
	case 0: // Herhangi bir versiyon (genel format)
		return uuidGenRegex.MatchString(uuid)
	default:
		return false
	}
}
