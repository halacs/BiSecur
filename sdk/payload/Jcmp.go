package payload

import "encoding/hex"

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

func (j *Jcmp) String() string {
	return "Jcmp"
}
