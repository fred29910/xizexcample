# 通信协议文档

## 概述

本文档详细说明斗牛游戏服务器的客户端/服务端通信协议，包括消息格式、消息ID定义、数据结构等。

## 协议格式

### 消息封装格式

采用 Zinx 框架的标准消息格式 + Protocol Buffers 序列化：

```
┌──────────────┬──────────────┬─────────────────────┐
│  Message ID  │  Data Length │   Message Data      │
│   (4 bytes)  │  (4 bytes)   │   (variable)        │
│   uint32     │   uint32     │   []byte (Protobuf) │
└──────────────┴──────────────┴─────────────────────┘
```

- **Message ID**: 消息类型标识（详见消息ID定义）
- **Data Length**: 消息体长度（字节数）
- **Message Data**: Protocol Buffers 序列化的消息体

### 字节序

- 网络字节序（Big Endian）
- Protocol Buffers 自动处理序列化

## 消息ID定义

### 客户端到服务端（C2S）

| 消息ID | 消息名称 | 描述 |
|--------|---------|------|
| 101 | C2S_JOIN_ROOM_REQ | 加入房间请求 |
| 102 | C2S_PLAYER_READY_REQ | 玩家准备请求 |
| 103 | C2S_BID_BANKER_REQ | 抢庄请求 |
| 104 | C2S_PLACE_BET_REQ | 下注请求 |
| 105 | C2S_SHOWDOWN_REQ | 摊牌请求 |
| 106 | C2S_LEAVE_ROOM_REQ | 离开房间请求 |

### 服务端到客户端（S2C）

#### 响应消息（ACK）

| 消息ID | 消息名称 | 描述 |
|--------|---------|------|
| 201 | S2C_JOIN_ROOM_ACK | 加入房间响应 |
| 210 | S2C_BID_BANKER_ACK | 抢庄响应 |
| 211 | S2C_PLACE_BET_ACK | 下注响应 |
| 212 | S2C_SHOWDOWN_ACK | 摊牌响应 |

#### 通知消息（NTF）

| 消息ID | 消息名称 | 描述 |
|--------|---------|------|
| 202 | S2C_SYNC_ROOM_STATE_NTF | 房间状态同步通知 |
| 203 | S2C_GAME_START_NTF | 游戏开始通知 |
| 204 | S2C_DEAL_CARDS_NTF | 发牌通知 |
| 205 | S2C_BID_BANKER_NTF | 抢庄阶段通知 |
| 206 | S2C_BET_NTF | 下注阶段通知 |
| 207 | S2C_SHOWDOWN_NTF | 摊牌阶段通知 |
| 208 | S2C_GAME_RESULT_NTF | 游戏结果通知 |
| 209 | S2C_PLAYER_LEAVE_NTF | 玩家离开通知 |

## 数据结构定义

### 基础数据类型

#### Card（扑克牌）

```protobuf
message Card {
  Suit suit = 1;  // 花色
  Rank rank = 2;  // 点数
}
```

#### Suit（花色）

```protobuf
enum Suit {
  SUIT_UNKNOWN = 0;
  SPADES = 1;      // ♠ 黑桃
  HEARTS = 2;      // ♥ 红心
  CLUBS = 3;       // ♣ 梅花
  DIAMONDS = 4;    // ♦ 方块
}
```

#### Rank（点数）

```protobuf
enum Rank {
  RANK_UNKNOWN = 0;
  ACE = 1;         // A
  TWO = 2;         // 2
  THREE = 3;       // 3
  FOUR = 4;        // 4
  FIVE = 5;        // 5
  SIX = 6;         // 6
  SEVEN = 7;       // 7
  EIGHT = 8;       // 8
  NINE = 9;        // 9
  TEN = 10;        // 10
  JACK = 11;       // J
  QUEEN = 12;      // Q
  KING = 13;       // K
}
```

#### PlayerStatus（玩家状态）

```protobuf
enum PlayerStatus {
  STATUS_UNKNOWN = 0;
  WAITING = 1;     // 等待状态
  READY = 2;       // 准备状态
  PLAYING = 3;     // 游戏中
}
```

#### GameState（游戏阶段）

```protobuf
enum GameState {
  STATE_UNKNOWN = 0;
  WAITING_FOR_PLAYERS = 1;  // 等待玩家
  DEALING = 2;              // 发牌
  BIDDING = 3;              // 抢庄
  BETTING = 4;              // 下注
  SHOWDOWN = 5;             // 摊牌
  SETTLEMENT = 6;           // 结算
}
```

