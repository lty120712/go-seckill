package main

import (
	"go-chat/internal/app"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @BasePath  /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	app.Start()
}
