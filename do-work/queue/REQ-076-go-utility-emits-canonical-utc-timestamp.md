---
id: REQ-076
title: Go utility emits the canonical UTC timestamp, preferred over date -u when built
status: pending
created_at: 2026-08-03T15:45:29Z
status_changed_at: 2026-08-03T15:45:29Z
user_request: UR-014
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
maintenance: false
review_generated: false
write_set: [tools/queue-kanban/main.go, tools/queue-kanban/timestamp.go, tools/queue-kanban/timestamp_test.go, actions/work-reference.md, CLAUDE.md]
---

# Go Utility Emits the Canonical UTC Timestamp, Preferred Over `date -u` When Built

## What

`tools/queue-kanban` already allocates the two other things the pipeline needs to get right and
cannot guess — REQ numbers (`next-req`) and version numbers (`next-version`). Timestamps are the
third, and they are still obtained by every action shelling out to `date -u +%Y-%m-%dT%H:%M:%SZ`
(the Timestamp rule, `actions/work-reference.md` → **Request File Schema — Full Frontmatter**).

Add a `now` subcommand that prints the current UTC instant in exactly the schema's ISO-8601 form,
and make the Timestamp rule name it as the **preferred** source when the binary is already built,
keeping `date -u` as the documented fallback.

Two things this buys beyond tidiness:

- **Windows `cmd` has no `date -u +FORMAT`** — the prescribed command silently produces garbage
  there today. `queue-kanban now` is the first portable path.
- **One implementation of the format**, next to the two readers that already parse it
  (`tools/queue-kanban/model.go`'s state timer, the future-stamp guard) — so the writer and the
  readers agree by construction rather than by convention.

## Why This Is Worth a REQ Rather Than a Sweep

The Go toolchain is **optional by design**: `CLAUDE.md` → Shipped Tooling calls the board "the one
capability that needs a compiler" and says not to reach for a compiled tool in any other action.
Timestamps are written by nearly every action, so a naive implementation — "call the Go tool for
every `*_at` stamp" — would promote a compiler to a hard dependency of the whole pipeline and break
"design for the floor" (`CLAUDE.md` → Agent Compatibility). That constraint, not the subcommand, is
what makes this a REQ: the work is mostly getting the preference order and the carve-out wording
right.

## Detailed Requirements

1. **Add `queue-kanban now`** — prints the current UTC instant as `YYYY-MM-DDTHH:MM:SSZ` and a
   trailing newline, nothing else, so a caller can use it directly the way `next-req` is used.
   Read-only: it touches no file, which means the "exactly two write surfaces" sentence in
   `CLAUDE.md` → Shipped Tooling stays true as written and needs no amendment.
2. **Amend the Timestamp rule in one place** (`actions/work-reference.md`, the canonical home) to
   state the preference order. Do **not** edit the ~11 sites in 8 files that cite the rule
   (`actions/capture-reference.md`, `actions/work.md`, `actions/memory-reference.md`,
   `actions/memory.md`, `actions/forensics.md`, `actions/clarify.md`, `actions/abandon.md`) — they
   point at the rule, and a copy per site is exactly the staleness `CLAUDE.md` → Closed Enumerations
   Go Stale warns about. State the trigger condition once; let the citations inherit it.
3. **The preference order must not prescribe a compile.** Use the already-built binary if it is
   there; otherwise `date -u`. Never `go run .` for a timestamp — that pays a Go compile per stamp
   and is strictly worse than the fallback it would replace. `actions/board.md` is the only thing
   that builds the binary, so "already built" is the honest condition, not "Go is installed".
4. **Name a Windows fallback** for the no-binary case, since `date -u +FORMAT` is a POSIX-ism:
   PowerShell's `Get-Date` in UTC round-trip form. One clause in the rule, not a platform matrix.
5. **Carve out the toolchain exception** in `CLAUDE.md` → Shipped Tooling: the "don't reach for a
   compiled tool in any other action" sentence now has one narrow, optional-by-construction
   exception, and the sentence must say so — an unamended prohibition plus a shipped violation is
   how a rule stops being read.
6. **Leave every other subcommand alone.** No change to `summary`, `generate`, `serve`, `next-req`,
   `next-version`, `verify`, or to the board model.

## Non-Goals

- Making the Go tool the only timestamp source (explicitly rejected during capture — see UR-014).
- Any change to how timestamps are *read* (`model.go`'s state timer, the 2-minute skew allowance,
  the future-stamp badge).
- Retro-fixing existing stamps in `do-work/`.

## Red-Green Proof

**RED prompt/case:** A Go test asserting `queue-kanban now` prints a value matching
`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$` that the tool's own frontmatter timestamp parser accepts,
and that round-trips to within a second of the injected clock. Fails today — the subcommand does not
exist and `main.go`'s dispatch exits 2 on an unknown subcommand.

**Why RED now:** `tools/queue-kanban/main.go:48` dispatches six subcommands; `now` is not one of
them.

**GREEN when:** That test passes, `go test ./...` in `tools/queue-kanban/` stays green, the
subcommand appears in `main.go`'s doc-comment usage block, and the Timestamp rule in
`actions/work-reference.md` states the preference order with the fallback intact.

Inject the clock rather than calling `time.Now()` inside the tested function — the existing tests
already thread a clock through (`LoadBoard(..., time.Now(), ...)`, `runVerifyProbes(..., time.Now())`),
so follow that shape.

## Origin

Raised by the user mid-session during `do-work clarify` on REQ-074, which was itself about a missing
`status_changed_at` stamp. Captured as its own REQ rather than folded into REQ-074: that one is a
one-line prose fix to Crash Recovery, this one adds Go code, amends the canonical Timestamp rule, and
touches the toolchain-exception wording.

## Full Context

See `do-work/user-requests/UR-014/input.md` for the verbatim ask and the recorded capture decision
(including the two rejected shapes).
