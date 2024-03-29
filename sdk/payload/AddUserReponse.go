package payload

import (
	"encoding/hex"
	"fmt"
)

type AddUserResponse struct {
	Payload
	userId byte
}

func (aur *AddUserResponse) Encode() []byte {
	return []byte(hex.EncodeToString(aur.data))
}

func AddUserResponsePayload(userId byte) PayloadInterface {
	pl := Payload{
		data:       []byte{userId},
		dataLength: byte(1),
	}

	return &AddUserResponse{
		Payload: pl,
		userId:  userId,
	}
}

func (aur *AddUserResponse) String() string {
	return fmt.Sprintf("AddUserResponse: %s", aur.data)
}

func (aur *AddUserResponse) GetUserId() byte {
	return aur.userId
}

func DecodeAddUserResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	if len(payloadBytes) != 1 {
		return nil, fmt.Errorf("invalid payload length for AddUserResponsePayload: %d", len(payloadBytes))
	}

	aurp := AddUserResponsePayload(payloadBytes[0])
	return aurp, nil
}
