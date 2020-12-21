package App

import (
	"github.com/alleswebdev/mail-owl/internal/models"
	"github.com/alleswebdev/mail-owl/internal/services/broker"
)

// handle msgs, if a msg cant handle return true and error
func (app *App) EmailHandler(msg broker.Message) (bool, error) {
	notice := models.SchedulerNotice{}
	err := notice.UnmarshalJSON(msg.Body)

	if err != nil {
		return false, err
	}

	app.Logger.Infof("received notice with state %s, id %d", notice.State, notice.Id)

	if notice.Debug {
		app.Logger.Infof("id:%d, %s", notice.Id, "notice in debug mode")
		app.PublishState(notice, models.Success, nil)
		return false, nil
	}

	err = app.Mailer.Send(notice)

	if err != nil {
		app.PublishState(notice, models.Error, err)
		return true, err
	}

	app.PublishState(notice, models.Success, nil)

	app.Logger.Infof("notice with id %d send", notice.Id)

	return false, nil
}