#### CardPattern（牌型）

```protobuf
enum CardPattern {
  PATTERN_UNKNOWN = 0;
  NO_NIU = 1;           // 无牛
  NIU_1 = 2;            // 牛一
  NIU_2 = 3;            // 牛二
  NIU_3 = 4;            // 牛三
  NIU_4 = 5;            // 牛四
  NIU_5 = 6;            // 牛五
  NIU_6 = 7;            // 牛六
  NIU_7 = 8;            // 牛七
  NIU_8 = 9;            // 牛八
  NIU_9 = 10;           // 牛九
  NIU_NIU = 11;         // 牛牛
  FIVE_FLOWER_NIU = 12; // 五花牛
  BOMB_NIU = 13;        // 炸弹牛
  FIVE_SMALL_NIU = 14;  // 五小牛
}
```

### 复合数据结构

#### PlayerInfo（玩家信息）

```protobuf
message PlayerInfo {
  int64 player_id = 1;          // 玩家ID
  string nickname = 2;          // 昵称
  int64 score = 3;              // 积分
  PlayerStatus status = 4;      // 玩家状态
  bool is_banker = 5;           // 是否是庄家
  repeated Card hand = 6;       // 手牌（仅发给玩家本人）
  CardPattern card_pattern = 7; // 牌型（摊牌后）
}
```

#### RoomInfo（房间信息）

```protobuf
message RoomInfo {
  int32 room_id = 1;              // 房间ID
  repeated PlayerInfo players = 2; // 玩家列表
  GameState game_state = 3;       // 游戏阶段
  int64 banker_id = 4;            // 庄家ID
}
```

#### PlayerResult（玩家结算结果）

```protobuf
message PlayerResult {
  int64 player_id = 1;         // 玩家ID
  repeated Card hand = 2;      // 手牌
  CardPattern card_pattern = 3; // 牌型
  int64 score_change = 4;      // 分数变化（正数=赢，负数=输）
  int64 final_score = 5;       // 最终分数
}
```

## 消息详细定义

### 1. 加入房间

#### C2S_JoinRoomReq（请求）

```protobuf
message C2S_JoinRoomReq {
  int32 room_id = 1;  // 房间ID，0表示随机分配
}
```

**字段说明**:
- `room_id`: 要加入的房间ID，如果为0则服务器自动分配房间

#### S2C_JoinRoomAck（响应）

```protobuf
message S2C_JoinRoomAck {
  int32 ret_code = 1;      // 返回码，0表示成功
  RoomInfo room_info = 2;  // 房间信息
}
```

**返回码定义**:
- `0`: 成功
- `1`: 房间已满
- `2`: 房间不存在且无法创建
- `3`: 玩家已在其他房间

**使用场景**:
- 新玩家加入房间
- 断线重连

### 2. 玩家准备

#### C2S_PlayerReadyReq（请求）

```protobuf
message C2S_PlayerReadyReq {
  bool is_ready = 1;  // true=准备, false=取消准备
}
```

**字段说明**:
- `is_ready`: 准备状态标识

**触发条件**:
- 游戏处于 `WAITING_FOR_PLAYERS` 阶段
- 至少2名玩家准备后自动开始游戏

### 3. 抢庄

#### C2S_BidBankerReq（请求）

```protobuf
message C2S_BidBankerReq {
  int32 multiple = 1;  // 抢庄倍数，0表示不抢
}
```

**字段说明**:
- `multiple`: 抢庄倍数，范围 0-10，0表示不抢庄

#### S2C_BidBankerAck（响应）

```protobuf
message S2C_BidBankerAck {
  int32 ret_code = 1;   // 返回码
  int64 player_id = 2;  // 成为庄家的玩家ID
}
```

**返回码定义**:
- `0`: 成功
- `1`: 游戏阶段错误
- `2`: 庄家已确定
- `3`: 无效的倍数

#### S2C_BidBankerNtf（通知）

```protobuf
message S2C_BidBankerNtf {
  int32 countdown = 1;  // 抢庄倒计时（秒）
}
```

**广播时机**:
- 进入抢庄阶段时广播给所有玩家

### 4. 下注

#### C2S_PlaceBetReq（请求）

```protobuf
message C2S_PlaceBetReq {
  int32 multiple = 1;  // 下注倍数
}
```

**字段说明**:
- `multiple`: 下注倍数，范围 1-10

