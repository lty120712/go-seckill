package manager

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go-chat/configs"
	"go-chat/internal/utils/logUtil"
	"mime/multipart"
	"path/filepath"
	"strings"
)

var MinioManagerInstance *MinioManager

type MinioManager struct {
	client *minio.Client
}

func InitMinIO() {
	minioConfig := configs.AppConfig.Minio
	MinioManagerInstance = &MinioManager{}
	client, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		logUtil.Errorf("MinIO 初始化失败: %v", err)
	}
	_, err = client.ListBuckets(context.Background())
	if err != nil {
		logUtil.Errorf("MinIO 连接失败: %v", err)
		return
	}
	MinioManagerInstance.client = client
	logUtil.Infof("MinIO 初始化成功")
}

// UploadToMinio 上传文件到 MinIO，并返回完整 URL
func (m *MinioManager) UploadToMinio(fileHeader *multipart.FileHeader) (string, error) {
	minioConfig := configs.AppConfig.Minio

	// 打开上传文件
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 获取后缀与分类目录
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename)) // 如 .jpg
	category := getCategoryByExt(ext)                         // 如 image

	// 生成随机文件名（UUID）
	newFileName := uuid.New().String() + ext
	objectPath := fmt.Sprintf("%s/%s", category, newFileName)

	// 上传
	_, err = m.client.PutObject(context.Background(), minioConfig.Bucket, objectPath, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	// 拼接 URL
	url := fmt.Sprintf("%s/%s/%s", minioConfig.BaseUrl, minioConfig.Bucket, objectPath)
	return url, nil
}

// getCategoryByExt 根据扩展名分类文件夹
func getCategoryByExt(ext string) string {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return "image"
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv":
		return "video"
	case ".mp3", ".wav", ".aac", ".flac":
		return "audio"
	case ".zip", ".rar", ".tar", ".gz", ".7z":
		return "archive"
	default:
		return "other"
	}
}
