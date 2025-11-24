package core

//
// -----------------------------------------------------------------------------
// ValidationResult
// -----------------------------------------------------------------------------
// Bu yapı, bir doğrulama sürecinin nihai çıktısını temsil eden merkezî veri modelidir.
// Doğrulama işleminden sonra elimizde iki önemli bilgi bulunur:
//
//   1. Hatalar (errors):
//      - Hangi alanlarda sorun olduğunu
//      - Hangi kuralların ihlal edildiğini
//      - Kullanıcıya hangi mesajların gösterilmesi gerektiğini içerir.
//
//   2. Geçerli Veri (validData):
//      - Doğrulamadan başarıyla geçen temizlenmiş/dönüştürülmüş verileri içerir.
//      - Böylece iş akışının sonraki aşamalarında yalnızca güvenli, doğrulanmış,
//        normlanmış veri kullanılır.
//
// ValidationResult, özellikle büyük projelerde hataların yönetimi, API çıktıları,
// form işleme süreçleri ve loglama sistemleri için temel karşılıktır.
// Modern çerçevelerdeki (Laravel ValidationResult, Symfony FormErrorBundle vb.)
// eşdeğer davranışı Go ekosisteminde minimal, okunabilir ve genişletilebilir
// bir yapıyla sunar.
//
// Metadata:
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

// ValidationResult, bir doğrulama işleminin sonucunu temsil eder.
type ValidationResult struct {
	// errors, alanlara göre gruplanmış hata mesajlarını tutar.
	errors map[string][]string

	// validData, doğrulamadan başarıyla geçen temiz veri setidir.
	validData map[string]any
}

// NewResult
// -----------------------------------------------------------------------------
// Yeni bir ValidationResult nesnesi oluşturur.
// - Başlangıçta boş hata listesi
// - Boş geçerli veri listesi
// ile temiz bir doğrulama sonucu üretmek için kullanılır.
func NewResult() *ValidationResult {
	return &ValidationResult{
		errors:    make(map[string][]string),
		validData: make(map[string]any),
	}
}

// AddError
// -----------------------------------------------------------------------------
// Belirtilen alana (field) bir hata mesajı ekler.
// Aynı alan için birden fazla hata oluşabilir; bu nedenle alan bazında liste tutulur.
func (r *ValidationResult) AddError(field, message string) {
	r.errors[field] = append(r.errors[field], message)
}

// HasErrors
// -----------------------------------------------------------------------------
// Doğrulama sonucunda herhangi bir hata olup olmadığını kontrol eder.
// - Hata varsa: true
// - Hiç hata yoksa: false
// Basit ama doğrulama akışının vazgeçilmez parçasıdır.
func (r *ValidationResult) HasErrors() bool {
	return len(r.errors) > 0
}

// Errors
// -----------------------------------------------------------------------------
// Tüm hataları olduğu gibi döndürür.
// Genellikle API yanıtlarında, debug ekranlarında veya kullanıcıya
// gösterilecek hata çıktılarında kullanılır.
func (r *ValidationResult) Errors() map[string][]string {
	return r.errors
}

// ValidData
// -----------------------------------------------------------------------------
// Geçerli (doğrulanmış ve dönüştürülmüş) veri setini döndürür.
// Doğrulama sonrasında uygulamanın iş katmanına iletilen temiz veri buradadır.
func (r *ValidationResult) ValidData() map[string]any {
	return r.validData
}

// SetValidData
// -----------------------------------------------------------------------------
// Dışarıdan doğrulanmış veri setinin tümünü bir kerede atamak için kullanılır.
// Özellikle schema bazlı doğrulamalarda, tüm alanlar tarandıktan sonra
// validData'nın topluca set edilmesi için idealdir.
func (r *ValidationResult) SetValidData(data map[string]any) {
	r.validData = data
}
