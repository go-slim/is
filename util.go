package is

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Error definitions for the is package
var (
	// ErrBadType is returned when a value has an unsupported or unexpected type
	ErrBadType = errors.New("bad value type")

	// ErrBadRange is returned when a range validation has invalid parameters (min > max)
	ErrBadRange = errors.New("bad value range")
)

// Pre-allocated reflect types for type checking optimization
var (
	// timeType holds the reflect.Type for time.Time for faster type comparisons
	timeType = reflect.TypeOf(time.Time{})

	// nilType holds the reflect.Type for nil byte slice for JSON validation
	nilType = reflect.TypeOf([]byte(nil))
)

// toFloat converts any value to float64.
// It supports numeric types, time.Duration, and string representations.
// Returns ErrBadType if the conversion is not possible.
//
// Supported input types:
//   - Numeric types (int, uint, float variants)
//   - time.Duration
//   - json.Number
//   - string (parsable as float)
//   - nil (returns 0)
func toFloat(in any) (f64 float64, err error) {
	switch tVal := in.(type) {
	case nil:
		f64 = 0
	case string:
		f64, err = strconv.ParseFloat(strings.TrimSpace(tVal), 64)
	case int:
		f64 = float64(tVal)
	case int8:
		f64 = float64(tVal)
	case int16:
		f64 = float64(tVal)
	case int32:
		f64 = float64(tVal)
	case int64:
		f64 = float64(tVal)
	case uint:
		f64 = float64(tVal)
	case uint8:
		f64 = float64(tVal)
	case uint16:
		f64 = float64(tVal)
	case uint32:
		f64 = float64(tVal)
	case uint64:
		f64 = float64(tVal)
	case float32:
		f64 = float64(tVal)
	case float64:
		f64 = tVal
	case time.Duration:
		f64 = float64(tVal)
	case json.Number:
		f64, err = tVal.Float64()
	default:
		err = ErrBadType
	}
	return
}

// toInt64 converts any value to int64.
// It supports numeric types, time.Duration, and string representations.
// Returns ErrBadType if the conversion is not possible.
//
// Supported input types:
//   - Numeric types (int, uint, float variants)
//   - time.Duration
//   - json.Number
//   - string (parsable as integer)
//   - nil (returns 0)
func toInt64(in any) (i64 int64, err error) {
	switch tVal := in.(type) {
	case nil:
		i64 = 0
	case string:
		i64, err = strconv.ParseInt(strings.TrimSpace(tVal), 10, 0)
	case int:
		i64 = int64(tVal)
	case int8:
		i64 = int64(tVal)
	case int16:
		i64 = int64(tVal)
	case int32:
		i64 = int64(tVal)
	case int64:
		i64 = tVal
	case uint:
		i64 = int64(tVal)
	case uint8:
		i64 = int64(tVal)
	case uint16:
		i64 = int64(tVal)
	case uint32:
		i64 = int64(tVal)
	case uint64:
		i64 = int64(tVal)
	case float32:
		i64 = int64(tVal)
	case float64:
		i64 = int64(tVal)
	case time.Duration:
		i64 = int64(tVal)
	case json.Number:
		i64, err = tVal.Int64()
	default:
		err = ErrBadType
	}
	return
}

// compString compares two strings using the specified operator.
// It uses Go's strings.Compare for lexicographical comparison.
// Returns the result of `first op second`.
func compString(first, second, op string) bool {
	rs := strings.Compare(first, second)
	if rs < 0 {
		return op == "<" || op == "<="
	} else if rs > 0 {
		return op == ">" || op == ">="
	} else {
		return op == ">=" || op == "<=" || op == "="
	}
}

// compTime compares two time.Time values using the specified operator.
// It supports temporal comparison operators and uses Go's built-in time comparison methods.
// Returns the result of `first op dstTime`.
func compTime(first, dstTime time.Time, op string) (ok bool) {
	switch op {
	case "<":
		return first.Before(dstTime)
	case "<=":
		return first.Before(dstTime) || first.Equal(dstTime)
	case ">":
		return first.After(dstTime)
	case ">=":
		return first.After(dstTime) || first.Equal(dstTime)
	case "=":
		return first.Equal(dstTime)
	case "!=":
		return !first.Equal(dstTime)
	}
	return
}

// compNum compares two numeric values using the specified operator.
// This is a generic function that works with int64, float64, and uint64 types.
// Returns the result of `first op second`.
func compNum[T int64 | float64 | uint64](first, second T, op string) bool {
	switch op {
	case "<":
		return first < second
	case "<=":
		return first <= second
	case ">":
		return first > second
	case ">=":
		return first >= second
	case "=":
		return first == second
	case "!=":
		return first != second
	}
	return false
}

// compBool compares two boolean values using the specified operator.
// Internally converts booleans to integers (false=0, true=1) for comparison.
// Returns the result of `first op second`.
func compBool(first, second bool, op string) bool {
	return compNum(boolToInt(first), boolToInt(second), op)
}

// boolToInt converts a boolean value to an integer representation.
// Returns 1 for true, 0 for false.
func boolToInt(a bool) int64 {
	if a {
		return 1
	} else {
		return 0
	}
}

// calcLength calculates the length of a value for validation purposes.
// Returns -1 if the value type doesn't support length calculation.
// Handles strings, arrays, maps, slices, and numeric types differently.
func calcLength(val any) int {
	v := reflect.Indirect(reflect.ValueOf(val))

	// (u)int use width.
	switch v.Kind() {
	case reflect.String:
		return len([]rune(v.String()))
	case reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return v.Len()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return len(strconv.FormatInt(int64(v.Uint()), 10))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return len(strconv.FormatInt(v.Int(), 10))
	case reflect.Float32, reflect.Float64:
		return len(fmt.Sprint(v.Interface()))
	default:
		// cannot get length
		return -1
	}
}

// getLength calculates the length of a value for validation purposes.
// Returns an error if the value type doesn't support length calculation.
// The rune parameter determines whether to count Unicode runes or bytes for strings.
// Handles strings, arrays, maps, slices, and pointer to arrays.
//
// When rune is true: counts Unicode runes (characters) for strings
// When rune is false: counts bytes for strings
func getLength(a any, rune bool) (int, error) {
	field := reflect.ValueOf(a)

	switch field.Kind() {
	case reflect.String:
		if rune {
			return utf8.RuneCountInString(field.String()), nil
		}
		return field.Len(), nil

	case reflect.Slice, reflect.Map, reflect.Array:
		return field.Len(), nil

	case reflect.Pointer:
		if field.Type().Elem().Kind() == reflect.Array {
			// Length from type declaration for pointer to arrays
			return field.Type().Elem().Len(), nil
		}
		return 0, ErrBadType

	default:
		return 0, ErrBadType
	}
}
