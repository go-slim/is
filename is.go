// Package is provides a comprehensive set of validation and type checking utilities for Go applications.
//
// The package offers simple, readable validation functions for common data types and patterns:
// - Email addresses, phone numbers, URLs
// - UUIDs, hashes, and identifiers
// - Colors, coordinates, and geographic data
// - File paths and directory checks
// - Type conversions and comparisons
// - Length and range validations
//
// Most functions follow a simple pattern: `Is.TypeName(value) bool` where TypeName represents
// the validation being performed. The package also provides generic versions of some functions
// for better type safety.
//
// Examples:
//
//	is.Email("user@example.com")        // true
//	is.PhoneNumber("13800138000")      // true
//	is.UUID("550e8400-e29b-41d4-a716-446655440000") // true
//	is.Between(25, 18, 65)              // true
//
// The package is designed to be used in form validation, API input validation,
// configuration validation, and general data verification scenarios.
package is

import (
	"encoding/json"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Email validates whether the given string is a valid email address.
// It supports international email addresses and follows RFC 5322 standards.
//
// Examples:
//
//	is.Email("user@example.com")     // true
//	is.Email("user+tag@domain.org")  // true
//	is.Email("invalid-email")        // false
func Email(str string) bool {
	return emailRegex.MatchString(str)
}

// E164 validates whether the given string conforms to the E.164 international
// telephone numbering plan format. E.164 numbers are in the format +[country code][subscriber number].
//
// Examples:
//
//	is.E164("+14155552671")         // true (US number)
//	is.E164("+8613800138000")       // true (China number)
//	is.E164("14155552671")          // false (missing +)
//	is.E164("+123")                 // false (too short)
func E164(str string) bool {
	return e164Regex.MatchString(str)
}

// PhoneNumber validates whether the given string is a valid Chinese mainland mobile phone number.
// It supports both domestic format (11 digits starting with 1) and international format (+86 prefix).
//
// Examples:
//
//	is.PhoneNumber("13800138000")   // true
//	is.PhoneNumber("+8613800138000") // true
//	is.PhoneNumber("12345678901")   // false (invalid pattern)
//	is.PhoneNumber("2800138000")    // false (invalid prefix)
func PhoneNumber(s string) bool {
	return phoneNumberRegex.MatchString(s)
}

// Semver validates whether the given string conforms to Semantic Versioning 2.0.0 specification.
// Semantic versions are in the format MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD].
//
// Examples:
//
//	is.Semver("1.0.0")              // true
//	is.Semver("2.1.3-beta.2+exp.sha.5114f85") // true
//	is.Semver("1.0")                // false (missing PATCH)
//	is.Semver("v1.0.0")             // false (invalid prefix)
func Semver(s string) bool {
	return semverRegex.MatchString(s)
}

// Label validates whether the given string follows variable naming conventions.
// Labels must start with a letter (A-F) followed by word characters, digits, or underscores.
// This is commonly used for environment variable names, configuration keys, etc.
//
// Examples:
//
//	is.Label("API_HOST")            // true
//	is.Label("Fconfig_debug")        // true
//	is.Label("1INVALID")            // false (starts with digit)
//	is.Label("api-host")            // false (contains hyphen)
func Label(s string) bool {
	return labelRegex.MatchString(s)
}

// Base64 validates whether the given string is valid Base64 encoded data.
// It supports standard Base64 encoding with proper padding.
//
// Examples:
//
//	is.Base64("SGVsbG8gV29ybGQ=")    // true ("Hello World")
//	is.Base64("YW55IGNhcm5hbCBwbGF5")  // true ("any carnal play")
//	is.Base64("Hello World")        // false (not encoded)
//	is.Base64("SGVsbG8=")           // false (invalid length)
func Base64(s string) bool {
	return base64Regex.MatchString(s)
}

// URL validates whether the given string is a valid URL.
// It supports various URL schemes and follows standard URL parsing rules.
// The function strips URL fragments (#) before validation to handle common use cases.
//
// Examples:
//
//	is.URL("https://example.com")       // true
//	is.URL("ftp://user:pass@host/path") // true
//	is.URL("http://localhost:8080")     // true
//	is.URL("example.com")              // false (missing scheme)
//	is.URL("")                         // false (empty string)
func URL(s string) bool {
	var i int
	// checks needed as of Go 1.6 because of change https://github.com/golang/go/commit/617c93ce740c3c3cc28cdd1a0d712be183d0b328#diff-6c2d018290e298803c0c9419d8739885L195
	// emulate browser and strip the '#' suffix prior to validation. see issue-#237
	if i = strings.Index(s, "#"); i > -1 {
		s = s[:i]
	}
	if len(s) == 0 {
		return false
	}
	u, err := url.ParseRequestURI(s)
	if err != nil || u.Scheme == "" {
		return false
	}
	return true
}

// Base64URL validates whether the given string is a valid URL-safe Base64 encoded string.
// It uses the URL-safe Base64 alphabet (A-Z, a-z, 0-9, -, _) instead of the standard Base64 alphabet.
// This format is commonly used in URLs and JWT tokens where '+' and '/' characters would need to be URL-encoded.
//
// Examples:
//
//	is.Base64URL("SGVsbG8gV29ybGQ")      // true ("Hello World")
//	is.Base64URL("eyJhbGciOiJIUzI1NiIs") // true (JWT header)
//	is.Base64URL("Hello World")        // false (not encoded)
//	is.Base64URL("SGVsbG8=")           // false (invalid length)
func Base64URL(str string) bool {
	return base64URLRegex.MatchString(str)
}

// JWT validates whether the given string is a valid JSON Web Token (JWT).
// JWTs have three parts separated by dots: header.payload.signature
// The header and payload are Base64URL encoded, while the signature is optional.
//
// Examples:
//
//	is.JWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkxNjQ1ODkyODkxNjk1NzU5IiwibmFtZSI6IkpvaG4ifQ") // true
//	is.JWT("invalid.jwt")              // false (missing dots)
//	is.JWT("A.B")                   // false (only two parts)
func JWT(str string) bool {
	return jWTRegex.MatchString(str)
}

// UUID5 validates whether the given string is a valid UUID version 5.
// UUIDv5 uses namespace and name hashing with SHA-1 to generate deterministic UUIDs.
// The version field in the UUID must be 5 (binary 0101).
//
// Examples:
//
//	is.UUID5("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // true
//	is.UUID5("550e8400-e29b-41d4-a716-446655440000") // false (v1 UUID)
//	is.UUID5("invalid-uuid")                // false
func UUID5(str string) bool {
	return uUID5Regex.MatchString(str)
}

// UUID4 validates whether the given string is a valid UUID version 4.
// UUIDv4 uses random or pseudo-random numbers to generate UUIDs.
// The version field in the UUID must be 4 (binary 0100) and the variant field must be 8, 9, A, or B.
//
// Examples:
//
//	is.UUID4("550e8400-e29b-41d4-a716-446655440000") // true
//	is.UUID4("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // false (v3 UUID)
//	is.UUID4("invalid-uuid")                // false
func UUID4(str string) bool {
	return uUID4Regex.MatchString(str)
}

// UUID3 validates whether the given string is a valid UUID version 3.
// UUIDv3 uses namespace and name hashing with MD5 to generate deterministic UUIDs.
// The version field in the UUID must be 3 (binary 0011) and the variant field must be 8 or B.
//
// Examples:
//
//	is.UUID3("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // true
//	is.UUID3("550e8400-e29b-41d4-a716-446655440000") // false (v1 UUID)
//	is.UUID3("invalid-uuid")                // false
func UUID3(str string) bool {
	return uUID3Regex.MatchString(str)
}

// UUID validates whether the given string is a valid UUID of any version (v1, v3, v4, or v5).
// UUIDs are 128-bit unique identifiers typically displayed as 32 hexadecimal digits
// with hyphens in the pattern 8-4-4-4-12.
//
// Examples:
//
//	is.UUID("550e8400-e29b-41d4-a716-446655440000") // true
//	is.UUID("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // true
//	is.UUID("550e8400-e29b-41d4-a716")          // false (too short)
//	is.UUID("not-a-uuid")                        // false
func UUID(str string) bool {
	return uUIDRegex.MatchString(str)
}

// ULID validates whether the given string is a valid ULID (Universally Unique Lexicographically Sortable Identifier).
// ULIDs are 26-character strings using Crockford's Base32 alphabet, designed to be sortable by time.
//
// Examples:
//
//	is.ULID("01ARZ3NDEKTSV4RRFFQ69G5FA")      // true
//	is.ULID("01BX5ZZKBKACTAV9WEVGEMMVR")      // true
//	is.ULID("invalid-ulid")                    // false
//	is.ULID("01ARZ3NDEKTSV4RRFFQ69G5FA0")     // false (wrong length)
func ULID(str string) bool {
	return uLIDRegex.MatchString(str)
}

// MD4 validates whether the given string is a valid MD4 hash digest.
// MD4 is a cryptographic hash function that produces a 128-bit (32-character) hash value.
// Note: MD4 is considered insecure for cryptographic purposes due to known vulnerabilities.
//
// Examples:
//
//	is.MD4("31d6cfe0d16ae931b73c59d7e0c089c0b") // true
//	is.MD4("invalid-hash")                 // false
func MD4(str string) bool {
	return md4Regex.MatchString(str)
}

// MD5 validates whether the given string is a valid MD5 hash digest.
// MD5 produces a 128-bit (32-character) hash value, commonly used for file integrity checks.
// Note: MD5 is considered insecure for cryptographic purposes due to collision vulnerabilities.
//
// Examples:
//
//	is.MD5("5d41402abc4b2a76b9719d911017c592") // true ("hello")
//	is.MD5("invalid-hash")                  // false
func MD5(str string) bool {
	return md5Regex.MatchString(str)
}

// SHA256 validates whether the given string is a valid SHA-256 hash digest.
// SHA-256 produces a 256-bit (64-character) hash value, part of the SHA-2 family.
// It's widely used in blockchain, digital signatures, and data integrity verification.
//
// Examples:
//
//	is.SHA256("2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824") // true ("hello")
//	is.SHA256("invalid-hash")                     // false
func SHA256(str string) bool {
	return sha256Regex.MatchString(str)
}

// SHA384 validates whether the given string is a valid SHA-384 hash digest.
// SHA-384 produces a 384-bit (96-character) hash value, part of the SHA-2 family.
// It's commonly used in applications requiring higher security than SHA-256.
//
// Examples:
//
//	is.SHA384("59e1742e6e2bc8c9a5cb46c6e9a7a3d551cedd73a4b17c4a9a4661d6f8a") // true ("hello")
//	is.SHA384("invalid-hash")                     // false
func SHA384(str string) bool {
	return sha384Regex.MatchString(str)
}

// SHA512 validates whether the given string is a valid SHA-512 hash digest.
// SHA-512 produces a 512-bit (128-character) hash value, part of the SHA-2 family.
// It's the strongest hash function in the SHA-2 family and is recommended for new applications.
//
// Examples:
//
//	is.SHA512("9b71d46469a7a33e2c9a5b5d9f2f436c14b862f3498e7b2a631f8aeeb8b9d7")
//	is.SHA512("invalid-hash")                     // false
func SHA512(str string) bool {
	return sha512Regex.MatchString(str)
}

// ASCII validates whether the given string contains only ASCII characters (0-127).
// This function checks that all characters in the string are within the standard ASCII range,
// which includes control characters, digits, letters, and basic punctuation.
//
// Examples:
//
//	is.ASCII("Hello World")         // true
//	is.ASCII("12345")               // true
//	is.ASCII("Héllo")               // false (contains é)
//	is.ASCII("你好")                // false (Chinese characters)
//	is.ASCII("")                    // true (empty string)
func ASCII(str string) bool {
	return aSCIIRegex.MatchString(str)
}

// Alpha validates whether the given string contains only alphabetic characters (A-Z, a-z).
// This function checks that all characters are letters from the English alphabet,
// without any numbers, spaces, or special characters.
//
// Examples:
//
//	is.Alpha("Hello")               // true
//	is.Alpha("WORLD")               // true
//	is.Alpha("HelloWorld")          // true
//	is.Alpha("Hello World")         // false (contains space)
//	is.Alpha("Hello123")            // false (contains numbers)
//	is.Alpha("")                    // true (empty string)
func Alpha(str string) bool {
	return alphaRegex.MatchString(str)
}

// Alphanumeric validates whether the given string contains only letters and numbers (A-Z, a-z, 0-9).
// This function checks that all characters are alphanumeric characters from the English alphabet,
// without any spaces, punctuation, or special characters.
//
// Examples:
//
//	is.Alphanumeric("Hello123")        // true
//	is.Alphanumeric("Test2023")         // true
//	is.Alphanumeric("abc")              // true
//	is.Alphanumeric("123")              // true
//	is.Alphanumeric("Hello World")      // false (contains space)
//	is.Alphanumeric("Hello-123")        // false (contains hyphen)
//	is.Alphanumeric("")                 // true (empty string)
func Alphanumeric(str string) bool {
	return alphaNumericRegex.MatchString(str)
}

// AlphaUnicode validates whether the given string contains only Unicode alphabetic characters.
// Unlike Alpha() which only supports English letters (A-Z, a-z), this function supports
// alphabetic characters from all languages including accented letters and non-Latin scripts.
//
// Examples:
//
//	is.AlphaUnicode("Hello")           // true
//	is.AlphaUnicode("Héllo")           // true (contains é)
//	is.AlphaUnicode("Привет")          // true (Cyrillic)
//	is.AlphaUnicode("你好")            // true (Chinese)
//	is.AlphaUnicode("مرحبا")           // true (Arabic)
//	is.AlphaUnicode("Hello123")        // false (contains numbers)
//	is.AlphaUnicode("Hello World")     // false (contains space)
func AlphaUnicode(str string) bool {
	return alphaUnicodeRegex.MatchString(str)
}

// AlphanumericUnicode validates whether the given string contains only Unicode letters and numbers.
// Unlike Alphanumeric() which only supports English letters (A-Z, a-z), this function supports
// alphabetic characters and numbers from all languages including accented letters and non-Latin scripts.
//
// Examples:
//
//	is.AlphanumericUnicode("Hello123")     // true
//	is.AlphanumericUnicode("Héllo123")     // true (contains é)
//	is.AlphanumericUnicode("Привет123")    // true (Cyrillic with numbers)
//	is.AlphanumericUnicode("你好123")      // true (Chinese with numbers)
//	is.AlphanumericUnicode("مرحبا123")     // true (Arabic with numbers)
//	is.AlphanumericUnicode("Hello 123")    // false (contains space)
//	is.AlphanumericUnicode("Hello-123")    // false (contains hyphen)
func AlphanumericUnicode(str string) bool {
	return alphaUnicodeNumericRegex.MatchString(str)
}

// Numeric is the validation function for validating if the current field's value is a valid numeric value.
func Numeric[T any](t T) bool {
	ctx := reflect.ValueOf(t)
	switch ctx.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return numericRegex.MatchString(ctx.String())
	}
}

