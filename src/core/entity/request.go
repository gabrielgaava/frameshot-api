package entity

import "time"

type RequestStatus string

const (
	Pending    RequestStatus = "PENDING"
	InProgress RequestStatus = "IN_PROGRESS"
	Completed  RequestStatus = "COMPLETED"
	Failed     RequestStatus = "FAILED"
)

type Request struct {
	ID           uint64
	UserId       string
	UserEmail    string
	VideoSize    int64
	VideoKey     string
	ZipOutputKey *string
	Status       RequestStatus
	CreatedAt    time.Time
	FinishedAt   *time.Time
}
