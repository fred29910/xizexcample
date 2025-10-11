# 架构说明文档

本文档详细描述了 `xizexcample` 项目的核心架构、组件交互和启动流程。

## 核心组件分析

以下是项目关键包的核心职责和功能分析：

- **`cmd` (main.go)**: 作为应用程序的入口点，`main.go` 负责初始化 Zinx 服务器实例，从 `internal/conf` 加载配置，设置服务器钩子，并最终启动服务。

- **`internal/conf`**: 此包负责加载和管理整个应用程序的配置。它通过 `init()` 函数自动从 `conf/zinx.json` 文件加载配置，并将其存储在全局可访问的 `AppConfig` 变量中。

- **`internal/router`**: 定义了消息路由规则。`InitRouter` 函数将不同的消息 ID 映射到相应的处理程序（Handler），例如 `JoinRoomHandler`、`PlaceBetHandler` 等。此外，它还负责设置连接创建和销毁时的生命周期钩子。

- **`internal/logic`**: 包含了游戏的核心业务逻辑和数据结构。
    - **`room.go`**: 定义了 `Room` 结构体，这是游戏的核心场景。它管理房间状态、玩家列表、牌局（Deck），并提供线程安全的方法来处理玩家加入/离开、发牌、设置庄家等操作。
    - **`player.go`**: 定义了 `Player` 结构体，用于表示一个玩家及其状态，如手牌、是否在线、是否为庄家等。
    - **`deck.go`**: 定义了 `Deck` 结构体，负责管理一副牌，包括初始化、洗牌和发牌等功能。
    - **`room_fsm.go`**: (推断) 可能包含游戏房间的状态机逻辑，用于管理游戏的不同阶段（如准备、下注、摊牌等）。

- **`internal/server`**: 包含与 Zinx 服务器生命周期相关的钩子函数。例如，`hooks.go` 中的 `OnConnStop` 函数定义了当客户端连接断开时的处理逻辑。

- **`api/proto`**: 此目录存放了用于客户端和服务端之间通信的 Protobuf (`.proto`) 文件。这些文件定义了消息的结构，并通过代码生成脚本（如 `scripts/gen_proto.sh`）生成相应的 Go 代码。

## 依赖关系与启动流程图

下图清晰地展示了各组件之间的依赖关系以及应用的完整启动流程。

```mermaid
graph TD
    subgraph "启动流程"
        A[main.go] --> B{加载配置 internal/conf};
        B --> C{创建 Zinx 服务器};
        C --> D{初始化路由 internal/router};
        D -- 注册 Handler --> E[消息处理程序];
        E -- 调用业务逻辑 --> F[internal/logic];
        F -- 操作 --> G[Room/Player/Deck];
        C --> H{设置连接钩子 internal/server};
        C --> I[启动服务器 s.Serve()];
    end

    subgraph "组件依赖关系"
        main.go -- 依赖 --> internal/conf;
        main.go -- 依赖 --> internal/server;
        main.go -- 依赖 --> internal/router;
        internal/router -- 依赖 --> internal/logic;
        internal/logic -- 依赖 --> internal/msg;
    end
