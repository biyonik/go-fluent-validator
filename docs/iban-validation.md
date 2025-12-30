# IBAN Validation

The `Iban()` validator handles International Bank Account Number validation with country-specific rules and checksum verification.

## Basic Usage

```go
schema := v.Make().Shape(map[string]v.Type{
    // Any country
    "iban": v.Iban().Required().Label("IBAN"),

    // Specific country
    "ibanTR": v.Iban().Country("TR").Required().Label("Turkish IBAN"),
    "ibanDE": v.Iban().Country("DE").Required().Label("German IBAN"),
})
```

## Country-Specific Validation

```go
schema := v.Make().Shape(map[string]v.Type{
    "bankAccount": v.Iban().
        Country("TR").    // Turkish IBAN format
        Required().
        Label("Bank Account"),
})
```

## Features

The validator performs:

1. **Format check** — Correct pattern for the country
2. **Length validation** — Country-specific length requirements
3. **Checksum verification** — MOD-97 algorithm validation
4. **Country code validation** — Valid ISO country code

## Supported Countries

Common supported countries include:

| Code | Country | Length |
|------|---------|--------|
| TR | Turkey | 26 |
| DE | Germany | 22 |
| GB | United Kingdom | 22 |
| FR | France | 27 |
| ES | Spain | 24 |
| IT | Italy | 27 |
| NL | Netherlands | 18 |

And many more...

## API Reference

| Method | Description | Example |
|--------|-------------|---------|
| `.Required()` | Field must be present | `.Required()` |
| `.Country(code)` | Specific country code | `.Country("TR")` |
| `.Label(name)` | Custom error label | `.Label("IBAN")` |
