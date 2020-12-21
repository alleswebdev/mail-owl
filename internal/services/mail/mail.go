package mail

import (
	"fmt"
	"github.com/alleswebdev/mail-owl/internal/config"
	"github.com/alleswebdev/mail-owl/internal/models"
	"gopkg.in/gomail.v2"
	"io"
)

type Mailer struct {
	Dialer     *gomail.Dialer
	Open       bool
	sendCloser gomail.SendCloser
}

func NewMailer(config config.Config) *Mailer {
	dialer := gomail.NewDialer(config.SmtpHost, config.SmtpPort, config.EmailFrom, config.SmtpPassword)
	return &Mailer{Dialer: dialer, Open: false}
}

func (m *Mailer) Shutdown() error {
	if m.Open {
		if err := m.sendCloser.Close(); err != nil {
			return err
		}

		m.Open = false
	}

	return nil
}

func (m *Mailer) Send(notice models.SchedulerNotice) error {
	var err error

	if notice.Debug {
		return fmt.Errorf("notice in debug state")
	}

	if !m.Open {
		if m.sendCloser, err = m.Dialer.Dial(); err != nil {
			return fmt.Errorf("smtp connection error:%s", err)
		}

		m.Open = true
	}

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", m.Dialer.Username, "Mail-owl")

	msg.SetHeader("To", notice.To...)
	msg.SetHeader("Subject", notice.Subject)

	if len(notice.Bcc) != 0 {
		msg.SetHeader("Bcc", notice.Bcc...)
	}

	if len(notice.Cc) != 0 {
		msg.SetHeader("Cc", notice.Cc...)
	}

	msg.AddAlternativeWriter("text/html", func(w io.Writer) error {
		_, err = w.Write(notice.Build)
		return err
	})

	if err := gomail.Send(m.sendCloser, msg); err != nil {
		return fmt.Errorf("smtp sending error:%s", err)
	}

	return nil
}
