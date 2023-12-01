package payload

import "encoding/hex"

type HmGetTransition struct {
	Payload
}

func HmGetTransitionPayload(data byte) PayloadInterface { // TODO change name of data variable to something more meaningful
	return &HmGetTransition{
		Payload{
			data: []byte{data},
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
