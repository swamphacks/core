package application

import (
	"testing"
	"time"
)

func TestAdmissionWaitlistCutoff(t *testing.T) {
	cases := []struct {
		name string
		at   string
		open bool
	}{
		{"last minute is open", "2026-10-16T03:59:59Z", true},
		{"midnight Eastern is closed", "2026-10-16T04:00:00Z", false},
		{"after cutoff is closed", "2026-10-16T04:00:01Z", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			at, err := time.Parse(time.RFC3339, tc.at)
			if err != nil {
				t.Fatal(err)
			}
			if got := admissionWaitlistIsOpen(at); got != tc.open {
				t.Fatalf("open = %v, want %v", got, tc.open)
			}
		})
	}
}
