// -----------------------------------------------------------------------------
// Message System for Internationalization (i18n)
// -----------------------------------------------------------------------------
// Bu dosya, doğrulama hata mesajlarının çoklu dil desteğini sağlayan merkezi
// mesaj yönetim sistemini içerir. Sistem, Laravel'in Lang facade'ine ve
// Symfony'nin Translation component'ine benzer bir yaklaşım sunar.
//
// Neyi, Nasıl ve Neden:
//   - Neyi: Tüm validation hatalarını kullanıcının diline çevirmek
//   - Nasıl: Key-value pairs ve placeholder replacement ile
//   - Neden: Global kullanıma açık, profesyonel bir kütüphane olması için
//
// Özellikler:
//   - Çoklu dil desteği (TR, EN, varsayılan olarak)
//   - Custom mesaj tanımlama
//   - Placeholder replacement (%s, %d gibi)
//   - Thread-safe mesaj yönetimi
//   - Fallback mekanizması (dil yoksa default'a düş)
//
// Kullanım:
//   messages.SetLocale("en")
//   msg := messages.Get("validation.required", "Email")
//   // Output: "Email is required"
//
// Metadata:
// @author   Ahmet ALTUN
// @github   github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email    ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

package i18n

import (
	"fmt"
	"sync"
)

// MessageKey, mesaj anahtarı için tip tanımı
type MessageKey string

// Validation message keys
const (
	KeyRequired          MessageKey = "validation.required"
	KeyMin               MessageKey = "validation.min"
	KeyMax               MessageKey = "validation.max"
	KeyEmail             MessageKey = "validation.email"
	KeyURL               MessageKey = "validation.url"
	KeyIP                MessageKey = "validation.ip"
	KeyUUID              MessageKey = "validation.uuid"
	KeyIBAN              MessageKey = "validation.iban"
	KeyCreditCard        MessageKey = "validation.credit_card"
	KeyPhone             MessageKey = "validation.phone"
	KeyInteger           MessageKey = "validation.integer"
	KeyNumeric           MessageKey = "validation.numeric"
	KeyBoolean           MessageKey = "validation.boolean"
	KeyArray             MessageKey = "validation.array"
	KeyObject            MessageKey = "validation.object"
	KeyDate              MessageKey = "validation.date"
	KeyDateFormat        MessageKey = "validation.date_format"
	KeyDateMin           MessageKey = "validation.date_min"
	KeyDateMax           MessageKey = "validation.date_max"
	KeyOneOf             MessageKey = "validation.one_of"
	KeyMinLength         MessageKey = "validation.min_length"
	KeyMaxLength         MessageKey = "validation.max_length"
	KeyMinElements       MessageKey = "validation.min_elements"
	KeyMaxElements       MessageKey = "validation.max_elements"
	KeyPassword          MessageKey = "validation.password"
	KeyPasswordUpper     MessageKey = "validation.password.uppercase"
	KeyPasswordLower     MessageKey = "validation.password.lowercase"
	KeyPasswordNumeric   MessageKey = "validation.password.numeric"
	KeyPasswordSpecial   MessageKey = "validation.password.special"
	KeyPasswordUnique    MessageKey = "validation.password.unique_chars"
	KeyPasswordRepeating MessageKey = "validation.password.repeating"
	KeyPasswordCommon    MessageKey = "validation.password.common"
	KeyPasswordKeyboard  MessageKey = "validation.password.keyboard"
	KeyPasswordWeak      MessageKey = "validation.password.weak"
	KeyTurkishChars      MessageKey = "validation.turkish_chars"
	KeyNoTurkishChars    MessageKey = "validation.no_turkish_chars"
	KeyDomain            MessageKey = "validation.domain"
	KeyCharSet           MessageKey = "validation.charset"
	KeyTransform         MessageKey = "validation.transform_error"
	KeyCrossValidation   MessageKey = "validation.cross_validation"
)

