# is

**[中文文档](README.zh-CN.md)** | English

A tiny, fast, dependency-free validation/detection toolkit for Go.

This package focuses on simple predicate-style helpers (returning bool) for common checks:

- Emails, IP/MAC addresses, phone numbers
- Numbers vs. numeric strings, booleans, case checks
- UUID/ULID, hashes, colors, encodings (HTML/URL/Base64)
- JWT tokens, semantic versions, labels
- Time formats and time zones
- Length comparisons, comparisons between numbers/strings/times
- File/dir existence

Module path: `go-slim.dev/is`

## Install

```bash
go get go-slim.dev/is
```

## Quick Start

```go
package main

import (
    "fmt"
    is "go-slim.dev/is"
)

func main() {
    fmt.Println(is.Email("user@example.com")) // true
    fmt.Println(is.PhoneNumber("13800138000")) // true (Chinese mobile)

    // Number vs Numeric
    fmt.Println(is.Number("123"))  // true (digits only)
    fmt.Println(is.Number("12.3")) // false (has decimal)
    fmt.Println(is.Numeric("12.3")) // true (numeric string allowed)

    // Boolean detection (strings/ints)
    fmt.Println(is.Boolean("yes")) // true
    fmt.Println(is.Boolean(1))      // true

    // IP, MAC
    fmt.Println(is.IP("127.0.0.1")) // true
    fmt.Println(is.IPv6("::1"))     // true
    fmt.Println(is.MAC("01:23:45:67:89:ab")) // true

    // Colors
    fmt.Println(is.Color("#fff"))        // true
    fmt.Println(is.RGBA("rgba(1,2,3,0.5)")) // true

    // Additional validations
    fmt.Println(is.JWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")) // true
    fmt.Println(is.Semver("1.2.3")) // true

    // Time & timezone
    fmt.Println(is.Datetime("2023-01-02", "2006-01-02")) // true
    fmt.Println(is.Timezone("UTC"))                        // true

    // Length helpers (string length by runes; numeric width by digits)
    fmt.Println(is.Length("abc", 3, "="))            // true
    fmt.Println(is.LengthBetween([]int{1,2,3}, 1, 3))  // true

    // Compare and Between/NotBetween
    fmt.Println(is.Compare(2, 1, ">")) // true
}
```

## API Highlights

- Strings and bytes
  - `Email(s)`, `Base64(s)`, `Base64URL(s)`, `URLEncoded(s)`, `HTMLEncoded(s)`
- Identifiers and hashes
  - `UUID(s)`, `UUID3/4/5(s)`, `ULID(s)`, `MD5/SHA256/SHA384/SHA512(s)`, `JWT(s)`, `Semver(s)`, `Label(s)`
- Numbers & booleans
  - `Number(v)` digits-only for strings; numbers for numeric types
  - `Numeric(v)` accepts numeric strings with decimals and numeric types
  - `Boolean(v)` accepts 1/0, yes/no, on/off, true/false (case-insensitive) and numeric 0/1
- IP, MAC & Phone
  - `IPv4(s)`, `IPv6(s)`, `IP(s)`, `MAC(s)`, `PhoneNumber(s)`, `E164(s)`
- Colors
  - `HEXColor(s)`, `RGB(s)`, `RGBA(s)`, `HSL(s)`, `HSLA(s)`, `Color(s)`
- Time
  - `Datetime(s, layout)`, `Timezone(name)`
- Length & comparison
  - `Length(v, n, op)`, `LengthBetween(v, min, max)`
  - `Compare(a, b, op)`, `GreaterThan/Equal/LessThan`, `Between/NotBetween`
- Filesystem
  - `File(path)`, `Dir(path)`

Notes:

- `Number()` vs `Numeric()` are intentionally different (digits-only vs numeric including decimals).
- `URLEncoded()` detector allows spaces or valid %XX sequences per current regex.
- `Between()`/`NotBetween()` guard semantics are covered by tests; prefer `Compare()` for direct relational checks.

## Examples

See tests for live examples:

- `is_test.go` for core API coverage
- `boundary_test.go` for edge cases

## Performance

Benchmarks (Apple M4 Pro, Go 1.24; indicative only):

```
BenchmarkCompare_Number-14   ~5.5 ns/op   0 B/op   0 allocs/op
BenchmarkIPv4-14             ~17 ns/op    0 B/op   0 allocs/op
BenchmarkURL-14              ~168 ns/op   192 B/op 2 allocs/op
```

Run locally:

```bash
go test -bench . -benchmem -run=^$
```

## Testing

```bash
go test -v ./...
```

## Versioning & Compatibility

- Go 1.20+ recommended (tested with Go 1.24)
- Public API aims to be stable; any breaking changes will bump the major version

## Contributing

Issues and PRs welcome. Please include:

- Focused changes, clear rationale
- Tests for fixes/features
- Benchmarks where performance is relevant

## License

MIT
