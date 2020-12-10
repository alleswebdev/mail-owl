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
		return true, err
	}

	return false, nil
}

func (app *App) HandleByState(notice models.SchedulerNotice) error {
	switch notice.State {
	case models.New:
		err := notice.Save(app.Storage)

		if err != nil {
			return err
		}

		// вынести эту часть в отдельеую функцию
		notice.State = models.Build
		noticeEncode, err := notice.MarshalJSON()

		if err != nil {
			return err
		}

		err = app.Broker.Publish(broker.Message{
			Body:    noticeEncode,
			Headers: nil,
		}, rabbitmq.BuilderQueue)

		if err != nil {
			return err
		}

		app.Logger.Infof(" notice with id %d send to building", notice.Id)
	}

	return nil
}
