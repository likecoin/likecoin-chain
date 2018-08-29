package stdouthook

import (
	"fmt"
	"log"
	"os"

	"github.com/gogap/config"
	"github.com/gogap/logrus_mate"
	"github.com/sirupsen/logrus"
)

// StdoutHook is a hook for sending log to stdout
type StdoutHook struct {
	formatter logrus.Formatter
	level     uint32
}

// NewStdoutHook creates a new StdoutHook
func NewStdoutHook(config config.Configuration) (hook logrus.Hook, err error) {
	hook = &StdoutHook{
		formatter: &logrus.JSONFormatter{},
		level:     uint32(config.GetInt32("level", 0)),
	}
	return hook, err
}

// Fire sends log to stdout
func (h *StdoutHook) Fire(entry *logrus.Entry) (err error) {
	if uint32(entry.Level) > h.level {
		return nil
	}

	msg, err := h.formatter.Format(entry)
	if err != nil {
		log.Printf("Unable to format JSON from logrus entry, %v", err)
		return err
	}
	msgStr := string(msg)

	switch entry.Level {
	case logrus.PanicLevel:
		fallthrough
	case logrus.FatalLevel:
		fallthrough
	case logrus.ErrorLevel:
		_, err = fmt.Fprintln(os.Stderr, msgStr)
	case logrus.WarnLevel:
		fallthrough
	case logrus.InfoLevel:
		fallthrough
	case logrus.DebugLevel:
		_, err = fmt.Println(msgStr)
	default:
		return nil
	}

	return err
}

// Levels set logging level for this hook
func (h *StdoutHook) Levels() []logrus.Level {
	return logrus.AllLevels[:h.level+1]
}

func init() {
	logrus_mate.RegisterHook("stdout", NewStdoutHook)
}
