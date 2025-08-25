package interfacesservice

import (
	"mime/multipart"
)

type FileServiceInterface interface {
	Upload(id uint, file *multipart.FileHeader) (url string, err error)
}