// Messages, bir dil için tüm mesajları içeren harita
type Messages map[MessageKey]string

// Translator, mesaj çeviri sistemi
type Translator struct {
	mu              sync.RWMutex
	currentLocale   string
	defaultLocale   string
	messages        map[string]Messages
	fallbackEnabled bool
}

var (
	// globalTranslator, singleton translator instance
	globalTranslator *Translator
	once             sync.Once
)

// init, translator'ı başlatır ve varsayılan dilleri yükler
func init() {
	once.Do(func() {
		globalTranslator = &Translator{
			currentLocale:   "en",
			defaultLocale:   "en",
			messages:        make(map[string]Messages),
			fallbackEnabled: true,
		}
		globalTranslator.loadDefaultMessages()
	})
}

// loadDefaultMessages, varsayılan İngilizce ve Türkçe mesajları yükler
func (t *Translator) loadDefaultMessages() {
	// English messages
	t.messages["en"] = Messages{
		KeyRequired:          "%s is required",
		KeyMin:               "%s must be at least %v",
		KeyMax:               "%s must be at most %v",
		KeyEmail:             "%s must be a valid email address",
		KeyURL:               "%s must be a valid URL",
		KeyIP:                "%s must be a valid IP address",
		KeyUUID:              "%s must be a valid UUID",
		KeyIBAN:              "%s must be a valid IBAN",
		KeyCreditCard:        "%s must be a valid credit card number",
		KeyPhone:             "%s must be a valid %s phone number",
		KeyInteger:           "%s must be an integer",
		KeyNumeric:           "%s must be a numeric value",
		KeyBoolean:           "%s must be a boolean value",
		KeyArray:             "%s must be an array",
		KeyObject:            "%s must be an object",
		KeyDate:              "%s must be a valid date",
		KeyDateFormat:        "%s is not in a valid date format. Expected: %s",
		KeyDateMin:           "%s cannot be before %s",
		KeyDateMax:           "%s cannot be after %s",
		KeyOneOf:             "%s must be one of: %s",
		KeyMinLength:         "%s must be at least %d characters long",
		KeyMaxLength:         "%s must be at most %d characters long",
		KeyMinElements:       "%s must contain at least %d elements",
		KeyMaxElements:       "%s must contain at most %d elements",
		KeyPassword:          "%s must meet password requirements",
		KeyPasswordUpper:     "%s must contain at least one uppercase letter",
		KeyPasswordLower:     "%s must contain at least one lowercase letter",
		KeyPasswordNumeric:   "%s must contain at least one number",
		KeyPasswordSpecial:   "%s must contain at least one special character (%s)",
		KeyPasswordUnique:    "%s must contain at least %d unique characters",
		KeyPasswordRepeating: "%s cannot have more than %d repeating characters",
		KeyPasswordCommon:    "%s is too common, please choose a more secure password",
		KeyPasswordKeyboard:  "%s cannot contain keyboard sequences",
		KeyPasswordWeak:      "%s is not strong enough, please choose a more complex password",
		KeyTurkishChars:      "%s must contain Turkish characters",
		KeyNoTurkishChars:    "%s must not contain Turkish characters",
		KeyDomain:            "%s must be a valid domain name",
		KeyCharSet:           "%s must contain only '%s' characters",
		KeyTransform:         "Transformation error: %s",
		KeyCrossValidation:   "Cross-field validation failed: %s",
	}

	// Turkish messages
	t.messages["tr"] = Messages{
		KeyRequired:          "%s alanı zorunludur",
		KeyMin:               "%s alanı en az %v olmalıdır",
		KeyMax:               "%s alanı en fazla %v olmalıdır",
		KeyEmail:             "%s alanı geçerli bir e-posta adresi olmalıdır",
		KeyURL:               "%s alanı geçerli bir URL olmalıdır",
		KeyIP:                "%s alanı geçerli bir IP adresi olmalıdır",
		KeyUUID:              "%s alanı geçerli bir UUID olmalıdır",
		KeyIBAN:              "%s alanı geçerli bir IBAN olmalıdır",
		KeyCreditCard:        "%s alanı geçerli bir kredi kartı numarası olmalıdır",
		KeyPhone:             "%s alanı geçerli bir %s telefon numarası olmalıdır",
		KeyInteger:           "%s alanı tamsayı olmalıdır",
		KeyNumeric:           "%s alanı sayısal bir değer olmalıdır",
		KeyBoolean:           "%s alanı boolean tipinde olmalıdır",
		KeyArray:             "%s alanı dizi (array) tipinde olmalıdır",
		KeyObject:            "%s alanı nesne (object) tipinde olmalıdır",
		KeyDate:              "%s alanı geçerli bir tarih olmalıdır",
		KeyDateFormat:        "%s geçerli bir tarih formatı değil. Beklenen: %s",
		KeyDateMin:           "%s alanı %s tarihinden önce olamaz",
		KeyDateMax:           "%s alanı %s tarihinden sonra olamaz",
		KeyOneOf:             "%s alanı şunlardan biri olmalıdır: %s",
		KeyMinLength:         "%s alanı en az %d karakter olmalıdır",
		KeyMaxLength:         "%s alanı en fazla %d karakter olmalıdır",
		KeyMinElements:       "%s alanında en az %d eleman olmalıdır",
		KeyMaxElements:       "%s alanında en fazla %d eleman olmalıdır",
		KeyPassword:          "%s şifre gereksinimlerini karşılamalıdır",
		KeyPasswordUpper:     "%s en az bir büyük harf içermelidir",
		KeyPasswordLower:     "%s en az bir küçük harf içermelidir",
		KeyPasswordNumeric:   "%s en az bir rakam içermelidir",
		KeyPasswordSpecial:   "%s en az bir özel karakter içermelidir (%s)",
		KeyPasswordUnique:    "%s en az %d farklı karakter içermelidir",
		KeyPasswordRepeating: "%s en fazla %d adet tekrar eden karakter içerebilir",
		KeyPasswordCommon:    "%s çok yaygın bir şifre, lütfen daha güvenli bir şifre seçin",
		KeyPasswordKeyboard:  "%s klavye düzeninde sıralı karakterler içeremez",
		KeyPasswordWeak:      "%s yeterince karmaşık değil, lütfen daha güçlü bir şifre seçin",
		KeyTurkishChars:      "%s alanında Türkçe karakter bulunmalıdır",
		KeyNoTurkishChars:    "%s alanında Türkçe karakter bulunmamalıdır",
		KeyDomain:            "%s alanı geçerli bir alan adı olmalıdır",
		KeyCharSet:           "%s alanı '%s' karakter setine uymalıdır",
		KeyTransform:         "Dönüşüm hatası: %s",
		KeyCrossValidation:   "Çapraz alan doğrulaması başarısız: %s",
	}
}

