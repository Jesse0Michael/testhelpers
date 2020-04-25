package testtime

import (
	"fmt"
	"testing"
	"time"
)

func ExampleParse() {
	t := &testing.T{}

	fmt.Println(Parse(t, time.RFC3339, "1990-06-04T01:02:03Z"))
	// Output:
	// 1990-06-04 01:02:03 +0000 UTC
}

func ExampleParseInLocation() {
	t := &testing.T{}

	fmt.Println(ParseInLocation(t, time.RFC3339, "1990-06-04T01:02:03Z", time.UTC))
	// Output:
	// 1990-06-04 01:02:03 +0000 UTC
}

func TestParse_Failure(t *testing.T) {
	tt := &testing.T{}
	Parse(tt, time.RFC3339, "June 4 1990")
	if !tt.Failed() {
		t.Error("expected Parse() to fail to parse time")
	}
}

func TestParseInLocation_Failure(t *testing.T) {
	tt := &testing.T{}
	ParseInLocation(tt, time.RFC3339, "June 4 1990", time.UTC)
	if !tt.Failed() {
		t.Error("expected ParseInLocation() to fail to parse time")
	}
}
