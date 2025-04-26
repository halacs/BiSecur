package bisecur

import (
	"bisecur/cli"
	"bisecur/sdk"
)

func GetName(localMac, mac [6]byte, host string, port int, token uint32) (string, error) {
	client := sdk.NewClient(cli.Log, localMac, mac, host, port, token)
	err := client.Open()
	if err != nil {
		return "", err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			cli.Log.Fatalf("%v", err2)
		}
	}()

	var name string
	err = Retry(func() error {
		var err2 error
		name, err2 = client.GetName()
		return err2
	})

	if err != nil {
		return "", err
	}

	return name, nil
}
