package sms

import (
	"encoding/json"
	"fmt"
	"github.com/alleswebdev/mail-owl/internal/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type SmsRu struct {
	SmsApiId string
}

func NewSms(smsApiId string) *SmsRu {
	return &SmsRu{SmsApiId: smsApiId}
}

type Response struct {
	Status     string `json:"status"`
	StatusCode string `json:"status_code"`
}

func (s *SmsRu) Send(notice models.SchedulerNotice) (error, string) {
	uri := fmt.Sprintf("https://sms.ru/sms/send?api_id=%s&to=%s&msg=%s&json=1",
		s.SmsApiId,
		strings.Join(notice.To, ","),
		url.QueryEscape(string(notice.Build)),
	)

	resp, err := http.Get(uri)

	if err != nil {
		return err, ""
	}

	defer resp.Body.Close()

	var response Response

	respBody, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)

	if resp.StatusCode != 200 {
		return fmt.Errorf("sms send status code:%d, status: %s", resp.StatusCode, resp.Status), ""
	}

	if response.Status != "OK" {
		return fmt.Errorf("%s", string(respBody)), ""
	}

	return nil, string(respBody)
}
