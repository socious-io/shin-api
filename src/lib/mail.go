package lib

import (
	"errors"
	"strconv"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	SendGridTemplates map[string]string
	SendGridClient    SendGridType
)

type SendGridType struct {
	ApiKey   string
	Url      string
	Disabled bool
}

func (sgc *SendGridType) SendWithTemplate(address string, name string, templateId string, items map[string]string) error {
	if sgc.Disabled {
		return nil
	}
	//Create Mail payload
	m := mail.NewV3Mail()
	m.SetFrom(mail.NewEmail("Shin", "no-replay@shinid.com"))
	m.SetTemplateID(templateId)

	//Adding Personalization
	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(name, address),
	}
	p.AddTos(tos...)
	for key, value := range items {
		p.SetDynamicTemplateData(key, value)
	}
	m.AddPersonalizations(p)

	//Setup the request
	request := sendgrid.GetRequest(sgc.ApiKey, "/v3/mail/send", sgc.Url)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	response, err := sendgrid.API(request)
	if err != nil {
		return err
	} else if strings.Split(strconv.Itoa(response.StatusCode), "")[0] != "2" {
		return errors.New(response.Body)
	}
	return nil
}

func InitSendGridLib(sgc SendGridType) {
	SendGridClient = sgc
	SendGridTemplates = map[string]string{
		"otp":                          "d-5ace79d4674b45d1bfb9b35c4d6eb8c0",
		"forget-password":              "d-d68cc8d8409942599f761261e5a7fbcb",
		"credentials-recipients":       "d-12fddb16345e4073a6741884237ed39c",
		"credentials-import-completed": "d-12fddb16345e4073a6741884237ed39c", //TODO: Add sendgrid template
	}
}
