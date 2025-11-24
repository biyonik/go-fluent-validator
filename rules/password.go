package rules

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

//
// -----------------------------------------------------------------------------
// PasswordRules ve Şifre Doğrulama Kuralları
// -----------------------------------------------------------------------------
// Bu dosya, şifre güvenliği ve doğrulama için gerekli tüm kuralları merkezi bir
// yerde toplar. Laravel ve Symfony’deki password validator mantığını Go’ya uyarlayan
// bir yapıdır.
//
// PasswordRules struct’ı, şifrelerin güvenlik politikalarını belirler:
//   - Minimum/maximum uzunluk
//   - Büyük/küçük harf, rakam ve özel karakter zorunluluğu
//   - Tekrarlayan karakter sınırlaması
//   - Klavye dizilim kontrolleri
//   - Ortak/kolay tahmin edilebilir şifrelerin engellenmesi
//   - Entropi (karmaşıklık) kontrolü
//
// Bu yapı sayesinde, uygulamada tüm şifre doğrulama işlemleri standart, güvenli ve
// kolay genişletilebilir hale gelir.
//
// Metadata:
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

// PasswordRules, PHP'deki $passwordRules dizisinin Go struct karşılığıdır.
type PasswordRules struct {
	MinLength         int
	MaxLength         int
	RequireUppercase  bool
	RequireLowercase  bool
	RequireNumeric    bool
	RequireSpecial    bool
	SpecialChars      string
	MinUniqueChars    int
	MaxRepeatingChars int
	DisallowCommon    bool
	DisallowKeyboard  bool
	MinEntropy        float64
}

// commonPasswords
// -----------------------------------------------------------------------------
// Çok yaygın kullanılan şifreleri listeleyen sabit harita.
// Bu şifreler engellenerek güvenlik artırılır.
var commonPasswords = map[string]bool{
	"password": true, "123456": true, "qwerty": true, "111111": true, "abc123": true,
	"letmein": true, "admin": true, "welcome": true, "monkey": true, "dragon": true,
}

// keyboardPatterns
// -----------------------------------------------------------------------------
// Klavye dizilimlerine göre kolay tahmin edilebilecek şifreleri tespit etmek için
// kullanılan diziler.
var keyboardPatterns = []string{
	"qwerty", "asdfgh", "zxcvbn",
	"123456", "654321",
	"abc", "cba", "xyz",
}

// stringReverse
// -----------------------------------------------------------------------------
// PHP’deki strrev() fonksiyonunun Go karşılığıdır.
// Verilen string’in karakterlerini tersine çevirir.
func stringReverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// hasKeyboardPattern
// -----------------------------------------------------------------------------
// Şifrede klavye düzeni ile sıralı karakterlerin bulunup bulunmadığını kontrol eder.
// PHP portudur. Dizilim ters okunsa dahi tespit edilir.
func hasKeyboardPattern(password string) bool {
	loweredPass := strings.ToLower(password)
	for _, pattern := range keyboardPatterns {
		if strings.Contains(loweredPass, pattern) || strings.Contains(loweredPass, stringReverse(pattern)) {
			return true
		}
	}
	return false
}

// hasRepeatingChars
// -----------------------------------------------------------------------------
// Şifrede ardışık tekrar eden karakterlerin belirtilen sınırı aşıp aşmadığını kontrol eder.
// maxRepeats <= 0 ise kontrol yapılmaz.
func hasRepeatingChars(password string, maxRepeats int) bool {
	if maxRepeats <= 0 {
		return false
	}
	chars := []rune(password)
	consecutive := 1
	var lastChar rune

	for i, char := range chars {
		if i > 0 && char == lastChar {
			consecutive++
			if consecutive > maxRepeats {
				return true
			}
		} else {
			consecutive = 1
		}
		lastChar = char
	}
	return false
}

