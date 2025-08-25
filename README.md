# Go Chat IM 系统模板项目

这是一个基于 **Gin** 框架的即时通讯系统（IM）模板项目，集成了常用后端中间件和组件，如 JWT、GORM、Redis、RabbitMQ、WebSocket、定时任务、Swagger 文档等，具备良好的分层结构与可扩展性，适合用于快速搭建中小型 Golang 后端项目。

## 🚀 项目特性

- ✅ Gin + GORM + Viper 基础架构
- ✅ JWT 登录认证中间件
- ✅ Redis 缓存/消息支持
- ✅ RabbitMQ 消息队列（Direct 模式）
- ✅ WebSocket 实时通信
- ✅ 定时任务（基于 robfig/cron）
- ✅ Swagger API 文档生成
- ✅ 完整目录结构和配置分离
- ✅ 热更新（使用 `fresh`）

---

## 📁 项目结构

```text
my-project/
├── configs/         # 配置文件和读取逻辑（使用 Viper）
├── internal/
│   ├── api/v1/      # 路由注册、API 分发
│   ├── app/         # 项目入口
│   ├── consumer/    # 消费者
│   ├── controller/  # 控制器
│   ├── db/          # 数据库初始化（MySQL/Redis）
│   ├── handler/     # 处理器(ws消息处理等)
│   ├── manager/     # 第三方服务管理（rabbitmq,websocket）
│   ├── middleware/  # JWT 等中间件
│   ├── model/       # 请求/响应结构体、实体模型
│   ├── service/     # 核心业务逻辑
│   ├── repository/  # 数据库访问封装
│   ├── utils/       # 工具函数（如 JWT 工具）
│   └── timer/       # 定时任务
├── docs/            # Swagger 文档
├── scripts/         # 启动脚本
├── tests/           # 单元测试
├── main.go          # 启动入口
└── go.mod
```
## 🛠️ 核心组件集成说明

### Gin 框架初始化

- 启动入口函数：`internal/app/bootstrap.go`
- 统一注册中间件、路由、文档、数据库、服务等
- 在 `main.go` 中通过 `app.Start()` 启动服务

---

### GORM + MySQL

- 数据源配置读取于 `configs/app.yaml`
- 通过 `db.InitMysql()` 初始化连接
- 支持连接池参数配置、结构体自动映射、日志打印等

---

### Viper 配置

- 支持热更新（文件变更自动生效）
- 环境变量自动读取
- 多级嵌套结构体绑定配置项
- `app.yaml` 支持配置端口、数据库、JWT、Redis、RabbitMQ、消息队列等模块

---

### JWT 登录认证

- 配置文件路径：`configs.AppConfig.Jwt`
- 中间件路径：`middleware/AuthMiddleware.go`
- 工具函数定义：`utils/jwt.go`
- 控制器中使用示例：

```go
userApi.POST("/update_info", AuthMiddleware(), UserController{}.UpdateInfo)
```
### Redis 缓存支持
- 支持连接池配置与上下文超时控制

- 常见命令封装（SET/GET/HSET/ZRANGE 等）

- 测试文件示例：tests/redis_test.go

### Swagger 接口文档
- 文档目录：docs/

- 注解写在 controller 结构体方法上

- main.go 添加 Swag 元信息注释：

```go
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
```
启动后访问：http://localhost:8080/swagger/index.html

### RabbitMQ 消息中间件
- 组件管理于 manager/rabbitmqManager.go

- 支持发送与消费（Direct 交换机）

- 消费者映射注册：consumer/consumerMap.go

- 配置消费队列列表于 app.yaml 下的 mq 字段

### WebSocket 实时通信
- 管理器：manager/WebSocketManager.go

- 启动监听 /ws 地址，基于用户 ID 建立连接

- 支持发送消息到：
指定用户 多个用户 所有用户

### 定时任务
- 调度器：timer/Timer.go

- 每个任务单独封装函数，如 HeartBeatTimer

- 使用 cron 表达式调度任务

## 📦 安装依赖与运行项目
拉取项目后安装依赖：

```bash
go mod tidy
```
安装开发工具：

```bash

go install github.com/pilu/fresh@latest # 实时热更新
go install github.com/swaggo/swag/cmd/swag@latest # 接口文档生成
```
生成接口文档：
```bash
swag init
```
启动服务：

```bash
./scripts/start.bat
```

启动 WebSocket（已在 Start() 中自动初始化）
默认监听地址：:80 使用示例：

```bash
ws://localhost/ws?id=用户ID
```