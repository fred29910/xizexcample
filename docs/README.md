# 斗牛游戏服务器

## 项目简介

这是一个基于 Go 语言开发的多人在线斗牛游戏服务器，采用 Zinx 网络框架实现高性能的 TCP 网络通信。项目实现了完整的斗牛游戏核心逻辑，包括房间管理、玩家管理、游戏状态机、牌型计算等功能模块。

## 主要特性

- **高性能网络通信**: 基于 Zinx 框架，支持高并发连接
- **完整的游戏逻辑**: 实现标准斗牛游戏规则，包括抢庄、下注、摊牌、结算
- **状态机管理**: 采用 FSM（有限状态机）管理游戏流程
- **房间系统**: 支持多房间并发，最多 5 人/房间
- **断线重连**: 支持玩家断线后重新连接
- **牌型计算**: 完整实现斗牛牌型计算，包括特殊牌型（五小牛、炸弹、金花等）
- **Protocol Buffers**: 使用 Protobuf 定义通信协议

## 技术栈

### 核心框架
- **Go**: 1.25.1
- **Zinx**: 1.2.7 - 高性能游戏服务器框架
- **Protocol Buffers**: v1.33.0 - 数据序列化

### 开发工具
- **testify**: v1.8.1 - 单元测试框架
- **Make**: 项目构建工具

## 项目结构

```
xizexcample/
├── api/                    # API 定义
│   └── proto/             # Protocol Buffers 定义文件
├── bin/                   # 编译输出目录
├── conf/                  # 配置文件
│   └── zinx.json         # Zinx 服务器配置
├── docs/                  # 文档目录
├── internal/              # 内部实现代码
│   ├── conf/             # 配置加载模块
│   ├── logic/            # 游戏核心逻辑
│   ├── msg/              # 消息定义（由 proto 生成）
│   ├── pkg/              # 通用工具包
│   ├── router/           # 消息路由处理
│   └── server/           # 服务器管理模块
├── scripts/               # 脚本工具
│   └── gen_proto.sh      # Protocol Buffers 生成脚本
├── tests/                 # 测试代码
│   ├── e2e/              # 端到端测试
│   └── unit/             # 单元测试
├── go.mod                 # Go 模块定义
├── main.go                # 程序入口
└── Makefile              # 构建脚本
```

## 快速开始

### 环境要求

- Go 1.25+ 
- Protocol Buffers 编译器（用于生成代码）
- Make（可选，用于便捷构建）

### 安装依赖

```bash
go mod download
```

### 编译项目

```bash
# 使用 Makefile
make build

# 或直接使用 go build
go build -o bin/xizexcample .
```

### 运行服务器

```bash
# 使用 Makefile
make run

# 或直接运行编译后的程序
./bin/xizexcample
```

服务器将在 `0.0.0.0:8999` 上启动（可在 `conf/zinx.json` 中修改配置）。

### 运行测试

```bash
# 运行所有测试
make test

# 运行特定测试
go test -v ./internal/logic/...

# 运行 E2E 测试
go test -v ./tests/e2e/...
```

## 配置说明

配置文件位于 `conf/zinx.json`：

```json
{
  "Name": "ZinxServerApp",
  "Host": "0.0.0.0",
  "TCPPort": 8999,
  "MaxConn": 12000,
  "WorkerPoolSize": 10
}
```

- **Name**: 服务器名称
- **Host**: 监听地址
- **TCPPort**: TCP 端口号
- **MaxConn**: 最大连接数
- **WorkerPoolSize**: 工作协程池大小

## 游戏规则

斗牛游戏基本规则：

1. **玩家数量**: 2-5 人
2. **发牌**: 每位玩家发 5 张牌
3. **抢庄**: 玩家竞争成为庄家
4. **下注**: 闲家根据手牌进行下注
5. **摊牌**: 所有玩家亮牌比较大小
6. **结算**: 计算输赢并更新分数

### 牌型大小

从大到小：
- 金花（同花顺）
- 炸弹（四张相同）
- 五小牛（五张点数之和 ≤ 10）
- 牛牛（10 的倍数）
- 牛九 ~ 牛一
- 无牛

## 文档导航

- [架构设计](./architecture.md) - 系统架构和设计理念
- [功能模块](./features.md) - 详细功能说明
- [通信协议](./protocol.md) - 客户端/服务端通信协议
- [开发指南](./development.md) - 开发、测试、部署指南

## License

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
