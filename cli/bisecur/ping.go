package bisecur

import (
	"halsecur/cli"
	"halsecur/sdk"
	"time"
)

func Ping(localMac [6]byte, mac [6]byte, host string, port int, count int, delay time.Duration, token uint32) error {
	return Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
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
	})
}
