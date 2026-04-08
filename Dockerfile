FROM golang:1.25.3-alpine AS builder
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOTOOLCHAIN=local

COPY . .

RUN go mod tidy
RUN go build -o im-server main.go

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

COPY --from=builder /app/im-server .
COPY --from=builder /app/config ./config

EXPOSE 8081
CMD ["./im-server"]