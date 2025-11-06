package is

import (
	"fmt"
	"testing"
)

func TestPhoneNumber(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		// Valid cases - domestic format
		{"valid domestic 13x", "13800138000", true},
		{"valid domestic 15x", "15912345678", true},
		{"valid domestic 18x", "18611112222", true},
		{"valid domestic 19x", "19800001111", true},
		{"valid domestic 14x", "14712345678", true},
		{"valid domestic 16x", "16512345678", true},
		{"valid domestic 17x", "17312345678", true},

		// Valid cases - international format
		{"valid international +86", "+8613800138000", true},
		{"valid international 86", "8613800138000", true},

		// Invalid cases - wrong length
		{"too short", "138001380", false},
		{"too long", "138001380000", false},

		// Invalid cases - wrong prefix
		{"starts with 0", "03800138000", false},
		{"starts with 2", "23800138000", false},
		{"starts with 9", "93800138000", false},

		// Invalid cases - non-numeric
		{"contains letters", "1380013800a", false},
		{"contains spaces", "138 0013 8000", false},
		{"contains dash", "138-0013-8000", false},

		// Edge cases
		{"empty string", "", false},
		{"only +86", "+86", false},
		{"only 86", "86", false},
		{"10 digits", "1380013800", false},
		{"12 digits", "138001380000", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PhoneNumber(tt.input); got != tt.want {
				t.Errorf("PhoneNumber(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIdcard(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		// Valid cases - 18-digit format with numeric check digit
		{"valid 18-digit numeric check", "110105194912310021", true},
		{"valid 18-digit numeric check 2", "440524188001010014", true},
		{"valid 18-digit numeric check 3", "320102198501010012", true},

		// Valid cases - 18-digit format with X check digit
		{"valid 18-digit X uppercase", "11010519491231002X", true},
		{"valid 18-digit x lowercase", "11010519491231002x", true},
		{"valid 18-digit X uppercase 2", "44052418800101001X", true},
		{"valid 18-digit x lowercase 2", "32010219850101001x", true},

		// Valid cases - 15-digit format (old format)
		{"valid 15-digit format", "110105491231002", true},
		{"valid 15-digit format 2", "440524800101001", true},
		{"valid 15-digit format 3", "320102850101001", true},

		// Invalid cases - wrong length
		{"14 digits", "11010549123100", false},
		{"16 digits", "1101054912310021", false},
		{"17 digits", "11010519491231002", false},
		{"19 digits", "1101051949123100212", false},

		// Invalid cases - invalid characters in 18-digit format
		{"18-digit with letter Y", "11010519491231002Y", false},
		{"18-digit with space", "110105194912310 21", false},
		{"18-digit with dash", "110105-1949-1231-002X", false},

		// Invalid cases - invalid characters in 15-digit format
		{"15-digit with letter", "11010549123100X", false},
		{"15-digit with space", "110105 49123100", false},

		// Invalid cases - X in wrong position (15-digit format)
		{"15-digit cannot have X", "11010549123100X", false},

		// Edge cases
		{"empty string", "", false},
		{"only X", "X", false},
		{"all zeros 15-digit", "000000000000000", true},
		{"all zeros 18-digit", "000000000000000000", true},
		{"all nines 15-digit", "999999999999999", true},
		{"all nines 18-digit", "999999999999999999", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Idcard(tt.input); got != tt.want {
				t.Errorf("Idcard(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func ExamplePhoneNumber() {
	fmt.Println(PhoneNumber("13800138000"))
	fmt.Println(PhoneNumber("+8613800138000"))
	fmt.Println(PhoneNumber("12345678901"))
	fmt.Println(PhoneNumber("2800138000"))
	// Output:
	// true
	// true
	// true
	// false
}

func ExampleIdcard() {
	fmt.Println(Idcard("11010519491231002X"))
	fmt.Println(Idcard("110105491231002"))
	fmt.Println(Idcard("110105194912310021"))
	fmt.Println(Idcard("123456789012345"))
	fmt.Println(Idcard("11010519491231002"))
	// Output:
	// true
	// true
	// true
	// true
	// false
}
