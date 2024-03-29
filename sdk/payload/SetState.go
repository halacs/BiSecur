package payload

import "encoding/hex"

type SetState struct {
	Payload
}

func SetStatePayload(portID byte) PayloadInterface {
	var state byte = 0xFF // it seems this is always the same
	return &SetState{
		Payload{
			data: []byte{portID, state},
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
