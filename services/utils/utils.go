package utils

import logger "github.com/likecoin/likechain/services/log"

var log = logger.L

// RetryIfPanic retries the operation if there is a panic
func RetryIfPanic(retryCount uint, f func()) {
	for retryCount > 0 {
		success := false
		func() {
			retryCount--
			defer func() {
				if err := recover(); err != nil {
					if retryCount == 0 {
						log.WithField("panic_value", err).Panic("Panic retry limit exceeded")
					} else {
						log.
							WithField("remaining_retry_count", retryCount).
							WithField("panic_value", err).
							Error("Caught panic, retrying")
					}
				}
			}()
			f()
			success = true
		}()
		if success {
			break
		}
	}
}
