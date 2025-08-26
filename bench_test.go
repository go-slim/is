package is

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkURL(b *testing.B) {
	s := "https://user:pass@example.com:8443/a/b/c?x=1#frag"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = URL(s)
	}
}

func BenchmarkEmail(b *testing.B) {
	s := "first.last+tag@example.co.uk"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Email(s)
	}
}

func BenchmarkJSON_String(b *testing.B) {
	s := "{\"a\":1,\"b\":[1,2,3],\"c\":{\"d\":true}}"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = JSON(s)
	}
}

func BenchmarkJSON_Bytes(b *testing.B) {
	s := []byte("{\"a\":1,\"b\":[1,2,3],\"c\":{\"d\":true}}")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = JSON(s)
	}
}

func BenchmarkCompare_Number(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Compare(123, 45, ">")
	}
}

func BenchmarkIPv4(b *testing.B) {
	s := "192.168.1.100"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IPv4(s)
	}
}

func BenchmarkIPv6(b *testing.B) {
	s := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IPv6(s)
	}
}

func BenchmarkColor_HEX(b *testing.B) {
	s := "#1a2b3c"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HEXColor(s)
	}
}

func BenchmarkColor_RGB(b *testing.B) {
	s := "rgb(12,34,56)"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RGB(s)
	}
}

func BenchmarkLength_String(b *testing.B) {
	s := "abcdefghijklmnopqrstuvwxyz"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Length(s, 26, "=")
	}
}

func BenchmarkEmpty_Mixed(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Empty(0)
		_ = Empty("")
		_ = Empty(false)
		_ = Empty([]int{})
		_ = Empty(time.Time{})
	}
}

func BenchmarkNumeric_String(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Numeric(strconv.Itoa(i))
	}
}
