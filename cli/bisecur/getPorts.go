package bisecur

import (
	"halsecur/sdk"
)

func GetPorts(localMac [6]byte, mac [6]byte, host string, port int, token uint32) ([]byte, error) {
	var portIds []byte

	return portIds, Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		var err error
		portIds, err = client.GetPorts()
		return err
	})
}

func GetType(localMac [6]byte, mac [6]byte, host string, port int, devicePort byte, token uint32) (byte, error) {
	var portType byte

	return portType, Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		var err error
		portType, err = client.GetType(devicePort)
		return err
	})
}
