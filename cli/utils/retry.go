package utils

import (
	"context"
	"halsecur/cli"
	"time"
)

func RetryAlways(retryCount int, f func() error) error {
	return Retry(context.Background(), retryCount, 0, f, func(err error) bool {
		return true // all errors are retriable
	})
}

func RetryAlwaysWithContext(ctx context.Context, retryCount int, f func() error) error {
	isRetriableErrorFunc := func(error) bool {
		if ctx.Err() != nil {
			// The context has been canceled or the deadline has passed.
			return false
		}

		// The context is still active so otherwise all errors are retriable by definition.
		return true
	}

	return Retry(ctx, retryCount, 0, f, isRetriableErrorFunc)
}

func RetryAlwaysWithDelay(retryCount int, delay time.Duration, f func() error) error {
	return Retry(context.Background(), retryCount, delay, f, func(err error) bool {
		return true // all errors are retriable
	})
}

func Retry(ctx context.Context, retryCount int, delay time.Duration, f func() error, isRetriableError func(error) bool) error {
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
