package retry

import (
	"time"
)

func RetryFunc(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return RetryFunc(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}

func RetryFuncWithCode(attempts int, sleep time.Duration, fn func() (bool, int, error)) (int, error) {
	if retry, code, err := fn(); err != nil {
		if retry {
			if attempts--; attempts > 0 {
				time.Sleep(sleep)
				return RetryFuncWithCode(attempts, 2*sleep, fn)
			}
			return code, err
		} else {
			return code, err
		}
	}
	return 0, nil
}
