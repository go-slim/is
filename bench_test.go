package is

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkURL(b *testing.B) {
	s := "https://user:pass@example.com:8443/a/b/c?x=1#frag"
	for b.Loop() {
		_ = URL(s)
	}
}

func BenchmarkEmail(b *testing.B) {
	s := "first.last+tag@example.co.uk"
	for b.Loop() {
		_ = Email(s)
	}
}

func BenchmarkJSON_String(b *testing.B) {
	s := "{\"a\":1,\"b\":[1,2,3],\"c\":{\"d\":true}}"
	for b.Loop() {
		_ = JSON(s)
	}
}

func BenchmarkJSON_Bytes(b *testing.B) {
	s := []byte("{\"a\":1,\"b\":[1,2,3],\"c\":{\"d\":true}}")
	for b.Loop() {
		_ = JSON(s)
	}
}

func BenchmarkCompare_Number(b *testing.B) {
	for b.Loop() {
		_ = Compare(123, 45, ">")
	}
}

func BenchmarkIPv4(b *testing.B) {
	s := "192.168.1.100"
	for b.Loop() {
		_ = IPv4(s)
	}
}

func BenchmarkIPv6(b *testing.B) {
	s := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	for b.Loop() {
		_ = IPv6(s)
	}
}

func BenchmarkColor_HEX(b *testing.B) {
	s := "#1a2b3c"
	for b.Loop() {
		_ = HEXColor(s)
	}
}

func BenchmarkColor_RGB(b *testing.B) {
	s := "rgb(12,34,56)"
	for b.Loop() {
		_ = RGB(s)
	}
}

func BenchmarkLength_String(b *testing.B) {
	s := "abcdefghijklmnopqrstuvwxyz"
	for b.Loop() {
		_ = Length(s, 26, "=")
	}
}

func BenchmarkEmpty_Mixed(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Empty(0)
		_ = Empty("")
		_ = Empty(false)
		_ = Empty([]int{})
		_ = Empty(time.Time{})
	}
}

func BenchmarkNumeric_String(b *testing.B) {
	for i := 0; b.Loop(); i++ {
		_ = Numeric(strconv.Itoa(i))
	}
}
