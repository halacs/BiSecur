package payload

import (
	"encoding/hex"
	"fmt"
)

type PortIdResponse struct {
	Payload
	portId byte
}

func (pr *PortIdResponse) Encode() []byte {
	return []byte(hex.EncodeToString(pr.data))
}

func PortIdResponsePayload(portId byte) PayloadInterface {
	pl := Payload{
		data:       []byte{portId},
		dataLength: byte(1),
	}

	return &PortIdResponse{
		Payload: pl,
		portId:  portId,
	}
}

func (pr *PortIdResponse) String() string {
	return fmt.Sprintf("PortIdResponse: %s", pr.data)
}

func (pr *PortIdResponse) GetPortId() byte {
	return pr.portId
}

func DecodePortIdResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	if len(payloadBytes) != 1 {
		return nil, fmt.Errorf("invalid payload length for PortIdResponsePayload: %d", len(payloadBytes))
	}

	pr := PortIdResponsePayload(payloadBytes[0])
	return pr, nil
}
