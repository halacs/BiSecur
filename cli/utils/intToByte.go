package utils

import (
	"fmt"
	"math"
)

func SafeLen(i int) (byte, error) {
	// G115 Fix: Ensure length fits in byte before casting
	if i > math.MaxUint8 {
		// You must decide how to handle this. Returning an error is usually best.
		return 0, fmt.Errorf("length overflow. Cannot cast to byte from int. Value: %d", i)
	}
	return byte(i), nil // #nosec G115 safe after validating the value
}
