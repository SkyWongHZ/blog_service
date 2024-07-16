

# 使用官方 Golang 1.22.5镜像作为基础镜像
FROM hub.atomgit.com/amd64/golang:1.19 as builder


# 设置工作目录
WORKDIR /app

# 安装 Git
RUN apt-get update && apt-get install -y git

# 设置 Go 代理
ENV GOPROXY=https://goproxy.cn,direct


# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 将源代码复制到工作目录
COPY . .

# 构建 Go 应用
RUN go build -o main .

# 使用一个更小的基础镜像
FROM hub.atomgit.com/amd64/alpine:3.15.10

# 安装必要的 ca-certificates
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从 builder 镜像复制构建的二进制文件
COPY --from=builder /app/main .

# 暴露应用程序端口（可选，根据你的应用需求）
EXPOSE 8080

# 运行 Go 应用
CMD ["./main"]

