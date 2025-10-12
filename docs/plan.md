# 项目优化与交付计划（后端）

更新时间：2025-10-11 22:22:24 +08:00

## 一、基于现有实现的分析摘要

- **核心玩法循环**
  - 状态机定义于 `internal/logic/room_fsm.go`，阶段依次为：`STATE_WAITING_FOR_PLAYERS` → `STATE_DEALING` → `STATE_BIDDING` → `STATE_BETTING` → `STATE_SHOWDOWN` → `STATE_SETTLEMENT` → 回到等待。
  - 发牌与状态推进由 `RoomFSM.DealCards()`、`RoomFSM.BidBanker()`、`RoomFSM.PlaceBet()`、`RoomFSM.Showdown()`、`RoomFSM.Settlement()` 驱动（部分 TODO 未实现）。
  - 牌型与比较逻辑在 `internal/logic/bull_logic.go`，`CalculateBull()`、`CompareHands()` 提供核心胜负判定。
  - 扑克与牌堆在 `internal/logic/deck.go`，房间管理在 `internal/logic/room.go`（玩家增删、离线、发牌、庄家设置、定时清理断线）。
- **协议与消息**
  - `api/proto/game.proto` 定义了 C2S/S2C 消息和 `GameState`、`PlayerStatus`、`CardPattern` 等枚举，产物 `internal/msg/game.pb.go`。
  - 关键消息流（示例）：C2S_JOIN_ROOM_REQ → S2C_JOIN_ROOM_ACK / S2C_SYNC_ROOM_STATE_NTF → S2C_GAME_START_NTF → S2C_DEAL_CARDS_NTF → S2C_BID_BANKER_NTF → S2C_BET_NTF → S2C_SHOWDOWN_NTF → S2C_GAME_RESULT_NTF。
- **用户操作路径（当前）**
  - 连接 TCP（Zinx），`main.go` 中通过 `znet.NewServer()` 启动，`internal/router` 注册路由（实现未在仓库，但在测试中引用）。
  - 玩家 JoinRoom → Ready（未见处理）→ 游戏开始（FSM）→ 发牌（FSM）→ 抢庄/下注/摊牌/结算（FSM）。
- **系统交互与数据流**
  - 当前状态与房间数据仅在内存结构 `Room/Player/Deck` 中维护，尚未持久化。断线重连通过 `tests/e2e/reconnect_test.go` 的模拟请求覆盖基础路径：`room.SetPlayerOffline()` → 再次 join → 恢复在线。
- **异常与容错（现状）**
  - `room_fsm.go` 对非法状态转换返回 error 并记录日志。
  - `room.go` 中对加人、发牌、移除等操作有错误返回；离线清理基于定时器（5 分钟）。
  - 缺少消息级别的入参校验、超时控制、回放/补偿机制、幂等性与重放防护。
- **潜在性能瓶颈**
  - 房间与玩家方法内部采用粗粒度锁；当玩家数、房间数增大时可能产生竞争。
  - 广播未实现，未来广播若采用逐连接同步写，易阻塞房间循环。
  - 所有状态推进在同一 goroutine/锁上下文中可能成为热点；消息编解码/校验无缓存；无对象池。

## 二、优化建议

- **体验**
  - 在 `game.proto` 增加明确的失败码和可读错误消息字段；所有 ACK 都需带 ret_code/ret_msg。
  - 增加 Ready、超时倒计时广播（已有 NTF 字段，可完善逻辑）。
  - 断线重连后的房间快照补发：`S2C_SyncRoomStateNtf` 应包含完整玩家状态/牌局阶段。
- **稳定性**
  - 引入“每房间一个事件循环（actor）”模型，所有房间内写操作通过 channel 串行处理，替代细粒度锁，消除竞态。
  - 为阶段操作引入 server 端超时（抢庄/下注/摊牌），到时自动推进。
  - 消息幂等与去重：基于 `player_id` + `client_msg_seq` 实现去重表（Redis/本地 LRU）。
