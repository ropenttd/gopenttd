package gopenttd

import (
	"fmt"
	"time"
)

// Below decoding stuff gleaned from the source, mostly https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/udp.h#L77
// OttdDateFormat takes the number of days from the OpenTTD Epoch (the 1st of January 0000) and returns a UTC time object.
func OttdDateFormat(date uint32) (t time.Time) {
	t, _ = time.ParseInLocation(time.RFC3339, "0000-01-01T00:00:00+00:00", time.UTC)

	// Add the date count to this epoch - note that the golang time library takes into account leap years for us automagically
	t = t.AddDate(0, 0, int(date))

	return t
}

func getByteString(byteArray []byte) string {
	return fmt.Sprintf("%x", byteArray)
}
