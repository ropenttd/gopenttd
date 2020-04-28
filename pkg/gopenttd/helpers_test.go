package gopenttd

import (
	"testing"
	"time"
)

func TestDateDecode(t *testing.T) {
	// 01-01-1940
	date := OttdDateFormat(708570)
	expect := time.Date(1940, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	if date != expect {
		t.Errorf("Date calculation was incorrect, got: %s expected: %s", date, expect)
	}
}
