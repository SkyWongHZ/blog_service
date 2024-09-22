# 使用 GitHub Actions 实现自动化部署到阿里云服务器

本文将介绍如何使用 GitHub Actions 来实现 Go 项目的自动化构建和部署到阿里云服务器。这种方法可以大大简化部署流程，提高开发效率。

## 1. 创建工作流文件

首先，在 GitHub 仓库中创建 `.github/workflows/deploy.yml` 文件。这个文件定义了自动化部署的工作流程。

```yaml
name: Deploy to Aliyun Server

on:
  workflow_dispatch:
    inputs:
      deployMessage:
        description: 'Deployment message'
        required: true
        default: 'Manual deployment triggered'

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Build
      run: go build -v ./...

    - name: Test
      run: go test -v ./...

    - name: Login to Aliyun Container Registry
      uses: docker/login-action@v2
      with:
        registry: <your-registry-url>
        username: ${{ secrets.ALIYUN_CR_USERNAME }}
        password: ${{ secrets.ALIYUN_CR_PASSWORD }}
    
    - name: Generate config file
      run: |
        envsubst < configs/config.yaml.template > configs/config.yaml
      env:
        DB_USERNAME: ${{ secrets.DB_USERNAME }}
        DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
        DB_HOST: ${{ secrets.DB_HOST }}
        DB_NAME: ${{ secrets.DB_NAME }}
        JWT_SECRET: ${{ secrets.JWT_SECRET }}
    
    - name: Build and push Docker image
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: <your-registry-url>/<your-project-name>/blog-service:${{ github.sha }},<your-registry-url>/<your-project-name>/blog-service:latest

    - name: Deploy to Aliyun Server
      uses: appleboy/ssh-action@master
      with:
        host: ${{ secrets.ALIYUN_HOST }}
        username: ${{ secrets.ALIYUN_USERNAME }}
        key: ${{ secrets.ALIYUN_SERVER_SSH_KEY }}
        script: |
          docker login --username=${{ secrets.ALIYUN_CR_USERNAME }} --password=${{ secrets.ALIYUN_CR_PASSWORD }} <your-registry-url>
          docker pull <your-registry-url>/<your-project-name>/blog-service:latest
          docker stop blog-service || true
          docker rm blog-service || true
          docker run -d --network my-network --name blog-service -p 8080:8080 --restart always <your-registry-url>/<your-project-name>/blog-service:latest
```

## 2. 配置工作流

工作流程包括以下主要步骤：

1. 检出代码
2. 设置 Go 环境
3. 构建和测试 Go 项目
4. 登录到阿里云容器镜像服务
5. 生成配置文件
6. 构建 Docker 镜像并推送到阿里云容器镜像仓库
7. 通过 SSH 连接到阿里云服务器并部署

为了保护敏感信息，我们使用了 `config.yaml.template` 文件和环境变量。确保在 Dockerfile 中正确复制生成的配置文件：

```dockerfile
COPY configs/config.yaml /root/configs/config.yaml
```

## 3. 设置密钥和环境变量

在 GitHub 仓库的 Settings -> Secrets and variables -> Actions 中添加以下密钥：

- `ALIYUN_HOST`: 阿里云 ECS 实例的公网 IP
- `ALIYUN_USERNAME`: SSH 用户名
- `ALIYUN_SERVER_SSH_KEY`: 用于 SSH 连接到阿里云服务器的私钥
- `ALIYUN_CR_USERNAME`: 阿里云容器镜像服务用户名
- `ALIYUN_CR_PASSWORD`: 阿里云容器镜像服务密码
- `DB_USERNAME`, `DB_PASSWORD`, `DB_HOST`, `DB_NAME`, `JWT_SECRET`: 数据库和 JWT 相关的配置

## 4. 配置 Dockerfile

在项目根目录创建 `Dockerfile`：

```dockerfile
# 构建阶段
FROM golang:1.21 as builder

WORKDIR /app

# 安装 Git（如果需要）
RUN apt-get update && apt-get install -y git

# 设置 Go 代理
ENV GOPROXY=https://goproxy.cn,direct

# 复制并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 运行阶段
FROM alpine:3.15

WORKDIR /root/

# 从构建阶段复制二进制文件和配置
COPY --from=builder /app/main .
COPY configs/config.yaml /root/configs/config.yaml

# 设置可执行权限
RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]
```

## 5. 触发部署

可以通过以下方式触发部署：

- 推送到特定分支（例如 `main`）
- 在 GitHub Actions 页面手动触发工作流

## 6. 验证部署

部署完成后，可以通过以下方式验证：

1. SSH 连接到阿里云服务器，检查 Docker 容器状态：

   ```bash
   docker ps
   ```

2. 访问应用程序的公网地址和端口，确认服务是否正常运行。


通过实施这些建议，可以建立一个更加健壮和可靠的自动化部署流程，提高开发效率和系统稳定性。

注意：在实际使用时，请将 `<your-registry-url>` 和 `<your-project-name>` 替换为自己项目的实际值。