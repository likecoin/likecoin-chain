package log

import (
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	appConf "github.com/likecoin/likechain/abci/config"
	"github.com/sirupsen/logrus"
)

var (
	// L is an instance of a logger
	L = logrus.New()

	config = appConf.GetConfig()
)

func init() {
	formatter := &runtime.Formatter{
		ChildFormatter: &logrus.TextFormatter{},
		// Enable line number logging
		Line: true,
		// Enable file name logging
		File: true,
	}

	// Replace the default Logrus Formatter with Banzai Cloud runtime Formatter
	L.Formatter = formatter

	// Output to stdout instead of the default stderr
	L.Out = os.Stdout

	logLevel, err := logrus.ParseLevel(config.Log.Level)
	if err == nil {
		L.Level = logLevel
	} else if config.IsProduction() {
		L.Level = logrus.ErrorLevel
	} else {
		L.Level = logrus.DebugLevel
	}
}