// Number is the validation function for validating if the current field's value is a valid number.
func Number[T any](t T) bool {
	rv := reflect.ValueOf(t)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return numberRegex.MatchString(rv.String())
	}
}

// Boolean is the validation function for validating if the current field's value can be safely converted to a boolean.
func Boolean[T any](t T) bool {
	ref := reflect.ValueOf(t)
	switch ref.Kind() {
	case reflect.String:
		switch ref.String() {
		case "1", "yes", "YES", "Yes", "on", "ON", "On", "true", "TRUE", "True",
			"0", "no", "NO", "No", "", "off", "OFF", "Off", "false", "FALSE", "False":
			return true
		default:
			return false
		}
	case reflect.Int, reflect.Int32, reflect.Int64:
		n := ref.Int()
		return n == 0 || n == 1
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		n := ref.Uint()
		return n == 0 || n == 1
	case reflect.Bool:
		return ref.Bool()
	default:
		return false
	}
}

// Default is the opposite of required aka HasValue
func Default(val any) bool {
	return !HasValue(val)
}

// HasValue is the validation function for validating if the current field's value is not the default static value.
func HasValue(val any) bool {
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !rv.IsNil()
	default:
		return rv.IsValid() && rv.Interface() != reflect.Zero(rv.Type()).Interface()
	}
}

