package core

import (
	"fmt"
	"github.com/biyonik/go-fluent-validator/i18n"
)

// -----------------------------------------------------------------------------
// BaseType
// -----------------------------------------------------------------------------
// Bu yapı, sistemdeki tüm tiplerin temelini oluşturan soyut bir çekirdek sınıftır.
// Bir doğrulama ve dönüştürme altyapısının kalbinde yer alır ve her bir veri tipinin
// ortak davranışlarını standart bir çerçeveye oturtur.
//
// Buradaki amaç şudur:
//   - Her tip için “zorunlu alan”, “varsayılan değer”, “etiket”, “dönüşüm fonksiyonları”
//     gibi tekrar tekrar yazılan özellikleri merkezi bir noktada toplamak.
//   - Tiplerin üzerinde çalıştığı doğrulama ve transformasyon zincirini güvenli,
//     okunabilir ve genişletilebilir bir yapıda sunmak.
//   - Hem geliştiricinin işini kolaylaştırmak (DRY prensibi), hem de tüm sistemde
//     uyumlu bir doğrulama davranışı elde etmek.
//
// Bu yapı, daha üst düzey tipler tarafından embed edilerek kullanılır. Böylece
// tek bir yerde oluşturulan bu mimari, yüzlerce farklı tip için standart haline gelir.
// Laravel ve Symfony’nin form/validation sistemlerindeki soyut sınıfların Go karşılığı
// gibi düşünülebilir.
//
// Metadata:
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------
type BaseType struct {
	isRequired      bool
	label           string
	defaultValue    any
	transformations []func(any) (any, error)
}

// SetRequired
// -----------------------------------------------------------------------------
// Bu fonksiyon, ilgili tipin “zorunlu” bir alan olduğunu işaretler. Form girişlerinde
// veya API isteklerinde veri gelmediğinde otomatik olarak hata üretilmesini sağlar.
// Bu sayede, hem daha temiz doğrulama kuralları yazılır hem de geliştirici,
// ilgili alanı takip etmek zorunda kalmaz.
func (b *BaseType) SetRequired() {
	b.isRequired = true
}

// SetLabel
// -----------------------------------------------------------------------------
// Form alanlarına okunabilir ve kullanıcı dostu bir başlık (etiket) tanımlamak için
// kullanılır. Bu sayede hata mesajları, ham değişken adı yerine daha estetik ve
// anlaşılır bir metinle gösterilir.
func (b *BaseType) SetLabel(label string) {
	b.label = label
}

// GetLabel
// -----------------------------------------------------------------------------
// Bu fonksiyon, alan için özel olarak atanmış bir etiket varsa onu döndürür,
// yoksa varsayılan alan adını kullanır.
// `types` paketindeki diğer tiplerin insana okunabilir hata mesajlarına sahip
// olabilmesinin temel parçasıdır.
func (b *BaseType) GetLabel(defaultLabel string) string {
	if b.label != "" {
		return b.label
	}
	return defaultLabel
}

// SetDefault
// -----------------------------------------------------------------------------
// Bu fonksiyon, bir alan için varsayılan değer atanmasını sağlar. Kullanıcı herhangi
// bir değer göndermezse veya nil gelirse, bu varsayılan değer otomatik olarak kullanılır.
// Özellikle optional alanların yönetimi için son derece kritiktir.
func (b *BaseType) SetDefault(value any) {
	b.defaultValue = value
}

// AddTransform
// -----------------------------------------------------------------------------
// Bir alanı kaydetmeden veya doğrulamadan önce üzerinde çalışılacak dönüşüm adımlarını
// tanımlar.
// Örnek: trim, lower-case, sayıya dönüştürme, format temizleme vb.
// Her dönüşüm bir fonksiyon olarak eklenir ve Transform sırasında sırayla uygulanır.
func (b *BaseType) AddTransform(fn func(any) (any, error)) {
	b.transformations = append(b.transformations, fn)
}

// Transform
// -----------------------------------------------------------------------------
// Bu fonksiyon, ilgili alana eklenmiş tüm dönüşüm fonksiyonlarını sırayla çalıştırır.
//   - Eğer değer boş (nil) ama varsayılan değer tanımlıysa, otomatik olarak varsayılan
//     değer uygulanır.
//   - Dönüşüm zinciri boyunca herhangi bir adım hata verirse işlem kesilir.
//
// Bu yapı sayesinde, her tipin kendi dönüşüm akışı sade ve okunabilir bir şekilde
// tanımlanabilir.
func (b *BaseType) Transform(value any) (any, error) {
	if value == nil && b.defaultValue != nil {
		value = b.defaultValue
	}
	if value == nil {
		return nil, nil
	}
	var err error
	for _, fn := range b.transformations {
		value, err = fn(value)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

// Validate
// -----------------------------------------------------------------------------
// Bu fonksiyon, ilgili alan için temel doğrulamayı yapar. Şimdilik sadece
// “zorunlu alan” kontrolü içerir; ancak BaseType diğer tipler tarafından
// genişletildiği için her tipin kendine ait gelişmiş doğrulama kuralları
// bu yapının üzerine inşa edilir.
//
// Eğer alan zorunlu ise:
// - Nil değer,
// - Boş string değer,
// hata olarak işlenir.
//
// Hata mesajları, geliştirici için değil kullanıcı için okunabilir şekilde
// `label` üzerinden üretilir.
func (b *BaseType) Validate(field string, value any, result *ValidationResult) {
	fieldName := b.GetLabel(field)
	if b.isRequired {
		if value == nil {
			result.AddError(field, i18n.Get(i18n.KeyRequired, fieldName))
			return
		}
		if str, ok := value.(string); ok && str == "" {
			result.AddError(field, i18n.Get(i18n.KeyRequired, fieldName))
			return
		}
	}
}
