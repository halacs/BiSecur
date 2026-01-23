package bisecur

import (
	"halsecur/sdk"
)

func Logout(localMac [6]byte, mac [6]byte, host string, port int, token uint32) error {
	if token == 0 {
		return nil //fmt.Errorf("invalid token value: 0x%X. Logout request ignored", token)
	}

	return Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		err := client.Logout()
		if err == nil {
			client.SetToken(0) // clear token in local client instance
		}
		return err
	})
}
