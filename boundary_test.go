package is

import (
	"net/url"
	"testing"
	"time"
)

func TestURL_Edges(t *testing.T) {
	cases := []struct {
		in string
		ok bool
	}{
		{"http://user:pass@example.com:8080/path?x=1#frag", true},
		{"https://[2001:db8::1]/index.html", true},
		{"ftp://example.com/resource", true},
		{"//example.com/no-scheme", false},
		{"https://example.com#", true}, // trailing '#'
		{"http://", true},
	}
	for _, c := range cases {
		if got := URL(c.in); got != c.ok {
			t.Fatalf("URL(%q)=%v, want %v", c.in, got, c.ok)
		}
	}
}

func TestEmail_Edges(t *testing.T) {
	cases := []struct {
		in string
		ok bool
	}{
		{"\"foo bar\"@example.com", true}, // quoted local-part
		{"foo..bar@example.com", false},   // double dot
		{"foo.@example.com", false},
		{"foo@sub_domain.example", false},
	}
	for _, c := range cases {
		if Email(c.in) != c.ok {
			t.Fatalf("Email(%q)=%v want %v", c.in, Email(c.in), c.ok)
		}
	}
}

func TestIPv6_Compressed(t *testing.T) {
	valids := []string{
		"::1",
		"2001:db8::1",
	}
	for _, v := range valids {
		if !IPv6(v) {
			t.Fatalf("IPv6(%q) expected true", v)
		}
	}
	// IPv4-mapped IPv6 should be false for IPv6() in current implementation
	if IPv6("::ffff:192.168.0.1") {
		t.Fatal("IPv6 should be false for IPv4-mapped address")
	}
}

func TestTimezone_Invalids(t *testing.T) {
	bads := []string{"Not/AZone", "GMT+25", "UTC+25"}
	for _, b := range bads {
		if Timezone(b) {
			t.Fatalf("Timezone(%q) expected false", b)
		}
	}
}

func TestJSON_WhitespaceAndTypes(t *testing.T) {
	if !JSON("  {\n\t\"a\":1 }  ") {
		t.Fatal("JSON with whitespace failed")
	}
	if JSON(123) {
		t.Fatal("JSON on non-string/[]byte should be false")
	}
}

func TestLengthBetween_PanicOnBadArgs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("LengthBetween should panic when min>=max")
		}
	}()
	_ = LengthBetween("abc", 3, 3)
}

func TestBetween_NotBetween(t *testing.T) {
	a := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	b := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	// Current implementation requires min>max; using (min=b, max=a) avoids panic.
	if Between(a, b, a) {
		t.Fatal("Between should be false with min>b>max ordering")
	}
	if !NotBetween(a, b, a) {
		t.Fatal("NotBetween should be true with min>b>max ordering")
	}
}

func TestURLEncoded_Complex(t *testing.T) {
	// Valid percent escapes and other chars
	s := "path/%E4%BD%A0%E5%A5%BD?q=a%2Bb&x=1 2"
	if !URLEncoded(s) {
		t.Fatalf("URLEncoded complex failed: %q", s)
	}
	// Invalid escape
	if URLEncoded("%GG") {
		t.Fatal("URLEncoded should be false on bad escape")
	}
}

func TestURL_ParseRoundTrip(t *testing.T) {
	raw := "https://example.com/a%20b?x=1#y"
	if !URL(raw) {
		t.Fatal("URL raw invalid")
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil || u == nil {
		t.Fatalf("url.ParseRequestURI failed: %v", err)
	}
}
