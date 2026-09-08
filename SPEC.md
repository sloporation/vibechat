# VibeChat Spec

## System Operation

VibeChat is organized around **chat channels**. A channel is the unit of membership, message delivery, ordering, and history.

A channel may be:

- **Direct:** exactly two participants.
- **Room:** two or more participants, with the same transaction-log behavior as a direct channel.

Every channel has its own append-only transaction log. The log is not global across the service.

When a client sends a message, the client submits an outgoing request to the server. The server authenticates the sender, authorizes membership, validates the payload, assigns the channel transaction ID, assigns the message ID, persists the transaction, and then delivers the resulting incoming event to channel participants.

The server is authoritative for:

- Channel membership.
- Transaction IDs and channel ordering.
- Message IDs.
- Message timestamps.
- Message persistence.

The client is authoritative only for its local UI state and its outgoing client request ID. A client can safely reconnect by asking for transactions after its last applied transaction ID for each channel.

The same committed transaction can be delivered through the live WebSocket stream and later returned by synchronization. Clients must apply transactions idempotently.

## Transport Model

The application has two communication directions:

- **Outgoing:** client to server. These are POST-style commands: the client submits an action for processing.
- **Incoming:** server to client. These are GET-style results/events: the client receives the server's accepted representation, including server-generated IDs and metadata.

### HTTP

HTTP uses normal request methods:

- `POST` submits an outgoing command.
- `GET` reads channel state or replays committed transactions.

HTTP request bodies contain the outgoing payload. HTTP response bodies contain the incoming payload. The response is not assumed to be identical to the request.

Suggested operations:

```text
POST /channels/{channel_id}/messages
GET  /channels/{channel_id}/messages?after={transaction_id}&limit=100
```

### WebSocket

WebSocket does not have HTTP-style POST and GET messages after the connection is established. It is a bidirectional stream:

- A client-to-server WebSocket frame is equivalent to `POST`.
- A server-to-client WebSocket frame is equivalent to an incoming `GET` result/event.
- A synchronization request sent over WebSocket is an outgoing command; the resulting events are incoming messages.

The payload includes a `d` direction marker so the same envelope is unambiguous across HTTP and WebSocket:

- `out` — client-to-server command.
- `in` — server-to-client response or event.

## Scope

VibeChat is a small direct and room-based messaging application.

The first version supports:

- User accounts.
- Exact username lookup.
- Direct channels and rooms.
- Plain-text messages.
- Persistent per-channel transaction logs.
- Real-time delivery over WebSocket.
- Reconnect and transaction replay.

Files, notifications, federation, bouncers, end-to-end encryption, and native clients are future work.

## Constants

- Every operation uses the payload spec.
- Every payload identifies its direction, action, sender, recipient, and channel where applicable.
- Outgoing and incoming payloads have the same envelope but different payload contents.
- JSON is the application format.
- Binary or opaque values inside a payload are base64 encoded.
- Message text is UTF-8 text, base64 encoded in the wire payload.
- The server validates every payload and authorizes every operation.
- The server assigns channel transaction IDs, message IDs, and server timestamps.

## Payload Spec

All payloads use the same top-level structure:

| Key | Purpose |
| --- | --- |
| `d` | Direction: `out` for outgoing or `in` for incoming. |
| `a` | Numeric action being performed or reported. |
| `s` | Sender identity. |
| `r` | Recipient identity, server identity, or channel identity. |
| `c` | Channel ID. Required for channel operations. |
| `t` | Channel transaction ID. Required on incoming committed channel events; included on outgoing sync cursors. |
| `m` | Server message ID when a message is being received; omitted from a new outgoing message. |
| `p` | Action-specific payload. |

The `t` field is the transaction ID for this channel. It is monotonically increasing within a channel and has no ordering meaning across channels.

The `m` field is the server-assigned message ID. A client may include a client idempotency ID inside `p`, but it must not use that value as the message ID.

### Outgoing envelope

An outgoing message request has no server transaction ID or server message ID yet:

