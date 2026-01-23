package bisecur

import (
	"halsecur/sdk"
)

func Login(localMac [6]byte, mac [6]byte, host string, port int, username string, password string) (uint32, error) {
	var token uint32
	return token, Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		err2 := client.Login(username, password)
		token = client.GetToken()
		return err2
	})
}
