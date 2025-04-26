package users

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
	"bisecur/sdk"
)

func ListUsers(localMac [6]byte, mac [6]byte, host string, port int, token uint32) error {
	client := sdk.NewClient(cli.Log, localMac, mac, host, port, token)
	err := client.Open()
	if err != nil {
		return err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			cli.Log.Errorf("%v", err)
		}
	}()

	var users *sdk.Users
	err = bisecur.Retry(func() error {
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
