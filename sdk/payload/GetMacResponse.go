package payload

import (
	"encoding/hex"
	"fmt"
)

type GetMac struct {
	Payload
}

func GetMacResponsePayload(macAddress [6]byte) PayloadInterface {
	gm := &GetMac{
		Payload: Payload{
			data:       macAddress[:],
			dataLength: byte(len(macAddress)),
		},
	}
	return gm
}

func (gm *GetMac) String() string {
	return fmt.Sprintf("GetMac: %X", gm.data)
}

func DecodeGetMacResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	macAddress := [6]byte{
		payloadBytes[0], payloadBytes[1], payloadBytes[2],
		payloadBytes[3], payloadBytes[4], payloadBytes[5],
	}
	lrp := GetMacResponsePayload(macAddress)
	return lrp, nil
}

func (gm *GetMac) Encode() []byte {
	hexStr := hex.EncodeToString(gm.data)
	return []byte(hexStr)
}

func (gm *GetMac) GetMac() [6]byte {
	macAddress := [6]byte{
		gm.data[0], gm.data[1], gm.data[2],
		gm.data[3], gm.data[4], gm.data[5],
	}
	return macAddress
}
