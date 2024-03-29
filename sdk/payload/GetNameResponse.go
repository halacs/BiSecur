package payload

import (
	"encoding/hex"
	"fmt"
)

type GetNameResponse struct {
	Payload
}

func (gnr *GetNameResponse) Encode() []byte {
	return []byte(hex.EncodeToString(gnr.data))
}

func GetNameResponsePayload(name string) PayloadInterface {
	return &GetNameResponse{
		Payload{
			data:       []byte(name),
			dataLength: byte(len(name)),
		},
	}
}

func (gnr *GetNameResponse) String() string {
	return fmt.Sprintf("GetNameResponse: %s", gnr.data)
}

func (gnr *GetNameResponse) GetName() string {
	return string(gnr.data)
}

func DecodeGetNameResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	name := string(payloadBytes)
	gnrp := GetNameResponsePayload(name)
	return gnrp, nil
}
