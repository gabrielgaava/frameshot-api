package notification

import (
	"context"
	"errors"
	"example/web-service-gin/src/infra/configuration"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type SNSClient interface {
	Publish(ctx context.Context, input *sns.PublishInput, opts ...func(*sns.Options)) (*sns.PublishOutput, error)
}

type SNSHandler struct {
	Client  SNSClient
	Configs *configuration.Aws
	Ctx     context.Context
}

func NewSNSHandler(conf *configuration.Aws, ctx context.Context) *SNSHandler {
	snsClient := sns.NewFromConfig(conf.Config)
	return &SNSHandler{
		Client:  snsClient,
		Configs: conf,
		Ctx:     ctx,
	}
}

// PublishMessage publishes a message to the given SNS topic.
func (h *SNSHandler) PublishMessage(topicArn, message string, attributes map[string]string) error {
	if topicArn == "" {
		return errors.New("topicArn is required")
	}
	if message == "" {
		return errors.New("message is required")
	}

	// Constructing message attributes
	msgAttributes := make(map[string]types.MessageAttributeValue)
	for key, value := range attributes {
		msgAttributes[key] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}

	// Publishing the message
	_, err := h.Client.Publish(h.Ctx, &sns.PublishInput{
		TopicArn:          aws.String(topicArn),
		Message:           aws.String(message),
		MessageAttributes: msgAttributes,
	})

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
