package payload

type GetName struct {
	Payload
}

func GetNamePayload() PayloadInterface {
	return &GetName{
		Payload{
			data: []byte{},
		},
	}
}

func (gn *GetName) String() string {
	return "GetName"
}
