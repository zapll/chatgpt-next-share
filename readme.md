# ChatGPT-Next-Share

一个简单的 `ChatGPT` 共享程序, 

基于 [ninja](https://github.com/gngpp/ninja) 提供 `ChatGPT` 反代能力, 

并在此基础上提供 `账号管理`, `会话隔离` 等基础功能, 方便进行分享与用户管理.

电报群: [TG](https://t.me/+ZmZ49HVYwU0zYjg1)

## 前置

一个解锁 `ChatGPT` 的网络环境. 

Vps推荐 [racknerd](https://my.racknerd.com/aff.php?aff=10886) 洛杉矶 地区, 最低配即可.

## 部署

### docker部署

```shell
git clone git@github.com:zapll/chatgpt-next-share.git
docker compse up
```

服务启动后, 默认情况下 

- chatgpt服务: [http://127.0.0.1:3000](http://127.0.0.1:3000)

- 后台管理服务: [http://127.0.0.1:3000](http://127.0.0.1:3000)

如何使用

1. 管理后台添加 ChatGPT 账号

- 准备账号, 在 chat.openai.com 官网登录你的账号, 右键检查, 打开调试工具, 点击 `Application` 选项卡, 
  找到 Name 为 `__Secure-next-auth.session-token` 的 Cookie, 并复制他的 Value

- 登录后台, 默认账号: `nextshare`, 默认密码: `cns@0001`

- 账号管理菜单下新建, 贴如上一步中复制的 `session-token` 即可

2. 登录使用 ChatGPT

此时可以使用 `cns0001` 这个测试 token 进行登录

### 编译部署

开发与编译本项目的环境依赖: `bun >= 1.0.26`, `go >= 1.20`, `一个启动好的 ninja`

1.  项目根目录下设置环境变量

```shell
export CNS_DATA=$PWD
export CNS_NINJA=http://127.0.0.1:7999
```

2. 启动后台服务

```shell
cd admin
bun install
bun run dev
```

3. 启动后台服务

```shell
cd share
go run .
```

## 支持本项目

如果本项目对你有帮助, 请不吝赞赏一下.

![](./qrcode.png)