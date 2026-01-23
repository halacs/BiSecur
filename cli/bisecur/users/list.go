package users

import (
	"halsecur/cli"
	"halsecur/cli/bisecur"
	"halsecur/sdk"
)

func ListUsers(localMac [6]byte, mac [6]byte, host string, port int, token uint32) error {
	var users *sdk.Users
	err := bisecur.Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		var err2 error
		users, err2 = client.GetUsers()
		return err2
	})

	if err != nil {
		return err
	}

	cli.Log.WithField("users", users).Infof("Success")

	return nil
}
