package port

import (
	"context"
	"example/web-service-gin/src/core/entity"
	"mime/multipart"
)

type RequestRepository interface {

	// CreateRequest creates a new Request on database and return it
	CreateRequest(ctx context.Context, request *entity.Request) (*entity.Request, error)

	//GetAllUserRequests returns a list of all user requests
	GetAllUserRequests(ctx context.Context, userId string) ([]entity.Request, error)

	//UpdateRequest updates the role request entity
	UpdateRequest(ctx context.Context, request *entity.Request) (*entity.Request, error)
}

type RequestService interface {
	Create(ctx context.Context, request *entity.Request, file *multipart.FileHeader) (*entity.Request, error)
	Update(ctx context.Context, request *entity.Request) (*entity.Request, error)
	List(ctx context.Context, userId string) ([]entity.Request, error)
	Get(ctx context.Context, id string) (*entity.Request, error)
}
