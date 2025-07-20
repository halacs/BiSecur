package bisecur

import (
	"bisecur/cli"
)

func Retry(f func() error) error {
	var err error

	for i := 0; i < retryCount; i++ {
		err = f()
		if err == nil {
			break
		}
		cli.Log.Debugf("Retriable error: %v", err)
	}

	return err
}