```json
{
  "d": "out",
  "a": 1,
  "s": "alice",
  "r": "chat-001",
  "c": "chat-001",
  "p": {
    "m": "SGVsbG8gQm9i",
    "i": "client-msg-001"
  }
}
```

### Incoming envelope

An incoming committed message includes the channel transaction ID and server message ID:

```json
{
  "d": "in",
  "a": 4,
  "s": "bob",
  "r": "chat-001",
  "c": "chat-001",
  "t": 12,
  "m": "msg-001",
  "p": {
    "m": "SGVsbG8gQWxpY2U=",
    "i": "client-msg-022",
    "created_at": "2026-09-08T00:00:00.000Z"
  }
}
```

The sender and recipient are application identities, not display names. The server resolves them to internal IDs. Clients must treat IDs as opaque strings.

## Actions

| Action | Name | Direction | Purpose |
| ---: | --- | --- | --- |
| `0` | `key_exchange` | out/in | Reserved for future encryption work. |
| `1` | `message.create` | out | Submit a new message. |
| `2` | `message.accepted` | in | Confirm that the server persisted or recognized a message. |
| `3` | `channel.sync` | out | Request channel events after a transaction ID. |
| `4` | `message.created` | in | Deliver a committed message transaction. |
| `5` | `error` | in | Report a rejected operation. |
| `6` | `channel.event` | in | Deliver a future non-message channel transaction. |

The MVP requires `message.create`, `message.accepted`, `channel.sync`, `message.created`, and `error`. `key_exchange` is reserved and is not an MVP encryption implementation.

## Outgoing Message

A new message uses action `1` and direction `out`.

| Key | Purpose |
| --- | --- |
| `m` | Base64-encoded UTF-8 message text. |
| `i` | Client-generated idempotency ID. |

```json
{
  "d": "out",
  "a": 1,
  "s": "alice",
  "r": "chat-001",
  "c": "chat-001",
  "p": {
    "m": "SGVsbG8gQm9i",
    "i": "client-msg-001"
  }
}
```

Rules:

- `c` identifies the channel receiving the message.
- `r` identifies the channel for channel operations; the server determines the channel participants.
- `t` and the top-level server message ID `m` are omitted from a new outgoing message.
- `p.m` must not be empty after decoding.
- The decoded message must be no longer than 4,000 characters.
- `p.i` is unique per sender and channel. Retrying the same `p.i` must not create a duplicate.
- The sender must be a member of the channel.
- The server assigns the durable message ID, channel transaction ID, and timestamp.

## Incoming Accepted Message

A successful persistence response uses action `2` and direction `in`.

```json
{
  "d": "in",
  "a": 2,
  "s": "server",
  "r": "alice",
  "c": "chat-001",
  "t": 12,
  "m": "msg-001",
  "p": {
    "i": "client-msg-001",
    "status": "sent",
    "created_at": "2026-09-08T00:00:00.000Z"
  }
}
```

| Field | Purpose |
| --- | --- |
| `c` | Channel ID. |
| `t` | Channel transaction ID. |
| `m` | Server-assigned message ID. |
| `p.i` | Original client idempotency ID. |
| `p.status` | `sent` when newly persisted, or `duplicate` when returning an existing message. |
| `p.created_at` | Server-created UTC timestamp. |

The acknowledgement means the server committed the message to the channel transaction log. It does not mean that any recipient has read it.

## Incoming Message Event

A committed message delivered to channel participants uses action `4` and direction `in`.

```json
{
  "d": "in",
  "a": 4,
  "s": "bob",
  "r": "chat-001",
  "c": "chat-001",
  "t": 12,
  "m": "msg-001",
  "p": {
    "m": "SGVsbG8gQWxpY2U=",
    "i": "client-msg-022",
    "created_at": "2026-09-08T00:00:00.000Z"
  }
}
```

The top-level `m` is the server message ID. The `p.m` value is the message body encoded as base64. This intentionally distinguishes message identity from message content.

Clients must upsert by top-level `m` and must mark `(c, t)` as applied. The same event may be received more than once through WebSocket and synchronization.

