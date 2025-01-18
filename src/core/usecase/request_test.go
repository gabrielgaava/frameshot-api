package usecase_test

import (
	"context"
	"errors"
	"example/web-service-gin/src/core/entity"
	"example/web-service-gin/src/core/usecase"
	"example/web-service-gin/src/utils"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
	"testing"
	"time"
)

type MockRequestRepository struct {
	mock.Mock
}

type MockStoragePort struct {
	mock.Mock
}

type MockRequestNotifications struct {
	mock.Mock
}

func (m *MockRequestRepository) CreateRequest(ctx context.Context, request *entity.Request) (*entity.Request, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*entity.Request), args.Error(1)
}

func (m *MockRequestRepository) UpdateRequest(ctx context.Context, request *entity.Request) (*entity.Request, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*entity.Request), args.Error(1)
}

func (m *MockRequestRepository) UpdateStatusByVideoKey(ctx context.Context, status string, videoKey string) (*entity.Request, error) {
	args := m.Called(ctx, status, videoKey)
	return args.Get(0).(*entity.Request), args.Error(1)
}

func (m *MockRequestRepository) GetAllUserRequests(ctx context.Context, userId string) ([]entity.Request, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]entity.Request), args.Error(1)
}

func (m *MockRequestRepository) GetById(ctx context.Context, id uint64) (*entity.Request, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Request), args.Error(1)
}

func (m *MockStoragePort) UploadFile(file *multipart.FileHeader, fileKey string) (string, error) {
	args := m.Called(file, fileKey)
	return args.String(0), args.Error(1)
}

func (m *MockStoragePort) DownloadFile(fileKey string) (*file.File, error) {
	args := m.Called(fileKey)
	return args.Get(0).(*file.File), args.Error(1)
}

func (m *MockRequestNotifications) SendVideoProccessToQueue(request *entity.Request) error {
	args := m.Called(request)
	return args.Error(0)
}

func setUp() (*MockRequestRepository, *MockStoragePort, *MockRequestNotifications, *usecase.RequestUseCase) {
	mockRepo := new(MockRequestRepository)
	mockStorage := new(MockStoragePort)
	mockNotification := new(MockRequestNotifications)
	requestUsecase := usecase.NewRequestUseCase(mockRepo, mockStorage, mockNotification)

	return mockRepo, mockStorage, mockNotification, requestUsecase
}

func TestCreateRequest_Success(t *testing.T) {

	mockRepo, mockStorage, _, requestUsecase := setUp()
	ctx := context.Background()
	request := &entity.Request{
		UserId: "user123",
	}
	videoFile := &multipart.FileHeader{
		Filename: "video.mp4",
		Size:     100 * 1024 * 1024, // 100 MB
	}

	fileKey := "videos_input/user123_2023-01-01-10-00-00.mp4"
	mockStorage.On("UploadFile", videoFile, mock.Anything).Return(fileKey, nil)
	mockRepo.On("CreateRequest", ctx, mock.Anything).Return(request, nil)

	createdRequest, err := requestUsecase.Create(ctx, request, videoFile)

	assert.NoError(t, err)
	assert.NotNil(t, createdRequest)
	mockStorage.AssertCalled(t, "UploadFile", videoFile, mock.AnythingOfType("string"))
	mockRepo.AssertCalled(t, "CreateRequest", ctx, mock.Anything)
}

func TestCreateRequest_InvalidFileExtension(t *testing.T) {
	_, _, _, requestUsecase := setUp()

	ctx := context.Background()
	request := &entity.Request{
		UserId: "user123",
	}
	videoFile := &multipart.FileHeader{
		Filename: "video.exe",
		Size:     100 * 1024 * 1024, // 100 MB
	}

	createdRequest, err := requestUsecase.Create(ctx, request, videoFile)

	assert.Error(t, err)
	assert.Nil(t, createdRequest)
	assert.EqualError(t, err, "file extension not allowed")
}

