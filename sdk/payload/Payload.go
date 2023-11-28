package payload

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type Payload struct {
	//PayloadLength uint16 // Length of the payload - 2 byte
	data       []byte
	dataLength byte
}

func (p *Payload) ToByteArray() []byte {
	return p.data
}

func (p *Payload) Length() int {
	return len(p.ToByteArray()) + binary.Size(p.dataLength) + 1 // TODO why +1 is needed???
}

func (p *Payload) Encode() []byte {
	data := make([]byte, 0)

	p.dataLength = 0
	if p.ToByteArray() != nil && len(p.ToByteArray()) > 0 {
		p.dataLength = byte(len(p.ToByteArray()))
	}

	if p.ToByteArray() != nil && len(p.ToByteArray()) > 0 {
		data = append(data, []byte{p.dataLength}...)
		data = append(data, p.ToByteArray()...)
	}

	return data
}

func getHormanEncodedUsernamePassword(username string, password string) []byte {
	buffer := new(bytes.Buffer)

	usernameLengthStr := fmt.Sprintf("%02X", len(username))
	_, err := buffer.WriteString(usernameLengthStr)
	if err != nil {
		return []byte{}
	}

	usernameHexStr := hex.EncodeToString([]byte(username))
	_, err = buffer.WriteString(usernameHexStr)
	if err != nil {
		return []byte{}
	}

	passwordHexStr := hex.EncodeToString([]byte(password))
	_, err = buffer.WriteString(passwordHexStr)
	if err != nil {
		return []byte{}
	}

	return buffer.Bytes()
}
