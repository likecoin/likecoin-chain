package log

import (
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/likecoin/likechain/abci/env"
	"github.com/sirupsen/logrus"
)

// L is an instance of a logger
var L = logrus.New()

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

	if env.IsProduction {
		L.Level = logrus.InfoLevel
	} else {
		L.Level = logrus.DebugLevel
	}
}