func TestCreateRequest_InvalidFileSize(t *testing.T) {
	_, _, _, requestUsecase := setUp()

	ctx := context.Background()
	request := &entity.Request{
		UserId: "user123",
	}
	videoFile := &multipart.FileHeader{
		Filename: "video.avi",
		Size:     800 * 1024 * 1024, // 100 MB
	}

	createdRequest, err := requestUsecase.Create(ctx, request, videoFile)

	assert.Error(t, err)
	assert.Nil(t, createdRequest)
	assert.EqualError(t, err, "file size is greater than 500Mb")
}

func TestUpdateRequest_Success(t *testing.T) {

	mockRepo, _, _, requestUsecase := setUp()
	ctx := context.Background()
	moment := time.Now()

	requestMock := &entity.Request{
		ID:         12,
		UserId:     "user123",
		Status:     entity.Completed,
		FinishedAt: moment,
	}

	mockRepo.On("UpdateRequest", ctx, requestMock).Return(requestMock, nil)
	updatedRequest, err := requestUsecase.Update(ctx, requestMock)

	assert.NoError(t, err)
	assert.NotNil(t, updatedRequest)
	mockRepo.AssertCalled(t, "UpdateRequest", ctx, requestMock)
	assert.Equal(t, requestMock.Status, updatedRequest.Status)
	assert.Equal(t, requestMock.FinishedAt, updatedRequest.FinishedAt)
}

func TestUpdateRequest_Error(t *testing.T) {

	mockRepo, _, _, requestUsecase := setUp()
	ctx := context.Background()
	expectedError := errors.New("mock error")
	requestMock := &entity.Request{
		ID: 12,
	}

	mockRepo.On("UpdateRequest", ctx, requestMock).Return((*entity.Request)(nil), expectedError)
	updatedRequest, err := requestUsecase.Update(ctx, requestMock)

	assert.Error(t, err)
	assert.Nil(t, updatedRequest)
	mockRepo.AssertCalled(t, "UpdateRequest", ctx, requestMock)
}

func TestListRequest_Success(t *testing.T) {
	mockRepo, _, _, requestUsecase := setUp()
	ctx := context.Background()

	// Given:
	userId := "user123"
	request1 := &entity.Request{ID: 1}
	request2 := &entity.Request{ID: 2}
	expectedList := []entity.Request{*request1, *request2}

	// When:
	mockRepo.On("GetAllUserRequests", ctx, userId).Return(expectedList, nil)
	data, err := requestUsecase.List(ctx, userId)

	// Then:
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, len(expectedList), len(data))
	mockRepo.AssertCalled(t, "GetAllUserRequests", ctx, userId)
}

func TestListRequest_Error(t *testing.T) {
	mockRepo, _, _, requestUsecase := setUp()
	ctx := context.Background()
	expectedError := errors.New("mock error")
	userId := "user123"

	mockRepo.On("GetAllUserRequests", ctx, userId).Return(([]entity.Request)(nil), expectedError)
	updatedRequest, err := requestUsecase.List(ctx, userId)

	assert.Error(t, err)
	assert.Nil(t, updatedRequest)
	mockRepo.AssertCalled(t, "GetAllUserRequests", ctx, userId)
}

func TestGetRequest_Success(t *testing.T) {
	mockRepo, _, _, requestUsecase := setUp()
	ctx := context.Background()

	// Given
	var requestId uint64 = 32
	expectedRequest := &entity.Request{ID: requestId}

	// When
	mockRepo.On("GetById", ctx, requestId).Return(expectedRequest, nil)
	getData, err := requestUsecase.Get(ctx, requestId)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, getData)
	assert.Equal(t, expectedRequest, getData)
	mockRepo.AssertCalled(t, "GetById", ctx, requestId)
}

