// -----------------------------------------------------------------------------
// StringType: Metin Doğrulama ve Dönüştürme Sınıfı
// -----------------------------------------------------------------------------
// Bu sınıf, string değerler üzerinde kapsamlı doğrulama ve dönüşüm işlevleri sağlar.
// Laravel/Symfony tarzında, her alan için esnek ve zincirlenebilir kurallar oluşturur.
// Neyi, Nasıl ve Neden:
//   - Neyi: Metin verilerini doğrulamak, düzenlemek ve standartlaştırmak
//   - Nasıl: Regex, kurallar ve dönüşüm fonksiyonları ile
//   - Neden: Kullanıcı girdilerini güvenli ve belirlenen standartlara uygun hale getirmek
//
// Desteklenen başlıca özellikler:
//   - Min/Max uzunluk, e-posta doğrulama, URL doğrulama, IP doğrulama
//   - Telefon numarası, şifre doğrulama, HTML etiket temizleme, trim işlemi
//
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

package types

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/biyonik/go-fluent-validator/core"
	"github.com/biyonik/go-fluent-validator/i18n"
	"github.com/biyonik/go-fluent-validator/rules"
)

var (
	emailRegex        = regexp.MustCompile(`^[a-zA-Z0-9]+([._+-][a-zA-Z0-9]+)*@[a-zA-Z0-9]+([.-][a-zA-Z0-9]+)*\.[a-zA-Z]{2,}$`)
	urlRegex          = regexp.MustCompile(`^https?://[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)+(:[0-9]+)?(/[^\s]*)?(\?[^\s]*)?$`)
	alphaRegex        = regexp.MustCompile(`^[a-zA-Z]+$`)
	alphanumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	numericRegex      = regexp.MustCompile(`^[0-9]+$`)
	macRegex          = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	hexRegex          = regexp.MustCompile(`^[0-9A-Fa-f]+$`)
	base64Regex       = regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`)
)

// StringType, string tipindeki veriler için doğrulama ve dönüşüm kurallarını tutar.
type StringType struct {
	core.BaseType
	minLength        *int
	maxLength        *int
	emailRegex       *regexp.Regexp
	urlRegex         *regexp.Regexp
	allowedValues    []string
	passwordRules    *rules.PasswordRules
	ipVersion        *int
	phoneCountry     *string
	customValidation *core.CustomValidation
	// New validators
	isAlpha          bool
	isAlphanumeric   bool
	isNumeric        bool
	startsWith       *string
	endsWith         *string
	contains         *string
	customRegex      *regexp.Regexp
	regexError       error
	isMAC            bool
	isHex            bool
	isBase64         bool
}

// Required, alanın zorunlu olmasını sağlar.
func (s *StringType) Required() *StringType {
	s.SetRequired()
	return s
}

// Label, alan için okunabilir bir isim tanımlar.
func (s *StringType) Label(label string) *StringType {
	s.SetLabel(label)
	return s
}

// Default, alan için varsayılan değer belirler.
func (s *StringType) Default(value string) *StringType {
	s.SetDefault(value)
	return s
}

// Min, string için minimum uzunluğu ayarlar.
func (s *StringType) Min(length int) *StringType {
	s.minLength = &length
	return s
}

// Max, string için maksimum uzunluğu ayarlar.
func (s *StringType) Max(length int) *StringType {
	s.maxLength = &length
	return s
}

// Email, alanın e-posta formatında olmasını sağlar.
func (s *StringType) Email() *StringType {
	s.emailRegex = emailRegex
	return s
}

// URL, alanın URL formatında olmasını sağlar.
func (s *StringType) URL() *StringType {
	s.urlRegex = urlRegex
	return s
}

// OneOf, alanın belirli bir değer listesi içinde olmasını sağlar.
func (s *StringType) OneOf(values []string) *StringType {
	s.allowedValues = values
	return s
}

// IP, alanın IP adresi formatında olmasını sağlar.
func (s *StringType) IP(version ...int) *StringType {
	v := 0
	if len(version) > 0 {
		v = version[0]
	}
	s.ipVersion = &v
	return s
}

// Phone, alanın belirli ülkeye ait telefon numarası formatında olmasını sağlar.
func (s *StringType) Phone(countryCode string) *StringType {
	s.phoneCountry = &countryCode
	return s
}

// Trim, string değerlerin başındaki ve sonundaki boşlukları temizler.
func (s *StringType) Trim() *StringType {
	s.AddTransform(func(value any) (any, error) {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("Trim sadece string değerler için uygulanabilir")
		}
		return strings.TrimSpace(str), nil
	})
	return s
}

// StripTags, HTML etiketlerini temizler, istenen etiketleri bırakabilir.
func (s *StringType) StripTags(allowedTags ...string) *StringType {
	s.AddTransform(func(value any) (any, error) {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("StripTags sadece string değerler için uygulanabilir")
		}
		return rules.StripHtmlTags(str, allowedTags...), nil
	})
	return s
}

// Password, alanın şifre doğrulama kurallarına uymasını sağlar.
func (s *StringType) Password(options ...PasswordOption) *StringType {
	defaults := &rules.PasswordRules{
		MinLength:         8,
		MaxLength:         72,
		RequireUppercase:  true,
		RequireLowercase:  true,
		RequireNumeric:    true,
		RequireSpecial:    true,
		SpecialChars:      `!@#$%^&*(),.?":{}|<>+-`,
		MinUniqueChars:    6,
		MaxRepeatingChars: 3,
		DisallowCommon:    true,
		DisallowKeyboard:  true,
		MinEntropy:        50.0,
	}
	for _, option := range options {
		option(defaults)
	}
	s.passwordRules = defaults
	return s
}

