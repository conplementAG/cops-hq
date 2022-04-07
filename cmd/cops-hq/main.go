package main

import (
	"github.com/denisbiondic/cops-hq/pkg/logging"
	"github.com/sirupsen/logrus"
)

// main in the cops-hq project is a simple "consumer" CLI, which can be used for manual code tests.
// Logic and code here should be kept to minimum. Use automated tests instead.
func main() {
	logging.Init("hq.log")
	logrus.Info("this is a test")
}
