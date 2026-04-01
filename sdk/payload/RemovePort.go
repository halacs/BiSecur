package payload

import (
	"encoding/hex"
	"fmt"
)

type RemovePort struct {
	Payload
	portId byte
}

func (rp *RemovePort) Encode() []byte {
	return []byte(hex.EncodeToString(rp.data))
}

func RemovePortPayload(portId byte) PayloadInterface {
	pl := Payload{
		data:       []byte{portId},
		dataLength: byte(1),
	}

	return &RemovePort{
		Payload: pl,
		portId:  portId,
	}
}

func (rp *RemovePort) String() string {
	return fmt.Sprintf("RemovePort: %s", rp.data)
}

func (rp *RemovePort) GetPortId() byte {
	return rp.portId
}

func DecodeRemovePortPayload(payloadBytes []byte) (PayloadInterface, error) {
	if len(payloadBytes) != 1 {
		return nil, fmt.Errorf("invalid payload length for RemovePortPayload: %d", len(payloadBytes))
	}

	rp := RemovePortPayload(payloadBytes[0])
	return rp, nil
}