**限制**:
- 仅闲家可以下注
- 庄家不需要下注
- 每个玩家只能下注一次

#### S2C_PlaceBetAck（响应）

```protobuf
message S2C_PlaceBetAck {
  int32 ret_code = 1;   // 返回码
  int32 multiple = 2;   // 下注倍数
}
```

**返回码定义**:
- `0`: 成功
- `1`: 游戏阶段错误
- `2`: 已下注
- `3`: 庄家不能下注
- `4`: 余额不足
- `5`: 倍数无效

#### S2C_BetNtf（通知）

```protobuf
message S2C_BetNtf {
  int64 banker_id = 1;  // 庄家ID
  int32 countdown = 2;  // 下注倒计时（秒）
}
```

**广播时机**:
- 进入下注阶段时广播给所有玩家

### 5. 摊牌

#### C2S_ShowdownReq（请求）

```protobuf
message C2S_ShowdownReq {
  repeated Card sorted_hand = 1;  // 客户端排序后的手牌
}
```

**字段说明**:
- `sorted_hand`: 客户端排序后的手牌（可选，服务端会重新计算）

**注意事项**:
- 服务端会验证手牌的合法性
- 服务端重新计算牌型，不信任客户端结果

#### S2C_ShowdownAck（响应）

```protobuf
message S2C_ShowdownAck {
  int32 ret_code = 1;  // 返回码
}
```

**返回码定义**:
- `0`: 成功
- `1`: 游戏阶段错误
- `2`: 手牌无效

#### S2C_ShowdownNtf（通知）

```protobuf
message S2C_ShowdownNtf {
  int32 countdown = 1;  // 摊牌倒计时（秒）
}
```

**广播时机**:
- 进入摊牌阶段时广播给所有玩家

### 6. 游戏开始

#### S2C_GameStartNtf（通知）

```protobuf
message S2C_GameStartNtf {
  int64 banker_id = 1;  // 庄家ID（如果已确定）
}
```

**广播时机**:
- 游戏从 `WAITING_FOR_PLAYERS` 转到 `DEALING` 时

### 7. 发牌

#### S2C_DealCardsNtf（通知）

```protobuf
message S2C_DealCardsNtf {
  repeated Card hand = 1;  // 玩家的手牌（5张）
}
```

**发送时机**:
- 发牌阶段，单独发送给每个玩家
- 仅包含该玩家自己的手牌（私密性）

### 8. 游戏结果

#### S2C_GameResultNtf（通知）

```protobuf
message S2C_GameResultNtf {
  repeated PlayerResult results = 1;  // 所有玩家的结算结果
}
```

**广播时机**:
- 结算阶段完成后广播给所有玩家

**内容包括**:
- 每个玩家的手牌
- 牌型
- 输赢金额
- 最终分数

### 9. 房间状态同步

#### S2C_SyncRoomStateNtf（通知）

```protobuf
message S2C_SyncRoomStateNtf {
  RoomInfo room_info = 1;  // 完整的房间信息
}
```

**广播时机**:
- 玩家加入/离开房间
- 游戏状态变化
- 断线重连后同步状态

### 10. 玩家离开

#### C2S_LeaveRoomReq（请求）

```protobuf
message C2S_LeaveRoomReq {
  // 空消息
}
```

#### S2C_PlayerLeaveNtf（通知）

```protobuf
message S2C_PlayerLeaveNtf {
  int64 player_id = 1;  // 离开的玩家ID
}
```

**广播时机**:
- 玩家主动离开房间
- 玩家断线超时被踢出

## 通信流程示例

### 完整游戏流程

```
1. 玩家加入房间
Client A → [101] C2S_JoinRoomReq
Server → [201] S2C_JoinRoomAck → Client A
Server → [202] S2C_SyncRoomStateNtf → All Clients

2. 玩家准备
Client A → [102] C2S_PlayerReadyReq {is_ready: true}
Client B → [102] C2S_PlayerReadyReq {is_ready: true}

3. 游戏开始
Server → [203] S2C_GameStartNtf → All Clients

4. 发牌
Server → [204] S2C_DealCardsNtf → Client A (A的手牌)
Server → [204] S2C_DealCardsNtf → Client B (B的手牌)

5. 抢庄阶段
Server → [205] S2C_BidBankerNtf → All Clients
Client A → [103] C2S_BidBankerReq {multiple: 5}
Server → [210] S2C_BidBankerAck → Client A
Server → [202] S2C_SyncRoomStateNtf → All Clients

6. 下注阶段
Server → [206] S2C_BetNtf → All Clients
Client B → [104] C2S_PlaceBetReq {multiple: 3}
Server → [211] S2C_PlaceBetAck → Client B

7. 摊牌阶段
Server → [207] S2C_ShowdownNtf → All Clients
Client A → [105] C2S_ShowdownReq
Client B → [105] C2S_ShowdownReq
Server → [212] S2C_ShowdownAck → Clients

8. 结算
Server → [208] S2C_GameResultNtf → All Clients
Server → [202] S2C_SyncRoomStateNtf → All Clients
```

