package is

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEmail(t *testing.T) {
	tests := []struct {
		in string
		ok bool
	}{
		{"user@example.com", true},
		{"user.name+tag@sub.example.co", true},
		{"bad@", false},
		{"@bad.com", false},
	}
	for _, tt := range tests {
		if got := Email(tt.in); got != tt.ok {
			t.Fatalf("Email(%q)=%v, want %v", tt.in, got, tt.ok)
		}
	}
}

func TestE164(t *testing.T) {
	tests := []struct {
		in string
		ok bool
	}{
		{"+14155552671", true},
		{"14155552671", false},
		{"+1 415 555 2671", false},
	}
	for _, tt := range tests {
		if got := E164(tt.in); got != tt.ok {
			t.Fatalf("E164(%q)=%v, want %v", tt.in, got, tt.ok)
		}
	}
}

func TestPhoneNumberCN(t *testing.T) {
	tests := []struct {
		in string
		ok bool
	}{
		{"13800138000", true},
		{"19912345678", true},
		{"123456", false},
	}
	for _, tt := range tests {
		if got := PhoneNumber(tt.in); got != tt.ok {
			t.Fatalf("PhoneNumber(%q)=%v, want %v", tt.in, got, tt.ok)
		}
	}
}

func TestSemver(t *testing.T) {
	tests := []struct {
		in string
		ok bool
	}{
		{"1.0.0", true},
		{"1.2.3-alpha+001", true},
		{"1.0", false},
	}
	for _, tt := range tests {
		if got := Semver(tt.in); got != tt.ok {
			t.Fatalf("Semver(%q)=%v, want %v", tt.in, got, tt.ok)
		}
	}
}

func TestLabel(t *testing.T) {
	tests := []struct {
		in string
		ok bool
	}{
		{"a_name1", true},
		{"_name1", false},
		{"中文", false},
		{"1abc", false},
	}
	for _, tt := range tests {
		if got := Label(tt.in); got != tt.ok {
			t.Fatalf("Label(%q)=%v, want %v", tt.in, got, tt.ok)
		}
	}
}

func TestBase64(t *testing.T) {
	if !Base64("SGVsbG8=") {
		t.Fatal("Base64(Hello) expected true")
	}
	if Base64("not-base64!!") {
		t.Fatal("Base64(bad) expected false")
	}
}

func TestURL(t *testing.T) {
	// with fragment
	if !URL("https://example.com/path#frag") {
		t.Fatal("URL expected true")
	}
	if URL("") {
		t.Fatal("URL empty expected false")
	}
	if URL("notaurl") {
		t.Fatal("URL invalid expected false")
	}
}

func TestBase64URL(t *testing.T) {
	// valid samples must be in groups of 4 or properly padded
	if !Base64URL("abcd") {
		t.Fatal("Base64URL valid expected true")
	}
	if !Base64URL("ab_-") {
		t.Fatal("Base64URL valid expected true")
	}
	if Base64URL("+/=") {
		t.Fatal("Base64URL invalid expected false")
	}
}

func TestUUIDsAndULID(t *testing.T) {
	u3 := "f47ac10b-58cc-3af2-8f5a-1234567890ab"
	u4 := "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	u5 := "f47ac10b-58cc-5ef2-8f5a-1234567890ab"
	ul := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

	if !UUID(u3) || !UUID3(u3) {
		t.Fatal("UUID3 invalid")
	}
	if !UUID(u4) || !UUID4(u4) {
		t.Fatal("UUID4 invalid")
	}
	if !UUID(u5) || !UUID5(u5) {
		t.Fatal("UUID5 invalid")
	}
	if !ULID(ul) {
		t.Fatal("ULID invalid")
	}
}

func TestHashes(t *testing.T) {
	if !MD5("d41d8cd98f00b204e9800998ecf8427e") {
		t.Fatal("MD5 invalid")
	}
	if !SHA256("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855") {
		t.Fatal("SHA256 invalid")
	}
	// SHA384 of empty string
	if !SHA384("38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da" +
		"274edebfe76f65fbd51ad2f14898b95b") {
		t.Fatal("SHA384 invalid")
	}
	if !SHA512("cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce" +
		"47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e") {
		t.Fatal("SHA512 invalid")
	}
}