- **安全性**
  - 接入认证（JWT/Session），所有 C2S 请求需携带 `auth_token`，服务器绑定 `player_id` 与连接。
  - 重要操作（下注、支付等）加入签名与服务端限流（IP、账号维度）。
  - 参数严格校验（下注范围、阶段校验、余额校验、牌面校验）。
- **可扩展性**
  - 将实时对局（Zinx TCP）与账户/物品/支付/好友等后台服务拆分：REST/gRPC 微服务化（后述设计）。
  - 房间分片：按 `room_id % N` 路由到不同进程/节点；跨节点通过消息总线（NATS/Kafka）同步事件或采用一致性哈希网关。

## 三、后端详细技术设计

### 3.1 服务与模块划分

- **Game Gateway（实时网关，现有 Zinx）**
  - 职责：TCP 连接、协议收发（MsgID + protobuf）、路由到房间 actor。
  - 依赖：Auth Service（校验 token）、Room Service（房间定位）、Match Service（可选）。
- **Room Service（对局/房间进程）**
  - 职责：房间生命周期、FSM、计时器、结算、广播。
  - 模式：每房间一个 goroutine + 事件 channel；无锁或仅连接字典读锁。
- **User Service（账户）**
  - REST/gRPC：注册、登录、资料、会话、风控。
- **Inventory/Trading Service（物品与交易）**
  - REST/gRPC：持久化道具、订单、交易撮合（若有）。
- **Payment Service（支付）**
  - 与支付渠道对接，异步回调，账本与对账。
- **Friend Service（好友/社交）**
  - 好友关系、申请、黑名单。
- **Shared Components**
  - `internal/pkg/logger`、`internal/pkg/errors`、`internal/pkg/auth`、`internal/pkg/codec`、`internal/pkg/rate`。

### 3.2 API 接口定义（示例）

- **REST（User/Inventory/Payment/Friend）**
  - POST /api/v1/users/register {username, password}
  - POST /api/v1/users/login {username, password} → {token}
  - GET /api/v1/users/me (Bearer) → {id, nickname, balance}
  - POST /api/v1/payments/orders {amount, channel} → {order_id}
  - POST /api/v1/friends/requests {to_user_id}
  - POST /api/v1/friends/accept {request_id}
  - ...
- **gRPC（Room/Match，可选）**
  - RoomAdminService: CreateRoom(room_id), CloseRoom(room_id), KickPlayer(room_id, player_id)
  - MatchService: Enqueue(user_id, mode) → room_id
- **Game TCP（现有）** MsgID 与 `game.proto`
  - C2S_JOIN_ROOM_REQ {room_id}
  - C2S_PLAYER_READY_REQ {is_ready}
  - C2S_BID_BANKER_REQ {multiple}
  - C2S_PLACE_BET_REQ {multiple}
  - C2S_SHOWDOWN_REQ {sorted_hand}
  - C2S_LEAVE_ROOM_REQ {}
  - 统一 ACK：{ret_code, ret_msg?, data?}

### 3.3 数据模型与存储

- **PostgreSQL（关系数据）**
  - users(id PK, username unique, password_hash, nickname, created_at)
  - sessions(id PK, user_id FK, jwt_id unique, expires_at)
  - friends(id PK, user_id, friend_user_id, status, created_at)
  - inventories(id PK, user_id, item_id, quantity, updated_at)
  - trades(id PK, buyer_id, seller_id, item_id, price, status, created_at)
  - payments(id PK, user_id, order_id unique, amount, channel, status, created_at)
  - ledgers(id PK, user_id, delta, reason, ref_id, created_at)
  - rooms(id PK, shard_key, status, created_at)
  - room_players(id PK, room_id, user_id, score, is_banker, last_active_at)