// Hexadecimal validates whether the given string contains only hexadecimal characters (0-9, A-F, a-f).
// This function is useful for validating hex values like color codes, binary data encoded as hex,
// or any other hexadecimal representation.
//
// Examples:
//
//	is.Hexadecimal("FF00AA")        // true
//	is.Hexadecimal("123abc")         // true
//	is.Hexadecimal("0xFF")           // false (contains 0x prefix)
//	is.Hexadecimal("FFGG")           // false (G is not hex)
//	is.Hexadecimal("hello")          // false (contains non-hex letters)
//	is.Hexadecimal("")               // true (empty string)
func Hexadecimal(str string) bool {
	return hexadecimalRegex.MatchString(str)
}

// HEXColor validates whether the given string is a valid HEX color code.
// Supports both 3-digit (#RGB), 6-digit (#RRGGBB), and 8-digit (#RRGGBBAA) hexadecimal color formats.
// The 8-digit format includes alpha/transparency channel.
//
// Examples:
//
//	is.HEXColor("#FF0000")          // true (red)
//	is.HEXColor("#00FF00")          // true (green)
//	is.HEXColor("#0000FF")          // true (blue)
//	is.HEXColor("#F00")             // true (short red)
//	is.HEXColor("#FF000080")        // true (red with 50% alpha)
//	is.HEXColor("FF0000")           // false (missing #)
//	is.HEXColor("#GGGGGG")          // false (invalid hex)
func HEXColor(str string) bool {
	return hexColorRegex.MatchString(str)
}

