---
id: REQ-076
title: Go utility emits the canonical UTC timestamp, preferred over date -u when built
status: completed
claimed_at: 2026-08-03T16:09:06Z
completed_at: 2026-08-03T16:14:26Z
kb_status: pending
route: C
created_at: 2026-08-03T15:45:29Z
status_changed_at: 2026-08-03T15:45:29Z
user_request: UR-014
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
maintenance: false
review_generated: false
write_set: [tools/queue-kanban/main.go, tools/queue-kanban/timestamp.go, tools/queue-kanban/timestamp_test.go, tools/queue-kanban/prime-do-kanban.md, actions/work-reference.md, CLAUDE.md]
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

---

## Triage

**Route: C** - Complex

**Reasoning:** New Go source plus tests, an amendment to a canonical rule with a preference order and a
platform fallback, and a carve-out in a prohibition that currently forbids the whole change. Three
subsystems (tool, pipeline instructions, maintainer contract) and a documented conflict to resolve — not
a single-location change.

**Planning:** Required

## Plan

1. Read `allocate.go` + `allocate_test.go` as the shape to copy (a read-only allocator subcommand with
   injected inputs) → verify: the new file mirrors that structure, not a new pattern.
2. Find the read-side parser so the writer can be proven compatible → verify: `parseTimestamp`
   (`model.go:916`) accepts what the writer emits, asserted in a test.
3. Write the failing tests first (`tdd: true`) → verify: `go test` fails on undefined symbols.
4. Implement `formatCanonicalTimestamp` / `writeCanonicalTimestamp` + `now` dispatch → verify: tests
   pass, `./queue-kanban now` output is byte-identical to `date -u +%Y-%m-%dT%H:%M:%SZ`.
5. Amend the Timestamp rule in its canonical home only → verify: the citing sites are untouched
   (`git diff --stat` names no other action file).
6. Carve out `CLAUDE.md`'s compiled-tool prohibition and update its read-only subcommand list → verify:
   the prohibition no longer contradicts a shipped instruction.

**Plan validation:** 6 tasks, above the 3-task comfort threshold — but tasks 1–2 are reads and 5–6 are
one sentence each, so the implementation surface is two files. No requirement is uncovered; no task
traces to something the REQ didn't ask for.

## Exploration

- **Pattern to copy:** `allocate.go` + `allocate_test.go`. A subcommand is a small file with a doc
  comment carrying the *why*, a pure function, and a `run*Command` in `main.go` that owns its FlagSet.
  Clocks are injected, never called inside the tested function — `LoadBoard(..., time.Now(), ...)` and
  `runVerifyProbes(..., time.Now())` both take theirs from the caller.
- **The read side to prove compatibility against:** `parseTimestamp` (`model.go:916`) tries
  `time.RFC3339` among other layouts; `coerceScalarToString` (`model.go:1300`) re-normalizes YAML-parsed
  times to RFC3339 UTC. A round-trip test through `parseTimestamp` is the honest proof that writer and
  readers agree.
- **`time.RFC3339` is the wrong layout to emit**, though it is right to parse: it renders a numeric
  offset for a non-UTC instant (`+09:00`) and carries sub-second digits — neither of which the schema's
  `*_at` shape accepts. Hence an explicit `2006-01-02T15:04:05Z` layout plus `.UTC().Truncate(time.Second)`.
- **Citing sites confirmed as pointer-only** (requirement 2): 11 `date -u` mentions across 8 files, all
  of which name the Timestamp rule rather than re-deriving it. None needed editing.
- **No user-facing doc lists the subcommands** (`grep next-req/next-version docs/` → empty), so the only
  enumerations to keep in sync are `main.go`'s usage block, the prime's header and file map, and
  `CLAUDE.md`'s read-only list.

## Scope

**Files I will touch:**
- `tools/queue-kanban/timestamp.go` (new) — the formatter and the writer
- `tools/queue-kanban/timestamp_test.go` (new) — the RED cases
- `tools/queue-kanban/main.go` — dispatch arm, usage block, write-surface note
- `actions/work-reference.md` — the Timestamp rule, one place
- `CLAUDE.md` — the toolchain carve-out and the read-only subcommand list
- `tools/queue-kanban/prime-do-kanban.md` — header subcommand list + file map (**added mid-build, D-01**)

**Acceptance criteria (restated from the REQ):**
- `queue-kanban now` prints `YYYY-MM-DDTHH:MM:SSZ` and a newline, nothing else; touches no file.
- The Timestamp rule states the preference order in one place; the ~11 citing sites are untouched.
- The order never prescribes a compile — the condition is "binary already built".
- A Windows fallback is named for the no-binary case.
- `CLAUDE.md`'s compiled-tool prohibition carries the exception.
- No other subcommand, and no board behavior, changes.

## Pre-Flight

