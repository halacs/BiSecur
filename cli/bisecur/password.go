package bisecur

import (
	"halsecur/sdk"
)

func UserPasswordChange(localMac [6]byte, mac [6]byte, host string, port int, token uint32, userId byte, newPassword string) error {
	return Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		return client.PasswordChange(userId, newPassword)
	})
}
