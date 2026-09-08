# Messaging App Specification

## 1. Product Summary

The product is a private, real-time messaging app for people who want to have direct conversations and small group conversations from a web or mobile client.

The first release should make the core loop reliable and understandable:

1. A person creates an account and signs in.
2. They find another person or accept an invitation.
3. They start a one-to-one conversation or create a small group.
4. Participants send and receive text messages in real time.
5. The client preserves message history and clearly communicates delivery state.

The MVP is intentionally narrower than a full social or collaboration platform. It will optimize for dependable conversations, not feature breadth.

## 2. Goals

### 2.1 MVP goals

- Support account creation, sign-in, sign-out, and session recovery.
- Allow users to find people by an approved identifier and begin a conversation.
- Support one-to-one conversations.
- Support small group conversations.
- Send, receive, edit, and delete text messages.
- Deliver new messages in real time when participants are online.
- Preserve message history locally on the client. When a bouncer is configured, it also stores encrypted history for the user's clients.
- Show useful message states: sending, sent, delivered, and read.
- Notify users about new messages when they are not actively viewing the conversation.
- Provide basic abuse and privacy controls.

### 2.2 Non-goals for the MVP

- Manual key management and recovery workflows beyond what clients handle automatically.
- Voice or video calls.
- Stories, channels, public communities, or a social feed.
- Bots, integrations, or an app marketplace.
- Payments or monetization.
- Advanced message search.
- Large-file storage or media processing.
- Multi-device conflict resolution beyond a normal signed-in session.

These may be considered after the core messaging experience is stable.

## 3. Users and Terminology

- **User:** An authenticated person with an account.
- **Conversation:** A persistent container for messages and participants.
- **Direct conversation:** A conversation with exactly two participants.
- **Group conversation:** A conversation with three or more participants.
- **Participant:** A user who belongs to a conversation.
- **Message:** A user-authored item in a conversation. The MVP supports text only.
- **Read cursor:** The point through which a participant has read messages in a conversation.
- **Online:** The server has a current connection or recent activity for the user. This is an approximate presence state, not a guarantee that the user is looking at the app.

## 4. User Stories

### 4.1 Account and identity

- As a new user, I can create an account with a unique public username, email address, and password.
- As a returning user, I can sign in and remain signed in across normal app restarts.
- As a user, I can sign out from the current session.
- As a user, I can change my display name, avatar, and basic profile settings.
- As a user, I can control whether other users can find me by username.

### 4.2 Starting conversations

- As a user, I can search for another user by exact username.
- As a user, I can start a direct conversation with a user who has not blocked me.
- As a user, I can create a group conversation by selecting users and assigning a group name.
- As a user, I can invite another user to a conversation when I have permission to do so.
- As a user, I can leave a group conversation.

### 4.3 Messaging

- As a participant, I can send a text message to the conversation.
- As a participant, I see new messages without manually refreshing the conversation.
- As a participant, I can edit my own message within the allowed edit window.
- As a participant, I can delete my own message and see that it was deleted.
- As a participant, I can see when a message is sending, sent, delivered, or read.
- As a participant, I can see when another participant is typing, with a privacy-preserving indicator.

### 4.4 Notifications and safety

- As a user, I receive an in-app notification for a new message outside the active conversation.
- As a user, I can mute a conversation.
- As a user, I can block another user.
- As a user, I can report a user or message.
- As a user, I can delete my account and request deletion of associated personal data.

## 5. Product Requirements

### 5.1 Authentication and sessions

- Registration must validate email format, password strength, and username uniqueness.
- Passwords must never be stored in plaintext.
- A user may have multiple active sessions unless a later security decision changes this.
- Sign-out must invalidate the current session.
- Authentication failures should not reveal whether an email or username exists.
- Account recovery must use a time-limited, single-use recovery mechanism.

### 5.2 Profiles and discovery

- Each user has a stable internal ID and a unique public username.
- A profile contains a display name and optional avatar.
- Username search returns exact matches in the MVP; prefix or full-text search is deferred.
- Search results must not reveal private account data.
- A blocked user must not be able to start a new conversation with the blocker or see the blocker's presence.

