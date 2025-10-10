# Tasks: “斗牛牛”棋牌游戏服务器

**Input**: Design documents from `/specs/001-/`
**Prerequisites**: plan.md, spec.md, data-model.md, api/proto/game.proto

**Tests**: Per the constitution, tests are MANDATORY. All new features must include unit and E2E tests.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions
- Paths shown below assume single project structure from `plan.md`.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure.

- [x] T001 [P] Initialize Go module with `go mod init`.
- [x] T002 [P] Add Zinx and Protobuf dependencies to `go.mod`.
- [x] T003 [P] Create the project directory structure as defined in `plan.md`.
- [x] T004 [P] Configure linting and formatting tools (`gofmt`, `golint`) and setup pre-commit hooks.
- [x] T005 [P] Configure unit and E2E testing frameworks.
- [x] T006 [P] Implement a script to compile `.proto` files into `internal/msg/game.pb.go`.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented.

- [x] T007 Implement the main server entry point in `main.go` to initialize and start the Zinx server.
- [x] T008 [P] Implement the global `RoomManager` singleton in `internal/server/room_manager.go` with methods for creating, finding, and deleting rooms.
- [x] T009 [P] Define the core `Player` struct in `internal/logic/player.go` to hold player state.
- [x] T010 [P] Define the `Card` and deck logic in `internal/logic/deck.go`, including shuffle functionality.
- [x] T011 Implement the basic `Room` struct in `internal/logic/room.go`, including player management (add/remove).
- [x] T012 Setup the basic Zinx message router in `internal/router/router.go` to link `MsgID`s to handlers.

**Checkpoint**: Foundation ready - user story implementation can now begin.

---

## Phase 3: User Story 1 - 玩家加入并开始游戏 (Priority: P1) 🎯 MVP

**Goal**: Players can join a room, and the game starts when enough players are ready.
**Independent Test**: Simulate multiple clients joining, sending "ready" status, and verify that the `S2C_GameStartNtf` and `S2C_DealCardsNtf` are broadcast correctly.

### Tests for User Story 1 (MANDATORY)

- [x] T013 [P] [US1] Unit test for `RoomManager.CreateRoom` and `RoomManager.GetRoom` in `internal/server/room_manager_test.go`.
- [x] T014 [P] [US1] Unit test for `Room.AddPlayer` and `Room.RemovePlayer` in `internal/logic/room_test.go`.
- [x] T015 [US1] E2E test for the full join-and-start flow in `tests/e2e/join_game_test.go`.

### Implementation for User Story 1

- [x] T016 [US1] Implement the handler for `C2S_JoinRoomReq` in `internal/router/join_room.go` to add a player to a room.
- [x] T017 [US1] Implement the handler for `C2S_PlayerReadyReq` in `internal/router/player_ready.go`.
- [x] T018 [US1] In the `PlayerReady` handler, add logic to check if the game can start (2-5 players ready).
- [x] T019 [US1] Implement the game state machine in `internal/logic/room_fsm.go` with initial states `WaitingForPlayers` and `Dealing`.
- [x] T020 [US1] When the game starts, transition the FSM to `Dealing`, shuffle the deck, and deal 5 cards to each player.
- [x] T021 [US1] Implement broadcasting logic within the `Room` to send `S2C_SyncRoomStateNtf` when players join/ready.
- [x] T022 [US1] Broadcast `S2C_GameStartNtf` and `S2C_DealCardsNtf` to all players in the room when the game starts.

**Checkpoint**: User Story 1 should be fully functional and testable independently.

---

## Phase 4: User Story 2 - 完成一局完整的游戏流程 (Priority: P1)

**Goal**: Players experience the full gameplay loop from bidding to settlement.
**Independent Test**: In a running game, simulate clients sending bid, bet, and showdown messages, and verify that the `S2C_GameResultNtf` contains the correct win/loss calculations.

### Tests for User Story 2 (MANDATORY)

- [x] T023 [P] [US2] Unit test for `RoomFSM.BidBanker` in `internal/logic/room_fsm_test.go`.
- [x] T024 [P] [US2] Unit test for `Deck.Shuffle` and `Deck.Deal` in `internal/logic/deck_test.go`.
- [x] T025 [US2] E2E test for the bidding and betting flow in `tests/e2e/betting_flow_test.go`.

### Implementation for User Story 2

- [x] T026 [US2] Implement the handler for `C2S_BidBankerReq` in `internal/router/bid_banker.go`.
- [x] T027 [US2] Implement the handler for `C2S_PlaceBetReq` in `internal/router/place_bet.go`.
- [x] T028 [US2] Implement the handler for `C2S_ShowdownReq` in `internal/router/showdown.go`.
- [x] T029 [US2] Implement the `CardType` and `BullValue` calculation logic in `internal/logic/bull_logic.go`.
- [x] T030 [US2] In the `Showdown` handler, use `bull_logic.go` to determine winners and calculate score changes.
- [x] T031 [US2] Implement broadcasting logic for `S2C_BankerNtf`, `S2C_BetNtf`, `S2C_ShowdownNtf`, and `S2C_GameResultNtf`.
- [x] T032 [US2] Add `Player.BetAmount` and `Player.HasBet()` methods to track individual bets.
- [x] T033 [US2] Add `Room.HasBanker()` and `Room.GetBankerID()` methods to manage the banker role.

**Checkpoint**: User Stories 1 AND 2 should both work.

---

## Phase 5: User Story 3 - 玩家断线重连 (Priority: P2)

**Goal**: Players who disconnect can rejoin an ongoing game.
**Independent Test**: Disconnect a client mid-game, then reconnect and verify that the server resends the complete current game state to that client only.

### Tests for User Story 3 (MANDATORY)

- [ ] T034 [US1] E2E test for disconnect and reconnect flow in `tests/e2e/reconnect_test.go`.

### Implementation for User Story 3

- [ ] T035 [US3] In the Zinx `OnConnStop` hook, mark the player as disconnected instead of removing them from the room immediately.
- [ ] T036 [US3] When a player sends `C2S_JoinRoomReq`, check if they have a disconnected session in that room.
- [ ] T037 [US3] If they are reconnecting, re-associate their new connection with their existing `Player` object.
- [ ] T038 [US3] Implement a function to generate and send a full room state synchronization message to the reconnected player.
- [ ] T039 [US3] Add a timer to clean up disconnected players from a room if they do not reconnect within a specified time limit (e.g., 5 minutes).

---

## Phase N: Polish & Cross-Cutting Concerns

- [ ] T040 [P] Add detailed logging for all major game events.
- [ ] T041 [P] Implement configuration management (`internal/conf/`) to handle server port, timeouts, etc.
- [ ] T042 Code cleanup and refactoring based on review feedback.
- [ ] T043 Run `go fmt` and `golint` across the entire codebase to ensure compliance with the constitution.