# ChatGPT-Next-Share

一个简单的 `ChatGPT` 共享程序, 

基于 [ninja](https://github.com/gngpp/ninja) 提供 `ChatGPT` 反代能力, 

并在此基础上提供 `账号管理`, `会话隔离` 等基础功能, 方便进行分享与用户管理.

电报群: [TG](https://t.me/+ZmZ49HVYwU0zYjg1)

演示站: token `cns0001` 登录 [https://gpt.daodao.run](https://gpt.daodao.run)

## 前置

一个解锁 `ChatGPT` 的网络环境. 

Vps推荐 [racknerd](https://my.racknerd.com/aff.php?aff=10886) 洛杉矶 地区, 最低配即可.

## 部署

### docker部署

```shell
git clone https://github.com/zapll/chatgpt-next-share.git
docker compose up
```

`docker-compose.yml` 的简介绍

```yml
version: '3'

services:
  chatgpt-next-share:
    image: ghcr.io/zapll/chatgpt-next-share:latest
    container_name: chatgpt-next-share
    restart: unless-stopped
    volumes:
      - ./data:/data  # 挂载数据目录, 如果你不是直接克隆的本仓库, 那么记得把仓库中 data/db.sqlite 复制到你的目录中
    ports:
      - "3001:3001" # 导出后台服务端口
      - "3000:3000" # 导出代理服务端口
    environment:
      - CNS_NINJA=http://ninja:7999  # ninja 服务地址, 任意能连接到的地址即可, 也就是下发的 ninja 服务是非必须的
      - CNS_DATA=/data  # 存放数据的目录, 须根上方挂载的数据目录相同, 确保该目录下有 db.sqlite 文件
    depends_on:
      - ninja
  ninja:
    image: gngpp/ninja:latest
    container_name: ninja
    restart: unless-stopped
    environment:
      - TZ=Asia/Shanghai
    command: run --enable-webui --arkose-endpoint http://127.0.0.1:3000
    # ninja 服务必须启动ui --enable-webui
    # --arkose-endpoint 参数可以替换为你的实际域名, 否则无法使用 gpt4/gpts 等 

```

服务启动后, 默认情况下 

- chatgpt服务: [http://127.0.0.1:3000](http://127.0.0.1:3000)

- 后台管理服务: [http://127.0.0.1:3001](http://127.0.0.1:3001)

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