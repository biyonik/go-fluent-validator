# Getting Started

## Installation

```bash
go get github.com/biyonik/go-fluent-validator
```

**Requirements:** Go 1.16 or higher

## Basic Usage

```go
package main

import (
    "fmt"
    v "github.com/biyonik/go-fluent-validator"
)

func main() {
    // Define a schema
    schema := v.Make().Shape(map[string]v.Type{
        "username": v.String().Required().Min(3).Max(20).Label("Username"),
        "email":    v.String().Required().Email().Trim().Label("Email"),
        "age":      v.Number().Min(18).Integer().Label("Age"),
        "role":     v.String().OneOf([]string{"admin", "user"}).Default("user"),
    })

    // Data to validate (e.g., from JSON request)
    data := map[string]any{
        "username": "biyonik",
        "email":    "  user@example.com  ", // Will be trimmed
        "age":      25,
    }

    // Validate
    result := schema.Validate(data)

    if result.HasErrors() {
        fmt.Println("Errors:", result.Errors())
        return
    }

    // Get sanitized data
    validData := result.ValidData()
    fmt.Printf("Email: '%s'\n", validData["email"]) // 'user@example.com'
    fmt.Printf("Role: '%s'\n", validData["role"])   // 'user' (default)
}
```

## Core Concepts

### 1. Schema Definition

Schemas are defined using `v.Make().Shape()` with a map of field names to validators:

```go
schema := v.Make().Shape(map[string]v.Type{
    "fieldName": v.String().Required(),
})
```

### 2. Validation Result

The `Validate()` method returns a result object with these methods:

| Method | Returns | Description |
|--------|---------|-------------|
| `HasErrors()` | `bool` | True if validation failed |
| `Errors()` | `map[string][]string` | Field-keyed error messages |
| `AllErrors()` | `[]string` | Flat list of all errors |
| `ValidData()` | `map[string]any` | Sanitized, validated data |

### 3. Method Chaining

All validators support fluent chaining:

```go
v.String().Required().Min(3).Max(100).Email().Trim().Label("Email")
```

### 4. Labels

Use `.Label()` for human-readable error messages:

```go
// Without label: "email must be a valid email address"
// With label: "Email Address must be a valid email address"
v.String().Email().Label("Email Address")
```

### 5. Default Values

Provide defaults for optional fields:

```go
v.String().Default("guest")  // If field is missing, use "guest"
v.Number().Default(0)
v.Boolean().Default(false)
```

## Next Steps

- [String Validation](./string-validation.md) — Most common validation type
- [Advanced String](./advanced-string.md) — Sanitization and XSS protection
- [Examples](./examples.md) — Real-world use cases
