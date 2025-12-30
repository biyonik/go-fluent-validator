// -----------------------------------------------------------------------------
// Basic Validation Examples
// -----------------------------------------------------------------------------
// Bu dosya, go-fluent-validator kütüphanesinin temel kullanımını gösteren
// örnekleri içerir. Gerçek dünya senaryolarında nasıl kullanılacağını gösterir.
//
// Örnekler:
//   - Simple form validation
//   - User registration
//   - API request validation
//   - Data transformation
//
// Metadata:
// @author   Ahmet ALTUN
// @github   github.com/biyonik
// @linkedin linkedin.com/in/biyonik
// @email    ahmet.altun60@gmail.com
// -----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"log"

	validation "github.com/biyonik/go-fluent-validator"
	"github.com/biyonik/go-fluent-validator/i18n"
)

// Example1_SimpleValidation demonstrates basic string validation
func Example1_SimpleValidation() {
	fmt.Println("=== Example 1: Simple Validation ===")

	// Define schema
	schema := validation.Make().Shape(map[string]validation.Type{
		"username": validation.String().Required().Min(3).Max(20).Label("Username"),
		"email":    validation.String().Required().Email().Trim().Label("Email Address"),
		"age":      validation.Number().Min(18).Integer().Label("Age"),
	})

	// Input data (e.g., from JSON body)
	data := map[string]any{
		"username": "johndoe",
		"email":    "  john@example.com  ", // Will be trimmed
		"age":      25,
	}

	// Validate
	result := schema.Validate(data)

	// Check results
	if result.HasErrors() {
		fmt.Println("❌ Validation failed:")
		for field, errors := range result.Errors() {
			fmt.Printf("  - %s: %v\n", field, errors)
		}
	} else {
		fmt.Println("✅ Validation successful!")

		// Get sanitized data
		validData := result.ValidData()
		fmt.Printf("Email: '%s'\n", validData["email"]) // Output: 'john@example.com' (trimmed)
		fmt.Printf("Age: %v\n", validData["age"])
	}
	fmt.Println()
}

// Example2_UserRegistration demonstrates a complete user registration form
func Example2_UserRegistration() {
	fmt.Println("=== Example 2: User Registration ===")

	schema := validation.Make().Shape(map[string]validation.Type{
		"username":         validation.String().Required().Min(3).Max(20).Label("Username"),
		"email":            validation.String().Required().Email().Trim().Label("Email"),
		"password":         validation.String().Required().Password().Label("Password"),
		"password_confirm": validation.String().Required().Label("Password Confirmation"),
		"age":              validation.Number().Min(13).Integer().Label("Age"),
		"role":             validation.String().OneOf([]string{"user", "admin"}).Default("user"),
	}).CrossValidate(func(data map[string]any) error {
		password, _ := data["password"].(string)
		confirm, _ := data["password_confirm"].(string)

		if password != confirm {
			return fmt.Errorf("passwords do not match")
		}
		return nil
	})

	// Valid registration
	validData := map[string]any{
		"username":         "johndoe",
		"email":            "john@example.com",
		"password":         "MyP@ssw0rd123",
		"password_confirm": "MyP@ssw0rd123",
		"age":              25,
	}

	result := schema.Validate(validData)
	if !result.HasErrors() {
		fmt.Println("✅ Registration successful!")
		fmt.Printf("User role: %s (default applied)\n", result.ValidData()["role"])
	}

	// Invalid registration - weak password
	invalidData := map[string]any{
		"username":         "johndoe",
		"email":            "john@example.com",
		"password":         "weak",
		"password_confirm": "weak",
		"age":              25,
	}

	result = schema.Validate(invalidData)
	if result.HasErrors() {
		fmt.Println("\n❌ Registration failed:")
		for field, errors := range result.Errors() {
			fmt.Printf("  - %s: %v\n", field, errors)
		}
	}
	fmt.Println()
}

// Example3_ConditionalValidation demonstrates When() conditional rules
func Example3_ConditionalValidation() {
	fmt.Println("=== Example 3: Conditional Validation (Payment Methods) ===")

	schema := validation.Make().Shape(map[string]validation.Type{
		"payment_method": validation.String().OneOf([]string{"credit_card", "paypal"}).Required(),
	}).When("payment_method", "credit_card", func() validation.Schema {
		return validation.Make().Shape(map[string]validation.Type{
			"card_number": validation.CreditCard().Required().Label("Card Number"),
			"cvv":         validation.String().Min(3).Max(4).Required().Label("CVV"),
		})
	}).When("payment_method", "paypal", func() validation.Schema {
		return validation.Make().Shape(map[string]validation.Type{
			"paypal_email": validation.String().Email().Required().Label("PayPal Email"),
		})
	})

	// Credit card payment
	creditCardData := map[string]any{
		"payment_method": "credit_card",
		"card_number":    "4532015112830366",
		"cvv":            "123",
	}

	result := schema.Validate(creditCardData)
	if !result.HasErrors() {
		fmt.Println("✅ Credit card payment validated successfully")
	}

	// PayPal payment
	paypalData := map[string]any{
		"payment_method": "paypal",
		"paypal_email":   "user@example.com",
	}

	result = schema.Validate(paypalData)
	if !result.HasErrors() {
		fmt.Println("✅ PayPal payment validated successfully")
	}

	// Invalid - missing required fields for credit card
	invalidData := map[string]any{
		"payment_method": "credit_card",
		// Missing card_number and cvv
	}

	result = schema.Validate(invalidData)
	if result.HasErrors() {
		fmt.Println("\n❌ Payment validation failed:")
		for field, errors := range result.Errors() {
			fmt.Printf("  - %s: %v\n", field, errors)
		}
	}
	fmt.Println()
}