// RGB validates whether the given string is a valid RGB color specification.
// RGB colors use the format rgb(red, green, blue) where each value is 0-255.
// The function accepts values with or without spaces between components.
//
// Examples:
//
//	is.RGB("rgb(255, 0, 0)")        // true (red)
//	is.RGB("rgb(0, 255, 0)")        // true (green)
//	is.RGB("rgb(0, 0, 255)")        // true (blue)
//	is.RGB("rgb(255,255,255)")      // true (white, no spaces)
//	is.RGB("rgb(128, 128, 128)")    // true (gray)
//	is.RGB("rgb(256, 0, 0)")        // false (value > 255)
//	is.RGB("rgb(-1, 0, 0)")         // false (negative value)
//	is.RGB("rgba(255, 0, 0)")       // false (wrong format)
func RGB(str string) bool {
	return rgbRegex.MatchString(str)
}

// RGBA validates whether the given string is a valid RGBA color specification.
// RGBA colors use the format rgba(red, green, blue, alpha) where RGB values are 0-255
// and alpha is 0.0-1.0 (0 = fully transparent, 1 = fully opaque).
// The function accepts values with or without spaces between components.
//
// Examples:
//
//	is.RGBA("rgba(255, 0, 0, 0.5)")     // true (red with 50% opacity)
//	is.RGBA("rgba(0, 255, 0, 1.0)")     // true (green, fully opaque)
//	is.RGBA("rgba(0, 0, 255, 0)")       // true (blue, fully transparent)
//	is.RGBA("rgba(255,255,255,0.8)")    // true (white with 80% opacity)
//	is.RGBA("rgba(255, 0, 0, 1.5)")     // false (alpha > 1.0)
//	is.RGBA("rgba(255, 0, 0, -0.1)")    // false (negative alpha)
//	is.RGBA("rgb(255, 0, 0)")           // false (missing alpha)
func RGBA(str string) bool {
	return rgbaRegex.MatchString(str)
}

// HSL validates whether the given string is a valid HSL color specification.
// HSL colors use the format hsl(hue, saturation, lightness) where:
// - hue: 0-360 degrees on the color wheel
// - saturation: 0-100% (0 = grayscale, 100 = fully saturated)
// - lightness: 0-100% (0 = black, 100 = white)
// The function accepts values with or without spaces between components.
//
// Examples:
//
//	is.HSL("hsl(0, 100%, 50%)")        // true (red)
//	is.HSL("hsl(120, 100%, 50%)")      // true (green)
//	is.HSL("hsl(240, 100%, 50%)")      // true (blue)
//	is.HSL("hsl(0,100%,50%)")          // true (no spaces)
//	is.HSL("hsl(180, 50%, 75%)")       // true (cyan, less saturated)
//	is.HSL("hsl(361, 100%, 50%)")      // false (hue > 360)
//	is.HSL("hsl(0, 101%, 50%)")        // false (saturation > 100%)
//	is.HSL("hsla(0, 100%, 50%)")       // false (wrong format)
func HSL(str string) bool {
	return hslRegex.MatchString(str)
}

// HSLA validates whether the given string is a valid HSLA color specification.
// HSLA colors use the format hsla(hue, saturation, lightness, alpha) where:
// - hue: 0-360 degrees on the color wheel
// - saturation: 0-100% (0 = grayscale, 100 = fully saturated)
// - lightness: 0-100% (0 = black, 100 = white)
// - alpha: 0.0-1.0 (0 = fully transparent, 1 = fully opaque)
// The function accepts values with or without spaces between components.
//
// Examples:
//
//	is.HSLA("hsla(0, 100%, 50%, 0.5)")   // true (red with 50% opacity)
//	is.HSLA("hsla(120, 100%, 50%, 1.0)")  // true (green, fully opaque)
//	is.HSLA("hsla(240, 100%, 50%, 0)")    // true (blue, fully transparent)
//	is.HSLA("hsla(0,100%,50%,0.8)")       // true (no spaces, 80% opacity)
//	is.HSLA("hsla(0, 100%, 50%, 1.5)")    // false (alpha > 1.0)
//	is.HSLA("hsla(0, 100%, 50%, -0.1)")   // false (negative alpha)
//	is.HSLA("hsl(0, 100%, 50%)")          // false (missing alpha)
func HSLA(str string) bool {
	return hslaRegex.MatchString(str)
}

// Color validates whether the given string represents a valid color value.
// It supports multiple color formats: HEX colors, RGB, RGBA, HSL, and HSLA.
//
// Examples:
//
//	is.Color("#FF0000")                 // true (HEX red)
//	is.Color("rgb(255, 0, 0)")         // true (RGB red)
//	is.Color("rgba(255, 0, 0, 0.5)")    // true (RGBA red with 50% opacity)
//	is.Color("hsl(0, 100%, 50%)")       // true (HSL red)
//	is.Color("hsla(0, 100%, 50%, 0.5)")  // true (HSLA red with 50% opacity)
//	is.Color("red")                     // false (color name not supported)
func Color(str string) bool {
	return HEXColor(str) || HSLA(str) || HSL(str) || RGB(str) || RGBA(str)
}

