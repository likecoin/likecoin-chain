package nethook

import (
	"io"
	"log"
	"net"

	"github.com/gogap/config"
	"github.com/gogap/logrus_mate"
	"github.com/sirupsen/logrus"
)

// NetHook is a hook for sending log to a specific network address
type NetHook struct {
	writer    io.Writer
	formatter logrus.Formatter
}

// NewNetHook creates a new NetHook
func NewNetHook(config config.Configuration) (hook logrus.Hook, err error) {
	network := config.GetString("network", "")
	address := config.GetString("address", "")
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	hook = &NetHook{
		writer:    conn,
		formatter: &logrus.JSONFormatter{},
	}
	return hook, err
}

// Fire sends log to network address
func (h *NetHook) Fire(entry *logrus.Entry) (err error) {
	msg, err := h.formatter.Format(entry)
	if err != nil {
		log.Printf("Unable to format JSON from logrus entry, %v", err)
		return err
	}
	_, err = h.writer.Write(msg)
	return err
}

// Levels set logging level for this hook
func (*NetHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func init() {
	logrus_mate.RegisterHook("net", NewNetHook)
}
