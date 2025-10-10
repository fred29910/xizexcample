# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

使用 Go 语言和 Zinx 框架构建一个功能完善、性能稳定的“斗牛牛”棋牌游戏服务器。服务器将管理完整的游戏生命周期，处理实时通信，并在内存中维护所有游戏状态。

## Technical Context

**Language/Version**: Go (latest stable version)
**Primary Dependencies**: Zinx, Protobuf
**Storage**: N/A (All game state is in-memory)
**Testing**: Go standard library testing, testify
**Target Platform**: Linux server
**Project Type**: Single project (game server)
**Performance Goals**: Support 100 concurrent rooms (500 users) with <200ms latency
**Constraints**: TCP long connections, custom binary protocol (MsgID + Protobuf)
**Scale/Scope**: Initial version supports the full "斗牛牛" game loop for up to 5 players per room.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [ ] **I. 严格的代码风格**: 所有代码是否遵循 `gofmt` 和 `golint` 规范？
- [ ] **II. 提交前静态检查**: 是否已配置并将在提交前运行 `lint` 检查？
- [ ] **III. 全面的单元测试**: 新功能是否计划了完整的单元测试覆盖？
- [ ] **IV. 端到端质量保证**: 新功能是否需要并计划了 E2E 测试？

## Project Structure

### Documentation (this feature)

```
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```
# Single project (game server)
.
├── api/              # Protobuf definition files (.proto)
│   └── proto/
│       └── game.proto
├── internal/         # Internal application code
│   ├── conf/         # Configuration
│   ├── logic/        # Core business logic (Room, Player, Game State Machine)
│   ├── msg/          # Protobuf generated code
│   ├── router/       # Zinx message routers
│   └── server/       # Server setup and management (RoomManager)
├── main.go           # Application entry point
└── tests/            # Test files
    ├── e2e/
    └── unit/
```

**Structure Decision**: A single Go project structure is chosen. The `internal` directory houses all core application logic, preventing it from being imported by other projects. The `api` directory contains the Protobuf definitions, which act as the contract between the server and client. `main.go` initializes and starts the server.

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