// Latitude validates whether the given value represents a valid geographic latitude coordinate.
// Valid latitude ranges from -90° (South Pole) to +90° (North Pole).
// The function accepts various numeric types and string representations.
//
// Examples:
//
//	is.Latitude(45.0)                  // true
//	is.Latitude(-90.0)                 // true (South Pole)
//	is.Latitude("90.0")                // true (North Pole)
//	is.Latitude(91.0)                  // false (beyond range)
//	is.Latitude("invalid")             // false
func Latitude[T any](t T) bool {
	ref := reflect.ValueOf(t)
	var v string
	switch ref.Kind() {
	case reflect.String:
		v = ref.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(ref.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(ref.Uint(), 10)
	case reflect.Float32:
		v = strconv.FormatFloat(ref.Float(), 'f', -1, 32)
	case reflect.Float64:
		v = strconv.FormatFloat(ref.Float(), 'f', -1, 64)
	default:
		//fmt.Errorf("bad ref type %T", ref.Interface())
		return false
	}
	return latitudeRegex.MatchString(v)
}

// Longitude is the validation function for validating if the field's value is a valid longitude coordinate.
func Longitude[T any](t T) bool {
	ref := reflect.ValueOf(t)
	var v string
	switch ref.Kind() {
	case reflect.String:
		v = ref.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(ref.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(ref.Uint(), 10)
	case reflect.Float32:
		v = strconv.FormatFloat(ref.Float(), 'f', -1, 32)
	case reflect.Float64:
		v = strconv.FormatFloat(ref.Float(), 'f', -1, 64)
	default:
		//fmt.Errorf("bad field type %T", ref.Interface())
		return false
	}
	return longitudeRegex.MatchString(v)
}

// JSON validates whether the given value is valid JSON data.
// It accepts both string and byte slice inputs, making it flexible for different use cases.
//
// Examples:
//
//	is.JSON(`{"name": "John", "age": 30}`) // true
//	is.JSON([]byte(`[1, 2, 3]`))        // true
//	is.JSON("invalid json")              // false
//	is.JSON(123)                        // false
func JSON[T any](t T) bool {
	rv := reflect.ValueOf(t)
	if rv.Type() == nilType {
		return json.Valid(rv.Bytes())
	}
	if rv.Kind() == reflect.String {
		return json.Valid([]byte(rv.String()))
	}
	return false
}

// Datetime validates whether the given string matches the specified time layout format.
// It uses Go's standard time parsing, making it compatible with all predefined layouts.
//
// Examples:
//
//	is.Datetime("2023-12-25", "2006-01-02")        // true
//	is.Datetime("10:30:45", "15:04:05")           // true
//	is.Datetime("2023-12-25 10:30:45", time.RFC3339) // true
//	is.Datetime("invalid", "2006-01-02")            // false
func Datetime(str, layout string) bool {
	_, err := time.Parse(layout, str)
	return err == nil
}

// Timezone validates whether the given string is a valid time zone identifier.
// It validates against the IANA Time Zone Database, which is the standard for time zones.
//
// Examples:
//
//	is.Timezone("America/New_York")  // true
//	is.Timezone("UTC")               // true
//	is.Timezone("Europe/London")     // true
//	is.Timezone("")                  // false (empty)
//	is.Timezone("Local")             // false (reserved keyword)
//	is.Timezone("Invalid/Zone")     // false
func Timezone(str string) bool {
	// empty value is converted to UTC by time.LoadLocation but disallow it as it is not a valid time zone name
	if str == "" {
		return false
	}

	// Local value is converted to the current system time zone by time.LoadLocation but disallow it as it is not a valid time zone name
	if strings.ToLower(str) == "local" {
		return false
	}

	_, err := time.LoadLocation(str)
	return err == nil
}

// IPv4 validates whether the given string is a valid IPv4 address.
// IPv4 addresses consist of four decimal numbers (0-255) separated by dots.
// Leading zeros are allowed but each octet must be in the valid range.
//
// Examples:
//
//	is.IPv4("192.168.1.1")          // true (private network)
//	is.IPv4("8.8.8.8")              // true (public DNS)
//	is.IPv4("127.0.0.1")            // true (localhost)
//	is.IPv4("255.255.255.255")      // true (broadcast)
//	is.IPv4("0.0.0.0")              // true (unspecified)
//	is.IPv4("192.168.1")            // false (missing octet)
//	is.IPv4("192.168.1.256")        // false (octet > 255)
//	is.IPv4("2001:db8::1")          // false (IPv6 format)
func IPv4(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() != nil
}

// IPv6 validates whether the given string is a valid IPv6 address.
// IPv6 addresses consist of eight groups of hexadecimal digits separated by colons.
// Supports compressed notation using double colons (::) and embedded IPv4 addresses.
//
// Examples:
//
//	is.IPv6("2001:db8::1")           // true (documentation prefix)
//	is.IPv6("::1")                   // true (localhost)
//	is.IPv6("fe80::1")               // true (link-local)
//	is.IPv6("2001:0db8:85a3:0000:0000:8a2e:0370:7334") // true (fully expanded)
//	is.IPv6("2001:db8::ffff:192.0.2.1") // true (embedded IPv4)
//	is.IPv6("192.168.1.1")           // false (IPv4 format)
//	is.IPv6("2001:db8:::1")          // false (invalid triple colon)
//	is.IPv6("2001:db8:12345::1")     // false (invalid hex group)
func IPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() == nil
}

// IP validates whether the given string is a valid IP address (IPv4 or IPv6).
// This function accepts both IPv4 and IPv6 address formats and validates
// them according to their respective standards using Go's net.ParseIP.
//
// Examples:
//
//	is.IP("192.168.1.1")            // true (IPv4)
//	is.IP("2001:db8::1")            // true (IPv6)
//	is.IP("::1")                    // true (IPv6 localhost)
//	is.IP("127.0.0.1")              // true (IPv4 localhost)
//	is.IP("invalid")                // false (invalid format)
//	is.IP("256.256.256.256")        // false (invalid IPv4)
//	is.IP("2001:db8:::1")           // false (invalid IPv6)
func IP(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil
}

// MAC validates whether the given string is a valid MAC address (Media Access Control address).
// MAC addresses are unique identifiers assigned to network interfaces.
// Supports both colon-separated (00:00:5e:00:53:01) and hyphen-separated (00-00-5e-00-53-01) formats.
//
// Examples:
//
//	is.MAC("00:00:5e:00:53:01")     // true (colon format)
//	is.MAC("00-00-5e-00-53-01")     // true (hyphen format)
//	is.MAC("aa:bb:cc:dd:ee:ff")     // true (lowercase hex)
//	is.MAC("AA:BB:CC:DD:EE:FF")     // true (uppercase hex)
//	is.MAC("00:00:5e:00:53")        // false (too few octets)
//	is.MAC("00:00:5e:00:53:01:02")  // false (too many octets)
//	is.MAC("gg:hh:ii:jj:kk:ll")     // false (invalid hex)
func MAC(str string) bool {
	_, err := net.ParseMAC(str)
	return err == nil
}

// Lowercase validates whether the given string contains only lowercase letters.
// The function returns true only if all alphabetic characters are lowercase.
// Non-alphabetic characters (numbers, symbols, spaces) are allowed but don't affect the result.
// Empty strings return false as they contain no characters to validate.
//
// Examples:
//
//	is.Lowercase("hello world")       // true
//	is.Lowercase("test123")           // true (numbers are ignored)
//	is.Lowercase("你好")              // true (Chinese characters are ignored)
//	is.Lowercase("Hello World")       // false (contains uppercase H and W)
//	is.Lowercase("HELLO")             // false (all uppercase)
//	is.Lowercase("")                  // false (empty string)
func Lowercase(str string) bool {
	if str == "" {
		return false
	}
	return str == strings.ToLower(str)
}

// Uppercase validates whether the given string contains only uppercase letters.
// The function returns true only if all alphabetic characters are uppercase.
// Non-alphabetic characters (numbers, symbols, spaces) are allowed but don't affect the result.
// Empty strings return false as they contain no characters to validate.
//
// Examples:
//
//	is.Uppercase("HELLO WORLD")       // true
//	is.Uppercase("TEST123")           // true (numbers are ignored)
//	is.Uppercase("你好")              // true (Chinese characters are ignored)
//	is.Uppercase("Hello World")       // false (contains lowercase h and w)
//	is.Uppercase("hello")             // false (all lowercase)
//	is.Uppercase("")                  // false (empty string)
func Uppercase(str string) bool {
	if str == "" {
		return false
	}
	return str == strings.ToUpper(str)
}

// Empty checks if a value is empty or not using a comprehensive set of rules.
// A value is considered empty if it meets any of the following conditions:
//   - String, Array, Map, Slice: length == 0
//   - Integer, Float: value == 0
//   - Boolean: value == false
//   - Interface, Pointer, Channel, Func: nil
//   - Time: zero time value
//   - Invalid: invalid reflect value
//
// Examples:
//
//	is.Empty("")                    // true
//	is.Empty(0)                     // true
//	is.Empty(false)                 // true
//	is.Empty([]string{})            // true
//	is.Empty(nil)                   // true
//	is.Empty(time.Time{})           // true
//	is.Empty("hello")               // false
//	is.Empty(42)                    // false
func Empty[T any](t T) bool {
	rv := reflect.ValueOf(t)
	switch rv.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Invalid:
		return true
	case reflect.Interface, reflect.Ptr:
		if rv.IsNil() {
			return true
		}
		return Empty(rv.Elem().Interface())
	case reflect.Struct:
		v, ok := rv.Interface().(time.Time)
		if ok && v.IsZero() {
			return true
		}
	}
	return false
}

