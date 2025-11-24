// -----------------------------------------------------------------------------
// PasswordOption Yardımcı Fonksiyonları
// -----------------------------------------------------------------------------
// Bu dosya, şifre doğrulama kuralları için seçenek (option) fonksiyonlarını içerir.
// Her fonksiyon, rules.PasswordRules üzerinde belirli bir özelliği ayarlamak için
// kullanılır. Böylece kullanıcı, şifre kurallarını esnek ve okunabilir bir şekilde
// yapılandırabilir.
// -----------------------------------------------------------------------------
//
// Örnek kullanım:
//   rules := &rules.PasswordRules{}
//   WithMinLength(8)(rules)
//   WithRequireUppercase(true)(rules)
//
// Neyi, Nasıl ve Neden:
//   - Neyi: Şifre kuralları (uzunluk, büyük harf, küçük harf, sayı, özel karakter)
//   - Nasıl: Fonksiyonlar PasswordRules yapısına closures ile değer atar.
//   - Neden: Daha okunabilir, zincirlenebilir ve esnek bir kurallar sistemi sağlamak.
//
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

package types

import "github.com/biyonik/go-fluent-validator/rules"

// PasswordOption, PasswordRules üzerinde bir ayarı uygulamak için kullanılan fonksiyon tipidir.
type PasswordOption func(*rules.PasswordRules)

// WithMinLength, şifre için minimum karakter uzunluğunu ayarlar.
func WithMinLength(length int) PasswordOption {
	return func(r *rules.PasswordRules) {
		r.MinLength = length
	}
}

// WithMaxLength, şifre için maksimum karakter uzunluğunu ayarlar.
func WithMaxLength(length int) PasswordOption {
	return func(r *rules.PasswordRules) {
		r.MaxLength = length
	}
}

// WithRequireUppercase, şifrede büyük harf gerekip gerekmediğini ayarlar.
func WithRequireUppercase(required bool) PasswordOption {
	return func(r *rules.PasswordRules) {
		r.RequireUppercase = required
	}
}

// WithRequireLowercase, şifrede küçük harf gerekip gerekmediğini ayarlar.
func WithRequireLowercase(required bool) PasswordOption {
	return func(r *rules.PasswordRules) {
		r.RequireLowercase = required
	}
}

// WithRequireNumeric, şifrede sayısal karakter gerekip gerekmediğini ayarlar.
func WithRequireNumeric(required bool) PasswordOption {
	return func(r *rules.PasswordRules) {
		r.RequireNumeric = required
	}
}

// WithRequireSpecial, şifrede özel karakter gerekip gerekmediğini ayarlar.
func WithRequireSpecial(required bool) PasswordOption {
	return func(r *rules.PasswordRules) {
		r.RequireSpecial = required
	}
}

// WithSpecialChars, şifrede kullanılabilecek özel karakterleri belirler.
func WithSpecialChars(chars string) PasswordOption {
	return func(r *rules.PasswordRules) {
		r.SpecialChars = chars
	}
}

// WithMinUniqueChars, şifrede bulunması gereken minimum farklı karakter sayısını ayarlar.
func WithMinUniqueChars(count int) PasswordOption {
	return func(r *rules.PasswordRules) {
		r.MinUniqueChars = count
	}
}
