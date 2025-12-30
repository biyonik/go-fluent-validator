# Advanced String (Sanitization)

`AdvancedString()` provides powerful sanitization and transformation for handling user input safely. Use this when you need XSS protection, HTML filtering, or content cleanup.

## When to Use

| Scenario | Use |
|----------|-----|
| Form input, user-generated content | `AdvancedString()` |
| Filenames from uploads | `AdvancedString()` |
| HTML content that needs filtering | `AdvancedString()` |
| Simple validation without sanitization | `String()` |

## Basic Sanitization

```go
schema := v.Make().Shape(map[string]v.Type{
    "bio": v.AdvancedString().
        Required().
        Trim().                            // Remove whitespace
        StripTags("<b>", "<i>", "<u>").    // Keep only these HTML tags
        EscapeHTML().                      // Escape remaining HTML entities
        MaxWords(100).                     // Limit word count
        Label("Biography"),
})
```

## XSS Protection

```go
schema := v.Make().Shape(map[string]v.Type{
    // Remove all HTML and escape special characters
    "comment": v.AdvancedString().
        StripTags().       // Remove ALL HTML tags
        EscapeHTML().      // Convert & < > " ' to entities
        Label("Comment"),

    // Allow some safe tags
    "content": v.AdvancedString().
        StripTags("<p>", "<br>", "<b>", "<i>", "<a>").
        EscapeHTML().
        Label("Content"),
})
```

## Filename Sanitization

```go
schema := v.Make().Shape(map[string]v.Type{
    "filename": v.AdvancedString().
        Required().
        SanitizeFilename().    // Remove dangerous chars (../, null bytes, etc.)
        StripPunctuation().    // Remove punctuation
        Label("File Name"),
})
```

## Character Filtering

```go
schema := v.Make().Shape(map[string]v.Type{
    // Only allow alphanumeric
    "slug": v.AdvancedString().
        Trim().
        CharSet("alphanumeric").
        Label("URL Slug"),

    // Remove emojis
    "text": v.AdvancedString().
        FilterEmoji(true).     // true = remove emojis
        Label("Text"),

    // Convert Turkish characters to ASCII
    "searchKey": v.AdvancedString().
        ReplaceTurkishChars(). // ş→s, ğ→g, ı→i, etc.
        Label("Search Key"),
})
```

## Domain Validation

```go
schema := v.Make().Shape(map[string]v.Type{
    "domain": v.AdvancedString().
        ValidateDomain().      // Check valid domain name format
        Label("Domain"),
})
```

## Complete Example

```go
// Blog post content with comprehensive sanitization
schema := v.Make().Shape(map[string]v.Type{
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
        Label("URL Slug"),

    "content": v.AdvancedString().
        Required().
        Trim().
        StripTags("<p>", "<br>", "<b>", "<i>", "<u>", "<a>", "<ul>", "<li>").
        Min(100).
        Label("Content"),

    "excerpt": v.AdvancedString().
        Trim().
        StripTags().
        MaxWords(50).
        Label("Excerpt"),
})
```

## API Reference

| Method | Description | Example |
|--------|-------------|---------|
| `.Trim()` | Remove leading/trailing whitespace | `.Trim()` |
| `.StripTags(allowed...)` | Remove HTML tags, optionally keep some | `.StripTags("<b>", "<i>")` |
| `.EscapeHTML()` | Escape HTML entities | `.EscapeHTML()` |
| `.FilterEmoji(remove)` | Remove or keep emojis | `.FilterEmoji(true)` |
| `.SanitizeFilename()` | Make safe for file systems | `.SanitizeFilename()` |
| `.CharSet(type)` | Allow only specific characters | `.CharSet("alphanumeric")` |
| `.MaxWords(n)` | Limit word count | `.MaxWords(50)` |
| `.StripPunctuation()` | Remove punctuation marks | `.StripPunctuation()` |
| `.ReplaceTurkishChars()` | Convert Turkish → ASCII | `.ReplaceTurkishChars()` |
| `.ValidateDomain()` | Validate domain name format | `.ValidateDomain()` |

All methods from [String Validation](./string-validation.md) are also available (`.Required()`, `.Min()`, `.Max()`, etc.).