### 5.3 Conversations

- A direct conversation is unique for a pair of users. Starting it again opens the existing conversation.
- A group conversation has an owner or administrator role and a configurable participant limit.
- Conversation ordering is based on the latest locally known message or relevant activity.
- The client must retain enough membership history to explain system events such as joining or leaving.
- Users who leave a group can no longer send messages, but retain access to history only if that policy is explicitly selected. The initial default should be no further access to new messages.
- Muting suppresses notifications but does not remove messages or change read state.

### 5.4 Messages

- MVP messages contain text, author, conversation, creation time, and lifecycle state. The client stores this data locally.
- Message text must have a maximum length defined before implementation and enforced by the relay and client.
- Empty or whitespace-only messages are invalid.
- Messages are ordered by a timestamp and stable message ID. Clients must tolerate clock differences and out-of-order delivery.
- The sender creates the message ID. Relays use it for deduplication and must not replace it.
- Retrying a request must not create duplicate messages.
- Editing a message updates its content and records that it was edited.
- Deleting a message replaces its visible content with a deletion marker rather than silently removing the timeline entry.
- An edit or deletion must be rejected if the user is no longer a participant or lacks permission.
- The MVP should define an edit/delete time window, with the initial recommendation of 15 minutes for edits and 15 minutes for user deletion.

### 5.5 Delivery, read state, and presence

- **Sending:** The client has submitted the message, but the server has not acknowledged it.
- **Sent:** A relay accepted the encrypted message for delivery.
- **Delivered:** The message has been delivered to at least one active recipient session. For a group, the UI should avoid implying that every member received it unless per-user delivery is implemented.
- **Read:** A participant has advanced their read cursor past the message.
- Read state is per participant and must not require a separate read receipt for every message.
- Typing indicators are ephemeral and must expire automatically.
- Presence is best effort and must not be treated as proof of availability.

### 5.6 Real-time behavior and offline recovery

- The client maintains a real-time connection when network conditions allow.
- Reconnect attempts use bounded backoff and do not lose locally composed unsent text.
- On reconnect, the client requests all conversation events after its last known cursor.
- Events must be safe to receive more than once and safe to apply out of order where practical.
- If real-time delivery fails, the client must still display local history and retry delivery when the connection returns.
- A message that cannot be sent must be visibly marked as failed and offer retry or discard.

### 5.7 Notifications

- In-app notifications are required for new messages in inactive conversations.
- A conversation must not notify the user while it is muted.
- The active conversation should not produce duplicate in-app notifications for messages already visible to the user.
- Push or email notifications are out of scope until the target platforms and permission model are selected.

### 5.8 Moderation and privacy

- Users can block and unblock other users.
- Users can report a message or account with a reason and optional description.
- Reports must be persisted with enough context for investigation, subject to retention policy.
- Authorization must be checked on every conversation and message operation server-side.
- The app must minimize exposure of email addresses and other private profile fields.
- Sensitive actions should be auditable without storing message content in application logs.

## 6. Core Screens and States

The first client should include these areas:

- **Sign up / sign in:** validation, recovery entry point, and error states.
- **Conversation list:** unread counts, latest message preview, mute state, loading, empty, and error states.
- **Conversation view:** message history, composer, send failure state, typing indicator, read state, and pagination/loading history.
- **New conversation:** exact username search and group participant selection.
- **Profile/settings:** identity settings, session controls, notification preferences, and account deletion.
- **Conversation details:** participants, mute, leave, block/report actions, and group administration where applicable.

Every network-backed view must define loading, empty, failure, retry, and unauthorized states before implementation.

## 7. Initial Data Model

The logical model should include at least:

- `User`: ID, username, email, display name, avatar reference, status, created time, updated time.
- `Session`: ID, user ID, token or credential reference, device metadata, created time, last-used time, expiry/revocation state.
- `DeviceKey`: device ID, user ID, public key, key type, created time, status, and revocation state. Private keys remain on clients.
- `Conversation`: ID, type, name, owner/admin reference, created time, updated time, archived/deleted state.
- `ConversationParticipant`: conversation ID, user ID, role, joined time, left time, mute state, read cursor, last active time.
- `Message`: ID, conversation ID, author ID, encrypted content, content type, created time, edited time, deleted time.
- `MessageReceipt` or equivalent read cursor representation if per-user delivery/read detail is required.
- `Block`: blocker ID, blocked ID, created time.
- `Report`: reporter ID, subject type and ID, reason, description, status, created time, resolved time.

