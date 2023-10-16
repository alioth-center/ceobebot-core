# CeobeBot QQChannel

## 使用前准备

1. 你需要在腾讯申请一个QQ机器人，获取appID, appToken
2. 你需要创建好频道，获取频道ID，用户(主人)ID
3. 你最好有一定的计算机基础，能够使用命令行

## 功能介绍

描述中标 `*` 的部分功能有额外说明，可以看对应目录下的 `readme.md` 文件

| 功能 |  描述  | 指令 | 描述 | 主要贡献者 | 目录 |
| :--: |:--:|:--:|:--:|:--:|:--:|
| base | 基本功能 | 无 | 基础功能 | [sunist-c](github.com/sunist-c)  |     /plugin     |
| example | 示例插件 | /hello | 一个发送问候的示例插件 | [sunist-c](github.com/sunist-c)  | /plugin/example |
| zerobot | zbp适配 | /zero [command] [args...] | 从zerobot-plugin迁移的部分功能 | [sunist-c](github.com/sunist-c) | /plugin/zerobot |
| mihoyo | 米哈游接口 | /mihoyo [command] [args...] | 使用enka.network接口的米哈游相关功能 | [sunist-c](github.com/sunist-c) | /plugin/mihoyo |
| ai | OpenAI接口 | /ai [command] [args...] | 使用OpenAI接口的对话/绘画功能 | [sunist-c](github.com/sunist-c) | /plugin/chat |

## 开发计划

在插件外，我们会对本机器人的基础功能进行拓展，开发路线如下：

- [x] 支持QQ频道，可以接受/发送QQ频道聊天室消息与私聊消息
- [ ] 支持 `go-cqhttp` / `shamrock` 等框架
- [ ] 通过本地 `rpc` 的方式支持如 `javascript` / `python` 等语言的插件
- [ ] 通过 `chrmoedp` / `puppeteer` 等依赖支持浏览器渲染图片

## 如何使用

我们准备了四种方式来运行这个机器人，你可以选择一种你喜欢的方式来运行

### a. 从源代码运行

1. 你需要一个golang开发环境，项目使用的是 `golang:1.21`，你可以在 `go.mod` 中找到最新版本的机器人所使用的 `golang` 版本
2. 你需要一个 `git` 客户端，用来下载源代码，或者你可以直接下载源代码的压缩包

    找到一个你喜欢的目录，用命令行下载源代码：
    
    ```bash
    git clone https://github.com/ceobebot/qqchannel.git
    ```
    
    或者你可以直接下载源代码的压缩包，解压到你喜欢的目录

    打开 `plugin/imports` 文件夹，在 `imports.go` 文件中，你可以打开或者关闭你需要的插件，然后保存

3. 在你准备好上面两步后，根据你的操作系统类型，在项目目录下，执行下面的指令：

    ```bash
    # 下载依赖
    go mod tidy
    
    # 编译
    ## 如果你是windows用户，执行下面的指令，windows只支持x86_64架构
    make build_windows
    
    ## 如果你是linux用户，执行下面的指令，linux只支持x86_64架构
    make build_linux
    
    ## 如果你是mac用户，执行下面的指令，mac支持x86_64和arm64架构，请根据你的机器选择
    ## arm64架构(使用 Apple M 芯片的计算机)的mac用户请使用下面的指令
    make build_mac_arm64
    
    ## x86_64架构(使用 Intel 芯片的计算机)的mac用户请使用下面的指令
    make build_mac_x86_64
    
    ## 如果你想全部都来一份，可以使用下面的指令
    make build_all
    ```

4. 编译完成后，你可以在 `build` 文件夹下找到编译好的机器人，你可以使用下面的指令来运行机器人
5. 这个时候你需要在 `data` 文件夹下准备好 `config.yaml` 文件，你可以在 `config.yaml.example` 中找到一个示例，你需要将其复制一份，并且修改其中的内容，需要注意的是部分插件需要你提供额外的配置，你可以在对应目录的 `config.yaml.example` 中找到示例
6. 你可以直接打开 `build` 文件夹下对应的应用程序，如果你是linux/mac用户，你可以直接 `./ceobebot-qqchanel_${对应架构}` 来运行机器人，如果你是windows用户，你可以直接双击运行
7. 享受机器人带来的乐趣吧

### b. 从二进制运行

如果你不想自己编译，你可以直接下载我们编译好的二进制文件，然后按照上面的 a.5. 以后的步骤来运行

### c. 从 docker 运行

1. 准备好 docker 与 docker-compose(如果你想使用 docker-compose 来运行)，你可以在 [docker](https://docs.docker.com/get-docker/) 与 [docker-compose](https://docs.docker.com/compose/install/) 中找到安装方法
2. 你需要一个 `git` 客户端，用来下载源代码，或者你可以直接下载源代码的压缩包

    找到一个你喜欢的目录，用命令行下载源代码：
    
    ```bash
    git clone https://github.com/ceobebot/qqchannel.git
    ```
    
    或者你可以直接下载源代码的压缩包，解压到你喜欢的目录
    
    打开 `plugin/imports` 文件夹，在 `imports.go` 文件中，你可以打开或者关闭你需要的插件，然后保存

3. 你需要在 `data` 文件夹下准备好 `config.yaml` 文件，你可以在 `config.yaml.example` 中找到一个示例，你需要将其复制一份，并且修改其中的内容，需要注意的是部分插件需要你提供额外的配置，你可以在对应目录的 `config.yaml.example` 中找到示例
4. 在项目目录下，执行下面的指令：

    ```bash
    # 使用 docker 运行
    docker build -t ceobebot-qqchannel .
    docekr run -d --name ceobebot-qqchannel \
        -v ./data:/app/data \
        ceobebot-qqchannel
    
    # 使用 docker-compose 运行
    docker-compose up -d
    ```

5. 享受机器人带来的乐趣吧

### d. 从构建好的 docker 镜像运行

1. 准备好 docker 与 docker-compose(如果你想使用 docker-compose 来运行)，你可以在 [docker](https://docs.docker.com/get-docker/) 与 [docker-compose](https://docs.docker.com/compose/install/) 中找到安装方法
2. 找一个你喜欢的目录，创建一个 `data` 目录，准备好 `config.yaml` 文件，你可以在 `config.yaml.example` 中找到一个示例，你需要将其复制一份，并且修改其中的内容，需要注意的是部分插件需要你提供额外的配置，你可以在对应目录的 `config.yaml.example` 中找到示例
3. 在 `data` 目录的父目录下，执行下面的指令：

    ```bash
    docker run -d --name ceobebot-qqchannel \
        -v ./data:/app/data \
        ghcr.io/ceobebot/qqchannel:latest
    ```

4. 享受机器人带来的乐趣吧