// Custom adds a custom validation function
func (s *StringType) Custom(validator func(string) error) *StringType {
	if s.customValidation == nil {
		s.customValidation = core.NewCustomValidation()
	}

	s.customValidation.AddSync(func(value any) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("value must be string")
		}
		return validator(str)
	})

	return s
}

// Alpha ensures the string contains only alphabetic characters
func (s *StringType) Alpha() *StringType {
	s.isAlpha = true
	return s
}

// Alphanumeric ensures the string contains only alphanumeric characters
func (s *StringType) Alphanumeric() *StringType {
	s.isAlphanumeric = true
	return s
}

// Numeric ensures the string contains only numeric characters
func (s *StringType) Numeric() *StringType {
	s.isNumeric = true
	return s
}

// StartsWith ensures the string starts with a specific prefix
func (s *StringType) StartsWith(prefix string) *StringType {
	s.startsWith = &prefix
	return s
}

// EndsWith ensures the string ends with a specific suffix
func (s *StringType) EndsWith(suffix string) *StringType {
	s.endsWith = &suffix
	return s
}

// Contains ensures the string contains a specific substring
func (s *StringType) Contains(substring string) *StringType {
	s.contains = &substring
	return s
}

// Regex validates the string against a custom regular expression
func (s *StringType) Regex(pattern string) *StringType {
	var err error
	s.customRegex, err = regexp.Compile(pattern)
	if err != nil {
		s.regexError = fmt.Errorf("invalid regex pattern: %w", err)
	}
	return s
}

// MAC ensures the string is a valid MAC address
func (s *StringType) MAC() *StringType {
	s.isMAC = true
	return s
}

// Hex ensures the string is a valid hexadecimal string
func (s *StringType) Hex() *StringType {
	s.isHex = true
	return s
}

// Base64 ensures the string is a valid base64 encoded string
func (s *StringType) Base64() *StringType {
	s.isBase64 = true
	return s
}

