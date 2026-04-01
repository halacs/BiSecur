package payload

import (
	"encoding/hex"
	"fmt"
)

type GetPortsResponse struct {
	Payload
	portIds []byte
}

func (gp *GetPortsResponse) Encode() []byte {
	return []byte(hex.EncodeToString(gp.data))
}

func GetPortsResponsePayload(portIds []byte) PayloadInterface {
	pl := Payload{
		data:       portIds,
		dataLength: byte(len(portIds)), // #nosec G115 intentionally truncate data instead of safe cast
	}

	return &GetPortsResponse{
		Payload: pl,
		portIds: portIds,
	}
}

func (gp *GetPortsResponse) String() string {
	return fmt.Sprintf("GetPortsResponse: %X", gp.data)
}

func (gp *GetPortsResponse) GetPortIds() []byte {
	return gp.portIds
}

func DecodeGetPortsResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	gp := GetPortsResponsePayload(payloadBytes)
	return gp, nil
}