- **Redis（高速数据/会话/限流）**
  - session:{jwt_id} → user_id (TTL=token_exp)
  - room:{room_id}:snapshot → RoomInfo JSON（断线重连）
  - room:{room_id}:events → 最近 N 个事件（可选）
  - idempotency:{user_id}:{seq} → 1 (TTL=5m)
  - rate:{ip} / rate:{user_id}（漏斗或令牌桶）

### 3.4 关键业务逻辑与算法（伪代码/流程）

- **房间 Actor 主循环**
```go
for evt := range room.ch {
  switch evt.Type {
  case Join:
    if len(room.players) >= MaxPlayers { reply(ErrRoomFull); continue }
    addPlayer(evt.Player); snapshotDirty = true
  case Ready:
    markReady(evt.Player)
    if allReady(room) { fsm.StartGame(); schedule(DealTimeout) }
  case DealTimeout:
    fsm.DealCards(); broadcast(DealCardsNtf)
    schedule(BidTimeout)
  case Bid:
    collectBid(evt.Player, evt.Multiple)
    if allBidOrTimeout() { resolveBanker(); fsm.BidBanker(); schedule(BetTimeout) }
  case Bet:
    validateBet(evt.Player, evt.Multiple)
    if allBetOrTimeout() { fsm.PlaceBet(); schedule(ShowdownTimeout) }
  case Showdown:
    lockHand(evt.Player, evt.Sorted)
    if allShowOrTimeout() { settle(); fsm.Showdown(); fsm.Settlement(); resetRound() }
  }
  if snapshotDirty { saveSnapshot(room); snapshotDirty=false }
}
```

- **抢庄决策**
```go
resolveBanker() {
  maxM := max(multiples)
  cand := players with multiple==maxM
  if len(cand) == 1 { banker = cand[0] } else { banker = random(cand) }
}
```

- **结算（使用 `CompareHands`）**
```go
settle() {
  foreach p != banker {
    result := CompareHands(p.Hand, banker.Hand)
    mult := betMultiple(p) * bankerMultiple()
    delta := 0
    if result > 0 { delta = +base * mult } else if result < 0 { delta = -base * mult }
    p.Score += delta; banker.Score -= delta
    writeLedger(p.ID, delta, RoundID)
  }
  broadcast(GameResultNtf)
}
```

### 3.5 并发处理策略

- 每房间一个 goroutine + channel（串行化房间内写操作）。
- 网络读写：
  - 读线程将 C2S 解析为事件投递到房间 channel。
  - 写（广播）通过房间 goroutine 聚合后异步扇出到连接写队列（每连接单独写协程，避免阻塞）。
- 连接管理：`player_id -> Conn` 原子映射，重连时热切换。

### 3.6 缓存策略

- 房间快照（Redis）用于重连快速恢复与观战视图。
- 热门玩家资料/好友列表短期缓存（TTL 5-10 分钟）。
- 幂等键与限流计数在 Redis 实现，使用 Lua 保证原子性。

### 3.7 错误码规范（示例）

- 0: OK
- 1001: InvalidParam
- 1002: Unauthorized
- 1003: RoomNotFound
- 1004: RoomFull
- 1005: InvalidState
- 1006: InsufficientBalance
- 1099: InternalError

所有 TCP ACK 均返回 `{ret_code, ret_msg?, data?}`；REST 使用 HTTP 状态码 + 业务码。

### 3.8 认证与授权

- 登录发放 JWT（HS256），`sub=user_id`，短期有效；刷新令牌机制。
- TCP 首次鉴权消息或在 JOIN 前附带 `auth_token`；服务端校验并绑定连接属性。
- 敏感操作（下注/支付）叠加人机校验/频控策略。

---

## 四、开发任务清单

> 负责人仅作示例（BE1/BE2/INFRA），可按团队实际调整。

