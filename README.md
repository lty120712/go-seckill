# Go-Seckill 秒杀商城项目

这是一个基于 **Gin** 框架的高并发秒杀商城项目，集成了常用后端中间件和组件，如 GORM、Redis、RabbitMQ、定时任务、Swagger 文档等，支持 **现时秒杀** 与 **优惠券功能**，具备良好的分层结构与可扩展性，适合用于快速搭建电商类促销场景的 Golang 后端项目。

## 🚀 项目特性

- ✅ Gin + GORM + Viper 基础架构
- ✅ Redis 秒杀库存预扣 + 原子脚本防超卖
- ✅ RabbitMQ 消息队列削峰填谷
- ✅ MySQL 订单落库 + 乐观锁扣减库存
- ✅ 优惠券功能（领取、使用、过期管理）
- ✅ 秒杀活动功能（活动时间、库存、价格）
- ✅ Swagger API 文档生成
- ✅ 定时任务（基于 robfig/cron）
- ✅ 完整目录结构和配置分离

---

## 📁 项目结构

```text
go-seckill/
├── bin/                 # 可执行文件/编译产物
├── configs/             # 配置文件 (app.yaml, mq.yaml 等)
├── docs/                # Swagger、架构图、接口文档
├── internal/            # 内部模块（核心业务）
│   ├── api/             # 路由注册与分发
│   │   └── v1/          # v1 版本 API
│   ├── app/             # 项目启动/初始化入口
│   ├── consumer/        # MQ 消费者（订单落库、券过期处理）
│   ├── controller/      # 控制器 (HTTP Handler)
│   ├── db/              # 数据库初始化 (MySQL/Redis)
│   ├── interfaces/      # 接口定义 (如 service 接口、repository 接口)
│   ├── manager/         # 第三方服务管理 (rabbitmq, websocket, minio)
│   ├── middleware/      # 鉴权、跨域、限流中间件
│   ├── model/           # 数据模型 (实体/DTO/请求响应)
│   ├── repository/      # 数据访问层 (DAO，GORM 操作封装)
│   ├── service/         # 核心业务逻辑
│   │   ├── seckill/     # 秒杀业务
│   │   ├── coupon/      # 优惠券业务
│   │   └── order/       # 订单业务
│   ├── timer/           # 定时任务 (券过期清理、活动结束处理)
│   ├── utils/           # 工具函数 (JWT, 加密, 日志)
│   └── ws/              # WebSocket 模块 (如果用得上)
├── scripts/             # 启动/部署脚本
├── tests/               # 单元测试/集成测试
├── tmp/                 # 临时文件
├── .gitignore
├── Dockerfile
├── generate_users.sql   # 初始化用户 SQL
├── go.mod
├── main.go              # 入口
└── README.md

```

## 🛠️ 功能模块说明

### 现时秒杀
- 支持活动配置（开始时间、结束时间、秒杀价格、库存数量）
- Redis Lua 脚本实现库存预扣与去重（防止超卖）
- RabbitMQ 队列异步下单，削峰填谷
- MySQL 落库，保证订单与库存一致性

### 优惠券功能
- 优惠券发放：支持数量、时间范围配置
- 优惠券领取：Redis 控制并发领取，防止超发
- 优惠券使用：下单时校验是否可叠加秒杀商品
- 优惠券过期：Redis TTL + 定时任务，最终一致性更新 MySQL 状态

### Swagger 接口文档
- 文档目录：docs/
- 注解写在 controller 方法上
- 启动后访问：http://localhost:8080/swagger/index.html

### 定时任务
- 调度器：timer/Timer.go
- 用于处理优惠券过期、秒杀活动结束等任务

---

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

默认监听地址：:8080
