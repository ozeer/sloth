### sloth
延迟队列。

---------------------------------------
#### 使用
```
cp config.ini.sample config.ini
go build 
./sloth -c config.ini
```

#### 特性
* 支持基本的延迟功能
* 资源文件可配置化
* 支持优雅重启

#### Todo
* 支持消息重试机制
* 支持消费功能，目前消息消费逻辑仍然在hhzhome老队列
* 增加prometheus监控功能
