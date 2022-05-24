package senders

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// SlackSender sender
type SlackSender struct {
	Sender
}

// SlackMessage slack message model
type SlackMessage struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	AsUser  bool   `json:"as_user"`
	Text    string `json:"text"`
	Token   string `json:"token"`
}

// SlackPusher pusher
var SlackPusher SlackSender

// init
func init() {
	SlackPusher = SlackSender{}
}

// is support
func (s *SlackSender) IsSupport() bool {
	return true
}

// Send send
func (s *SlackSender) Send(notifications []*Notification) {
	//if !utils.Config.SenderConfig.Slack.IsEnabled {
	//	return
	//}
	for _, item := range notifications {
		s.SingleSend(item)
	}
}

// SingleSend send notification
func (s *SlackSender) SingleSend(notification *Notification) {
	message := SlackMessage{
		AsUser:  true,
		Channel: "utils.Config.SenderConfig.Slack.Channel",
		Text:    s.BuildMessage(notification),
	}
	data := url.Values{}
	data.Set("token", "utils.Config.SenderConfig.Slack.RobotToken")
	data.Add("channel", message.Channel)
	data.Add("text", message.Text)
	data.Add("as_user", strconv.FormatBool(message.AsUser))

	body, err := http.Post("https://slack.com/api/chat.postMessage", "application/x-www-form-urlencoded",
		bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Error(err)
		return
	}
	content, err := ioutil.ReadAll(body.Body)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(string(content))
}

// BuildMessage build message
func (s *SlackSender) BuildMessage(notification *Notification) string {
	return fmt.Sprintf("status:%d,type:%s,url:%s", notification.HTTPStatus, notification.Type, notification.URL)
}
