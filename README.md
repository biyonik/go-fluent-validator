# 🚀 Go Fluent Validator

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/biyonik/go-fluent-validator.svg)](https://pkg.go.dev/github.com/biyonik/go-fluent-validator)
[![Go Report Card](https://goreportcard.com/badge/github.com/biyonik/go-fluent-validator)](https://goreportcard.com/report/github.com/biyonik/go-fluent-validator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?logo=go)](https://golang.org/dl/)

</div>

---

**Stop writing 50 lines of validation boilerplate.** Go Fluent Validator lets you validate, sanitize, and transform data in a single readable chain — with zero dependencies.

```go
// Before: struct tags + custom validators + manual sanitization + error handling = 50+ lines

// After:
schema := v.Make().Shape(map[string]v.Type{
    "email":    v.String().Required().Email().Trim(),
    "password": v.String().Required().Password(v.WithMinLength(10), v.WithRequireSpecial(true)),
    "age":      v.Number().Min(18).Integer(),
})

result := schema.Validate(data)
if result.HasErrors() {
    return result.Errors() // map[string][]string with localized messages
}
cleanData := result.ValidData() // sanitized & validated
```

## Why This Library?

| Problem | Go Fluent Validator |
|---------|---------------------|
| Struct tags are hard to read | Fluent, chainable API |
| Validation and sanitization are separate | Built-in XSS protection, HTML stripping, emoji filtering |
| i18n is an afterthought | 11 languages out of the box |
| External dependencies | Zero. Pure Go standard library |

## Installation

```bash
go get github.com/biyonik/go-fluent-validator
```

## Quick Examples

### User Registration

```go
schema := v.Make().Shape(map[string]v.Type{
    "username": v.String().Required().Min(3).Max(20).AlphaNumeric(),
    "email":    v.String().Required().Email().Trim(),
    "password": v.String().Required().Password(
        v.WithMinLength(10),
        v.WithRequireUppercase(true),
        v.WithRequireSpecial(true),
        v.WithRejectCommon(true),  // blocks "password123" etc.
    ),
    "age":      v.Number().Min(18).Integer(),
}).CrossValidate(func(data map[string]any) error {
    // Custom cross-field validation
    return nil
})
```

### Dynamic Payment Form

```go
schema := v.Make().Shape(map[string]v.Type{
    "paymentMethod": v.String().OneOf([]string{"card", "iban"}).Required(),
    "cardNumber":    v.CreditCard(),
    "iban":          v.Iban(),
}).When("paymentMethod", "card", func() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "cardNumber": v.CreditCard().Required(), // Required only when card selected
    })
}).When("paymentMethod", "iban", func() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "iban": v.Iban().Country("TR").Required(),
    })
})
```

### Sanitize User Input (XSS Protection)

```go
schema := v.Make().Shape(map[string]v.Type{
    "bio": v.AdvancedString().
        Trim().
        StripTags("<b>", "<i>").  // Allow only <b> and <i>
        EscapeHTML().
        MaxWords(100),
    "filename": v.AdvancedString().SanitizeFilename(),
})
```

## Available Types

| Type | Use Case |
|------|----------|
| `v.String()` | Text, emails, URLs, patterns |
| `v.AdvancedString()` | User input with sanitization |
| `v.Number()` | Age, price, ratings |
| `v.Boolean()` | Flags, toggles |
| `v.Date()` | Dates with custom formats |
| `v.Array()` | Lists with element validation |
| `v.Object()` | Nested structures |
| `v.Uuid()` | UUID v1-5 |
| `v.Iban()` | Bank accounts (country-specific) |
| `v.CreditCard()` | Visa, MasterCard, Amex... |

## Localization

```go
import "github.com/biyonik/go-fluent-validator/i18n"

i18n.SetLocale("tr") // Turkish error messages
i18n.SetLocale("de") // German
i18n.SetLocale("es") // Spanish
// + 8 more languages
```

## Documentation

📚 **[Full Documentation](./docs/)** — Detailed guides for every feature

📖 **[API Reference](https://pkg.go.dev/github.com/biyonik/go-fluent-validator)** — Complete API docs

## Comparison

| Feature | Go Fluent Validator | validator/v10 | ozzo-validation |
|---------|:------------------:|:-------------:|:---------------:|
| Fluent API | ✅ | ❌ | ✅ |
| Zero Dependencies | ✅ | ❌ | ❌ |
| Built-in Sanitization | ✅ | ❌ | ❌ |
| Conditional Rules (.When) | ✅ | ❌ | ✅ |
| Password Strength | ✅ | ❌ | ❌ |
| i18n (11 languages) | ✅ | ✅ | ❌ |
| IBAN/Credit Card | ✅ | ✅ | ❌ |

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License — see [LICENSE](LICENSE) for details.

---

<div align="center">

**[⭐ Star on GitHub](https://github.com/biyonik/go-fluent-validator)** if you find it useful!

</div>
