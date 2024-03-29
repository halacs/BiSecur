package payload

import (
	"encoding/hex"
	"fmt"
)

type RemoveUser struct {
	Payload
	userId byte
}

func (ru *RemoveUser) Encode() []byte {
	return []byte(hex.EncodeToString(ru.data))
}

func RemoveUserPayload(userId byte) PayloadInterface {
	pl := Payload{
		data:       []byte{userId},
		dataLength: byte(1),
	}

	return &RemoveUser{
		Payload: pl,
		userId:  userId,
	}
}

func (ru *RemoveUser) String() string {
	return fmt.Sprintf("RemoveUser: %s", ru.data)
}

func (ru *RemoveUser) GetUserId() byte {
	return ru.userId
}

func DecodeRemoveUserPayload(payloadBytes []byte) (PayloadInterface, error) {
	if len(payloadBytes) != 1 {
		return nil, fmt.Errorf("invalid payload length for RemoveUserPayload: %d", len(payloadBytes))
	}

	ru := RemoveUserPayload(payloadBytes[0])
	return ru, nil
}
