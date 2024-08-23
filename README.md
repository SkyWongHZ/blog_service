# go_blog_service

### 用户模块
待解决问题:
* ~~更新接口有问题(用户user,文章article,标签Tag都有类似问题)~~
* 注册用户接口密码做加密措施
* user/:id以及查询参数user  get请求的用法
* 用户名去重逻辑((用户user,文章article,标签Tag都有类似问题))

### 整个系统
* ~~部署到云服务器,上nginx和docker~~
* 修改日志为zap
* 本地yaml配置和远程服务器配置支持
* 建自动化流水线  CI/CD(docker  compose或直接挂载)
* 写脚本deploy.sh支持自动化部署
* yaml支持远程和本地数据库配置
* API网关  
* 鉴权服务   

启动命令  
go run  main.go



 

