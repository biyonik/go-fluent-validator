package types

import (
	"fmt"
	"time"

	"github.com/biyonik/go-fluent-validator/core"
)

// DateType
//
// Tarih tabanlı doğrulamaların yönetilmesini sağlayan gelişmiş tip sınıfıdır.
// Bu tip, gelen veriyi hem dönüştürür (string → time.Time) hem de belirli
// koşullar altında doğrular. Kullanıcı bir tarih formatı belirleyebilir,
// minimum ve maksimum tarih aralığı tanımlayabilir.
//
// Sağladığı özellikler:
//   - string veya time.Time değerlerini kabul eder
//   - Kullanıcı özel tarih formatı belirleyebilir (Go time layout formatı)
//   - min() ve max() ile tarih aralığı doğrulaması yapılabilir
//   - Varsayılan format “2006-01-02” olarak ayarlanmıştır
//
// Bu tip genellikle API validasyonlarında, form verilerinde, DTO modellerinde
// tarih alanlarının kontrollü ve güvenli bir şekilde işlenmesi için kullanılır.
//
// Yazar Bilgileri:
//   - @author  Ahmet Altun
//   - @github  https://github.com/biyonik
//   - @linkedin https://linkedin.com/in/biyonik
//   - @email   ahmet.altun60@gmail.com
type DateType struct {
	core.BaseType
	format     string  // Beklenen tarih formatı (Go time layout)
	minDateStr *string // Minimum tarih sınırı (string formatında)
	maxDateStr *string // Maksimum tarih sınırı (string formatında)
}

// Required, alanın boş geçilemeyeceğini belirtir.
//
// Döndürür:
//   - *DateType: zincirleme kullanım için aynı örnek geri döner.
func (d *DateType) Required() *DateType {
	d.SetRequired()
	return d
}

// Label, doğrulama hatalarında görünecek kullanıcı dostu alan adını belirler.
//
// Parametreler:
//   - label (string): Kullanıcıya gösterilecek alan etiketi.
//
// Döndürür:
//   - *DateType
func (d *DateType) Label(label string) *DateType {
	d.SetLabel(label)
	return d
}

// Default, alan için varsayılan bir değer belirler.
//
// Parametreler:
//   - value (string): Varsayılan tarih değeri.
//
// Döndürür:
//   - *DateType
func (d *DateType) Default(value string) *DateType {
	d.SetDefault(value)
	return d
}

// Format, tarihin hangi Go time layout formatına göre parse edileceğini belirler.
//
// Parametreler:
//   - goTimeFormat (string): Örneğin: "2006-01-02", "02.01.2006", "2006-01-02 15:04"
//
// Döndürür:
//   - *DateType
func (d *DateType) Format(goTimeFormat string) *DateType {
	d.format = goTimeFormat
	return d
}

// Min, alanın belirtilen tarihten önce olmamasını sağlar.
//
// Parametreler:
//   - dateStr (string): Minimum tarih (format, seçilen format ile uyumlu olmalıdır)
//
// Döndürür:
//   - *DateType
func (d *DateType) Min(dateStr string) *DateType {
	d.minDateStr = &dateStr
	return d
}

// Max, alanın belirtilen tarihten sonra olmamasını sağlar.
//
// Parametreler:
//   - dateStr (string): Maksimum tarih.
//
// Döndürür:
//   - *DateType
func (d *DateType) Max(dateStr string) *DateType {
	d.maxDateStr = &dateStr
	return d
}

// Transform, gelen değeri string → time.Time formatına dönüştürür.
//
// Dönüşüm Süreci:
//  1. BaseType.Transform uygulanır
//  2. Eğer değer nil ise → nil döner
//  3. Eğer değer already time.Time ise → direkt döner
//  4. Eğer değer string ise → seçilen format ile parse edilir
//
// Parametreler:
//   - value (any): ham veri
//
// Döndürür:
//   - any: dönüştürülmüş time.Time değeri
//   - error: geçersiz format varsa döner
func (d *DateType) Transform(value any) (any, error) {
	value, err := d.BaseType.Transform(value)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	if _, ok := value.(time.Time); ok {
		return value, nil
	}
	str, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("tarih alanı string veya time.Time tipinde olmalıdır")
	}

	layout := d.format
	if layout == "" {
		layout = "2006-01-02"
	}

	parsedDate, err := time.Parse(layout, str)
	if err != nil {
		return nil, fmt.Errorf("geçerli bir tarih formatı değil. Beklenen: %s", layout)
	}
	return parsedDate, nil
}

// Validate, tarih için tüm doğrulama kurallarını çalıştırır.
//
// Gerçekleştirilen kontroller:
//  1. BaseType doğrulamaları (required, nullable kontrolü)
//  2. Değerin time.Time olup olmadığı
//  3. min() kontrolü
//  4. max() kontrolü
//
// Parametreler:
//   - field (string): alan adı (path)
//   - value (any): doğrulanacak değer
//   - result (*core.ValidationResult): sonuç nesnesi
func (d *DateType) Validate(field string, value any, result *core.ValidationResult) {
	d.BaseType.Validate(field, value, result)
	if result.HasErrors() {
		return
	}
	if value == nil {
		return
	}

	parsedDate, ok := value.(time.Time)
	if !ok {
		result.AddError(field, fmt.Sprintf("%s alanı geçerli bir tarih olmalıdır", d.GetLabel(field)))
		return
	}

	fieldName := d.GetLabel(field)
	layout := d.format
	if layout == "" {
		layout = "2006-01-02"
	}

	// Min Date Kontrolü
	if d.minDateStr != nil {
		minDate, err := time.Parse(layout, *d.minDateStr)
		if err != nil {
			result.AddError(field, fmt.Sprintf("%s için tanımlanan min() kuralı geçersiz formatta", fieldName))
		} else if parsedDate.Before(minDate) {
			result.AddError(field, fmt.Sprintf("%s alanı %s tarihinden önce olamaz", fieldName, *d.minDateStr))
		}
	}

	// Max Date Kontrolü
	if d.maxDateStr != nil {
		maxDate, err := time.Parse(layout, *d.maxDateStr)
		if err != nil {
			result.AddError(field, fmt.Sprintf("%s için tanımlanan max() kuralı geçersiz formatta", fieldName))
		} else if parsedDate.After(maxDate) {
			result.AddError(field, fmt.Sprintf("%s alanı %s tarihinden sonra olamaz", fieldName, *d.maxDateStr))
		}
	}
}
