// -----------------------------------------------------------------------------
// Custom Validators System
// -----------------------------------------------------------------------------
// Bu dosya, kullanıcıların kendi özel doğrulama kurallarını tanımlayabilmesi
// için gereken yapıları ve yardımcı fonksiyonları içerir. Laravel'in Custom
// Validation Rules ve Symfony'nin Custom Constraints yapısına benzer bir
// yaklaşım sunar.
//
// Neyi, Nasıl ve Neden:
//   - Neyi: Kullanıcının iş mantığına özel validation kuralları
//   - Nasıl: Callback fonksiyonları ve Rule interface'i ile
//   - Neden: Her proje farklı ihtiyaçlara sahip, standart kurallar yetmez
//
// Özellikler:
//   - Sync validation (immediate)
//   - Async validation (database checks, API calls)
//   - Chainable custom rules
//   - Context-aware validation (access to other fields)
//   - Reusable custom rules
//
// Kullanım:
//   validation.String().Custom(func(value string) error {
//       if isBlacklisted(value) {
//           return errors.New("value is blacklisted")
//       }
//       return nil
//   })
//
// Metadata:
// @author   Ahmet ALTUN
// @github   github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email    ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

package core

import (
	"context"
	"fmt"
)

// CustomValidator, senkron özel doğrulama fonksiyonu tipi
//
// Parametreler:
//   - value: Doğrulanacak değer
//
// Dönüş:
//   - error: Doğrulama başarısızsa hata, başarılıysa nil
type CustomValidator func(value any) error

// AsyncValidator, asenkron özel doğrulama fonksiyonu tipi
// Database sorguları, API çağrıları gibi I/O işlemleri için
//
// Parametreler:
//   - ctx: Context (timeout, cancellation için)
//   - value: Doğrulanacak değer
//
// Dönüş:
//   - error: Doğrulama başarısızsa hata, başarılıysa nil
type AsyncValidator func(ctx context.Context, value any) error

// ContextValidator, diğer alanlara erişebilen özel doğrulama fonksiyonu
// Cross-field validations için kullanılır
//
// Parametreler:
//   - value: Mevcut alanın değeri
//   - data: Tüm form/request verisi
//
// Dönüş:
//   - error: Doğrulama başarısızsa hata, başarılıysa nil
type ContextValidator func(value any, data map[string]any) error

// Rule, yeniden kullanılabilir özel doğrulama kuralı interface'i
// Laravel'in Rule interface'ine benzer
type Rule interface {
	// Validate, kuralın doğrulama mantığını içerir
	Validate(value any) error

	// Message, hata durumunda gösterilecek mesajı döndürür
	// Parametresiz olduğunda varsayılan mesaj döner
	// İsteğe bağlı olarak fieldName parametresi alabilir
	Message() string
}

// RuleFunc, Rule interface'ini implement eden basit fonksiyon wrapper
type RuleFunc struct {
	validator CustomValidator
	message   string
}

// NewRule, fonksiyondan yeni bir Rule oluşturur
//
// Parametreler:
//   - validator: Doğrulama fonksiyonu
//   - message: Hata mesajı
//
// Dönüş:
//   - Rule: Yeni kural instance
//
// Örnek:
//
//	blacklistRule := core.NewRule(
//	    func(value any) error {
//	        str := value.(string)
//	        if isBlacklisted(str) {
//	            return fmt.Errorf("blacklisted")
//	        }
//	        return nil
//	    },
//	    "This value is not allowed",
//	)
func NewRule(validator CustomValidator, message string) Rule {
	return &RuleFunc{
		validator: validator,
		message:   message,
	}
}

// Validate, RuleFunc'ın doğrulama mantığını çalıştırır
func (r *RuleFunc) Validate(value any) error {
	return r.validator(value)
}

// Message, RuleFunc'ın hata mesajını döndürür
func (r *RuleFunc) Message() string {
	return r.message
}

