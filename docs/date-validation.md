# Date Validation

The `Date()` validator parses and validates dates with flexible format support and range constraints.

## Basic Usage

```go
schema := v.Make().Shape(map[string]v.Type{
    // Parse with custom format
    "birthDate": v.Date().
        Format("2006-01-02").     // Go time format
        Required().
        Before(time.Now()).       // Must be in the past
        Label("Birth Date"),

    // Default format (RFC3339)
    "timestamp": v.Date().
        Required().
        Label("Timestamp"),
})
```

## Date Ranges

```go
schema := v.Make().Shape(map[string]v.Type{
    // Min and max dates
    "startDate": v.Date().
        Format("2006-01-02").
        Min(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)).
        Max(time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)).
        Label("Start Date"),

    // Before specific date (past)
    "birthDate": v.Date().
        Format("2006-01-02").
        Before(time.Now()).
        Label("Birth Date"),

    // After specific date (future)
    "eventDate": v.Date().
        Format("2006-01-02").
        After(time.Now()).
        Label("Event Date"),
})
```

## Date with Time

```go
schema := v.Make().Shape(map[string]v.Type{
    "appointment": v.Date().
        Format("2006-01-02 15:04:05").
        Required().
        After(time.Now()).
        Label("Appointment"),
})
```

## Go Time Format Reference

Go uses a reference time for format strings: `Mon Jan 2 15:04:05 MST 2006`

| Pattern | Meaning | Example |
|---------|---------|---------|
| `2006` | Year (4 digits) | 2024 |
| `06` | Year (2 digits) | 24 |
| `01` | Month (2 digits) | 03 |
| `1` | Month (no padding) | 3 |
| `Jan` | Month (short) | Mar |
| `02` | Day (2 digits) | 09 |
| `2` | Day (no padding) | 9 |
| `15` | Hour (24h) | 14 |
| `03` | Hour (12h) | 02 |
| `04` | Minute | 30 |
| `05` | Second | 45 |

Common formats:
- `2006-01-02` → 2024-03-15
- `02/01/2006` → 15/03/2024
- `2006-01-02 15:04:05` → 2024-03-15 14:30:45

## API Reference

| Method | Description | Example |
|--------|-------------|---------|
| `.Format(layout)` | Go time format layout | `.Format("2006-01-02")` |
| `.Required()` | Field must be present | `.Required()` |
| `.Min(date)` | Minimum date | `.Min(time.Now())` |
| `.Max(date)` | Maximum date | `.Max(deadline)` |
| `.Before(date)` | Must be before date | `.Before(time.Now())` |
| `.After(date)` | Must be after date | `.After(startDate)` |
| `.Label(name)` | Custom error label | `.Label("Due Date")` |
