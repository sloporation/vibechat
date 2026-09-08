# Agent Instructions

## Git workflow

When committing changes, always push to a new branch and create a PR. Always put the agent name as a tag in the commit, for example `[AGENT] Some PR name`. Always link the agent's GitHub user to the PR.

The configured Git identity for repositories under `~/Projects/github` uses the GitHub user `jsevenb`.

## Specification conventions

- `SPEC.md` is the active protocol and system specification.
- `old/` contains superseded Markdown references; consult it for historical context but do not treat it as the active design.
- VibeChat is organized around channels. A channel can be a direct 1-to-1 chat or a multi-participant room.
- Each channel has its own append-only transaction log and transaction cursor. There is no global transaction cursor in the MVP.
- Outgoing client commands are POST-style; incoming server responses/events are GET-style. Over WebSocket, client-to-server frames are outgoing/POST-equivalent and server-to-client frames are incoming/GET-equivalent.
- Payloads use direction markers: `d: "out"` for client-to-server and `d: "in"` for server-to-client.
- Incoming committed channel events include the channel transaction ID and server-generated message ID. Outgoing messages include a client-generated idempotency ID instead.
- Changes to the protocol must update `SPEC.md` and preserve idempotent delivery and replay semantics.

Do not modify `AGENTS.md` unless the user explicitly asks for an update.
