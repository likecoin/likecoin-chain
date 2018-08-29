package log

import (
	"log"
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/gogap/logrus_mate"
	appConf "github.com/likecoin/likechain/abci/config"
	"github.com/sirupsen/logrus"

	// For using file log
	_ "github.com/gogap/logrus_mate/writers/rotatelogs"
	// For using net log hook
	_ "github.com/likecoin/likechain/abci/log/hooks/net"
	// For using stdout log hook
	_ "github.com/likecoin/likechain/abci/log/hooks/stdout"
)

var (
	// L is an instance of a logger
	L = logrus.New()

	config = appConf.GetConfig()
)

func init() {
	err := logrus_mate.Hijack(L, logrus_mate.ConfigString(config.LogConfig))
	if err != nil {
		log.Printf("Error occurs when applying config for logging, %v, using default config for logging\n", err)

		formatter := &runtime.Formatter{
			ChildFormatter: &logrus.JSONFormatter{},
			// Enable line number logging
			Line: true,
			// Enable file name logging
			File: true,
		}

		// Replace the default Logrus Formatter with Banzai Cloud runtime Formatter
		L.Formatter = formatter

		// Output to stdout instead of the default stderr
		L.Out = os.Stdout

		if config.IsProduction() {
			L.Level = logrus.ErrorLevel
		} else {
			L.Level = logrus.DebugLevel
		}
	}
}
