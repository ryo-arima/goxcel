package util

import (
	"math"
	"strconv"
	"strings"
)

// ParseColWidth parses a column width string and returns Excel "character" width units (float64).
// Supported units:
// - plain number (interpreted as Excel width in characters)
// - cm, mm, in (inches), pt (points), px (pixels), ch (characters)
// Approximations used:
// - Pixels per inch = 96
// - Character width ~= (px - 5) / 7 (Excel approximation)
func ParseColWidth(s string) (float64, error) {
	val, unit := splitNumberUnit(s)
	if unit == "" || unit == "ch" {
		return val, nil
	}
	switch unit {
	case "cm":
		px := cmToPx(val)
		return pxToColWidth(px), nil
	case "mm":
		px := cmToPx(val / 10.0)
		return pxToColWidth(px), nil
	case "in":
		px := inchesToPx(val)
		return pxToColWidth(px), nil
	case "pt":
		px := pointsToPx(val)
		return pxToColWidth(px), nil
	case "px":
		return pxToColWidth(val), nil
	default:
		// Unknown unit: try numeric fallback
		return val, nil
	}
}

// ParseRowHeight parses a row height string and returns points (float64).
// Supported units: plain number (points), pt, cm, mm, in, px
func ParseRowHeight(s string) (float64, error) {
	val, unit := splitNumberUnit(s)
	if unit == "" || unit == "pt" {
		return val, nil
	}
	switch unit {
	case "cm":
		return cmToPoints(val), nil
	case "mm":
		return cmToPoints(val / 10.0), nil
	case "in":
		return val * 72.0, nil
	case "px":
		return pxToPoints(val), nil
	default:
		// Unknown unit: treat as points
		return val, nil
	}
}

func splitNumberUnit(s string) (float64, string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, ""
	}
	// find first non-number/decimal
	idx := -1
	for i, r := range s {
		if !(r >= '0' && r <= '9') && r != '.' && r != '-' {
			idx = i
			break
		}
	}
	if idx == -1 {
		v, _ := strconv.ParseFloat(s, 64)
		return v, ""
	}
	v, _ := strconv.ParseFloat(strings.TrimSpace(s[:idx]), 64)
	unit := strings.ToLower(strings.TrimSpace(s[idx:]))
	return v, unit
}

func cmToPx(cm float64) float64     { return inchesToPx(cm / 2.54) }
func inchesToPx(in float64) float64 { return in * 96.0 }
func pointsToPx(pt float64) float64 { return (pt / 72.0) * 96.0 }
func pxToPoints(px float64) float64 { return (px / 96.0) * 72.0 }
func cmToPoints(cm float64) float64 { return (cm / 2.54) * 72.0 }

// pxToColWidth converts pixels to Excel column width (characters approx)
func pxToColWidth(px float64) float64 {
	// Excel approximates with padding; ensure non-negative
	w := (px - 5.0) / 7.0
	if w < 0 {
		w = 0
	}
	return math.Round(w*100) / 100 // 2 decimals
}
