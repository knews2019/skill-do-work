---
id: REQ-098
title: "Verify probes: assigned-elsewhere-claimed-here and ur-archived-with-live-member"
status: completed
created_at: 2026-08-04T19:44:17Z
claimed_at: 2026-08-04T21:12:28Z
completed_at: 2026-08-04T21:20:00Z
commit: 47cd408
kb_status: pending
user_request: UR-018
domain: backend
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-097]
maintenance: false
related: [REQ-097]
batch: parallel-building
write_set: [tools/queue-kanban/verify.go, tools/queue-kanban/verify_test.go, actions/forensics.md, CLAUDE.md]
---

# Verify Probes: Assignment Drift and UR Closure Drift

## What

Two new read-only probes in `queue-kanban verify`, extending the report-and-route contract (verify never repairs — fixes belong to `actions/cleanup.md`):

1. `assigned-elsewhere-claimed-here` — a REQ carrying `assigned_to` is sitting in `do-work/working/` (someone claimed work earmarked for another session without clearing the field).
2. `ur-archived-with-live-member` — an archived UR has a member REQ (by `user_request:` scan) still in `queue/` or `working/` (a silent-merge class also reachable from a botched cleanup).

## Detailed Requirements

- Implement in `tools/queue-kanban/verify.go`, ~30 lines each, following the existing probe structure (finding code, message, routing suggestion). Tests in `verify_test.go` with fixture trees for both the firing and non-firing case.
- Probe 2's membership scan uses `user_request:` frontmatter across `queue/`, `working/`, `archive/` root and `archive/UR-NNN/` — the same closure predicate `actions/work.md` Step 8 evaluates (the UR `requests:` array is capture-time-only, never the closure predicate).
- Route findings to `actions/cleanup.md` in the report text, matching existing probes' phrasing.
- `go test ./...` green.

## Red-Green Proof

**RED prompt/case:** A fixture with an `assigned_to` REQ in `working/`, and one with an archived UR whose member REQ sits in `queue/` — `queue-kanban verify` today reports neither.
**Why RED now:** No probe covers assignment drift (the field is new) or archived-UR/live-member drift.
**GREEN when:** Both probes fire on their fixtures, stay silent on clean trees, and the full Go test suite passes.
**Validation:** User confirmed (approved plan, Phase 2 item 7).

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 2).

---
*Source: approved plan, Phase 2*

## Triage

**Route: B** - Medium

