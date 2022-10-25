package internal

import (
	"github.com/conplementag/cops-hq/v2/pkg/error_handling"
	"github.com/sirupsen/logrus"
)

func ReturnErrorOrPanic(err error) error {
	if err != nil && error_handling.PanicOnAnyError {
		// we log the error, so it ends up in the log file as well. Consequence: it will be shown twice in the stdout, but
		// this we have to live with
		logrus.Error(err)
		panic(err)
	}

	return err
}
