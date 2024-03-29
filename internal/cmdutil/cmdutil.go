package cmdutil

import (
	"github.com/avast/retry-go/v4"
	"github.com/conplementag/cops-hq/v2/pkg/error_handling"
	"github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

// ExecuteFunctionWithRetry - reruns a function in case of error and logs error
func ExecuteFunctionWithRetry(function func() error, maxAttempts uint) error {
	if error_handling.PanicOnAnyError {
		defer func() {
			logrus.Infof("Reenable errorhandling panic on error")
			error_handling.PanicOnAnyError = true
		}()

		logrus.Info("Disable error handling panic on error for retryable function")
		error_handling.PanicOnAnyError = false
	}

	return retry.Do(function,
		retry.Delay(time.Second),
		retry.DelayType(retry.BackOffDelay),
		retry.OnRetry(func(n uint, err error) {
			logrus.Infof("Retry %d - happend because of %s", n+1, err)
		}),
		retry.Attempts(maxAttempts),
	)
}

// ExecuteWithRetry reruns a function in case of error and logs error
func ExecuteWithRetry[T string | *exec.Cmd](function func(T) (string, error), command T, maxAttempts uint) error {
	err := ExecuteFunctionWithRetry(func() error {
		var retryError error
		_, retryError = function(command)
		return retryError
	}, maxAttempts)

	return err
}
