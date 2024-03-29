package payload

import "encoding/hex"

type HmGetTransition struct {
	Payload
}

func HmGetTransitionPayload(portID byte) PayloadInterface {
	return &HmGetTransition{
		Payload{
			data: []byte{portID},
		},
	}
}

func (hgt *HmGetTransition) Encode() []byte {
	encodedBytes := make([]byte, 2)
	hex.Encode(encodedBytes, hgt.ToByteArray())
	return encodedBytes
}

func (hgt *HmGetTransition) String() string {
	return "HmGetTransition"
}
