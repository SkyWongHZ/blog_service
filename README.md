# go_blog_service

### 业务问题

- [x] 更新接口有问题(用户 user,文章 article,标签 Tag 都有类似问题)
- [ ] 注册用户接口密码做加密措施
- [x] 用户名去重逻辑((用户 user,文章 article,标签 Tag 都有类似问题))
- [x] 接口返回有问题 做成 data code msg 形式
- [ ] 评论模块
- [x] 添加 minio 后文章模块的业务逻辑修改,使用阿里云 oss
- [] 文章模块逻辑修复,搭配上 minio 后,redis 和其他逻辑修复

### 技术栈

- [x] 部署到云服务器,上 nginx 和 docker
- [x] 修改日志为 zap
- [x] 写脚本 deploy.sh 支持自动化部署
- [x] yaml 支持远程和本地数据库配置
- [ ] API 网关
- [x] 鉴权服务
- [ ] grpc
- [ ] websocket
- [x] github actinon 自动化部署
- [x] 缓存 redis
- [ ] 消息队列：可以先用 Redis，之后可以考虑 RabbitMQ 或 Kafka
- [ ] 全文搜索：Elasticsearch（如果您需要高级搜索功能）
- [ ] 日志管理：ELK stack（Elasticsearch, Logstash, Kibana）或 Prometheus + Grafana
- [x] 对象存储：MinIO 或云服务（如阿里云 OSS、AWS S3）用于存储用户上传的图片等文件
- [ ] 权限控制：Casbin
- [ ] 定时任务
- [x] docker-compose 使用
- [ ] Oauth2.0

swagger 地址
http://118.178.184.13:8080/swagger/index.html