**Reasoning:** Both probes, their firing conditions, their routing target and their membership predicate are specified. What needed discovery was the existing probe structure (finding shape, category constants, the runner's call list), whether a `UserRequestTicket` carries a tree section (it does not — path is the only evidence), and which prose enumerates the probe set. `tdd: true`, so five tests came first.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided, TDD

*Skipped by work action*

## Exploration

**Probe structure:** each probe is an `append<Name>Findings(&report, …)` function called from `runVerifyProbes`, appending `VerifyFinding{Category, Detail, Fixable, Remedy}`. Category constants are a single `const` block. `renderVerifyReport` is fully generic — it iterates findings and prints `Category`, `[fixable]`, `Detail` and `Remedy`, so a new probe needs **no** renderer change. `FixableCount` drives the trailing `run do-work cleanup` line.

**`Fixable` is a claim about cleanup, not about severity.** `appendStrandedFinishedFindings` sets it (Pass 0 sweeps mechanically); the version-mismatch and unmerged-worktree findings do not. Both new probes are **not** fixable: whose claim stands, and whether to un-archive a UR or force-resolve its live member, are decisions with different consequences.

**`UserRequestTicket` has `FilePath` but no `TreeSection`** (unlike `RequestTicket`), so archived-ness has to come from the path. Matched on the `/do-work/archive/` segment with separators on both sides, so a project directory merely named `archive` elsewhere cannot satisfy it.

**Membership is already computed the right way.** `linkRequestsToUserRequests` fills `UserRequestTicket.RequestIds` from each REQ's `user_request:` frontmatter — which is exactly Step 8's closure predicate, and exactly *not* the UR's own `requests:` array. So probe 2 needed no new scan, only the discipline to use `RequestIds` and not parse the array. A test was written specifically to pin that, with a fixture where the two disagree.

**Prose enumerations of the probe set:** `actions/forensics.md` Check 14 carries a probe table (the only place the set is listed), and `CLAUDE.md` describes what `verify` covers. Both go stale on a new probe — the Closed-Enumerations pattern.

## Scope

**Files I will touch:**
- `tools/queue-kanban/verify_test.go` — five tests, written first
- `tools/queue-kanban/verify.go` — two category constants, two probes, one path helper, two runner calls
- `actions/forensics.md` — Check 14's probe table (two rows) plus a note that the table is not a contract
- `CLAUDE.md` — the one-line description of what `verify` covers

**Acceptance criteria (restated from REQ):**
- [ ] `assigned-elsewhere-claimed-here` fires on an `assigned_to` REQ in `working/`, silent otherwise
- [ ] `ur-archived-with-live-member` fires on an archived UR with a member in `queue/`/`working/`, silent otherwise
- [ ] Probe 2's membership uses the `user_request:` frontmatter scan, never the UR's `requests:` array
- [ ] Both follow the existing probe structure (category, message, routing suggestion)
- [ ] Both route to `actions/cleanup.md` in the report text, matching existing phrasing
- [ ] Tests cover the firing and non-firing case for both
- [ ] `go test ./...` green; verify stays read-only and never repairs

## Pre-Flight

- **WARN — baseline suite red before any change:** the same 8 `chmod 500`-versus-root failures inherited by every REQ in this batch.
- `go test ./...` green at claim time.

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/verify_test.go` (modified) — five tests: probe 1 firing (asserting the assignee is named verbatim and the finding is **not** fixable) and not firing on the normal earmarked-and-left-alone case; probe 2 firing (asserting both ids are named, the resolved member is **not** listed, and it is not fixable), not firing on a properly-closed UR, and `TestVerifyLiveMemberProbeScansUserRequestFrontmatterNotTheUrArray` — a fixture whose `requests:` array omits the live member, which a probe reading the array would pass silently.
- `tools/queue-kanban/verify.go` (modified) — `verifyCategoryAssignedElsewhereClaimedHere` and `verifyCategoryArchivedUserRequestLiveMember`; `appendAssignedElsewhereFindings` and `appendArchivedUserRequestLiveMemberFindings`, both wired into `runVerifyProbes` between the stranded-finished and worktree probes; `isArchivedUserRequestPath` helper.
- `actions/forensics.md` (modified) — two rows in Check 14's probe table, each naming the definition it reuses, plus a sentence marking the table as the set-as-it-stands rather than a contract, so the next probe does not silently make it wrong.
- `CLAUDE.md` (modified) — `verify` now described as covering assignment and UR-closure invariants alongside the release, queue and worktree ones.

**What was done:** Added the two read-only probes. Neither is fixable, both name their remedy as something cleanup *asks* about, and the renderer needed no change because it was already generic.

## Testing

### Red-green validation (`tdd: true`)

**RED**, from the five tests written first:

```
$ go test ./...
./verify_test.go:793:49: undefined: verifyCategoryAssignedElsewhereClaimedHere
./verify_test.go:818:41: undefined: verifyCategoryAssignedElsewhereClaimedHere
./verify_test.go:844:51: undefined: verifyCategoryArchivedUserRequestLiveMember
./verify_test.go:874:41: undefined: verifyCategoryArchivedUserRequestLiveMember
./verify_test.go:901:51: undefined: verifyCategoryArchivedUserRequestLiveMember
FAIL	github.com/knews2019/skill-do-work/queue-kanban [build failed]
```

**GREEN**, after the probes:

```
$ gofmt -w verify.go && go test ./...
ok  	github.com/knews2019/skill-do-work/queue-kanban	3.994s
```

This matches the REQ's `## Red-Green Proof` — its RED is "a fixture with an `assigned_to` REQ in `working/`, and one with an archived UR whose member REQ sits in `queue/` — verify today reports neither" — and adds three cases the proof did not name: both non-firing cases, and the array-versus-frontmatter divergence.

### Verify stays read-only

`TestVerifyWritesNothing` (pre-existing) passes unchanged, which is the assertion that matters most for a probe addition: neither new probe writes, and neither is marked `Fixable`, so `FixableCount` — and therefore the `run do-work cleanup` line — is unaffected by them.

### Against this repo's real tree

```
$ (cd tools/queue-kanban && go build -o queue-kanban .) && tools/queue-kanban/queue-kanban verify --repo-root /home/user/skill-do-work
queue-kanban verify — /home/user/skill-do-work
  OK: no findings
```

Clean, which is the correct answer: UR-018 is still in `do-work/user-requests/` with live members (so probe 2 must stay silent — it fires only on an *archived* UR), and `do-work/working/` is empty at the time of the run.

### Contract suite

```
$ bash _dev/tests/contract-regressions.sh 2>&1 | grep -c '^FAIL'
8
```

The pre-existing eight, name-for-name.

## Lessons Learned

**What worked:**
- Writing the array-versus-frontmatter divergence test *before* the probe. `RequestIds` is already the right data, so the probe would have been correct by accident; the test is what makes it correct on purpose, and it is the one that would catch a future "optimization" to read the `requests:` array instead.
- Asserting the **negative** in probe 2's firing test — that the terminally-resolved sibling is *not* named. A probe that lists every member of the UR would have passed a test that only checked for the live one.

**What didn't:**
- Rebuilding the binary from inside `tools/queue-kanban` and then invoking it by a repo-relative path in the same command line — the `cd` had already moved the shell, so the invocation failed with `No such file or directory` and read like a build failure. Wrap the build in a subshell (`(cd … && go build …)`) so the outer shell never moves. Same class of stale/wrong-path confusion REQ-097 hit twice.

**Worth knowing:**
- `renderVerifyReport` is generic over findings, so a new probe is three edits (constant, function, runner call) and never a renderer change. The temptation to add per-category formatting is what would break that.
- `Fixable` means *cleanup can resolve this mechanically*, not *this is minor*. Setting it on a human-decision finding would make `do-work cleanup` advertise a repair it must not perform unasked.
- `UserRequestTicket` carries no `TreeSection`. Anything that needs a UR's location reads `FilePath` — and should match `/do-work/archive/` with both separators, not the bare word.

## Orientation

`do-work verify` now catches two drift shapes the claim-anywhere model makes reachable: a REQ being built here while still earmarked for somewhere else, and a UR archived while one of its members is still live. Both are read-only findings that route to `do-work cleanup`, and neither is marked mechanically fixable, because both resolutions are human calls. Lives in `tools/queue-kanban/verify.go` with the user-facing description in `actions/forensics.md` Check 14. `prime_files` is empty; `tools/queue-kanban/prime-do-kanban.md`'s referenced paths were spot-checked and all still exist.

## Qualification

**Passed** — 4 files verified, 7 acceptance criteria traced, red-green confirmed from actual `go test` output.

- **Substantive:** ~60 lines of Go across two probes and a helper, five tests with real fixture trees, two prose rows.
- **Wired:** both probes are called from `runVerifyProbes`; verified by running the built binary against this repo rather than by reading the call site.
- **Flowing, not hollow:** each probe's firing test asserts on the `Detail` string's content, so a probe that fired with an empty or generic message would fail.
- **Requirements traced:** all seven criteria map to a test or a diff hunk.

## Review

**Overall: 95%** | 2026-08-04T21:20:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — both probes fire on their fixtures and stay silent on clean trees, the membership predicate is pinned by a divergence test, the full Go suite is green, and verify remains read-only.
**Suggested testing:** 1 item
**Follow-ups created:** None

**Requirements checklist:** all seven `## Scope` criteria delivered; evidence in `## Testing`.

**Minor:**
- Scope grew by two files beyond the REQ's `write_set` (`actions/forensics.md`, `CLAUDE.md`), both declared before editing. Both are enumerations of the probe set that a new probe silently falsifies — the Closed-Enumerations pattern — so leaving them would have shipped a probe table that reads as complete and is not. The forensics row additions also came with a note marking the table as non-authoritative, which is the durable fix rather than the per-probe one.

**Scope drift:** none against the declaration — four files declared, four touched.

**Restatement sweep (MUST):** run. The diff adds two probes, so the sweep looks for text that restates *what verify covers*: `actions/forensics.md` Check 14's table (the only enumeration of the set) and `CLAUDE.md`'s description. Both updated. Nothing else in the shipped tree enumerates probe categories — checked by grepping the category strings across `actions/`, `docs/`, and `tools/queue-kanban/*.md`, which returned nothing, confirming the tool's output is the only other place they appear.

**Suggested additional testing:**
- A probe-2 fixture where the archived UR's live member sits in `working/` *and* carries `assigned_to`, so both new probes fire on one tree — confirms they compose rather than shadowing each other in the report.

*Reviewed by review-work action (pipeline mode, in-session)*
