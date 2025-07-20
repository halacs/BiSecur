package bisecur

import (
	"bisecur/cli"
	"bisecur/sdk"
	"fmt"
)

func Logout(localMac [6]byte, mac [6]byte, host string, port int, token uint32) error {
	if token == 0 {
		return fmt.Errorf("invalid token value: 0x%X", token)
	}

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

	client.SetToken(token)

	err = client.Logout()
	if err != nil {
		return err
	}

	return nil
}
