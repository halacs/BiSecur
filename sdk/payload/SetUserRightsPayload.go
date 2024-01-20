package payload

import (
	"encoding/hex"
	"fmt"
)

type SetUserRights struct {
	Payload
	userId byte
}

func (sur *SetUserRights) Encode() []byte {
	return []byte(hex.EncodeToString(sur.data))
}

func SetUserRightsPayload(userId byte) PayloadInterface {
	pl := Payload{
		data:       []byte{userId},
		dataLength: byte(1),
	}

	return &SetUserRights{
		Payload: pl,
		userId:  userId,
	}
}

func (sur *SetUserRights) String() string {
	return fmt.Sprintf("SetUserRights: %s", sur.data)
}

func (sur *SetUserRights) GetUserId() byte {
	return sur.userId
}

func DecodeSetUserRightsPayload(payloadBytes []byte) (PayloadInterface, error) {
	if len(payloadBytes) != 1 {
		return nil, fmt.Errorf("invalid payload length for SetUserRightsPayload: %d", len(payloadBytes))
	}

	sur := SetUserRightsPayload(payloadBytes[0])
	return sur, nil
}
