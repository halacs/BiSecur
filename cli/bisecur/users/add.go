package users

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
	"bisecur/sdk"
)

func UserAdd(localMac [6]byte, mac [6]byte, host string, port int, token uint32, userName string, password string) (byte, error) {
	client := sdk.NewClient(cli.Log, localMac, mac, host, port, token)
	err := client.Open()
	if err != nil {
		return 0, err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			cli.Log.Errorf("%v", err)
		}
	}()

	var userId byte
	err = bisecur.Retry(func() error {
		var err2 error
		userId, err2 = client.AddUser(userName, password)
		return err2
	})

	if err != nil {
		return 0, err
	}

	return userId, nil
}
