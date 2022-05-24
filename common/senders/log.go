package senders

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// LogSender sender
type LogSender struct {
	Sender
}

// LogPusher log pusher
var LogPusher LogSender

func init() {
	LogPusher = LogSender{}
}

// IsSupport is support
func (f *LogSender) IsSupport() bool {
	return true
}

// Send send
func (f *LogSender) Send(notifications []*Notification) {
	//if !utils.Config.SenderConfig.Log.IsEnabled {
	//	return
	//}
	for _, item := range notifications {
		f.SingleSend(item)
	}
}

// SingleSend send notification
func (f *LogSender) SingleSend(notification *Notification) {
	message := f.BuildMessage(notification)
	log.Info(message)
}

// BuildMessage build message
func (f *LogSender) BuildMessage(notification *Notification) string {
	return fmt.Sprintf("status:%d,type:%s,monitor:%s,url:%s", notification.HTTPStatus,
		notification.Reason, notification.URL)
}