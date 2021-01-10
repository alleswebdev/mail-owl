package App

import (
	"bytes"
	"github.com/alleswebdev/mail-owl/internal/models"
	"github.com/alleswebdev/mail-owl/internal/services/broker"
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
		app.PublishState(notice, models.Error, err)
		return false, err
	}

	err = tmpl.Execute(&tpl, notice.Params)

	if err != nil {
		app.PublishState(notice, models.Error, err)
		return false, err
	}

	notice.Build = tpl.Bytes()

	app.PublishState(notice, models.Builded, nil)

	app.Logger.Infof("notice with id %d builded and send to scheduler", notice.Id)

	return false, nil
}
