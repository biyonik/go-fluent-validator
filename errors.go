package validation

import "fmt"

//
// -----------------------------------------------------------------------------
// Validation Hataları (FieldError)
// -----------------------------------------------------------------------------
// Bu dosya, doğrulama sırasında oluşabilecek hata mesajlarını yönetmek
// için FieldError struct'ı ve yardımcı fonksiyonları içerir.
//
// Laravel/Symfony Validator'daki FieldError mantığına benzer şekilde,
// her alan için özel hata mesajları üretilebilir.
//
// Metadata:
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

// FieldError
// -----------------------------------------------------------------------------
// Belirli bir alan (field) için validation hatasını temsil eder.
//
// Alanlar:
//   - Field: Hata oluşan alan adı
//   - Message: Hata mesajı
type FieldError struct {
	Field   string
	Message string
}

// Error
// -----------------------------------------------------------------------------
// error interface implementasyonu. FieldError struct'ını error olarak döndürür.
//
// Dönüş:
//   - string: "<field>: <message>" formatında hata mesajı
func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewFieldError
// -----------------------------------------------------------------------------
// Yeni bir FieldError nesnesi oluşturur.
//
// Parametreler:
//   - field: Hata oluşan alan adı
//   - message: Hata mesajı
//
// Dönüş:
//   - error: FieldError interface'i
//
// Örnek Kullanım:
//
//	err := validation.NewFieldError("password_confirm", "Şifreler eşleşmiyor")
func NewFieldError(field, message string) error {
	return &FieldError{
		Field:   field,
		Message: message,
	}
}

// Custom validator fonksiyonlarında kullanılmak üzere basit bir validation
// error tipidir. FieldError'dan farklı olarak sadece mesaj içerir, field bilgisi
// validasyon sistemi tarafından otomatik olarak eklenir.
type ValidationError struct {
	Message string
}

// Error interface implementasyonu
func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError
// -----------------------------------------------------------------------------
// Custom validator fonksiyonlarında kullanılmak üzere yeni bir ValidationError
// oluşturur. Bu fonksiyon custom validation'larda hata döndürmek için kullanılır.
//
// Örnek kullanım:
//
//	v.String().Custom(func(value string) error {
//	    if value == "forbidden" {
//	        return core.NewValidationError("Bu değer kullanılamaz")
//	    }
//	    return nil
//	})
func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}
