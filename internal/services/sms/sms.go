package sms

import "github.com/alleswebdev/mail-owl/internal/models"

type Service interface {
	Send(notice models.SchedulerNotice) (error, string)
}
