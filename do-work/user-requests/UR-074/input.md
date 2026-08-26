---
id: UR-074
title: 'Show the implementation duration on recently done cards'
created_at: 2026-08-26T13:02:22Z
requests: [REQ-374]
word_count: 75
---

# Show the Implementation Duration on Recently Done Cards

## Full Verbatim Input

> ````text
> Do-work capture-request. So when we show the recently done cards, please also show the duration, how long it took since it was started until it is finished to implement that card to make it delivered. By making it delivered, I mean it was moved to the Done column, completed status.
>
> This is in the Kanban, so when we run, just run Kanban, this interface shows up via Go utility.
>
> after capture, run do-work verify-request
> ````

The third line is a session directive — run `do-work verify-requests` once capture lands — carried out in-session. It is not queued work and no REQ covers it.

A screenshot accompanied the first line. It could not be persisted as an asset in this session (no file-backed attachment path was available); REQ-374's Assets section carries the full text description instead.

## Capture-Time Decisions

- **Odd spans are shown and marked, not hidden.** Asked during capture, the user chose: a normal card reads its span plainly, a span over four hours renders with a paused marker, a negative span renders as a broken-stamp warning, and a card with no parseable start stamp shows no duration at all.

---
*Captured: 2026-08-26T13:02:22Z*
