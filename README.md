### sloth
延迟队列。

---------------------------------------
#### 使用
方式1（本地调试）:
```
go run main.go -c config_dev.ini
或
go build 
./sloth -c config_dev.ini
```
方式2（直接以docker运行方式调试）:
```
docker-compose up -d
```

#### 特性
* 支持基本的延迟功能
* 资源文件可配置化
* 支持优雅重启

#### Todo
* 支持消息重试机制
* 支持消费功能，目前消息消费逻辑仍然在hhzhome老队列
* 增加prometheus监控功能
