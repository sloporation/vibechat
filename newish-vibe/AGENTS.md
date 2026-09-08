# Project Context

## Working directory

This project is worked on from `~/Projects/github/sloporation/vibechat`.

## Context files

- `SPEC.md` is the broad product and architecture specification. Use it for background and context.
- `SUB_SPEC.md` is the deliberately reduced MVP specification. Use it as the active implementation scope unless a decision explicitly overrides it.
- This file, `AGENTS.md`, is durable project context for work performed in this directory.
- `SUB_SPEC.md` Section 11 defines the standardized JSON envelope for MVP HTTP responses and WebSocket messages.

## MVP direction

Prioritize a small, reliable direct-messaging web MVP over feature breadth. The target loop is:

1. Register and sign in.
2. Find another user by exact username.
3. Start or reopen a direct conversation.
4. Send and receive plain-text messages in real time.
5. Persist and reload message history.
6. Sign out and sign back in.

Keep groups, end-to-end encryption, federation, bouncers, native mobile clients, files, notifications, moderation, advanced delivery state, and multi-device synchronization out of the MVP unless explicitly reintroduced.

## Working convention

Before making implementation decisions, inspect `AGENTS.md` and `SUB_SPEC.md`; consult `SPEC.md` when broader context is needed. Record important scope or architecture decisions in this file or the relevant spec so they survive future sessions.
