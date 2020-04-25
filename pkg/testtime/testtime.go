package testtime

import (
	"testing"
	"time"
)

// Parse will parse a time value in a given format while handling any errors in the test
func Parse(t *testing.T, layout, value string) time.Time {
	v, err := time.Parse(layout, value)
	if err != nil {
		t.Errorf("failed to parse time layout: %s value: %s", layout, value)
	}

	return v
}

// ParseInLocation will parse a time value in a given format and location while handling any errors in the test
func ParseInLocation(t *testing.T, layout, value string, loc *time.Location) time.Time {
	v, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		t.Errorf("failed to parse time layout: %s value: %s, location: %s", layout, value, loc.String())
	}

	return v
}
