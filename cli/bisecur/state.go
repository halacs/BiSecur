package bisecur

import (
	"bisecur/cli"
	"bisecur/sdk"
	"bisecur/sdk/payload"
)

func GetStatus(localMac [6]byte, mac [6]byte, host string, port int, devicePort byte, token uint32) (*payload.HmGetTransitionResponse, error) {
	client := sdk.NewClient(cli.Log, localMac, mac, host, port, token)
	err := client.Open()
	if err != nil {
		return nil, err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			cli.Log.Errorf("%v", err2)
		}
	}()

	var status *payload.HmGetTransitionResponse
	err = Retry(func() error {
		var err2 error
		status, err2 = client.GetTransition(devicePort)
		return err2
	})

	if err != nil {
		return nil, err
	}

	return status, nil
}

func SetState(localMac [6]byte, mac [6]byte, host string, port int, devicePort byte, token uint32) error {
	client := sdk.NewClient(cli.Log, localMac, mac, host, port, token)
	err := client.Open()
	if err != nil {
		return err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			cli.Log.Errorf("%v", err2)
		}
	}()

	err = Retry(func() error {
		err2 := client.SetState(devicePort)
		return err2
	})

	if err != nil {
		return err
	}

	return nil
}