// calculatePasswordEntropy
// -----------------------------------------------------------------------------
// Şifrenin tahmin edilebilirlik/karmaşıklık düzeyini hesaplar (entropi).
// PHP portudur, Go’da math.Log2 ile uygulanır.
// Harf, rakam ve özel karakter çeşitliliğine göre entropi hesaplanır.
func calculatePasswordEntropy(password string) float64 {
	length := float64(len(password))
	if length == 0 {
		return 0
	}
	charPool := 0.0

	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		charPool += 26
	}
	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		charPool += 26
	}
	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		charPool += 10
	}
	if regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		charPool += 32 // PHP'deki varsayılan özel karakter sayısı
	}

	if charPool == 0 {
		return 0
	}

	// PHP'deki log(pool, 2) Go'da math.Log2(pool) demektir.
	return length * math.Log2(charPool)
}

// ValidatePassword
// -----------------------------------------------------------------------------
// Verilen şifreyi belirtilen kurallara göre doğrular ve hata mesajlarını döndürür.
//
// Parametreler:
//   - password: doğrulanacak şifre
//   - rules: uygulanacak şifre kuralları
//
// Dönüş:
//   - []string: tüm hata mesajları, kurallar sağlanıyorsa boş slice
//
// Açıklama:
// - Minimum ve maksimum uzunluk kontrolleri
// - Büyük, küçük harf, rakam ve özel karakter kontrolleri
// - Benzersiz karakter sayısı ve tekrar eden karakterler
// - Klavye düzeni ve yaygın şifre kontrolü
// - Entropi (karmaşıklık) kontrolü
func ValidatePassword(password string, rules *PasswordRules) []string {
	errors := []string{}
	if rules == nil {
		return errors
	}

	passLen := len(password)

	if passLen < rules.MinLength {
		errors = append(errors, fmt.Sprintf("en az %d karakter uzunluğunda olmalıdır", rules.MinLength))
	}
	if passLen > rules.MaxLength {
		errors = append(errors, fmt.Sprintf("en fazla %d karakter uzunluğunda olmalıdır", rules.MaxLength))
	}
	if rules.RequireUppercase && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		errors = append(errors, "en az bir büyük harf içermelidir")
	}
	if rules.RequireLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
		errors = append(errors, "en az bir küçük harf içermelidir")
	}
	if rules.RequireNumeric && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		errors = append(errors, "en az bir rakam içermelidir")
	}
	if rules.RequireSpecial {
		specialChars := regexp.QuoteMeta(rules.SpecialChars)
		if !regexp.MustCompile(fmt.Sprintf("[%s]", specialChars)).MatchString(password) {
			errors = append(errors, fmt.Sprintf("en az bir özel karakter içermelidir (%s)", rules.SpecialChars))
		}
	}

	// Gelişmiş kontroller
	uniqueChars := make(map[rune]bool)
	for _, char := range password {
		uniqueChars[char] = true
	}
	if len(uniqueChars) < rules.MinUniqueChars {
		errors = append(errors, fmt.Sprintf("en az %d farklı karakter içermelidir", rules.MinUniqueChars))
	}

	if rules.DisallowKeyboard && hasKeyboardPattern(password) {
		errors = append(errors, "klavye düzeninde sıralı karakterler içeremez")
	}
	if rules.MaxRepeatingChars > 0 && hasRepeatingChars(password, rules.MaxRepeatingChars) {
		errors = append(errors, fmt.Sprintf("en fazla %d adet tekrar eden karakter içerebilir", rules.MaxRepeatingChars))
	}
	if rules.DisallowCommon && commonPasswords[strings.ToLower(password)] {
		errors = append(errors, "çok yaygın bir şifre, lütfen daha güvenli bir şifre seçin")
	}

	entropy := calculatePasswordEntropy(password)
	if entropy < rules.MinEntropy {
		errors = append(errors, "yeterince karmaşık değil, lütfen daha güçlü bir şifre seçin")
	}

	return errors
}
