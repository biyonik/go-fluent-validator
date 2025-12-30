# Array Validation

The `Array()` validator handles lists with element-level validation, uniqueness checks, and size constraints.

## Basic Usage

```go
schema := v.Make().Shape(map[string]v.Type{
    "tags": v.Array().
        Required().
        Min(1).           // At least 1 element
        Max(10).          // At most 10 elements
        NotEmpty().       // Cannot be empty
        Label("Tags"),
})
```

## Element Validation

Validate each element in the array:

```go
schema := v.Make().Shape(map[string]v.Type{
    // Array of emails
    "emails": v.Array().
        Elements(v.String().Email()).
        Label("Email List"),

    // Array of numbers in range
    "scores": v.Array().
        Elements(v.Number().Between(0, 100)).
        Min(1).
        Label("Scores"),

    // Array of strings with length constraints
    "tags": v.Array().
        Elements(v.String().Min(2).Max(20)).
        Max(10).
        Label("Tags"),
})
```

## Uniqueness

```go
schema := v.Make().Shape(map[string]v.Type{
    // All elements must be unique
    "userIds": v.Array().
        Elements(v.String()).
        Unique().
        Label("User IDs"),
})
```

## Contains Value

```go
schema := v.Make().Shape(map[string]v.Type{
    // Must contain "admin"
    "roles": v.Array().
        Contains("admin").
        Label("Roles"),
})
```

## Array of Objects

```go
schema := v.Make().Shape(map[string]v.Type{
    "users": v.Array().
        Elements(v.Object().Shape(map[string]v.Type{
            "name":  v.String().Required(),
            "email": v.String().Email().Required(),
            "age":   v.Number().Min(0),
        })).
        NotEmpty().
        Label("Users"),
})
```

## Complete Example

```go
// E-commerce cart items
schema := v.Make().Shape(map[string]v.Type{
    "items": v.Array().
        Required().
        NotEmpty().
        Min(1).
        Max(50).
        Elements(v.Object().Shape(map[string]v.Type{
            "productId": v.Uuid().Required(),
            "quantity":  v.Number().Integer().Min(1).Max(99).Required(),
            "price":     v.Number().Positive().Required(),
        })).
        Label("Cart Items"),

    "tags": v.Array().
        Elements(v.String().Min(2).Max(30)).
        Unique().
        Max(20).
        Label("Tags"),
})
```

## API Reference

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
