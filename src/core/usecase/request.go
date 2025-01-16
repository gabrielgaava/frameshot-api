package usecase

import (
	"context"
	"errors"
	"example/web-service-gin/src/core/entity"
	"example/web-service-gin/src/core/port"
	"mime/multipart"
	"strings"
	"time"
)

type RequestUseCase struct {
	repository port.RequestRepository
	storage    port.StoragePort
}

// NewRequestUseCase creates a new user service instance
func NewRequestUseCase(repo port.RequestRepository, storage port.StoragePort) *RequestUseCase {
	return &RequestUseCase{
		repo,
		storage,
	}
}

func (service *RequestUseCase) Create(ctx context.Context, request *entity.Request, file *multipart.FileHeader) (*entity.Request, error) {

	_, err := validateFileRules(file)
	fileKeyName := generateFileKey(request.UserId, file)

	// Is a Valid File
	if err != nil {
		return nil, err
	}

	request.Status = entity.Pending
	request.CreatedAt = time.Now()
	fileKey, err := service.storage.UploadFile(file, fileKeyName)

	// Storage Error
	if err != nil {
		return nil, err
	}

	request.VideoKey = fileKey
	request.VideoSize = file.Size
	createdRequest, err := service.repository.CreateRequest(ctx, request)

	// Repository Error
	if err == nil {
		return nil, err
	}

	return createdRequest, nil

}

func (service *RequestUseCase) Update(ctx context.Context, request *entity.Request) (*entity.Request, error) {
	return nil, nil
}

func (service *RequestUseCase) List(ctx context.Context, userId string) ([]entity.Request, error) {
	requestList, err := service.repository.GetAllUserRequests(ctx, userId)

	if err != nil {
		return nil, err
	}

	return requestList, nil
}

func (service *RequestUseCase) Get(ctx context.Context, id string) (*entity.Request, error) {
	return nil, nil
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

// KB to MB
// (file.Size / 1024) / 1000
