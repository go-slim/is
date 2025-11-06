package is

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
