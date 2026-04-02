# 适配你的IM项目 最终版（解决Go版本错误）
FROM golang:1.25-alpine AS builder
WORKDIR /app

# 国内代理 + 固定Go版本
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOTOOLCHAIN=local

# 复制代码
COPY . .

# 安装依赖并编译
RUN go mod tidy
RUN go build -o im-server main.go

# 最小运行镜像
FROM alpine:latest
WORKDIR /app

# 时区
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

# 复制程序和配置
COPY --from=builder /app/im-server .
COPY --from=builder /app/config ./config

EXPOSE 8081
CMD ["./im-server"]