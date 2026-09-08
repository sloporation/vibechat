# Chat MVP Sub-Spec

## 1. Purpose

This is the deliberately reduced implementation spec for the first usable version of the chat app. `SPEC.md` remains the broader product and architecture reference. When the two conflict, this file defines the MVP scope unless a decision is explicitly recorded elsewhere.

The goal is to get two people reliably exchanging text messages in a browser with minimal infrastructure and few unresolved design decisions.

## 2. MVP outcome

A new user can:

1. Create an account and sign in.
2. Find another user by exact username.
3. Start or reopen a direct conversation.
4. Send and receive text messages without manually refreshing.
5. See conversation history after closing and reopening the app.
6. Sign out and sign back in.

The MVP should be usable by a small private test group. It does not need to be federation-ready, mobile-native, end-to-end encrypted, or production-scale.

## 3. Explicitly out of scope

Defer all of the following until the basic direct-message loop is stable:

- End-to-end encryption and device-key management.
- Relays, federation, bouncers, and opaque cross-server envelopes.
- Native mobile clients.
- Group conversations.
- File uploads, images, audio, video, and links previews.
- Message editing, deletion, reactions, replies, and threads.
- Typing indicators and presence.
- Push, email, and browser notifications.
- Blocking, reporting, moderation tooling, and account deletion workflows.
- Advanced search, public profiles, channels, bots, and integrations.
- Multi-device synchronization and offline outbound queues.
- Formal delivery/read receipts.
- High-availability, multi-region, and large-scale deployment work.

TLS must still be used for deployed environments. Deferring end-to-end encryption is an MVP tradeoff, not a claim that the product is private by default.

## 4. Product decisions for the MVP

- **Client:** One responsive web client for desktop and mobile-sized screens.
- **Server:** One application server providing the web app, HTTP API, and real-time connection.
- **Database:** SQLite for local development and the first private deployment, unless the implementation environment already requires PostgreSQL.
- **Identity:** Unique username, email address, and password.
- **Conversation type:** Direct conversations between exactly two users.
- **Message type:** Plain text only.
- **Maximum message length:** 4,000 characters.
- **Real-time transport:** WebSocket; reconnect and refresh are fallback behavior.
- **History:** Server-persisted message history, paginated when needed.
- **Authentication:** Cookie-based session or an equivalent server-managed session. Passwords use an adaptive password hash and are never stored in plaintext.

## 5. Required user flows

### 5.1 Account

- Register with username, email, and password.
- Validate required fields, email format, username uniqueness, and minimum password strength.
- Sign in with email and password.
- Remain signed in across normal browser restarts.
- Sign out the current session.
- Show generic authentication errors that do not disclose whether an account exists.

Password recovery is not required for the first private MVP unless the chosen deployment needs it before testing.

### 5.2 Find and start a conversation

- Search by exact username.
- Show only the matching user's public username and display name.
- Open the existing direct conversation if one already exists.
- Create a direct conversation otherwise.
- Do not allow a user to start a conversation with themselves.

### 5.3 Messaging

- Display messages in chronological order.
- Show author, message text, and timestamp.
- Reject empty or whitespace-only messages.
- Enforce the 4,000-character limit on both client and server.
- Send a message over the API and update the UI without a full-page refresh.
- Deliver new messages over WebSocket to the other participant when connected.
- Persist every accepted message in the database.
- Reconnect the WebSocket after temporary failure and refresh conversation state safely.
- Do not create duplicates when a client retries an operation.

For this MVP, the visible states are `sending`, `sent`, and `failed`. `sent` means the server accepted and persisted the message; it does not mean the recipient read it.

## 6. Minimum screens

- Sign up.
- Sign in.
- Conversation list.
- Direct conversation view with message history and composer.
- New conversation / exact username search.
- Basic account menu with display name and sign-out.

Each network-backed screen must have loading, empty, error, and retry states where applicable. Failed sends must remain visible and provide retry or discard.

## 7. Minimal data model

### User

- `id`
- `username` (unique)
- `email` (unique)
- `password_hash`
- `display_name`
- `created_at`
- `updated_at`

### Session

- `id`
- `user_id`
- `expires_at`
- `created_at`
- `last_used_at`

### Conversation

- `id`
- `created_at`
- `updated_at`

A direct conversation must be unique for a pair of users.

### ConversationParticipant

- `conversation_id`
- `user_id`
- `created_at`

### Message

- `id` (server-generated unique ID)
- `conversation_id`
- `author_id`
- `body`
- `created_at`

The server is authoritative for MVP message history. No client-side database or bouncer is required.

## 8. Minimum API and real-time behavior

