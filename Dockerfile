# 适配你的IM项目专用Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app

# 解决国内依赖下载超时
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV CGO_ENABLED=0

# 复制项目所有文件
COPY . .

# 下载依赖 + 编译项目
RUN go mod tidy
RUN go build -o im-server main.go

# 最小运行镜像
FROM alpine:latest
WORKDIR /app

# 同步时区(解决时间错误)
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

# 复制编译好的程序和配置文件
COPY --from=builder /app/im-server .
COPY --from=builder /app/config ./config

# 你的IM项目端口
EXPOSE 8081

# 启动项目
CMD ["./im-server"]