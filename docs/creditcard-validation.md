# Credit Card Validation

The `CreditCard()` validator handles payment card validation with Luhn algorithm verification and card type detection.

## Basic Usage

```go
schema := v.Make().Shape(map[string]v.Type{
    // Any card type
    "cardNumber": v.CreditCard().Required().Label("Card Number"),

    // Specific card type
    "visaCard": v.CreditCard().Type("visa").Required().Label("Visa Card"),
})
```

## Card Type Validation

```go
schema := v.Make().Shape(map[string]v.Type{
    "visa":       v.CreditCard().Type("visa").Label("Visa"),
    "mastercard": v.CreditCard().Type("mastercard").Label("MasterCard"),
    "amex":       v.CreditCard().Type("amex").Label("American Express"),
    "discover":   v.CreditCard().Type("discover").Label("Discover"),
})
```

## Supported Card Types

| Type | Prefix | Length |
|------|--------|--------|
| `visa` | 4 | 16 |
| `mastercard` | 51-55, 2221-2720 | 16 |
| `amex` | 34, 37 | 15 |
| `discover` | 6011, 65 | 16 |
| `diners` | 36, 38 | 14 |
| `jcb` | 35 | 16 |

## Features

The validator performs:

1. **Luhn algorithm** — Checksum validation
2. **Card type detection** — Based on prefix patterns
3. **Length validation** — Card type specific length
4. **Format cleaning** — Handles spaces and dashes

## Complete Example

```go
// Payment form validation
schema := v.Make().Shape(map[string]v.Type{
    "cardNumber": v.CreditCard().Required().Label("Card Number"),
    "cvv":        v.String().Min(3).Max(4).Numeric().Required().Label("CVV"),
    "expiryDate": v.String().Regex(`^(0[1-9]|1[0-2])\/\d{2}$`).Required().Label("Expiry"),
    "cardHolder": v.String().Required().Min(2).Label("Card Holder"),
})
```

## API Reference

| Method | Description | Example |
|--------|-------------|---------|
| `.Required()` | Field must be present | `.Required()` |
| `.Type(cardType)` | Specific card type | `.Type("visa")` |
| `.Label(name)` | Custom error label | `.Label("Card Number")` |
