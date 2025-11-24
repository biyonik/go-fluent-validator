package core

import "fmt"

// -----------------------------------------------------------------------------
// FieldError
// -----------------------------------------------------------------------------
// Bu yapı, doğrulama süreçlerinde ortaya çıkan alan bazlı hataları temsil eden
// özel bir hata türüdür. Sistem içinde herhangi bir alan (field) ile ilgili
// bir geçersizlik durumu oluştuğunda, basit string hataları yerine daha anlamlı,
// daha yapısal ve kolay yönetilebilir bir hata modeli sunar.
//
// Bu yaklaşım, modern backend mimarilerinde olduğu gibi hataları hem geliştiriciye
// hem de üst katman servislerine okunabilir ve ayırt edilebilir bir formatta
// iletmeyi sağlar. Böylece:
//   - Hangi alanda hata olduğu net şekilde bilinir.
//   - Kullanıcıya gösterilecek mesaj ile teknik hata birbirinden ayrılabilir.
//   - Hatalar loglanırken, izleme sistemlerinde gruplanırken veya API yanıtlarında
//     işlenirken standart bir model kullanılmış olur.
//
// Laravel ve Symfony’deki form/validation error objelerinin Go’daki karşılığı
// gibi düşünülebilir. Bu sayede büyük doğrulama zincirleri çok daha temiz,
// anlaşılır ve sürdürülebilir hale gelir.
//
// Metadata:
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------
type FieldError struct {
	// Field, hatanın hangi alanda meydana geldiğini belirtir.
	Field string

	// Message, kullanıcıya gösterilecek veya loglanacak açıklayıcı hata mesajıdır.
	Message string
}

// Error
// -----------------------------------------------------------------------------
// Go'nun standart `error` arayüzünü (interface) karşılamak için kullanılan bu fonksiyon,
// FieldError'u okunabilir bir metne dönüştürür.
// Örnek çıktı: "email: Geçerli bir e-posta adresi giriniz"
//
// Bu fonksiyon sayesinde FieldError, Go ekosisteminde doğal bir hata gibi
// davranabilir ve tüm built-in hata mekanizmalarıyla uyumlu çalışır.
func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewFieldError
// -----------------------------------------------------------------------------
// Yeni bir FieldError nesnesi üretmek için kullanılan yardımcı fonksiyondur.
// Alan adı ve mesajı alır, bunları objeye dönüştürerek geliştiricinin kullanımına
// sunar.
//
// Bu fonksiyon sık kullanılan bir pattern olan “constructor helper” mantığını
// takip eder. Daha temiz, okunabilir ve standart bir hata oluşturma akışı sağlar.
func NewFieldError(field, message string) error {
	return &FieldError{Field: field, Message: message}
}
