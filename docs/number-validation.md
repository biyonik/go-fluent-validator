# Number Validation

The `Number()` validator handles numeric values with range constraints, integer checks, and divisibility rules.

## Basic Usage

```go
schema := v.Make().Shape(map[string]v.Type{
    // Basic number
    "age": v.Number().
        Required().
        Integer().        // Must be whole number
        Min(18).          // Minimum value
        Max(120).         // Maximum value
        Label("Age"),

    // Price with decimals allowed
    "price": v.Number().
        Required().
        Positive().       // Must be > 0
        Max(999999.99).
        Label("Price"),

    // Default value
    "page": v.Number().
        Integer().
        Default(1),
})
```

## Sign Constraints

```go
schema := v.Make().Shape(map[string]v.Type{
    // Must be positive (> 0)
    "quantity": v.Number().Positive().Label("Quantity"),

    // Must be negative (< 0)
    "temperature": v.Number().Negative().Label("Temperature"),
})
```

## Range Validation

```go
schema := v.Make().Shape(map[string]v.Type{
    // Between min and max (inclusive)
    "rating": v.Number().
        Between(1, 5).
        Label("Rating"),

    // Separate min/max
    "score": v.Number().
        Min(0).
        Max(100).
        Label("Score"),
})
```

## Divisibility

```go
schema := v.Make().Shape(map[string]v.Type{
    // Must be divisible by 5
    "quantity": v.Number().
        MultipleOf(5).
        Label("Quantity"),

    // Must be even
    "seats": v.Number().
        MultipleOf(2).
        Integer().
        Label("Seats"),
})
```

## API Reference

| Method | Description | Example |
|--------|-------------|---------|
| `.Required()` | Field must be present | `.Required()` |
| `.Min(n)` | Minimum value | `.Min(0)` |
| `.Max(n)` | Maximum value | `.Max(100)` |
| `.Between(min, max)` | Range (inclusive) | `.Between(1, 10)` |
| `.Integer()` | Must be whole number | `.Integer()` |
| `.Positive()` | Must be > 0 | `.Positive()` |
| `.Negative()` | Must be < 0 | `.Negative()` |
| `.MultipleOf(n)` | Divisible by n | `.MultipleOf(5)` |
| `.Default(value)` | Default if missing | `.Default(0)` |
| `.Label(name)` | Custom error label | `.Label("Age")` |
| `.Custom(fn)` | Custom validator | `.Custom(func(v float64) error {...})` |

## Notes

- All numbers are handled as `float64` internally
- Use `.Integer()` when you need whole numbers
- `Between(min, max)` is inclusive on both ends
