package payload

import "encoding/hex"

func MockPayload(dataHex string) Payload {
	data, err := hex.DecodeString(dataHex)
	if err != nil {
		panic(err)
	}

	return Payload{
		data: data,
	}
}
