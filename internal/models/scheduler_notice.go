package models

import (
	"context"
	"github.com/alleswebdev/mail-owl/internal/storage"
)

type SchedulerState string

const (
	New     SchedulerState = "new"
	Build                  = "build"
	Builded                = "builded"
	Error                  = "error"
	Success                = "success"
)

type SchedulerNotice struct {
	Id          int                    `json:"id"`
	Type        string                 `json:"type"`
	To          []string               `json:"to"`
	Bcc         []string               `json:"bcc"`
	Cc          []string               `json:"cc"`
	Subject     string                 `json:"subject"`
	Template    string                 `json:"template"`
	Attachments []Attachment           `json:"attachment"`
	Params      map[string]interface{} `json:"params"`
	Debug       bool                   `json:"debug"`
	Raw         []byte                 `json:"-"`
	Build       []byte                 `json:"build"`
	State       SchedulerState         `json:"state"`
	Error       string                 `json:"error"`
}

type Attachment struct {
	Filename string
	Body     []byte
}

func (n *SchedulerNotice) Save(db storage.DBStorage) error {
	var (
		id  int
		err error = nil
	)

	if n.Id > 0 {
		err = db.Db.QueryRow(context.Background(),
			`UPDATE  public.notifications 
				SET status = $1,
					status_message = $2,
					updated_at = now()
				WHERE id = $3
				RETURNING id
			`, n.State, n.Error, n.Id).Scan(&id)
	} else {
		err = db.Db.QueryRow(context.Background(), `INSERT INTO public.notifications(status, type, raw, created_at, debug)
			VALUES('new', $1, $2, now(), $3) RETURNING id`, n.Type, string(n.Raw), n.Debug).Scan(&id)
		n.Id = id
	}

	return err
}
