# 使用 Golang 官方镜像作为构建环境
FROM golang:1.24 AS builder

# 设置工作目录
WORKDIR /go-chat

# 复制 Go 项目的文件到容器内
COPY . .

# 获取依赖并编译 Go 项目
RUN go mod tidy
RUN go build -o go-chat .

# 使用 jrottenberg/ffmpeg 镜像，包含 FFmpeg 环境
FROM jrottenberg/ffmpeg:4.3

# 设置工作目录
WORKDIR /go-chat

# 从构建阶段复制编译好的二进制文件到最终镜像
COPY --from=builder /go-chat/go-chat /usr/local/bin/

# 配置容器的启动命令
ENTRYPOINT ["go-chat"]

# 配置容器开放端口
EXPOSE 80 9090
