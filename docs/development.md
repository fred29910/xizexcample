# 开发指南

## 环境搭建

### 系统要求

- **操作系统**: Linux / macOS / Windows
- **Go 版本**: 1.25.1 或更高
- **Protocol Buffers**: 3.0+ （可选，用于修改协议）
- **Make**: 用于便捷构建（可选）
- **Git**: 版本控制

### 安装 Go

#### Linux / macOS

```bash
# 下载 Go
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz

# 解压
sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz

# 配置环境变量（添加到 ~/.bashrc 或 ~/.zshrc）
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# 验证安装
go version
```

#### Windows

1. 下载安装包：https://go.dev/dl/go1.25.1.windows-amd64.msi
2. 运行安装程序
3. 验证安装：`go version`

### 安装 Protocol Buffers（可选）

仅在需要修改 `.proto` 文件时安装。

#### Linux

```bash
# Debian/Ubuntu
sudo apt-get install -y protobuf-compiler

# 或从源码安装
wget https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip
unzip protoc-21.12-linux-x86_64.zip -d protoc
sudo mv protoc/bin/protoc /usr/local/bin/
sudo mv protoc/include/* /usr/local/include/
```

#### macOS

```bash
brew install protobuf
```

#### 安装 Go 插件

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

### 克隆项目

```bash
# 克隆仓库
git clone <repository-url>
cd xizexcample

# 或者如果已有本地项目
cd /opt/codes/workspace/games/xizexcample
```

### 安装依赖

```bash
# 下载所有依赖
go mod download

# 验证依赖
go mod verify

# 整理依赖（可选）
go mod tidy
```

## 项目配置

### 配置文件说明

**conf/zinx.json**:

```json
{
  "Name": "ZinxServerApp",      // 服务器名称
  "Host": "0.0.0.0",             // 监听地址
  "TCPPort": 8999,               // TCP 端口
  "MaxConn": 12000,              // 最大连接数
  "WorkerPoolSize": 10           // 工作协程池大小
}
```

**修改配置**:
- 开发环境：使用 `127.0.0.1` 或 `0.0.0.0`
- 生产环境：根据实际网络配置
- 端口：避免与其他服务冲突

### 环境变量（可扩展）

```bash
# 设置日志级别
export LOG_LEVEL=debug

# 设置配置文件路径
export CONFIG_PATH=./conf/zinx.json
```

## 编译与运行

### 使用 Makefile

```bash
# 构建项目
make build

# 运行项目
make run

# 运行测试
make test

# 清理构建产物
make clean

# 整理依赖
make tidy
```

### 使用 Go 命令

#### 开发环境运行

```bash
# 直接运行（带版本信息）
go run -ldflags="-X main.version=dev -X main.commit=$(git rev-parse --short HEAD)" .

# 或简单运行
go run main.go
```

#### 编译可执行文件

```bash
# 编译到 bin/ 目录
mkdir -p bin
go build -v -o bin/xizexcample .

# 编译时注入版本信息
VERSION=1.0.0
COMMIT=$(git rev-parse --short HEAD)
go build -ldflags="-s -w -X main.version=$VERSION -X main.commit=$COMMIT" -o bin/xizexcample .

# 运行编译后的程序
./bin/xizexcample
```

#### 交叉编译

```bash
# 编译 Linux 版本
GOOS=linux GOARCH=amd64 go build -o bin/xizexcample-linux .

# 编译 macOS 版本
GOOS=darwin GOARCH=amd64 go build -o bin/xizexcample-darwin .

# 编译 Windows 版本
GOOS=windows GOARCH=amd64 go build -o bin/xizexcample.exe .
```

### 启动服务器

```bash
# 前台运行
./bin/xizexcample

# 后台运行（Linux/macOS）
nohup ./bin/xizexcample > server.log 2>&1 &

# 使用 systemd（生产环境推荐）
sudo systemctl start xizexcample
```

### 停止服务器

```bash
# 查找进程
ps aux | grep xizexcample

# 优雅停止
kill -TERM <PID>

# 强制停止
kill -9 <PID>

# 使用 systemd
sudo systemctl stop xizexcample
```

## 开发流程

### 代码组织

```
internal/
├── conf/       # 配置管理
├── logic/      # 核心游戏逻辑（无网络依赖）
├── msg/        # 消息定义（自动生成）
├── pkg/        # 通用工具
├── router/     # 消息路由处理
└── server/     # 服务管理
```