Message history is not stored by the relay. Each client stores its own local history. If a user uses a bouncer, the bouncer stores the user's encrypted message envelopes and may serve as the history source for that user's clients. The relay may retain only transient delivery queues, subject to queue policy.

## 8. Relay Protocol

The protocol connects clients, relay servers, and optional bouncers. A relay forwards opaque, encrypted application payloads; it must not need to decrypt message or file content. “Relay” is the provisional name for the user-owned home server described in `IDEA.md`.

### 8.1 Transport and envelope

- All connections use an encrypted transport.
- Every application frame is JSON encoded, then base64 encoded for transport.
- The decoded JSON is a small, readable envelope. `type` identifies the operation; `id` correlates a request and response; `payload` contains operation data.
- `sender` and `recipient` use a provisional federated identifier format such as `@user#relay.example`. The final identifier syntax remains open.
- Encrypted application data is represented as base64 in `payload.ciphertext`; the relay routes it without inspecting it.

```json
{
  "version": 1,
  "type": "message.send",
  "id": "req_01",
  "sender": "@alice#relay.example",
  "recipient": "@bob#other.example",
  "payload": {
    "message_id": "msg_01",
    "conversation_id": "conv_01",
    "ciphertext": "BASE64_ENCRYPTED_JSON",
    "content_type": "text/plain"
  }
}
```

Responses use the same envelope and reference the request with `in_reply_to`:

```json
{
  "version": 1,
  "type": "message.accepted",
  "id": "evt_01",
  "in_reply_to": "req_01",
  "payload": { "message_id": "msg_01", "status": "queued" }
}
```

Malformed, unsupported, unauthorized, or expired requests receive an error envelope. Error codes are stable machine-readable strings; human-readable text is supplementary.

```json
{
  "version": 1,
  "type": "error",
  "id": "err_01",
  "in_reply_to": "req_01",
  "payload": { "code": "UNAUTHORIZED", "message": "Authentication required." }
}
```

### 8.2 Client connects and authenticates

1. The client opens an encrypted connection to its relay.
2. The client sends `auth.request` with its stable user ID, public key, protocol version, and proof of possession of its private key or session credential.
3. The relay validates the credential and responds with `auth.accepted` or `error`.
4. The client sends `presence.set` with `online` and an optional device/session ID.
5. The relay responds with `presence.accepted` and begins routing messages for that session.

```json
{
  "version": 1,
  "type": "auth.request",
  "id": "req_auth_01",
  "sender": "@alice#relay.example",
  "payload": {
    "public_key": "BASE64_PUBLIC_KEY",
    "credential": "SESSION_OR_SIGNATURE",
    "client": "desktop",
    "protocol": 1
  }
}
```

### 8.3 Relay queues and delivers messages

1. A sender encrypts the application message for the recipient and stores the message locally. For direct messages, the current draft uses the recipient’s public key; the eventual wire format may use that key to establish a temporary symmetric session key rather than encrypting every message directly.
2. The sender sends `message.send` to the sender’s relay.
3. The sender’s relay validates the envelope, authorization, destination, and size without decrypting `ciphertext`.
4. If the destination is connected, the relay forwards the envelope to the destination relay or client.
5. If the destination is offline, the relay queues the opaque envelope for later delivery. Queue retention and limits are open decisions.
6. The destination client sends `message.received` after accepting the envelope.
7. The destination relay returns delivery status to the sender’s relay, which forwards `message.delivered` when supported.

```json
{
  "version": 1,
  "type": "message.received",
  "id": "evt_recv_01",
  "in_reply_to": "req_01",
  "payload": { "message_id": "msg_01", "device_id": "phone_01" }
}
```

The relay must deduplicate by `(sender, message_id)` and preserve delivery ordering where possible. Acknowledging a message means the relay accepted it; it does not mean the recipient stored or read it.

