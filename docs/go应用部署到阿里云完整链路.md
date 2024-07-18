### 部署博客系统到云服务器的完整流程
#### 1. 准备工作
##### 连接到云服务器：
```shell
ssh root@<your-cloud-server-ip>
```
##### 更新包管理器：
```shell
apt-get update -y && apt-get upgrade -y
```
#### 2. 设置远程服务器环境
##### 更新好升级系统包
```go
sudo apt update
sudo apt upgrade -y
```
##### 安装 Docker（CentOS系统用yum来安装）：
```shell
yum install -y docker-ce docker-ce-cli containerd.io
```
##### 启动 Docker 并设置开机自启：
```shell
systemctl start docker
sudo systemctl enable docker
```
##### 验证 Docker 安装：
```shell
docker --version
```
##### 登录docker
```shell
docker login
```
##### 安装mysql

1. **拉取 MySQL Docker 镜像**：
```shell
docker pull mysql:latest
```

2. **运行 MySQL 容器**：
```shell
docker run --name mysql-container -e MYSQL_ROOT_PASSWORD=yourpassword -d mysql:latest
```
替换 `yourpassword` 为你设置的 MySQL 根密码

####  3 编写 Dockerfile
在项目根目录下创建一个 `Dockerfile`：
```dockerfile


# 使用官方 Golang 1.22.5镜像作为基础镜像
FROM hub.atomgit.com/amd64/golang:1.21 as builder


# 设置工作目录
WORKDIR /app

# 安装 Git
RUN apt-get update && apt-get install -y git

RUN git --version

# 设置 Go 代理
ENV GOPROXY=https://goproxy.cn,direct

# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod go.sum ./

# 清理和验证依赖项
RUN go mod tidy

# 下载依赖
RUN go mod download

# 将源代码复制到工作目录
COPY . .

# 构建 Go 应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 使用 Debian 作为运行阶段的基础镜像
# FROM debian:bullseye-slim

# 使用一个更小的基础镜像
FROM hub.atomgit.com/amd64/alpine:3.15.10


# 创建配置文件目录并复制配置文件
RUN mkdir -p /root/configs
COPY configs /root/configs

# 安装必要的 ca-certificates
# RUN apk update && apk add --no-cache ca-certificates
# RUN apk add --no-cache ca-certificates

# 设置工作目录
WORKDIR /root/

# 从 builder 镜像复制构建的二进制文件
COPY --from=builder /app/main .
COPY configs/config.yaml /root/config.yaml

# 确保二进制文件可执行
RUN chmod +x ./main

# 暴露应用程序端口（可选，根据你的应用需求）
EXPOSE 8080

# 运行 Go 应用
CMD ["./main"]


```
ps:云服务器安装RUN apk update && apk add --no-cache ca-certificates这一步骤太慢可使用国内镜像,
云服务器是http，不需要这步骤
```go
# 更换镜像源为阿里云，加快下载速度
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache ca-certificates
```

docker hub下载镜像慢的时候，可以用hub.atomgit.com国内镜像替换,如hub.atomgit.com/amd64/golang:1.21
#### 4  将代码推送到远程服务器
使用 `scp` 来复制代码到远程服务器
```go
scp -r /path/to/your/local/project user@remote_server:/path/to/remote/directory
```
#### 5. 构建 Docker 镜像
在项目根目录下，运行以下命令构建 Docker 镜像：
```shell
docker build -t my-go-app:latest .
```
#### 6. 运行 Docker 容器
有时候云服务器docker连接不上数据库,需要创建一个自定义Docker网络,并将相关容器连接到该网络.

1. 创建一个自定义网络
```dockerfile
docker network create my-network
```
2 启动 MySQL 容器并连接到自定义网络
```dockerfile
docker run -d --name my-mysql --network my-network \
    -e MYSQL_ROOT_PASSWORD=123456 \
    -e MYSQL_DATABASE=blog_service \
    -p 3306:3306 \
    mysql:5.7
```
 3. 确认 MySQL 容器正在运行
```dockerfile
docker ps
```
确保 `my-mysql` 容器在自定义网络 `my-network` 上运行。

运行以下命令启动 Docker 容器：
```shell
docker run -d --name my-go-app --network my-network -p 8080:8080 my-go-app
```
#### 7 管理容器
验证服务是否正常运行
```dockerfile
docker logs my-mysql 
docker logs my-go-app
```
停止容器
```shell
docker stop my-go-app-container
```
启动已停止的容器：
```shell
docker start my-go-app-container
```
删除容器：
```shell
docker rm my-go-app-container
```
查看所有容器（包括停止的）：
```shell
docker ps -a
```
#### 8 验证部署
在浏览器中访问云服务器的 IP 地址和端口 `8080`，例如 `http://<your-cloud-server-ip>:8080`，你应该会看到 "Hello, World!"。
```dockerfile
curl http://你的服务器IP:8080
```
#### 9 配置反向代理
用nginx配置反向代理
### 
