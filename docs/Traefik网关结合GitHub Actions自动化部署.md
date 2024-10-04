以下是关于Traefik结合GitHub Actions进行自动化部署的Markdown文本,已屏蔽敏感信息:

# Traefik与GitHub Actions自动化部署

## Traefik配置

1. 创建Traefik项目目录结构:

```
traefik-project/
├── docker-compose.yml
└── config/
    ├── traefik.yml
    └── acme.json
```

2. `docker-compose.yml`内容:

```yaml
version: '3'

services:
  traefik:
    image: traefik:v2.5
    container_name: traefik
    restart: unless-stopped
    security_opt:
      - no-new-privileges:true
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./config/traefik.yml:/traefik.yml:ro
      - ./config/acme.json:/acme.json
    networks:
      - traefik-public

networks:
  traefik-public:
    external: true
```

3. `config/traefik.yml`内容:

```yaml
api:
  dashboard: true

entryPoints:
  http:
    address: ":80"
  https:
    address: ":443"

providers:
  docker:
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false

certificatesResolvers:
  letsencrypt:
    acme:
      email: your-email@example.com
      storage: acme.json
      httpChallenge:
        entryPoint: http

http:
  middlewares:
    https-redirect:
      redirectScheme:
        scheme: https
        permanent: true

  routers:
    http-catchall:
      rule: "hostregexp(`{host:.+}`)"
      entrypoints:
        - http
      middlewares:
        - https-redirect
```

4. `config/acme.json`文件:

这是一个用于存储Let's Encrypt证书的JSON文件。初始时为空,Traefik会自动填充。确保权限设置为600:

```bash
touch config/acme.json
chmod 600 config/acme.json
```

## GitHub Actions配置

1. 在你的GitHub仓库中创建`.github/workflows/deploy.yml`文件:

```yaml
name: Deploy to Server

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Install SSH Key
      uses: shimataro/ssh-key-action@v2
      with:
        key: ${{ secrets.SERVER_SSH_KEY }}
        known_hosts: ${{ secrets.KNOWN_HOSTS }}
    
    - name: Deploy to Server
      run: |
        ssh user@your-server-ip << EOF
          cd /path/to/your/project
          git pull origin main
          docker-compose up -d --build
        EOF
```

2. 在GitHub仓库设置中添加以下Secrets:

- `SERVER_SSH_KEY`: 你的服务器SSH私钥
- `KNOWN_HOSTS`: 服务器的known_hosts内容

## 部署流程

1. 开发者推送代码到GitHub main分支。
2. GitHub Actions自动触发部署工作流。
3. Actions使用SSH连接到服务器。
4. 在服务器上拉取最新代码并重新构建Docker容器。
5. Traefik自动检测新的容器并更新路由配置。

通过这种方式,你可以实现代码推送后的自动部署,同时利用Traefik进行反向代理和自动HTTPS配置。

请记得根据实际情况调整配置文件中的域名、邮箱等信息。