---
id: UR-014
title: The Go utility should generate timestamps too
created_at: 2026-08-03T15:45:29Z
requests: [REQ-076]
word_count: 12
---

# The Go Utility Should Generate Timestamps Too

## Summary

Raised mid-session during `do-work clarify` on REQ-074 (the crash-recovery `status_changed_at`
question). The user wants `tools/queue-kanban` to be a timestamp source alongside its existing
`next-req` / `next-version` allocators, so the skill's `*_at` stamps stop depending solely on
`date -u`.

## Verbatim Input

> also use the go utility to generate the time as well

## Clarification Taken During Capture

The Go toolchain is deliberately optional in this skill — the board is "the one capability that
needs a compiler," and no other action may require one. Making it the sole timestamp source would
turn a compiler into a hard dependency of the whole pipeline, since nearly every action stamps a
timestamp. The user was shown three shapes and chose:

- [x] How should the Go utility's new timestamp command relate to the existing `date -u` rule?
  → **Confirmed: preferred when available, with `date -u` as the documented fallback** (user, via
  `do-work clarify`, 2026-08-03). The skill keeps working on agents with no Go toolchain, and
  Windows `cmd` — which has no `date -u +FORMAT` — gains a working path for the first time.

Rejected alternatives, recorded so they are not re-litigated: making the Go tool the only source
(compiler becomes a hard dependency), and folding the work into REQ-074 (different write set,
different review — REQ-074 is a one-line prose fix).

## Extracted Requests

| REQ | Title | Origin in the input |
| --- | --- | --- |
| REQ-076 | Go utility emits the canonical UTC timestamp, preferred over `date -u` when built | "also use the go utility to generate the time as well" |
