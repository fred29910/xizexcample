# 游戏核心设计（FSM/Actor/计时器/结算/幂等）

更新时间：2025-10-11 22:22:24 +08:00

## 1. 状态机（FSM）

- 阶段：`WAITING_FOR_PLAYERS` → `DEALING` → `BIDDING` → `BETTING` → `SHOWDOWN` → `SETTLEMENT` → `WAITING_FOR_PLAYERS`
- 约束：
  - 只有在 `WAITING_FOR_PLAYERS` 且满足最小玩家数且全部 `READY` 才能 `StartGame()`
  - `DEALING` 完成后进入 `BIDDING`
  - `BIDDING` 需在倒计时内收集“抢庄倍数”，超时按默认 0；并列最高者随机为庄
  - `BETTING` 需在倒计时内收集下注，超时自动下注最小额或视为放弃（可配置）
  - `SHOWDOWN` 收集摊牌顺序（或由服务端按手牌计算），超时自动摊牌
  - `SETTLEMENT` 计算输赢与分数变化并广播

## 2. Actor（每房间一个事件循环）

- 事件结构：
```go
 type EventType int
 const (
   EJoin EventType = iota; ELeave; EReady; EDealTimeout; EBid; EBidTimeout; EBet; EBetTimeout; EShow; EShowTimeout; EShutdown
 )
 type Event struct { Type EventType; PlayerID int64; Payload any; At time.Time }
```
- 主循环：串行处理事件，所有写操作仅在 Actor 内执行；避免锁竞争。
- 入口：TCP 路由（`internal/router/*`）将 C2S 转为 `Event` 投递到房间 `chan Event`。

## 3. 计时器

- 每阶段开始时 `schedule(after, EventTimeout)`；到时投递 `*Timeout` 事件。
- 计时广播：通过 `S2C_*_Ntf.countdown` 值向在线玩家周期推送。

## 4. 抢庄与下注

- 抢庄：
```go
 collectBid(pID, multiple)
 // 结束条件：全部提交或超时
 resolveBanker():
   maxM := max(bid[p])
   cand := {p | bid[p]==maxM}
   banker := cand[ rand(0..len(cand)-1) ]
```
- 下注校验：
  - 金额>0、≤余额、在允许倍数集合内（1/2/3/5 等可配）
  - 阶段校验只允许在 `BETTING`
  - 幂等：同一 `client_seq` 只生效一次

## 5. 摊牌与结算

- 手牌计算：可以使用客户端排序后上报，也可服务端独立计算校验（防作弊优先服务端计算）。
- 结算：使用 `internal/logic/bull_logic.go: CompareHands()`，公式：
```go
 delta = base * playerBetMultiple * bankerMultiple * sign(result)
 p.Score += delta; banker.Score -= delta
 ledger(user_id, delta, round_id, reason)
```
- 结果广播：`S2C_GameResultNtf{results[]}`，并重置局内状态。

## 6. 数据一致性与重连

- 快照：进入每阶段/关键事件后，写入 Redis `room:{room_id}:snapshot`（含玩家列表、阶段、庄家、倒计时、可选显示手牌）。
- 重连：JOIN 时若发现离线标记，则先恢复连接与在线状态，再下发最新快照。

## 7. 幂等与安全

- 幂等键：`idemp:{user_id}:{client_seq}` 设置 TTL；重复消息直接 ACK 成功。
- 限流：IP/账号维度令牌桶；超限返回错误码 429（REST）或业务码。
- 鉴权：JOIN 前校验 JWT，写入连接属性 `playerID`。

## 8. 广播与背压

- 每连接单独写协程 + 有界队列；超出阈值丢弃非关键 NTF 或降采样；慢连接可断开。
- 聚合：短时间内的状态同步合并为一次 `S2C_SyncRoomStateNtf`。

## 9. 错误码与 ACK

- 统一 `ret_code/ret_msg`；典型：InvalidParam/Unauthorized/RoomFull/InvalidState/InsufficientBalance/InternalError。

## 10. 任务落地建议（实现顺序）

1) Actor 框架与事件定义 → 2) 计时器/倒计时 → 3) 抢庄/下注/摊牌/结算 → 4) 快照与重连 → 5) 广播扇出与背压 → 6) 幂等与限流。
