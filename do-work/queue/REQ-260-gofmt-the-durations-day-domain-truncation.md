---
id: REQ-260
title: Gofmt the Durations day-domain truncation expression
status: pending-answers
created_at: 2026-08-18T18:41:26Z
user_request: UR-051
addendum_to: REQ-251
domain: general
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/durations.go
---

# Gofmt the Durations Day-Domain Truncation Expression

## What

`skills/do-work-board/tools/queue-kanban/durations.go:340` fails `gofmt -l`: missing spaces around `*` in `rangeEnd.UTC().Truncate(24*time.Hour)`. Introduced by REQ-248, cosmetic only, and not enforced by the maintainer gate (which runs `go vet`, not `gofmt`). One-character fix; optionally worth asking whether the gate should run `gofmt -l` so this class cannot recur.

## Context

Discovered by REQ-251's builder ([low], pre-existing at its branch point; production file, so the test-hygiene carve-out does not apply and this takes the consent flow).

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-251: a gofmt formatting miss in `durations.go` from REQ-248. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — or fold it into REQ-252, which already owns the file.
