# Go Fluent Validator

[![Go Reference](https://pkg.go.dev/badge/github.com/biyonik/go-fluent-validator.svg)](https://pkg.go.dev/github.com/biyonik/go-fluent-validator)
[![Go Report Card](https://goreportcard.com/badge/github.com/biyonik/go-fluent-validator)](https://goreportcard.com/report/github.com/biyonik/go-fluent-validator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Go Fluent Validator** is a type-safe, chainable, and zero-dependency validation library for Go, heavily inspired by **Zod** and **Laravel Validation**. It allows you to build complex validation schemas programmatically with a clean, fluent API.

It supports **data transformation** (sanitization) alongside validation, ensuring your data is not only valid but also clean and safe to use.

---

### 🌍 Language / Dil
- [🇬🇧 English](#-english)
- [🇹🇷 Türkçe](#-türkçe)

---

## 🇬🇧 English

### ✨ Features

- **Fluent API:** Construct schemas by chaining methods (`v.String().Email().Required()`).
- **Type-Safe:** Distinct validators for `String`, `Number`, `Boolean`, `Date`, `Array`, and `Object`.
- **Rich Rule Set:** Built-in support for `UUID`, `IBAN`, `CreditCard`, `IP`, `Phone`, `Password` strength, and more.
- **Sanitization (Transformation):** Cleanse your data *before* validation (e.g., `Trim()`, `StripTags()`, `SanitizeFilename()`).
- **Cross-Field Validation:** Validate fields dependent on others (e.g., password confirmation).
- **Conditional Rules:** Apply rules dynamically based on other field values using `.When()`.
- **Zero Dependencies:** Built using only the Go standard library.

### 📦 Installation

```bash
go get [github.com/biyonik/go-fluent-validator](https://github.com/biyonik/go-fluent-validator)
```

### 🚀 Usage

#### 1. Basic Validation

Define a schema, validate a map (e.g., decoded JSON), and get sanitized results.

```go
package main

import (
	"fmt"
	v "[github.com/biyonik/go-fluent-validator](https://github.com/biyonik/go-fluent-validator)"
)

func main() {
	// 1. Define the Schema
	userSchema := v.Make().Shape(map[string]v.Type{
		"username": v.String().Required().Min(3).Max(20).Label("Username"),
		"email":    v.String().Required().Email().Trim().Label("Email Address"),
		"age":      v.Number().Min(18).Integer().Label("Age"),
		"role":     v.String().OneOf([]string{"admin", "user", "editor"}).Default("user"),
	})

	// 2. Input Data (e.g., from JSON body)
	data := map[string]any{
		"username": "biyonik",
		"email":    "  example@domain.com  ", // Will be trimmed automatically
		"age":      25,
	}

	// 3. Validate
	result := userSchema.Validate(data)

	// 4. Check Results
	if result.HasErrors() {
		fmt.Println("Validation failed:", result.Errors())
	} else {
		fmt.Println("Validation successful!")
		
		// Get sanitized data
		validData := result.ValidData()
		fmt.Printf("Email: '%s'\n", validData["email"]) // Output: 'example@domain.com'
		fmt.Printf("Role: '%s'\n", validData["role"])   // Output: 'user' (default value)
	}
}
```

#### 2. Advanced String & Sanitization

Use `AdvancedString()` to apply powerful sanitization filters before validation logic runs.

```go
v.AdvancedString().
    Required().
    Trim().                   // Remove whitespace
    StripTags("<b>", "<i>").  // Remove all HTML tags except <b> and <i>
    FilterEmoji(true).        // Remove emojis
    SanitizeFilename().       // Make safe for file systems
    CharSet("alphanumeric").  // Allow only a-z, A-Z, 0-9
    Label("Bio")
```

#### 3. Password Validation

Built-in robust password policy enforcement.

```go
v.String().Password(
    v.WithMinLength(8),
    v.WithRequireUppercase(true),
    v.WithRequireLowercase(true),
    v.WithRequireNumeric(true),
    v.WithRequireSpecial(true),
).Required()
```

#### 4. Cross-Field Validation

Validate fields that depend on each other (like password confirmation) using `.CrossValidate()`.

```go
schema := v.Make().Shape(map[string]v.Type{
    "password":        v.String().Required().Min(8).Label("Password"),
    "passwordConfirm": v.String().Required().Label("Confirm Password"),
}).CrossValidate(func(data map[string]any) error {
    pass, _ := data["password"].(string)
    confirm, _ := data["passwordConfirm"].(string)

    if pass != confirm {
        // This error is added to the "_cross_validation" field
        return fmt.Errorf("Passwords do not match")
    }
    return nil
})
```

#### 5. Conditional Validation (`When`)

Apply rules dynamically based on the value of another field.

```go
paymentSchema := v.Make().Shape(map[string]v.Type{
    "paymentMethod": v.String().OneOf([]string{"credit_card", "paypal"}).Required(),
    "paypalEmail":   v.String().Email(), // Optional by default
    
}).When("paymentMethod", "paypal", func() v.Schema {
    // This schema is merged ONLY if paymentMethod is "paypal"
    return v.Make().Shape(map[string]v.Type{
        "paypalEmail": v.String().Required().Label("PayPal Email"),
    })
}).When("paymentMethod", "credit_card", func() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "cardNumber": v.CreditCard().Required(),
        "cvv":        v.String().Min(3).Max(4).Required(),
    })
})
```

### 📚 Available Types

| Type | Description | Example Methods |
|------|-------------|-----------------|
| `v.String()` | Text validation | `.Email()`, `.URL()`, `.Min()`, `.Password()` |
| `v.AdvancedString()` | Text + Sanitization | `.StripTags()`, `.EscapeHTML()`, `.FilterEmoji()` |
| `v.Number()` | Numeric validation | `.Min()`, `.Max()`, `.Integer()` |
| `v.Boolean()` | Boolean validation | `.Default(false)` |
| `v.Date()` | Date parsing & validation | `.Format("2006-01-02")`, `.Min()`, `.Max()` |
| `v.Array()` | List validation | `.Min()`, `.Max()`, `.Elements(schema)` |
| `v.Object()` | Nested object validation | `.Shape(map...)` |
| `v.Uuid()` | UUID validation | `.Version(4)` |
| `v.Iban()` | IBAN validation | `.Country("TR")` |
| `v.CreditCard()` | Payment card validation | `.Type("visa")` |

---

## 🇹🇷 Türkçe

**Go Fluent Validator**, Go için geliştirilmiş, **Zod** ve **Laravel Validation**'dan ilham alan, tip güvenli ve dış bağımlılık içermeyen bir doğrulama kütüphanesidir. Karmaşık doğrulama şemalarını zincirleme (chainable) metotlarla oluşturmanızı sağlar.

Sadece doğrulamakla kalmaz, veriyi **temizler (sanitize)** ve dönüştürür.

### ✨ Özellikler

- **Akıcı (Fluent) API:** Metotları zincirleyerek okunabilir şemalar oluşturun (`v.String().Email().Required()`).
- **Tip Güvenli:** `String`, `Number`, `Boolean`, `Date`, `Array` ve `Object` için özelleşmiş doğrulayıcılar.
- **Zengin Kural Seti:** `UUID`, `IBAN`, `Kredi Kartı`, `IP`, `Telefon`, `Parola Gücü` ve daha fazlası yerleşik olarak gelir.
- **Veri Temizleme (Sanitization):** Veriyi doğrulamadan *önce* temizleyin (`Trim()`, `HTML Etiketlerini Kaldır` vb.).
- **Çapraz Alan Doğrulama:** Alanları birbirine göre doğrulayın (örn: Şifre tekrarı kontrolü).
- **Koşullu Kurallar:** `.When()` metodu ile bir alanın değerine göre diğer alanlara dinamik kurallar uygulayın.
- **Bağımlılık Yok:** Sadece Go standart kütüphanesi kullanılarak geliştirilmiştir.

### 📦 Kurulum

```bash
go get [github.com/biyonik/go-fluent-validator](https://github.com/biyonik/go-fluent-validator)
```

### 🚀 Kullanım Örnekleri

#### 1. Temel Doğrulama

Şemayı tanımlayın, gelen veriyi (örn: JSON) doğrulayın ve temizlenmiş veriyi alın.

```go
package main

import (
	"fmt"
	v "[github.com/biyonik/go-fluent-validator](https://github.com/biyonik/go-fluent-validator)"
)

func main() {
	// 1. Şema Tanımı
	kullaniciSemasi := v.Make().Shape(map[string]v.Type{
		"kullanici_adi": v.String().Required().Min(3).Max(20).Label("Kullanıcı Adı"),
		"email":         v.String().Required().Email().Trim().Label("E-posta Adresi"),
		"yas":           v.Number().Min(18).Integer().Label("Yaş"),
		"rol":           v.String().OneOf([]string{"admin", "user"}).Default("user"),
	})

	// 2. Gelen Veri
	data := map[string]any{
		"kullanici_adi": "biyonik",
		"email":         "  ornek@domain.com  ", // Otomatik olarak Trim() uygulanır
		"yas":           25,
	}

	// 3. Doğrulama
	sonuc := kullaniciSemasi.Validate(data)

	// 4. Sonuç Kontrolü
	if sonuc.HasErrors() {
		fmt.Println("Hatalar:", sonuc.Errors())
	} else {
		fmt.Println("Başarılı!")
		
		// Temizlenmiş ve doğrulanmış veriyi al
		temizVeri := sonuc.ValidData()
		fmt.Printf("Email: '%s'\n", temizVeri["email"]) // Çıktı: 'ornek@domain.com'
		fmt.Printf("Rol: '%s'\n", temizVeri["rol"])     // Çıktı: 'user' (varsayılan)
	}
}
```

#### 2. Gelişmiş String & Temizleme

Input temizliği ve güvenliği için `AdvancedString()` kullanın.

```go
v.AdvancedString().
    Required().
    Trim().                   // Boşlukları temizle
    StripTags("<b>", "<i>").  // <b> ve <i> hariç tüm HTML'i temizle (XSS koruması)
    FilterEmoji(true).        // Emojileri kaldır
    SanitizeFilename().       // Dosya ismi için güvenli hale getir
    CharSet("alphanumeric").  // Sadece a-z, A-Z, 0-9 karakterlerine izin ver
    Label("Biyografi")
```

#### 3. Parola Kuralları

Güçlü parola politikalarını kolayca uygulayın.

```go
v.String().Password(
    v.WithMinLength(8),
    v.WithRequireUppercase(true), // Büyük harf zorunlu
    v.WithRequireLowercase(true), // Küçük harf zorunlu
    v.WithRequireNumeric(true),   // Rakam zorunlu
    v.WithRequireSpecial(true),   // Özel karakter zorunlu
).Required()
```

#### 4. Çapraz Alan Doğrulama (Şifre Eşleştirme)

İki alanı birbiriyle karşılaştırmak için `.CrossValidate()` kullanın.

```go
sema := v.Make().Shape(map[string]v.Type{
    "sifre":        v.String().Required().Min(8).Label("Şifre"),
    "sifre_tekrar": v.String().Required().Label("Şifre Tekrar"),
}).CrossValidate(func(data map[string]any) error {
    sifre, _ := data["sifre"].(string)
    tekrar, _ := data["sifre_tekrar"].(string)

    if sifre != tekrar {
        return fmt.Errorf("Şifreler eşleşmiyor")
    }
    return nil
})
```

#### 5. Koşullu Doğrulama (`When`)

Bir kuralı sadece belirli bir koşul sağlandığında devreye sokun. Dinamik formlar için idealdir.

```go
odemeSemasi := v.Make().Shape(map[string]v.Type{
    "odeme_yontemi": v.String().OneOf([]string{"kredi_karti", "havale"}).Required(),
    "iban":          v.Iban().Country("TR"), // Varsayılan olarak zorunlu değil
    
}).When("odeme_yontemi", "havale", func() v.Schema {
    // Sadece ödeme yöntemi "havale" ise IBAN zorunlu olsun
    return v.Make().Shape(map[string]v.Type{
        "iban": v.Iban().Country("TR").Required().Label("IBAN Numarası"),
    })
}).When("odeme_yontemi", "kredi_karti", func() v.Schema {
    // Kredi kartı seçilirse kart bilgileri zorunlu
    return v.Make().Shape(map[string]v.Type{
        "kart_no": v.CreditCard().Required(),
        "cvv":     v.String().Min(3).Max(4).Required(),
    })
})
```

### 📚 Kullanılabilir Tipler

| Tip | Açıklama | Örnek Metotlar |
|-----|----------|----------------|
| `v.String()` | Metin doğrulama | `.Email()`, `.URL()`, `.Min()`, `.Password()` |
| `v.AdvancedString()` | Metin + Temizleme | `.StripTags()`, `.EscapeHTML()`, `.FilterEmoji()` |
| `v.Number()` | Sayısal doğrulama | `.Min()`, `.Max()`, `.Integer()` |
| `v.Boolean()` | Mantıksal doğrulama | `.Default(false)` |
| `v.Date()` | Tarih ayrıştırma ve doğrulama | `.Format("2006-01-02")`, `.Min()`, `.Max()` |
| `v.Array()` | Dizi/Liste doğrulama | `.Min()`, `.Max()`, `.Elements(sema)` |
| `v.Object()` | İç içe nesne doğrulama | `.Shape(map...)` |
| `v.Uuid()` | UUID doğrulama | `.Version(4)` |
| `v.Iban()` | IBAN doğrulama | `.Country("TR")` |
| `v.CreditCard()` | Kredi kartı doğrulama | `.Type("visa")` |

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.