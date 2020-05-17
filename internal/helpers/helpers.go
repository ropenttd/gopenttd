package helpers

import (
	"fmt"
)

func PackString(input string) (output []byte) {
	output = []byte(input)

	// Zero pad
	output = append(output, byte(0))

	return
}

func GetByteString(byteArray []byte) string {
	return fmt.Sprintf("%x", byteArray)
}
