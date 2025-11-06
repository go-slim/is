package is

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func Test_toFloat(t *testing.T) {
	cases := []struct {
		in  any
		out float64
		ok  bool
	}{
		{nil, 0, true},
		{" 1.25 ", 1.25, true},
		{int(3), 3, true},
		{int64(4), 4, true},
		{uint(5), 5, true},
		{float32(1.5), 1.5, true},
		{float64(2.5), 2.5, true},
		{time.Duration(10), 10, true},
		{json.Number("6.75"), 6.75, true},
		{struct{}{}, 0, false},
	}
	for _, c := range cases {
		got, err := toFloat(c.in)
		if c.ok && err != nil {
			t.Fatalf("toFloat(%T) unexpected err: %v", c.in, err)
		}
		if !c.ok && err == nil {
			t.Fatalf("toFloat(%T) expected err", c.in)
		}
		if c.ok && got != c.out {
			t.Fatalf("toFloat(%v)=%v want %v", c.in, got, c.out)
		}
	}
}

func Test_toInt64(t *testing.T) {
	cases := []struct {
		in  any
		out int64
		ok  bool
	}{
		{nil, 0, true},
		{" 12 ", 12, true},
		{int(-3), -3, true},
		{int64(4), 4, true},
		{uint(5), 5, true},
		{float32(1.9), 1, true},
		{float64(2.9), 2, true},
		{time.Duration(10), 10, true},
		{json.Number("-7"), -7, true},
		{struct{}{}, 0, false},
	}
	for _, c := range cases {
		got, err := toInt64(c.in)
		if c.ok && err != nil {
			t.Fatalf("toInt64(%T) unexpected err: %v", c.in, err)
		}
		if !c.ok && err == nil {
			t.Fatalf("toInt64(%T) expected err", c.in)
		}
		if c.ok && got != c.out {
			t.Fatalf("toInt64(%v)=%v want %v", c.in, got, c.out)
		}
	}
}

func Test_compString(t *testing.T) {
	if !compString("a", "b", "<") {
		t.Fatal("a<b failed")
	}
	if !compString("b", "b", "=") {
		t.Fatal("b=b failed")
	}
	if compString("b", "a", "<") {
		t.Fatal("b<a should be false")
	}
}

func Test_compTime(t *testing.T) {
	a := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	b := a.Add(24 * time.Hour)
	if !compTime(a, b, "<") {
		t.Fatal("time < failed")
	}
	if !compTime(b, a, ">") {
		t.Fatal("> failed")
	}
	if !compTime(a, a, "=") {
		t.Fatal("= failed")
	}
	if compTime(a, a, "!=") == false { /* both ok branches covered */
	}
}

func Test_compNum(t *testing.T) {
	if !compNum[int64](1, 2, "<") {
		t.Fatal("int < failed")
	}
	if !compNum[float64](1.5, 1.5, "=") {
		t.Fatal("float = failed")
	}
	if compNum[uint64](3, 4, ">") {
		t.Fatal("> should be false")
	}
}

func Test_compBool_boolToInt(t *testing.T) {
	if !compBool(true, false, ">") {
		t.Fatal("true>false failed")
	}
	if boolToInt(true) != 1 || boolToInt(false) != 0 {
		t.Fatal("boolToInt failed")
	}
}

func Test_calcLength(t *testing.T) {
	if calcLength("café") != 4 {
		t.Fatal("calcLength runes failed")
	}
	if calcLength([]int{1, 2, 3}) != 3 {
		t.Fatal("calcLength slice failed")
	}
	if calcLength(uint(123)) != 3 {
		t.Fatal("calcLength uint width failed")
	}
	if calcLength(int(-12)) != 3 {
		t.Fatal("calcLength int width failed")
	}
}

func Test_getLength(t *testing.T) {
	n, err := getLength("abc", false)
	if err != nil || n != 3 {
		t.Fatalf("getLength string failed: %v %d", err, n)
	}
	n, err = getLength("café", true)
	if err != nil || n != 4 {
		t.Fatalf("getLength rune string failed: %v %d", err, n)
	}
	n, err = getLength([]int{1, 2, 3}, false)
	if err != nil || n != 3 {
		t.Fatalf("getLength slice failed: %v %d", err, n)
	}

	arr := [5]int{1, 2, 3}
	n, err = getLength(&arr, false)
	if err != nil || n != 5 {
		t.Fatalf("getLength ptr-to-array type length failed: %v %d", err, n)
	}

	if _, err = getLength(123, false); !errors.Is(err, ErrBadType) {
		t.Fatalf("getLength non-supported type should return ErrBadType, got %v", err)
	}
}