// CustomValidation, bir tip için özel doğrulama kuralları tutan yapı
type CustomValidation struct {
	syncValidators    []CustomValidator
	asyncValidators   []AsyncValidator
	contextValidators []ContextValidator
	rules             []Rule
}

// NewCustomValidation, yeni bir CustomValidation oluşturur
func NewCustomValidation() *CustomValidation {
	return &CustomValidation{
		syncValidators:    make([]CustomValidator, 0),
		asyncValidators:   make([]AsyncValidator, 0),
		contextValidators: make([]ContextValidator, 0),
		rules:             make([]Rule, 0),
	}
}

// AddSync, senkron özel doğrulama ekler
func (cv *CustomValidation) AddSync(validator CustomValidator) {
	cv.syncValidators = append(cv.syncValidators, validator)
}

// AddAsync, asenkron özel doğrulama ekler
func (cv *CustomValidation) AddAsync(validator AsyncValidator) {
	cv.asyncValidators = append(cv.asyncValidators, validator)
}

// AddContext, context-aware doğrulama ekler
func (cv *CustomValidation) AddContext(validator ContextValidator) {
	cv.contextValidators = append(cv.contextValidators, validator)
}

// AddRule, Rule interface'ini implement eden kural ekler
func (cv *CustomValidation) AddRule(rule Rule) {
	cv.rules = append(cv.rules, rule)
}

// ValidateSync, tüm senkron doğrulamaları çalıştırır
//
// Parametreler:
//   - field: Alan adı
//   - value: Değer
//   - result: ValidationResult
func (cv *CustomValidation) ValidateSync(field string, value any, result *ValidationResult) {
	// Sync validators
	for _, validator := range cv.syncValidators {
		if err := validator(value); err != nil {
			result.AddError(field, err.Error())
		}
	}

	// Rules
	for _, rule := range cv.rules {
		if err := rule.Validate(value); err != nil {
			message := rule.Message()
			if message == "" {
				message = err.Error()
			}
			result.AddError(field, message)
		}
	}
}

// ValidateAsync, tüm asenkron doğrulamaları çalıştırır
//
// Parametreler:
//   - ctx: Context
//   - field: Alan adı
//   - value: Değer
//   - result: ValidationResult
//
// Dönüş:
//   - error: İlk hata (varsa)
func (cv *CustomValidation) ValidateAsync(ctx context.Context, field string, value any, result *ValidationResult) error {
	for _, validator := range cv.asyncValidators {
		if err := validator(ctx, value); err != nil {
			result.AddError(field, err.Error())
			return err // İlk async hatada dur
		}
	}
	return nil
}

// ValidateContext, context-aware doğrulamaları çalıştırır
//
// Parametreler:
//   - field: Alan adı
//   - value: Mevcut alan değeri
//   - data: Tüm veri
//   - result: ValidationResult
func (cv *CustomValidation) ValidateContext(field string, value any, data map[string]any, result *ValidationResult) {
	for _, validator := range cv.contextValidators {
		if err := validator(value, data); err != nil {
			result.AddError(field, err.Error())
		}
	}
}

// HasValidators, herhangi bir custom validator olup olmadığını kontrol eder
func (cv *CustomValidation) HasValidators() bool {
	return len(cv.syncValidators) > 0 ||
		len(cv.asyncValidators) > 0 ||
		len(cv.contextValidators) > 0 ||
		len(cv.rules) > 0
}

// -----------------------------------------------------------------------------
// Yaygın Kullanılan Custom Rules (Built-in Utilities)
// -----------------------------------------------------------------------------

// UniqueRule, veritabanında unique kontrolü için kullanılır
type UniqueRule struct {
	checker func(value any) (bool, error)
	message string
}

