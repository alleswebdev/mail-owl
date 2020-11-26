package NoticeBuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"vcs.torbor.ru/notizer/workers/spectacler/internal/config"
	"vcs.torbor.ru/notizer/workers/spectacler/internal/models"
	"vcs.torbor.ru/notizer/workers/spectacler/internal/storage"
)

func TestBuilder_ParseNoticeSms(t *testing.T) {
	cfg := *config.LoadConfig("../../..")
	tplStorage, err := storage.NewTemplateStorage(cfg)

	assert.Nil(t, err)

	b := New(tplStorage, &cfg)

	notice := models.Notification{
		Type:        "sms",
		Data:        `{"recipients":["779183581117"],"template":"order-is-created","params":{"id":"test"},"debug":false,"project":"portal","prefix":"test-files\/"}`,
		Project:     "",
		Prefix:      "",
		Layout:      "",
		NoLayout:    false,
		Subject:     "",
		Template:    "",
		Attachments: nil,
		Params:      nil,
	}

	err, r := b.ParseNotice(&notice)

	assert.Nil(t, err)

	assert.Equal(t, r.String(), "Ваша заявка #test принята. Спасибо, что выбрали нас!")
}

const testRsponse = `<h1>test was here!</h1>`

func TestBuilder_ParseNoticeEmail(t *testing.T) {
	cfg := *config.LoadConfig("../../..")
	tplStorage, err := storage.NewTemplateStorage(cfg)

	assert.Nil(t, err)

	b := New(tplStorage, &cfg)

	notice := models.Notification{
		Type:        "email",
		Data:        `{"to":["aalles@torbor.ru"],"bcc":["aalles@torbor.ru"],"cc":["aalles@torbor.ru"],"attachments":[],"subject":"\u0421\u0442\u0440\u0435\u0441\u0441 \u0442\u0435\u0441\u0442!","template":"letter","params":{"unsubcribelink":"https:\/\/torbor.ru\/fsdgdsfg","text":"<h1>test was here!<\/h1>"},"debug":true,"project":"portal","prefix":"test-files\/","layout":"layout","noLayout":true}`,
		Project:     "",
		Prefix:      "",
		Layout:      "",
		NoLayout:    false,
		Subject:     "",
		Template:    "",
		Attachments: nil,
		Params:      nil,
	}

	err, r := b.ParseNotice(&notice)

	assert.Nil(t, err)
	assert.Contains(t, r.String(), testRsponse)
}
