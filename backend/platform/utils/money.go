package utils

import (
	"math"
	"strconv"
	"strings"
)

// FormatETB writes a whole-birr amount with thousands separators: 12500 → "12,500".
func FormatETB(v float64) string {
	n := int64(math.Round(v))
	neg := n < 0
	if neg {
		n = -n
	}
	s := strconv.FormatInt(n, 10)
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	if neg {
		return "-" + b.String()
	}
	return b.String()
}

// FormatBudgetRange is the one way budgets read in the app:
// "ETB 5,000 – 10,000", "ETB 5,000+", "Up to ETB 10,000", or "" when unknown.
func FormatBudgetRange(min, max *float64) string {
	switch {
	case min == nil && max == nil:
		return ""
	case min != nil && max != nil:
		return "ETB " + FormatETB(*min) + " – " + FormatETB(*max)
	case min != nil:
		return "ETB " + FormatETB(*min) + "+"
	default:
		return "Up to ETB " + FormatETB(*max)
	}
}