### 断线重连流程

```
1. 玩家断线
Client disconnect
Server → SetPlayerOffline()

2. 玩家重连
Client → [101] C2S_JoinRoomReq (same room_id)
Server detects existing player
Server → [201] S2C_JoinRoomAck (with full room state)
Server → [202] S2C_SyncRoomStateNtf → Reconnected Client

3. 恢复游戏
Client receives current game state
Continue from current phase
```

## 错误处理

### 通用错误码

所有响应消息（ACK）都包含 `ret_code` 字段：

| 返回码 | 含义 | 说明 |
|--------|------|------|
| 0 | 成功 | 请求成功处理 |
| 1 | 通用错误 | 未分类的错误 |
| 2 | 游戏阶段错误 | 当前阶段不允许此操作 |
| 3 | 参数无效 | 请求参数不合法 |
| 4 | 房间不存在 | 指定的房间不存在 |
| 5 | 房间已满 | 房间人数已达上限 |
| 6 | 权限不足 | 没有执行该操作的权限 |
| 7 | 余额不足 | 玩家积分不足 |
| 8 | 重复操作 | 已执行过的操作不能重复 |

### 客户端异常处理

**连接异常**:
- 网络断开：尝试重连
- 服务器关闭：提示并退出

**协议异常**:
- 消息解析失败：记录日志并忽略
- 收到未知消息ID：记录日志并忽略

**业务异常**:
- 操作失败（ret_code != 0）：根据错误码提示用户

## 协议版本管理

### 当前版本

- 协议版本: `v1.0`
- Protobuf版本: `proto3`

### 版本兼容性

- 向后兼容：新增字段使用 optional 或 repeated
- 废弃字段：保留字段编号，标记为 deprecated
- 重大变更：增加新的消息ID或版本号

### 版本协商（待实现）

```protobuf
message C2S_HandshakeReq {
  string protocol_version = 1;
  string client_version = 2;
}

message S2C_HandshakeAck {
  int32 ret_code = 1;
  string server_version = 2;
}
```

## 安全性考虑

### 数据验证

**服务端验证**:
- 所有客户端输入必须验证
- 手牌合法性检查
- 下注金额范围检查
- 操作权限验证

**客户端验证**:
- 基本输入合法性检查
- 不信任服务端关键数据（金币等）

### 防作弊

**服务端权威**:
- 所有游戏逻辑在服务端执行
- 客户端仅负责展示和输入
- 牌型计算完全由服务端完成
- 手牌不在客户端之间传递

**数据加密（可扩展）**:
- TLS/SSL 加密通信
- 敏感数据加密存储
- 防重放攻击

## 性能优化

### 消息压缩

- Protocol Buffers 本身高效
- 可选：启用 gzip 压缩大消息

### 批量操作

- 广播消息批量发送
- 减少网络往返次数

### 心跳机制（可扩展）

```protobuf
message C2S_HeartbeatReq {
  int64 timestamp = 1;
}

message S2C_HeartbeatAck {
  int64 timestamp = 1;
  int64 server_time = 2;
}
```

## 调试支持

### 消息日志

建议记录以下信息：
- 消息ID
- 发送/接收时间
- 消息体（JSON格式）
- 玩家ID
- 房间ID

### 抓包分析

- 使用 Wireshark + Protobuf 插件
- 导入 .proto 文件解析消息

## 附录

### Protocol Buffers 文件

完整的 .proto 文件位于：`api/proto/game.proto`

### 代码生成

```bash
# 生成 Go 代码
protoc --go_out=. api/proto/game.proto

# 或使用项目脚本
./scripts/gen_proto.sh
```

### 客户端SDK（待实现）

建议为不同平台提供SDK：
- Go Client SDK
- Unity C# SDK
- Cocos Creator TypeScript SDK