func TestASCIIAlphaNumeric(t *testing.T) {
	if !ASCII("Hello123") {
		t.Fatal("ASCII failed")
	}
	if ASCII("你好") {
		t.Fatal("ASCII expected false")
	}
	if !Alpha("Hello") || Alpha("Hello1") {
		t.Fatal("Alpha checks failed")
	}
	if !Alphanumeric("abc123") || Alphanumeric("abc-123") {
		t.Fatal("Alphanumeric checks failed")
	}
	if !AlphaUnicode("你好") {
		t.Fatal("AlphaUnicode failed")
	}
	if !AlphanumericUnicode("你好123") {
		t.Fatal("AlphanumericUnicode failed")
	}
}

func TestNumericAndNumberGeneric(t *testing.T) {
	if !Numeric(123) || !Numeric(1.23) || !Numeric("123") {
		t.Fatal("Numeric true cases failed")
	}
	if Numeric("abc") {
		t.Fatal("Numeric false case failed")
	}
	if !Number(123) || !Number(1.23) {
		t.Fatal("Number true cases failed")
	}
	if Number("-12.3") {
		t.Fatal("Number string with decimal should be false (digits only)")
	}
	if Number("abc") {
		t.Fatal("Number false case failed")
	}
}

func TestBooleanGeneric(t *testing.T) {
	// valid strings
	valids := []string{"1", "yes", "YES", "Yes", "on", "ON", "On", "true", "TRUE", "True", "0", "no", "NO", "No", "", "off", "OFF", "Off", "false", "FALSE", "False"}
	for _, s := range valids {
		if !Boolean(s) {
			t.Fatalf("Boolean(%q) expected true", s)
		}
	}
	if Boolean("maybe") {
		t.Fatal("Boolean(" + "maybe" + ") expected false")
	}
	if !Boolean(1) || !Boolean(0) {
		t.Fatal("Boolean int 0/1 failed")
	}
	if !Boolean(uint(1)) || !Boolean(uint(0)) {
		t.Fatal("Boolean uint 0/1 failed")
	}
	if Boolean(2) {
		t.Fatal("Boolean int=2 should be false")
	}
}

func TestDefaultHasValue(t *testing.T) {
	if !Default(0) || HasValue(0) {
		t.Fatal("Default/HasValue int zero failed")
	}
	if Default(1) || !HasValue(1) {
		t.Fatal("Default/HasValue int non-zero failed")
	}
	var s string
	if !Default(s) || HasValue(s) {
		t.Fatal("Default/HasValue string empty failed")
	}
	s = "x"
	if Default(s) || !HasValue(s) {
		t.Fatal("Default/HasValue string non-empty failed")
	}
}

func TestColors(t *testing.T) {
	if !HEXColor("#fff") || !HEXColor("#ffffff") {
		t.Fatal("HEXColor failed")
	}
	if !RGB("rgb(0,0,0)") || !RGBA("rgba(0,0,0,0.5)") {
		t.Fatal("RGB/A failed")
	}
	if !HSL("hsl(120,100%,50%)") || !HSLA("hsla(120,100%,50%,0.3)") {
		t.Fatal("HSL/A failed")
	}
	if !Color("#fff") || !Color("rgba(1,2,3,0.2)") {
		t.Fatal("Color failed")
	}
}

func TestLatLong(t *testing.T) {
	if !Latitude("-12.34") || !Longitude("56.78") {
		t.Fatal("Lat/Long string failed")
	}
	if !Latitude(12.34) || !Longitude(-56.78) {
		t.Fatal("Lat/Long float failed")
	}
	if Latitude("999") || Longitude("999") {
		t.Fatal("Lat/Long out of range should fail (regex)")
	}
}

func TestJSONGeneric(t *testing.T) {
	if !JSON("{\"a\":1}") {
		t.Fatal("JSON string failed")
	}
	if JSON("{bad}") {
		t.Fatal("JSON bad string expected false")
	}
	b := []byte("{\"b\":2}")
	if !JSON(b) {
		t.Fatal("JSON []byte failed")
	}
}

func TestDatetimeTimezone(t *testing.T) {
	if !Datetime("2023-01-02", "2006-01-02") {
		t.Fatal("Datetime failed")
	}
	if Datetime("2023/01/02", "2006-01-02") {
		t.Fatal("Datetime expected false")
	}
	if !Timezone("UTC") {
		t.Fatal("Timezone UTC should be valid")
	}
	if Timezone("") || Timezone("local") {
		t.Fatal("Timezone empty/local should be invalid")
	}
}

