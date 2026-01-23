package users

import (
	"halsecur/cli/bisecur"
	"halsecur/sdk"
)

func UserRemove(localMac [6]byte, mac [6]byte, host string, port int, token uint32, userId byte) error {
	return bisecur.Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		return client.RemoveUser(userId)
	})
}
