package bisecur

import (
	"halsecur/sdk"
)

func AddPort(localMac [6]byte, mac [6]byte, host string, port int, token uint32) (portId byte, err error) {
	return portId, Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		var err2 error
		portId, err2 = client.AddPort()
		return err2
	})
}

func InheritPort(localMac [6]byte, mac [6]byte, host string, port int, token uint32) (portId byte, err error) {
	return portId, Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		var err2 error
		portId, err2 = client.InheritPort()
		return err2
	})
}

func RemovePort(localMac [6]byte, mac [6]byte, host string, port int, devicePort byte, token uint32) error {
	return Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		return client.RemovePort(devicePort)
	})
}
