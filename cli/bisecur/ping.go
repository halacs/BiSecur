package bisecur

import (
	"bisecur/cli"
	"bisecur/sdk"
	"time"
)

func Ping(localMac [6]byte, mac [6]byte, host string, port int, count int, delay time.Duration, token uint32) error {
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

	received := 0
	for i := 0; i < count; i++ {
		sentTimestamp, receivedTimestamp, err := client.Ping()

		if err != nil {
			cli.Log.Errorf("%v", err)
			continue
		}

		received = received + 1
		rtt := receivedTimestamp - sentTimestamp
		cli.Log.Infof("Response %d of %d received in %d ms", received, count, rtt)

		if i < count {
			time.Sleep(delay)
		}
	}

	return nil
}
