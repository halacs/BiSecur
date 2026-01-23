package bisecur

import (
	"halsecur/sdk"
)

func GetName(localMac, mac [6]byte, host string, port int, token uint32) (string, error) {
	var name string
	return name, Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		var err2 error
		name, err2 = client.GetName()
		return err2
	})
}
