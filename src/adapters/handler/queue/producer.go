package queue

import (
	"encoding/json"
	"example/web-service-gin/src/core/entity"
	"log/slog"
)

type SQSProducer struct {
	SQSHandler *SQSHandler
	QueueUrl   string
}

func NewSQSProducer(handler *SQSHandler, url string) *SQSProducer {
	return &SQSProducer{SQSHandler: handler, QueueUrl: url}
}

// SQSHandler structure
func (h *SQSProducer) SendVideoProccessToQueue(request *entity.Request) error {

	// Map domain to reciver contract
	bodyData := SnapVideoRequest{
		Id:           request.ID,
		IdUser:       request.UserId,
		FileSize:     request.VideoSize,
		S3FileKey:    request.VideoKey,
		CreationDate: request.CreatedAt,
	}

	jsonData, parseError := json.Marshal(bodyData)

	if parseError != nil {
		slog.Error("Error trying to conver entity to JSON", "error", parseError)
	}

	err := h.SQSHandler.SendMessage(h.QueueUrl, string(jsonData))

	if err != nil {
		slog.Error("Error trying to send message", "destination", h.QueueUrl)
	}

	return nil

}
