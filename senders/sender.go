package senders

// Notification notification model
type Notification struct {
	HTTPStatus int
	Reason     string
	Type       string
	URL        string
}

// sender
type Sender interface {
	Send([]*Notification)
	IsSupport() bool
}