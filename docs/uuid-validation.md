# UUID Validation

The `Uuid()` validator handles UUID validation with version-specific checks.

## Basic Usage

```go
schema := v.Make().Shape(map[string]v.Type{
    // Any UUID version
    "id": v.Uuid().Required().Label("ID"),

    // Specific version
    "sessionId": v.Uuid().Version(4).Required().Label("Session ID"),
})
```

## Version-Specific Validation

```go
schema := v.Make().Shape(map[string]v.Type{
    // UUID v1 (timestamp-based)
    "eventId": v.Uuid().Version(1).Label("Event ID"),

    // UUID v4 (random)
    "userId": v.Uuid().Version(4).Label("User ID"),

    // UUID v5 (namespace-based)
    "resourceId": v.Uuid().Version(5).Label("Resource ID"),
})
```

## Supported Versions

| Version | Type | Use Case |
|---------|------|----------|
| 1 | Time-based | When creation time matters |
| 2 | DCE Security | Rarely used |
| 3 | MD5 hash | Reproducible IDs from names |
| 4 | Random | Most common, general purpose |
| 5 | SHA-1 hash | Reproducible IDs from names |

## API Reference

| Method | Description | Example |
|--------|-------------|---------|
| `.Required()` | Field must be present | `.Required()` |
| `.Version(n)` | Specific UUID version (1-5) | `.Version(4)` |
| `.Label(name)` | Custom error label | `.Label("User ID")` |
