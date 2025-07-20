package bisecur

import (
	"bisecur/cli"
	"bisecur/sdk"
)

func Login(localMac [6]byte, mac [6]byte, host string, port int, username string, password string) (uint32, error) {
	client := sdk.NewClient(cli.Log, localMac, mac, host, port, 0)
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

	err = client.Login(username, password)
	if err != nil {
		return 0, err
	}

	return client.GetToken(), nil
}
