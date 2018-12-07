package log

import "github.com/sirupsen/logrus"

// L is an instance of a logger
var L = logrus.New()

func init() {
	L.Level = logrus.InfoLevel
}
