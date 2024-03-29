package payload

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type LoginResponse struct {
	Payload
	senderID byte // TODO what is the purpose of this field?
	token    uint32
}

func DecodeLoginResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	id := payloadBytes[0]
	token := binary.BigEndian.Uint32(payloadBytes[1:])

	lrp := LoginResponsePayload(id, token)

	return lrp, nil
}

func (l *LoginResponse) Encode() []byte {
	data := hex.EncodeToString(l.ToByteArray())
	return []byte(data)
}

func (l *LoginResponse) ToByteArray() []byte {
	return toByteArray(l.senderID, l.token)
}

func toByteArray(senderID byte, token uint32) []byte {
	var buffer = new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, senderID)
	if err != nil {
		panic(err)
	}

	err = binary.Write(buffer, binary.BigEndian, token)
	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func LoginResponsePayload(senderID byte, token uint32) PayloadInterface {
	lrp := &LoginResponse{
		Payload: Payload{
			data:       toByteArray(senderID, token),
			dataLength: byte(len(toByteArray(senderID, token))),
		},
		senderID: senderID,
		token:    token,
	}
	return lrp
}

func (l *LoginResponse) String() string {
	return fmt.Sprintf("SenderID: 0x%X, Token: 0x%X", l.senderID, l.token)
}

func (l *LoginResponse) GetToken() uint32 {
	return l.token
}

func (l *LoginResponse) GetSenderID() byte {
	return l.senderID
}
