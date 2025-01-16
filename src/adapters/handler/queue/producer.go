package queue

import (
	"encoding/json"
	"example/web-service-gin/src/core/entity"
	"log/slog"
	"time"
)

type SQSProducer struct {
	Handler  *SQSHandler
	QueueUrl string
}

func NewSQSProducer(handler *SQSHandler, url string) *SQSProducer {
	return &SQSProducer{Handler: handler, QueueUrl: url}
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

	err := h.Handler.SendMessage(h.QueueUrl, string(jsonData))

	if err != nil {
		slog.Error("Error trying to send message", "destination", h.QueueUrl)
	}

	return nil

}

type SnapVideoRequest struct {
	Id           uint64    `json:"id" example:"1"`
	IdUser       string    `json:"id_user" example:"1231231231"`
	FileSize     int64     `json:" file_size" example:"1048576"`
	S3FileKey    string    `json:"s3_file_key" example:"https://google.com"`
	CreationDate time.Time `json:"creation_date" example:"1970-01-01T00:00:00Z"`
}
