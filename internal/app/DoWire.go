package app

import (
	"github.com/sirupsen/logrus"
	controllers "go-chat/internal/controller"
	"go-chat/internal/manager"
	"go-chat/internal/repository"
	"go-chat/internal/service"
)

func doWire() {
	//repo
	repository.InitUserRepository()

	repository.InitFileRepository()
	//ws

	//service
	service.InitUserService(repository.UserRepositoryInstance)

	service.InitFileService(repository.FileRepositoryInstance, manager.MinioManagerInstance)
	//controller
	controllers.InitUserController(service.UserServiceInstance)
	controllers.InitFileController(service.FileServiceInstance)
	//延迟注入

	logrus.Info("=======================依赖注入完成=====================")
}
