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

// HTTPHookMsg represents a message to be sent to the HTTP log endpoint
type HTTPHookMsg struct {
	Format  string
	Type    string
	Content string
}

// HTTPHook is the logrus hook which sends logs with level warning or above to the specified HTTP endpoint
type HTTPHook struct {
	id            string
	endpoint      string
	cleanupSignal chan bool
	msgs          chan HTTPHookMsg
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
			requestContent := map[string]string{
				"from":    hook.id,
				"format":  msg.Format,
				"type":    msg.Type,
				"content": msg.Content,
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
func (hook *HTTPHook) Push(msg HTTPHookMsg) {
	hook.msgs <- msg
}

// Cleanup tells the HTTPHook to send all remaining messages, then shutdown
func (hook *HTTPHook) Cleanup() {
	hook.cleanupSignal <- true
	<-hook.cleanupSignal
}

// NewHTTPHook initializes a HTTPHook
func NewHTTPHook(id, endpoint string, ethEndpointHosts []string) *HTTPHook {
	hook := &HTTPHook{
		id:            id,
		endpoint:      endpoint,
		cleanupSignal: make(chan bool),
		msgs:          make(chan HTTPHookMsg, 64),
	}
	go hook.run()
	initMsg := HTTPHookMsg{
		Type: "init",
	}
	initContent := map[string]interface{}{
		"msg":              "HTTP endpoint initialized",
		"timestamp":        time.Now(),
		"ethEndpointHosts": ethEndpointHosts,
	}
	initContentJSON, err := json.Marshal(initContent)
	if err == nil {
		initMsg.Format = "json"
		initMsg.Content = string(initContentJSON)
	} else {
		initMsg.Format = "fallback"
		initMsg.Content = fmt.Sprintf("HTTP endpoint initialized (timestamp: %v, ethEndpointHosts: %v)", time.Now(), ethEndpointHosts)
	}
	hook.Push(initMsg)
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
	msg := HTTPHookMsg{Type: "log"}
	var formatter logrus.Formatter
	formatter = &logrus.JSONFormatter{}
	bs, err := formatter.Format(entry)
	if err == nil {
		msg.Format = "json"
		msg.Content = string(bs)
	} else {
		formatter = &logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		}
		bs, err = formatter.Format(entry)
		if err == nil {
			msg.Format = "text"
			msg.Content = string(bs)
		} else {
			msg.Format = "fallback"
			msg.Content = fmt.Sprintf(
				"All formatters failed (timestamp: {{%v}}, msg: {{%s}}, data: {{%v}})",
				entry.Time, entry.Message, entry.Data,
			)
		}
	}
	hook.Push(msg)
	return nil
}

func init() {
	L.Level = logrus.InfoLevel
}
