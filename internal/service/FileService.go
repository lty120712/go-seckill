package service

import (
	interfacemanager "go-chat/internal/interfaces/manager"
	interfacerepository "go-chat/internal/interfaces/repository"
	fileutil "go-chat/internal/utils/fileUtil"
	"mime/multipart"
)

type FileService struct {
	fileRepository interfacerepository.FileRepositoryInterface
	minioManager   interfacemanager.MinioManager
}

var FileServiceInstance *FileService

func InitFileService(fileRepository interfacerepository.FileRepositoryInterface,
	minioManager interfacemanager.MinioManager) {
	FileServiceInstance = &FileService{
		fileRepository: fileRepository,
		minioManager:   minioManager,
	}
}

func (s *FileService) Upload(id uint, file *multipart.FileHeader) (url string, err error) {
	parseFile, err := fileutil.ParseFile(file)
	if err != nil {
		return "", err
	}
	parseFile.ID = id
	//上传minio获取url
	url, err = s.minioManager.UploadToMinio(file)
	parseFile.Url = url
	//保存数据库
	err = s.fileRepository.Create(parseFile)
	return url, nil
}
