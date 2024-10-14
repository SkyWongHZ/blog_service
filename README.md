# go_blog_service

### 业务问题
- [ ] 注册用户接口密码做加密措施
- [ ] 评论模块
- [ ] 点赞模块
- [ ] cobar 命令行
- [ ] nullString 




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
