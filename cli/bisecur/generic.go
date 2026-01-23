package bisecur

import (
	"halsecur/cli"
	"halsecur/cli/utils"
	"halsecur/sdk"
	"time"
)

func GenericWithRetryAlways(localMac [6]byte, mac [6]byte, host string, port int, token uint32, retryCount int, f func(c *sdk.Client) error) error {
	return Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		return utils.RetryAlways(retryCount, func() error {
			return f(client)
		})
	})
}
func GenericWithRetry(localMac [6]byte, mac [6]byte, host string, port int, token uint32, retryCount int, delayIfError time.Duration, isRetriableError func(error) bool, f func(c *sdk.Client) error) error {
	return Generic(localMac, mac, host, port, token, func(client *sdk.Client) error {
		return utils.Retry(retryCount, delayIfError, func() error {
			return f(client)
		}, isRetriableError)
	})
}

func Generic(localMac [6]byte, mac [6]byte, host string, port int, token uint32, f func(client *sdk.Client) error) error {
	client := sdk.NewClient(cli.Log, localMac, mac, host, port, token)
	defer func() {
		err2 := client.Close()
		if err2 != nil {
			cli.Log.Errorf("%v", err2)
		}
	}()

	err := client.Open()
	if err != nil {
		return err
	}

	err = f(client)

	if err != nil {
		return err
	}

	return nil
}
