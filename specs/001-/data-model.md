# Data Model: “斗牛牛”棋牌游戏服务器

This document outlines the core data entities for the game server, based on the feature specification.

## Entity: Player

Represents a single user connected to the server and participating in a game.

| Field | Type | Description | Validation Rules |
|---|---|---|---|
| `PlayerID` | `int64` | Unique identifier for the player. | Must be > 0. |
| `Nickname` | `string` | Player's display name. | |
| `Score` | `int64` | Player's current score or currency. | Cannot be negative. |
| `RoomID` | `int32` | The ID of the room the player is currently in. `0` if not in a room. | |
| `Hand` | `[]Card` | The player's 5 cards. | Must contain exactly 5 cards during gameplay. |
| `BetAmount` | `int32` | The amount the player has bet in the current round. | |
| `Status` | `PlayerStatus` | The current status of the player (e.g., `Waiting`, `Ready`, `Playing`). | |
| `IsBanker` | `bool` | `true` if the player is the banker for the current round. | |

### State Transitions (PlayerStatus)

- `Waiting` -> `Ready` (Player clicks ready)
- `Ready` -> `Playing` (Game starts)
- `Playing` -> `Waiting` (Game ends)

## Entity: Room

Represents a single game instance, containing a group of players and the game state.

| Field | Type | Description | Validation Rules |
|---|---|---|---|
| `RoomID` | `int32` | Unique identifier for the room. | Must be > 0. |
| `Players` | `map[int64]*Player` | A map of players currently in the room, keyed by `PlayerID`. | Max 5 players. |
| `GameState` | `GameState` | The current state of the game FSM (e.g., `Dealing`, `Bidding`, `Settling`). | |
| `BankerID` | `int64` | The `PlayerID` of the current banker. | |
| `BaseBet` | `int32` | The base bet amount for the room. | |
| `Deck` | `[]Card` | The deck of cards for the current game. | |

### State Transitions (GameState)

- `WaitingForPlayers` -> `Dealing` (Enough players are ready)
- `Dealing` -> `Bidding` (Cards have been dealt)
- `Bidding` -> `Betting` (Banker has been decided)
- `Betting` -> `Showdown` (All non-banker players have placed bets)
- `Showdown` -> `Settlement` (All players have revealed their hands)
- `Settlement` -> `WaitingForPlayers` (Results have been sent, new round begins)

## Entity: Card

Represents a single playing card.

| Field | Type | Description |
|---|---|---|
| `Suit` | `Suit` | The suit of the card (Spades, Hearts, Clubs, Diamonds). |
| `Rank` | `Rank` | The rank of the card (Ace, 2, 3, ..., King). |

## Enums

### PlayerStatus
- `STATUS_UNKNOWN`
- `WAITING`
- `READY`
- `PLAYING`

### GameState
- `STATE_UNKNOWN`
- `WAITING_FOR_PLAYERS`
- `DEALING`
- `BIDDING`
- `BETTING`
- `SHOWDOWN`
- `SETTLEMENT`

### Suit
- `SUIT_UNKNOWN`
- `SPADES`
- `HEARTS`
- `CLUBS`
- `DIAMONDS`

### Rank
- `RANK_UNKNOWN`
- `ACE`
- `TWO`
- ...
- `KING`