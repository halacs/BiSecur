package cli

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func ParesMacString(mac string) ([6]byte, error) {
	bytes := strings.Split(mac, ":")
	if len(bytes) != 6 {
		return [6]byte{}, fmt.Errorf("invalid mac address length: %s", mac)
	}

	clearedMacStr := strings.ReplaceAll(mac, ":", "")
	macBytes, err := hex.DecodeString(clearedMacStr)
	return [6]byte{macBytes[0], macBytes[1], macBytes[2], macBytes[3], macBytes[4], macBytes[5]}, err
}
