package bisecur

import (
	"bisecur/cli"
	"bisecur/sdk"
)

func UserPasswordChange(localMac [6]byte, mac [6]byte, host string, port int, token uint32, userId byte, newPassword string) error {
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

	err = Retry(func() error {
		err2 := client.PasswordChange(userId, newPassword)
		return err2
	})

	if err != nil {
		return err
	}

	return nil
}