Baseline recorded at the start of the run (`tools/checks/preflight.sh "bash _dev/tests/contract-regressions.sh"`):
tree clean outside `do-work/`, baseline passing, `launched: true`. One note, not a blocker: the baseline
covers the contract suite only — the Go suite is this REQ's own surface and was confirmed green
separately before the first edit (`go test ./...` → `ok`).

## Decisions

**D-01 — Added `tools/queue-kanban/prime-do-kanban.md` to scope. DECIDE & STATE.**
The prime is what an agent reads before editing the tool, and it carries both a subcommand list in its
header and a per-file map under "Read first". A new source file and a new subcommand that appear in
neither would make the prime wrong on its first line — the same staleness class REQ-075 just finished
clearing. One line added to each list.

**D-02 — Emit an explicit layout constant rather than `time.RFC3339`. DECIDE & STATE.**
`time.RFC3339` is correct on the *parse* side and wrong on the *emit* side: for a non-UTC instant it
writes `+09:00` instead of converting, and it can carry sub-second digits. Both would produce a stamp the
schema's shape rejects and the board's future-stamp guard could flag. The constant carries a comment
saying so, because "why not just RFC3339" is the obvious next edit.

**D-03 — Truncate sub-second precision, don't round. DECIDE & STATE.**
Rounding `16:09:06.999` up to `:07` moves a stamp forward, and the board treats a stamp past its own
`now` (plus 2 minutes of skew) as clock skew. Truncation can only ever be a fraction of a second stale,
which nothing reads. Pinned by a test.

