package groups

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
	"bisecur/sdk"
)

func ListGroups(localMac [6]byte, mac [6]byte, host string, port int, token uint32) (*sdk.Groups, error) {
	client := sdk.NewClient(cli.Log, localMac, mac, host, port, token)
	err := client.Open()
	if err != nil {
		return nil, err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			cli.Log.Errorf("%v", err)
		}
	}()

	var groups *sdk.Groups
	err = bisecur.Retry(func() error {
		var err2 error
		groups, err2 = client.GetGroups()
		return err2
	})

	if err != nil {
		return nil, err
	}

	return groups, nil
}
