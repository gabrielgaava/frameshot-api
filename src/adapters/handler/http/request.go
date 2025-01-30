package http

import (
	"example/web-service-gin/src/core/entity"
	"example/web-service-gin/src/core/port"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	service port.RequestService
}

func NewRequestHandler(service port.RequestService) *RequestHandler {
	return &RequestHandler{
		service,
	}
}

type CreateRequestBody struct {
	UserId   string `json:"user_id" binding:"required" example:"123456"`
	VideoUrl string `json:"video_url" binding:"required" example:"https://google.com"`
}

func (r *RequestHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok", "message": "service is running"})
}

func (handler *RequestHandler) Register(ctx *gin.Context) {

	user := getAuthUser(ctx)
	file, _ := ctx.FormFile("video_file")
	log.Println(file.Filename)

	request := entity.Request{
		UserId:    user.Id,
		UserEmail: user.Email,
	}

	createdRequest, createError := handler.service.Create(ctx, &request, file)

	if createError != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "request not saved. Try again later."})
		return
	}

	rsp := newRequestResponse(createdRequest)
	ctx.JSON(http.StatusCreated, rsp)
}

func (handler *RequestHandler) ListUsers(ctx *gin.Context) {

	user := getAuthUser(ctx)

	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var requestList []requestResponse

	requests, err := handler.service.List(ctx, user.Id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "request not saved. Try again later."})
		return
	}

	if len(requests) == 0 {
		ctx.JSON(http.StatusNoContent, "")
		return
	}

	for _, request := range requests {
		requestList = append(requestList, newRequestResponse(&request))
	}

	ctx.JSON(http.StatusOK, requestList)
}

func getAuthUser(ctx *gin.Context) *entity.User {
	jwtServiceInterface, _ := ctx.Get("jwtService")
	jwtService := jwtServiceInterface.(port.JwtService)

	jwtToken := ctx.Request.Header.Get("Authorization")
	user, err := jwtService.GetUser(jwtToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return nil
	}

	return user
}

type requestResponse struct {
	ID           uint64               `json:"id" example:"1"`
	UserId       string               `json:"user_id" example:"1231231231"`
	UserEmail    string               `json:"user_email" example:"user@example.com"`
	VideoSize    int64                `json:"video_size" example:"1048576"`
	VideoKey     string               `json:"video_url" example:"https://google.com"`
	ZipOutputKey string               `json:"zip_output_key" example:"123456"`
	Status       entity.RequestStatus `json:"status" example:"PENDING"`
	CreatedAt    time.Time            `json:"created_at" example:"1970-01-01T00:00:00Z"`
	FinishedAt   time.Time            `json:"finished_at" example:"1970-01-01T00:00:00Z"`
}

func newRequestResponse(request *entity.Request) requestResponse {
	return requestResponse{
		ID:           request.ID,
		UserId:       request.UserId,
		UserEmail:    request.UserEmail,
		VideoSize:    request.VideoSize,
		VideoKey:     request.VideoKey,
		ZipOutputKey: request.ZipOutputKey,
		Status:       request.Status,
		CreatedAt:    request.CreatedAt,
		FinishedAt:   request.FinishedAt,
	}
}
