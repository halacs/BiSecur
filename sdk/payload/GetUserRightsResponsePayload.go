package payload

import (
	"encoding/hex"
	"fmt"
)

type GetUserRightsResponse struct {
	Payload
	d byte // TODO what is this exactly?
}

func (gurr *GetUserRightsResponse) Encode() []byte {
	return []byte(hex.EncodeToString(gurr.data))
}

func GetUserRightsResponsePayload(d byte) PayloadInterface {
	pl := Payload{
		data:       []byte{d},
		dataLength: byte(1),
	}

	return &GetUserRightsResponse{
		Payload: pl,
		d:       d,
	}
}

func (gurr *GetUserRightsResponse) String() string {
	return fmt.Sprintf("GetUserRightsResponse: %s", gurr.data)
}

func (gurr *GetUserRightsResponse) GetD() byte {
	return gurr.d
}

func DecodeGetUserRightsResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	if len(payloadBytes) != 1 {
		return nil, fmt.Errorf("invalid payload length for GetUserRightsResponsePayload: %d", len(payloadBytes))
	}

	gurr := GetUserRightsResponsePayload(payloadBytes[0])
	return gurr, nil
}
