# Quickstart: “斗牛牛”棋牌游戏服务器

This guide provides the basic steps to get the game server running.

## Prerequisites

- Go (latest stable version)
- Protocol Buffers Compiler (`protoc`)
- `protoc-gen-go` plugin

## 1. Compile Protocol Buffers

The project uses Protocol Buffers for client-server communication. You must compile the `.proto` file to generate the Go code for messages.

From the project root directory, run the following command:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    api/proto/game.proto
```

This will generate a `game.pb.go` file in the `internal/msg/` directory.

## 2. Build and Run the Server

Once the protocol code is generated, you can build and run the server.

```bash
# Navigate to the project root
go build -o niuniu_server .
./niuniu_server
```

The server will start and listen for incoming TCP connections on the configured port.

## 3. Project Structure Overview

- `api/proto/`: Contains the `.proto` definition files.
- `internal/conf/`: Configuration files.
- `internal/logic/`: Core game logic (Room, Player, FSM).
- `internal/msg/`: Generated Protobuf Go code.
- `internal/router/`: Zinx message routers.
- `internal/server/`: Server initialization and `RoomManager`.
- `main.go`: Main application entry point.
- `tests/`: Unit and E2E tests.