package utils

import "testing"

func TestFormatBudgetRange(t *testing.T) {
	f := func(v float64) *float64 { return &v }
	cases := []struct {
		min, max *float64
		want     string
	}{
		{f(5000), f(10000), "ETB 5,000 – 10,000"},
		{f(1250000), nil, "ETB 1,250,000+"},
		{nil, f(999), "Up to ETB 999"},
		{nil, nil, ""},
	}
	for _, c := range cases {
		if got := FormatBudgetRange(c.min, c.max); got != c.want {
			t.Errorf("got %q, want %q", got, c.want)
		}
	}
}
