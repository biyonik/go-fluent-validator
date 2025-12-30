# String Validation

The `String()` validator handles text validation with format checks, length constraints, and pattern matching.

## Basic Rules

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

    // Default value
    "status": v.String().Default("active"),
})
```

## Format Validation

```go
schema := v.Make().Shape(map[string]v.Type{
    // Email
    "email": v.String().Email().Required().Label("Email"),

    // URL (http/https)
    "website": v.String().URL().Label("Website"),

    // IP address
    "ipv4": v.String().IP("v4").Label("IPv4 Address"),
    "ipv6": v.String().IP("v6").Label("IPv6 Address"),
    "anyIP": v.String().IP("").Label("IP Address"),

    // Phone number
    "phoneUS": v.String().Phone("US").Label("US Phone"),
    "phoneTR": v.String().Phone("TR").Label("TR Phone"),

    // MAC address
    "macAddr": v.String().MAC().Label("MAC Address"),

    // Hexadecimal
    "hexColor": v.String().Hex().Label("Hex Color"),

    // Base64
    "encoded": v.String().Base64().Label("Base64 Data"),
})
```

## Character Set Validation

```go
schema := v.Make().Shape(map[string]v.Type{
    // Only letters (a-z, A-Z)
    "letters": v.String().Alpha().Label("Letters Only"),

    // Letters and numbers (a-z, A-Z, 0-9)
    "alphanum": v.String().AlphaNumeric().Label("Alphanumeric"),

    // Numbers only (0-9)
    "numbers": v.String().Numeric().Label("Numbers Only"),
})
```

## Enumeration

```go
schema := v.Make().Shape(map[string]v.Type{
    // Must be one of these values
    "role": v.String().OneOf([]string{"admin", "user", "editor"}).Label("Role"),

    // Must NOT be one of these values
    "username": v.String().NotOneOf([]string{"admin", "root", "system"}).Label("Username"),
})
```

## Transformation

```go
schema := v.Make().Shape(map[string]v.Type{
    // Remove leading/trailing whitespace
    "email": v.String().Trim().Email().Label("Email"),
})
```

> **Note:** For advanced sanitization (XSS protection, HTML stripping), see [Advanced String](./advanced-string.md).

## API Reference

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
| `.NotOneOf(values)` | Value not in list | `.NotOneOf([]string{"x"})` |
| `.Trim()` | Remove whitespace | `.Trim()` |
| `.Default(value)` | Default if missing | `.Default("guest")` |
| `.Label(name)` | Custom error label | `.Label("Username")` |
| `.Custom(fn)` | Custom validator | `.Custom(func(v string) error {...})` |
