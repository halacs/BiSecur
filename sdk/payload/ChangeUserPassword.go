package payload

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type ChangeUserPassword struct {
	Payload
	userId   byte
	password string
}

func (cup *ChangeUserPassword) Encode() []byte {
	return []byte(hex.EncodeToString(cup.data))
}

func ChangeUserPasswordPayload(userId byte, password string) PayloadInterface {
	var buff bytes.Buffer
	buff.Write([]byte{userId})
	buff.Write([]byte(password))

	pl := Payload{
		data:       buff.Bytes(),
		dataLength: byte(len(buff.Bytes())),
	}

	return &ChangeUserPassword{
		Payload:  pl,
		userId:   userId,
		password: password,
	}
}

func (cup *ChangeUserPassword) String() string {
	return fmt.Sprintf("ChangeUserPassword: %s", cup.data)
}

func (cup *ChangeUserPassword) GetUserId() byte {
	return cup.userId
}

func (cup *ChangeUserPassword) GetPassword() string {
	return cup.password
}

func DecodeChangeUserPasswordPayload(payloadBytes []byte) (PayloadInterface, error) {
	cup := ChangeUserPasswordPayload(payloadBytes[0], string(payloadBytes[1:]))
	return cup, nil
}
