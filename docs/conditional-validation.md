# Conditional Validation

Conditional validation applies rules dynamically based on other field values using `.When()`. This is perfect for forms where required fields change based on user selection.

## Basic Usage

```go
schema := v.Make().Shape(map[string]v.Type{
    "accountType": v.String().OneOf([]string{"personal", "business"}).Required(),
    "companyName": v.String(),  // Optional by default
    "taxId":       v.String(),  // Optional by default

}).When("accountType", "business", func() v.Schema {
    // Make these required ONLY for business accounts
    return v.Make().Shape(map[string]v.Type{
        "companyName": v.String().Required().Label("Company Name"),
        "taxId":       v.String().Required().Label("Tax ID"),
    })
})
```

## Payment Method Example

```go
schema := v.Make().Shape(map[string]v.Type{
    "paymentMethod": v.String().
        OneOf([]string{"credit_card", "paypal", "bank_transfer"}).
        Required().
        Label("Payment Method"),

    // All optional by default
    "cardNumber":  v.CreditCard(),
    "cvv":         v.String(),
    "paypalEmail": v.String(),
    "iban":        v.Iban(),

}).When("paymentMethod", "credit_card", func() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "cardNumber": v.CreditCard().Required().Label("Card Number"),
        "cvv":        v.String().Min(3).Max(4).Required().Label("CVV"),
    })

}).When("paymentMethod", "paypal", func() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "paypalEmail": v.String().Email().Required().Label("PayPal Email"),
    })

}).When("paymentMethod", "bank_transfer", func() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "iban": v.Iban().Country("TR").Required().Label("IBAN"),
    })
})
```

## Boolean Conditions

```go
schema := v.Make().Shape(map[string]v.Type{
    "needsShipping": v.Boolean().Default(false),
    "address":       v.String(),
    "city":          v.String(),
    "zipCode":       v.String(),

}).When("needsShipping", true, func() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "address": v.String().Required().Label("Address"),
        "city":    v.String().Required().Label("City"),
        "zipCode": v.String().Length(5).Required().Label("ZIP Code"),
    })
})
```

## Multiple Conditions

You can chain multiple `.When()` calls:

```go
schema := v.Make().Shape(map[string]v.Type{
    "role": v.String().OneOf([]string{"user", "admin", "moderator"}).Required(),
    // ...
}).When("role", "admin", func() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "adminCode": v.String().Required(),
    })
}).When("role", "moderator", func() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "region": v.String().Required(),
    })
})
```

## Subscription Example

```go
schema := v.Make().Shape(map[string]v.Type{
    "planType": v.String().OneOf([]string{"free", "pro", "enterprise"}).Required(),
    "billingEmail": v.String(),
    "companySize":  v.Number(),
    "contractLength": v.Number(),

}).When("planType", "pro", func() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "billingEmail": v.String().Email().Required().Label("Billing Email"),
    })

}).When("planType", "enterprise", func() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "billingEmail":   v.String().Email().Required().Label("Billing Email"),
        "companySize":    v.Number().Min(50).Required().Label("Company Size"),
        "contractLength": v.Number().OneOf([]float64{12, 24, 36}).Required().Label("Contract Length"),
    })
})
```

## How It Works

1. Base schema defines all possible fields
2. Fields are optional by default (unless marked `.Required()`)
3. `.When(field, value, schemaFn)` adds/overrides rules when condition matches
4. Multiple `.When()` calls are evaluated in order
5. Only matching conditions apply their rules

## Best Practices

1. **Define all fields in base schema** — Even if optional, declare them
2. **Use clear field names** — The condition field should be obvious
3. **Keep conditions simple** — One field, one value per `.When()`
4. **Order matters** — Later `.When()` calls can override earlier ones
