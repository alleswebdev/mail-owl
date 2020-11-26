package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFill(t *testing.T) {
	nt := &Notification{Type: "email", Data: `{"to":["support@torbor.ru"],"bcc":[],"cc":[],"attachments":[],"subject":"sonya_test@torbor.ru:test subject","template":"contact","params":{"text":"test content"},"debug":true}`}
	err := nt.Fill()
	assert.Nil(t, err)
	assert.Equal(t, nt.Template, "contact")
}
