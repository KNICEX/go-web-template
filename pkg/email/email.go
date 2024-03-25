package email

import "errors"

type Driver interface {
	Close()
	Send(to, title, body string) error
}

var (
	ErrChanNotOpen    = errors.New("email queue is not started")
	ErrNoActiveDriver = errors.New("no available email provider")
)

func Send(to, title, body string) error {
	return Client.Send(to, title, body)
}
