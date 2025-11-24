// Package rules, network doğrulama kurallarını içerir.
// Bu dosya IP ve telefon numarası doğrulaması gibi network odaklı kuralları barındırır.
package rules

import (
	"net"    // IP doğrulaması için standart kütüphane
	"regexp" // Regex işlemleri için
)

//
// -----------------------------------------------------------------------------
// Network Doğrulama Kuralları
// -----------------------------------------------------------------------------
// Bu dosya, modern uygulamalarda sıkça ihtiyaç duyulan IP adresi ve telefon numarası
// doğrulama işlemlerini merkezi bir yerde toplar. Mikroservis mimarisine sahip
// yapılarda, API güvenliğinde, form doğrulamaları ve kayıt süreçlerinde kritik
// önem taşıyan bu kontroller; tekrar yazılmayı engellemek, standart oluşturmak
// ve genişletilebilir bir doğrulama zinciri sağlamak amacıyla soyutlanmıştır.
//
// IsValidIP:
//   - Kullanıcının girdiği IP adresinin geçerli bir IPv4 veya IPv6 formatında
//     olup olmadığını kontrol eder.
//   - Net kütüphanesinin düşük seviyeli, güvenilir parse mekanizmasını kullanır.
//
// IsValidPhoneNumber:
//   - Ülke bazlı telefon numarası doğrulama yapar.
//   - Türkiye ve ABD için hazır regex şablonları içerir.
//   - Geliştirici isterse aynı map’e yeni ülke kuralları ekleyerek sistemi genişletebilir.
//
// Bu mimari, Laravel'in rule sınıflarını andırır; yalın ama güçlü bir doğrulama
// altyapısı sunar.
//
// Metadata:
// @author    Ahmet Altun
// @email     ahmet.altun60@gmail.com
// @github    github.com/biyonik
// @linkedin  linkedin.com/in/biyonik
// -----------------------------------------------------------------------------

// IsValidIP
// -----------------------------------------------------------------------------
// Verilen IP adresinin geçerli olup olmadığını kontrol eder.
//
// Parametreler:
//   - ip: Doğrulanacak IP adresi (string).
//   - version: IP versiyonu.
//     4 = IPv4
//     6 = IPv6
//     0 = Her ikisi de kabul edilir
//
// Dönüş:
//   - bool → IP geçerliyse true, değilse false.
//
// Açıklama:
// net.ParseIP, Go'nun yerleşik IP parsing mekanizmasıdır ve RFC uyumlu,
// güvenilir bir kontrol sağlar. Version seçimi ile daha sıkı doğrulama
// yapılabilir.
func IsValidIP(ip string, version int) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	switch version {
	case 4:
		return parsedIP.To4() != nil
	case 6:
		return parsedIP.To4() == nil && parsedIP.To16() != nil
	case 0:
		return true
	default:
		return false
	}
}

// phonePatterns
// -----------------------------------------------------------------------------
// Ülke bazlı telefon numarası doğrulama regex kalıplarını tutan harita.
// Bu yapı isteğe bağlı olarak genişletilebilir.
//
// TR → Türkiye GSM numaraları için
// US → Amerika birleşik devletleri telefon formatı için
var phonePatterns = map[string]*regexp.Regexp{
	"TR": regexp.MustCompile(`^(05|5)[0-9]{9}$`),                    // Türkiye GSM
	"US": regexp.MustCompile(`^(\+1|1)?[2-9]\d{2}[2-9]\d{2}\d{4}$`), // ABD
}

// IsValidPhoneNumber
// -----------------------------------------------------------------------------
// Verilen telefon numarasının geçerli olup olmadığını kontrol eder.
//
// Parametreler:
//   - phone: Kullanıcının girdiği telefon numarası (string).
//   - country: Ülke kodu (örn: "TR", "US").
//
// Dönüş:
//   - bool → Numara geçerliyse true, geçersizse false.
//
// Açıklama:
// Bu fonksiyon, boşluk, tire, parantez gibi yaygın formatlandırma karakterlerini
// temizler ve ardından ilgili ülke regex'i ile doğrulama yapar.
// Böylece kullanıcı “(0532) 123-45-67” gibi farklı formatlarda girse de,
// normalize edilip tutarlı bir doğrulama yapılır.
func IsValidPhoneNumber(phone string, country string) bool {
	pattern, ok := phonePatterns[country]
	if !ok {
		return false
	}

	// Boşluk, tire ve parantezleri temizle
	cleanNumber := regexp.MustCompile(`\s+|-|\(|\)`).ReplaceAllString(phone, "")
	return pattern.MatchString(cleanNumber)
}
