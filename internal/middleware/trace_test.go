package middleware

import "testing"

func TestParseTraceID(t *testing.T) {
	traceID := "4bf92f3577b34da6a3ce929d0e0e4736"
	got := parseTraceID("00-" + traceID + "-00f067aa0ba902b7-01")
	if got != traceID {
		t.Fatalf("expected %s, got %s", traceID, got)
	}
}
