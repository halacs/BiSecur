package payload

type Received struct {
	Payload
}

func ReceivedPayload(receivedData []byte) PayloadInterface {
	return &Received{
		Payload{
			data: receivedData,
		},
	}
}
