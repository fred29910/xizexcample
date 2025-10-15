# 系统架构设计

## 概述

本项目采用分层架构设计，将网络通信、业务逻辑、数据管理等模块清晰分离，便于维护和扩展。整体架构遵循 MVC 模式的变体，适合游戏服务器的开发需求。

## 总体架构

```
┌─────────────────────────────────────────────────────────┐
│                      Client Layer                        │
│                    (TCP Clients)                         │
└─────────────────────┬───────────────────────────────────┘
                      │ TCP Connection
                      │ Protocol Buffers
┌─────────────────────▼───────────────────────────────────┐
│                   Network Layer                          │
│                  (Zinx Framework)                        │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────┐ │
│  │ Connection  │  │   Message    │  │  Worker Pool   │ │
│  │  Manager    │  │    Router    │  │                │ │
│  └─────────────┘  └──────────────┘  └────────────────┘ │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                   Router Layer                           │
│              (internal/router/)                          │
│  ┌──────────┐ ┌───────────┐ ┌──────────┐ ┌──────────┐  │
│  │JoinRoom  │ │BidBanker  │ │PlaceBet  │ │Showdown  │  │
│  │Handler   │ │Handler    │ │Handler   │ │Handler   │  │
│  └──────────┘ └───────────┘ └──────────┘ └──────────┘  │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                  Service Layer                           │
│              (internal/server/)                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │            RoomManager (Singleton)                │   │
│  │  - Room Management                                │   │
│  │  - Player-Room Mapping                            │   │
│  │  - Lifecycle Hooks                                │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                  Logic Layer                             │
│              (internal/logic/)                           │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│  │  Room    │ │ Player   │ │   Deck   │ │Bull Logic│   │
│  │          │ │          │ │          │ │          │   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘   │
│  ┌────────────────────────────────────────────────┐     │
│  │              RoomFSM (State Machine)           │     │
│  └────────────────────────────────────────────────┘     │
└──────────────────────────────────────────────────────────┘
```

## 分层设计

### 1. Network Layer（网络层）

**职责**: 处理 TCP 连接、消息编解码、连接管理

**核心组件**:
- **Zinx Server**: 基于 Zinx 框架的 TCP 服务器
- **Connection Manager**: 管理所有客户端连接
- **Message Router**: 将消息路由到对应的处理器
- **Worker Pool**: 工作协程池，处理并发请求

**技术选型**:
- Zinx v1.2.7 - 高性能游戏服务器框架
- 支持 12000+ 并发连接
- 10 个工作协程处理消息

### 2. Router Layer（路由层）

**位置**: `internal/router/`

**职责**: 接收网络层转发的消息，解析请求，调用业务逻辑

**主要处理器**:

| 处理器 | 消息ID | 功能 |
|--------|--------|------|
| JoinRoomHandler | 101 | 处理玩家加入房间请求 |
| PlayerReadyHandler | 102 | 处理玩家准备状态 |
| BidBankerHandler | 103 | 处理抢庄请求 |
| PlaceBetHandler | 104 | 处理下注请求 |
| ShowdownHandler | 105 | 处理摊牌请求 |
| LeaveRoomHandler | 106 | 处理离开房间请求 |

**设计模式**:
- Handler Pattern: 每个消息类型对应一个独立的处理器
- Base Router: 所有处理器继承自 BaseRouter，提供公共功能

### 3. Service Layer（服务层）

**位置**: `internal/server/`

**职责**: 管理全局资源、协调业务逻辑

**核心组件**:

#### RoomManager（单例模式）
- 管理所有游戏房间的创建、销毁
- 维护玩家与房间的映射关系
- 提供房间查询接口
- 处理连接断开钩子

**数据结构**:
```go
type RoomManager struct {
    rooms      map[int32]*logic.Room  // 房间映射表
    playerRoom map[int64]int32        // 玩家-房间映射
    mu         sync.RWMutex           // 读写锁
}
```

### 4. Logic Layer（逻辑层）

**位置**: `internal/logic/`

**职责**: 实现游戏核心逻辑，不依赖网络层

#### 主要组件

##### Room（房间）
- 管理房间内的玩家列表（最多 5 人）
- 维护牌堆（Deck）
- 关联状态机（FSM）
- 发牌、设置庄家等操作
- 定时清理断线玩家（5 分钟超时）

**数据结构**:
```go
type Room struct {
    ID      int32
    Players map[int64]*Player
    Deck    *Deck
    FSM     *RoomFSM
    mu      sync.RWMutex
}
```

##### Player（玩家）
- 存储玩家基本信息（ID、昵称、分数）
- 管理手牌
- 跟踪游戏状态（等待、准备、游戏中、离线）
- 维护连接对象
- 支持断线重连

**数据结构**:
```go
type Player struct {
    ID             int64
    Nickname       string
    Score          int64
    RoomID         int32
    Hand           []Card
    BetAmount      int32
    Status         PlayerStatus
    isBanker       bool
    isOnline       bool
    Conn           ziface.IConnection
    DisconnectTime int64
    mu             sync.RWMutex
}
```

##### RoomFSM（房间状态机）
实现游戏流程控制，确保状态转换的合法性。

**状态定义**:
```go
const (
    STATE_WAITING_FOR_PLAYERS  // 等待玩家
    STATE_DEALING              // 发牌
    STATE_BIDDING              // 抢庄
    STATE_BETTING              // 下注
    STATE_SHOWDOWN             // 摊牌
    STATE_SETTLEMENT           // 结算
)
```

