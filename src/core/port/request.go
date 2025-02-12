package port

import (
	"context"
	"example/web-service-gin/src/core/entity"
	"mime/multipart"
)

type RequestRepository interface {

	// CreateRequest creates a new Request on database and return it
	CreateRequest(ctx context.Context, request *entity.Request) (*entity.Request, error)

	// GetById searchs for a request with informed ID
	GetById(ctx context.Context, id uint64) (*entity.Request, error)

	//GetAllUserRequests returns a list of all user requests
	GetAllUserRequests(ctx context.Context, userId string) ([]entity.Request, error)

	//UpdateRequest updates the role request entity
	UpdateRequest(ctx context.Context, request *entity.Request) (*entity.Request, error)

	//UpdateStatusByVideoKey updates the status of a request in database looking for the file key name
	UpdateStatusByVideoKey(ctx context.Context, status string, videoKey string) (*entity.Request, error)
}

type RequestService interface {
	Create(ctx context.Context, request *entity.Request, file *multipart.FileHeader) (*entity.Request, error)
	Update(ctx context.Context, request *entity.Request) (*entity.Request, error)
	List(ctx context.Context, userId string) ([]entity.Request, error)
	Get(ctx context.Context, id uint64) (*entity.Request, error)
	HandleUploadNotification(ctx context.Context, msg entity.EventMessage)
	HandleVideoOutputNotification(ctx context.Context, msg entity.EventMessage)
}

type QueuePort interface {
	//SendVideoProccessToQueue insert a new conversion request to the queue
	SendVideoProccessToQueue(request *entity.Request) error
}

type MailServicePort interface {
	// NotifyRequestStatus notify when a video request is converted with
	//success by the service Or with any error
	NotifyRequestStatus(request *entity.Request, status string) error
}
