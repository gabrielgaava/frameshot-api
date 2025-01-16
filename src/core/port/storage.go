package port

import (
	"github.com/golang-migrate/migrate/v4/source/file"
	"mime/multipart"
)

type StoragePort interface {
	UploadFile(file *multipart.FileHeader, fileKey string) (string, error)
	DownloadFile(fileKey string) (*file.File, error)
}
