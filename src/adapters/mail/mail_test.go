package mail_test

import (
	"example/web-service-gin/src/adapters/mail"
	"example/web-service-gin/src/core/entity"
	"example/web-service-gin/src/infra/configuration"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setUp() *mail.MailService {
	configs := &configuration.Mail{
		Key:        "mock-api-key",
		TemplateId: "mock-template-id",
	}

	return mail.NewMailService(configs)
}

func TestNotifyRequestStatus_Success(t *testing.T) {
	mockSendGrid := setUp()

	request := &entity.Request{
		ID:        22,
		UserEmail: "usuario@email.com",
	}

	gock.New("https://api.sendgrid.com").
		Post("/v3/mail/send").
		Reply(202).
		JSON(map[string]interface{}{})

	err := mockSendGrid.NotifyRequestStatus(request, "sucesso")

	assert.NoError(t, err)
}

func TestNotifyRequestStatus_Error(t *testing.T) {
	mockSendGrid := setUp()

	request := &entity.Request{
		ID:        22,
		UserEmail: "xxxxxxx.com",
	}

	gock.New("https://api.sendgrid.com").
		Post("/v3/mail/send").
		Reply(500).
		JSON(map[string]interface{}{})

	err := mockSendGrid.NotifyRequestStatus(request, "error")

	assert.Error(t, err)
}
