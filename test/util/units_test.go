package util_test

import (
	"math"
	"testing"

	"github.com/ryo-arima/goxcel/pkg/util"
)

func almostEqual(a, b, eps float64) bool { return math.Abs(a-b) <= eps }

func TestParseColWidth(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"10", 10},
		{"96px", (96.0 - 5.0) / 7.0},   // approx 13
		{"2.54cm", (96.0 - 5.0) / 7.0}, // 1 inch
		{"14pt", 1.95},                 // approx ( (14/72*96) - 5 ) / 7
	}
	for _, c := range cases {
		got, err := util.ParseColWidth(c.in)
		if err != nil {
			t.Fatalf("ParseColWidth(%q) error: %v", c.in, err)
		}
		if !almostEqual(got, c.want, 0.01) {
			t.Errorf("ParseColWidth(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseRowHeight(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"10", 10},
		{"72px", 54},   // 72 px -> 54 pt
		{"2.54cm", 72}, // 1 inch in points
	}
	for _, c := range cases {
		got, err := util.ParseRowHeight(c.in)
		if err != nil {
			t.Fatalf("ParseRowHeight(%q) error: %v", c.in, err)
		}
		// Allow small epsilon
		if !almostEqual(got, c.want, 0.5) {
			t.Errorf("ParseRowHeight(%q) = %v, want ~%v", c.in, got, c.want)
		}
	}
}
