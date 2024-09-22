# go_blog_service

### 业务问题
- [x] 更新接口有问题(用户user,文章article,标签Tag都有类似问题)
- [ ] 注册用户接口密码做加密措施
- [x] 用户名去重逻辑((用户user,文章article,标签Tag都有类似问题))
- [x] 接口返回有问题  做成data  code   msg形式

### 技术栈
- [x] 部署到云服务器,上nginx和docker
- [x] 修改日志为zap
- [x] 写脚本deploy.sh支持自动化部署
- [x] yaml支持远程和本地数据库配置
- [ ] API网关  
- [x] 鉴权服务 
- [ ] grpc 
- [ ] websocket
- [x] github  actinon自动化部署
- [ ] 缓存redis
- [ ] 消息队列：可以先用 Redis，之后可以考虑 RabbitMQ 或 Kafka
- [ ] 全文搜索：Elasticsearch（如果您需要高级搜索功能）
- [ ] 日志管理：ELK stack（Elasticsearch, Logstash, Kibana）或 Prometheus + Grafana
- [ ] 对象存储：MinIO 或云服务（如阿里云 OSS、AWS S3）用于存储用户上传的图片等文件
- [ ] 权限控制：Casbin
- [ ] 定时任务

swagger地址
http://118.178.184.13:8080/swagger/index.html








 

