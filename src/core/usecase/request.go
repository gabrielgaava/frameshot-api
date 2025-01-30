package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"example/web-service-gin/src/adapters/handler/queue"
	"example/web-service-gin/src/adapters/storage/bucket"
	"example/web-service-gin/src/core/entity"
	"example/web-service-gin/src/core/port"
	"fmt"
	"mime/multipart"
	"strings"
	"time"
)

type RequestUseCase struct {
	repository port.RequestRepository
	storage    port.StoragePort
	queue      port.QueuePort
	mail       port.MailServicePort
}

// NewRequestUseCase creates a new user service instance
func NewRequestUseCase(repo port.RequestRepository, storage port.StoragePort, queue port.QueuePort, notif port.MailServicePort) *RequestUseCase {
	return &RequestUseCase{repo, storage, queue, notif}
}

func (usecase *RequestUseCase) Create(ctx context.Context, request *entity.Request, file *multipart.FileHeader) (*entity.Request, error) {

	_, err := validateFileRules(file)
	fileKeyName := generateFileKey(request.UserId, file)

	// Is a Valid File
	if err != nil {
		return nil, err
	}

	request.Status = entity.Pending
	request.CreatedAt = time.Now()
	fileKey, err := usecase.storage.UploadFile(file, fileKeyName)

	// Storage Error
	if err != nil {
		return nil, err
	}

	request.VideoKey = fileKey
	request.VideoSize = file.Size
	request, err = usecase.repository.CreateRequest(ctx, request)

	// Repository Error
	if err != nil {
		return nil, err
	}

	return request, nil

}

func (usecase *RequestUseCase) Update(ctx context.Context, request *entity.Request) (*entity.Request, error) {

	updatedRequest, err := usecase.repository.UpdateRequest(ctx, request)

	if err != nil {
		return nil, err
	}

	return updatedRequest, nil

}

func (usecase *RequestUseCase) List(ctx context.Context, userId string) ([]entity.Request, error) {
	requestList, err := usecase.repository.GetAllUserRequests(ctx, userId)

	if err != nil {
		return nil, err
	}

	if requestList == nil {
		return []entity.Request{}, nil
	}

	return requestList, nil
}

func (usecase *RequestUseCase) Get(ctx context.Context, id uint64) (*entity.Request, error) {

	request, err := usecase.repository.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return request, nil
}

func (usecase *RequestUseCase) HandleUploadNotification(ctx context.Context, msg entity.EventMessage) {

	var event bucket.S3Event
	var bodyMessage string = msg.Body

	err := json.Unmarshal([]byte(bodyMessage), &event)
	if err != nil {
		fmt.Println("Error converting body message: ", err)
		return
	}

	// Loop for each record of the S3 Event
	for _, record := range event.Records {
		if record.S3.ConfigurationID == "VideoUploaded" {
			var fileKey string = record.S3.Object.Key
			fmt.Println("Bucket:", record.S3.Bucket.Name)
			fmt.Println("Key:", record.S3.Object.Key)
			fmt.Println("Size:", record.S3.Object.Size)

			// Update status on Database
			request, _ := usecase.repository.UpdateStatusByVideoKey(ctx, string(entity.InProgress), fileKey)

			// Sent Message to SQS to Start Upload
			usecase.queue.SendVideoProccessToQueue(request)
		}
	}

}

func (usecase *RequestUseCase) HandleVideoOutputNotification(ctx context.Context, msg entity.EventMessage) {

	var notification queue.SnapVideoResponse
	var bodyMessage string = msg.Body
	var statusMessage string

	fmt.Println("Message from Queue: ", msg.Body)

	err := json.Unmarshal([]byte(bodyMessage), &notification)
	if err != nil {
		fmt.Println("Error converting body message: ", err)
		return
	}

	var isSuccess bool = notification.Status == "OK"

	videoRequest, getError := usecase.Get(ctx, notification.Id)

	if getError != nil {
		fmt.Println("Invalid request ID: ", err)
		return
	}

	videoRequest.FinishedAt = time.Now()
	videoUrl := usecase.storage.GetFileUrl(notification.S3ZipFileKey)

	if isSuccess {
		videoRequest.Status = entity.Completed
		videoRequest.ZipOutputKey = videoUrl
		statusMessage = "sucesso"
	} else {
		videoRequest.Status = entity.Failed
		statusMessage = "erro"
	}

	_, err = usecase.repository.UpdateRequest(ctx, videoRequest)

	if err != nil {
		return
	}

	fmt.Println("Sucesso: ", statusMessage)
	_ = usecase.mail.NotifyRequestStatus(videoRequest, statusMessage)
}

func validateFileRules(file *multipart.FileHeader) (bool, error) {

	var allowedExtensions = [...]string{"mp4", "mkv", "avi", "webm", "mov"}
	var fileExtension = strings.Split(file.Filename, ".")[1]
	var fileSize = (file.Size / 1024) / 1000 // Mbs

	for _, extension := range allowedExtensions {
		if extension == fileExtension {
			if fileSize <= 500 {
				return true, nil
			} else {
				return false, errors.New("file size is greater than 500Mb")
			}
		}
	}

	return false, errors.New("file extension not allowed")

}

func generateFileKey(userId string, file *multipart.FileHeader) string {
	now := time.Now().UTC().Format("2006-01-02-15-04-05")
	fileExtension := strings.Split(file.Filename, ".")[1]
	fileKey := "videos_input/" + userId + "_" + now + "." + fileExtension

	return fileKey
}