### 添加新功能

#### 1. 定义消息协议

编辑 `api/proto/game.proto`:

```protobuf
// 添加新消息 ID
enum MsgID {
  // ...
  C2S_NEW_FEATURE_REQ = 107;
  S2C_NEW_FEATURE_ACK = 213;
}

// 定义消息结构
message C2S_NewFeatureReq {
  int32 param = 1;
}

message S2C_NewFeatureAck {
  int32 ret_code = 1;
  string result = 2;
}
```

#### 2. 生成 Go 代码

```bash
# 使用脚本
./scripts/gen_proto.sh

# 或手动执行
protoc --go_out=. api/proto/game.proto
```

#### 3. 实现业务逻辑

在 `internal/logic/` 中添加逻辑：

```go
// internal/logic/new_feature.go
package logic

func (r *Room) NewFeature(param int32) error {
    // 实现业务逻辑
    return nil
}
```

#### 4. 添加消息处理器

在 `internal/router/` 中添加处理器：

```go
// internal/router/new_feature.go
package router

import (
    "encoding/json"
    "github.com/aceld/zinx/ziface"
    "xizexcample/internal/msg"
    "xizexcample/internal/pkg/logger"
)

type NewFeatureHandler struct {
    BaseRouter
}

func (h *NewFeatureHandler) Handle(request ziface.IRequest) {
    // 1. 解析请求
    var req msg.C2S_NewFeatureReq
    err := json.Unmarshal(request.GetData(), &req)
    if err != nil {
        logger.ErrorLogger.Printf("Failed to unmarshal: %v", err)
        sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_NEW_FEATURE_ACK), "Invalid request")
        return
    }

    // 2. 调用业务逻辑
    player, room, err := GetPlayerAndRoom(request)
    if err != nil {
        sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_NEW_FEATURE_ACK), err.Error())
        return
    }

    err = room.NewFeature(req.Param)
    if err != nil {
        sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_NEW_FEATURE_ACK), err.Error())
        return
    }

    // 3. 发送响应
    ack := &msg.S2C_NewFeatureAck{
        RetCode: 0,
        Result:  "success",
    }
    ackData, _ := json.Marshal(ack)
    request.GetConnection().SendMsg(uint32(msg.MsgID_S2C_NEW_FEATURE_ACK), ackData)
}
```

#### 5. 注册路由

在 `internal/router/router.go` 中注册：

```go
func InitRouter(server ziface.IServer) {
    // ...
    server.AddRouter(107, &NewFeatureHandler{})
}
```

### 修改现有功能

1. 定位相关代码（router/logic）
2. 修改业务逻辑
3. 更新测试用例
4. 运行测试验证
5. 更新文档

### 提交代码

```bash
# 查看修改
git status
git diff

# 暂存修改
git add .

# 提交（使用有意义的消息）
git commit -m "feat: add new feature for xxx"

# 推送到远程
git push origin feature-branch
```

## 测试

### 单元测试

#### 运行所有测试

```bash
# 使用 Makefile
make test

# 或使用 go test
go test ./...

# 详细输出
go test -v ./...

# 并发测试（检测竞态条件）
go test -race ./...
```

#### 运行特定测试

```bash
# 测试特定包
go test -v ./internal/logic/

# 测试特定文件
go test -v ./internal/logic/room_test.go

# 测试特定函数
go test -v -run TestRoomFSM ./internal/logic/
```

#### 查看覆盖率

```bash
# 生成覆盖率报告
go test -cover ./...

# 生成详细覆盖率文件
go test -coverprofile=coverage.out ./...

# 查看 HTML 报告
go tool cover -html=coverage.out
```

### 集成测试

```bash
# 运行集成测试（如果存在）
go test -v ./tests/integration/...
```

### E2E 测试

```bash
# 启动服务器（后台）
./bin/xizexcample &
SERVER_PID=$!

# 运行 E2E 测试
go test -v ./tests/e2e/...

# 停止服务器
kill $SERVER_PID
```

### 性能测试

```bash
# 运行基准测试
go test -bench=. ./internal/logic/

# 生成 CPU profile
go test -bench=. -cpuprofile=cpu.prof ./internal/logic/

# 查看 profile
go tool pprof cpu.prof
```

