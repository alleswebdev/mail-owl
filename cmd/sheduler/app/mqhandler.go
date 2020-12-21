package App

import (
	"github.com/alleswebdev/mail-owl/internal/models"
	"github.com/alleswebdev/mail-owl/internal/services/broker"
	"github.com/alleswebdev/mail-owl/internal/services/broker/rabbitmq"
)

// handle msgs, if a msg cant handle return true and error
func (app *App) SchedulerHandler(msg broker.Message) (bool, error) {
	notice := models.SchedulerNotice{Raw: msg.Body, State: models.New}
	err := notice.UnmarshalJSON(msg.Body)

	app.Logger.Infof("received notice with state %s, id %d", notice.State, notice.Id)

	if err != nil {
		return false, err
	}

	err = app.HandleByState(notice)

	if err != nil {
		return false, err
	}

	return false, nil
}

func (app *App) HandleByState(notice models.SchedulerNotice) error {
	err := notice.Save(app.Storage)
	if err != nil {
		return err
	}

	switch notice.State {
	case models.New:
		app.PublishState(PublishState{
			Notice: notice,
			State:  models.Build,
			Queue:  rabbitmq.BuilderQueue,
		})

		app.Logger.Infof(" notice with id %d send to building", notice.Id)
	case models.Builded:
		queue := rabbitmq.EmailQueue
		if notice.Type == "sms" {
			queue = rabbitmq.SmsQueue
		}

		app.PublishState(PublishState{
			Notice: notice,
			State:  models.Builded,
			Queue:  queue,
		})

		app.Logger.Infof(" notice with id %d send to sending", notice.Id)
	}

	return nil
}