### 8.4 Client receives queued messages

1. After `presence.set`, the relay sends `queue.available` with a count or cursor, not message content.
2. The client sends `queue.request` with its last acknowledged cursor.
3. The relay sends queued encrypted envelopes as `message.deliver` events.
4. The client decrypts and validates each envelope, then sends `message.received`.
5. The relay removes or marks each item delivered only after acknowledgement. The relay queue is not message history; retry and retention policy is open.

```json
{
  "version": 1,
  "type": "queue.request",
  "id": "req_queue_01",
  "sender": "@alice#relay.example",
  "payload": { "after": "queue_cursor_42", "limit": 100 }
}
```

### 8.5 Presence and typing

1. A connected client sends `presence.set` on connect and when its state changes.
2. The relay publishes `presence.changed` to authorized interested clients.
3. A client sends `typing.start` or `typing.stop` for a conversation.
4. The relay forwards typing events without persisting them; events expire automatically.
5. The client sends `presence.set` with `offline` before an intentional disconnect. The relay also marks the session offline on timeout.

Presence is best effort and must not be treated as proof that a user is available or viewing a conversation.

### 8.6 Read receipts

1. The client advances its local read cursor when the user views messages.
2. The client sends `conversation.read` with the conversation and highest read message/sequence.
3. The relay forwards the encrypted or opaque receipt to relevant participants’ relays.
4. Recipients apply the receipt and may display `read` state.

```json
{
  "version": 1,
  "type": "conversation.read",
  "id": "req_read_01",
  "sender": "@bob#other.example",
  "payload": { "conversation_id": "conv_01", "through": "msg_01" }
}
```

### 8.7 Bouncer proxy

The bouncer is an always-on user agent that acts as the user’s connected client and stores encrypted envelopes and local history for that user. It does not decrypt message or file content. A client must know when it is connected through a bouncer so it can request history, upload local records, and obtain device keys through the bouncer protocol.

1. The bouncer authenticates to the user’s relay as the user’s account representative. The relay needs to authorize the bouncer, not each client behind it.
2. The relay marks the user online through the bouncer and routes messages to it.
3. The bouncer stores received encrypted envelopes and acknowledges them to the relay.
4. A phone or desktop client connects to the bouncer and authenticates to the bouncer. Bouncer authentication may be separate from relay authentication; upstream relay authentication is handled by the bouncer.
5. The bouncer identifies itself and exposes its history and key-provisioning capabilities to the client.
6. The client requests stored envelopes, local history, current cursors, and any keys needed to decrypt them.
7. The bouncer forwards encrypted records and encrypted key material to the client. The client decrypts content locally.
8. The bouncer may proxy outbound `message.send` requests from the client through the relay.

The bouncer is therefore a protocol participant, not merely a cache. It provides minimal multi-client history handoff, not a general synchronization service. Its storage format, retention, delegation, revocation, and recovery behavior require separate decisions.

Each client generates its own device key. There is no primary client. When a user has a bouncer, each client uses a per-user history cursor and device key identity:

1. The client sends newly accepted encrypted message records and local state changes to the bouncer.
2. The bouncer stores records by message/event ID and acknowledges the highest stored cursor.
3. A newly connected client requests records after its local cursor.
4. The bouncer returns missing encrypted records in cursor order.
5. The client applies records idempotently and acknowledges its new cursor.

This handoff covers missing history and delivery state only. It does not merge conflicting edits, reconcile independent local deletions, or make the bouncer authoritative over a client’s local presentation.

### 8.8 Multi-device keys

The bouncer distributes encrypted key material but does not store private keys. It may relay private keys between authenticated user devices when those devices explicitly authorize the transfer.

1. A new client generates its device key and authenticates to the bouncer.
2. The bouncer discovers other authenticated clients belonging to the user that hold the required keys.
3. The bouncer requests approval from an eligible client, using a user-visible challenge such as matching a code or emoji on both devices.
4. The approving client negotiates a temporary encrypted channel with the new client.
5. The approving client sends the required private or conversation keys through that channel. The bouncer may relay the transfer but does not store the keys.
6. The new client proves possession of its private key and decrypts the transferred key material locally.
7. If no authenticated client holds the required keys, the related conversations are lost unless the user has made a separate backup.
8. Key re-rolling is an explicit user action, not an automatic recovery behavior.