// SetLocale, aktif dili ayarlar
//
// Parametreler:
//   - locale: Dil kodu (örn: "en", "tr", "de")
//
// Örnek:
//
//	i18n.SetLocale("tr")
func SetLocale(locale string) {
	globalTranslator.SetLocale(locale)
}

// SetLocale, translator için aktif dili ayarlar
func (t *Translator) SetLocale(locale string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.currentLocale = locale
}

// GetLocale, aktif dili döndürür
func GetLocale() string {
	return globalTranslator.GetLocale()
}

// GetLocale, translator'ın aktif dilini döndürür
func (t *Translator) GetLocale() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.currentLocale
}

// AddMessages, belirli bir dil için mesajları ekler veya günceller
//
// Parametreler:
//   - locale: Dil kodu
//   - messages: Mesaj haritası
//
// Örnek:
//
//	i18n.AddMessages("de", i18n.Messages{
//	    i18n.KeyRequired: "%s ist erforderlich",
//	    i18n.KeyEmail: "%s muss eine gültige E-Mail-Adresse sein",
//	})
func AddMessages(locale string, messages Messages) {
	globalTranslator.AddMessages(locale, messages)
}

// AddMessages, translator'a yeni dil mesajları ekler
func (t *Translator) AddMessages(locale string, messages Messages) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.messages[locale]; !exists {
		t.messages[locale] = make(Messages)
	}

	for key, value := range messages {
		t.messages[locale][key] = value
	}
}