// Example4_InternationalizationSupport demonstrates i18n usage
func Example4_InternationalizationSupport() {
	fmt.Println("=== Example 4: Internationalization (i18n) ===")

	schema := validation.Make().Shape(map[string]validation.Type{
		"email": validation.String().Required().Email().Label("Email"),
	})

	invalidData := map[string]any{
		"email": "", // Empty email
	}

	// English messages
	i18n.SetLocale("en")
	result := schema.Validate(invalidData)
	if result.HasErrors() {
		fmt.Println("English errors:")
		for field, errors := range result.Errors() {
			fmt.Printf("  - %s: %v\n", field, errors)
		}
	}

	// Turkish messages
	i18n.SetLocale("tr")
	result = schema.Validate(invalidData)
	if result.HasErrors() {
		fmt.Println("\nTurkish errors:")
		for field, errors := range result.Errors() {
			fmt.Printf("  - %s: %v\n", field, errors)
		}
	}
	fmt.Println()
}

// Example5_AdvancedStringValidation demonstrates AdvancedString features
func Example5_AdvancedStringValidation() {
	fmt.Println("=== Example 5: Advanced String Sanitization ===")

	schema := validation.Make().Shape(map[string]validation.Type{
		"bio": validation.AdvancedString().
			Required().
			Trim().
			StripTags("<b>", "<i>"). // Only allow <b> and <i> tags
			FilterEmoji(true).       // Remove emojis
			Label("Biography"),
		"slug": validation.AdvancedString().
			Required().
			SanitizeFilename().
			Label("URL Slug"),
	})

	data := map[string]any{
		"bio":  "  <p>Hello! 😊 I love <b>coding</b> and <script>alert('xss')</script></p>  ",
		"slug": "My Awesome Post!!!",
	}

	result := schema.Validate(data)
	if !result.HasErrors() {
		validData := result.ValidData()
		fmt.Println("✅ Sanitization successful!")
		fmt.Printf("Bio: %s\n", validData["bio"])   // HTML and emojis removed
		fmt.Printf("Slug: %s\n", validData["slug"]) // Sanitized for filesystem
	}
	fmt.Println()
}

// Example6_APIRequestValidation demonstrates validating HTTP API requests
func Example6_APIRequestValidation() {
	fmt.Println("=== Example 6: API Request Validation ===")

	// Simulate JSON request body
	jsonBody := `{
		"title": "My Blog Post",
		"content": "This is the content",
		"tags": ["golang", "validation"],
		"status": "published",
		"author": {
			"name": "John Doe",
			"email": "john@example.com"
		}
	}`

	// Define schema
	schema := validation.Make().Shape(map[string]validation.Type{
		"title":   validation.String().Required().Min(5).Max(200),
		"content": validation.String().Required().Min(10),
		"tags":    validation.Array().Min(1).Max(5).Elements(validation.String()),
		"status":  validation.String().OneOf([]string{"draft", "published"}).Default("draft"),
		"author": validation.Object().Required().Shape(map[string]validation.Type{
			"name":  validation.String().Required(),
			"email": validation.String().Email().Required(),
		}),
	})

	// Parse JSON
	var data map[string]any
	if err := json.Unmarshal([]byte(jsonBody), &data); err != nil {
		log.Fatal(err)
	}

	// Validate
	result := schema.Validate(data)
	if result.HasErrors() {
		fmt.Println("❌ API request validation failed:")
		for field, errors := range result.Errors() {
			fmt.Printf("  - %s: %v\n", field, errors)
		}
	} else {
		fmt.Println("✅ API request validated successfully!")
		// You can now safely use result.ValidData()
		validData := result.ValidData()
		fmt.Printf("Post title: %s\n", validData["title"])
		fmt.Printf("Status: %s\n", validData["status"])
	}
	fmt.Println()
}

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║      go-fluent-validator - Basic Examples                 ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	Example1_SimpleValidation()
	Example2_UserRegistration()
	Example3_ConditionalValidation()
	Example4_InternationalizationSupport()
	Example5_AdvancedStringValidation()
	Example6_APIRequestValidation()

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║      All examples completed!                               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}
