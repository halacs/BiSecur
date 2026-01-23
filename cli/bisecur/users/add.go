package users

import (
	"halsecur/cli/bisecur"
	"halsecur/sdk"
)

func UserAdd(localMac [6]byte, mac [6]byte, host string, port int, token uint32, userName string, password string) (byte, error) {
	var userId byte
	return userId, bisecur.Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		var err2 error
		userId, err2 = client.AddUser(userName, password)
		return err2
	})
}
