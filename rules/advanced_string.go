package rules

import "regexp"

//
// -----------------------------------------------------------------------------
// Türkçe Karakter ve Alan Adı Doğrulama Kuralları
// -----------------------------------------------------------------------------
// Bu dosya, validation sisteminde kullanılabilecek iki pratik kuralı içerir:
//   - HasTurkishChars: Metin içinde Türkçe karakter bulunup bulunmadığını kontrol eder.
//   - IsValidDomain: Bir string'in geçerli bir alan adı olup olmadığını doğrular.
//
// Bu kurallar, tipler (Type) ve şemalar (Schema) içinde ek doğrulama adımları
// olarak kullanılabilir. Böylece proje genelinde ortaklaşa kullanılan, tekrar
// yazılması gerekmeyen, merkezi doğrulama fonksiyonları oluşturulmuş olur.
//
// HasTurkishChars:
//   PHP tarafındaki “hasTurkishChars” fonksiyonunun Go'ya birebir davranış
//   portudur. Go’nun rune tabanlı karakter modeli sayesinde çok baytlı Unicode
//   karakterleri güvenli ve doğru şekilde ele alır.
//
// IsValidDomain:
//   PHP’deki domain doğrulama regex’inin Go karşılığıdır. Alt alan adları
//   opsiyoneldir ve fonksiyona parametre olarak verilebilir. Böylece hem
//   “example.com” hem de “blog.example.com” gibi formatlar gerektiğinde esnek
//   şekilde doğrulanabilir.
//
// Bu tür kurallar özellikle form doğrulama, kullanıcı kayıt aşamaları,
// SEO kontrolleri, e-posta doğrulama süreçleri, input sanitizasyonu,
// API endpoint parametre kontrolleri gibi geniş bir alanda tekrar tekrar
// kullanılabilir hale gelir.
//
// Metadata:
// author: Ahmet ALTUN
// @github github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

// HasTurkishChars
// -----------------------------------------------------------------------------
// Bu fonksiyon, verilen metin içinde Türkçe karakter bulunup bulunmadığını kontrol eder.
// Karakter listesi PHP'deki karşılığından birebir alınmıştır: çÇğĞıİöÖşŞüÜ
//
// Geri dönüş:
//   - true  → Metinde en az bir Türkçe karakter var
//   - false → Hiç Türkçe karakter yok
//
// Not: regexp.MatchString kullanılarak hızlı bir tarama yapılır.
func HasTurkishChars(text string) bool {
	match, _ := regexp.MatchString(`[çÇğĞıİöÖşŞüÜ]`, text)
	return match
}

// IsValidDomain
// -----------------------------------------------------------------------------
// Bu fonksiyon, verilen domain string'inin geçerli olup olmadığını kontrol eder.
//
// Parametreler:
//   - domain: kontrol edilecek alan adı
//   - allowSubdomain: subdomain'e izin verilip verilmediğini belirler
//
// Davranış:
//   - Eğer allowSubdomain = true → blog.example.com gibi alt alan adları kabul edilir.
//   - Eğer allowSubdomain = false → yalnızca example.com gibi kök alan adları kabul edilir.
//
// Bu yapı, PHP'deki domain doğrulama regex'inin sadeleştirilmiş Go karşılığıdır.
func IsValidDomain(domain string, allowSubdomain bool) bool {
	var pattern *regexp.Regexp
	if allowSubdomain {
		pattern = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	} else {
		pattern = regexp.MustCompile(`^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`)
	}
	return pattern.MatchString(domain)
}
