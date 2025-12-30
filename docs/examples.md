# Real-World Examples

Complete validation schemas for common use cases.

## User Registration

```go
func UserRegistrationSchema() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "username": v.String().
            Required().
            Min(3).Max(20).
            AlphaNumeric().
            Label("Username"),

        "email": v.String().
            Required().
            Email().
            Trim().
            Label("Email Address"),

        "password": v.String().
            Required().
            Password(
                v.WithMinLength(10),
                v.WithRequireUppercase(true),
                v.WithRequireLowercase(true),
                v.WithRequireNumeric(true),
                v.WithRequireSpecial(true),
                v.WithRejectCommon(true),
            ).
            Label("Password"),

        "passwordConfirm": v.String().
            Required().
            Label("Confirm Password"),

        "birthDate": v.Date().
            Format("2006-01-02").
            Required().
            Before(time.Now().AddDate(-13, 0, 0)). // Must be at least 13
            Label("Birth Date"),

        "termsAccepted": v.Boolean().
            Required().
            Label("Terms & Conditions"),

    }).CrossValidate(func(data map[string]any) error {
        pass, _ := data["password"].(string)
        confirm, _ := data["passwordConfirm"].(string)
        if pass != confirm {
            return fmt.Errorf("Passwords do not match")
        }
        return nil
    })
}
```

## E-Commerce Checkout

```go
func CheckoutSchema() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        // Payment selection
        "paymentMethod": v.String().
            OneOf([]string{"credit_card", "paypal", "bank_transfer"}).
            Required().
            Label("Payment Method"),

        // Credit card fields (conditionally required)
        "cardNumber": v.CreditCard(),
        "cvv":        v.String(),
        "expiryDate": v.String(),

        // PayPal (conditionally required)
        "paypalEmail": v.String(),

        // Bank transfer (conditionally required)
        "iban": v.Iban(),

        // Shipping address
        "shipping": v.Object().Shape(map[string]v.Type{
            "firstName": v.String().Required().Label("First Name"),
            "lastName":  v.String().Required().Label("Last Name"),
            "street":    v.String().Required().Label("Street Address"),
            "city":      v.String().Required().Label("City"),
            "zipCode":   v.String().Length(5).Required().Label("ZIP Code"),
            "country":   v.String().Required().Label("Country"),
            "phone":     v.String().Phone("US").Label("Phone"),
        }).Required().Label("Shipping Address"),

        // Cart items
        "items": v.Array().
            NotEmpty().
            Elements(v.Object().Shape(map[string]v.Type{
                "productId": v.Uuid().Required(),
                "quantity":  v.Number().Integer().Min(1).Max(99).Required(),
                "price":     v.Number().Positive().Required(),
            })).
            Label("Cart Items"),

        // Optional
        "couponCode": v.String(),
        "notes":      v.AdvancedString().StripTags().MaxWords(100),

    }).When("paymentMethod", "credit_card", func() v.Schema {
        return v.Make().Shape(map[string]v.Type{
            "cardNumber": v.CreditCard().Required().Label("Card Number"),
            "cvv":        v.String().Min(3).Max(4).Numeric().Required().Label("CVV"),
            "expiryDate": v.String().Regex(`^(0[1-9]|1[0-2])\/\d{2}$`).Required().Label("Expiry Date"),
        })
    }).When("paymentMethod", "paypal", func() v.Schema {
        return v.Make().Shape(map[string]v.Type{
            "paypalEmail": v.String().Email().Required().Label("PayPal Email"),
        })
    }).When("paymentMethod", "bank_transfer", func() v.Schema {
        return v.Make().Shape(map[string]v.Type{
            "iban": v.Iban().Required().Label("IBAN"),
        })
    })
}
```

## Blog Post

