---
id: REQ-072
title: Go utility allocates REQ ids and version numbers and verifies release consistency
status: pending
created_at: 2026-08-03T11:41:15Z
user_request: UR-013
domain: backend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-071, REQ-073]
batch: parallel-builds
write_set: [tools/queue-kanban/main.go, tools/queue-kanban/*.go, actions/capture.md, actions/forensics.md, actions/work.md, CLAUDE.md, _dev/tests/contract-regressions.sh]
---

# Go Utility Allocates REQ Ids and Version Numbers and Verifies Release Consistency

## What

Add a subcommand to the shipped Go tool that allocates the next REQ number, allocates and writes the
next version, and verifies the release-ritual invariants that are currently checked by hand. It
**never** writes the changelog body.

## AI Execution State (P-A-U Loop)

- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Two of the pre-commit rules are cross-file consistency checks a human is asked to perform by hand
every single commit, and both are documented as *already having been gotten wrong* — the version must
be strictly greater than the newest `CHANGELOG.md` entry, and the entry title must not already be in
use. Those are the checks worth automating. The allocation is the easy half; the verification is the
half that is actually failing.

## Context

- `tools/queue-kanban/main.go:16-30` — a hand-rolled `os.Args[1]` subcommand switch with no external
  CLI library, each subcommand owning its own `flag.FlagSet`. Existing subcommands: `summary`,
  `generate`, `serve`. This is where the new one slots in.
- `tools/queue-kanban/frontmatter.go`, `walk.go`, `model.go` — the queue scan already exists. Reuse it;
  do not add a second parser.
- `actions/capture.md:73` — today's REQ allocation: scan `REQ-*.md` across `do-work/queue/`,
  `do-work/working/`, and `do-work/archive/` (including inside `do-work/archive/UR-*/`), then
  increment from the highest. REQ and UR use separate sequences; start at 1 if none exist.
- `actions/version.md:5` — `**Current version**: 0.163.3`. Note this lives inside an *action file* (a
  prompt the agent reads), not a dedicated data file. Anchor on the `**Current version**: ` prefix.
- `actions/forensics.md:23` — `## Checks`, thirteen existing read-only probes. Checks 9, 11, 12 and 13
  overlap the verify surface; reuse their definitions rather than re-deriving them.
- `actions/forensics.md:3` — forensics is **read-only by contract**, which is why `verify` fits there
  and the allocating subcommands do not.
- `actions/cleanup.md` — the consent-gated fixer. Fixes belong there, not in this tool.

## Detailed Requirements

1. **`next-req`** — returns the next free REQ number using `actions/capture.md:73`'s existing rule
   (max across queue, working, archive, including `archive/UR-*/`; start at 1 when none). Read-only
   with respect to the queue.
2. **`next-version <patch|minor|major>`** — computes the next semantic version from
   `actions/version.md:5` and writes it back, returning the new number. The bump size is an
   **argument, not an inference** — patch vs minor vs major is a human judgment the tool must not make.
3. **`verify`** — read-only. Reports, at minimum:
   - the current version is strictly greater than the newest `CHANGELOG.md` entry
   - the newest entry's title is not reused by an earlier entry
   - duplicate REQ numbers across queue / working / archive
   - orphan `worktree-agent-*` worktrees and branches
   - `do-work/CHECKPOINT.md` naming a REQ that no longer exists
   - claims past the staleness threshold, and unparseable or future-dated `claimed_at`
   - any `worktree-agent-*` worktree whose `do-work/` differs from the main tree's (an owner
     impersonation — a builder that wrote queue state)
4. **Never write `CHANGELOG.md`.** Unique version numbers do not make a shared prepend safe, so the
   changelog stays an owner-only, human-authored write.
5. **`verify` exits non-zero on findings and routes the fixable ones.** Its report must name which
   findings are mechanically fixable and point at `do-work cleanup` for them — e.g.
   `3 fixable: run do-work cleanup`. This is how the user's "any fix that needs to be done, can be
   performed" is honored without this tool doing the fixing (see `## Open Questions`).
6. **Allocate, then read back to confirm ownership.** After writing a version, re-read the file and
   confirm the value landed before reporting success.
7. **Gaps are acceptable and require no special handling.** The max+1 rule already tolerates them and
   nothing in the skill walks a contiguous sequence. Do not add gap-filling, gap-detection, or
   compaction.
8. **Wire the three call sites** — `next-req` into `actions/capture.md:73`, `next-version` into the
   Step 9 commit ritual in `actions/work.md`, and `verify` into `actions/forensics.md`'s `## Checks`.
   **No new action file and no new SKILL.md routing row.**
