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
	"fmt"
	"regexp"
	"strings"

	"github.com/biyonik/go-fluent-validator/core"
	"github.com/biyonik/go-fluent-validator/rules"
)

// StringType, string tipindeki veriler için doğrulama ve dönüşüm kurallarını tutar.
type StringType struct {
	core.BaseType
	minLength     *int
	maxLength     *int
	emailRegex    *regexp.Regexp
	urlRegex      *regexp.Regexp
	allowedValues []string
	passwordRules *rules.PasswordRules
	ipVersion     *int
	phoneCountry  *string
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
	s.emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return s
}

// URL, alanın URL formatında olmasını sağlar.
func (s *StringType) URL() *StringType {
	s.urlRegex = regexp.MustCompile(`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`)
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
		result.AddError(field, fmt.Sprintf("%s alanı metin tipinde olmalıdır", s.GetLabel(field)))
		return
	}

	fieldName := s.GetLabel(field)

	if s.minLength != nil && len(str) < *s.minLength {
		result.AddError(field, fmt.Sprintf("%s alanı en az %d karakter olmalıdır", fieldName, *s.minLength))
	}
	if s.maxLength != nil && len(str) > *s.maxLength {
		result.AddError(field, fmt.Sprintf("%s alanı en fazla %d karakter olmalıdır", fieldName, *s.maxLength))
	}
	if s.emailRegex != nil && !s.emailRegex.MatchString(str) {
		result.AddError(field, fmt.Sprintf("%s alanı geçerli bir e-posta formatında değil", fieldName))
	}
	if s.passwordRules != nil && str != "" {
		passwordErrors := rules.ValidatePassword(str, s.passwordRules)
		for _, err := range passwordErrors {
			result.AddError(field, fmt.Sprintf("%s %s", fieldName, err))
		}
	}
	if s.ipVersion != nil {
		if !rules.IsValidIP(str, *s.ipVersion) {
			result.AddError(field, fmt.Sprintf("%s alanı geçerli bir IP adresi olmalıdır", fieldName))
		}
	}
	if s.phoneCountry != nil {
		if !rules.IsValidPhoneNumber(str, *s.phoneCountry) {
			result.AddError(field, fmt.Sprintf("%s alanı geçerli bir %s telefon numarası olmalıdır", fieldName, *s.phoneCountry))
		}
	}
}
