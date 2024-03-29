package payload

import (
	"encoding/hex"
	"fmt"
)

type GetUserRights struct {
	Payload
	userId byte
}

func (gur *GetUserRights) Encode() []byte {
	return []byte(hex.EncodeToString(gur.data))
}

func GetUserRightsPayload(userId byte) PayloadInterface {
	pl := Payload{
		data:       []byte{userId},
		dataLength: byte(1),
	}

	return &GetUserRights{
		Payload: pl,
		userId:  userId,
	}
}

func (gur *GetUserRights) String() string {
	return fmt.Sprintf("GetUserRights: %s", gur.data)
}

func (gur *GetUserRights) GetUserId() byte {
	return gur.userId
}

func DecodeGetUserRightsPayload(payloadBytes []byte) (PayloadInterface, error) {
	if len(payloadBytes) != 1 {
		return nil, fmt.Errorf("invalid payload length for AddUserResponsePayload: %d", len(payloadBytes))
	}

	gurp := GetUserRightsPayload(payloadBytes[0])
	return gurp, nil
}