The bouncer may authorize a new client and is the only device the relay needs to recognize for a bouncer-backed user. Clients normally do not connect directly to the relay when using a bouncer. Whether this authority boundary is required for all deployments remains an open design decision. The bouncer must not be able to decrypt messages solely because it relays key material.

### 8.9 Encrypted files

Files use the same message envelope but are transferred separately from the real-time channel.

1. The client generates a file key and encrypts the file locally.
2. The client uploads the encrypted bytes to an authorized upload endpoint or relay-managed object store.
3. The upload service returns an opaque `file_id`, size, hash, and download reference.
4. The client encrypts a file descriptor containing that reference, key, content type, and filename for the recipient(s).
5. The client sends the descriptor as `message.send` with `content_type` such as `image/jpeg` or `application/octet-stream`.
6. The recipient decrypts the descriptor, downloads the encrypted bytes, verifies the hash, and decrypts locally.
7. Uploads support resumable chunks and expiration. Chunk size, maximum file size, retention, and malware handling are open decisions.

The relay must never require plaintext file content. File metadata may still reveal size, timing, or routing information unless the deployment adds metadata protection.

### 8.10 Reconnect and replay

1. The client reconnects after transport failure using bounded backoff.
2. The client re-authenticates to its bouncer, or directly to the relay when no bouncer is configured, and sends its last received cursor for each stream.
3. The relay replays missed events and queued messages, safely deduplicated by event/message ID. The bouncer may replay locally stored envelopes or history records.
4. The client acknowledges each accepted event and refreshes state if a cursor is no longer available.
5. The client marks unacknowledged outbound messages as `failed` or `unknown` and offers retry without creating duplicates.

All events must be idempotent. A client must tolerate duplicate delivery, and a relay must tolerate duplicate acknowledgement.

### 8.11 Protocol decisions still open

- Final name and authority model for user-owned relays.
- Federated identifier syntax: `@user#relay`, `user@relay`, or another form.
- Whether relays may persist queued envelopes, and for how long.
- How a bouncer provisions encrypted keys across devices without accessing plaintext private keys.
- Key exchange, rotation, revocation, and group-key management.
- Relay-to-relay authentication and discovery.
- Maximum envelope, message, and file sizes.
- Upload authorization, resumable upload protocol, retention, and abuse scanning.
- Whether edits, deletes, membership changes, and moderation reports are encrypted application events or relay-visible control events.
- Whether the relay recognizes only the bouncer or also individual clients for bouncer-backed users.

## 9. Control-Plane API and Event Expectations

The relay protocol above handles delivery. A separate control plane may handle account, directory, conversation membership, moderation, and device/key management. The implementation may choose REST, RPC, GraphQL, or another interface, but it must support equivalent operations:

- Register, sign in, sign out, recover account, and manage sessions.
- Read and update the current user's profile and settings.
- Search users by exact username.
- List conversations with pagination and unread/read state.
- Create a direct or group conversation.
- Add, remove, or leave participants according to role rules.
- List local message history from the client or bouncer, if that storage is exposed through a protocol API.
- Send, edit, delete, and retry messages safely.
- Mark a conversation read.
- Block, unblock, and report.

Control-plane events and relay events should include enough information to apply changes idempotently. The initial event vocabulary should cover:

- Conversation created or updated.
- Participant added, removed, or left.
- Message created, edited, or deleted.
- Read cursor advanced.
- Typing started or stopped.
- Presence changed.

The relay is not authoritative for message history. The client is authoritative for its local presentation; when configured, the bouncer is the shared encrypted history source and key-provisioning service for that user's clients. Clients using the same bouncer must be able to retrieve missing records by cursor. Conflict resolution remains out of scope.

## 10. Security and Reliability Requirements

