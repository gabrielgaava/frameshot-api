package http_test

import (
	"bytes"
	"context"
	"errors"
	controller "example/web-service-gin/src/adapters/handler/http"
	"example/web-service-gin/src/core/entity"
	"example/web-service-gin/src/utils/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mocks

type MockRequestService struct {
	mock.Mock
}

func (m *MockRequestService) Create(ctx context.Context, request *entity.Request, file *multipart.FileHeader) (*entity.Request, error) {
	args := m.Called(ctx, request, file)
	return args.Get(0).(*entity.Request), args.Error(1)
}

func (m *MockRequestService) Update(ctx context.Context, request *entity.Request) (*entity.Request, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*entity.Request), args.Error(1)
}

func (m *MockRequestService) List(ctx context.Context, userId string) ([]entity.Request, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]entity.Request), args.Error(1)
}

func (m *MockRequestService) Get(ctx context.Context, id uint64) (*entity.Request, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Request), args.Error(1)
}

func (m *MockRequestService) HandleUploadNotification(ctx context.Context, msg entity.EventMessage) {
	m.Called(ctx, msg)
	return
}

func (m *MockRequestService) HandleVideoOutputNotification(ctx context.Context, msg entity.EventMessage) {
	m.Called(ctx, msg)
	return
}

// Testing Cases

func setUp(passAuth bool) (*controller.RequestHandler, *gin.Engine, *MockRequestService) {
	mockJwtService := new(mocks.MockJwtService)
	mockService := new(MockRequestService)
	handler := controller.NewRequestHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("jwtService", mockJwtService)
		c.Next()
	})

	if passAuth {
		mockJwtService.On("GetUser", "valid-token").Return(&entity.User{
			Id:    "123456",
			Email: "user@example.com",
		}, nil)
	} else {
		mockJwtService.On("GetUser", "not-valid-token").Return(
			(*entity.User)(nil),
			errors.New("invalid token"))
	}

	return handler, router, mockService
}

func generateFileForm(t *testing.T) (*bytes.Buffer, string) {

	var contentType string
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	contentType = writer.FormDataContentType()
	fileWriter, _ := writer.CreateFormFile("video_file", "test.mp4")
	fileContent := []byte("conteúdo fictício do arquivo")

	_, err := fileWriter.Write(fileContent)

	if err != nil {
		assert.Fail(t, "Should not get an error")
	}

	err = writer.Close()

	if err != nil {
		assert.Fail(t, "Should not get an error")
	}

	return body, contentType
}

func TestRequestHandler_Register(t *testing.T) {

	handler, router, service := setUp(true)

	router.POST("/requests", handler.Register)

	body, contentType := generateFileForm(t)
	service.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&entity.Request{}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/requests", body)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRequestHandler_RegisterError(t *testing.T) {

	handler, router, service := setUp(true)

	router.POST("/requests", handler.Register)

	body, contentType := generateFileForm(t)
	service.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return((*entity.Request)(nil), errors.New("error"))

	req, _ := http.NewRequest(http.MethodPost, "/requests", body)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRequestHandler_ListUsers(t *testing.T) {

	handler, router, service := setUp(true)
	router.GET("/requests", handler.ListUsers)

	request1 := mocks.MockGetRequest()
	request2 := mocks.MockGetRequest()
	request2.ID = 2
	requestList := []entity.Request{request1, request2}

	service.On("List", mock.Anything, mock.Anything).Return(requestList, nil)
	req, _ := http.NewRequest(http.MethodGet, "/requests", nil)
	req.Header.Add("Authorization", "valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestHandler_DatabaseError(t *testing.T) {

	handler, router, service := setUp(true)
	router.GET("/requests", handler.ListUsers)

	service.On("List", mock.Anything, mock.Anything).Return(([]entity.Request)(nil), errors.New("DB Error"))
	req, _ := http.NewRequest(http.MethodGet, "/requests", nil)
	req.Header.Add("Authorization", "valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRequestHandler_InvalidToken(t *testing.T) {

	handler, router, service := setUp(false)
	router.GET("/requests", handler.ListUsers)

	request1 := mocks.MockGetRequest()
	request2 := mocks.MockGetRequest()
	request2.ID = 2
	requestList := []entity.Request{request1, request2}

	service.On("List", mock.Anything, mock.Anything).Return(requestList, nil)
	req, _ := http.NewRequest(http.MethodGet, "/requests", nil)
	req.Header.Add("Authorization", "not-valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
