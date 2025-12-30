# Custom Validators

Custom validators let you implement your own validation logic for business rules, database lookups, or any specialized requirements.

## Basic Usage

Use `.Custom()` on any validator type:

```go
schema := v.Make().Shape(map[string]v.Type{
    "username": v.String().
        Required().
        Custom(func(value string) error {
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

## String Custom Validator

```go
// Check against reserved words
v.String().Custom(func(value string) error {
    if strings.Contains(value, "badword") {
        return fmt.Errorf("contains prohibited content")
    }
    return nil
})
```

## Number Custom Validator

```go
// Business logic validation
v.Number().Custom(func(value float64) error {
    if value > 1000 && value < 5000 {
        return fmt.Errorf("amount must be either under 1000 or over 5000")
    }
    return nil
})
```

## Database Lookup

```go
schema := v.Make().Shape(map[string]v.Type{
    "email": v.String().
        Email().
        Required().
        Custom(func(value string) error {
            exists, err := db.EmailExists(value)
            if err != nil {
                return fmt.Errorf("database error: %v", err)
            }
            if exists {
                return fmt.Errorf("email already registered")
            }
            return nil
        }).
        Label("Email"),
})
```

## External Service Validation

```go
schema := v.Make().Shape(map[string]v.Type{
    "couponCode": v.String().
        Custom(func(value string) error {
            coupon, err := couponService.GetByCode(value)
            if err != nil {
                return fmt.Errorf("invalid coupon code")
            }

            if coupon.ExpiresAt.Before(time.Now()) {
                return fmt.Errorf("coupon has expired")
            }

            if coupon.UsageCount >= coupon.MaxUsage {
                return fmt.Errorf("coupon usage limit reached")
            }

            return nil
        }).
        Label("Coupon Code"),
})
```

## Combining with Built-in Rules

Custom validators run after built-in rules:

```go
v.String().
    Required().           // Runs first
    Min(3).               // Runs second
    Max(20).              // Runs third
    AlphaNumeric().       // Runs fourth
    Custom(func(v string) error {  // Runs last
        // Only called if all above pass
        return nil
    })
```

## Error Messages

Return descriptive errors:

```go
// ❌ Bad
return fmt.Errorf("invalid")

// ✅ Good
return fmt.Errorf("username '%s' is already taken", value)

// ✅ Good
return fmt.Errorf("price must be a multiple of 0.25, got %.2f", value)
```

## Reusable Validators

Create reusable validator functions:

```go
func notReservedUsername(value string) error {
    reserved := []string{"admin", "root", "system", "api", "www"}
    for _, r := range reserved {
        if strings.EqualFold(value, r) {
            return fmt.Errorf("'%s' is a reserved username", value)
        }
    }
    return nil
}

func uniqueEmail(db *Database) func(string) error {
    return func(value string) error {
        if db.EmailExists(value) {
            return fmt.Errorf("email already registered")
        }
        return nil
    }
}

// Usage
schema := v.Make().Shape(map[string]v.Type{
    "username": v.String().Custom(notReservedUsername),
    "email":    v.String().Custom(uniqueEmail(db)),
})
```

## Best Practices

1. **Return `nil` for success** — Only return error when validation fails
2. **Be specific** — Include the problematic value in error messages
3. **Keep it focused** — One custom validator, one concern
4. **Handle edge cases** — Empty strings, nil values, etc.
5. **Consider performance** — Avoid expensive operations in validators
