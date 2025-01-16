package queue

import (
	"context"
	"example/web-service-gin/src/core/entity"
	"log"
	"log/slog"
	"time"
)

// MessageProcessor it's a function tha process each message
type MessageProcessor func(ctx context.Context, msg entity.EventMessage)

// StartQueueConsumer inicia o consumo de uma fila
func StartQueueConsumer(handler *SQSHandler, queueURL string, processor MessageProcessor, ctx context.Context) {
	slog.Info("Starting queue consumer:", "queueUrl", queueURL)

	for {
		messages, err := handler.ReceiveMessages(queueURL, 5, 5)
		if err != nil {
			log.Printf("Error receiving messages from %s: %v", queueURL, err)
			time.Sleep(5 * time.Second) // Retry with backoff
			continue
		}

		for _, message := range messages {

			toDomain := entity.EventMessage{
				MessageID: *message.MessageId,
				Source:    queueURL,
				Body:      *message.Body,
				Date:      message.Attributes["SentTimestamp"],
			}

			processor(ctx, toDomain)
			// Delete the message from queue after finish the process
			err := handler.DeleteMessage(queueURL, *message.ReceiptHandle)
			if err != nil {
				log.Printf("Error deleting message from %s: %v", queueURL, err)
			}
		}
	}
}