// Validate, string değer üzerinde tüm kuralları uygular ve hata durumlarını result'a ekler.
func (s *StringType) Validate(field string, value any, result *core.ValidationResult) {
	s.BaseType.Validate(field, value, result)
	if result.HasErrors() {
		return
	}
	if value == nil {
		return
	}

	str, ok := value.(string)
	if !ok {
		result.AddError(field, i18n.Get(i18n.KeyString, s.GetLabel(field)))
		return
	}

	fieldName := s.GetLabel(field)

	if s.minLength != nil && len(str) < *s.minLength {
		result.AddError(field, i18n.Get(i18n.KeyMinLength, fieldName, *s.minLength))
	}

	if s.maxLength != nil && len(str) > *s.maxLength {
		result.AddError(field, i18n.Get(i18n.KeyMaxLength, fieldName, *s.maxLength))
	}

	if s.emailRegex != nil {
		if strings.Contains(str, "..") {
			result.AddError(field, i18n.Get(i18n.KeyEmail, fieldName))
			return
		}

		parts := strings.Split(str, "@")
		if len(parts) == 2 {
			domainParts := strings.Split(parts[1], ".")
			if len(domainParts) > 0 {
				tld := domainParts[len(domainParts)-1]
				if len(tld) < 2 {
					result.AddError(field, i18n.Get(i18n.KeyEmail, fieldName))
					return
				}
			}
		}

		if !s.emailRegex.MatchString(str) {
			result.AddError(field, i18n.Get(i18n.KeyEmail, fieldName))
		}
	}

	if s.urlRegex != nil {
		if strings.Contains(str, " ") {
			result.AddError(field, i18n.Get(i18n.KeyURL, fieldName))
			return
		}

		if !strings.HasPrefix(str, "http://") && !strings.HasPrefix(str, "https://") {
			result.AddError(field, i18n.Get(i18n.KeyURL, fieldName))
			return
		}

		withoutProtocol := strings.TrimPrefix(strings.TrimPrefix(str, "https://"), "http://")
		if len(withoutProtocol) == 0 {
			result.AddError(field, i18n.Get(i18n.KeyURL, fieldName))
			return
		}

		if !s.urlRegex.MatchString(str) {
			result.AddError(field, i18n.Get(i18n.KeyURL, fieldName))
		}
	}

	if len(s.allowedValues) > 0 {
		found := false
		for _, allowed := range s.allowedValues {
			if str == allowed {
				found = true
				break
			}
		}
		if !found {
			result.AddError(field, i18n.Get(i18n.KeyOneOf, fieldName, fmt.Sprintf("%v", s.allowedValues)))
		}
	}

	if s.passwordRules != nil && str != "" {
		passwordErrors := rules.ValidatePassword(str, s.passwordRules)
		for _, err := range passwordErrors {
			result.AddError(field, fmt.Sprintf("%s %s", fieldName, err))
		}
	}
	if s.ipVersion != nil {
		if !rules.IsValidIP(str, *s.ipVersion) {
			result.AddError(field, i18n.Get(i18n.KeyIP, fieldName))
		}
	}
	if s.phoneCountry != nil {
		if !rules.IsValidPhoneNumber(str, *s.phoneCountry) {
			result.AddError(field, i18n.Get(i18n.KeyPhone, fieldName, *s.phoneCountry))
		}
	}

	// New validators
	if s.isAlpha && !alphaRegex.MatchString(str) {
		result.AddError(field, i18n.Get(i18n.KeyAlpha, fieldName))
	}

	if s.isAlphanumeric && !alphanumericRegex.MatchString(str) {
		result.AddError(field, i18n.Get(i18n.KeyAlphanumeric, fieldName))
	}

	if s.isNumeric && !numericRegex.MatchString(str) {
		result.AddError(field, i18n.Get(i18n.KeyNumericString, fieldName))
	}

	if s.startsWith != nil && !strings.HasPrefix(str, *s.startsWith) {
		result.AddError(field, i18n.Get(i18n.KeyStartsWith, fieldName, *s.startsWith))
	}

	if s.endsWith != nil && !strings.HasSuffix(str, *s.endsWith) {
		result.AddError(field, i18n.Get(i18n.KeyEndsWith, fieldName, *s.endsWith))
	}

	if s.contains != nil && !strings.Contains(str, *s.contains) {
		result.AddError(field, i18n.Get(i18n.KeyContains, fieldName, *s.contains))
	}

	if s.regexError != nil {
		result.AddError(field, fmt.Sprintf("%s: %s", fieldName, s.regexError.Error()))
		return
	}

	if s.customRegex != nil && !s.customRegex.MatchString(str) {
		result.AddError(field, i18n.Get(i18n.KeyRegex, fieldName))
	}

	if s.isMAC && !macRegex.MatchString(str) {
		result.AddError(field, i18n.Get(i18n.KeyMAC, fieldName))
	}

	if s.isHex && !hexRegex.MatchString(str) {
		result.AddError(field, i18n.Get(i18n.KeyHex, fieldName))
	}

	if s.isBase64 {
		// Check if it's valid base64 by trying to decode it
		if _, err := base64.StdEncoding.DecodeString(str); err != nil {
			result.AddError(field, i18n.Get(i18n.KeyBase64, fieldName))
		}
	}

	if s.customValidation != nil && s.customValidation.HasValidators() {
		s.customValidation.ValidateSync(field, value, result)
	}
}
