package util

import (
	"time"
)

// DateFormat takes the number of days from the OpenTTD Epoch (the 1st of January 0000) and returns a UTC time object.
func DateFormat(date uint32) (t time.Time) {
	// Below decoding stuff gleaned from the source, mostly https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/udp.h#L77
	t, _ = time.ParseInLocation(time.RFC3339, "0000-01-01T00:00:00+00:00", time.UTC)

	// Add the date count to this epoch - note that the golang time library takes into account leap years for us automagically
	t = t.AddDate(0, 0, int(date))

	return t
}
