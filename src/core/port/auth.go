package port

import "example/web-service-gin/src/core/entity"

type JwtService interface {
	GetUser(token string) (*entity.User, error)
}
