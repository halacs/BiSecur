package payload

import (
	"encoding/hex"
	"fmt"
)

type GetType struct {
	Payload
	portId byte
}

func (gt *GetType) Encode() []byte {
	return []byte(hex.EncodeToString(gt.data))
}

func GetTypePayload(portId byte) PayloadInterface {
	pl := Payload{
		data:       []byte{portId},
		dataLength: byte(1),
	}

	return &GetType{
		Payload: pl,
		portId:  portId,
	}
}

func (gt *GetType) String() string {
	return fmt.Sprintf("GetType: %X", gt.data)
}

func DecodeGetTypePayload(payloadBytes []byte) (PayloadInterface, error) {
	if len(payloadBytes) != 1 {
		return nil, fmt.Errorf("invalid payload length for GetTypePayload: %d", len(payloadBytes))
	}

	gt := GetTypePayload(payloadBytes[0])
	return gt, nil
}
