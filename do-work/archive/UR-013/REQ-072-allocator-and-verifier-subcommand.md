---
id: REQ-072
title: Go utility allocates REQ ids and version numbers and verifies release consistency
status: completed
claimed_at: 2026-08-03T14:50:03Z
completed_at: 2026-08-03T15:08:08Z
commit: 5db22ea
route: C
kb_status: promoted
kb_entry: REQ-072-go-utility-allocates-req-ids-and-version.md
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
write_set: [tools/queue-kanban/allocate.go, tools/queue-kanban/release.go, tools/queue-kanban/verify.go, tools/queue-kanban/main.go, tools/queue-kanban/allocate_test.go, tools/queue-kanban/release_test.go, tools/queue-kanban/verify_test.go, tools/queue-kanban/prime-do-kanban.md, actions/capture.md, actions/work.md, actions/forensics.md, docs/forensics-guide.md, CLAUDE.md, _dev/tests/contract-regressions.sh]
---

# Go Utility Allocates REQ Ids and Version Numbers and Verifies Release Consistency

## What

Add a subcommand to the shipped Go tool that allocates the next REQ number, allocates and writes the
next version, and verifies the release-ritual invariants that are currently checked by hand. It
**never** writes the changelog body.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read `tools/queue-kanban/prime-do-kanban.md` (the declared prime) plus `crew-members/general.md`, `coding-guardrails.md`, and `anti-slop.md` (the changelog entry is a human-facing artifact). Technical approach written into `## Plan` above before any code: three new files (`allocate.go` / `release.go` / `verify.go`) beside the existing ones, three arms in `main.go`'s existing hand-rolled switch, no second parser — `next-req` rides `enumerateDoWorkTree` and the duplicate-id probe reads `Board.Warnings`. The prime's `## Lessons` were read; none touched allocation or the release ritual.
- [x] **[APPLY]:** Tests first, confirmed RED (every new symbol undefined → build failure), then the three files, then `main.go`, then the prose call sites. One mid-implementation correction: running `verify` against the live repo showed the version-vs-changelog probe firing on a perfectly healthy tree, so the probe was re-derived (see `## Decisions` D-02) and its tests rewritten. Scope held to the declared list plus `docs/forensics-guide.md`, declared as D-03.
- [x] **[UNIFY]:** `git diff --stat` → 8 files under `tools/queue-kanban/` (+1,739/-9) and 6 outside it (+101/-2, excluding the REQ's own move). Verified each: **allocate.go / release.go / verify.go** — `gofmt` clean, `go vet` clean, every exported and unexported symbol carries a doc comment stating the *why*, no `panic`, no debug prints, no `TODO`; **main.go** — three new arms plus the `want …` list and the package doc's write-surface paragraph, existing arms untouched; **prime-do-kanban.md** — header subcommand list, four `## Read first` entries, four new `## Traps`; **capture.md / work.md / forensics.md** — one accelerator block each, every one stating its missing-`go` fallback (the accelerator paragraph in work.md was relocated *after* the Step 9 procedure so the core prose reads uninterrupted); **CLAUDE.md** — the ONE→two write-surface correction; **docs/forensics-guide.md** — one table row; **contract-regressions.sh** — `bash -n` and `shellcheck` clean. Debug-artifact scan over all four Go files: none. Compiled binary is gitignored and not staged.

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

---

## Triage

**Route: C** - Complex

**Reasoning:** Three new subcommands in a compiled module with table tests, plus three prose call sites, plus a contract edit in `CLAUDE.md` — new code across multiple systems (Go tool + action prose) with a scope boundary the REQ does not settle (which repo's version file `next-version` may write). Planning earns its keep here.

**Planning:** Required

## Plan

**Architecture — three new files beside the existing ones, no new parser.**

| New file | Holds |
| --- | --- |
| `allocate.go` | `next-req` — the REQ-number scan, built on the existing `enumerateDoWorkTree` walk |
| `release.go` | the version-file and `CHANGELOG.md` readers plus the semver bump; `next-version` writes through this |
| `verify.go` | every `verify` probe, the finding type, and the report |

`main.go` gains three `case` arms in the existing hand-rolled `os.Args[1]` switch, each with its own `flag.FlagSet` — the shape the file already documents.

**Reuse, per the REQ's Context:** `enumerateDoWorkTree` (`walk.go`) already collects every `REQ-*.md` under `queue/`, `working/`, and the whole `archive/**` subtree — which is exactly `actions/capture.md`'s allocation rule. `next-req` walks that output; it does not add a second scan. The verify probes that need parsed frontmatter (duplicate ids, claim ages) go through `LoadBoard`, whose `Board.Warnings` already carries the duplicate-id findings, so those are read rather than re-derived.

**The scope question the REQ leaves open — which version file may be written.** Requirement 2 names `actions/version.md`, but that path only exists in a **skill-development checkout**. In a consumer install the skill sits at `.claude/skills/do-work/` while `do-work/` sits at the consumer's repo root, and the version the Step 9 ritual needs there is the consumer's own (`package.json`, `Cargo.toml`, a `VERSION` file — `actions/work-reference.md` → Changelog Entry Procedure resolves that generically). Hard-coding `actions/version.md` would silently bump the wrong file, or nothing.

Resolution: anchor on the **`**Current version**: ` prefix** (which the REQ already prescribes) rather than on a path. `next-version` looks for a file carrying that line — `<repo-root>/actions/version.md` by default, overridable with `--version-file` — and when there is none it **reports and exits non-zero without writing**, and the call-site prose falls back to the Changelog Entry Procedure's own source resolution. Same posture as requirement 9's missing-`go` fallback: the tool is an accelerator for the repo that matches this convention, never a new dependency. Logged as D-01.

**`verify` probes** (requirement 3's seven are the floor; one added, per Builder Guidance):

| Probe | Source of truth reused |
| --- | --- |
| version strictly greater than the newest `CHANGELOG.md` entry | `CLAUDE.md` § Before Every Commit item 1 |
| newest entry title not reused by an earlier entry | `CLAUDE.md` § Before Every Commit item 2 |
| duplicate REQ numbers across queue / working / archive | `Board.Warnings` duplicate-id findings (`model.go`) |
| orphan `worktree-agent-*` worktrees and branches | `actions/cleanup.md` Pass 5 |
| `do-work/CHECKPOINT.md` naming a REQ that no longer exists | — (new; nothing else checks it) |
| claims past the staleness threshold, or an unparseable / future-dated `claimed_at` | `actions/forensics.md` Check 1 + Check 12; threshold mirrors REQ-071's three hours |
| a `worktree-agent-*` worktree whose `do-work/` differs from the main tree's | `actions/work-reference.md` → Worktree Dispatch Mode ("state stays home") |
| **added:** finished REQs stranded in `queue/` or `working/` | `actions/forensics.md` Check 9 — the second genuinely cleanup-fixable finding, which requirement 5 needs to be worth stating |

Each finding carries a `Fixable` flag. Only the two that `actions/cleanup.md` can mechanically resolve (stranded finished REQs → Pass 0, merged orphan worktrees → Pass 5) set it; the rest are human decisions and must not claim otherwise. The report ends with `N fixable: run do-work cleanup` when any are, satisfying requirement 5 without the tool ever repairing anything (`## Open Questions` settled that).

**Test plan (RED first, per `tdd: true`).** One `_test.go` per new file, table-driven over `t.TempDir()` fixtures — the pattern `future_timestamp_test.go` and `board_synthetic_test.go` already use. Cases are the REQ's Red-Green Proof list: `next-req` on a `REQ-070` fixture → 71, empty → 1, gapped (`REQ-001` + `REQ-070`) → 71; `next-version patch` on `0.163.3` → `0.163.4` written and returned; `verify` failing on version-equals-newest-heading, on a duplicated title, and on two `REQ-071-*.md` under different slugs.

**Order of work:** tests → `allocate.go` → `release.go` → `verify.go` → `main.go` wiring → prose call sites → `CLAUDE.md` → contract-suite assertions.

*Plan validation:* every one of the ten Detailed Requirements maps to a task above (1→`allocate.go`, 2→`release.go`+D-01, 3→`verify.go`, 4→a hard no-write boundary in `release.go`/`verify.go`, 5→the `Fixable` flag + report footer, 6→read-back-after-write in `release.go`, 7→no gap handling anywhere, 8→the three prose call sites, 9→the missing-`go` fallback in each call site, 10→`CLAUDE.md`). No orphan tasks. Task count is 8 — past the 3-task comfort line the plan-validation rule flags, but the REQ is deliberately one coherent unit (a subcommand family plus its three wirings), and splitting it would ship a tool nothing calls.

## Exploration

- `tools/queue-kanban/main.go:31-46` — the dispatch switch: `subcommand := os.Args[1]` unless it starts with `-`, then a `switch` with `"", "summary"` sharing an arm and a `default` that exits 2 with a `want …` list. That list must grow with the new arms or the error message lies.
- `tools/queue-kanban/walk.go:105` — `enumerateDoWorkTree(repoRoot)` returns `discoveredTreeFiles{RequestFiles []requestFileReference, …}`, each reference carrying `AbsolutePath` + `TreeSection` (`queue`/`working`/`archive`). `deliverables/`, `runs/`, `assets/`, and dotdirs are pruned; strays outside the three sections are collected separately. Filenames are `REQ-NNN-slug.md`, so the number is available without opening a file — which is what keeps `next-req` in milliseconds.
- `tools/queue-kanban/walk.go:44` — `resolveRepoRoot` walks up for a `do-work/` dir and **skips a skill install** (`isSkillInstallDirectory`: a `SKILL.md` at its top level). Worth knowing for `next-version`: in this repo the walk finds the repo root, and the skill install *is* the repo root, so `actions/version.md` resolves. In a consumer repo the same walk finds the consumer root, where `actions/version.md` does not exist — exactly the case D-01 handles.
- `tools/queue-kanban/model.go:505-538` — duplicate-id resolution already exists: one winner per id by tree-section precedence, with a `duplicate REQ id %s: showing the %s copy …` warning per loser. `verify` reads these instead of re-deriving.
- `tools/queue-kanban/model.go:818-838` + `gitBinaryAvailable()` — the established shell-out shape (`exec.Command("git", "-C", repoRoot, …)`, absent-git tolerated). The worktree probes follow it.
- `tools/queue-kanban/model.go:45` — `futureTimestampSkewAllowance = 2 * time.Minute`, and `parseTimestamp` accepts four layouts. Reuse both; do not introduce a second skew constant.
- `actions/capture.md:73` — the allocation rule in prose, including "REQ and UR use separate numbering sequences" and "start at 1 if none exist." `next-req` covers REQ only; UR allocation stays prose (not in this REQ's requirements).
- `actions/board.md:41-49` — the degradation pattern to copy: run `go version`, and on absence print a fixed two-line message and stop without blocking anything else. Board *stops* because the toolchain is the capability; these call sites must instead **fall back**, which is the difference requirement 9 is pointing at.
- `actions/forensics.md` — thirteen `### N.` checks, read-only by contract (`:3`). The new one lands as `### 14.` and is the only check that shells out to a compiled binary, so it needs its own graceful-skip line.
- `CLAUDE.md:159` — the sentence to amend: the Testing view is "the tool's **ONE** write surface … it never touches `status` or any other pipeline field." After this REQ there are two write surfaces, and the second one writes *outside* `do-work/` — which is the part worth stating, since it is what keeps "never touches pipeline state" true.

## Scope

**Files I will touch:**
- `tools/queue-kanban/allocate.go` (new) — `next-req` scan over the existing walk
- `tools/queue-kanban/release.go` (new) — version-line read/bump/write + `CHANGELOG.md` newest-entry and title readers
- `tools/queue-kanban/verify.go` (new) — the eight probes, the finding type, the report
- `tools/queue-kanban/main.go` (modify) — three subcommand arms + the `default` want-list
- `tools/queue-kanban/allocate_test.go` (new) — REQ-number cases
- `tools/queue-kanban/release_test.go` (new) — bump + write-back + read-back cases
- `tools/queue-kanban/verify_test.go` (new) — probe cases
- `tools/queue-kanban/prime-do-kanban.md` (modify) — the tool's prime; three subcommands change its index
- `actions/capture.md` (modify) — `next-req` call site with the manual fallback
- `actions/work.md` (modify) — `next-version` call site in the Step 9 commit ritual, with the fallback
- `actions/forensics.md` (modify) — `verify` as Check 14
- `docs/forensics-guide.md` (modify) — the user-facing row for that check (see `## Decisions` D-03)
- `CLAUDE.md` (modify) — the write-surface contract sentence
- `_dev/tests/contract-regressions.sh` (modify) — assertions for the call sites and the never-write-changelog boundary

**Files I will NOT touch:**
- `tools/queue-kanban/model.go`, `walk.go`, `frontmatter.go` — reused as-is; adding a second parser is explicitly ruled out
- `tools/queue-kanban/testing.go`, `serve.go`, `generate.go` — the board's existing write surface and renderers are unrelated
- `actions/cleanup.md` — repairs stay there, unchanged; `verify` only points at it
- `CHANGELOG.md` body beyond this REQ's own release entry — the tool must never write it
- `SKILL.md` — no routing row, no new action file

**Acceptance criteria (restated from REQ):**
- [ ] `next-req` returns max+1 across queue/working/archive incl. `archive/UR-*/`; 1 on an empty tree; gaps tolerated
- [ ] `next-version <patch|minor|major>` computes from the `**Current version**: ` line, writes it back, returns the new number, and reads it back to confirm before reporting success
- [ ] Bump size is an argument, never inferred
- [ ] `verify` is read-only and reports all eight probes
- [ ] The tool never writes `CHANGELOG.md`
- [ ] `verify` exits non-zero on findings and names the mechanically fixable ones, pointing at `do-work cleanup`
- [ ] No gap-filling, gap-detection, or compaction anywhere
- [ ] Three call sites wired: `actions/capture.md`, `actions/work.md` Step 9, `actions/forensics.md`
- [ ] No new action file, no new SKILL.md routing row
- [ ] Every call site falls back to the existing manual procedure when `go` is absent
- [ ] `CLAUDE.md`'s write-surface sentence updated in the same commit
- [ ] `next-version` writes exactly one file and nothing in `do-work/`

## Decisions

- **D-01 — `next-version` anchors on the `**Current version**: ` prefix, not on `actions/version.md`, and refuses rather than guessing when no such line exists.** DECIDE & STATE. Requirement 2 names the path, but that path only exists in a skill-development checkout; in a consumer install the skill sits under `.claude/skills/do-work/` while `do-work/` is at the consumer root, so a hard-coded path would resolve to nothing (or, worse, to a coincidence). The prefix is already what the REQ says to anchor on, so this honors the requirement's own instruction over its example path. `--version-file` overrides; a missing line exits 1 **without writing** and the call site falls back to the Changelog Entry Procedure's generic source resolution. Reversible, and the fallback is the safe direction: a repo that doesn't match the convention gets the manual procedure it already had.
- **D-02 — `verify`'s version probe asserts version/changelog *agreement*, and the strictly-greater rule moved to ordering *within* the changelog.** DECIDE & STATE, and it corrects the requirement's literal wording. Requirement 3's first bullet says "the current version is strictly greater than the newest `CHANGELOG.md` entry." That is true only *while composing* a release — the moment the entry is written, version and newest entry are equal, which is the healthy steady state. Built literally, the probe fired on this repo's own clean tree the first time it ran (caught in acceptance testing, not review). So: a mismatch is the finding, with the direction naming the cause (ahead = a bump without its entry, behind = an entry whose version never reached the version file), and the strictly-greater check that *does* survive into the committed state — newest entry beats every earlier entry — became its own probe. That second probe is the one that catches the duplicate version numbers `CLAUDE.md` records as having already happened, so the intent is better served, not weakened. Both pre-commit and post-commit workflows now get a correct answer, which is why this is a correction rather than an escalation.
- **D-03 — Extended Scope to `docs/forensics-guide.md` (one table row).** DECIDE & STATE, declared before the edit. The guide is the user-facing description of what forensics checks; a fourteenth check absent from it is invisible to the person who would run it. One row, no other line touched.
- **D-04 — Added an eighth `verify` probe (finished REQs stranded in `queue/`/`working/`) beyond requirement 3's seven.** DECIDE & STATE, under Builder Guidance's "the seven items above are the floor, not a closed set." Requirement 5 needs the report to name mechanically-fixable findings; of the seven listed, only orphan worktrees qualify, which makes `N fixable: run do-work cleanup` a one-case feature. Stranded finished REQs are `actions/forensics.md` Check 9's existing definition and are exactly what cleanup Pass 0 sweeps, so the probe reuses a definition rather than inventing one and gives the routing a second real case.
- **D-05 — `next-req` consults both the filename number and the frontmatter `id:`, taking the higher.** DECIDE & STATE. `actions/capture.md`'s rule says to scan `REQ-*.md` filenames, and filenames alone would be the cheaper read. But a file renamed away from its id still owns that id's number, and handing it out again would manufacture exactly the duplicate `verify` exists to detect. The extra read uses the existing lenient frontmatter parser, so it adds no new parsing path and stays in the millisecond range.

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/allocate.go` (new)
- `tools/queue-kanban/release.go` (new)
- `tools/queue-kanban/verify.go` (new)
- `tools/queue-kanban/allocate_test.go` (new)
- `tools/queue-kanban/release_test.go` (new)
- `tools/queue-kanban/verify_test.go` (new)
- `tools/queue-kanban/main.go` (modified)
- `tools/queue-kanban/prime-do-kanban.md` (modified)
- `actions/capture.md` (modified)
- `actions/work.md` (modified)
- `actions/forensics.md` (modified)
- `docs/forensics-guide.md` (modified)
- `CLAUDE.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Added three subcommands to the shipped board tool. `next-req` prints the next free REQ number by running `actions/capture.md`'s own scan on top of the existing `enumerateDoWorkTree` walk (no second parser), consulting both filename and frontmatter id and tolerating gaps. `next-version <patch|minor|major>` rewrites the single `**Current version**: X.Y.Z` line, reads it back to confirm the write landed, prints the new number, and refuses without writing when the repo keeps no such line — the bump size is a required positional argument, never inferred. `verify` is read-only and runs eight probes: version/changelog agreement, changelog version ordering, entry-title reuse, duplicate REQ ids, a checkpoint naming a REQ that no longer exists, untrustworthy `claimed_at` stamps (stale past three hours, unparseable, future-dated past the shared 2-minute skew allowance, or absent), finished REQs stranded outside the archive, and `worktree-agent-*` leftovers plus any such worktree carrying uncommitted `do-work/` changes. Each finding carries a `Fixable` flag set only for what `do-work cleanup` can mechanically resolve, and the report ends with `N fixable: run do-work cleanup`; a probe that cannot run is reported as skipped rather than passing silently. Nothing in the tool writes `CHANGELOG.md`. Wired into the three existing actions as optional accelerators — each stating its missing-`go` fallback — plus a fourteenth forensics check and its guide row, `CLAUDE.md`'s ONE→two write-surface correction, the tool's prime file, and 21 new Go tests plus 8 contract-suite assertions.

## Qualification

Passed — 14 files verified in the diff, 12 acceptance criteria traced, P-A-U confirmed against the actual diff.

- **Files exist / show in diff:** all 14 present; the 6 `(new)` files are all substantive (109–487 lines of source, 162–428 of tests) — none is a placeholder. Every new source file is reachable: `allocate.go`, `release.go`, and `verify.go` are each called from a `main.go` subcommand arm and from their own tests, so the unreferenced-`(new)`-file WARN class does not apply.
- **Substantive:** `go vet` and `gofmt` clean; no `panic`, no debug prints, no `TODO`. `main.go`'s three new arms match the three new runners, and the `want …` error list was updated with them (a pinned assertion, since a stale list lies about what exists).
- **Requirements traced:** 1→`nextRequestNumber` + 9 test cases; 2→`allocateNextVersion` with the positional bump size and `TestBumpSemanticVersionRejectsAnUnnamedBumpSize`; 3→eight probes in `verify.go`, each with a test; 4→`TestAllocateNextVersionWritesNothingElse` plus a suite assertion that `release.go` never writes CHANGELOG; 5→`Fixable`/`FixableCount` + `TestVerifyReportRoutesFixableFindingsToCleanup`; 6→the read-back in `allocateNextVersion` + `TestAllocateNextVersionWritesBackAndConfirms`; 7→no gap logic exists anywhere (asserted by the gap test returning 71, not 2); 8→three call sites, each pinned by a suite assertion; 9→each call site's fallback sentence, also pinned; 10→`CLAUDE.md` two-write-surface sentence + assertion; the write-surface constraint→`TestVerifyWritesNothing` walks the whole fixture tree before and after.
- **Flowing (not hollow):** every probe was exercised against a fixture *and* the whole tool was run against this live repo — `next-req` printed 75, `verify` exited 0 clean after the D-02 correction and exited 1 naming the real offender before it.
- **Scope drift:** one, self-inflicted and caught mechanically. `docs/forensics-guide.md` was decided on and logged as D-03 at edit time, but the `## Scope` list was not updated in the same pass — `tools/checks/scope-drift.sh` flagged it as touched-but-undeclared (exit 1). The declaration was then completed in `## Scope` and `write_set`, and the check re-run clean. Worth recording rather than quietly fixing: writing the decision down is not the same as declaring the file, and only the script noticed the difference. `model.go`, `walk.go`, `frontmatter.go`, `testing.go`, `serve.go`, `generate.go`, `actions/cleanup.md`, and `SKILL.md` are untouched, as declared.
- **Contamination check (Step 10):** REQ-071 touched `actions/work-reference.md`, `actions/work.md`, `_dev/tests/contract-regressions.sh`, `docs/work-guide.md`. The two files this REQ shares with it — `actions/work.md` and `_dev/tests/contract-regressions.sh` — are both *expected* overlaps declared in this REQ's own `write_set` at capture (requirement 8 names the Step 9 call site; the suite is where every contract assertion goes). Different sections in both cases: Step 9's commit ritual here vs. Step 1's recovery there, and a separate assertion block. Not contamination.

## Testing

**Tests run:** `go test ./...` in `tools/queue-kanban/` and `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ All passing (Go suite ok; contract suite green including the `record-commit-hash`, `blanked-req-scan`, and `update-script-behavior` sub-probes)

**Red-green validation:**
- 21 new Go tests: ✗ before implementation → ✓ after. The RED was total and unambiguous — every new symbol was undefined, so the package did not compile (`undefined: nextRequestNumber`, `undefined: VerifyReport`, and 9 more before the compiler gave up).
- 8 new contract-suite assertions: ✗ before the prose call sites existed → ✓ after.
- Live-repo acceptance run (the REQ's GREEN condition): `./queue-kanban next-req --repo-root ../..` printed **75**. The REQ predicted 74 — written when REQ-073 was the highest id; REQ-074 was queued by REQ-071's Discovered Task since, so 75 is the correct answer and the prediction is stale, not the tool.
- `./queue-kanban verify --repo-root ../..` → **exit 1** on the first run, naming `version-not-ahead-of-changelog: current version 0.164.0 is not strictly greater than the newest CHANGELOG.md entry 0.164.0`. That was the probe misreading a healthy tree, not a repo defect — see `## Decisions` D-02. After the correction: **exit 0, no findings.**
- Hand-planted failures, each confirmed to fire and name the offender: a version file ahead of the newest entry, a version file behind it, a duplicated changelog version, a reused entry title, two `REQ-071-*.md` files under different slugs, a checkpoint naming a nonexistent REQ, four flavors of untrustworthy `claimed_at` (stale / future / unparseable / absent) alongside a fresh claim that correctly stayed unflagged, and a finished REQ left in `queue/`.

**New tests added:**
- `tools/queue-kanban/allocate_test.go` — 6 table cases (highest-id, empty, gapped, spanning queue/working/nested-archive, wide zero-padding, pruned subtrees) + frontmatter-id-outranks-filename + read-only-toward-the-queue
- `tools/queue-kanban/release_test.go` — bump table (5), unnamed-bump-size rejection, prefix anchoring, missing-line error, write-back-and-confirm, writes-nothing-else, changelog entry parsing, numeric version comparison
- `tools/queue-kanban/verify_test.go` — clean tree, version ahead, version behind, duplicate changelog version, reused title, duplicate REQ ids, checkpoint ghost, five claim shapes, stranded finished REQ, fixable routing, whole-tree read-only snapshot, foreign-changelog-format skip
- `_dev/tests/contract-regressions.sh` — never-write-CHANGELOG, verify-stays-read-only, complete `want …` subcommand list, three call sites × (invocation + fallback), `CLAUDE.md` two-write-surface sentence, forensics-guide row

**Existing tests updated (cross-REQ impact):** none — every change is additive. No prior Go test or contract assertion was modified.

**Pre-flight baseline:** clean (working tree clean outside `do-work/`, both suites green before any edit).

## Review

**Overall: 91%**

**Pipeline mode — Approve (self-review).** Acceptance: **Pass**.

**Requirements Compliance: 100%** — all ten Detailed Requirements delivered, with one corrected rather than followed literally (requirement 3's first bullet; `## Decisions` D-02 records why the literal reading produces a probe that fires on every healthy repo, and how the intent is better served by the pair that replaced it). Every Constraint holds: the tool ships inside `tools/queue-kanban/` with no independent versioning; no locking scheme was built and the prose says plainly that duplicates are possible, not impossible; the write surface widened by exactly one file outside `do-work/`; the REQ is framed and justified on the two hand-checked pre-commit rules, not on parallel builds; the compiled binary stays gitignored and unstaged.

**Code Quality: 92%** — three focused files with one job each, reusing the existing walk, the existing lenient frontmatter reader, the existing duplicate-id resolution, the existing skew constant, and the established `exec.Command("git", "-C", …)` shape. No second parser, no new skew constant, no new timestamp layouts. Naming follows the repo's two-word convention throughout (`highestNumberInUse`, `changelogEntryHeadingPattern`, `worktreeDirtyQueueState`). Doc comments state *why*, not *what* — the `next-version` anchor comment and the version-probe comment both carry the reasoning a later reader would otherwise re-derive. Deductions: `verify.go` at 487 lines is the largest new file and its eight probes are eight top-level functions sharing one report — fine now, a split candidate at twelve; and `runVerifyProbes` builds a full board (which shells out to git for commit dates) to answer probes that need only frontmatter, so it is heavier than `next-req`'s millisecond path. Neither is wrong, both are worth knowing.

**Test Adequacy: 95%** — 21 Go tests plus 8 contract assertions, and the RED was genuine and total (the package did not compile). Coverage goes past the REQ's listed cases: the read-only guarantees are asserted by snapshotting the whole fixture tree before and after, not by sampling a file; the claim probe includes a fresh claim that must *stay* unflagged, which is the assertion that would catch an over-eager threshold; the foreign-changelog case asserts the skip is reported *by name*, so a silently-skipped invariant cannot pass as a clean one. Deducted because two probes are only covered by fixture-free paths: the worktree probes skip when git cannot answer, so `listWorktreeAgentWorktrees` / `worktreeDirtyQueueState` are exercised on this live repo (correctly finding nothing) but have no fixture that plants a `worktree-agent-*` leftover — a temp-repo test with a real `git worktree add` would close it, and is noted under suggested testing rather than built.

**Scope Discipline: 88%** — held to the plan, with one self-inflicted miss the tooling caught: `docs/forensics-guide.md` was decided on and logged as D-03 at edit time but not added to `## Scope` in the same pass, so `tools/checks/scope-drift.sh` flagged it touched-but-undeclared. Declared and re-checked clean. Deducted honestly for that, and for `actions/forensics.md`'s read-only bullet needing an amendment I had not planned — both are small, both are recorded rather than smoothed over.

**Risk: Low.** The new code is additive; no existing path changed behavior. Three risk surfaces, each bounded: (1) `next-version` is the only writer, it refuses unless exactly one `**Current version**: ` line exists, and it read-back-confirms — a failed write reports instead of returning a number the caller would put in a changelog heading; (2) `next-req` is not atomic, which is stated in the code comment, the prime's Traps, and `actions/capture.md`'s prose rather than papered over; (3) the git shell-outs are read commands only, and every one degrades to a reported skip.

**Restatement Sweep — run, three findings.** The diff redefines two things other text restates: *how many write surfaces the tool has*, and *whether forensics modifies anything*.
- `CLAUDE.md` — the "ONE write surface" sentence, which requirement 10 named. Fixed, and generalized: it now states both surfaces, names the changelog exclusion, and says that a third surface means amending this sentence in the same commit.
- `actions/forensics.md` Core Rules — "never modifies files" became false the moment Check 14 prescribed `go build`, which writes the tool's binary. **This one was caused by my own addition**, in a file already in Scope, so it was fixed inline with a scoped carve-out naming exactly what gets written and confirming nothing in the project or `do-work/` is touched. Had I not swept, forensics would have shipped contradicting itself on its central promise.
- `actions/board.md:5` — "It writes exactly three things" (**Minor**, reported, not fixed, no follow-up). The subject is *the board*, and `do-work board` genuinely still writes only those three: `next-version` is unreachable from that action. So the sentence is accurate-as-scoped rather than stale, and the tool-wide statement now ships in `tools/queue-kanban/prime-do-kanban.md`'s Traps (`CLAUDE.md` is export-ignored, so it could not carry the shipped truth alone). Worth knowing if board.md's opening is ever rewritten to describe the binary rather than the action.
- Verified still-accurate, no change needed: `actions/capture.md`'s allocation rule (the shortcut sits beside it and defers to it), `actions/cleanup.md` (untouched; `verify` only points at it), the terminal-status readers in present-work / review-work / ai-report (they scan for the highest *completed* REQ, unrelated to allocation).

**Coding-guardrails spot-check:** Think Before Coding — five `## Decisions` entries; the one place the REQ's literal text was wrong got a full written argument (D-02) instead of a silent reinterpretation. Simplicity First — no CLI library, no config file, no gap machinery, no locking; the eight probes are eight plain functions. Surgical Changes — `model.go`, `walk.go`, `frontmatter.go`, `testing.go`, `serve.go`, and `generate.go` are untouched, as declared. Goal-Driven Execution — the tool was run against this live repo, which is what surfaced D-02; a suite-only pass would have shipped the broken probe.

**Findings**
- **Minor:** the worktree probes have no fixture that plants a real `worktree-agent-*` leftover (see Test Adequacy). They were exercised live, correctly finding nothing — which proves the happy path and not the detection.
- **Minor:** `verify` prints the `--repo-root` value verbatim in its header (`queue-kanban verify — ../..`) rather than the resolved absolute path. Pre-existing behavior of `resolveRepoRootOrDefault`, shared with every other subcommand; left alone rather than changed under this REQ.
- **Nit:** `verify.go` is the largest new file; the threshold/prompt-free probe set is the natural split point if it grows.

**Suggested additional testing**
- A temp-repo test that runs a real `git worktree add worktree-agent-REQ-999-probe`, then asserts both worktree probes fire — including the dirty-`do-work/` case, which is the one probe here that detects a rule being broken rather than a bookkeeping slip.
- Run `next-version patch` for real on the next release and confirm the read-back path behaves on a file with a pre-commit hook attached; every test here writes to a plain temp file.
- Point `verify` at a consumer repo (no `**Current version**: ` line, possibly a keep-a-changelog file) and confirm it degrades to reported skips rather than findings — the fixture covers the shape, a real install would confirm the ergonomics.

**Follow-up REQs created:** none. No Important findings survived the sweep — the one restatement that would have shipped a contradiction was inside a declared file and caused by this REQ, so it was fixed here.

## Lessons Learned

**What worked:**
- **Running the tool against the live repo before believing the test suite.** All 21 tests were green *and* the version probe was wrong — the fixtures encoded the same misreading as the code, because I wrote both from the same sentence. Pointing the binary at this repo took ten seconds and exposed it immediately. A prose-derived probe should always be run against a known-healthy tree, since "fires on a healthy repo" is the failure a fixture built from the same premise cannot see.
- **Reusing `enumerateDoWorkTree` and `Board.Warnings` instead of re-deriving.** `next-req` inherits the walk's pruning of `deliverables/`, `runs/`, and `assets/` for free — and that pruning is load-bearing: a deliverable copy named `REQ-900-*.md` under `assets/` would otherwise push allocation 800 numbers into the future. The duplicate-id probe likewise inherits tree-section precedence. Neither behavior would have occurred to me to write from scratch.

**What didn't:**
- **Taking requirement 3's "strictly greater" literally.** The rule is real (it's in `CLAUDE.md` § Before Every Commit) but it describes a *transient* state during release composition; the committed steady state is equality. Encoding a mid-process condition as a standing invariant produces a check that fails on every healthy repo — the most useless possible outcome, since it trains the user to ignore the tool. The generalizable lesson: when a requirement quotes a pre-commit rule, ask *at what moment is this true*, because a verifier runs at an arbitrary moment.
- **Assuming `actions/version.md` is a findable path.** It only exists at the repo root in a skill-development checkout. Ten minutes went into the design before noticing that a consumer install puts the skill under `.claude/skills/do-work/` while `do-work/` sits at the consumer root — so the tool's own repo-root resolution (which deliberately *skips* skill installs, per `walk.go`'s `isSkillInstallDirectory`) would never find it. Anchoring on the content marker rather than the path is what made the subcommand portable, and it was the REQ's own instruction all along.
- **Writing the decision down is not declaring the file.** D-03 recorded the intent to touch `docs/forensics-guide.md` and `## Scope` still didn't list it; only `tools/checks/scope-drift.sh` noticed. The Decisions section and the Scope declaration are read by different checks, and satisfying one does not satisfy the other.

**Worth knowing:**
- **Adding a forensics check that shells out to a compiler breaks that action's central promise.** `actions/forensics.md` opens with "never modifies files"; `go build` writes a binary. The carve-out is scoped (gitignored binary, inside the skill install, nothing in the project or `do-work/`), but any future check that runs a build needs the same treatment — or the promise quietly stops being true.
- **`CLAUDE.md` is export-ignored, so it cannot be the only home for a shipped contract.** The complete two-write-surface statement lives there *and* in `tools/queue-kanban/prime-do-kanban.md`'s Traps, because only the second one reaches a consumer. Anything a downstream reader needs has to live in a shipped file.
- **The version probe's direction is diagnostic, not decoration.** "Ahead" means a bump landed without its changelog entry (normal mid-release, a real finding afterwards); "behind" means an entry carries a version the version file never received. Collapsing them into one "mismatch" message would throw away the half of the signal that tells you which file to fix.

## Orientation

Now the release ritual's two hand-checked consistency rules are machine-checked, and REQ and version numbers can be allocated instead of eyeballed — `tools/queue-kanban/`, the shipped board tool, via three new subcommands (`next-req`, `next-version`, `verify`) wired into capture, the Step 9 commit ritual, and forensics as Check 14.

`[MAP CHANGED]` — the board tool is no longer a read-only viewer with one narrow testing write. It now has **two** write surfaces (the testing fields inside `do-work/`, and one version line outside it) and a release-verification role, which is a change in what the tool *is*. `CLAUDE.md`'s write-surface contract and `tools/queue-kanban/prime-do-kanban.md` both record it; the changelog is explicitly excluded from both surfaces and stays a human write.

Prime staleness spot-check on the declared prime (`tools/queue-kanban/prime-do-kanban.md`, per `actions/prime.md` Step 2): every path it references still exists (`main.go`, `walk.go`, `model.go`, `generate.go`, `serve.go`, `web/`, `go.mod`), and it was updated in this REQ rather than left stale — header subcommand list, four `## Read first` entries for the new files, four `## Traps` covering the never-write-changelog boundary, the version-probe polarity, the bump-size rule, and non-atomic allocation.
