package core

//
// -----------------------------------------------------------------------------
// Type & Schema Arayüzleri
// -----------------------------------------------------------------------------
// Bu dosya, doğrulama altyapısının iki temel yapı taşı olan `Type` ve `Schema`
// arayüzlerini tanımlar. Bir veri doğrulama sisteminin nasıl davranması gerektiğini
// belirleyen çatıyı oluşturur.
//
// Bu mimari Laravel'in Validator yapısını, Symfony’nin Form/Constraint sistemini
// anımsatacak şekilde tasarlanmıştır. Amaç; hem sade hem de genişletilebilir,
// modüler, okunabilir ve geliştirici dostu bir doğrulama ekosistemi sunmaktır.
//
// Type Arayüzü:
//   - Her bir veri tipinin (string, int, email, credit card vb.) nasıl doğrulanacağını
//     ve doğrulama öncesi hangi dönüşümlerden geçeceğini tanımlar.
//   - Uygulamada kullanılan bütün tipler bu arayüzü implement ederek ortak
//     davranış gösterir.
//
// Schema Arayüzü:
//   - Bir veri setinin tüm alanlarını aynı anda ele alan, "form şeması" gibi
//     çalışan kapsayıcı doğrulama mekanizmasıdır.
//   - Tiplerin birlikte çalışmasını sağlar, çapraz doğrulamalar (cross validation),
//     koşullu kurallar (when) gibi gelişmiş özellikler sunar.
//   - Böylece kompleks doğrulama senaryoları bile sade ve okunabilir kalır.
//
// Bu tasarım ile elde edilen faydalar:
//   - Projede yüzlerce farklı alanda doğrulama yapıldığında bile aynı standart korunur.
//   - Alan bazlı doğrulama (Type) ile form bazlı doğrulama (Schema) birbirinden bağımsız,
//     ama koordineli şekilde çalışır.
//   - Kod tekrarı azalır, genişletme kolaylaşır.
//
// Metadata:
// @author Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

// Type her veri tipi doğrulayıcısının uyması gereken arayüzdür.
type Type interface {
	// Validate, asıl doğrulama mantığını çalıştırır.
	// Her tip kendi doğrulama kurallarını burada uygular.
	Validate(field string, value any, result *ValidationResult)

	// Transform, doğrulama öncesinde veriyi temizler ve dönüştürür.
	// Trim, parse, normalizasyon gibi işlemler burada yapılır.
	Transform(value any) (any, error)
}

// Schema veri setini doğrulayan yapıdır.
type Schema interface {
	// Validate, verilen data haritası üzerinden tüm doğrulamayı çalıştırır
	// ve tek bir ValidationResult döner.
	Validate(data map[string]any) *ValidationResult

	// Shape, doğrulanacak veri yapısının hangi alanlardan oluştuğunu tanımlar.
	// Her alan bir Type örneği ile ilişkilendirilir.
	Shape(shape map[string]Type) Schema

	// CrossValidate, tüm veri seti üzerinde çalışan yüksek seviyeli doğrulama kuralları sağlar.
	// Örneğin: "start_date < end_date" gibi ilişkisel kontroller.
	CrossValidate(fn func(data map[string]any) error) Schema

	// When, belirli bir alan belirli bir değere sahipse ek kurallar eklemek için kullanılır.
	// Laravel'in "sometimes" veya "required_if" kurallarına benzer bir mantık sunar.
	When(field string, expectedValue any, callback func() Schema) Schema
}