func TestIPAndMAC(t *testing.T) {
	if !IPv4("192.168.0.1") || IPv6("192.168.0.1") {
		t.Fatal("IPv4/IPv6 mismatch")
	}
	if !IPv6("2001:0db8:85a3:0000:0000:8a2e:0370:7334") {
		t.Fatal("IPv6 failed")
	}
	if !IP("127.0.0.1") || !IP("::1") {
		t.Fatal("IP failed")
	}
	if !MAC("01:23:45:67:89:ab") {
		t.Fatal("MAC failed")
	}
}

func TestCaseChecks(t *testing.T) {
	if !Lowercase("abc") || Lowercase("Abc") {
		t.Fatal("Lowercase failed")
	}
	if !Uppercase("ABC") || Uppercase("Abc") {
		t.Fatal("Uppercase failed")
	}
}

func TestEmptyNotEmpty(t *testing.T) {
	if !Empty(0) || !Empty("") || !Empty(false) {
		t.Fatal("Empty basic failed")
	}
	if Empty(1) || Empty("x") || Empty(true) {
		t.Fatal("Empty basic negatives failed")
	}
	var p *int
	if !Empty(p) {
		t.Fatal("Empty nil pointer failed")
	}
	now := time.Time{}
	if !Empty(now) {
		t.Fatal("Empty zero time failed")
	}
	if !NotEmpty(1) || NotEmpty(0) {
		t.Fatal("NotEmpty failed")
	}
}

func TestEncodingDetectors(t *testing.T) {
	if !URLEncoded("a%20b") {
		t.Fatal("URLEncoded failed")
	}
	// regex allows any non-% bytes or valid %xx sequences; space is allowed
	if !URLEncoded("a b") {
		t.Fatal("URLEncoded space should be allowed by detector")
	}
	if !HTMLEncoded("&lt;div&gt;") {
		t.Fatal("HTMLEncoded failed")
	}
	if !HTML("<div>") {
		t.Fatal("HTML tag detection failed")
	}
}

func TestFileAndDir(t *testing.T) {
	dir := t.TempDir()
	// file inside dir
	fp := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(fp, []byte("hi"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	if !Dir(dir) {
		t.Fatal("Dir should be true")
	}
	if Dir(fp) {
		t.Fatal("Dir on file should be false")
	}

	if !File(fp) {
		t.Fatal("File should be true")
	}
	if File(dir) {
		t.Fatal("File on dir should be false")
	}
	if File(filepath.Join(dir, "nope")) {
		t.Fatal("File on missing path should be false")
	}
}

func TestCompareAndWrappers(t *testing.T) {
	if !GreaterThan(2, 1) || GreaterThan(1, 2) {
		t.Fatal("GreaterThan failed")
	}
	if !GreaterEqualThan(2, 2) || !LessThan(1, 2) || !LessEqualThan(2, 2) {
		t.Fatal("GE/LT/LE failed")
	}
	if !Equal("a", "a") || Equal("a", "b") {
		t.Fatal("Equal string failed")
	}
	if !NotEqual(1, 2) || NotEqual(2, 2) {
		t.Fatal("NotEqual failed")
	}

	// string compare with op in Compare
	if !Compare("a", "b", "<") || Compare("b", "b", "<") {
		t.Fatal("Compare strings failed")
	}

	// time compare
	a := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	b := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	if !Compare(a, b, "<") {
		t.Fatal("time comparisons failed: a<b")
	}
	// Between/NotBetween: current implementation guards on min>max and yields false/true respectively
	if Between(a, b, a) {
		t.Fatal("Between should be false with min>b>max ordering per implementation")
	}
	if !NotBetween(a, b, a) {
		t.Fatal("NotBetween should be true with min>b>max ordering per implementation")
	}
}

func TestLengthHelpers(t *testing.T) {
	if !Length("abc", 3, "=") {
		t.Fatal("Length string failed")
	}
	if !Length(123, 3, "=") {
		t.Fatal("Length int digits failed")
	}
	if !LengthBetween([]int{1, 2, 3}, 1, 3) {
		t.Fatal("LengthBetween slice failed")
	}
}
