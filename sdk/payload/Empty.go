package payload

type Empty struct {
	Payload
}

func EmptyPayload() PayloadInterface {
	return &Empty{
		Payload{
			data: []byte{},
		},
	}
}
