package rules

import (
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

//
// -----------------------------------------------------------------------------
// Ödeme ve Finans Kuralları
// -----------------------------------------------------------------------------
// Bu dosya, kredi kartı ve IBAN doğrulama gibi finansal işlemlerde kullanılan
// kuralları merkezi bir yerde toplar. Laravel/Symfony’deki PaymentValidator
// mantığının Go portudur.
//
// Öne çıkan özellikler:
//   - Luhn algoritması ile kredi kartı doğrulaması
//   - Kart tipi (Visa, Mastercard, Amex) kontrolü
//   - IBAN doğrulaması (ülke kodu ve uzunluk kontrolü)
//   - Büyük sayılar için math/big kullanımı
//
// Bu yapılar, ödeme işlemleri, online form doğrulamaları, banka API entegrasyonları
// gibi kritik finansal alanlarda güvenli ve standart bir doğrulama altyapısı sağlar.
//
// Metadata:
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

// luhnCheck
// -----------------------------------------------------------------------------
// Verilen sayısal string’in Luhn algoritmasına uygun olup olmadığını kontrol eder.
//
// Luhn algoritması:
//   - Kredi kartı numaralarının doğrulanmasında kullanılan standart bir yöntemdir.
//   - PHP portudur.
//
// Geri dönüş:
//   - true → geçerli
//   - false → geçersiz
func luhnCheck(number string) bool {
	var sum int
	isEvenIndex := len(number)%2 == 0

	for _, digitChar := range number {
		digit, err := strconv.Atoi(string(digitChar))
		if err != nil {
			return false // Sayısal olmayan karakter
		}

		if isEvenIndex {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isEvenIndex = !isEvenIndex
	}
	return sum%10 == 0
}

// IsValidCreditCard
// -----------------------------------------------------------------------------
// Verilen kredi kartı numarasının geçerli olup olmadığını kontrol eder.
//
// Parametreler:
//   - cardNumber: doğrulanacak kredi kartı numarası
//   - cardType: opsiyonel, "visa", "mastercard", "amex" gibi türü belirtir
//
// Dönüş:
//   - bool → geçerliyse true, değilse false
//
// Açıklama:
//   - Önce boşluk ve tireler temizlenir
//   - Kart tipi regex ile kontrol edilir
//   - Luhn algoritması ile sayısal doğrulama yapılır
func IsValidCreditCard(cardNumber string, cardType string) bool {
	// Boşluk ve tireleri kaldır (PHP'deki preg_replace)
	number := regexp.MustCompile(`\D`).ReplaceAllString(cardNumber, "")

	// Kart tipi desenleri
	patterns := map[string]*regexp.Regexp{
		"visa":       regexp.MustCompile(`^4[0-9]{12}(?:[0-9]{3})?$`),
		"mastercard": regexp.MustCompile(`^5[1-5][0-9]{14}$`),
		"amex":       regexp.MustCompile(`^3[47][0-9]{13}$`),
	}

	// Kart tipi kontrolü
	if cardType != "" {
		pattern, ok := patterns[cardType]
		if !ok || !pattern.MatchString(number) {
			return false // Belirtilen tip değil
		}
	}

	// Luhn algoritması
	return luhnCheck(number)
}

// IsValidIBAN
// -----------------------------------------------------------------------------
// Verilen IBAN numarasının geçerli olup olmadığını kontrol eder.
//
// Parametreler:
//   - iban: doğrulanacak IBAN stringi
//   - countryCode: opsiyonel, ülke kodu (örn: "TR", "DE").
//
// Dönüş:
//   - bool → geçerliyse true, değilse false
//
// Açıklama:
//   - Boşluklar temizlenir ve büyük harfe çevrilir
//   - Ülke kodu verilmişse uzunluk kontrolü yapılır
//   - Regex ile format kontrolü yapılır
//   - IBAN 4 karakter öne alınarak rakamsal forma çevrilir
//   - math/big ile büyük sayı mod 97 işlemi yapılır
//   - PHP’deki bcmod karşılığıdır
func IsValidIBAN(iban string, countryCode string) bool {
	iban = strings.ToUpper(strings.ReplaceAll(iban, " ", ""))

	if countryCode != "" {
		expectedLength, ok := ibanCountryLengths[countryCode]
		if !ok || len(iban) != expectedLength {
			return false
		}
	}

	if match, _ := regexp.MatchString(`^[A-Z]{2}\d{2}[A-Z0-9]{4,}$`, iban); !match {
		return false
	}

	rearranged := iban[4:] + iban[:4]

	converted := ""
	for _, char := range rearranged {
		if unicode.IsLetter(char) {
			converted += strconv.Itoa(int(char - 'A' + 10))
		} else {
			converted += string(char)
		}
	}

	ibanNum := new(big.Int)
	ibanNum, ok := ibanNum.SetString(converted, 10)
	if !ok {
		return false
	}

	mod97 := new(big.Int)
	mod97.SetInt64(97)

	remainder := new(big.Int)
	remainder.Mod(ibanNum, mod97)

	return remainder.Int64() == 1
}

// ibanCountryLengths
// -----------------------------------------------------------------------------
// Ülke kodlarına göre IBAN uzunluklarını tutan harita.
// Bu yapı ile her ülke için doğru IBAN uzunluğu kontrol edilebilir.
var ibanCountryLengths = map[string]int{
	"TR": 26, "DE": 22, "GB": 22, "FR": 27, "IT": 27, "NL": 18,
}
