// -----------------------------------------------------------------------------
// AdvancedStringType
// -----------------------------------------------------------------------------
// Bu dosya, gelişmiş string doğrulama ve dönüştürme ihtiyaçları için özel olarak
// tasarlanmış AdvancedStringType yapısını içerir. Laravel veya Symfony tarzı modern
// framework'lerin sunduğu esnek doğrulama olanaklarını Go ekosistemine taşımak
// amacıyla geliştirilmiştir.
//
// Bu yapı sayesinde:
// - HTML etiketlerini temizleme,
// - XSS’e karşı HTML kaçışlama,
// - Emoji filtreleme,
// - Dosya adı sanitize etme,
// - Türkçe karakter zorunluluğu veya yasağı,
// - Domain format kontrolü,
// - Özel karakter seti doğrulama
//
// gibi gelişmiş string operasyonları, akıcı ve zincirlenebilir bir API ile kolayca
// uygulanabilir hale gelir.
//
// AdvancedStringType, temel StringType'ın tüm özelliklerini devralır ve üzerine
// dönüşüm (transform) ve doğrulama (validation) aşamalarının genişletilmiş
// versiyonlarını ekler. Böylece bir form inputunun hem normalize edilmesi hem de
// iş kurallarına uygunluğunun kontrolü tek bir doğrulama akışı içinde çözülebilir.
//
// Bu sınıf, yüksek güvenlik gerektiren API'larda, kullanıcı girdilerinin
// temizlenmesinin zorunlu olduğu durumlarda veya içerik filtreleme ihtiyacının
// bulunduğu projelerde kritik rol oynar.
//
// Yazar Bilgileri:
// @author   Ahmet ALTUN
// @github   github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email    ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

package types

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/biyonik/go-fluent-validator/core"
	"github.com/biyonik/go-fluent-validator/i18n"
	"github.com/biyonik/go-fluent-validator/rules"
)

// AdvancedStringType, gelişmiş doğrulama ve transform işlemleri için kullanılan
// string türüdür. Standart StringType üzerine ek güvenlik ve formatlama kuralları
// getirmek için tasarlanmıştır.
type AdvancedStringType struct {
	StringType
	turkishChars *bool   // Türkçe karakter içermeli mi / içermemeli mi?
	domainCheck  *bool   // Domain doğrulaması yapılacak mı?
	charSet      *string // Belirli bir karakter seti zorunluluğu
	maxWords     *int
}

// StripTags, verilen string içindeki HTML etiketlerini (izin verilenler hariç)
// temizler. Kullanıcı girdisini normalize etmek için sıkça kullanılır.
func (as *AdvancedStringType) StripTags(allowedTags ...string) *AdvancedStringType {
	as.AddTransform(func(value any) (any, error) {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("StripTags sadece string'lere uygulanabilir")
		}
		return rules.StripHtmlTags(str, allowedTags...), nil
	})
	return as
}

// EscapeHTML, XSS saldırılarına karşı HTML karakterlerini güvenli hale getirir.
// Güvenlik seviyesi yüksek projelerde mutlaka kullanılmalıdır.
func (as *AdvancedStringType) EscapeHTML() *AdvancedStringType {
	as.AddTransform(func(value any) (any, error) {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("EscapeHTML sadece string'lere uygulanabilir")
		}
		return rules.PreventXss(str), nil
	})
	return as
}

// SanitizeFilename, bir dosya adını güvenli hale getirir. Sistem çağrılarına veya
// dosya işlemlerine zarar verebilecek karakterleri temizler.
func (as *AdvancedStringType) SanitizeFilename() *AdvancedStringType {
	as.AddTransform(func(value any) (any, error) {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("SanitizeFilename sadece string'lere uygulanabilir")
		}
		return rules.SanitizeFilename(str), nil
	})
	return as
}

// FilterEmoji, metindeki emoji karakterlerini isteğe bağlı olarak temizler veya
// filtreler. Chat sistemi, loglama veya özel karakter sınırlaması olan sistemlerde
// kullanılır.
func (as *AdvancedStringType) FilterEmoji(remove bool) *AdvancedStringType {
	as.AddTransform(func(value any) (any, error) {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("FilterEmoji sadece string'lere uygulanabilir")
		}
		return rules.FilterEmoji(str, remove), nil
	})
	return as
}

// TurkishChars, string'in Türkçe karakter içerip içermemesi gerektiğini belirler.
// allow=true => Türkçe karakter zorunlu
// allow=false => Türkçe karakter yasak
func (as *AdvancedStringType) TurkishChars(allow bool) *AdvancedStringType {
	as.turkishChars = &allow
	return as
}

