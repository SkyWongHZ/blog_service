#### docker查看卷
```plain
docker  volume   ls
```



#### docker删除无用卷
```plain
docker volume  prune
```

#### 列出当前正在使用的容器
```plain
docker ps
```



#### 创建一个命名卷
```plain
docker volume create <卷名>
```



#### 使用该卷启动 MySQL 容器(使用了自定义网格)
```dockerfile
docker run -d --name <new_mysql_container_name> --network <自定义网格名> -e MYSQL_ROOT_PASSWORD=your_new_password -e MYSQL_DATABASE=blog_service -v blog_service_mysql_data:/var/lib/mysql -p 3306:3306 --restart always mysql:latest
```

- `<new_mysql_container_name>`: 新的 MySQL 容器的名称
- `<your_root_password>`: MySQL root 用户的密码
- `<your_mysql_volume_name>`: 您的 MySQL 数据 volume 的名称
- `<host_port>`: 主机上映射到 MySQL 3306 端口的端口号

#### <font style="color:rgb(0, 0, 0);">docker下登录mysql </font>
```plain
docker exec -it <mysql模块名> mysql -u root -p
```

#### <font style="color:rgb(0, 0, 0);background-color:rgb(247, 247, 247);">docker打包镜像</font>
```plain
docker build -t <镜像名:标签>
```

<font style="color:rgb(0, 0, 0);"></font>

<font style="color:rgb(0, 0, 0);background-color:rgb(247, 247, 247);"></font>

<font style="color:rgb(0, 0, 0);background-color:rgb(247, 247, 247);"></font>