**状态转换图**:
```
WAITING_FOR_PLAYERS
    ↓ (2+ 玩家准备)
DEALING
    ↓ (发牌完成)
BIDDING
    ↓ (选出庄家)
BETTING
    ↓ (所有闲家下注完成)
SHOWDOWN
    ↓ (所有玩家摊牌)
SETTLEMENT
    ↓ (结算完成)
WAITING_FOR_PLAYERS (循环)
```

##### Deck（牌堆）
- 标准 52 张扑克牌
- 支持洗牌、发牌操作
- 线程不安全（由 Room 层加锁保护）

##### BullLogic（牛牛逻辑）
- 牌型计算：识别牛一到牛牛、特殊牌型
- 特殊牌型：五小牛、炸弹、金花（同花顺）
- 手牌比较：按牌型 → 牛值 → 最大单张 → 花色

**牌型优先级**:
```
金花 > 炸弹 > 五小牛 > 牛牛 > 牛九 > ... > 牛一 > 无牛
```

## 数据流

### 1. 玩家加入房间流程

```
Client → TCP → Router → RoomManager → Room
                ↓
              Player
                ↓
         Broadcast to all players
```

1. 客户端发送 `C2S_JoinRoomReq`
2. `JoinRoomHandler` 解析请求
3. `RoomManager` 获取或创建房间
4. 创建 `Player` 对象并加入房间
5. 发送 `S2C_JoinRoomAck` 给客户端
6. 广播房间状态给所有玩家

### 2. 游戏流程

```
WAITING → 玩家准备 → DEALING → 发牌
    ↓
BIDDING → 抢庄 → BETTING → 下注
    ↓
SHOWDOWN → 摊牌 → SETTLEMENT → 结算
    ↓
WAITING (下一局)
```

### 3. 断线重连流程

```
Client disconnect → OnConnStop Hook
    ↓
Set Player offline (keep in room)
    ↓
Start 5-minute timeout timer
    ↓
Client reconnect → JoinRoomHandler
    ↓
Detect existing player → Re-associate connection
    ↓
Send full room state sync
```

## 并发控制

### 1. 锁策略

- **RoomManager**: 使用 `sync.RWMutex` 保护全局房间映射
- **Room**: 使用 `sync.RWMutex` 保护房间内的玩家列表和状态
- **Player**: 使用 `sync.RWMutex` 保护玩家的游戏状态

### 2. 读写分离

- 查询操作使用 `RLock()`/`RUnlock()`
- 修改操作使用 `Lock()`/`Unlock()`
- 避免在锁内执行耗时操作

### 3. 协程安全

- 状态机转换串行化执行
- 消息处理通过 Worker Pool 并发执行
- 房间清理定时器独立协程运行

## 扩展性设计

### 1. 水平扩展

当前单机架构可扩展为分布式：
- 引入 Redis 存储房间和玩家数据
- 使用消息队列（如 NATS）实现跨节点通信
- 添加网关层进行负载均衡

### 2. 模块化

各层之间通过接口隔离：
- Router 层只依赖 Service 层接口
- Logic 层完全独立，可单独测试
- 便于替换网络框架或扩展游戏逻辑

### 3. 配置化

- 服务器参数通过配置文件管理
- 游戏规则参数可配置化（如房间人数上限、超时时间等）

## 性能考虑

### 1. 内存管理

- 对象池复用频繁创建的对象（如 Message）
- 及时清理断线玩家，避免内存泄漏
- 使用 `sync.Pool` 优化临时对象分配

### 2. 网络优化

- Protocol Buffers 高效序列化
- 批量消息广播优化
- 连接池管理

### 3. 计算优化

- 牌型计算结果缓存
- 避免重复计算
- 使用位运算优化牌型判断

## 安全性

### 1. 输入验证

- 所有客户端输入进行合法性检查
- 状态机确保游戏流程合法性
- 防止作弊和恶意请求

### 2. 连接管理

- 连接超时检测
- 异常连接断开处理
- 限制单 IP 连接数（可扩展）

### 3. 数据一致性

- 使用锁保证数据一致性
- 状态转换原子化操作
- 错误恢复机制

## 监控与日志

### 1. 日志系统

- **InfoLogger**: 记录正常业务流程
- **ErrorLogger**: 记录错误和异常
- 日志包含关键上下文（房间 ID、玩家 ID 等）

### 2. 监控指标（可扩展）

- 在线玩家数
- 活跃房间数
- 消息处理延迟
- 错误率统计

## 技术债务

当前实现中的 TODO 和改进点：

1. **广播机制**: 部分广播通知尚未实现
2. **结算逻辑**: 需要完善分数计算和更新
3. **玩家认证**: 缺少真实的玩家身份验证
4. **持久化**: 玩家数据未持久化存储
5. **优雅关闭**: 服务器关闭时需要通知玩家
6. **房间匹配**: 缺少自动匹配机制
7. **观战功能**: 未实现观战者模式

## 测试策略

### 1. 单元测试

- Logic 层组件（Room、Player、Deck、BullLogic）
- 状态机转换逻辑
- 牌型计算准确性

### 2. 集成测试

- RoomManager 与 Room 交互
- 完整游戏流程测试

### 3. E2E 测试

- 模拟多客户端连接
- 测试完整游戏流程
- 断线重连场景测试

## 部署架构

```
┌───────────────┐
│  Load Balance │
│   (Nginx)     │
└───────┬───────┘
        │
   ┌────┴────┐
   │         │
┌──▼───┐ ┌──▼───┐
│ App  │ │ App  │  (可水平扩展)
│ Node1│ │ Node2│
└──┬───┘ └──┬───┘
   │         │
   └────┬────┘
        │
   ┌────▼────┐
   │  Redis  │  (未来扩展)
   │ Cluster │
   └─────────┘
```

当前为单机部署，未来可扩展为分布式集群。
