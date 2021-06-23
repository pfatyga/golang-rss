package main

import "testing"

func TestSanitizeTitle(t *testing.T) {
	test := "Kats: One/ stop shop for time series analysis in Python"
	want := "Kats--One--stop-shop-for-time-series-analysis-in-Python"
	if got := sanitizeTitle(test); got != want {
		t.Errorf("sanitizeTitle() = %q, want %q", got, want)
	}
}