// Get, belirtilen anahtar için mesajı döndürür ve args ile formatlama yapar
//
// Parametreler:
//   - key: Mesaj anahtarı
//   - args: Placeholder'lara yerleştirilecek değerler
//
// Dönüş:
//   - string: Formatlanmış mesaj
//
// Örnek:
//
//	msg := i18n.Get(i18n.KeyRequired, "Email")
//	// English: "Email is required"
//	// Turkish: "Email alanı zorunludur"
func Get(key MessageKey, args ...any) string {
	return globalTranslator.Get(key, args...)
}

// Get, mesajı döndürür ve placeholder'ları doldurur
func (t *Translator) Get(key MessageKey, args ...any) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Önce aktif dilde ara
	if messages, exists := t.messages[t.currentLocale]; exists {
		if msg, found := messages[key]; found {
			return fmt.Sprintf(msg, args...)
		}
	}

	// Fallback enabled ise default dili dene
	if t.fallbackEnabled && t.currentLocale != t.defaultLocale {
		if messages, exists := t.messages[t.defaultLocale]; exists {
			if msg, found := messages[key]; found {
				return fmt.Sprintf(msg, args...)
			}
		}
	}

	// Hiçbir şey bulunamazsa raw key döndür
	return fmt.Sprintf("[%s]", key)
}

// T, Get fonksiyonunun kısa alias'ı (Laravel'deki __() veya t() gibi)
func T(key MessageKey, args ...any) string {
	return Get(key, args...)
}

// SetFallback, fallback mekanizmasını açar/kapatır
//
// Parametreler:
//   - enabled: true ise aktif dilde bulunamazsa default dile bakılır
func SetFallback(enabled bool) {
	globalTranslator.SetFallback(enabled)
}

// SetFallback, fallback mekanizmasını yapılandırır
func (t *Translator) SetFallback(enabled bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.fallbackEnabled = enabled
}

// SetDefaultLocale, varsayılan dili ayarlar (fallback için kullanılır)
//
// Parametreler:
//   - locale: Varsayılan dil kodu
func SetDefaultLocale(locale string) {
	globalTranslator.SetDefaultLocale(locale)
}

// SetDefaultLocale, translator için varsayılan dili ayarlar
func (t *Translator) SetDefaultLocale(locale string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.defaultLocale = locale
}

// HasLocale, belirtilen dilin yüklenip yüklenmediğini kontrol eder
//
// Parametreler:
//   - locale: Kontrol edilecek dil kodu
//
// Dönüş:
//   - bool: Dil mevcutsa true
func HasLocale(locale string) bool {
	return globalTranslator.HasLocale(locale)
}

// HasLocale, dilin mevcut olup olmadığını kontrol eder
func (t *Translator) HasLocale(locale string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, exists := t.messages[locale]
	return exists
}

// GetAvailableLocales, yüklü tüm dilleri döndürür
//
// Dönüş:
//   - []string: Mevcut dil kodları listesi
func GetAvailableLocales() []string {
	return globalTranslator.GetAvailableLocales()
}

// GetAvailableLocales, tüm yüklü dilleri listeler
func (t *Translator) GetAvailableLocales() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	locales := make([]string, 0, len(t.messages))
	for locale := range t.messages {
		locales = append(locales, locale)
	}
	return locales
}
