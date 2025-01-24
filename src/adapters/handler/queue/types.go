package queue

import "time"

type SnapVideoRequest struct {
	Id           uint64    `json:"id" example:"1"`
	IdUser       string    `json:"id_user" example:"1231231231"`
	FileSize     int64     `json:" file_size" example:"1048576"`
	S3FileKey    string    `json:"s3_file_key" example:"https://google.com"`
	CreationDate time.Time `json:"creation_date" example:"1970-01-01T00:00:00Z"`
}

type SnapVideoResponse struct {
	Id           uint64    `json:"id" example:"1"`
	IdUser       string    `json:"id_user" example:"1231231231"`
	Status       string    `json:"status" example:"SUCCESS"`
	S3ZipFileKey string    `json:"s3_zip_file_key" example:"https://google.com"`
	CreationDate time.Time `json:"creation_date" example:"2025-01-23T20:38:08.792075"`
	FinishedDate time.Time `json:"finished_date" example:"2025-01-23T20:38:08.792075"`
}
