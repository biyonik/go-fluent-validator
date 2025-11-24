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
	KeyString            MessageKey = "validation.string"
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
	// New string validators
	KeyAlpha          MessageKey = "validation.alpha"
	KeyAlphanumeric   MessageKey = "validation.alphanumeric"
	KeyNumericString  MessageKey = "validation.numeric_string"
	KeyStartsWith     MessageKey = "validation.starts_with"
	KeyEndsWith       MessageKey = "validation.ends_with"
	KeyContains       MessageKey = "validation.contains"
	KeyRegex          MessageKey = "validation.regex"
	KeyMAC            MessageKey = "validation.mac"
	KeyHex            MessageKey = "validation.hex"
	KeyBase64         MessageKey = "validation.base64"
	// New number validators
	KeyPositive   MessageKey = "validation.positive"
	KeyNegative   MessageKey = "validation.negative"
	KeyMultipleOf MessageKey = "validation.multiple_of"
	KeyBetween    MessageKey = "validation.between"
	// New array validators
	KeyUnique         MessageKey = "validation.unique"
	KeyArrayContains  MessageKey = "validation.array_contains"
	KeyNotEmpty       MessageKey = "validation.not_empty"
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
		KeyString:            "%s must be a string",
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
		// New string validators
		KeyAlpha:         "%s must contain only alphabetic characters",
		KeyAlphanumeric:  "%s must contain only alphanumeric characters",
		KeyNumericString: "%s must contain only numeric characters",
		KeyStartsWith:    "%s must start with '%s'",
		KeyEndsWith:      "%s must end with '%s'",
		KeyContains:      "%s must contain '%s'",
		KeyRegex:         "%s does not match the required pattern",
		KeyMAC:           "%s must be a valid MAC address",
		KeyHex:           "%s must be a valid hexadecimal string",
		KeyBase64:        "%s must be a valid base64 string",
		// New number validators
		KeyPositive:   "%s must be a positive number",
		KeyNegative:   "%s must be a negative number",
		KeyMultipleOf: "%s must be a multiple of %v",
		KeyBetween:    "%s must be between %v and %v",
		// New array validators
		KeyUnique:        "%s must contain only unique elements",
		KeyArrayContains: "%s must contain the value '%v'",
		KeyNotEmpty:      "%s must not be empty",
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
		KeyString:            "%s alanı metin tipinde olmalıdır",
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
		// New string validators
		KeyAlpha:         "%s alanı sadece alfabetik karakterler içermelidir",
		KeyAlphanumeric:  "%s alanı sadece alfanümerik karakterler içermelidir",
		KeyNumericString: "%s alanı sadece sayısal karakterler içermelidir",
		KeyStartsWith:    "%s alanı '%s' ile başlamalıdır",
		KeyEndsWith:      "%s alanı '%s' ile bitmelidir",
		KeyContains:      "%s alanı '%s' içermelidir",
		KeyRegex:         "%s alanı gerekli desene uymuyor",
		KeyMAC:           "%s alanı geçerli bir MAC adresi olmalıdır",
		KeyHex:           "%s alanı geçerli bir onaltılık (hexadecimal) dize olmalıdır",
		KeyBase64:        "%s alanı geçerli bir base64 dizesi olmalıdır",
		// New number validators
		KeyPositive:   "%s alanı pozitif bir sayı olmalıdır",
		KeyNegative:   "%s alanı negatif bir sayı olmalıdır",
		KeyMultipleOf: "%s alanı %v'nin katı olmalıdır",
		KeyBetween:    "%s alanı %v ile %v arasında olmalıdır",
		// New array validators
		KeyUnique:        "%s alanı sadece benzersiz elemanlar içermelidir",
		KeyArrayContains: "%s alanı '%v' değerini içermelidir",
		KeyNotEmpty:      "%s alanı boş olmamalıdır",
	}

	// German messages
	t.messages["de"] = Messages{
		KeyRequired:          "%s ist erforderlich",
		KeyMin:               "%s muss mindestens %v sein",
		KeyMax:               "%s darf höchstens %v sein",
		KeyEmail:             "%s muss eine gültige E-Mail-Adresse sein",
		KeyURL:               "%s muss eine gültige URL sein",
		KeyIP:                "%s muss eine gültige IP-Adresse sein",
		KeyUUID:              "%s muss eine gültige UUID sein",
		KeyIBAN:              "%s muss eine gültige IBAN sein",
		KeyCreditCard:        "%s muss eine gültige Kreditkartennummer sein",
		KeyPhone:             "%s muss eine gültige %s Telefonnummer sein",
		KeyInteger:           "%s muss eine ganze Zahl sein",
		KeyNumeric:           "%s muss ein numerischer Wert sein",
		KeyBoolean:           "%s muss ein boolescher Wert sein",
		KeyArray:             "%s muss ein Array sein",
		KeyObject:            "%s muss ein Objekt sein",
		KeyDate:              "%s muss ein gültiges Datum sein",
		KeyDateFormat:        "%s hat kein gültiges Datumsformat. Erwartet: %s",
		KeyDateMin:           "%s darf nicht vor %s liegen",
		KeyDateMax:           "%s darf nicht nach %s liegen",
		KeyOneOf:             "%s muss einer der folgenden Werte sein: %s",
		KeyMinLength:         "%s muss mindestens %d Zeichen lang sein",
		KeyMaxLength:         "%s darf höchstens %d Zeichen lang sein",
		KeyMinElements:       "%s muss mindestens %d Elemente enthalten",
		KeyMaxElements:       "%s darf höchstens %d Elemente enthalten",
		KeyPassword:          "%s muss die Passwortanforderungen erfüllen",
		KeyPasswordUpper:     "%s muss mindestens einen Großbuchstaben enthalten",
		KeyPasswordLower:     "%s muss mindestens einen Kleinbuchstaben enthalten",
		KeyPasswordNumeric:   "%s muss mindestens eine Zahl enthalten",
		KeyPasswordSpecial:   "%s muss mindestens ein Sonderzeichen enthalten (%s)",
		KeyPasswordUnique:    "%s muss mindestens %d eindeutige Zeichen enthalten",
		KeyPasswordRepeating: "%s darf nicht mehr als %d aufeinanderfolgende Zeichen enthalten",
		KeyPasswordCommon:    "%s ist zu gebräuchlich, bitte wählen Sie ein sichereres Passwort",
		KeyPasswordKeyboard:  "%s darf keine Tastatursequenzen enthalten",
		KeyPasswordWeak:      "%s ist nicht stark genug, bitte wählen Sie ein komplexeres Passwort",
		KeyTurkishChars:      "%s muss türkische Zeichen enthalten",
		KeyNoTurkishChars:    "%s darf keine türkischen Zeichen enthalten",
		KeyDomain:            "%s muss ein gültiger Domainname sein",
		KeyCharSet:           "%s darf nur '%s' Zeichen enthalten",
		KeyTransform:         "Transformationsfehler: %s",
		KeyCrossValidation:   "Feldübergreifende Validierung fehlgeschlagen: %s",
		// New string validators
		KeyAlpha:         "%s darf nur alphabetische Zeichen enthalten",
		KeyAlphanumeric:  "%s darf nur alphanumerische Zeichen enthalten",
		KeyNumericString: "%s darf nur numerische Zeichen enthalten",
		KeyStartsWith:    "%s muss mit '%s' beginnen",
		KeyEndsWith:      "%s muss mit '%s' enden",
		KeyContains:      "%s muss '%s' enthalten",
		KeyRegex:         "%s entspricht nicht dem erforderlichen Muster",
		KeyMAC:           "%s muss eine gültige MAC-Adresse sein",
		KeyHex:           "%s muss eine gültige hexadezimale Zeichenfolge sein",
		KeyBase64:        "%s muss eine gültige Base64-Zeichenfolge sein",
		// New number validators
		KeyPositive:   "%s muss eine positive Zahl sein",
		KeyNegative:   "%s muss eine negative Zahl sein",
		KeyMultipleOf: "%s muss ein Vielfaches von %v sein",
		KeyBetween:    "%s muss zwischen %v und %v liegen",
		// New array validators
		KeyUnique:        "%s darf nur eindeutige Elemente enthalten",
		KeyArrayContains: "%s muss den Wert '%v' enthalten",
		KeyNotEmpty:      "%s darf nicht leer sein",
	}

	// French messages
	t.messages["fr"] = Messages{
		KeyRequired:          "%s est requis",
		KeyMin:               "%s doit être au moins %v",
		KeyMax:               "%s doit être au maximum %v",
		KeyEmail:             "%s doit être une adresse e-mail valide",
		KeyURL:               "%s doit être une URL valide",
		KeyIP:                "%s doit être une adresse IP valide",
		KeyUUID:              "%s doit être un UUID valide",
		KeyIBAN:              "%s doit être un IBAN valide",
		KeyCreditCard:        "%s doit être un numéro de carte de crédit valide",
		KeyPhone:             "%s doit être un numéro de téléphone %s valide",
		KeyInteger:           "%s doit être un nombre entier",
		KeyNumeric:           "%s doit être une valeur numérique",
		KeyBoolean:           "%s doit être une valeur booléenne",
		KeyArray:             "%s doit être un tableau",
		KeyObject:            "%s doit être un objet",
		KeyDate:              "%s doit être une date valide",
		KeyDateFormat:        "%s n'est pas dans un format de date valide. Attendu : %s",
		KeyDateMin:           "%s ne peut pas être antérieur à %s",
		KeyDateMax:           "%s ne peut pas être postérieur à %s",
		KeyOneOf:             "%s doit être l'un des suivants : %s",
		KeyMinLength:         "%s doit contenir au moins %d caractères",
		KeyMaxLength:         "%s doit contenir au maximum %d caractères",
		KeyMinElements:       "%s doit contenir au moins %d éléments",
		KeyMaxElements:       "%s doit contenir au maximum %d éléments",
		KeyPassword:          "%s doit répondre aux exigences de mot de passe",
		KeyPasswordUpper:     "%s doit contenir au moins une lettre majuscule",
		KeyPasswordLower:     "%s doit contenir au moins une lettre minuscule",
		KeyPasswordNumeric:   "%s doit contenir au moins un chiffre",
		KeyPasswordSpecial:   "%s doit contenir au moins un caractère spécial (%s)",
		KeyPasswordUnique:    "%s doit contenir au moins %d caractères uniques",
		KeyPasswordRepeating: "%s ne peut pas avoir plus de %d caractères répétitifs",
		KeyPasswordCommon:    "%s est trop courant, veuillez choisir un mot de passe plus sécurisé",
		KeyPasswordKeyboard:  "%s ne peut pas contenir de séquences de clavier",
		KeyPasswordWeak:      "%s n'est pas assez fort, veuillez choisir un mot de passe plus complexe",
		KeyTurkishChars:      "%s doit contenir des caractères turcs",
		KeyNoTurkishChars:    "%s ne doit pas contenir de caractères turcs",
		KeyDomain:            "%s doit être un nom de domaine valide",
		KeyCharSet:           "%s doit contenir uniquement des caractères '%s'",
		KeyTransform:         "Erreur de transformation : %s",
		KeyCrossValidation:   "La validation entre champs a échoué : %s",
		// New string validators
		KeyAlpha:         "%s ne doit contenir que des caractères alphabétiques",
		KeyAlphanumeric:  "%s ne doit contenir que des caractères alphanumériques",
		KeyNumericString: "%s ne doit contenir que des caractères numériques",
		KeyStartsWith:    "%s doit commencer par '%s'",
		KeyEndsWith:      "%s doit se terminer par '%s'",
		KeyContains:      "%s doit contenir '%s'",
		KeyRegex:         "%s ne correspond pas au motif requis",
		KeyMAC:           "%s doit être une adresse MAC valide",
		KeyHex:           "%s doit être une chaîne hexadécimale valide",
		KeyBase64:        "%s doit être une chaîne base64 valide",
		// New number validators
		KeyPositive:   "%s doit être un nombre positif",
		KeyNegative:   "%s doit être un nombre négatif",
		KeyMultipleOf: "%s doit être un multiple de %v",
		KeyBetween:    "%s doit être entre %v et %v",
		// New array validators
		KeyUnique:        "%s ne doit contenir que des éléments uniques",
		KeyArrayContains: "%s doit contenir la valeur '%v'",
		KeyNotEmpty:      "%s ne doit pas être vide",
	}

	// Spanish messages
	t.messages["es"] = Messages{
		KeyRequired:          "%s es requerido",
		KeyMin:               "%s debe ser al menos %v",
		KeyMax:               "%s debe ser como máximo %v",
		KeyEmail:             "%s debe ser una dirección de correo electrónico válida",
		KeyURL:               "%s debe ser una URL válida",
		KeyIP:                "%s debe ser una dirección IP válida",
		KeyUUID:              "%s debe ser un UUID válido",
		KeyIBAN:              "%s debe ser un IBAN válido",
		KeyCreditCard:        "%s debe ser un número de tarjeta de crédito válido",
		KeyPhone:             "%s debe ser un número de teléfono %s válido",
		KeyInteger:           "%s debe ser un número entero",
		KeyNumeric:           "%s debe ser un valor numérico",
		KeyBoolean:           "%s debe ser un valor booleano",
		KeyArray:             "%s debe ser un array",
		KeyObject:            "%s debe ser un objeto",
		KeyDate:              "%s debe ser una fecha válida",
		KeyDateFormat:        "%s no está en un formato de fecha válido. Esperado: %s",
		KeyDateMin:           "%s no puede ser anterior a %s",
		KeyDateMax:           "%s no puede ser posterior a %s",
		KeyOneOf:             "%s debe ser uno de los siguientes: %s",
		KeyMinLength:         "%s debe tener al menos %d caracteres",
		KeyMaxLength:         "%s debe tener como máximo %d caracteres",
		KeyMinElements:       "%s debe contener al menos %d elementos",
		KeyMaxElements:       "%s debe contener como máximo %d elementos",
		KeyPassword:          "%s debe cumplir con los requisitos de contraseña",
		KeyPasswordUpper:     "%s debe contener al menos una letra mayúscula",
		KeyPasswordLower:     "%s debe contener al menos una letra minúscula",
		KeyPasswordNumeric:   "%s debe contener al menos un número",
		KeyPasswordSpecial:   "%s debe contener al menos un carácter especial (%s)",
		KeyPasswordUnique:    "%s debe contener al menos %d caracteres únicos",
		KeyPasswordRepeating: "%s no puede tener más de %d caracteres repetitivos",
		KeyPasswordCommon:    "%s es demasiado común, por favor elija una contraseña más segura",
		KeyPasswordKeyboard:  "%s no puede contener secuencias de teclado",
		KeyPasswordWeak:      "%s no es lo suficientemente fuerte, por favor elija una contraseña más compleja",
		KeyTurkishChars:      "%s debe contener caracteres turcos",
		KeyNoTurkishChars:    "%s no debe contener caracteres turcos",
		KeyDomain:            "%s debe ser un nombre de dominio válido",
		KeyCharSet:           "%s debe contener solo caracteres '%s'",
		KeyTransform:         "Error de transformación: %s",
		KeyCrossValidation:   "La validación entre campos falló: %s",
		// New string validators
		KeyAlpha:         "%s debe contener solo caracteres alfabéticos",
		KeyAlphanumeric:  "%s debe contener solo caracteres alfanuméricos",
		KeyNumericString: "%s debe contener solo caracteres numéricos",
		KeyStartsWith:    "%s debe comenzar con '%s'",
		KeyEndsWith:      "%s debe terminar con '%s'",
		KeyContains:      "%s debe contener '%s'",
		KeyRegex:         "%s no coincide con el patrón requerido",
		KeyMAC:           "%s debe ser una dirección MAC válida",
		KeyHex:           "%s debe ser una cadena hexadecimal válida",
		KeyBase64:        "%s debe ser una cadena base64 válida",
		// New number validators
		KeyPositive:   "%s debe ser un número positivo",
		KeyNegative:   "%s debe ser un número negativo",
		KeyMultipleOf: "%s debe ser un múltiplo de %v",
		KeyBetween:    "%s debe estar entre %v y %v",
		// New array validators
		KeyUnique:        "%s debe contener solo elementos únicos",
		KeyArrayContains: "%s debe contener el valor '%v'",
		KeyNotEmpty:      "%s no debe estar vacío",
	}

	// Japanese messages
	t.messages["ja"] = Messages{
		KeyRequired:          "%sは必須です",
		KeyMin:               "%sは最小%v以上である必要があります",
		KeyMax:               "%sは最大%v以下である必要があります",
		KeyEmail:             "%sは有効なメールアドレスである必要があります",
		KeyURL:               "%sは有効なURLである必要があります",
		KeyIP:                "%sは有効なIPアドレスである必要があります",
		KeyUUID:              "%sは有効なUUIDである必要があります",
		KeyIBAN:              "%sは有効なIBANである必要があります",
		KeyCreditCard:        "%sは有効なクレジットカード番号である必要があります",
		KeyPhone:             "%sは有効な%s電話番号である必要があります",
		KeyInteger:           "%sは整数である必要があります",
		KeyNumeric:           "%sは数値である必要があります",
		KeyBoolean:           "%sはブール値である必要があります",
		KeyArray:             "%sは配列である必要があります",
		KeyObject:            "%sはオブジェクトである必要があります",
		KeyDate:              "%sは有効な日付である必要があります",
		KeyDateFormat:        "%sは有効な日付形式ではありません。期待される形式: %s",
		KeyDateMin:           "%sは%sより前にすることはできません",
		KeyDateMax:           "%sは%sより後にすることはできません",
		KeyOneOf:             "%sは次のいずれかである必要があります: %s",
		KeyMinLength:         "%sは最低%d文字である必要があります",
		KeyMaxLength:         "%sは最大%d文字である必要があります",
		KeyMinElements:       "%sは最低%d個の要素を含む必要があります",
		KeyMaxElements:       "%sは最大%d個の要素を含む必要があります",
		KeyPassword:          "%sはパスワード要件を満たす必要があります",
		KeyPasswordUpper:     "%sは少なくとも1つの大文字を含む必要があります",
		KeyPasswordLower:     "%sは少なくとも1つの小文字を含む必要があります",
		KeyPasswordNumeric:   "%sは少なくとも1つの数字を含む必要があります",
		KeyPasswordSpecial:   "%sは少なくとも1つの特殊文字を含む必要があります (%s)",
		KeyPasswordUnique:    "%sは少なくとも%d個のユニークな文字を含む必要があります",
		KeyPasswordRepeating: "%sは%d個以上の繰り返し文字を含むことはできません",
		KeyPasswordCommon:    "%sは一般的すぎます。より安全なパスワードを選択してください",
		KeyPasswordKeyboard:  "%sはキーボード配列を含むことはできません",
		KeyPasswordWeak:      "%sは十分に強力ではありません。より複雑なパスワードを選択してください",
		KeyTurkishChars:      "%sはトルコ語の文字を含む必要があります",
		KeyNoTurkishChars:    "%sはトルコ語の文字を含んではいけません",
		KeyDomain:            "%sは有効なドメイン名である必要があります",
		KeyCharSet:           "%sは'%s'文字のみを含む必要があります",
		KeyTransform:         "変換エラー: %s",
		KeyCrossValidation:   "フィールド間の検証に失敗しました: %s",
		// New string validators
		KeyAlpha:         "%sはアルファベット文字のみを含む必要があります",
		KeyAlphanumeric:  "%sは英数字のみを含む必要があります",
		KeyNumericString: "%sは数字のみを含む必要があります",
		KeyStartsWith:    "%sは'%s'で始まる必要があります",
		KeyEndsWith:      "%sは'%s'で終わる必要があります",
		KeyContains:      "%sは'%s'を含む必要があります",
		KeyRegex:         "%sは必要なパターンと一致しません",
		KeyMAC:           "%sは有効なMACアドレスである必要があります",
		KeyHex:           "%sは有効な16進数文字列である必要があります",
		KeyBase64:        "%sは有効なbase64文字列である必要があります",
		// New number validators
		KeyPositive:   "%sは正の数である必要があります",
		KeyNegative:   "%sは負の数である必要があります",
		KeyMultipleOf: "%sは%vの倍数である必要があります",
		KeyBetween:    "%sは%vと%vの間である必要があります",
		// New array validators
		KeyUnique:        "%sは一意の要素のみを含む必要があります",
		KeyArrayContains: "%sは値'%v'を含む必要があります",
		KeyNotEmpty:      "%sは空であってはいけません",
	}

	// Chinese (Simplified) messages
	t.messages["zh"] = Messages{
		KeyRequired:          "%s是必填项",
		KeyMin:               "%s必须至少为%v",
		KeyMax:               "%s最多为%v",
		KeyEmail:             "%s必须是有效的电子邮件地址",
		KeyURL:               "%s必须是有效的URL",
		KeyIP:                "%s必须是有效的IP地址",
		KeyUUID:              "%s必须是有效的UUID",
		KeyIBAN:              "%s必须是有效的IBAN",
		KeyCreditCard:        "%s必须是有效的信用卡号",
		KeyPhone:             "%s必须是有效的%s电话号码",
		KeyInteger:           "%s必须是整数",
		KeyNumeric:           "%s必须是数值",
		KeyBoolean:           "%s必须是布尔值",
		KeyArray:             "%s必须是数组",
		KeyObject:            "%s必须是对象",
		KeyDate:              "%s必须是有效的日期",
		KeyDateFormat:        "%s不是有效的日期格式。期望格式：%s",
		KeyDateMin:           "%s不能早于%s",
		KeyDateMax:           "%s不能晚于%s",
		KeyOneOf:             "%s必须是以下之一：%s",
		KeyMinLength:         "%s必须至少为%d个字符",
		KeyMaxLength:         "%s最多为%d个字符",
		KeyMinElements:       "%s必须包含至少%d个元素",
		KeyMaxElements:       "%s最多包含%d个元素",
		KeyPassword:          "%s必须符合密码要求",
		KeyPasswordUpper:     "%s必须包含至少一个大写字母",
		KeyPasswordLower:     "%s必须包含至少一个小写字母",
		KeyPasswordNumeric:   "%s必须包含至少一个数字",
		KeyPasswordSpecial:   "%s必须包含至少一个特殊字符 (%s)",
		KeyPasswordUnique:    "%s必须包含至少%d个不同的字符",
		KeyPasswordRepeating: "%s不能包含超过%d个重复字符",
		KeyPasswordCommon:    "%s太常见了，请选择更安全的密码",
		KeyPasswordKeyboard:  "%s不能包含键盘序列",
		KeyPasswordWeak:      "%s不够强，请选择更复杂的密码",
		KeyTurkishChars:      "%s必须包含土耳其字符",
		KeyNoTurkishChars:    "%s不能包含土耳其字符",
		KeyDomain:            "%s必须是有效的域名",
		KeyCharSet:           "%s只能包含'%s'字符",
		KeyTransform:         "转换错误：%s",
		KeyCrossValidation:   "字段间验证失败：%s",
		// New string validators
		KeyAlpha:         "%s必须只包含字母字符",
		KeyAlphanumeric:  "%s必须只包含字母数字字符",
		KeyNumericString: "%s必须只包含数字字符",
		KeyStartsWith:    "%s必须以'%s'开头",
		KeyEndsWith:      "%s必须以'%s'结尾",
		KeyContains:      "%s必须包含'%s'",
		KeyRegex:         "%s不匹配所需模式",
		KeyMAC:           "%s必须是有效的MAC地址",
		KeyHex:           "%s必须是有效的十六进制字符串",
		KeyBase64:        "%s必须是有效的base64字符串",
		// New number validators
		KeyPositive:   "%s必须是正数",
		KeyNegative:   "%s必须是负数",
		KeyMultipleOf: "%s必须是%v的倍数",
		KeyBetween:    "%s必须在%v和%v之间",
		// New array validators
		KeyUnique:        "%s必须只包含唯一元素",
		KeyArrayContains: "%s必须包含值'%v'",
		KeyNotEmpty:      "%s不能为空",
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