```go
func BlogPostSchema() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "title": v.AdvancedString().
            Required().
            Trim().
            StripTags().
            Min(10).Max(200).
            Label("Title"),

        "slug": v.AdvancedString().
            Required().
            Trim().
            ReplaceTurkishChars().
            CharSet("alphanumeric").
            Min(5).Max(100).
            Label("URL Slug"),

        "content": v.AdvancedString().
            Required().
            Trim().
            StripTags("<p>", "<br>", "<b>", "<i>", "<u>", "<a>", "<ul>", "<ol>", "<li>", "<h2>", "<h3>").
            Min(100).
            Label("Content"),

        "excerpt": v.AdvancedString().
            Trim().
            StripTags().
            MaxWords(50).
            Label("Excerpt"),

        "category": v.String().
            Required().
            OneOf([]string{"technology", "lifestyle", "business", "health"}).
            Label("Category"),

        "tags": v.Array().
            Min(1).Max(10).
            Elements(v.String().Min(2).Max(30)).
            Unique().
            Label("Tags"),

        "featuredImage": v.String().URL().Label("Featured Image"),

        "publishDate": v.Date().
            Format("2006-01-02").
            Required().
            Label("Publish Date"),

        "status": v.String().
            OneOf([]string{"draft", "published", "archived"}).
            Default("draft"),

        "author": v.Object().Shape(map[string]v.Type{
            "id":   v.Uuid().Required(),
            "name": v.String().Required(),
        }).Required(),
    })
}
```

## API Request Validation

```go
func CreateUserAPISchema() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "user": v.Object().Shape(map[string]v.Type{
            "email":    v.String().Email().Required().Trim(),
            "password": v.String().Required().Min(8),
            "profile": v.Object().Shape(map[string]v.Type{
                "firstName": v.String().Required().Min(1).Max(50),
                "lastName":  v.String().Required().Min(1).Max(50),
                "avatar":    v.String().URL(),
            }),
        }).Required(),
        "metadata": v.Object().Shape(map[string]v.Type{
            "source":    v.String().OneOf([]string{"web", "mobile", "api"}),
            "ipAddress": v.String().IP(""),
            "userAgent": v.String().Max(500),
        }),
    })
}

// HTTP Handler
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var data map[string]any
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    result := CreateUserAPISchema().Validate(data)

    if result.HasErrors() {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]any{
            "success": false,
            "errors":  result.Errors(),
        })
        return
    }

    // Use validated data
    validData := result.ValidData()
    // ... create user
}
```

## Contact Form

```go
func ContactFormSchema() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "name": v.String().
            Required().
            Min(2).Max(100).
            Label("Name"),

        "email": v.String().
            Required().
            Email().
            Trim().
            Label("Email"),

        "subject": v.String().
            Required().
            Min(5).Max(200).
            Label("Subject"),

        "message": v.AdvancedString().
            Required().
            Trim().
            StripTags().
            Min(20).Max(5000).
            Label("Message"),

        "phone": v.String().
            Phone("US").
            Label("Phone"),

        "preferredContact": v.String().
            OneOf([]string{"email", "phone"}).
            Default("email"),
    })
}
```

## File Upload Metadata

```go
func FileUploadSchema() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "filename": v.AdvancedString().
            Required().
            SanitizeFilename().
            Min(1).Max(255).
            Label("File Name"),

        "description": v.AdvancedString().
            Trim().
            StripTags().
            MaxWords(100).
            Label("Description"),

        "category": v.String().
            OneOf([]string{"document", "image", "video", "audio", "other"}).
            Required().
            Label("Category"),

        "tags": v.Array().
            Elements(v.String().Min(2).Max(30)).
            Max(20).
            Unique().
            Label("Tags"),

        "isPublic": v.Boolean().Default(false),

        "expiresAt": v.Date().
            Format("2006-01-02").
            After(time.Now()).
            Label("Expiration Date"),
    })
}
```

## Settings/Preferences

```go
func UserSettingsSchema() v.Schema {
    return v.Make().Shape(map[string]v.Type{
        "notifications": v.Object().Shape(map[string]v.Type{
            "email":    v.Boolean().Default(true),
            "push":     v.Boolean().Default(true),
            "sms":      v.Boolean().Default(false),
            "digest":   v.String().OneOf([]string{"daily", "weekly", "never"}).Default("weekly"),
        }),

        "privacy": v.Object().Shape(map[string]v.Type{
            "profileVisibility": v.String().OneOf([]string{"public", "private", "friends"}).Default("public"),
            "showEmail":         v.Boolean().Default(false),
            "showPhone":         v.Boolean().Default(false),
        }),

        "display": v.Object().Shape(map[string]v.Type{
            "theme":    v.String().OneOf([]string{"light", "dark", "auto"}).Default("auto"),
            "language": v.String().OneOf([]string{"en", "tr", "de", "fr", "es"}).Default("en"),
            "timezone": v.String().Default("UTC"),
        }),
    })
}
```