## Per-Channel Transaction Log

Each channel has an independent, append-only transaction sequence.

```text
chat-001: 1, 2, 3, 4, ...
room-001: 1, 2, 3, 4, ...
chat-002: 1, 2, 3, 4, ...
```

A transaction contains:

| Field | Purpose |
| --- | --- |
| `channel_id` | Channel that owns the transaction. |
| `t` | Monotonically increasing transaction ID within that channel. |
| `id` | Unique transaction ID, represented by top-level `t` together with `c`. |
| `a` | Action represented by the transaction. |
| `s` | Sender identity. |
| `r` | Channel or recipient identity. |
| `m` | Server message ID when the transaction represents a message. |
| `p` | Event payload. |
| `created_at` | Server-created UTC timestamp. |

The transaction log is the source for replay and reconnect. A client stores the last applied transaction for each channel:

```text
chat-001 -> t=11
```

After reconnect, the client requests transactions after `t=11` for `chat-001`. The server returns transaction `12` onward in order. The log is scoped to the channel; there is no single global cursor in the MVP.

## Channel Synchronization

A client requests missed transactions with action `3` and direction `out`.

```json
{
  "d": "out",
  "a": 3,
  "s": "alice",
  "r": "server",
  "c": "chat-001",
  "t": 11,
  "p": {
    "limit": 100
  }
}
```

Here, `t` means “return transactions after this channel transaction ID.” The server responds with incoming action `4` messages, in channel transaction order. If the cursor is absent, the server starts from the oldest available transaction or the server's configured history boundary.

Synchronization over HTTP is a `GET`:

```text
GET /channels/chat-001/messages?after=11&limit=100
```

Synchronization over WebSocket is an outgoing `channel.sync` frame followed by incoming events. The semantic operation is the same; only the transport differs.

## Error

A rejected operation uses action `5` and direction `in`.

```json
{
  "d": "in",
  "a": 5,
  "s": "server",
  "r": "alice",
  "c": "chat-001",
  "p": {
    "code": "INVALID_MESSAGE",
    "message": "Message text must not be empty.",
    "request_id": "client-msg-001"
  }
}
```

Initial error codes:

- `INVALID_PAYLOAD`
- `UNAUTHENTICATED`
- `FORBIDDEN`
- `NOT_FOUND`
- `INVALID_MESSAGE`
- `DUPLICATE_REQUEST`
- `RATE_LIMITED`
- `INTERNAL_ERROR`

Error messages must not include passwords, session credentials, or unnecessary private data.

## MVP Data Model

The minimum persisted model is:

### User

- `id`
- `username` (unique)
- `email` (unique)
- `password_hash`
- `created_at`

### Channel

- `id`
- `type` (`direct` or `room`)
- `created_at`
- `updated_at`

A direct channel is unique for a pair of users. A room has two or more participants.

### Channel Participant

- `channel_id`
- `user_id`
- `created_at`

### Channel Transaction

- `channel_id`
- `t` (unique within `channel_id`)
- `action`
- `sender_id`
- `recipient_id` or channel recipient
- `message_id` when applicable
- `payload` (JSON)
- `created_at`

The database must enforce uniqueness for `(channel_id, t)` and for the sender/channel/client idempotency pair. Message IDs must be globally unique.

## MVP Acceptance Criteria

- Two users can register and sign in.
- One user can find the other by exact username.
- They can open one unique direct channel.
- A room can use the same channel and transaction-log model when room support is enabled.
- A user can send a non-empty text message using an outgoing POST-style operation.
- The server returns an incoming response containing the server message ID and channel transaction ID.
- The message is persisted in that channel's transaction log.
- Channel participants receive the incoming message event over WebSocket without refreshing.
- A reconnecting client can request events after its last per-channel transaction ID.
- Repeating a message request does not create a duplicate transaction or message.
- Unauthorized users cannot read or write the channel.
- Messages are limited to 4,000 decoded characters.
- Failed operations return a structured incoming error payload.