func NotEmpty[T any](t T) bool {
	return !Empty(t)
}

// URLEncoded validates whether the given string is properly URL-encoded.
// URL encoding (also known as percent encoding) converts characters to a format
// safe for transmission in URLs using % followed by two hexadecimal digits.
//
// Examples:
//
//	is.URLEncoded("Hello%20World")           // true (encoded space)
//	is.URLEncoded("name%3Dvalue")            // true (encoded equals)
//	is.URLEncoded("Hello%2GWorld")           // false (invalid hex)
//	is.URLEncoded("Hello World")             // false (not encoded)
//	is.URLEncoded("%")                       // false (incomplete encoding)
func URLEncoded(str string) bool {
	return uRLEncodedRegex.MatchString(str)
}

// HTMLEncoded validates whether the given string contains HTML entities.
// HTML entities are used to display reserved characters in HTML, starting with & and ending with ;.
// Common examples include &amp; (&), &lt; (<), &gt; (>), &quot; ("), and &#39; (').
//
// Examples:
//
//	is.HTMLEncoded("Hello&amp;World")          // true (encoded &)
//	is.HTMLEncoded("&lt;div&gt;content&lt;/div&gt;") // true (encoded tags)
//	is.HTMLEncoded("Hello&quot;World")          // true (encoded quotes)
//	is.HTMLEncoded("Hello & World")            // false (not properly encoded)
//	is.HTMLEncoded("&")                        // false (incomplete entity)
func HTMLEncoded(str string) bool {
	return hTMLEncodedRegex.MatchString(str)
}

// HTML validates whether the given string contains HTML tags.
// HTML tags are markup elements enclosed in angle brackets (< >).
// This function checks for common HTML tags including opening, closing, and self-closing tags.
//
// Examples:
//
//	is.HTML("<div>content</div>")              // true (basic div)
//	is.HTML("<p>Hello World</p>")               // true (paragraph tag)
//	is.HTML("<img src='image.jpg' />")          // true (self-closing tag)
//	is.HTML("<a href='link.html'>Click</a>")    // true (anchor tag)
//	is.HTML("Hello World")                      // false (no tags)
//	is.HTML("<div")                            // false (incomplete tag)
//	is.HTML(">content<")                       // false (malformed tags)
func HTML(str string) bool {
	return hTMLRegex.MatchString(str)
}

