package bisecur

import (
	"halsecur/sdk"
	"halsecur/sdk/payload"
)

func GetStatus(localMac [6]byte, mac [6]byte, host string, port int, devicePort byte, token uint32) (*payload.HmGetTransitionResponse, error) {
	var status *payload.HmGetTransitionResponse = nil

	return status, Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		var err error
		status, err = client.GetTransition(devicePort)
		return err
	})
}

func SetState(localMac [6]byte, mac [6]byte, host string, port int, devicePort byte, token uint32) error {
	return Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		return client.SetState(devicePort)
	})
}
