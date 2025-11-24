package rules

import (
	"regexp"
	"strings"
)

//
// -----------------------------------------------------------------------------
// Güvenlik ve Temizlik Kuralları
// -----------------------------------------------------------------------------
// Bu dosya, kullanıcı girdilerini temizlemek, XSS ve HTML enjeksiyonlarını önlemek,
// dosya isimlerini güvenli hale getirmek ve karakter setlerini doğrulamak için
// kullanılan fonksiyonları içerir.
//
// Özellikler:
//   - HTML etiketlerini temizleme (strip_tags taklidi)
//   - XSS önleme (PreventXss)
//   - Dosya adlarını sanitize etme
//   - Emoji filtreleme
//   - Karakter seti doğrulama (Latin, alphanumeric, numeric, alpha)
//
// Laravel/Symfony’deki InputSanitizer veya Validator mantığına benzer bir yapı
// sunar.
//
// Metadata:
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

// allowedTagsPattern ve regexler
// -----------------------------------------------------------------------------
// HTML etiketlerini ayıklamak ve izin verilen etiketleri korumak için regexler
// Emoji ve dosya isimleri için güvenli filtreler
var (
	// Tüm HTML etiketlerini yakalayan regex
	htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
	// İzin verilen etiketleri yakalamak için (basit hali)
	allowedTagRegex *regexp.Regexp

	// FilterEmoji
	emojiRegex = regexp.MustCompile(`[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}]`)

	// SanitizeFilename için: Güvensiz karakterler
	unsafeFilenameChars = regexp.MustCompile(`[^a-zA-Z0-9\-\_\.]`)
	consecutiveDots     = regexp.MustCompile(`\.{2,}`)
	leadingDot          = regexp.MustCompile(`^\.`)
	trailingDot         = regexp.MustCompile(`\.$`)

	// CharSetPatterns
	CharSetPatterns = map[string]*regexp.Regexp{
		"latin":        regexp.MustCompile(`^[\p{Latin}]+$`),
		"alphanumeric": regexp.MustCompile(`^[a-zA-Z0-9]+$`),
		"numeric":      regexp.MustCompile(`^[0-9]+$`),
		"alpha":        regexp.MustCompile(`^[a-zA-Z]+$`),
	}

	// SanitizeFilename için Türkçe karakter haritası
	turkishCharReplacer = strings.NewReplacer(
		"ç", "c", "Ç", "C",
		"ğ", "g", "Ğ", "G",
		"ı", "i", "İ", "I",
		"ö", "o", "Ö", "O",
		"ş", "s", "Ş", "S",
		"ü", "u", "Ü", "U",
	)
)

// StripHtmlTags
// -----------------------------------------------------------------------------
// Verilen string içerisindeki tüm HTML etiketlerini temizler.
// allowedTags parametresi şimdilik göz ardı edilir (PHP strip_tags taklidi).
//
// Parametreler:
//   - input: temizlenecek string
//   - allowedTags: opsiyonel, izin verilen etiketler (şimdilik desteklenmez)
//
// Dönüş:
//   - string: temizlenmiş metin
func StripHtmlTags(input string, allowedTags ...string) string {
	if len(allowedTags) == 0 {
		return htmlTagRegex.ReplaceAllString(input, "")
	}

	// TODO: allowedTags parametresini destekleyen gelişmiş HTML parser eklenmeli
	return htmlTagRegex.ReplaceAllString(input, "")
}

// PreventXss
// -----------------------------------------------------------------------------
// XSS ataklarını önlemek için temel HTML karakterlerini encode eder.
// PHP'deki preventXss fonksiyonunun portudur.
//
// Parametreler:
//   - input: XSS riski olan string
//
// Dönüş:
//   - string: güvenli string
func PreventXss(input string) string {
	input = strings.ReplaceAll(input, "&", "&amp;")
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, `"`, "&quot;")
	input = strings.ReplaceAll(input, "'", "&#39;")
	return input
}

// SanitizeFilename
// -----------------------------------------------------------------------------
// Dosya isimlerini güvenli hale getirir:
//   - Türkçe karakterleri latin harflerine çevirir
//   - Güvensiz karakterleri temizler
//   - Ardışık veya baş/son noktaları temizler
//   - Maksimum uzunluk 255 karakter ile sınırlıdır
//
// Parametreler:
//   - filename: temizlenecek dosya adı
//
// Dönüş:
//   - string: güvenli dosya adı
func SanitizeFilename(filename string) string {
	filename = turkishCharReplacer.Replace(filename)
	filename = unsafeFilenameChars.ReplaceAllString(filename, "")
	filename = consecutiveDots.ReplaceAllString(filename, ".")
	filename = leadingDot.ReplaceAllString(filename, "")
	filename = trailingDot.ReplaceAllString(filename, "")
	if len(filename) > 255 {
		filename = filename[:255]
	}
	return filename
}

// FilterEmoji
// -----------------------------------------------------------------------------
// Metin içerisindeki emojileri kaldırır veya bırakır.
//
// Parametreler:
//   - input: işlenecek metin
//   - remove: true ise emojileri temizle, false ise olduğu gibi bırak
//
// Dönüş:
//   - string: işlenmiş metin
func FilterEmoji(input string, remove bool) string {
	if remove {
		return emojiRegex.ReplaceAllString(input, "")
	}
	return input
}

// ValidateCharSet
// -----------------------------------------------------------------------------
// Verilen string’in belirli bir karakter setine uyup uymadığını kontrol eder.
//
// Parametreler:
//   - input: doğrulanacak string
//   - charSet: "latin", "alphanumeric", "numeric", "alpha"
//
// Dönüş:
//   - bool: karakter setine uyuyorsa true, değilse false
func ValidateCharSet(input string, charSet string) bool {
	pattern, ok := CharSetPatterns[charSet]
	if !ok {
		return false // Bilinmeyen karakter seti
	}
	return pattern.MatchString(input)
}
