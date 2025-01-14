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
}

// NewRequestUseCase creates a new user service instance
func NewRequestUseCase(repo port.RequestRepository) *RequestUseCase {
	return &RequestUseCase{
		repo,
	}
}

func (service *RequestUseCase) Create(ctx context.Context, request *entity.Request, file *multipart.FileHeader) (*entity.Request, error) {

	_, err := validateFileRules(file)

	// Is a Valid File
	if err != nil {
		return nil, err
	}

	request.Status = entity.Pending
	request.CreatedAt = time.Now()
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

	requestList, err := service.repository.GetAllUserRequests(ctx, "123123123asd")

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

// KB to MB
// (file.Size / 1024) / 1000