// File validates whether the given path points to an existing file (not a directory).
// This function checks the filesystem to determine if the path exists and is a regular file.
// Accepts string paths and returns false for directories or non-existent paths.
//
// Examples:
//
//	is.File("/etc/hosts")                    // true (if file exists)
//	is.File("./config.json")                 // true (if file exists in current directory)
//	is.File("/tmp")                          // false (directory)
//	is.File("/nonexistent/file")             // false (doesn't exist)
//	is.File(123)                             // false (not a string)
func File(val any) bool {
	field := reflect.ValueOf(val)
	switch field.Kind() {
	case reflect.String:
		fileInfo, err := os.Stat(field.String())
		if err != nil {
			return false
		}
		return !fileInfo.IsDir()
	}
	//fmt.Errorf("bad value type %T", field.Interface())
	return false
}

// Dir validates whether the given path points to an existing directory.
// This function checks the filesystem to determine if the path exists and is a directory.
// Accepts string paths and returns false for regular files or non-existent paths.
//
// Examples:
//
//	is.Dir("/tmp")                          // true (if directory exists)
//	is.Dir("./")                            // true (current directory)
//	is.Dir("/etc/hosts")                    // false (regular file)
//	is.Dir("/nonexistent/directory")        // false (doesn't exist)
//	is.Dir(123)                             // false (not a string)
func Dir(val any) bool {
	field := reflect.ValueOf(val)
	if field.Kind() == reflect.String {
		fileInfo, err := os.Stat(field.String())
		if err != nil {
			return false
		}
		return fileInfo.IsDir()
	}
	//fmt.Errorf("bad field type %T", field.Interface())
	return false
}

// OneOf validates whether the given value is one of the allowed values.
// This function checks if the provided value matches any value in the allowed slice.
// Useful for validating input against a predefined set of options.
//
// Examples:
//
//	is.OneOf("apple", []any{"apple", "banana", "orange"})  // true
//	is.OneOf(5, []any{1, 2, 3, 4, 5})                     // true
//	is.OneOf("red", []any{"green", "blue"})               // false
//	is.OneOf("test", []any{})                             // false (empty allowed values)
//	is.OneOf("hello", []any{"Hello", "WORLD"})            // false (case-sensitive)
func OneOf(val any, vals []any) bool {
	if len(vals) == 0 {
		return false
	}
	for _, a := range vals {
		if Equal(val, a) {
			return true
		}
	}
	return false
}

// Length validates whether the given value's length matches the specified criteria.
// Supports various value types: strings, arrays, slices, maps, and channels.
// Uses comparison operators: = (equals), != (not equals), > (greater than), < (less than), >= (greater or equal), <= (less or equal).
//
// Examples:
//
//	is.Length("hello", 5, "=")                 // true (exact length match)
//	is.Length([]int{1, 2, 3}, 3, "=")          // true (array length)
//	is.Length("world", 5, ">=")                // true (at least 5 characters)
//	is.Length("test", 5, "<")                  // true (less than 5 characters)
//	is.Length("hello", 10, "<=")               // true (10 or fewer characters)
//	is.Length(123, 3, "=")                     // false (numbers not supported)
func Length(val any, length int, op string) bool {
	if n := calcLength(val); n == -1 {
		return false
	} else {
		return Compare(n, length, op)
	}
}

// LengthBetween validates whether the given value's length falls within the specified range.
// Supports various value types: strings, arrays, slices, maps, and channels.
// The range is inclusive, checking that length is >= min and <= max.
//
// Examples:
//
//	is.LengthBetween("hello", 3, 8)            // true (length 5 is within range)
//	is.LengthBetween([]int{1, 2, 3}, 2, 5)     // true (length 3 is within range)
//	is.LengthBetween("world", 1, 4)           // false (length 5 exceeds max)
//	is.LengthBetween("test", 6, 10)           // false (length 4 below min)
//	is.LengthBetween("hello", 5, 5)           // true (exact match within range)
//	is.LengthBetween(12345, 3, 8)             // false (numbers not supported)
func LengthBetween(val any, min, max int) bool {
	if !Compare(min, max, "<") {
		panic(ErrBadType)
	} else if n, err := getLength(val, false); err != nil {
		return false
	} else {
		return Compare(n, min, ">=") && Compare(n, max, "<=")
	}
}

// Compare performs a comparison operation between two values using the specified operator.
// It supports various data types including numbers, strings, booleans, and time.Time values.
// The function returns the result of `srcVal op dstVal` where op can be "=", "!=", "<", "<=", ">", or ">=".
//
// Supported operators:
//   - "="   : equal to
//   - "!="  : not equal to
//   - "<"   : less than
//   - "<="  : less than or equal to
//   - ">"   : greater than
//   - ">="  : greater than or equal to
//
// Examples:
//
//	is.Compare(2, 3, ">")           // false
//	is.Compare(2, 1.3, ">")         // true
//	is.Compare(2.2, 1.3, ">")       // true
//	is.Compare("apple", "banana", "<") // true
//	is.Compare(true, false, ">")   // true
//	is.Compare(time.Now(), time.Now().Add(-time.Hour), ">") // true
func Compare(srcVal, dstVal any, op string) bool {
	srv := reflect.ValueOf(srcVal)

	switch srv.Kind() {
	case reflect.Struct:
		if srv.Type().ConvertibleTo(timeType) {
			drv := reflect.ValueOf(dstVal)
			if drv.Type().ConvertibleTo(timeType) {
				at := srv.Convert(timeType).Interface().(time.Time)
				bt := drv.Convert(timeType).Interface().(time.Time)
				return compTime(at, bt, op)
			}
		}
	case reflect.Bool:
		drv := reflect.ValueOf(dstVal)
		switch drv.Kind() {
		case reflect.Bool:
			return compBool(srv.Bool(), drv.Bool(), op)
		case reflect.String:
			if bl, err := strconv.ParseBool(drv.String()); err == nil {
				return compBool(srv.Bool(), bl, op)
			}
		}
	default:
		if srcStr, ok := srcVal.(string); ok {
			if dstStr, ok2 := dstVal.(string); ok2 {
				return compString(srcStr, dstStr, op)
			}
			break
		}
		// float
		if srcFlt, ok := srcVal.(float64); ok {
			if dstFlt, err := toFloat(dstVal); err == nil {
				return compNum(srcFlt, dstFlt, op)
			}
			break
		}
		if srcFlt, ok := srcVal.(float32); ok {
			if dstFlt, err := toFloat(dstVal); err == nil {
				return compNum(float64(srcFlt), dstFlt, op)
			}
			break
		}
		// as int64
		if srcInt, err := toInt64(srcVal); err != nil {
			break
		} else if dstInt, ex := toInt64(dstVal); ex != nil {
			break
		} else {
			return compNum(srcInt, dstInt, op)
		}
	}

	switch op {
	case "=":
		return srcVal == dstVal
	case "!=":
		return srcVal != dstVal
	default:
		//ErrBadType
		return false
	}
}

