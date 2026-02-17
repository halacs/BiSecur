package payload

import (
	"encoding/hex"
	"fmt"
	"halsecur/cli/utils"
)

type GetNameResponse struct {
	Payload
}

func (gnr *GetNameResponse) Encode() []byte {
	return []byte(hex.EncodeToString(gnr.data))
}

func GetNameResponsePayload(name string) PayloadInterface {
	dataLength, err := utils.SafeLen(len(name))
	if err != nil {
		return nil
	}
	return &GetNameResponse{
		Payload{
			data:       []byte(name),
			dataLength: dataLength,
		},
	}
}

func (gnr *GetNameResponse) String() string {
	return fmt.Sprintf("GetNameResponse: %s", gnr.data)
}

func (gnr *GetNameResponse) GetName() string {
	return string(gnr.data)
}

func DecodeGetNameResponsePayload(payloadBytes []byte) (PayloadInterface, error) {
	name := string(payloadBytes)
	gnrp := GetNameResponsePayload(name)
	return gnrp, nil
}
