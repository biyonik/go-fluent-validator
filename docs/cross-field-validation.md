# Cross-Field Validation

Cross-field validation allows you to validate relationships between multiple fields, such as password confirmation, date ranges, or budget constraints.

## Basic Usage

Use `.CrossValidate()` to add a function that receives all field data:

```go
schema := v.Make().Shape(map[string]v.Type{
    "password":        v.String().Required().Min(8),
    "passwordConfirm": v.String().Required(),
}).CrossValidate(func(data map[string]any) error {
    pass, _ := data["password"].(string)
    confirm, _ := data["passwordConfirm"].(string)

    if pass != confirm {
        return fmt.Errorf("Passwords do not match")
    }
    return nil
})
```

## Password Confirmation

```go
schema := v.Make().Shape(map[string]v.Type{
    "password": v.String().Required().Password(
        v.WithMinLength(10),
        v.WithRequireSpecial(true),
    ).Label("Password"),
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

## Date Range Validation

```go
schema := v.Make().Shape(map[string]v.Type{
    "startDate": v.Date().Format("2006-01-02").Required().Label("Start Date"),
    "endDate":   v.Date().Format("2006-01-02").Required().Label("End Date"),
}).CrossValidate(func(data map[string]any) error {
    start, _ := data["startDate"].(time.Time)
    end, _ := data["endDate"].(time.Time)

    if end.Before(start) {
        return fmt.Errorf("End date must be after start date")
    }
    return nil
})
```

## Budget Validation

```go
schema := v.Make().Shape(map[string]v.Type{
    "minBudget": v.Number().Required().Label("Minimum Budget"),
    "maxBudget": v.Number().Required().Label("Maximum Budget"),
}).CrossValidate(func(data map[string]any) error {
    min, _ := data["minBudget"].(float64)
    max, _ := data["maxBudget"].(float64)

    if max < min {
        return fmt.Errorf("Maximum budget must be greater than minimum")
    }
    return nil
})
```

## Multiple Cross-Validations

You can chain multiple cross-validations:

```go
schema := v.Make().Shape(map[string]v.Type{
    // fields...
}).CrossValidate(func(data map[string]any) error {
    // First validation
    return nil
}).CrossValidate(func(data map[string]any) error {
    // Second validation
    return nil
})
```

## Error Handling

Cross-validation errors are stored under the `_cross_validation` key:

```go
result := schema.Validate(data)
if result.HasErrors() {
    errors := result.Errors()
    if crossErrors, ok := errors["_cross_validation"]; ok {
        // Handle cross-validation errors
        fmt.Println(crossErrors)
    }
}
```

## Best Practices

1. **Keep it simple** — Each cross-validation should check one relationship
2. **Clear error messages** — Be specific about what failed
3. **Type assertions** — Always handle type assertion failures gracefully
4. **Order matters** — Cross-validation runs after individual field validation
