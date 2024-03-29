package payload

import (
	"encoding/hex"
	"fmt"
)

type Jcmp struct {
	Payload
}

func JcmpPayload(json string) PayloadInterface {
	return &Jcmp{
		Payload{
			data: []byte(json),
		},
	}
}

func (j *Jcmp) Encode() []byte {
	data := hex.EncodeToString(j.data)
	return []byte(data)
}

func DecodeJcmpPayload(payloadBytes []byte) (PayloadInterface, error) {
	return &Jcmp{
		Payload{
			data:       payloadBytes, // json request/response
			dataLength: byte(len(payloadBytes)),
		},
	}, nil
}

func (j *Jcmp) String() string {
	return fmt.Sprintf("Jcmp: %s", j.ToByteArray())
}
