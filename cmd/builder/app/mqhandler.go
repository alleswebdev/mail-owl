package App

import (
	"bytes"
	"github.com/alleswebdev/mail-owl/internal/models"
	"github.com/alleswebdev/mail-owl/internal/services/broker"
	"github.com/alleswebdev/mail-owl/internal/services/broker/rabbitmq"
	"html/template"
)

// handle msgs, if a msg cant handle return true and error
func (app *App) BuilderHandler(msg broker.Message) (bool, error) {
	notice := models.SchedulerNotice{}
	err := notice.UnmarshalJSON(msg.Body)

	app.Logger.Infof("received notice with state %s, id %d", notice.State, notice.Id)
	if err != nil {
		return false, err
	}

	var tpl bytes.Buffer

	tmpl, err := template.ParseFiles("internal/templates/" + notice.Template + ".html")

	if err != nil {
		return false, err
	}

	err = tmpl.Execute(&tpl, notice.Params)

	if err != nil {
		return false, err
	}

	notice.Build = tpl.Bytes()

	notice.State = models.Builded
	noticeEncode, err := notice.MarshalJSON()

	if err != nil {
		return false, err
	}

	err = app.Broker.Publish(broker.Message{
		Body:    noticeEncode,
		Headers: nil,
	}, rabbitmq.SchedulerQueue)

	if err != nil {
		return false, err
	}

	app.Logger.Infof("notice with id %d builded and send to sheduler", notice.Id)

	return false, nil
}
