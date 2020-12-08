package models

import (
	"context"
	"github.com/alleswebdev/mail-owl/internal/storage"
)

type SchedulerState string

const (
	New     SchedulerState = "new"
	Builded                = "builded"
	Error                  = "error"
	Success                = "success"
)

type SchedulerNotice struct {
	Id          int          `json:"id"`
	Type        string       `json:"type"`
	To          []string     `json:"to"`
	Bcc         []string     `json:"bcc"`
	Cc          []string     `json:"cc"`
	Subject     string       `json:"subject"`
	Template    string       `json:"template"`
	Attachments []Attachment `json:"attachment"`
	Params      map[string]interface{}
	Debug       bool
	Raw         []byte         `json:"-"`
	State       SchedulerState `json:"state"`
}

type Attachment struct {
	Filename string
	Body     []byte
}

func (n *SchedulerNotice) Save(db storage.DBStorage) error {
	var (
		id int
	)

	err := db.Db.QueryRow(context.Background(), `INSERT INTO public.notifications(status, type, raw, created_at, debug)
			VALUES('new', $1, $2, now(), $3) RETURNING id`, n.Type, string(n.Raw), n.Debug).Scan(&id)

	n.Id = id

	return err
}
