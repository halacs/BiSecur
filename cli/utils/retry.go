package utils

import (
	"halsecur/cli"
	"time"
)

func RetryAlways(retryCount int, f func() error) error {
	return Retry(retryCount, 0, f, func(err error) bool {
		return true // all errors are retriable
	})
}

func RetryAlwaysWithDelay(retryCount int, delay time.Duration, f func() error) error {
	return Retry(retryCount, delay, f, func(err error) bool {
		return true // all errors are retriable
	})
}

func Retry(retryCount int, delay time.Duration, f func() error, isRetriableError func(error) bool) error {
	var err error

	for i := 0; i < retryCount; i++ {
		err = f()
		if err == nil {
			break
		}
		if !isRetriableError(err) {
			cli.Log.Errorf("Not retriable error: %v", err)
			break
		}
		cli.Log.Warnf("Retriable error: %v", err)

		time.Sleep(delay)
	}

	return err
}
