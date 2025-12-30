// -----------------------------------------------------------------------------
// Advanced String Validation Comprehensive Tests
// -----------------------------------------------------------------------------
// Bu test dosyası, AdvancedString tipinin XSS, sanitization ve güvenlik
// özelliklerini hardcore şekilde test eder.
// -----------------------------------------------------------------------------

package tests

import (
	"strings"
	"testing"

	validation "github.com/biyonik/go-fluent-validator"
)

// TestAdvancedString_MaxWords tests word limit validation
func TestAdvancedString_MaxWords(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		maxWords  int
		wantError bool
	}{
		{
			name:      "within limit - 5 words",
			value:     "This is a test message",
			maxWords:  10,
			wantError: false,
		},
		{
			name:      "exactly at limit - 10 words",
			value:     "One two three four five six seven eight nine ten",
			maxWords:  10,
			wantError: false,
		},
		{
			name:      "exceeds limit - 15 words",
			value:     "This is a very long message with too many words that exceeds the maximum allowed limit",
			maxWords:  10,
			wantError: true,
		},
		{
			name:      "empty string - 0 words",
			value:     "",
			maxWords:  5,
			wantError: false,
		},
		{
			name:      "single word",
			value:     "Hello",
			maxWords:  1,
			wantError: false,
		},
		{
			name:      "multiple spaces between words",
			value:     "Hello    world    test",
			maxWords:  5,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"text": validation.AdvancedString().
					MaxWords(tt.maxWords).
					Label("Content"),
			})

			result := schema.Validate(map[string]any{
				"text": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// TestAdvancedString_StripPunctuation tests punctuation removal
func TestAdvancedString_StripPunctuation(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "remove basic punctuation",
			value:    "Hello, World!",
			expected: "Hello World",
		},
		{
			name:     "remove all special chars",
			value:    "Test@#$%String&*()",
			expected: "TestString",
		},
		{
			name:     "keep alphanumeric",
			value:    "ABC123xyz",
			expected: "ABC123xyz",
		},
		{
			name:     "complex sentence",
			value:    "Hello, how are you? I'm fine! Thanks.",
			expected: "Hello how are you Im fine Thanks",
		},
		{
			name:     "only punctuation",
			value:    "!@#$%^&*()",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"text": validation.AdvancedString().
					StripPunctuation().
					Label("Text"),
			})

			result := schema.Validate(map[string]any{
				"text": tt.value,
			})

			if result.HasErrors() {
				t.Errorf("unexpected error: %v", result.Errors())
			}

			validData := result.ValidData()
			got := validData["text"].(string)
			if got != tt.expected {
				t.Errorf("got '%s', want '%s'", got, tt.expected)
			}
		})
	}
}

// TestAdvancedString_ReplaceTurkishChars tests Turkish to ASCII conversion
func TestAdvancedString_ReplaceTurkishChars(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "all lowercase turkish chars",
			value:    "şğüöçı",
			expected: "sguoci",
		},
		{
			name:     "all uppercase turkish chars",
			value:    "ŞĞÜÖÇİ",
			expected: "SGUOCI",
		},
		{
			name:     "mixed case",
			value:    "İstanbul Şehri",
			expected: "Istanbul Sehri",
		},
		{
			name:     "SEO friendly slug",
			value:    "Türkçe Başlık Örneği",
			expected: "Turkce Baslik Ornegi",
		},
		{
			name:     "no turkish chars",
			value:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "common turkish words",
			value:    "çiçek göğüs ışık",
			expected: "cicek gogus isik",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"text": validation.AdvancedString().
					ReplaceTurkishChars().
					Label("Text"),
			})

			result := schema.Validate(map[string]any{
				"text": tt.value,
			})

			if result.HasErrors() {
				t.Errorf("unexpected error: %v", result.Errors())
			}

			validData := result.ValidData()
			got := validData["text"].(string)
			if got != tt.expected {
				t.Errorf("got '%s', want '%s'", got, tt.expected)
			}
		})
	}
}

// TestAdvancedString_ValidateDomain tests domain validation
func TestAdvancedString_ValidateDomain(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{
			name:      "valid domain",
			value:     "example.com",
			wantError: false,
		},
		{
			name:      "valid subdomain",
			value:     "sub.example.com",
			wantError: false,
		},
		{
			name:      "valid multi-level subdomain",
			value:     "a.b.c.example.com",
			wantError: false,
		},
		{
			name:      "valid hyphenated domain",
			value:     "my-domain.com",
			wantError: false,
		},
		{
			name:      "invalid - no TLD",
			value:     "example",
			wantError: true,
		},
		{
			name:      "invalid - starts with hyphen",
			value:     "-example.com",
			wantError: true,
		},
		{
			name:      "invalid - ends with hyphen",
			value:     "example-.com",
			wantError: true,
		},
		{
			name:      "invalid - contains space",
			value:     "exam ple.com",
			wantError: true,
		},
		{
			name:      "invalid - special chars",
			value:     "exam@ple.com",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"domain": validation.AdvancedString().
					ValidateDomain().
					Required().
					Label("Domain"),
			})

			result := schema.Validate(map[string]any{
				"domain": tt.value,
			})

			if (result.HasErrors()) != tt.wantError {
				t.Errorf("got error = %v, want error = %v", result.HasErrors(), tt.wantError)
				if result.HasErrors() {
					t.Logf("errors: %v", result.Errors())
				}
			}
		})
	}
}

