# 构建阶段
FROM golang:1.22.0 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
# 设置 Go 模块代理并下载依赖
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download


# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates，这对于 HTTPS 请求是必要的
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件和静态文件
COPY --from=builder /app/config ./config
COPY --from=builder /app/index.html .

# 设置环境变量
ENV GIN_MODE=release

# 暴露端口
EXPOSE 10025

# 运行应用
CMD ["./main --env dev"] 