// Domain, string'in geçerli bir domain olup olmadığını doğrular.
// allowSubdomain=true => alt domainlere izin verir.
func (as *AdvancedStringType) Domain(allowSubdomain bool) *AdvancedStringType {
	as.domainCheck = &allowSubdomain
	return as
}

// CharSet, bu string'in belirli bir karakter setine uygun olması zorunluluğunu ayarlar.
// Örn: "alpha", "alphanumeric", "numeric", "hex" vb.
func (as *AdvancedStringType) CharSet(set string) *AdvancedStringType {
	as.charSet = &set
	return as
}

// MaxWords, metindeki maksimum kelime sayısını sınırlar.
func (as *AdvancedStringType) MaxWords(count int) *AdvancedStringType {
	as.maxWords = &count
	return as
}

// StripPunctuation, noktalama işaretlerini temizler.
func (as *AdvancedStringType) StripPunctuation() *AdvancedStringType {
	as.AddTransform(func(value any) (any, error) {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("StripPunctuation sadece string'lere uygulanabilir")
		}
		// Noktalama işaretlerini kaldır
		reg := regexp.MustCompile(`[^\w\s]`)
		return reg.ReplaceAllString(str, ""), nil
	})
	return as
}

// ReplaceTurkishChars, Türkçe karakterleri ASCII karşılıklarıyla değiştirir.
// ş→s, ğ→g, ü→u, ö→o, ç→c, İ→I, ı→i
func (as *AdvancedStringType) ReplaceTurkishChars() *AdvancedStringType {
	as.AddTransform(func(value any) (any, error) {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("ReplaceTurkishChars sadece string'lere uygulanabilir")
		}

		replacer := strings.NewReplacer(
			"ş", "s", "Ş", "S",
			"ğ", "g", "Ğ", "G",
			"ü", "u", "Ü", "U",
			"ö", "o", "Ö", "O",
			"ç", "c", "Ç", "C",
			"ı", "i", "İ", "I",
		)
		return replacer.Replace(str), nil
	})
	return as
}

// ValidateDomain, string'in geçerli bir domain olup olmadığını doğrular.
func (as *AdvancedStringType) ValidateDomain() *AdvancedStringType {
	trueVal := true
	as.domainCheck = &trueVal
	return as
}

// Trim override - AdvancedStringType için
func (as *AdvancedStringType) Trim() *AdvancedStringType {
	as.StringType.Trim()
	return as
}

// Required override
func (as *AdvancedStringType) Required() *AdvancedStringType {
	as.StringType.Required()
	return as
}

// Label override
func (as *AdvancedStringType) Label(label string) *AdvancedStringType {
	as.StringType.Label(label)
	return as
}

// Min override
func (as *AdvancedStringType) Min(length int) *AdvancedStringType {
	as.StringType.Min(length)
	return as
}

// Max override
func (as *AdvancedStringType) Max(length int) *AdvancedStringType {
	as.StringType.Max(length)
	return as
}

// Validate, tüm gelişmiş string kurallarını uygular ve olası hataları ValidationResult'a
// ekler. Önce temel StringType doğrulaması yapılır, ardından gelişmiş kontroller çalışır.
func (as *AdvancedStringType) Validate(field string, value any, result *core.ValidationResult) {
	as.StringType.Validate(field, value, result)
	if result.HasErrors() || value == nil {
		return
	}

	str, _ := value.(string)
	fieldName := as.GetLabel(field)

	if as.maxWords != nil {
		words := strings.Fields(str)
		if len(words) > *as.maxWords {
			result.AddError(field, i18n.Get(i18n.KeyMaxWords, fieldName, *as.maxWords))
		}
	}

	if as.turkishChars != nil {
		hasTurkish := rules.HasTurkishChars(str)
		if *as.turkishChars && !hasTurkish {
			result.AddError(field, i18n.Get(i18n.KeyRequiredTurkishChars, fieldName, *as.turkishChars))
		} else if !*as.turkishChars && hasTurkish {
			result.AddError(field, i18n.Get(i18n.KeyDontRequiredTurkishChars, fieldName, *as.turkishChars))
		}
	}

	if as.domainCheck != nil {
		if !rules.IsValidDomain(str, *as.domainCheck) {
			result.AddError(field, i18n.Get(i18n.KeyMustBeValidDomain, fieldName, *as.domainCheck))
		}
	}

	if as.charSet != nil {
		if !rules.ValidateCharSet(str, *as.charSet) {
			result.AddError(field, i18n.Get(i18n.KeyFieldMustMatchField, fieldName, *as.charSet))
		}
	}
}