// GreaterThan validates whether value 'a' is greater than value 'b'.
// This function supports various data types including numbers, strings, booleans, and time.Time.
// Uses the Compare function internally with ">" operator.
//
// Examples:
//
//	is.GreaterThan(10, 5)             // true
//	is.GreaterThan("z", "a")           // true (lexicographic comparison)
//	is.GreaterThan(3.14, 2.71)         // true
//	is.GreaterThan(time.Now(), time.Now().Add(-time.Hour)) // true
func GreaterThan(a, b any) bool {
	return Compare(a, b, ">")
}

// GreaterEqualThan validates whether value 'a' is greater than or equal to value 'b'.
// This function supports various data types including numbers, strings, booleans, and time.Time.
// Uses the Compare function internally with ">=" operator.
//
// Examples:
//
//	is.GreaterEqualThan(10, 10)        // true (equal)
//	is.GreaterEqualThan(15, 5)         // true (greater)
//	is.GreaterEqualThan("hello", "hello") // true (equal strings)
//	is.GreaterEqualThan(3.14, 2.71)    // true
func GreaterEqualThan(a, b any) bool {
	return Compare(a, b, ">=")
}

// LessThan validates whether value 'a' is less than value 'b'.
// This function supports various data types including numbers, strings, booleans, and time.Time.
// Uses the Compare function internally with "<" operator.
//
// Examples:
//
//	is.LessThan(5, 10)                 // true
//	is.LessThan("a", "z")              // true (lexicographic comparison)
//	is.LessThan(2.71, 3.14)            // true
//	is.LessThan(time.Now().Add(-time.Hour), time.Now()) // true
func LessThan(a, b any) bool {
	return Compare(a, b, "<")
}

// LessEqualThan validates whether value 'a' is less than or equal to value 'b'.
// This function supports various data types including numbers, strings, booleans, and time.Time.
// Uses the Compare function internally with "<=" operator.
//
// Examples:
//
//	is.LessEqualThan(10, 10)           // true (equal)
//	is.LessEqualThan(5, 15)            // true (less)
//	is.LessEqualThan("hello", "hello") // true (equal strings)
//	is.LessEqualThan(2.71, 3.14)       // true
func LessEqualThan(a, b any) bool {
	return Compare(a, b, "<=")
}

// Equal validates whether value 'a' is equal to value 'b'.
// This function supports various data types including numbers, strings, booleans, and time.Time.
// Uses the Compare function internally with "=" operator for strict equality comparison.
//
// Examples:
//
//	is.Equal(10, 10)                   // true
//	is.Equal("hello", "hello")         // true
//	is.Equal(true, true)               // true
//	is.Equal(3.14, 3.14)               // true
//	is.Equal("hello", "Hello")         // false (case-sensitive)
//	is.Equal(10, "10")                 // false (different types)
func Equal(a, b any) bool {
	return Compare(a, b, "=")
}

// NotEqual validates whether value 'a' is not equal to value 'b'.
// This function supports various data types including numbers, strings, booleans, and time.Time.
// Uses the Compare function internally with "!=" operator for inequality comparison.
//
// Examples:
//
//	is.NotEqual(10, 5)                 // true
//	is.NotEqual("hello", "world")      // true
//	is.NotEqual(true, false)           // true
//	is.NotEqual(3.14, 2.71)            // true
//	is.NotEqual("hello", "hello")      // false (equal values)
//	is.NotEqual(10, "10")              // true (different types)
func NotEqual(a, b any) bool {
	return Compare(a, b, "!=")
}

// Between validates whether value 'val' is within the inclusive range [min, max].
// The function checks that min <= val <= max, requiring min to be less than max.
// If min >= max, the function panics with ErrBadRange.
// Supports various data types including numbers, strings, and time.Time.
//
// Examples:
//
//	is.Between(5, 1, 10)               // true (within range)
//	is.Between(10, 1, 10)              // true (equal to max)
//	is.Between(1, 1, 10)               // true (equal to min)
//	is.Between(0, 1, 10)               // false (below min)
//	is.Between(11, 1, 10)              // false (above max)
//	is.Between("m", "a", "z")          // true (lexicographic)
func Between(val, min, max any) bool {
	if !Compare(min, max, ">") {
		panic(ErrBadRange)
	}

	return Compare(val, min, ">=") && Compare(val, max, "<=")
}

// NotBetween validates whether value 'val' is outside the inclusive range [min, max].
// The function checks that val < min OR val > max, requiring min to be less than max.
// If min >= max, the function panics with ErrBadRange.
// This is the opposite of the Between function.
// Supports various data types including numbers, strings, and time.Time.
//
// Examples:
//
//	is.NotBetween(5, 1, 10)            // false (within range)
//	is.NotBetween(0, 1, 10)            // true (below min)
//	is.NotBetween(11, 1, 10)           // true (above max)
//	is.NotBetween(1, 1, 10)            // false (equal to min)
//	is.NotBetween(10, 1, 10)           // false (equal to max)
//	is.NotBetween("a", "m", "z")       // true (before range)
func NotBetween(val, min, max any) bool {
	if !Compare(min, max, ">") {
		panic(ErrBadRange)
	}

	return Compare(val, min, "<") || Compare(val, max, ">")
}
