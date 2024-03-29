package payload

import (
	"encoding/hex"
	"fmt"
)

type SetUserRightsResponse struct {
	Payload
	userId byte
	stg    byte // TODO find out what is this field
}

func (surr *SetUserRightsResponse) Encode() []byte {
	return []byte(hex.EncodeToString(surr.data))
}

func SetUserRightsResponsePayload(userId byte, stg byte) PayloadInterface {
	pl := Payload{
		data:       []byte{userId, stg},
		dataLength: byte(2),
	}

	return &SetUserRightsResponse{
		Payload: pl,
		userId:  userId,
		stg:     stg,
	}
}

func (surr *SetUserRightsResponse) String() string {
	return fmt.Sprintf("SetUserRightsResponse: %s", surr.data)
}

func (surr *SetUserRightsResponse) GetUserId() byte {
	return surr.userId
}

func (surr *SetUserRightsResponse) GetSomething() byte {
	return surr.stg
}

func DecodeSetUserRightsResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	if len(payloadBytes) != 2 {
		return nil, fmt.Errorf("invalid payload length for SetUserRightsResponsePayload: %d", len(payloadBytes))
	}

	gurr := SetUserRightsResponsePayload(payloadBytes[0], payloadBytes[1])
	return gurr, nil
}