9. **Degrade gracefully when `go` is absent**, exactly as `actions/board.md` does — the call sites must
   fall back to the existing manual procedure rather than failing. The Go toolchain is a documented
   exception for this one tool, never a hard dependency of the pipeline.
10. **Update `CLAUDE.md:159` in the same commit.** It currently states the board's Testing view is the
   tool's "**ONE** write surface… it never touches `status` or any other pipeline field." An allocator
   makes it two. This is the co-location rule applied to itself; it is prose-only with no assertion
   behind it, which is exactly why it gets forgotten.

## Constraints

- Ships inside `tools/queue-kanban/` (its own Go module) and travels in the tarball like the rest of
  the tool. No independent versioning — changes get a normal root `CHANGELOG.md` entry and a skill
  version bump.
- **Not trivially atomic, and that is accepted.** Exclusive-create on the new filename does not
  protect the *number*: two allocators both computing 172 with different title slugs both succeed.
  True atomicity needs a number-keyed marker file, i.e. new durable state, which this batch forbids.
  Allocation is human-initiated at capture time and runs in milliseconds — do not build a locking
  scheme, and do not claim in prose that duplicates are impossible.
- **Never mutates queue or REQ state.** `next-version` writes `actions/version.md` and nothing else —
  no REQ field, no `status`, no file in `do-work/`. The tool's write surface widens by exactly one
  file. Repairs are `actions/cleanup.md`'s job, which asks before it acts.
- **Justify this REQ on its own merits, independent of parallel builds.** Its value is that two
  pre-commit consistency checks are performed by hand today and have already been gotten wrong
  (`CLAUDE.md` § Before Every Commit, items 1 and 2). Do not frame it as concurrency support.
- Never commit the compiled binary (`tools/queue-kanban/.gitignore` already covers it).

## Dependencies

None. Independent of REQ-071 and REQ-073, though its `verify` probes cover state both of them
produce.

## Builder Guidance

**Certainty level: Firm** on the subcommand split, the never-write-changelog boundary, the
bump-size-as-argument rule, and the `CLAUDE.md` contract edit — all settled with the user.

**Mixed** on the `verify` probe list: the seven items above are the floor, not a closed set. Where a
probe duplicates an existing forensics check, reuse that check's definition and say so rather than
writing a second, subtly different one.

Latitude on: subcommand naming, output format (human-readable is fine; no machine-readable format is
required), and the staleness threshold's plumbing. Keep it small — this is a few hundred lines of Go
in an existing module.

## Open Questions

- [x] The original request said "any fix that needs to be done, can be performed" — should this tool
  perform repairs, or only detect and route them? → **Report only; fixes stay in
  `actions/cleanup.md`.** Resolved by the user at verify time (`do-work verify-requests`, 2026-08-03),
  choosing "Report only, fixes via cleanup" over a fix mode and over a queue-excluded partial fix mode.
  Rationale accepted: `actions/forensics.md:3` is read-only by contract and the board tool has one
  narrow write surface, so a repairing binary would change two contracts, while `actions/cleanup.md`
  already asks before it acts. **The intent is honored, not dropped** — requirement 5 makes `verify`
  name its fixable findings and point at `do-work cleanup`, so a single cheap invocation still tells
  you what to run.

## Red-Green Proof

**RED prompt/case:** In `tools/queue-kanban/`, add table tests against a synthetic queue fixture (the
pattern `board_synthetic_test.go` and `future_timestamp_test.go` already use):

1. `next-req` on a fixture whose highest id is `REQ-070` returns `71`; on an empty fixture returns `1`;
   on a fixture with a gap (`REQ-001`, `REQ-070`) still returns `71`.
2. `next-version patch` on `**Current version**: 0.163.3` writes and returns `0.163.4`.
3. `verify` fails on a fixture where `actions/version.md` equals the newest `CHANGELOG.md` heading, and
   on one where the newest entry title duplicates an earlier one.
4. `verify` fails on a fixture with two `REQ-071-*.md` files under different slugs.

**Why RED now:** No such subcommand exists — `tools/queue-kanban/main.go:16-30` dispatches only
`summary`, `generate`, and `serve`, so every case above fails to compile or exits unrecognized.

**GREEN when:** `go test ./...` passes in `tools/queue-kanban/`, and on this live repo
`next-req` prints `74`, `verify` passes clean, and a hand-planted duplicate version or reused title
makes `verify` fail naming the offending entry.

**Validation:** User confirmed — the user proposed the utility ("the golang utility that we talked
about, can also update a version number and return the new version number, we can do that with the
REQ numbers as well… we just eliminated duplicates in effect") and then chose the "allocate + verify,
never write the changelog" scope from an explicit option prompt.

## Full Context

See `do-work/user-requests/UR-013/input.md` for complete verbatim input.
