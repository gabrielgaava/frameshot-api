package notification_test

import (
	"context"
	"errors"
	"example/web-service-gin/src/adapters/handler/notification"
	"example/web-service-gin/src/infra/configuration"
	"github.com/aws/aws-sdk-go-v2/aws"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSNSClient implementa um mock para o cliente SNS.
type MockSNSClient struct {
	mock.Mock
}

func (m *MockSNSClient) Publish(ctx context.Context, input *sns.PublishInput, opts ...func(*sns.Options)) (*sns.PublishOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}

func TestNewSNSHandler(t *testing.T) {
	ctx := context.Background()
	config := configuration.Aws{
		Config: aws.Config{},
	}
	service := notification.NewSNSHandler(&config, ctx)
	assert.NotNil(t, service)
}

func TestPublishMessage_Success(t *testing.T) {
	mockSNSClient := new(MockSNSClient)
	mockCtx := context.Background()

	// Configuração do mock para retornar sucesso
	mockSNSClient.On("Publish", mockCtx, mock.AnythingOfType("*sns.PublishInput")).
		Return(&sns.PublishOutput{}, nil)

	conf := &configuration.Aws{}
	handler := &notification.SNSHandler{
		Client:  mockSNSClient,
		Configs: conf,
		Ctx:     mockCtx,
	}

	topicArn := "arn:aws:sns:us-east-1:123456789012:example-topic"
	message := "Test message"
	attributes := map[string]string{"Attribute1": "Value1", "Attribute2": "Value2"}

	err := handler.PublishMessage(topicArn, message, attributes)

	assert.NoError(t, err)
	mockSNSClient.AssertExpectations(t)
}

func TestPublishMessage_MissingTopicArn(t *testing.T) {
	mockSNSClient := new(MockSNSClient)
	mockCtx := context.Background()

	conf := &configuration.Aws{}
	handler := &notification.SNSHandler{
		Client:  mockSNSClient,
		Configs: conf,
		Ctx:     mockCtx,
	}

	message := "Test message"
	attributes := map[string]string{"Attribute1": "Value1"}

	err := handler.PublishMessage("", message, attributes)

	assert.EqualError(t, err, "topicArn is required")
	mockSNSClient.AssertNotCalled(t, "Publish")
}

func TestPublishMessage_MissingMessage(t *testing.T) {
	mockSNSClient := new(MockSNSClient)
	mockCtx := context.Background()

	conf := &configuration.Aws{}
	handler := &notification.SNSHandler{
		Client:  mockSNSClient,
		Configs: conf,
		Ctx:     mockCtx,
	}

	topicArn := "arn:aws:sns:us-east-1:123456789012:example-topic"
	attributes := map[string]string{"Attribute1": "Value1"}

	err := handler.PublishMessage(topicArn, "", attributes)

	assert.EqualError(t, err, "message is required")
	mockSNSClient.AssertNotCalled(t, "Publish")
}

func TestPublishMessage_Failure(t *testing.T) {
	mockSNSClient := new(MockSNSClient)
	mockCtx := context.Background()

	mockSNSClient.On("Publish", mockCtx, mock.AnythingOfType("*sns.PublishInput")).
		Return((*sns.PublishOutput)(nil), errors.New("SNS publish error"))

	conf := &configuration.Aws{}
	handler := &notification.SNSHandler{
		Client:  mockSNSClient,
		Configs: conf,
		Ctx:     mockCtx,
	}

	topicArn := "arn:aws:sns:us-east-1:123456789012:example-topic"
	message := "Test message"
	attributes := map[string]string{"Attribute1": "Value1"}

	err := handler.PublishMessage(topicArn, message, attributes)

	assert.EqualError(t, err, "failed to publish message: SNS publish error")
	mockSNSClient.AssertExpectations(t)
}
