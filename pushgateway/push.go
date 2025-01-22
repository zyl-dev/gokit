package pushgateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Sender struct {
	url         string
	token       string
	serviceName string
	action      string // optional
	pusher      string // optional
}

func NewPushGatewaySender(url, token, serviceName, action, pusher string) Sender {
	return Sender{url, token, serviceName, action, pusher}
}

// SendInstantNotice 发送即时通知, 目前支持通道：slack, dingtalk
func (s *Sender) SendInstantNotice(msg, chatGroupName, webhook string) error {
	content, data := make(map[string]string), make(map[string]interface{})
	content["message"] = msg
	content["pusher"] = s.pusher
	content["chat_group_name"] = chatGroupName
	content["webhook_url"] = webhook
	content["service_name"] = s.serviceName
	content["action"] = s.action
	data["data"] = content
	b, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", s.url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s", s.token))
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

// SendEmail 发送邮件, 目前支持服务商：mailgun
func (s *Sender) SendEmail(msg, subject, receiver string) error {
	content, data := make(map[string]string), make(map[string]interface{})
	content["message"] = msg
	content["pusher"] = s.pusher
	content["subject"] = subject
	content["receiver"] = receiver
	content["service_name"] = s.serviceName
	content["action"] = s.action
	data["data"] = content
	b, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", s.url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s", s.token))
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

// SendSMS 发送短信, 目前支持服务商：rrz, nxcloud, yunpian
func (s *Sender) SendSMS(msg, receiver string) error {
	content, data := make(map[string]string), make(map[string]interface{})
	content["message"] = msg
	content["pusher"] = s.pusher
	content["receiver"] = receiver
	content["service_name"] = s.serviceName
	content["action"] = s.action
	data["data"] = content
	b, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", s.url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s", s.token))
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}