// TestAdvancedString_XSS_Protection tests XSS attack prevention
func TestAdvancedString_XSS_Protection(t *testing.T) {
	maliciousInputs := []struct {
		name     string
		input    string
		contains []string // Sanitized output should NOT contain these
	}{
		{
			name:     "script tag injection",
			input:    "<script>alert('XSS')</script>",
			contains: []string{"<script>", "alert"},
		},
		{
			name:     "img onerror injection",
			input:    `<img src=x onerror="alert('XSS')">`,
			contains: []string{"onerror", "alert"},
		},
		{
			name:     "javascript protocol",
			input:    `<a href="javascript:alert('XSS')">Click</a>`,
			contains: []string{"javascript:", "alert"},
		},
		{
			name:     "event handler",
			input:    `<div onload="alert('XSS')">Test</div>`,
			contains: []string{"onload", "alert"},
		},
		{
			name:     "style with expression",
			input:    `<div style="background:url('javascript:alert(1)')">`,
			contains: []string{"javascript:", "alert"},
		},
	}

	for _, tt := range maliciousInputs {
		t.Run(tt.name, func(t *testing.T) {
			schema := validation.Make().Shape(map[string]validation.Type{
				"content": validation.AdvancedString().
					StripTags().
					EscapeHTML().
					Label("User Content"),
			})

			result := schema.Validate(map[string]any{
				"content": tt.input,
			})

			if result.HasErrors() {
				t.Errorf("sanitization should not produce errors: %v", result.Errors())
			}

			validData := result.ValidData()
			sanitized := validData["content"].(string)

			// Verify malicious content is removed/escaped
			for _, dangerous := range tt.contains {
				if strings.Contains(sanitized, dangerous) {
					t.Errorf("sanitized output still contains dangerous content: '%s' in '%s'", dangerous, sanitized)
				}
			}
		})
	}
}

// TestAdvancedString_RealWorld_BlogPost tests blog post sanitization
func TestAdvancedString_RealWorld_BlogPost(t *testing.T) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"title": validation.AdvancedString().
			Required().
			Trim().
			Min(10).Max(200).
			StripTags().
			Label("Title"),

		"slug": validation.AdvancedString().
			Required().
			Trim().
			ReplaceTurkishChars().
			StripPunctuation().
			Label("URL Slug"),

		"content": validation.AdvancedString().
			Required().
			Trim().
			StripTags("<p>", "<br>", "<b>", "<i>", "<u>", "<a>").
			MaxWords(500).
			Label("Content"),

		"excerpt": validation.AdvancedString().
			Trim().
			StripTags().
			MaxWords(50).
			Label("Excerpt"),
	})

	t.Run("valid blog post", func(t *testing.T) {
		data := map[string]any{
			"title":   "  How to Build a Go Application  ",
			"slug":    "Nasıl Bir Go Uygulaması Oluşturulur?",
			"content": "<p>This is a <b>great</b> article about Go.</p>",
			"excerpt": "Learn about Go programming in this comprehensive guide.",
		}

		result := schema.Validate(data)
		if result.HasErrors() {
			t.Errorf("valid post should not have errors: %v", result.Errors())
		}

		validData := result.ValidData()

		// Verify trimming
		title := validData["title"].(string)
		if title != "How to Build a Go Application" {
			t.Errorf("title not trimmed correctly: '%s'", title)
		}

		// Verify Turkish chars replaced in slug
		slug := validData["slug"].(string)
		if strings.Contains(slug, "ı") || strings.Contains(slug, "ş") {
			t.Errorf("slug still contains Turkish chars: '%s'", slug)
		}
	})

	t.Run("XSS attack attempt", func(t *testing.T) {
		data := map[string]any{
			"title":   "<script>alert('XSS')</script>Hacked Title",
			"slug":    "hacked-post",
			"content": `<p>Normal content</p><script>stealCookies()</script>`,
			"excerpt": "Safe excerpt",
		}

		result := schema.Validate(data)
		if result.HasErrors() {
			t.Logf("XSS validation errors (expected): %v", result.Errors())
		}

		validData := result.ValidData()
		content := validData["content"].(string)

		// Verify script tags removed
		if strings.Contains(content, "<script>") {
			t.Error("script tags not removed from content")
		}
	})

	t.Run("content too long", func(t *testing.T) {
		// Generate 501 words
		longContent := strings.Repeat("word ", 501)

		data := map[string]any{
			"title":   "Valid Title Here",
			"slug":    "valid-slug",
			"content": longContent,
			"excerpt": "Short excerpt",
		}

		result := schema.Validate(data)
		if !result.HasErrors() {
			t.Error("content with 501 words should fail MaxWords(500)")
		}
	})
}

// BenchmarkAdvancedString_XSSSanitization benchmarks XSS sanitization
func BenchmarkAdvancedString_XSSSanitization(b *testing.B) {
	schema := validation.Make().Shape(map[string]validation.Type{
		"content": validation.AdvancedString().
			StripTags().
			EscapeHTML(),
	})

	data := map[string]any{
		"content": `<script>alert('XSS')</script><p>Hello <b>World</b></p>`,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schema.Validate(data)
	}
}
