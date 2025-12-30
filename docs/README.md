# Documentation

Welcome to Go Fluent Validator documentation.

## Getting Started

- [Installation & Quick Start](./getting-started.md)

## Validation Types

- [String Validation](./string-validation.md) — Text, email, URL, patterns, character sets
- [Advanced String (Sanitization)](./advanced-string.md) — XSS protection, HTML stripping, emoji filtering
- [Number Validation](./number-validation.md) — Integer, range, positive/negative
- [Date Validation](./date-validation.md) — Custom formats, ranges, before/after
- [Array Validation](./array-validation.md) — Element validation, uniqueness, size
- [Object Validation](./object-validation.md) — Nested structures, deep validation

## Specialized Types

- [UUID Validation](./uuid-validation.md) — Version-specific UUID validation
- [IBAN Validation](./iban-validation.md) — Country-specific bank account validation
- [Credit Card Validation](./creditcard-validation.md) — Luhn algorithm, card type detection
- [Password Validation](./password-validation.md) — Strength policies, common password rejection

## Advanced Features

- [Cross-Field Validation](./cross-field-validation.md) — Password confirmation, date ranges
- [Conditional Validation](./conditional-validation.md) — Dynamic rules with `.When()`
- [Custom Validators](./custom-validators.md) — Implement your own logic
- [Internationalization (i18n)](./i18n.md) — Multi-language error messages

## Examples

- [Real-World Examples](./examples.md) — User registration, e-commerce, blog posts

## API Reference

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/biyonik/go-fluent-validator).
