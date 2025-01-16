package queue

import (
	"context"
	"example/web-service-gin/src/infra/configuration"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQSHandler structure
type SQSHandler struct {
	Configs *configuration.Aws
	Client  *sqs.Client
}

// NewSQSHandler creates a new instance of SQSHandler
func NewSQSHandler(configs *configuration.Aws) *SQSHandler {
	sqsClient := sqs.NewFromConfig(configs.Config)
	return &SQSHandler{Configs: configs, Client: sqsClient}
}

// ReceiveMessages reads the messages from queue
func (h *SQSHandler) ReceiveMessages(queueURL string, maxMessages int32, waitTime int32) ([]types.Message, error) {
	output, err := h.Client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitTime,
	})
	if err != nil {
		return nil, err
	}
	return output.Messages, nil
}

// SendMessage sends a message to a specif queue
func (h *SQSHandler) SendMessage(queueURL, messageBody string) error {
	_, err := h.Client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
	})
	return err
}

// DeleteMessage deletes a message from the queue
func (h *SQSHandler) DeleteMessage(queueURL, receiptHandle string) error {
	_, err := h.Client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}