- Enforce authentication and authorization on the server for every protected operation.
- Use TLS in transit and encryption at rest provided by the selected infrastructure.
- Hash passwords with a current adaptive password hashing algorithm.
- Apply rate limits to sign-in, registration, username search, message sending, and report creation.
- Validate and encode user-provided text at all display boundaries.
- Protect against duplicate submissions, replayed requests, and unauthorized object access.
- Do not put access tokens, passwords, message content, or recovery secrets in logs.
- Clients and bouncers should document backup and restore behavior for local history before production use.
- Define retention and deletion behavior for accounts, messages, reports, logs, and backups.

## 11. Accessibility and UX Requirements

- The client must be usable with keyboard navigation.
- Interactive controls must have accessible names and visible focus states.
- Messages and status changes must be understandable to screen readers.
- Color must not be the only signal for unread, failed, or selected states.
- Layouts must work on narrow mobile screens and larger desktop screens.
- Dates and times must be displayed in the user's locale while retaining an unambiguous underlying timestamp.
- The composer must preserve draft text during temporary connection failures.

## 12. Observability and Success Metrics

The system should measure reliability without collecting unnecessary message content:

- Registration and sign-in success/error rates.
- Time to open a conversation.
- Message send success rate and acknowledgement latency.
- Real-time reconnect rate and event replay failures.
- Notification delivery and read rates where supported.
- API error rates by operation and status class.
- Report volume and moderation resolution time.

MVP quality targets should be set before launch. At minimum, the team should agree on acceptable message loss tolerance, acknowledgement latency, and service availability.

## 13. Acceptance Criteria for MVP

The MVP is ready for a controlled release when:

- Two test users can register, sign in, find each other, and start a direct conversation.
- A message sent by one user appears for the other without a manual refresh when both are connected.
- Messages remain available from local client history after sign-out and back-in. A client using a bouncer can retrieve history through the bouncer.
- Temporary disconnection does not silently lose acknowledged messages or composed drafts.
- Duplicate retries do not create duplicate messages.
- Users can edit and delete their own eligible messages, and other users see the correct state.
- Read state and unread counts remain consistent across refresh and reconnect.
- Group creation, joining, leaving, and permission checks work according to the selected policy.
- Blocking prevents the defined contact and messaging behaviors.
- Unauthorized users cannot read or mutate conversations they do not belong to.
- Core flows have automated tests and manual checks for mobile, keyboard, and screen-reader behavior.
- Error states provide a recovery path rather than leaving the user in an indeterminate state.

## 14. Design Callouts

These are the remaining design decisions to work through before implementation. They are listed in dependency order.

### 14.1 Key hierarchy

- Each client generates a device key. There is no primary client.
- Clients use their keys to encrypt and decrypt message content; routing metadata remains visible.
- Private keys remain on clients. A bouncer may relay an explicitly authorized private-key transfer between authenticated user devices, but does not store the transferred key.
- The exact distinction between user keys, device keys, message keys, and file keys remains to be formalized.

### 14.2 Device enrollment and verification

- A bouncer may authorize a new device after authenticating the user.
- A client that holds the required key should approve key transfer through a user-visible challenge, such as matching a code or emoji on both screens.
- The exact challenge and device-discovery mechanism remain open.
- The client must know whether it is connected through a bouncer and must be able to authenticate the bouncer separately from relay authentication.

### 14.3 Key recovery

- If no authenticated client holds the required keys, the related conversations are lost unless the user has made a separate key backup.
- Key re-rolling is an explicit user action and is not automatic recovery.
- Backup format, account-recovery behavior, and whether recovery restores old messages remain open.
- A bouncer may relay keys between authenticated devices without decrypting messages.

### 14.4 Bouncer trust model

- The bouncer stores and forwards ciphertext and does not decrypt message content.
- The bouncer may authorize new clients and is the only device the relay needs to recognize for a bouncer-backed user.
- Clients normally do not connect directly to the relay when using a bouncer.
- The authority boundary between bouncer, relay, and direct clients remains an open deployment decision.
- A compromised bouncer may impersonate the authenticated user, request key transfers, initiate key exchanges, and withhold or replay encrypted records.

### 14.5 Room key management

