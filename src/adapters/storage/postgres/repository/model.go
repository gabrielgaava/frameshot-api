package repository

import (
	"database/sql"
	"time"
)

type RequestModel struct {
	ID           uint64
	UserId       string
	UserEmail    string
	VideoSize    int64
	VideoKey     string
	ZipOutputKey sql.NullString
	Status       string
	CreatedAt    time.Time
	FinishedAt   sql.NullTime
}