func TestHandleUploadNotification_Success(t *testing.T) {
	repo, _, notify, use := setUp()
	ctx := context.Background()

	// Given
	fileKey := "video_input/test.mp4"
	s3BodyMock := utils.GetMockS3EventBody()
	event := entity.EventMessage{
		MessageID: "123",
		Body:      s3BodyMock,
		Source:    "Test",
		Date:      time.Now().String(),
	}

	updatedRequest := &entity.Request{Status: entity.InProgress}

	// When
	repo.On("UpdateStatusByVideoKey", ctx, string(entity.InProgress), fileKey).Return(updatedRequest, nil)
	notify.On("SendVideoProccessToQueue", updatedRequest).Return(nil)
	use.HandleUploadNotification(ctx, event)

	// Then
	repo.AssertCalled(t, "UpdateStatusByVideoKey", ctx, string(entity.InProgress), fileKey)
	notify.AssertCalled(t, "SendVideoProccessToQueue", updatedRequest)
}

func TestHandleUploadNotification_InvalidBody(t *testing.T) {
	repo, _, notify, use := setUp()
	ctx := context.Background()

	// Given
	s3BodyMock := "teste"
	event := entity.EventMessage{Body: s3BodyMock}

	// When
	use.HandleUploadNotification(ctx, event)

	// Then
	repo.AssertNotCalled(t, "UpdateStatusByVideoKey")
	notify.AssertNotCalled(t, "SendVideoProccessToQueue")
}

func TestHandleVideoOutputNotification_UploadError(t *testing.T) {
	repo, _, _, use := setUp()
	ctx := context.Background()

	// Given
	var id uint64 = 1
	notificationBody := utils.GetMockOutputVideoEventBody("ERROR") // ID = 1
	message := entity.EventMessage{Body: notificationBody}
	request := utils.GetRequest()

	// When
	repo.On("GetById", ctx, id).Return(&request, nil)
	repo.On("UpdateRequest", ctx, mock.Anything).Return(&request, nil)
	use.HandleVideoOutputNotification(ctx, message)

	// Then
	repo.AssertCalled(t, "GetById", ctx, id)
	repo.AssertCalled(t, "UpdateRequest", ctx, mock.AnythingOfType("*entity.Request"))

}

func TestHandleVideoOutputNotification_UpdateError(t *testing.T) {
	repo, _, _, use := setUp()
	ctx := context.Background()

	// Given
	var id uint64 = 1
	notificationBody := utils.GetMockOutputVideoEventBody("ERROR") // ID = 1
	message := entity.EventMessage{Body: notificationBody}
	request := utils.GetRequest()

	// When
	repo.On("GetById", ctx, id).Return(&request, nil)
	repo.On("UpdateRequest", ctx, mock.Anything).Return((*entity.Request)(nil), errors.New("mock error"))
	use.HandleVideoOutputNotification(ctx, message)

	// Then
	repo.AssertCalled(t, "GetById", ctx, id)
	repo.AssertNotCalled(t, "UpdateRequest")

}

func TestHandleVideoOutputNotification_InvalidBody(t *testing.T) {
	repo, _, _, use := setUp()
	ctx := context.Background()

	// Given
	notificationBody := "invalid body message" // ID = nil
	message := entity.EventMessage{Body: notificationBody}

	// When
	use.HandleVideoOutputNotification(ctx, message)

	// Then
	repo.AssertNotCalled(t, "GetById")
	repo.AssertNotCalled(t, "UpdateRequest")
}

func TestHandleVideoOutputNotification_InvalidId(t *testing.T) {
	repo, _, _, use := setUp()
	ctx := context.Background()

	// Given
	var id uint64 = 1
	notificationBody := utils.GetMockOutputVideoEventBody("OK") // ID = 1
	message := entity.EventMessage{Body: notificationBody}

	// When
	repo.On("GetById", ctx, id).Return((*entity.Request)(nil), errors.New("invalid id"))
	use.HandleVideoOutputNotification(ctx, message)

	// Then
	repo.AssertNotCalled(t, "GetById")
	repo.AssertNotCalled(t, "UpdateRequest")
}
