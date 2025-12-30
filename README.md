# 🚀 Go Fluent Validator

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/biyonik/go-fluent-validator.svg)](https://pkg.go.dev/github.com/biyonik/go-fluent-validator)
[![Go Report Card](https://goreportcard.com/badge/github.com/biyonik/go-fluent-validator)](https://goreportcard.com/report/github.com/biyonik/go-fluent-validator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?logo=go)](https://golang.org/dl/)

**Type-safe, chainable, zero-dependency validation library for Go**

*Inspired by Zod and Laravel Validation*

[English](#-english) • [Türkçe](#-türkçe) • [Features](#-features) • [Installation](#-installation) • [Documentation](#-documentation)

</div>

---

## 🌍 Language / Dil

- [🇬🇧 English](#-english)
- [🇹🇷 Türkçe](#-türkçe)

---

## 🇬🇧 English

**Go Fluent Validator** is a powerful, type-safe validation library for Go that combines the elegance of **Zod** with the practicality of **Laravel Validation**. Build complex validation schemas with a clean, fluent API while maintaining zero external dependencies.

Not only does it validate your data, but it also **transforms and sanitizes** it, ensuring your data is clean, safe, and ready to use.

---

## ✨ Features

### 🎯 Core Capabilities

- **🔗 Fluent API**: Chain methods for readable, declarative validation schemas
- **🛡️ Type-Safe**: Distinct validators for String, Number, Boolean, Date, Array, Object, UUID, IBAN, and CreditCard
- **🧹 Built-in Sanitization**: Transform and clean data before validation (XSS protection, HTML stripping, emoji filtering, etc.)
- **🔄 Cross-Field Validation**: Validate interdependent fields (password confirmation, date ranges, etc.)
- **⚡ Conditional Rules**: Apply validation rules dynamically based on other field values using `.When()`
- **🌍 Multi-language Support**: Built-in localization for 11 Languages
- **📦 Zero Dependencies**: Uses only Go standard library
- **🎨 Custom Validators**: Implement your own validation logic easily
- **🔍 Rich Rule Set**: 50+ built-in validation rules

### 🎭 Advanced Features

- **UUID Validation**: Support for versions 1-5 with version-specific validation
- **IBAN Validation**: Country-specific IBAN validation with checksum verification
- **Credit Card Validation**: Luhn algorithm + card type detection (Visa, MasterCard, Amex, etc.)
- **Password Strength**: Comprehensive password policy enforcement (entropy, patterns, common passwords)
- **Network Validation**: IP addresses (v4/v6), Phone numbers (multiple countries), MAC addresses
- **String Sanitization**: XSS protection, HTML tag stripping, filename sanitization, emoji filtering
- **Date Handling**: Flexible date parsing with custom formats, min/max ranges, before/after validation
- **Array Validation**: Element-level validation, uniqueness checks, size constraints
- **Nested Objects**: Deep validation of complex data structures

---

## 📦 Installation

```bash
go get github.com/biyonik/go-fluent-validator
```

**Requirements**: Go 1.16 or higher

---

## 🚀 Quick Start

```go
package main

import (
	"fmt"
	v "github.com/biyonik/go-fluent-validator"
)

func main() {
	// Define a validation schema
	userSchema := v.Make().Shape(map[string]v.Type{
		"username": v.String().Required().Min(3).Max(20).Label("Username"),
		"email":    v.String().Required().Email().Trim().Label("Email"),
		"age":      v.Number().Min(18).Integer().Label("Age"),
		"role":     v.String().OneOf([]string{"admin", "user", "editor"}).Default("user"),
	})

	// Input data (e.g., from JSON request)
	data := map[string]any{
		"username": "biyonik",
		"email":    "  user@example.com  ", // Will be automatically trimmed
		"age":      25,
	}

	// Validate
	result := userSchema.Validate(data)

	// Check results
	if result.HasErrors() {
		fmt.Println("Validation failed:", result.Errors())
	} else {
		fmt.Println("✓ Validation successful!")

		// Get sanitized and validated data
		validData := result.ValidData()
		fmt.Printf("Email: '%s'\n", validData["email"]) // Output: 'user@example.com'
		fmt.Printf("Role: '%s'\n", validData["role"])   // Output: 'user' (default)
	}
}
```

---

## 📚 Documentation

### 📋 Table of Contents

- [Available Types](#available-types)
- [String Validation](#string-validation)
- [Advanced String Sanitization](#advanced-string-sanitization)
- [Number Validation](#number-validation)
- [Date Validation](#date-validation)
- [Array Validation](#array-validation)
- [Object Validation](#object-validation)
- [UUID Validation](#uuid-validation)
- [IBAN Validation](#iban-validation)
- [Credit Card Validation](#credit-card-validation)
- [Password Validation](#password-validation)
- [Cross-Field Validation](#cross-field-validation)
- [Conditional Validation](#conditional-validation)
- [Custom Validators](#custom-validators)
- [Internationalization](#internationalization)
- [Error Handling](#error-handling)

---

### Available Types

| Type | Description | Use Case |
|------|-------------|----------|
| `v.String()` | Basic text validation | Names, descriptions, general text |
| `v.AdvancedString()` | Text with sanitization | User input, HTML content, file names |
| `v.Number()` | Numeric validation | Age, price, quantity, ratings |
| `v.Boolean()` | Boolean validation | Flags, checkboxes, toggles |
| `v.Date()` | Date parsing & validation | Birth dates, deadlines, timestamps |
| `v.Array()` | List validation | Tags, categories, multiple selections |
| `v.Object()` | Nested object validation | Address, profile, complex structures |
| `v.Uuid()` | UUID validation | IDs, unique identifiers |
| `v.Iban()` | IBAN validation | Bank accounts |
| `v.CreditCard()` | Payment card validation | Payment processing |

---

### String Validation

The `String()` validator provides comprehensive text validation capabilities.

#### Basic String Rules

```go
schema := v.Make().Shape(map[string]v.Type{
	// Required field
	"name": v.String().Required().Label("Full Name"),

	// Length constraints
	"username": v.String().Min(3).Max(20).Label("Username"),

	// Exact length
	"zipCode": v.String().Length(5).Label("ZIP Code"),

	// Pattern matching
	"code": v.String().StartsWith("USR-").EndsWith("-END").Label("User Code"),

	// Contains substring
	"description": v.String().Contains("important").Label("Description"),

	// Regular expression
	"customPattern": v.String().Regex(`^[A-Z]{3}-\d{4}$`).Label("Pattern"),

	// Default value for missing fields
	"status": v.String().Default("active"),
})
```

#### Format Validation

```go
schema := v.Make().Shape(map[string]v.Type{
	// Email validation
	"email": v.String().Email().Required().Label("Email"),

	// URL validation (http/https)
	"website": v.String().URL().Label("Website"),

	// IP address (IPv4, IPv6, or both)
	"ipv4": v.String().IP("v4").Label("IPv4 Address"),
	"ipv6": v.String().IP("v6").Label("IPv6 Address"),
	"anyIP": v.String().IP("").Label("IP Address"),

	// Phone number (TR or US)
	"phoneUS": v.String().Phone("US").Label("US Phone"),
	"phoneTR": v.String().Phone("TR").Label("TR Phone"),

	// MAC address
	"macAddr": v.String().MAC().Label("MAC Address"),

	// Hexadecimal string
	"hexColor": v.String().Hex().Label("Hex Color"),

	// Base64 encoded
	"encoded": v.String().Base64().Label("Base64 Data"),
})
```

#### Character Set Validation

```go
schema := v.Make().Shape(map[string]v.Type{
	// Only alphabetic characters (a-z, A-Z)
	"letters": v.String().Alpha().Label("Letters Only"),

	// Alphanumeric only (a-z, A-Z, 0-9)
	"alphanum": v.String().AlphaNumeric().Label("Alphanumeric"),

	// Numeric only (0-9)
	"numbers": v.String().Numeric().Label("Numbers Only"),
})
```

#### String API Reference

| Method | Description | Example |
|--------|-------------|---------|
| `.Required()` | Field must be present and non-empty | `.Required()` |
| `.Min(n)` | Minimum length | `.Min(3)` |
| `.Max(n)` | Maximum length | `.Max(100)` |
| `.Length(n)` | Exact length | `.Length(5)` |
| `.Email()` | Valid email format | `.Email()` |
| `.URL()` | Valid URL (http/https) | `.URL()` |
| `.IP(version)` | IP address ("v4", "v6", "") | `.IP("v4")` |
| `.Phone(country)` | Phone number ("US", "TR") | `.Phone("US")` |
| `.MAC()` | MAC address | `.MAC()` |
| `.Hex()` | Hexadecimal string | `.Hex()` |
| `.Base64()` | Base64 encoded | `.Base64()` |
| `.Alpha()` | Letters only | `.Alpha()` |
| `.AlphaNumeric()` | Letters and numbers | `.AlphaNumeric()` |
| `.Numeric()` | Numbers only | `.Numeric()` |
| `.StartsWith(prefix)` | Starts with string | `.StartsWith("USR-")` |
| `.EndsWith(suffix)` | Ends with string | `.EndsWith(".com")` |
| `.Contains(substring)` | Contains substring | `.Contains("admin")` |
| `.Regex(pattern)` | Matches regex | `.Regex("^[A-Z]+$")` |
| `.OneOf(values)` | Value in list | `.OneOf([]string{"a", "b"})` |
| `.NotOneOf(values)` | Value not in list | `.NotOneOf([]string{"x", "y"})` |
| `.Trim()` | Remove whitespace | `.Trim()` |
| `.Default(value)` | Default if missing | `.Default("guest")` |
| `.Label(name)` | Custom error label | `.Label("Username")` |
| `.Custom(fn)` | Custom validator | `.Custom(func(v string) error {...})` |

---

### Advanced String Sanitization

`AdvancedString()` provides powerful sanitization and transformation capabilities for handling user input safely.

```go
schema := v.Make().Shape(map[string]v.Type{
	"bio": v.AdvancedString().
		Required().
		Trim().                            // Remove leading/trailing whitespace
		StripTags("<b>", "<i>", "<u>").    // Remove HTML except allowed tags (XSS protection)
		EscapeHTML().                      // Escape remaining HTML entities
		FilterEmoji(true).                 // Remove emoji characters
		SanitizeFilename().                // Make safe for filesystems
		CharSet("alphanumeric").           // Allow only a-z, A-Z, 0-9
		MaxWords(100).                     // Limit word count
		Label("Biography"),

	"filename": v.AdvancedString().
		Required().
		SanitizeFilename().                // Remove dangerous characters
		StripPunctuation().                // Remove punctuation
		Label("File Name"),

	"htmlContent": v.AdvancedString().
		StripTags().                       // Remove all HTML
		EscapeHTML().                      // Escape special characters
		Label("Content"),

	"slug": v.AdvancedString().
		Trim().
		ReplaceTurkishChars().             // Convert Turkish chars to ASCII (ş→s, ğ→g)
		CharSet("alphanumeric").
		Label("URL Slug"),

	"domain": v.AdvancedString().
		ValidateDomain().                  // Check valid domain name
		Label("Domain"),
})
```

#### Sanitization Methods

| Method | Description | Example |
|--------|-------------|---------|
| `.Trim()` | Remove leading/trailing whitespace | `.Trim()` |
| `.StripTags(allowed...)` | Remove HTML tags (keep allowed) | `.StripTags("<b>", "<i>")` |
| `.EscapeHTML()` | Escape HTML entities | `.EscapeHTML()` |
| `.FilterEmoji(remove)` | Remove or keep emoji | `.FilterEmoji(true)` |
| `.SanitizeFilename()` | Safe for file systems | `.SanitizeFilename()` |
| `.CharSet(type)` | Allow only char set | `.CharSet("alphanumeric")` |
| `.MaxWords(n)` | Limit word count | `.MaxWords(50)` |
| `.StripPunctuation()` | Remove punctuation | `.StripPunctuation()` |
| `.ReplaceTurkishChars()` | Turkish → ASCII | `.ReplaceTurkishChars()` |
| `.ValidateDomain()` | Validate domain name | `.ValidateDomain()` |

---

### Number Validation

Validate numeric values with constraints.

```go
schema := v.Make().Shape(map[string]v.Type{
	// Basic number validation
	"age": v.Number().
		Required().
		Integer().        // Must be integer
		Min(18).          // Minimum value
		Max(120).         // Maximum value
		Label("Age"),

	// Price with decimals
	"price": v.Number().
		Required().
		Positive().       // Must be > 0
		Max(999999.99).
		Label("Price"),

	// Negative numbers allowed
	"temperature": v.Number().
		Negative().       // Must be < 0
		Label("Temperature"),

	// Range validation
	"rating": v.Number().
		Between(1, 5).    // Between min and max (inclusive)
		Label("Rating"),

	// Multiple of
	"quantity": v.Number().
		MultipleOf(5).    // Must be divisible by 5
		Label("Quantity"),

	// Default value
	"page": v.Number().
		Integer().
		Default(1),
})
```

#### Number API Reference

| Method | Description | Example |
|--------|-------------|---------|
| `.Required()` | Field must be present | `.Required()` |
| `.Min(n)` | Minimum value | `.Min(0)` |
| `.Max(n)` | Maximum value | `.Max(100)` |
| `.Between(min, max)` | Range (inclusive) | `.Between(1, 10)` |
| `.Integer()` | Must be integer | `.Integer()` |
| `.Positive()` | Must be > 0 | `.Positive()` |
| `.Negative()` | Must be < 0 | `.Negative()` |
| `.MultipleOf(n)` | Divisible by n | `.MultipleOf(5)` |
| `.Default(value)` | Default if missing | `.Default(0)` |
| `.Label(name)` | Custom error label | `.Label("Age")` |
| `.Custom(fn)` | Custom validator | `.Custom(func(v float64) error {...})` |

---

### Date Validation

Parse and validate dates with flexible format support.

```go
schema := v.Make().Shape(map[string]v.Type{
	// Parse with custom format
	"birthDate": v.Date().
		Format("2006-01-02").     // Go time format
		Required().
		Before(time.Now()).       // Must be in the past
		Label("Birth Date"),

	// Date range
	"startDate": v.Date().
		Format("2006-01-02").
		Min(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)).
		Max(time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)).
		Label("Start Date"),

	// After specific date
	"eventDate": v.Date().
		Format("2006-01-02 15:04:05").
		After(time.Now()).        // Must be in the future
		Label("Event Date"),

	// Default format (RFC3339)
	"timestamp": v.Date().
		Required().
		Label("Timestamp"),
})
```

#### Date API Reference

| Method | Description | Example |
|--------|-------------|---------|
| `.Format(layout)` | Go time format layout | `.Format("2006-01-02")` |
| `.Required()` | Field must be present | `.Required()` |
| `.Min(date)` | Minimum date | `.Min(time.Now())` |
| `.Max(date)` | Maximum date | `.Max(deadline)` |
| `.Before(date)` | Must be before date | `.Before(time.Now())` |
| `.After(date)` | Must be after date | `.After(startDate)` |
| `.Label(name)` | Custom error label | `.Label("Due Date")` |

---

### Array Validation

Validate arrays with element-level rules.

```go
schema := v.Make().Shape(map[string]v.Type{
	// Basic array
	"tags": v.Array().
		Required().
		Min(1).                   // At least 1 element
		Max(10).                  // At most 10 elements
		NotEmpty().               // Cannot be empty array
		Label("Tags"),

	// Array with element validation
	"emails": v.Array().
		Elements(v.String().Email()).  // Each element must be valid email
		Unique().                      // All elements must be unique
		Label("Email List"),

	// Array of numbers
	"scores": v.Array().
		Elements(v.Number().Between(0, 100)).
		Min(1).
		Label("Scores"),

	// Array of objects
	"users": v.Array().
		Elements(v.Object().Shape(map[string]v.Type{
			"name":  v.String().Required(),
			"email": v.String().Email().Required(),
		})).
		Label("Users"),

	// Contains specific value
	"roles": v.Array().
		Contains("admin").        // Must contain "admin"
		Label("Roles"),
})
```

#### Array API Reference

| Method | Description | Example |
|--------|-------------|---------|
| `.Required()` | Field must be present | `.Required()` |
| `.Min(n)` | Minimum element count | `.Min(1)` |
| `.Max(n)` | Maximum element count | `.Max(100)` |
| `.NotEmpty()` | Must have at least 1 element | `.NotEmpty()` |
| `.Elements(schema)` | Validate each element | `.Elements(v.String())` |
| `.Unique()` | All elements must be unique | `.Unique()` |
| `.Contains(value)` | Must contain value | `.Contains("admin")` |
| `.Label(name)` | Custom error label | `.Label("Items")` |

---

### Object Validation

Validate nested objects and complex structures.

```go
schema := v.Make().Shape(map[string]v.Type{
	// Nested object
	"address": v.Object().Shape(map[string]v.Type{
		"street":  v.String().Required().Label("Street"),
		"city":    v.String().Required().Label("City"),
		"zipCode": v.String().Length(5).Label("ZIP"),
		"country": v.String().Required().Label("Country"),
	}).Required().Label("Address"),

	// Deeply nested
	"profile": v.Object().Shape(map[string]v.Type{
		"firstName": v.String().Required(),
		"lastName":  v.String().Required(),
		"contact": v.Object().Shape(map[string]v.Type{
			"email": v.String().Email().Required(),
			"phone": v.String().Phone("US"),
		}),
		"preferences": v.Object().Shape(map[string]v.Type{
			"theme":    v.String().OneOf([]string{"light", "dark"}),
			"language": v.String().Default("en"),
		}),
	}).Label("Profile"),
})
```

---

### UUID Validation

Validate UUIDs with version-specific checks.

```go
schema := v.Make().Shape(map[string]v.Type{
	// Any UUID version
	"id": v.Uuid().Required().Label("ID"),

	// Specific version
	"uuidV4": v.Uuid().Version(4).Required().Label("UUID v4"),
	"uuidV5": v.Uuid().Version(5).Required().Label("UUID v5"),
})
```

Supported versions: 1, 2, 3, 4, 5

---

### IBAN Validation

Validate International Bank Account Numbers with country-specific rules.

```go
schema := v.Make().Shape(map[string]v.Type{
	// Any country
	"iban": v.Iban().Required().Label("IBAN"),

	// Specific country
	"ibanTR": v.Iban().Country("TR").Required().Label("Turkish IBAN"),
	"ibanDE": v.Iban().Country("DE").Required().Label("German IBAN"),
	"ibanGB": v.Iban().Country("GB").Required().Label("UK IBAN"),
})
```

Features:
- Checksum validation (mod-97 algorithm)
- Country-specific length validation
- Format verification

---

### Credit Card Validation

Validate payment cards with Luhn algorithm and card type detection.

```go
schema := v.Make().Shape(map[string]v.Type{
	// Any card type
	"cardNumber": v.CreditCard().Required().Label("Card Number"),

	// Specific card type
	"visaCard": v.CreditCard().Type("visa").Required().Label("Visa Card"),
	"mastercard": v.CreditCard().Type("mastercard").Required().Label("MasterCard"),
	"amex": v.CreditCard().Type("amex").Required().Label("Amex Card"),
})
```

Supported card types:
- Visa
- MasterCard
- American Express
- Discover
- Diners Club
- JCB

Features:
- Luhn algorithm validation
- Card type detection
- Format verification

---

### Password Validation

Enforce strong password policies with comprehensive rules.

```go
schema := v.Make().Shape(map[string]v.Type{
	"password": v.String().Password(
		v.WithMinLength(10),              // Minimum 10 characters
		v.WithMaxLength(128),             // Maximum 128 characters
		v.WithRequireUppercase(true),     // At least 1 uppercase letter
		v.WithRequireLowercase(true),     // At least 1 lowercase letter
		v.WithRequireNumeric(true),       // At least 1 digit
		v.WithRequireSpecial(true),       // At least 1 special character
		v.WithMinUniqueChars(5),          // At least 5 unique characters
		v.WithRejectCommon(true),         // Reject common passwords
		v.WithCheckKeyboardPatterns(true), // Detect keyboard patterns (qwerty, 123456)
	).Required().Label("Password"),
})
```

#### Password Policy Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithMinLength(n)` | Minimum length | 8 |
| `WithMaxLength(n)` | Maximum length | 128 |
| `WithRequireUppercase(bool)` | Require uppercase letters | false |
| `WithRequireLowercase(bool)` | Require lowercase letters | false |
| `WithRequireNumeric(bool)` | Require digits | false |
| `WithRequireSpecial(bool)` | Require special chars | false |
| `WithMinUniqueChars(n)` | Minimum unique characters | 0 |
| `WithRejectCommon(bool)` | Reject common passwords | false |
| `WithCheckKeyboardPatterns(bool)` | Detect keyboard patterns | false |

Common password detection includes:
- "password", "123456", "qwerty", etc.
- Keyboard patterns (qwertyuiop, asdfghjkl, 12345678)
- Simple sequences

---

### Cross-Field Validation

Validate relationships between fields.

#### Password Confirmation

```go
schema := v.Make().Shape(map[string]v.Type{
	"password":        v.String().Required().Min(8).Label("Password"),
	"passwordConfirm": v.String().Required().Label("Confirm Password"),
}).CrossValidate(func(data map[string]any) error {
	pass, _ := data["password"].(string)
	confirm, _ := data["passwordConfirm"].(string)

	if pass != confirm {
		return fmt.Errorf("Passwords do not match")
	}
	return nil
})
```

#### Date Range Validation

```go
schema := v.Make().Shape(map[string]v.Type{
	"startDate": v.Date().Format("2006-01-02").Required().Label("Start Date"),
	"endDate":   v.Date().Format("2006-01-02").Required().Label("End Date"),
}).CrossValidate(func(data map[string]any) error {
	start, _ := data["startDate"].(time.Time)
	end, _ := data["endDate"].(time.Time)

	if end.Before(start) {
		return fmt.Errorf("End date must be after start date")
	}
	return nil
})
```

#### Budget Validation

```go
schema := v.Make().Shape(map[string]v.Type{
	"minBudget": v.Number().Required().Label("Minimum Budget"),
	"maxBudget": v.Number().Required().Label("Maximum Budget"),
}).CrossValidate(func(data map[string]any) error {
	min, _ := data["minBudget"].(float64)
	max, _ := data["maxBudget"].(float64)

	if max < min {
		return fmt.Errorf("Maximum budget must be greater than minimum")
	}
	return nil
})
```

**Note**: Cross-validation errors are stored in the `_cross_validation` field.

---

### Conditional Validation

Apply rules dynamically based on other field values.

#### Payment Method Example

```go
schema := v.Make().Shape(map[string]v.Type{
	"paymentMethod": v.String().
		OneOf([]string{"credit_card", "paypal", "bank_transfer"}).
		Required().
		Label("Payment Method"),

	// These are optional by default
	"cardNumber": v.CreditCard(),
	"cvv":        v.String(),
	"paypalEmail": v.String(),
	"iban":       v.Iban(),

}).When("paymentMethod", "credit_card", func() v.Schema {
	// Make credit card fields required ONLY if payment method is credit_card
	return v.Make().Shape(map[string]v.Type{
		"cardNumber": v.CreditCard().Required().Label("Card Number"),
		"cvv":        v.String().Min(3).Max(4).Required().Label("CVV"),
	})

}).When("paymentMethod", "paypal", func() v.Schema {
	// Make PayPal email required ONLY if payment method is paypal
	return v.Make().Shape(map[string]v.Type{
		"paypalEmail": v.String().Email().Required().Label("PayPal Email"),
	})

}).When("paymentMethod", "bank_transfer", func() v.Schema {
	// Make IBAN required ONLY if payment method is bank_transfer
	return v.Make().Shape(map[string]v.Type{
		"iban": v.Iban().Country("TR").Required().Label("IBAN"),
	})
})
```

#### Shipping Address Example

```go
schema := v.Make().Shape(map[string]v.Type{
	"needsShipping": v.Boolean().Default(false),
	"address":       v.String(),
	"city":          v.String(),
	"zipCode":       v.String(),

}).When("needsShipping", true, func() v.Schema {
	// Require shipping fields ONLY if needsShipping is true
	return v.Make().Shape(map[string]v.Type{
		"address": v.String().Required().Label("Address"),
		"city":    v.String().Required().Label("City"),
		"zipCode": v.String().Length(5).Required().Label("ZIP Code"),
	})
})
```

#### Account Type Example

```go
schema := v.Make().Shape(map[string]v.Type{
	"accountType": v.String().OneOf([]string{"personal", "business"}).Required(),
	"companyName": v.String(),
	"taxId":       v.String(),

}).When("accountType", "business", func() v.Schema {
	// Require business fields ONLY for business accounts
	return v.Make().Shape(map[string]v.Type{
		"companyName": v.String().Required().Label("Company Name"),
		"taxId":       v.String().Required().Label("Tax ID"),
	})
})
```

---

### Custom Validators

Implement your own validation logic.

#### Simple Custom Validator

```go
schema := v.Make().Shape(map[string]v.Type{
	"username": v.String().
		Required().
		Custom(func(value string) error {
			// Check against reserved usernames
			reserved := []string{"admin", "root", "system"}
			for _, r := range reserved {
				if value == r {
					return fmt.Errorf("username '%s' is reserved", value)
				}
			}
			return nil
		}).
		Label("Username"),
})
```

#### Database Lookup Validator

```go
schema := v.Make().Shape(map[string]v.Type{
	"email": v.String().
		Email().
		Required().
		Custom(func(value string) error {
			// Check if email exists in database
			exists, err := db.EmailExists(value)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("email already registered")
			}
			return nil
		}).
		Label("Email"),
})
```

#### Complex Business Logic

```go
schema := v.Make().Shape(map[string]v.Type{
	"couponCode": v.String().
		Custom(func(value string) error {
			// Validate coupon code
			coupon, err := couponService.GetByCode(value)
			if err != nil {
				return fmt.Errorf("invalid coupon code")
			}

			// Check expiration
			if coupon.ExpiresAt.Before(time.Now()) {
				return fmt.Errorf("coupon has expired")
			}

			// Check usage limit
			if coupon.UsageCount >= coupon.MaxUsage {
				return fmt.Errorf("coupon usage limit reached")
			}

			return nil
		}).
		Label("Coupon Code"),
})
```

---

### Internationalization

The library supports multiple languages for error messages.

#### Setting Locale

```go
import "github.com/biyonik/go-fluent-validator/i18n"

// Set to Turkish
i18n.SetLocale("tr")

// Set to German
i18n.SetLocale("de")

// Set to English (default)
i18n.SetLocale("en")
```

#### Supported Languages

- **English** (`en`) - Default
- **Turkish** (`tr`) - Full Turkish error messages
- **German** (`de`) - Full German error messages

#### Example with Turkish

```go
i18n.SetLocale("tr")

schema := v.Make().Shape(map[string]v.Type{
	"email": v.String().Email().Required().Label("E-posta"),
	"yas":   v.Number().Min(18).Label("Yaş"),
})

result := schema.Validate(map[string]any{
	"email": "invalid-email",
	"yas":   15,
})

if result.HasErrors() {
	// Turkish error messages:
	// - "E-posta geçerli bir e-posta adresi olmalıdır"
	// - "Yaş en az 18 olmalıdır"
	fmt.Println(result.Errors())
}
```

---

### Error Handling

#### Validation Result Structure

```go
result := schema.Validate(data)

// Check if validation failed
if result.HasErrors() {
	// Get all errors
	errors := result.Errors()
	// Returns: map[string][]string
	// Example: {
	//   "email": ["Email must be a valid email address"],
	//   "age": ["Age must be at least 18"]
	// }

	// Get errors for specific field
	emailErrors := errors["email"]

	// Get all error messages as flat array
	allMessages := result.AllErrors()
	// Returns: []string
}

// Get validated and sanitized data
if !result.HasErrors() {
	validData := result.ValidData()
	// Returns: map[string]any
	// All transformations applied (trim, defaults, etc.)
}
```

#### HTTP Response Example

```go
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var data map[string]any
	json.NewDecoder(r.Body).Decode(&data)

	schema := v.Make().Shape(map[string]v.Type{
		"email":    v.String().Email().Required().Label("Email"),
		"password": v.String().Min(8).Required().Label("Password"),
	})

	result := schema.Validate(data)

	if result.HasErrors() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"errors":  result.Errors(),
		})
		return
	}

	// Use validated data
	validData := result.ValidData()
	user := createUser(validData["email"].(string), validData["password"].(string))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    user,
	})
}
```

---

## 🎯 Real-World Examples

### User Registration

```go
func ValidateUserRegistration() v.Schema {
	return v.Make().Shape(map[string]v.Type{
		"username": v.String().
			Required().
			Min(3).Max(20).
			AlphaNumeric().
			Label("Username"),

		"email": v.String().
			Required().
			Email().
			Trim().
			Label("Email Address"),

		"password": v.String().
			Required().
			Password(
				v.WithMinLength(10),
				v.WithRequireUppercase(true),
				v.WithRequireLowercase(true),
				v.WithRequireNumeric(true),
				v.WithRequireSpecial(true),
			).
			Label("Password"),

		"passwordConfirm": v.String().
			Required().
			Label("Confirm Password"),

		"birthDate": v.Date().
			Format("2006-01-02").
			Required().
			Before(time.Now().AddDate(-18, 0, 0)).
			Label("Birth Date"),

		"terms": v.Boolean().
			Required().
			Label("Terms & Conditions"),

	}).CrossValidate(func(data map[string]any) error {
		pass, _ := data["password"].(string)
		confirm, _ := data["passwordConfirm"].(string)
		if pass != confirm {
			return fmt.Errorf("Passwords do not match")
		}
		return nil
	})
}
```

### E-commerce Checkout

```go
func ValidateCheckout() v.Schema {
	return v.Make().Shape(map[string]v.Type{
		// Payment method selection
		"paymentMethod": v.String().
			OneOf([]string{"credit_card", "paypal", "bank_transfer"}).
			Required().
			Label("Payment Method"),

		// Shipping info
		"shippingAddress": v.Object().Shape(map[string]v.Type{
			"firstName": v.String().Required().Label("First Name"),
			"lastName":  v.String().Required().Label("Last Name"),
			"street":    v.String().Required().Label("Street"),
			"city":      v.String().Required().Label("City"),
			"zipCode":   v.String().Length(5).Required().Label("ZIP Code"),
			"country":   v.String().Required().Label("Country"),
		}).Required().Label("Shipping Address"),

		// Cart items
		"items": v.Array().
			NotEmpty().
			Elements(v.Object().Shape(map[string]v.Type{
				"productId": v.Uuid().Required(),
				"quantity":  v.Number().Integer().Min(1).Required(),
			})).
			Label("Cart Items"),

	}).When("paymentMethod", "credit_card", func() v.Schema {
		return v.Make().Shape(map[string]v.Type{
			"cardNumber": v.CreditCard().Required().Label("Card Number"),
			"cvv":        v.String().Min(3).Max(4).Required().Label("CVV"),
			"expiryDate": v.String().
				Regex(`^(0[1-9]|1[0-2])\/\d{2}$`).
				Required().
				Label("Expiry Date"),
		})
	}).When("paymentMethod", "bank_transfer", func() v.Schema {
		return v.Make().Shape(map[string]v.Type{
			"iban": v.Iban().Required().Label("IBAN"),
		})
	})
}
```

### Blog Post

```go
func ValidateBlogPost() v.Schema {
	return v.Make().Shape(map[string]v.Type{
		"title": v.AdvancedString().
			Required().
			Trim().
			Min(10).Max(200).
			StripTags().
			Label("Title"),

		"slug": v.AdvancedString().
			Required().
			Trim().
			ReplaceTurkishChars().
			CharSet("alphanumeric").
			Label("URL Slug"),

		"content": v.AdvancedString().
			Required().
			Trim().
			StripTags("<p>", "<br>", "<b>", "<i>", "<u>", "<a>").
			Min(100).
			Label("Content"),

		"excerpt": v.AdvancedString().
			Trim().
			StripTags().
			MaxWords(50).
			Label("Excerpt"),

		"category": v.String().
			Required().
			OneOf([]string{"tech", "lifestyle", "business"}).
			Label("Category"),

		"tags": v.Array().
			Min(1).Max(10).
			Elements(v.String().Min(2).Max(20)).
			Unique().
			Label("Tags"),

		"publishDate": v.Date().
			Format("2006-01-02").
			Required().
			Label("Publish Date"),

		"featured": v.Boolean().
			Default(false),
	})
}
```

### API Request with File Upload

```go
func ValidateFileUpload() v.Schema {
	return v.Make().Shape(map[string]v.Type{
		"filename": v.AdvancedString().
			Required().
			SanitizeFilename().
			StripPunctuation().
			Label("File Name"),

		"description": v.AdvancedString().
			Trim().
			StripTags().
			MaxWords(100).
			Label("Description"),

		"category": v.String().
			OneOf([]string{"document", "image", "video", "other"}).
			Required().
			Label("Category"),

		"tags": v.Array().
			Elements(v.String().Min(2).Max(30)).
			Max(20).
			Unique().
			Label("Tags"),

		"isPublic": v.Boolean().
			Default(false),
	})
}
```

---

## 🔄 Comparison with Other Libraries

| Feature | Go Fluent Validator | validator/v10 | ozzo-validation | govalidator |
|---------|---------------------|---------------|-----------------|-------------|
| Fluent API | ✅ | ❌ | ✅ | ❌ |
| Type-Safe | ✅ | ✅ | ✅ | ❌ |
| Zero Dependencies | ✅ | ❌ | ❌ | ❌ |
| Sanitization | ✅ | ❌ | ❌ | ✅ |
| Conditional Rules | ✅ | ❌ | ✅ | ❌ |
| Cross-Field Validation | ✅ | ✅ | ✅ | ❌ |
| Custom Validators | ✅ | ✅ | ✅ | ✅ |
| i18n Support | ✅ | ✅ | ❌ | ❌ |
| Password Strength | ✅ | ❌ | ❌ | ❌ |
| IBAN/Credit Card | ✅ | ✅ | ❌ | ✅ |
| Nested Objects | ✅ | ✅ | ✅ | ❌ |
| Array Elements | ✅ | ✅ | ✅ | ❌ |

---

## 🤔 FAQ

### Q: Can I use this with struct tags?
**A:** No, this library is designed for map-based validation (e.g., decoded JSON). For struct tag validation, consider `validator/v10`.

### Q: How do I validate nested JSON?
**A:** Use `v.Object()` and `v.Array()` with `.Shape()` and `.Elements()`:
```go
v.Object().Shape(map[string]v.Type{
	"nested": v.Object().Shape(...),
})
```

### Q: Can I add my own custom rules?
**A:** Yes! Use `.Custom(func(value T) error {...})` on any validator type.

### Q: Is it production-ready?
**A:** Yes! The library is battle-tested, has comprehensive test coverage, and is used in production environments.

### Q: What about performance?
**A:** The library is designed for efficiency with minimal allocations. Sanitization happens in-place where possible.

### Q: Can I validate structs directly?
**A:** The library works with `map[string]any`. Convert your structs to maps or use JSON marshaling:
```go
jsonData, _ := json.Marshal(myStruct)
var data map[string]any
json.Unmarshal(jsonData, &data)
result := schema.Validate(data)
```

### Q: How do I handle file uploads?
**A:** Validate filenames with `v.AdvancedString().SanitizeFilename()`. File content validation should be done separately.

---

## 📝 Best Practices

### 1. Always Use Labels
```go
// ❌ Bad
v.String().Required().Email()

// ✅ Good
v.String().Required().Email().Label("Email Address")
```

### 2. Sanitize User Input
```go
// ✅ Good - Protect against XSS
v.AdvancedString().
	Trim().
	StripTags().
	EscapeHTML().
	Label("User Bio")
```

### 3. Use Cross-Validation for Related Fields
```go
// ✅ Good - Validate password confirmation
schema.CrossValidate(func(data map[string]any) error {
	// Validate relationship between fields
	return nil
})
```

### 4. Define Reusable Schemas
```go
// ✅ Good - Create schema factory functions
func UserSchema() v.Schema {
	return v.Make().Shape(map[string]v.Type{
		// ... fields
	})
}
```

### 5. Use Conditional Validation for Dynamic Forms
```go
// ✅ Good - Apply rules based on context
schema.When("accountType", "business", func() v.Schema {
	// Business-specific rules
	return v.Make().Shape(...)
})
```

---

## 🛣️ Roadmap

- [ ] Support for async validators
- [ ] Schema composition and reuse helpers
- [ ] JSON Schema export
- [ ] OpenAPI schema generation
- [ ] Form generation from schemas
- [ ] GraphQL integration
- [ ] MongoDB validation integration

---

## 🇹🇷 Türkçe

**Go Fluent Validator**, Go için geliştirilmiş, **Zod** ve **Laravel Validation**'dan ilham alan, tip güvenli ve dış bağımlılık içermeyen güçlü bir doğrulama kütüphanesidir. Karmaşık doğrulama şemalarını temiz ve okunabilir bir API ile oluşturmanızı sağlar.

Sadece doğrulamakla kalmaz, aynı zamanda veriyi **dönüştürür ve temizler**, böylece verinizin hem geçerli hem de güvenli olmasını sağlar.

---

## ✨ Özellikler

### 🎯 Temel Yetenekler

- **🔗 Akıcı (Fluent) API**: Okunabilir ve deklaratif doğrulama şemaları için metotları zincirleyin
- **🛡️ Tip Güvenli**: String, Number, Boolean, Date, Array, Object, UUID, IBAN ve CreditCard için özelleşmiş doğrulayıcılar
- **🧹 Yerleşik Temizleme**: Doğrulama öncesi veri dönüştürme (XSS koruması, HTML temizleme, emoji filtreleme vb.)
- **🔄 Çapraz Alan Doğrulama**: Birbirine bağımlı alanları doğrulayın (şifre onayı, tarih aralıkları vb.)
- **⚡ Koşullu Kurallar**: `.When()` kullanarak diğer alan değerlerine göre dinamik kurallar uygulayın
- **🌍 Çoklu Dil Desteği**: 11 dilde yerleşik mesaj desteği
- **📦 Bağımlılık Yok**: Sadece Go standart kütüphanesi kullanılır
- **🎨 Özel Doğrulayıcılar**: Kendi doğrulama mantığınızı kolayca uygulayın
- **🔍 Zengin Kural Seti**: 50+ yerleşik doğrulama kuralı

### 🎭 İleri Seviye Özellikler

- **UUID Doğrulama**: Versiyon 1-5 için versiyona özgü doğrulama desteği
- **IBAN Doğrulama**: Ülkeye özgü IBAN doğrulama ve checksum kontrolü
- **Kredi Kartı Doğrulama**: Luhn algoritması + kart tipi algılama (Visa, MasterCard, Amex vb.)
- **Şifre Gücü**: Kapsamlı şifre politikası uygulama (entropi, desenler, yaygın şifreler)
- **Ağ Doğrulama**: IP adresleri (v4/v6), Telefon numaraları (çoklu ülke), MAC adresleri
- **String Temizleme**: XSS koruması, HTML etiket temizleme, dosya adı güvenliği, emoji filtreleme
- **Tarih İşleme**: Özel formatlarla esnek tarih ayrıştırma, min/max aralıkları, öncesi/sonrası doğrulama
- **Dizi Doğrulama**: Element seviyesinde doğrulama, tekil kontrolleri, boyut kısıtlamaları
- **İç İçe Nesneler**: Karmaşık veri yapılarının derin doğrulaması

---

## 📦 Kurulum

```bash
go get github.com/biyonik/go-fluent-validator
```

**Gereksinimler**: Go 1.16 veya üzeri

---

## 🚀 Hızlı Başlangıç

```go
package main

import (
	"fmt"
	v "github.com/biyonik/go-fluent-validator"
)

func main() {
	// Bir doğrulama şeması tanımlayın
	kullaniciSemasi := v.Make().Shape(map[string]v.Type{
		"kullanici_adi": v.String().Required().Min(3).Max(20).Label("Kullanıcı Adı"),
		"email":         v.String().Required().Email().Trim().Label("E-posta"),
		"yas":           v.Number().Min(18).Integer().Label("Yaş"),
		"rol":           v.String().OneOf([]string{"admin", "user", "editor"}).Default("user"),
	})

	// Gelen veri (örn: JSON isteğinden)
	data := map[string]any{
		"kullanici_adi": "biyonik",
		"email":         "  kullanici@example.com  ", // Otomatik olarak temizlenecek
		"yas":           25,
	}

	// Doğrula
	sonuc := kullaniciSemasi.Validate(data)

	// Sonuçları kontrol et
	if sonuc.HasErrors() {
		fmt.Println("Doğrulama başarısız:", sonuc.Errors())
	} else {
		fmt.Println("✓ Doğrulama başarılı!")

		// Temizlenmiş ve doğrulanmış veriyi al
		gecerliVeri := sonuc.ValidData()
		fmt.Printf("Email: '%s'\n", gecerliVeri["email"]) // Çıktı: 'kullanici@example.com'
		fmt.Printf("Rol: '%s'\n', gecerliVeri["rol"])     // Çıktı: 'user' (varsayılan)
	}
}
```

---

## 📚 Kullanılabilir Tipler

| Tip | Açıklama | Kullanım Alanı |
|-----|----------|----------------|
| `v.String()` | Temel metin doğrulama | İsimler, açıklamalar, genel metinler |
| `v.AdvancedString()` | Temizleme ile metin | Kullanıcı girdisi, HTML içerik, dosya adları |
| `v.Number()` | Sayısal doğrulama | Yaş, fiyat, miktar, puanlar |
| `v.Boolean()` | Mantıksal doğrulama | Bayraklar, onay kutuları, açma/kapama |
| `v.Date()` | Tarih ayrıştırma ve doğrulama | Doğum tarihleri, son tarihler, zaman damgaları |
| `v.Array()` | Liste doğrulama | Etiketler, kategoriler, çoklu seçimler |
| `v.Object()` | İç içe nesne doğrulama | Adres, profil, karmaşık yapılar |
| `v.Uuid()` | UUID doğrulama | ID'ler, benzersiz tanımlayıcılar |
| `v.Iban()` | IBAN doğrulama | Banka hesapları |
| `v.CreditCard()` | Kredi kartı doğrulama | Ödeme işleme |

---

## 🎯 Gerçek Dünya Örnekleri

### Kullanıcı Kaydı

```go
func KullaniciKayitDogrula() v.Schema {
	return v.Make().Shape(map[string]v.Type{
		"kullanici_adi": v.String().
			Required().
			Min(3).Max(20).
			AlphaNumeric().
			Label("Kullanıcı Adı"),

		"email": v.String().
			Required().
			Email().
			Trim().
			Label("E-posta Adresi"),

		"sifre": v.String().
			Required().
			Password(
				v.WithMinLength(10),
				v.WithRequireUppercase(true),
				v.WithRequireLowercase(true),
				v.WithRequireNumeric(true),
				v.WithRequireSpecial(true),
			).
			Label("Şifre"),

		"sifre_tekrar": v.String().
			Required().
			Label("Şifre Tekrarı"),

		"dogum_tarihi": v.Date().
			Format("2006-01-02").
			Required().
			Before(time.Now().AddDate(-18, 0, 0)).
			Label("Doğum Tarihi"),

		"sartlar": v.Boolean().
			Required().
			Label("Kullanım Şartları"),

	}).CrossValidate(func(data map[string]any) error {
		sifre, _ := data["sifre"].(string)
		tekrar, _ := data["sifre_tekrar"].(string)
		if sifre != tekrar {
			return fmt.Errorf("Şifreler eşleşmiyor")
		}
		return nil
	})
}
```

### E-Ticaret Ödeme

```go
func OdemeDogrula() v.Schema {
	return v.Make().Shape(map[string]v.Type{
		"odeme_yontemi": v.String().
			OneOf([]string{"kredi_karti", "havale", "paypal"}).
			Required().
			Label("Ödeme Yöntemi"),

		"teslimat_adresi": v.Object().Shape(map[string]v.Type{
			"ad":       v.String().Required().Label("Ad"),
			"soyad":    v.String().Required().Label("Soyad"),
			"adres":    v.String().Required().Label("Adres"),
			"sehir":    v.String().Required().Label("Şehir"),
			"posta_kodu": v.String().Length(5).Required().Label("Posta Kodu"),
		}).Required().Label("Teslimat Adresi"),

	}).When("odeme_yontemi", "kredi_karti", func() v.Schema {
		return v.Make().Shape(map[string]v.Type{
			"kart_no": v.CreditCard().Required().Label("Kart Numarası"),
			"cvv":     v.String().Min(3).Max(4).Required().Label("CVV"),
		})
	}).When("odeme_yontemi", "havale", func() v.Schema {
		return v.Make().Shape(map[string]v.Type{
			"iban": v.Iban().Country("TR").Required().Label("IBAN"),
		})
	})
}
```

### Blog Yazısı

```go
func BlogYazisiDogrula() v.Schema {
	return v.Make().Shape(map[string]v.Type{
		"baslik": v.AdvancedString().
			Required().
			Trim().
			Min(10).Max(200).
			StripTags().
			Label("Başlık"),

		"icerik": v.AdvancedString().
			Required().
			Trim().
			StripTags("<p>", "<br>", "<b>", "<i>", "<u>", "<a>").
			Min(100).
			Label("İçerik"),

		"kategori": v.String().
			Required().
			OneOf([]string{"teknoloji", "yasam", "is"}).
			Label("Kategori"),

		"etiketler": v.Array().
			Min(1).Max(10).
			Elements(v.String().Min(2).Max(20)).
			Unique().
			Label("Etiketler"),

		"yayın_tarihi": v.Date().
			Format("2006-01-02").
			Required().
			Label("Yayın Tarihi"),
	})
}
```

---

## 📘 API Metodları

Tüm doğrulayıcılarda kullanılabilen genel metodlar:

| Metod | Açıklama | Örnek |
|-------|----------|-------|
| `.Required()` | Alan zorunlu | `.Required()` |
| `.Default(value)` | Varsayılan değer | `.Default("misafir")` |
| `.Label(name)` | Hata mesajı etiketi | `.Label("Kullanıcı Adı")` |
| `.Custom(fn)` | Özel doğrulayıcı | `.Custom(func(v T) error {...})` |

Daha detaylı API referansı için [İngilizce dokümantasyon](#string-validation) bölümüne bakınız.

---

## 🌍 Yerelleştirme

Türkçe hata mesajlarını kullanmak için:

```go
import "github.com/biyonik/go-fluent-validator/i18n"

// Türkçe'ye geç
i18n.SetLocale("tr")
```

**Desteklenen Diller:**
- **İngilizce** (`en`) - Varsayılan
- **Türkçe** (`tr`) - Tam Türkçe hata mesajları
- **Almanca** (`de`) - Tam Almanca hata mesajları

---

## 🤝 Katkıda Bulunma

Katkılarınızı bekliyoruz! Lütfen Pull Request göndermekten çekinmeyin.

1. Projeyi fork edin
2. Feature branch'i oluşturun (`git checkout -b feature/HarikaOzellik`)
3. Değişikliklerinizi commit edin (`git commit -m 'Harika bir özellik ekle'`)
4. Branch'inizi push edin (`git push origin feature/HarikaOzellik`)
5. Pull Request açın

---

## 📄 Lisans

MIT Lisansı altında dağıtılmaktadır. Daha fazla bilgi için `LICENSE` dosyasına bakınız.

---

## 🙏 Teşekkürler

Bu kütüphane şu harika projelerden ilham almıştır:
- [Zod](https://github.com/colinhacks/zod) - TypeScript-first schema validation
- [Laravel Validation](https://laravel.com/docs/validation) - The PHP Framework for Web Artisans

---

## 📞 İletişim

- **GitHub Issues**: [github.com/biyonik/go-fluent-validator/issues](https://github.com/biyonik/go-fluent-validator/issues)
- **Email**: ahmet.altun60@gmail.com

---

<div align="center">

**⭐ Projeyi beğendiyseniz GitHub'da yıldız vermeyi unutmayın! ⭐**

Made with ❤️ by [biyonik](https://github.com/biyonik)

</div>
