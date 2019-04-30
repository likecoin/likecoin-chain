package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// L is an instance of a logger
var L = logrus.New()

type httpHookMsg struct {
	format    string
	logString string
}

// HTTPHook is the logrus hook which sends logs with level warning or above to the specified HTTP endpoint
type HTTPHook struct {
	id            string
	endpoint      string
	cleanupSignal chan bool
	msgs          chan httpHookMsg
}

func (hook *HTTPHook) run() {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	isDone := false
	for {
		select {
		case <-hook.cleanupSignal:
			isDone = true
		case msg := <-hook.msgs:
			requestContent := map[string]interface{}{
				"from":   hook.id,
				"format": msg.format,
				"log":    msg.logString,
			}
			requestBody, err := json.Marshal(requestContent)
			if err != nil {
				L.
					WithError(err).
					Info("Failed to marshal HTTP log hook request content")
				break
			}
			bodyReader := bytes.NewReader(requestBody)
			res, err := client.Post(hook.endpoint, "application/json", bodyReader)
			if err != nil {
				L.
					WithField("endpoint", hook.endpoint).
					WithError(err).
					Info("Failed to POST to HTTP log hook endpoint")
				break
			}
			defer func() {
				if !res.Close {
					res.Body.Close()
				}
			}()
			if res.StatusCode != 200 {
				body, _ := ioutil.ReadAll(res.Body)
				L.
					WithField("code", res.StatusCode).
					WithField("status", res.Status).
					WithField("body", string(body)).
					Info("HTTP log hook endpoint returns failure")
				break
			}
		}
		if isDone && len(hook.msgs) == 0 {
			break
		}
	}
	hook.cleanupSignal <- true
}

// Push pushes a log message into queue, which will be sent to the HTTP endpoint asynchronously
func (hook *HTTPHook) Push(entry *logrus.Entry) {
	msg := httpHookMsg{}
	var formatter logrus.Formatter
	formatter = &logrus.JSONFormatter{}
	bs, err := formatter.Format(entry)
	if err == nil {
		msg.format = "json"
		msg.logString = string(bs)
	} else {
		formatter = &logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		}
		bs, err = formatter.Format(entry)
		if err == nil {
			msg.format = "text"
			msg.logString = string(bs)
		} else {
			msg.format = "fallback"
			msg.logString = fmt.Sprintf("All formatters failed (msg: {{%s}}, timestamp: {{%v}})", entry.Message, entry.Time)
		}
	}
	hook.msgs <- msg
}

// Cleanup tells the HTTPHook to send all remaining messages, then shutdown
func (hook *HTTPHook) Cleanup() {
	hook.cleanupSignal <- true
	<-hook.cleanupSignal
}

// NewHTTPHook initializes a HTTPHook
func NewHTTPHook(id, endpoint string) *HTTPHook {
	hook := &HTTPHook{
		id:            id,
		endpoint:      endpoint,
		cleanupSignal: make(chan bool),
		msgs:          make(chan httpHookMsg, 64),
	}
	hook.msgs <- httpHookMsg{
		format:    "init",
		logString: fmt.Sprintf("HTTP Endpoint Initialized at %v", time.Now()),
	}
	go hook.run()
	return hook
}

// Levels returns the levels to be hooked, implements logrus.Hook
func (hook *HTTPHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}

// Fire is triggered when there are logrus logs with matching level, implements logrus.Hook
func (hook *HTTPHook) Fire(entry *logrus.Entry) error {
	hook.Push(entry)
	return nil
}

func init() {
	L.Level = logrus.InfoLevel
}
