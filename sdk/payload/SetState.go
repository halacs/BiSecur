package payload

import "encoding/hex"

type SetState struct {
	Payload
}

func SetStatePayload() PayloadInterface {
	return &SetState{
		Payload{
			data: []byte{0x00, 0xFF}, // it seems this is always the same
		},
	}
}

func (st *SetState) Encode() []byte {
	data := hex.EncodeToString(st.data)
	return []byte(data)
}

func (st *SetState) String() string {
	return "SetState"
}
