package http

import (
	"example/web-service-gin/src/core/entity"
	"example/web-service-gin/src/core/port"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
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

// Register godoc
//
//	@Summary		Register a new video conversion request
//	@Description	create a new video conversion request with default data
//	@Tags			Requests
//	@Accept			json
//	@Produce		json
//	@Param			registerRequest	body		registerRequest	true	"Register request"
//	@Success		200				{object}	userResponse	"User created"
//	@Failure		400				{object}	errorResponse	"Validation error"
//	@Failure		401				{object}	errorResponse	"Unauthorized error"
//	@Failure		404				{object}	errorResponse	"Data not found error"
//	@Failure		409				{object}	errorResponse	"Data conflict error"
//	@Failure		500				{object}	errorResponse	"Internal server error"
//	@Router			/users [post]
func (handler *RequestHandler) Register(ctx *gin.Context) {
	form, _ := ctx.MultipartForm()
	file, _ := ctx.FormFile("video_file")
	log.Println(file.Filename)
	log.Println(form)

	request := entity.Request{
		UserId:   form.Value["user_id"][0],
		VideoUrl: form.Value["video_url"][0],
	}

	_, err := handler.service.Create(ctx, &request, file)

	if err != nil {
		// handleError(ctx, err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rsp := newRequestResponse(&request)

	handleSuccess(ctx, rsp)
}

func (handler *RequestHandler) ListUsers(ctx *gin.Context) {

	log.Println("Controler do GET")

	var requestList []requestResponse

	requests, err := handler.service.List(ctx, "123123123asd")
	if err != nil {
		//handleError(ctx, err)
		return
	}

	for _, request := range requests {
		requestList = append(requestList, newRequestResponse(&request))
	}

	total := uint64(len(requestList))
	log.Printf("Registros: %d\n", total)

	handleSuccess(ctx, requestList)
}

// response represents a response body format
type response struct {
	Success bool `json:"success" example:"true"`
	Error   any  `json:"error,omitempty"`
	Data    any  `json:"data,omitempty"`
}

// userResponse represents a user response body
type requestResponse struct {
	ID         uint64               `json:"id" example:"1"`
	UserId     string               `json:"user_id" example:"1231231231"`
	VideoUrl   string               `json:"video_url" example:"https://google.com"`
	Status     entity.RequestStatus `json:"status" example:"PENDING"`
	CreatedAt  time.Time            `json:"created_at" example:"1970-01-01T00:00:00Z"`
	FinishedAt *time.Time           `json:"finished_at" example:"1970-01-01T00:00:00Z"`
}

// validationError sends an error response for some specific request validation error
func validationError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, err.Error())
}

func newResponse(data any) response {
	return response{
		Success: true,
		Error:   nil,
		Data:    data,
	}
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(ctx *gin.Context, data any) {
	rsp := newResponse(data)
	ctx.JSON(http.StatusOK, rsp)
}

func newRequestResponse(request *entity.Request) requestResponse {
	return requestResponse{
		ID:         request.ID,
		UserId:     request.UserId,
		VideoUrl:   request.VideoUrl,
		Status:     request.Status,
		CreatedAt:  request.CreatedAt,
		FinishedAt: request.FinishedAt,
	}
}
