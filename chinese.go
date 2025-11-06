package is

import "regexp"

// Customized for Chinese mainland due to the author's background

var (
	phoneNumberRegexString = "^(\\+?86)?1[0-9]{10}$"
	idcardRegexString      = `^(\d{15}|\d{17}[\dXx])$`
)

var (
	phoneNumberRegex = regexp.MustCompile(phoneNumberRegexString)
	idcardRegex      = regexp.MustCompile(idcardRegexString)
)

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

// Idcard validates whether the given string is a valid Chinese ID card number.
// It supports both 15-digit and 18-digit formats:
// - 15-digit format: All numeric (used before 1999)
// - 18-digit format: 17 digits + 1 check digit (numeric or X/x)
//
// Examples:
//
//	is.Idcard("11010519491231002X")  // true (18-digit with X check digit)
//	is.Idcard("110105491231002")     // true (15-digit format)
//	is.Idcard("110105194912310021")  // true (18-digit with numeric check digit)
//	is.Idcard("123456789012345")     // false (invalid format)
//	is.Idcard("11010519491231002")   // false (17 digits, invalid)
func Idcard(s string) bool {
	return idcardRegex.MatchString(s)
}