### 编写测试用例

#### 单元测试示例

```go
// internal/logic/room_test.go
package logic

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestRoom_AddPlayer(t *testing.T) {
    room := NewRoom(1)
    player := NewPlayer(123, "TestPlayer", nil)

    // 测试添加玩家
    err := room.AddPlayer(player)
    assert.NoError(t, err)
    assert.Equal(t, 1, room.GetPlayerCount())

    // 测试重复添加
    err = room.AddPlayer(player)
    assert.Error(t, err)
}

func TestRoom_IsFull(t *testing.T) {
    room := NewRoom(1)

    // 添加 5 个玩家
    for i := 0; i < 5; i++ {
        player := NewPlayer(int64(i), "Player", nil)
        room.AddPlayer(player)
    }

    assert.True(t, room.IsFull())
}
```

## 调试

### 使用日志

```go
// 添加日志
logger.InfoLogger.Printf("Room %d: player %d joined", roomID, playerID)
logger.ErrorLogger.Printf("Failed to deal cards: %v", err)
```

### 使用 Go 调试器（Delve）

```bash
# 安装 delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试运行
dlv debug .

# 设置断点
(dlv) break internal/router/join_room.go:20

# 运行
(dlv) continue

# 查看变量
(dlv) print player

# 退出
(dlv) quit
```

### 使用 IDE 调试

#### VSCode

创建 `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}",
            "env": {},
            "args": []
        }
    ]
}
```

#### GoLand

1. 打开 `Run` → `Edit Configurations`
2. 添加 `Go Build` 配置
3. 设置 `Package path` 为 `.`
4. 点击调试按钮

### 网络调试

#### 使用 Wireshark

```bash
# 捕获本地回环
sudo wireshark -i lo -f "tcp port 8999"
```

#### 使用 tcpdump

```bash
# 捕获并保存
sudo tcpdump -i any -s 0 -w capture.pcap 'tcp port 8999'

# 查看捕获文件
tcpdump -r capture.pcap -X
```

### 性能分析

#### CPU Profiling

```go
// 在 main.go 中添加
import _ "net/http/pprof"
import "net/http"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

访问 `http://localhost:6060/debug/pprof/`

#### 内存分析

```bash
# 获取内存快照
curl http://localhost:6060/debug/pprof/heap > heap.prof

# 分析
go tool pprof heap.prof
```

## 代码规范

### Go 代码风格

