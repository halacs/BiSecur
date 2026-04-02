package payload

import (
	"encoding/hex"
	"fmt"
)

type GetTypeResponse struct {
	Payload
	portId   byte
	portType byte
}

func (gt *GetTypeResponse) Encode() []byte {
	return []byte(hex.EncodeToString(gt.data))
}

func GetTypeResponsePayload(portId byte, portType byte) PayloadInterface {
	pl := Payload{
		data:       []byte{portId, portType},
		dataLength: byte(2),
	}

	return &GetTypeResponse{
		Payload:  pl,
		portId:   portId,
		portType: portType,
	}
}

func (gt *GetTypeResponse) String() string {
	return fmt.Sprintf("GetTypeResponse: %X", gt.data)
}

func (gt *GetTypeResponse) GetPortId() byte {
	return gt.portId
}

func (gt *GetTypeResponse) GetPortType() byte {
	return gt.portType
}

func DecodeGetTypeResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	if len(payloadBytes) != 2 {
		return nil, fmt.Errorf("invalid payload length for GetTypeResponsePayload: %d", len(payloadBytes))
	}

	gt := GetTypeResponsePayload(payloadBytes[0], payloadBytes[1])
	return gt, nil
}