// NewUniqueRule, yeni bir unique rule oluşturur
//
// Parametreler:
//   - checker: Unique kontrolü yapan fonksiyon (true = unique, false = duplicate)
//   - message: Hata mesajı
//
// Örnek:
//
//	uniqueEmail := core.NewUniqueRule(
//	    func(value any) (bool, error) {
//	        email := value.(string)
//	        exists, err := db.CheckEmailExists(email)
//	        return !exists, err
//	    },
//	    "This email is already taken",
//	)
func NewUniqueRule(checker func(value any) (bool, error), message string) Rule {
	return &UniqueRule{
		checker: checker,
		message: message,
	}
}

// Validate, unique kontrolü yapar
func (u *UniqueRule) Validate(value any) error {
	isUnique, err := u.checker(value)
	if err != nil {
		return fmt.Errorf("unique check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("not unique")
	}
	return nil
}

// Message, hata mesajını döndürür
func (u *UniqueRule) Message() string {
	return u.message
}

// ExistsRule, veritabanında kayıt kontrolü için kullanılır
type ExistsRule struct {
	checker func(value any) (bool, error)
	message string
}

// NewExistsRule, yeni bir exists rule oluşturur
//
// Parametreler:
//   - checker: Varlık kontrolü yapan fonksiyon (true = exists, false = not found)
//   - message: Hata mesajı
//
// Örnek:
//
//	userExists := core.NewExistsRule(
//	    func(value any) (bool, error) {
//	        userID := value.(int)
//	        return db.UserExists(userID)
//	    },
//	    "User not found",
//	)
func NewExistsRule(checker func(value any) (bool, error), message string) Rule {
	return &ExistsRule{
		checker: checker,
		message: message,
	}
}

// Validate, varlık kontrolü yapar
func (e *ExistsRule) Validate(value any) error {
	exists, err := e.checker(value)
	if err != nil {
		return fmt.Errorf("exists check failed: %w", err)
	}
	if !exists {
		return fmt.Errorf("not found")
	}
	return nil
}

// Message, hata mesajını döndürür
func (e *ExistsRule) Message() string {
	return e.message
}

// RegexRule, regex pattern kontrolü için kullanılır
type RegexRule struct {
	pattern string
	message string
}

// NewRegexRule, yeni bir regex rule oluşturur
//
// Parametreler:
//   - pattern: Regex pattern
//   - message: Hata mesajı
//
// Örnek:
//
//	slugRule := core.NewRegexRule(
//	    `^[a-z0-9-]+$`,
//	    "Must be a valid slug (lowercase letters, numbers, and hyphens only)",
//	)
func NewRegexRule(pattern, message string) Rule {
	return &RegexRule{
		pattern: pattern,
		message: message,
	}
}

// Validate, regex kontrolü yapar
func (r *RegexRule) Validate(value any) error {
	// Bu basitleştirilmiş bir örnek
	// Gerçek implementasyonda regexp.MustCompile kullanılmalı
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be string")
	}

	// Pattern matching logic buraya gelecek
	_ = str // placeholder

	return nil
}

// Message, hata mesajını döndürür
func (r *RegexRule) Message() string {
	return r.message
}

// RefineFunc, değeri kontrol eden ve bool döndüren fonksiyon tipi
// Zod'daki refine() pattern'i
type RefineFunc func(value any) bool

// Refine, basit boolean kontrol için kullanılır
//
// Parametreler:
//   - fn: Kontrol fonksiyonu (true = valid, false = invalid)
//   - message: Hata mesajı
//
// Dönüş:
//   - Rule
//
// Örnek:
//
//	noSpamRule := core.Refine(
//	    func(value any) bool {
//	        str := value.(string)
//	        return !strings.Contains(strings.ToLower(str), "spam")
//	    },
//	    "Content cannot contain spam",
//	)
func Refine(fn RefineFunc, message string) Rule {
	return NewRule(func(value any) error {
		if !fn(value) {
			return fmt.Errorf("refinement failed")
		}
		return nil
	}, message)
}