The exact framework and URL naming are implementation choices, but the server must provide equivalent operations:

- Register, sign in, sign out, and get current session.
- Search users by exact username.
- List the current user's conversations.
- Create or open a direct conversation.
- List messages for a conversation with cursor-based pagination or a simple bounded initial page.
- Send a message idempotently.
- Open an authenticated WebSocket connection.

WebSocket events should minimally cover:

- `message.created`
- `conversation.updated` or an equivalent conversation-list update
- connection/authentication failure

Every protected operation must verify that the authenticated user belongs to the relevant conversation. The client must never be trusted for authorization.

## 9. Security baseline

- Use TLS outside local development.
- Hash passwords with a current adaptive password-hashing algorithm.
- Use secure, HttpOnly, SameSite session cookies if cookies are used.
- Validate and length-limit all input server-side.
- Escape user-provided text at render boundaries.
- Rate-limit login, registration, username search, and message sending enough for a private deployment.
- Do not log passwords, session credentials, or message bodies.
- Return generic errors for authentication and authorization failures.

## 10. Acceptance criteria

The MVP is acceptable when all of the following work in a clean test environment:

- Two users can register and sign in.
- User A can find User B by exact username.
- User A can create a direct conversation and send a message.
- User B sees the message without manually refreshing while both users are connected.
- Either user can close and reopen the app and still see persisted history.
- A disconnected client can reconnect and recover the current conversation state.
- Empty, oversized, unauthenticated, and unauthorized requests are rejected.
- Retrying a send does not create duplicate messages.
- Both users can sign out and sign back in successfully.
- Automated tests cover authentication, authorization, conversation uniqueness, message persistence, message validation, and the real-time send/receive path.
- A manual smoke test passes on a desktop viewport and a narrow mobile-sized viewport.

## 11. Standard JSON payload

All MVP HTTP JSON responses and WebSocket messages use the same top-level envelope. HTTP requests may use a direct request body, but HTTP responses and every WebSocket message use this envelope.

```json
{
  "version": 1,
  "id": "01J...",
  "type": "message.created",
  "timestamp": "2026-09-08T00:00:00.000Z",
  "request_id": "01J...",
  "payload": {},
  "error": null
}
```

Fields:

- `version`: integer protocol version; MVP value is `1`.
- `id`: unique opaque string generated by the sender.
- `type`: lowercase dot-separated message name, such as `message.create` or `message.created`.
- `timestamp`: UTC RFC 3339 timestamp with a `Z` suffix.
- `request_id`: optional ID of the request that caused this response or event.
- `payload`: operation-specific object; use `{}` when there is no data.
- `error`: `null` on success, otherwise an error object.

Unknown fields must be ignored. Unsupported versions must be rejected rather than processed. IDs are opaque and must not be parsed by clients. JSON is UTF-8 and must not be base64-encoded.

Requests use operation types such as `auth.register`, `auth.sign_in`, `user.search`, `conversation.open`, `message.list`, and `message.create`. Responses use the corresponding `*.response` type. Server-pushed events use `message.created` and `conversation.updated`.

A message-create request includes a client-generated idempotency value:

```json
{
  "version": 1,
  "id": "req_01",
  "type": "message.create",
  "timestamp": "2026-09-08T00:00:00.000Z",
  "payload": {
    "conversation_id": "conv_01",
    "body": "Hello",
    "client_message_id": "client_msg_01"
  },
  "error": null
}
```

The server validates this value together with the authenticated user. Retrying the same request must return the original result and must not create a duplicate message. Clients must upsert received messages by server-assigned message ID.

Errors have stable machine-readable codes:

```json
{
  "version": 1,
  "id": "res_01",
  "type": "message.create.response",
  "timestamp": "2026-09-08T00:00:00.000Z",
  "request_id": "req_01",
  "payload": {},
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Message body must not be empty.",
    "fields": { "body": "must not be empty" }
  }
}
```

Initial error codes are `INVALID_PAYLOAD`, `VALIDATION_ERROR`, `UNAUTHENTICATED`, `FORBIDDEN`, `NOT_FOUND`, `CONFLICT`, `RATE_LIMITED`, and `INTERNAL_ERROR`. Error messages must not expose secrets or unnecessary private data.

## 12. Post-MVP sequence

After this spec passes, consider features in this order:

1. Better account recovery and profile settings.
2. Message edits and deletion.
3. Read state and unread counts.
4. Basic groups.
5. Notifications.
6. Blocking and reporting.
7. Local client history and offline behavior.
8. End-to-end encryption and the relay/bouncer architecture from `SPEC.md`.
