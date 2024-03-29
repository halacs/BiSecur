package payload

import (
	"encoding/hex"
	"fmt"
)

type RemoveUserResponse struct {
	Payload
	userId byte
}

func (ru *RemoveUserResponse) Encode() []byte {
	return []byte(hex.EncodeToString(ru.data))
}

func RemoveUserResponsePayload(userId byte) PayloadInterface {
	pl := Payload{
		data:       []byte{userId},
		dataLength: byte(1),
	}

	return &RemoveUserResponse{
		Payload: pl,
		userId:  userId,
	}
}

func (ru *RemoveUserResponse) String() string {
	return fmt.Sprintf("RemoveUserResponse: %s", ru.data)
}

func (ru *RemoveUserResponse) GetUserId() byte {
	return ru.userId
}

func DecodeRemoveUserResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	if len(payloadBytes) != 1 {
		return nil, fmt.Errorf("invalid payload length for RemoveUserResponsePayload: %d", len(payloadBytes))
	}

	ru := RemoveUserResponsePayload(payloadBytes[0])
	return ru, nil
}