| 任务名称 | 优先级 | 预估工时 | 负责人 | 前置依赖 | 验收标准 | 标签 |
|---|---|---:|---|---|---|---|
| 定义统一错误码与 ACK 结构（TCP/REST） | 高 | 0.5d | BE1 | 无 | 所有 ACK 携带 ret_code/ret_msg；错误码文档化 | [STATUS: PARTIAL] [AREA: CORE] |
| 引入 JWT 鉴权与连接绑定 | 高 | 1.5d | BE1 | 用户登录接口 | 未鉴权请求被拒；JOIN 绑定 user_id 到连接 | [STATUS: PENDING] [AREA: SECURITY] |
| 房间 Actor 化改造（channel 事件循环） | 最高 | 3d | BE2 | 错误码 | FSM、入房、发牌、下注等通过事件串行；race 测试过 | [STATUS: PENDING] [AREA: CORE] |
| 阶段计时器与超时推进 | 高 | 1d | BE2 | Actor 化 | 抢庄/下注/摊牌超时自动推进；计时广播 | [STATUS: PENDING] [AREA: CORE] |
| 抢庄决策与下注校验实现 | 高 | 1.5d | BE2 | Actor + 计时器 | 正常/边界用例通过单测与集成测试 | [STATUS: PARTIAL] [AREA: GAMEPLAY] |
| 结算与账本写入（内存） | 高 | 1d | BE2 | 上述完成 | 使用 `CompareHands` 得出结果并广播；分数一致 | [STATUS: PENDING] [AREA: GAMEPLAY] |
| Redis 快照与重连恢复 | 中 | 1.5d | BE1 | Actor 化 | 断线重连后立刻拿到 `S2C_SyncRoomStateNtf` 完整视图 | [STATUS: PARTIAL] [AREA: RELIABILITY] |
| REST 用户注册/登录 | 高 | 1d | BE1 | DB 初始化 | 能注册/登录并签发 JWT，密码加盐哈希 | [STATUS: PENDING] [AREA: PLATFORM] |
| REST 好友增删/列表 | 中 | 1d | BE1 | 用户/鉴权 | 好友申请、接受、拉黑接口可用，带鉴权 | [STATUS: PENDING] [AREA: SOCIAL] |
| REST 物品清单与增减 | 中 | 1d | BE1 | 用户/鉴权 | 查询/修改库存的接口与权限校验 | [STATUS: PENDING] [AREA: ECONOMY] |
| 支付下单与回调框架 | 中 | 2d | BE1 | 用户/鉴权 | 能创建订单、接收模拟回调、更新状态与账本 | [STATUS: PENDING] [AREA: PAYMENT] |
| 限流与幂等（Redis + Lua） | 高 | 1d | BE2 | Redis | 热路径有频控；重复消息不重复执行 | [STATUS: PENDING] [AREA: RELIABILITY] |
| 广播扇出写队列与背压 | 高 | 1.5d | BE2 | Actor 化 | 高并发广播不阻塞房间循环；背压稳定 | [STATUS: PENDING] [AREA: PERFORMANCE] |
| 单元/集成测试完善（FSM/断线重连/边界） | 高 | 2d | QA/BE | 上述核心功能 | 覆盖率提升；关键路径测试绿色 | [STATUS: PENDING] [AREA: QUALITY] |
| 监控与日志结构化 | 中 | 1d | INFRA | 功能稳定 | QPS、延迟、错误率、各阶段停留时长可观测 | [STATUS: PENDING] [AREA: OBSERVABILITY] |

---

## 五、后续里程碑

- 里程碑1：核心对局稳定（Actor、计时、结算、重连）
- 里程碑2：账号与社交、道具、支付最小可用
- 里程碑3：分片与横向扩展、观战与回放

---

## 附：设计产物索引

- **[OpenAPI 合同]** `docs/contracts/openapi.yaml`
- **[架构设计]** `docs/design/architecture.md`
- **[游戏核心设计]** `docs/design/game-core.md`
- **[集成测试计划]** `docs/tests/integration.md`
