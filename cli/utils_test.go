package cli

import (
	"reflect"
	"testing"
)

func TestTransmissionContainerEncode(t *testing.T) {
	expectedMac := [6]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0xef}
	mac, err := ParesMacString("00:01:02:03:04:ef")
	if err != nil {
		t.Logf("Unexpected error: %v", err)
		t.Fail()
	}

	if !reflect.DeepEqual(mac, expectedMac) {
		t.Logf("Expected MAC: %X, Actual MAC: %X", expectedMac, mac)
		t.Fail()
	}
}
