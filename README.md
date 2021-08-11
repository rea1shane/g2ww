# Grafana to Wechat-Work
该服务的转发功能基于企业微信群机器人，通过 webhook 将告警信息转发到企业微信。

## 环境
- go version go1.16.7 linux/amd64
- Grafana v8.1.0 (62e720c06b)
- CentOS Linux release 7.4.1708 (Core)

## 安装与运行
1. 下载 `.zip` 包；
2. 解压到目标文件夹（推荐 `gopath` 路径）；
3. 进入文件夹，运行 `go mod tidy` 下载所需依赖；
4. 运行 `go build -o g2ww` 编译为可执行二进制文件；
5. `./g2ww [-port 3001]` 运行二进制文件，可指定服务启动端口；

## 配置请求方式
1. 在 grafana 中创建 `Notification channel`，类型为 `webhook`；
2. 在 `Url` 里添写 `http://{host}:{port}/send?key={bot_key}`；
3. 点击 `Send Test`，能够正确看到监控信息后点击 `Save`；

## 消息类型
实现了发送企业微信的两种类型的消息：
- `markdown`（默认）
- `news`（通过 grafana webhook url 的 `type=news` 参数启用）

## 运行状态
通过 `GET http://{host}:{port}/` 接口可以获取 g2ww 服务的运行状态, 获取 `发送成功` 和 `发送失败` 的消息数。

## 相关文档
[企业微信群机器人配置说明](https://work.weixin.qq.com/api/doc/90000/90136/91770)