遵循官方 [Effective Go](https://go.dev/doc/effective_go) 指南：

1. **命名规范**:
   - 包名：小写，简短
   - 变量名：驼峰命名
   - 导出：首字母大写
   - 私有：首字母小写

2. **注释**:
   ```go
   // NewRoom 创建一个新房间
   // roomID: 房间唯一标识
   func NewRoom(roomID int32) *Room {
       // ...
   }
   ```

3. **错误处理**:
   ```go
   if err != nil {
       logger.ErrorLogger.Printf("Failed: %v", err)
       return err
   }
   ```

4. **格式化**:
   ```bash
   # 格式化代码
   go fmt ./...
   
   # 或使用 goimports
   goimports -w .
   ```

### 代码检查

```bash
# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 运行检查
golangci-lint run

# 自动修复
golangci-lint run --fix
```

### Git Commit 规范

使用 [Conventional Commits](https://www.conventionalcommits.org/)：

```
feat: 添加新功能
fix: 修复 bug
docs: 更新文档
style: 代码格式调整
refactor: 重构
test: 添加测试
chore: 构建/工具变动
```

示例：
```bash
git commit -m "feat: add player ready handler"
git commit -m "fix: resolve race condition in room manager"
git commit -m "docs: update API protocol documentation"
```

## 部署

### 单机部署

#### 直接运行

```bash
# 复制配置文件
sudo mkdir -p /opt/xizexcample
sudo cp -r conf /opt/xizexcample/
sudo cp bin/xizexcample /opt/xizexcample/

# 运行
cd /opt/xizexcample
./xizexcample
```

#### 使用 systemd

创建服务文件 `/etc/systemd/system/xizexcample.service`:

```ini
[Unit]
Description=Xizexcample Game Server
After=network.target

[Service]
Type=simple
User=gameserver
WorkingDirectory=/opt/xizexcample
ExecStart=/opt/xizexcample/xizexcample
Restart=always
RestartSec=10

# 日志
StandardOutput=journal
StandardError=journal

# 限制
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable xizexcample
sudo systemctl start xizexcample

# 查看状态
sudo systemctl status xizexcample

# 查看日志
sudo journalctl -u xizexcample -f
```

### Docker 部署

#### Dockerfile

创建 `Dockerfile`:

```dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o xizexcample .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/xizexcample .
COPY --from=builder /app/conf ./conf

EXPOSE 8999
CMD ["./xizexcample"]
```

#### 构建和运行

```bash
# 构建镜像
docker build -t xizexcample:latest .

# 运行容器
docker run -d \
  --name xizexcample \
  -p 8999:8999 \
  -v /opt/xizexcample/conf:/root/conf \
  xizexcample:latest

# 查看日志
docker logs -f xizexcample

# 停止容器
docker stop xizexcample
```

#### Docker Compose

创建 `docker-compose.yml`:

```yaml
version: '3.8'

services:
  gameserver:
    build: .
    ports:
      - "8999:8999"
    volumes:
      - ./conf:/root/conf
    restart: unless-stopped
    environment:
      - LOG_LEVEL=info
```

运行：

```bash
docker-compose up -d
docker-compose logs -f
docker-compose down
```

### 负载均衡（可扩展）

#### Nginx 配置

```nginx
upstream gameservers {
    least_conn;
    server 192.168.1.10:8999;
    server 192.168.1.11:8999;
    server 192.168.1.12:8999;
}

server {
    listen 8999;
    
    location / {
        proxy_pass http://gameservers;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 监控

#### 健康检查

```bash
# 简单端口检查
nc -zv localhost 8999

# 或使用 TCP 连接测试
timeout 1 bash -c 'cat < /dev/null > /dev/tcp/localhost/8999'
```

#### 日志监控

```bash
# 实时查看错误日志
tail -f server.log | grep ERROR

# 统计错误数量
grep ERROR server.log | wc -l
```

## 故障排查

### 常见问题

#### 1. 端口被占用

```bash
# 查找占用端口的进程
lsof -i:8999
netstat -tulpn | grep 8999

# 杀死进程
kill -9 <PID>
```

#### 2. 依赖问题

```bash
# 清理并重新下载
go clean -modcache
go mod download
```

#### 3. 编译错误

```bash
# 更新依赖
go get -u ./...
go mod tidy
```

#### 4. 性能问题

```bash
# 检查协程泄漏
curl http://localhost:6060/debug/pprof/goroutine?debug=2

# 检查内存使用
free -h
top
```

### 日志分析

```bash
# 查看最近 100 行
tail -100 server.log

# 搜索特定错误
grep "Failed to" server.log

# 统计错误类型
grep ERROR server.log | awk '{print $NF}' | sort | uniq -c
```

## 持续集成/部署（CI/CD）

### GitHub Actions 示例

创建 `.github/workflows/ci.yml`:

```yaml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.25'
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3

  build:
    runs-on: ubuntu-latest
    needs: test
    steps:
    - uses: actions/checkout@v3
    
    - name: Build
      run: make build
    
    - name: Upload artifact
      uses: actions/upload-artifact@v3
      with:
        name: xizexcample
        path: bin/xizexcample
```

## 资源链接

### 官方文档

- [Go 官方文档](https://go.dev/doc/)
- [Zinx 框架](https://github.com/aceld/zinx)
- [Protocol Buffers](https://protobuf.dev/)

### 学习资源

- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [The Go Programming Language](https://www.gopl.io/)

### 工具推荐

- **IDE**: VSCode, GoLand
- **调试**: Delve
- **性能分析**: pprof
- **代码检查**: golangci-lint
- **API测试**: Postman, curl

## 最佳实践

1. **版本控制**: 始终使用 Git 管理代码
2. **测试驱动**: 编写测试覆盖关键逻辑
3. **代码审查**: Pull Request 前进行 Code Review
4. **文档更新**: 代码变更同步更新文档
5. **日志记录**: 记录关键操作和错误
6. **性能监控**: 定期进行性能分析
7. **安全意识**: 验证所有输入，防范攻击
8. **优雅关闭**: 处理信号，保存状态
9. **错误处理**: 完善的错误处理机制
10. **持续改进**: 定期重构和优化代码
