package mail

import (
	"example/web-service-gin/src/core/entity"
	"example/web-service-gin/src/infra/configuration"
	"fmt"
	"strconv"

	"github.com/sendgrid/sendgrid-go"
	mailer "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService struct {
	Config *configuration.Mail
}

func NewMailService(conf *configuration.Mail) *MailService {
	return &MailService{Config: conf}
}

func (service *MailService) NotifyRequestStatus(data *entity.Request, status string) error {
	idString := strconv.FormatUint(data.ID, 10)

	m := mailer.NewV3Mail()
	from := mailer.NewEmail("Frameshot Notification", "no_reply@frameshot.com.br")
	m.SetFrom(from)
	m.SetTemplateID(service.Config.TemplateId)

	p := mailer.NewPersonalization()
	tos := []*mailer.Email{
		mailer.NewEmail("User", data.UserEmail),
	}
	p.AddTos(tos...)
	p.SetDynamicTemplateData("request_id", idString)
	p.SetDynamicTemplateData("status_text", status)
	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(service.Config.Key, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mailer.GetRequestBody(m)
	response, err := sendgrid.API(request)

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println("E-mail enviado com sucesso!")
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		return nil
	}
}