**D-04 — `now` takes no `--repo-root`. DECIDE & STATE.**
Every other subcommand resolves a repo root because it reads the tree; this one reads a clock. Accepting
the flag would imply the stamp depends on the project, and it doesn't. It still parses a FlagSet so an
unrecognized flag is rejected rather than silently swallowed.

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/timestamp.go` (new) — `canonicalTimestampLayout`, `formatCanonicalTimestamp`,
  `writeCanonicalTimestamp`; the file's doc comment carries why the tool has a timestamp writer at all
  and why it is not the required source
- `tools/queue-kanban/timestamp_test.go` (new) — five cases: exact shape, round trip through the board's
  own `parseTimestamp`, non-UTC conversion, sub-second truncation, exact stdout bytes
- `tools/queue-kanban/main.go` (modified) — `now` dispatch arm, `runNowCommand`, usage block, unknown-
  subcommand message, and the write-surface note (now records `now` as read-only and says why it takes
  no `--repo-root`)
- `actions/work-reference.md` (modified) — the Timestamp rule now states a three-option preference order
  with the never-build-to-stamp guard
- `CLAUDE.md` (modified) — the compiled-tool prohibition carries the narrow exception, and the read-only
  subcommand list includes `now`
- `tools/queue-kanban/prime-do-kanban.md` (modified) — header subcommand list + `timestamp.go` in the
  file map (D-01)

**What was done:** Added a `now` subcommand that prints the current UTC instant in exactly the schema's
shape, and made the Timestamp rule prefer it *when the binary already exists* while keeping
`date -u` as the floor and naming a PowerShell fallback for Windows `cmd`. The rule was amended in one
place; the ~11 citing sites across 8 action files were deliberately left alone. `CLAUDE.md`'s "don't
reach for a compiled tool in any other action" now states the exception's gate (already-built plus a
documented fallback) instead of silently contradicting a shipped instruction.

## Qualification

Passed — 6 files verified, 6 acceptance criteria traced, scope-drift clean after D-01's declared
extension. `gofmt -l .` empty, `go vet ./...` clean. Mechanical checks green via
`tools/checks/qualify.sh`. The two `(new)` files are substantive (48 and 80 lines) and both are
referenced: `timestamp.go` by `main.go`'s dispatch, `timestamp_test.go` by the test runner.

## Testing

**Tests run:**
- `cd tools/queue-kanban && gofmt -l . && go vet ./... && go test ./...` → ✓ `ok` (all suites, 5 new cases)
- `bash _dev/tests/contract-regressions.sh` → ✓ passed
- `./tools/queue-kanban/queue-kanban verify --repo-root .` → ✓ no findings

**Red-green validation:** *(`tdd: true` — tests written before `timestamp.go` existed)*
- `TestFormatCanonicalTimestampMatchesSchemaShape`: ✗ build failure, `undefined: formatCanonicalTimestamp` → ✓ pass
- `TestFormatCanonicalTimestampRoundTripsThroughBoardParser`: ✗ → ✓ pass
- `TestFormatCanonicalTimestampConvertsNonUTCZones`: ✗ → ✓ pass
- `TestFormatCanonicalTimestampTruncatesSubSecondPrecision`: ✗ → ✓ pass
- `TestWriteCanonicalTimestampEmitsOnlyTheStampAndOneNewline`: ✗ `undefined: writeCanonicalTimestamp` → ✓ pass

**Acceptance check against the live binary:**
- `./queue-kanban now` → `2026-08-03T16:10:43Z`; `date -u +%Y-%m-%dT%H:%M:%SZ` in the same block →
  `2026-08-03T16:10:43Z`. **Byte-identical**, which is the substitutability the rule's option (1) claims.
- `./queue-kanban bogus` → the usage message now lists `now`.

**New tests added:** the five above, in `tools/queue-kanban/timestamp_test.go`.

## Review

**Overall: 92%** | 2026-08-03T16:13:45Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 90% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 2 minor

- **Minor — the preferred source will rarely be reached in practice, by design.** Requirement 2 forbade
  touching the citing sites, and 10 of the 11 don't merely cite the rule — they **inline
  `date -u +%Y-%m-%dT%H:%M:%SZ`** as the command to run (`actions/capture-reference.md` ×3,
  `actions/work.md` ×2, `actions/memory-reference.md` ×2, `actions/abandon.md`, `actions/clarify.md`,
  `actions/forensics.md`). An agent following any of those never consults the preference order, so it
  never reaches option 1 — and, more importantly, a Windows agent still gets the broken command. The
  requirement's reasoning (don't make 11 copies of the order) is sound; the fix that serves both is
  **subtraction, not duplication** — strip the inline command from the citing sites so they name the rule
  only, leaving exactly one place that says how to get a stamp. That is a maintenance pass on its own
  (`crew-members/maintenance.md` § delete before you add) and is offered rather than taken here, because
  narrowing 10 sites is not what this REQ was scoped to do.
- **Minor — the Windows path is prescribed but unverified.** `Get-Date -AsUTC` requires PowerShell 7+
  (`-AsUTC` was added in 7.0; on Windows PowerShell 5.1 the equivalent is
  `(Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")`). Nothing in this environment can run
  either, so the fallback is reasoned, not tested. Flagged rather than silently shipped as verified.

**Restatement sweep:** Run. The redefined element is the Timestamp rule's *how to obtain it* clause — the
shape (`YYYY-MM-DDTHH:MM:SSZ`) is unchanged, so every reader of the format stays correct. Enumerations of
the tool's subcommands were the real consumers and all four were updated in the same commit: `main.go`'s
usage block, `CLAUDE.md`'s read-only list, and the prime's header list + file map. `docs/` lists no
subcommands (verified by grep), and no test fixture encodes the old count. The citing sites are the
finding above, not a stale restatement — option 2 remains valid, just no longer the only option.

**Acceptance:** Pass — `./queue-kanban now` and `date -u +%Y-%m-%dT%H:%M:%SZ` produced byte-identical
output in the same shell block, which is the substitutability claim the rule makes.
**Suggested testing:** 2 items — (1) run the prescribed Windows fallback on a real Windows `cmd`/PowerShell
host and confirm the exact shape, then correct the rule if 5.1 support matters; (2) exercise the
preference order end-to-end by deleting the built binary and confirming an action still stamps correctly
via option 2 (the fallback is the floor, so it is the path that must never break).
**Follow-up REQs created:** None — both findings are Minor, and the first is a deliberate scope boundary
the user set during capture rather than a defect.

## Lessons Learned

**What worked:** Testing the writer against the project's own reader (`parseTimestamp`) rather than
against a hand-written expectation. That is the assertion that would actually catch a future "let's just
use `time.RFC3339`" edit, because it fails on the offset form a bare shape test would accept.

**What didn't:** The first instinct was `time.RFC3339` for both directions — it is the obvious constant
and it is wrong on the emit side (numeric offset for non-UTC, sub-second digits). Caught while writing
the non-UTC test, not while writing the code, which is the argument for writing that test first.

**Worth knowing:** Adding a *preferred* source to a rule does not make it used — the sites that inline
the old command keep teaching it. When a canonical rule grows a preference order, the citing sites either
have to stop restating the mechanism or they quietly pin the old one. That is the same failure shape
REQ-075 just fixed for `write_set`, arriving from the opposite direction: there the restatements carried
a stale *reason*, here they carry a stale *mechanism*.

## Orientation

`do-work` can now get its timestamps from the Go utility instead of the shell: `queue-kanban now` prints
the exact stamp every `*_at` field wants, and the Timestamp rule names it as the preferred source
whenever the binary is already built — with `date -u` still the documented floor, so nothing depends on a
compiler. This closes a real gap on Windows `cmd`, where the prescribed `date -u +FORMAT` does not exist
at all. Lives in the queue-kanban tool alongside the other two release-ritual allocators (`next-req`,
`next-version`), plus the pipeline's Timestamp rule and the maintainer-side toolchain exception it needed
carved out. Not `[MAP CHANGED]` — a fourth read-only subcommand in an established pattern, no new
module, contract, or data flow.

**Prime staleness spot-check** (`tools/queue-kanban/prime-do-kanban.md`): the change made it stale in two
places (header subcommand list, "Read first" file map) and both were updated in this REQ (D-01). All its
other referenced paths still exist. Not stale as of this commit.
