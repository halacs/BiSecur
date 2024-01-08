package cmd

func retry(f func() error) error {
	var err error

	for i := 0; i < retryCount; i++ {
		err = f()
		if err == nil {
			break
		}
		log.Debugf("Retriable error: %v", err)
	}

	return err
}
