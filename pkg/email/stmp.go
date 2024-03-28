package email

import (
	"go-web-template/pkg/conf"
	"go-web-template/pkg/logger"
	"gopkg.in/gomail.v2"
	"time"
)

type SMTP struct {
	ch     chan *gomail.Message
	chOpen bool
}

func NewSMTPClient() *SMTP {
	client := &SMTP{
		ch:     make(chan *gomail.Message, 30),
		chOpen: false,
	}
	client.Init()
	return client
}

func (client *SMTP) Send(to, title, body string) error {
	if !client.chOpen {
		return ErrChanNotOpen
	}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", conf.EmailConf.User, conf.EmailConf.Name)
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)
	client.ch <- m
	return nil
}

func (client *SMTP) Close() {
	if client.ch != nil {
		close(client.ch)
	}
}

func (client *SMTP) Init() {
	go func() {

		defer func() {
			if err := recover(); err != nil {
				client.chOpen = false
				logger.L().Error("Email sending queue crashed: ", err, "Resetting in 10 seconds.")
				time.Sleep(10 * time.Second)
				client.Init()
			}
		}()

		d := gomail.NewDialer(conf.EmailConf.Host,
			conf.EmailConf.Port,
			conf.EmailConf.User,
			conf.EmailConf.Password,
		)
		client.chOpen = true
		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-client.ch:
				if !ok {
					logger.L().Debug("Email queue closing...")
					client.chOpen = false
					return
				}
				if !open {
					// 尝试连接
					if s, err = d.Dial(); err != nil {
						panic(err)
					}
					open = true
				}
				if err = gomail.Send(s, m); err != nil {
					logger.L().Warn("Failed to send email: ", err)
				} else {
					logger.L().Debug("Email sent.")
				}
			//	长时间没有发送邮件，关闭连接
			case <-time.After(time.Duration(conf.EmailConf.Keepalive) * time.Second):
				if open {
					if err = s.Close(); err != nil {
						logger.L().Warn("Failed to close SMTP connection: ", err)
					}
					open = false
				}
			}
		}
	}()
}
