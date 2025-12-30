# Object Validation

The `Object()` validator handles nested objects and complex data structures.

## Basic Usage

```go
schema := v.Make().Shape(map[string]v.Type{
    "address": v.Object().Shape(map[string]v.Type{
        "street":  v.String().Required().Label("Street"),
        "city":    v.String().Required().Label("City"),
        "zipCode": v.String().Length(5).Label("ZIP"),
        "country": v.String().Required().Label("Country"),
    }).Required().Label("Address"),
})
```

## Deep Nesting

```go
schema := v.Make().Shape(map[string]v.Type{
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

## Optional vs Required Objects

```go
schema := v.Make().Shape(map[string]v.Type{
    // Required object - must be present
    "billing": v.Object().Shape(map[string]v.Type{
        "name": v.String().Required(),
        "card": v.CreditCard().Required(),
    }).Required().Label("Billing Info"),

    // Optional object - if present, must be valid
    "shipping": v.Object().Shape(map[string]v.Type{
        "address": v.String().Required(),
        "city":    v.String().Required(),
    }).Label("Shipping Info"),
})
```

## Complete Example

```go
// User profile with nested structures
schema := v.Make().Shape(map[string]v.Type{
    "user": v.Object().Shape(map[string]v.Type{
        "id":       v.Uuid().Required(),
        "username": v.String().Required().Min(3).Max(20),
        "email":    v.String().Email().Required(),

        "profile": v.Object().Shape(map[string]v.Type{
            "firstName": v.String().Required(),
            "lastName":  v.String().Required(),
            "bio":       v.AdvancedString().StripTags().MaxWords(200),
            "avatar":    v.String().URL(),
        }),

        "settings": v.Object().Shape(map[string]v.Type{
            "notifications": v.Boolean().Default(true),
            "theme":         v.String().OneOf([]string{"light", "dark", "auto"}).Default("auto"),
            "language":      v.String().Default("en"),
        }),

        "addresses": v.Array().Elements(v.Object().Shape(map[string]v.Type{
            "type":    v.String().OneOf([]string{"home", "work", "other"}).Required(),
            "street":  v.String().Required(),
            "city":    v.String().Required(),
            "country": v.String().Required(),
        })),
    }).Required(),
})
```

## Notes

- Nested objects preserve their structure in `ValidData()`
- Errors are keyed by dot notation: `address.street`, `profile.contact.email`
- Default values work at any nesting level
