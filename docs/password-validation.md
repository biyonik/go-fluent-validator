# Password Validation

The `Password()` method on `String()` provides comprehensive password policy enforcement with strength checks, pattern detection, and common password rejection.

## Basic Usage

```go
schema := v.Make().Shape(map[string]v.Type{
    "password": v.String().Password(
        v.WithMinLength(10),
        v.WithRequireUppercase(true),
        v.WithRequireLowercase(true),
        v.WithRequireNumeric(true),
        v.WithRequireSpecial(true),
    ).Required().Label("Password"),
})
```

## Policy Options

### Length Requirements

```go
v.String().Password(
    v.WithMinLength(10),   // Minimum 10 characters
    v.WithMaxLength(128),  // Maximum 128 characters
)
```

### Character Requirements

```go
v.String().Password(
    v.WithRequireUppercase(true),  // At least 1 A-Z
    v.WithRequireLowercase(true),  // At least 1 a-z
    v.WithRequireNumeric(true),    // At least 1 0-9
    v.WithRequireSpecial(true),    // At least 1 special char
)
```

### Uniqueness

```go
v.String().Password(
    v.WithMinUniqueChars(5),  // At least 5 different characters
)
```

### Security Checks

```go
v.String().Password(
    v.WithRejectCommon(true),         // Reject "password", "123456", etc.
    v.WithCheckKeyboardPatterns(true), // Detect "qwerty", "asdfgh", etc.
)
```

## Complete Example

```go
// High-security password policy
schema := v.Make().Shape(map[string]v.Type{
    "password": v.String().Password(
        v.WithMinLength(12),
        v.WithMaxLength(128),
        v.WithRequireUppercase(true),
        v.WithRequireLowercase(true),
        v.WithRequireNumeric(true),
        v.WithRequireSpecial(true),
        v.WithMinUniqueChars(6),
        v.WithRejectCommon(true),
        v.WithCheckKeyboardPatterns(true),
    ).Required().Label("Password"),

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

## Common Passwords Blocked

When `WithRejectCommon(true)` is enabled, these are blocked:

- password, password1, password123
- 123456, 12345678, 123456789
- qwerty, qwertyuiop
- letmein, welcome, admin
- And many more...

## Keyboard Patterns Blocked

When `WithCheckKeyboardPatterns(true)` is enabled:

- qwertyuiop, asdfghjkl, zxcvbnm
- 1234567890
- !@#$%^&*()

## API Reference

| Option | Description | Default |
|--------|-------------|---------|
| `WithMinLength(n)` | Minimum length | 8 |
| `WithMaxLength(n)` | Maximum length | 128 |
| `WithRequireUppercase(bool)` | Require A-Z | false |
| `WithRequireLowercase(bool)` | Require a-z | false |
| `WithRequireNumeric(bool)` | Require 0-9 | false |
| `WithRequireSpecial(bool)` | Require special chars | false |
| `WithMinUniqueChars(n)` | Minimum unique characters | 0 |
| `WithRejectCommon(bool)` | Reject common passwords | false |
| `WithCheckKeyboardPatterns(bool)` | Detect keyboard patterns | false |
