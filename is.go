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

// Base64URL 判断给出的字符串是否为有效且安全的 base64URL
func Base64URL(str string) bool {
	return base64URLRegex.MatchString(str)
}

// JWT is the validation function for validating if the current field's value is a valid JWT string.
func JWT(str string) bool {
	return jWTRegex.MatchString(str)
}

// UUID5 is the validation function for validating if the field's value is a valid v5 UUID.
func UUID5(str string) bool {
	return uUID5Regex.MatchString(str)
}

// UUID4 is the validation function for validating if the field's value is a valid v4 UUID.
func UUID4(str string) bool {
	return uUID4Regex.MatchString(str)
}

// UUID3 is the validation function for validating if the field's value is a valid v3 UUID.
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

// MD4 is the validation function for validating if the field's value is a valid MD4.
func MD4(str string) bool {
	return md4Regex.MatchString(str)
}

// MD5 is the validation function for validating if the field's value is a valid MD5.
func MD5(str string) bool {
	return md5Regex.MatchString(str)
}

// SHA256 is the validation function for validating if the field's value is a valid SHA256.
func SHA256(str string) bool {
	return sha256Regex.MatchString(str)
}

// SHA384 is the validation function for validating if the field's value is a valid SHA384.
func SHA384(str string) bool {
	return sha384Regex.MatchString(str)
}

// SHA512 is the validation function for validating if the field's value is a valid SHA512.
func SHA512(str string) bool {
	return sha512Regex.MatchString(str)
}

// ASCII is the validation function for validating if the field's value is a valid ASCII character.
func ASCII(str string) bool {
	return aSCIIRegex.MatchString(str)
}

// Alpha is the validation function for validating if the current field's value is a valid alpha value.
func Alpha(str string) bool {
	return alphaRegex.MatchString(str)
}

// Alphanumeric is the validation function for validating if the current field's value is a valid alphanumeric value.
func Alphanumeric(str string) bool {
	return alphaNumericRegex.MatchString(str)
}

// AlphaUnicode is the validation function for validating if the current field's value is a valid alpha unicode value.
func AlphaUnicode(str string) bool {
	return alphaUnicodeRegex.MatchString(str)
}

// AlphanumericUnicode is the validation function for validating if the current field's value is a valid alphanumeric unicode value.
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

// Hexadecimal is the validation function for validating if the current field's value is a valid hexadecimal.
func Hexadecimal(str string) bool {
	return hexadecimalRegex.MatchString(str)
}

// HEXColor is the validation function for validating if the current field's value is a valid HEX color.
func HEXColor(str string) bool {
	return hexColorRegex.MatchString(str)
}

// RGB is the validation function for validating if the current field's value is a valid RGB color.
func RGB(str string) bool {
	return rgbRegex.MatchString(str)
}

// RGBA is the validation function for validating if the current field's value is a valid RGBA color.
func RGBA(str string) bool {
	return rgbaRegex.MatchString(str)
}

// HSL is the validation function for validating if the current field's value is a valid HSL color.
func HSL(str string) bool {
	return hslRegex.MatchString(str)
}

// HSLA is the validation function for validating if the current field's value is a valid HSLA color.
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

// IPv4 is the validation function for validating if a value is a valid v4 IP address.
func IPv4(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() != nil
}

// IPv6 is the validation function for validating if the field's value is a valid v6 IP address.
func IPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() == nil
}

// IP is the validation function for validating if the field's value is a valid v4 or v6 IP address.
func IP(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil
}

// MAC is the validation function for validating if the field's value is a valid MAC address.
func MAC(str string) bool {
	_, err := net.ParseMAC(str)
	return err == nil
}

// Lowercase is the validation function for validating if the current field's value is a lowercase string.
func Lowercase(str string) bool {
	if str == "" {
		return false
	}
	return str == strings.ToLower(str)
}

// Uppercase is the validation function for validating if the current field's value is an uppercase string.
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

func URLEncoded(str string) bool {
	return uRLEncodedRegex.MatchString(str)
}

func HTMLEncoded(str string) bool {
	return hTMLEncodedRegex.MatchString(str)
}

func HTML(str string) bool {
	return hTMLRegex.MatchString(str)
}

// File is the validation function for validating if the current field's value is a valid file path.
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

// Dir is the validation function for validating if the current field's value is a valid directory.
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

func Length(val any, length int, op string) bool {
	if n := calcLength(val); n == -1 {
		return false
	} else {
		return Compare(n, length, op)
	}
}

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

// GreaterThan is the validation function for validating if the current field's value is greater than the param's value.
func GreaterThan(a, b any) bool {
	return Compare(a, b, ">")
}

// GreaterEqualThan is the validation function for validating if the current field's value is greater than or equal to the param's value.
func GreaterEqualThan(a, b any) bool {
	return Compare(a, b, ">=")
}

// LessThan is the validation function for validating if the current field's value is less than the param's value.
func LessThan(a, b any) bool {
	return Compare(a, b, "<")
}

// LessEqualThan is the validation function for validating if the current field's value is less than or equal to the param's value.
func LessEqualThan(a, b any) bool {
	return Compare(a, b, "<=")
}

// Equal is the validation function for validating if the current field's value is equal to the param's value.
func Equal(a, b any) bool {
	return Compare(a, b, "=")
}

func NotEqual(a, b any) bool {
	return Compare(a, b, "!=")
}

func Between(val, min, max any) bool {
	if !Compare(min, max, ">") {
		panic(ErrBadRange)
	}

	return Compare(val, min, ">=") && Compare(val, max, "<=")
}

func NotBetween(val, min, max any) bool {
	if !Compare(min, max, ">") {
		panic(ErrBadRange)
	}

	return Compare(val, min, "<") || Compare(val, max, ">")
}
