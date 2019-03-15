package utils

import logger "github.com/likecoin/likechain/services/log"

var log = logger.L

// RetryIfPanic retries the operation if there is a panic
func RetryIfPanic(retryCount uint, f func()) {
	panicCount := uint(0)
	for {
		success := false
		func() {
			defer func() {
				if err := recover(); err != nil {
					panicCount++
					if panicCount > retryCount {
						log.WithField("panic_value", err).Panic("Panic retry limit exceeded")
					} else {
						log.
							WithField("panic_count", panicCount).
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