- Direct conversations currently use recipient public keys to establish encrypted communication.
- Rooms need a shared-message-key design; encrypting every message separately to every participant may be too expensive.
- A Matrix-like approach is to use pairwise/device public-key encryption to distribute room session keys, then encrypt room messages with those symmetric session keys.
- Room encryption may eventually offer an advanced mode where clients exchange keys directly, and a basic mode where room key pairs are stored by the server and delegated to clients.
- Room keys should rotate when membership changes, especially when a participant leaves or is removed.
- Decide whether removed members can decrypt messages sent before their removal and whether they can decrypt future messages.

### 14.6 Message and file encryption format

- Which algorithms and envelope fields are required?
- How are nonces, authentication tags, key IDs, and protocol versions represented?
- How are encrypted file keys represented?
- How are key rotation and algorithm migration represented?
- Direct-message encryption should support key rolling so users can explicitly re-encrypt or replace communication keys.

### 14.7 Relay authority and federation

- The relay authenticates the user or bouncer and authorizes delivery, but does not inspect message content.
- How do relays discover and authenticate one another?
- Can a relay reject delivery using only visible metadata?
- Sender, recipient, and routing metadata remain visible so relays can deliver messages; message subject and data remain encrypted.
- Whether clients are always locked out of direct relay connections when a bouncer exists remains open.

### 14.8 Queue behavior

- How long does a relay retain undelivered ciphertext?
- What are the queue size and envelope-size limits?
- When does a queued message expire or become permanently failed?
- What does the sender see when a destination relay or bouncer rejects delivery?

### 14.9 Local history behavior

- What exactly does a client store locally?
- How are edits and deletes represented in local history?
- What happens when local history conflicts with a record received from the bouncer?
- Which client is allowed to make a local record authoritative for that user’s other clients?

### 14.10 Identity format

- Should identities use `@user#homeserver`, `user@homeserver`, or another format?
- How are users and rooms uniquely resolved?
- Can a user move to another relay without changing identity?

### 14.11 Encrypted moderation

- If messages are end-to-end encrypted, what does reporting include?
- Can a user submit decrypted evidence with a report?
- Can moderators inspect anything without weakening the encryption model?
- Which metadata may be retained for abuse prevention?

### 14.12 Large file transfer

- What is the maximum file size?
- Which resumable upload protocol is used?
- Where are encrypted files stored?
- How do expiration and cleanup work?
- Which file metadata is visible to relays and storage providers?

## 15. Decisions Required Before Implementation

The following decisions are intentionally open and should be resolved before the data model and API are finalized:

1. **Target platforms:** web only, native mobile, or web plus mobile.
2. **Identity model:** email plus username, phone number, invite-only accounts, or another model.
3. **Message privacy:** the current design uses end-to-end encrypted messages from the start.
4. **Group size:** maximum number of participants and whether groups may be converted or renamed.
5. **Membership history:** whether users who leave can retain access to prior messages.
6. **Message limits:** maximum text length, edit window, delete behavior, and attachment roadmap.
7. **Delivery semantics:** whether the product needs per-user delivery receipts in groups.
8. **Notifications:** in-app only, browser push, native push, email, or a staged approach.
9. **Moderation model:** admin-only review, automated screening, user reports, or a combination.
10. **Data retention:** account, message, report, log, and backup retention periods.
11. **Availability target:** expected users, regions, uptime, and recovery-point/recovery-time objectives.
12. **Visual direction:** product name, tone, branding, and primary client interaction model.

## 16. Suggested First Decisions

To keep implementation focused, the recommended initial choices are:

- Web client first, responsive on desktop and mobile.
- Email/password authentication with a unique public username.
- End-to-end encrypted messages from the start. Clients decrypt content; bouncers store encrypted history and provision encrypted keys across authorized devices.
- Direct conversations and groups capped at 50 participants.
- In-app notifications first; platform push notifications later.
- Cursor-based pagination for local client and bouncer history.
- Fifteen-minute edit and delete windows, with deleted messages retained as timeline markers.
- User-created groups with one owner and basic participant management.
- Reports and blocks included in the first private beta, even if moderation operations are initially manual.

These recommendations are defaults, not requirements.

