package interfaces

import "mime/multipart"

type MinioManager interface {
	UploadToMinio(fileHeader *multipart.FileHeader) (string, error)
}
