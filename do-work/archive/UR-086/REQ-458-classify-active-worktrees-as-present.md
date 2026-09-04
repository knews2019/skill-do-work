---
id: REQ-458
title: 'Addendum: classify active worktrees as present and non-fixable'
status: completed
created_at: 2026-08-31T21:38:14Z
user_request: UR-086
addendum_to: REQ-083
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-09-03T00:59:09Z
route: B
write_set:
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
  - skills/do-work/actions/forensics.md
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-09-03T00:59:54Z
  basis:
    - Route B
    - 3-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
    - cross-route regression gates
completed_at: 2026-09-03T02:30:00Z
commit: ea2bab0
release_at: 2026-09-03T02:30:00Z
---

# Addendum: Classify Active Worktrees as Present and Non-Fixable

## What

Correct REQ-083 (Verify reports every builder worktree as a fixable orphan, including active and unmerged ones) so a branch being merged into the integration branch is not, by itself, enough to call its worktree a leftover or mechanically fixable. A dirty worktree or a worktree belonging to an unfinished run must be reported as present and non-fixable; only clean merged residue from finished work may be reported as a fixable leftover.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read both `prime_files` and both `lessons-*` satellites, `CLAUDE.md` § Kanban Board Write Surfaces, and `actions/cleanup.md` → Pass 5. Approach: keep `classifyWorktreeMergeState` as-is and layer a `worktreeLeftoverDisposition` over it, fed by two repository reads — worktree dirtiness and whether the leftover's REQ id is still in `do-work/working/` — so merged-ness becomes necessary but not sufficient for `Fixable`.
- [x] **[APPLY]:** Two files. `skills/do-work/actions/forensics.md` was declared in Scope and deliberately not written — see the Qualification note and D-05.
- [x] **[UNIFY]:** `git diff --stat` → 2 files, +357/-40. Orchestrator independently re-ran `go build ./...`, `go vet ./...` (clean) and `gofmt -l .` (silent), and confirmed the five relevant tests pass, `TestVerifyWritesNothing` among them. Read the added declarations and functions in the diff — no debug artifacts, and the only subprocess added is a read-only `git status --porcelain`.

## Context

This corrects the incomplete merged-state classification introduced by REQ-083. The board screenshot showed two `merged-worktree-leftover` findings marked `cleanup can fix`:

- REQ-412 had a merged branch tip while its builder worktree contained uncommitted implementation changes. The current merge-base-only classifier therefore claimed “nothing is lost” even though cleanup's non-forced removal would refuse the dirty worktree.
- REQ-436 remained `claimed` before review and possible remediation had completed. Its clean merged worktree still belonged to the active pipeline; normal worktree cleanup occurs only after the REQ reaches its final path.

The accepted validation finding recorded `Surface-cost: N/A` because this is a direct classification correction, not a new guard, retry, fallback, or warning apparatus.

## Prior Implementation

REQ-083 added `classifyWorktreeMergeState` and `routeWorktreeLeftover`, split the old orphan category into merged, unmerged, and undetermined categories, and set `Fixable: true` for the merged branch state. Its implementation and fixture tests now live under `skills/do-work-board/tools/queue-kanban/verify.go` and `verify_test.go`; the accompanying forensics description lives under `skills/do-work/actions/forensics.md`. The recorded implementation commit is `f6c1514`.

## Requirements

- Preserve the existing no-liveness-signal decision: do not add a heartbeat, lock, PID probe, mtime heuristic, claim registry, or other process-liveness guess.
- Distinguish ordinary worktree dirtiness and unfinished pipeline state from clean merged residue using evidence the repository already records.
- Report dirty or unfinished-run builder worktrees as present and non-fixable. Do not describe them as leftovers, advertise `cleanup can fix`, or claim that nothing can be lost.
- Continue reporting genuinely finished, clean, merged residue as a mechanically fixable leftover.
- Keep `verify` read-only and preserve the current protection for developer-owned worktrees outside the `worktree-agent-*` naming convention.
- Keep all rendered category/remedy text and any user-facing verify documentation consistent with the corrected classifier.

## Red-Green Proof

**RED prompt/case:** Extend the real-Git worktree fixture in `skills/do-work-board/tools/queue-kanban/verify_test.go` with (1) a branch already merged into the integration branch whose worktree has an uncommitted ordinary source-file change and (2) a clean merged worktree whose matching REQ is still `claimed` before review/remediation finishes. Run the verify probes and inspect both findings.

**Why RED now:** Both cases route solely from `git merge-base --is-ancestor`, so each currently becomes `merged-worktree-leftover [fixable]` with the claim that cleanup is mechanical and nothing is lost.

**GREEN when:** Both regression cases report the worktree as present and non-fixable, with no `cleanup can fix` marker and no “nothing is lost” claim; a separate clean merged worktree from finished work remains a fixable leftover. The regression test must close both REQ-412 and REQ-436 instances without introducing a liveness signal.

**Validation:** User confirmed through the accepted `do-work validate-feedback` finding and this capture request.

## Assets

- `do-work/user-requests/UR-086/assets/REQ-458-screenshot-1-active-worktrees-labelled-leftovers.png` — queue board generated during the active run. Its Verify Findings strip shows two cards, for REQ-412 and REQ-436, both labelled `MERGED-WORKTREE-LEFTOVER`, both marked `cleanup can fix`, and both stating that the branch is contained in `HEAD` so nothing is lost.

---
*Source: user-approved `do-work validate-feedback` finding; full verbatim input in `do-work/user-requests/UR-086/input.md`.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The correction is precisely specified and `## Prior Implementation` names the exact functions, but the two evidence sources the fix must consult — worktree dirtiness and the REQ's own pipeline state — are not currently reachable from `appendWorktreeFindings`, so where they come from had to be discovered.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

`classifyWorktreeMergeState` (`verify.go:819`) runs exactly one probe — `git -C repoRoot merge-base --is-ancestor <branch> HEAD` — and collapses everything else into three states. `routeWorktreeLeftover` (`verify.go:851`) then maps `worktreeMergeStateMerged` to `Fixable: true` with the remedy "the branch is already contained in HEAD, so nothing is lost". That single ancestry bit is the whole basis for both claims, which is exactly the defect: ancestry says the *commits* are safe, never that the *worktree* is.

Two facts the repository already records, neither of them a liveness signal:

1. **Worktree dirtiness** — `appendWorktreeFindings` already holds `worktreePath` from `listWorktreeAgentWorktrees` (it uses it only for `locationDetail`). `git -C <worktreePath> status --porcelain --untracked-files=all` answers whether uncommitted work is present. This is the REQ-412 case: cleanup Pass 5's non-forced `git worktree remove` would itself refuse this worktree, so calling it mechanically fixable contradicts the very command the remedy names.
2. **Unfinished pipeline state** — a `worktree-agent-REQ-NNN-*` name carries its REQ id, and a REQ still in flight is exactly the one sitting in `do-work/working/`. That is the REQ-436 case: clean, merged, and still owned by a run that has not reached review or remediation.

Both are ordinary repository reads, so the no-liveness-signal constraint from REQ-073 holds: neither asks whether a process is alive, only what the repository already says.

`Fixable`'s doc comment (`verify.go:66`) defines it as "`do-work cleanup` can mechanically resolve it", and `routeWorktreeLeftover`'s own comment says anything landing on Pass 5's consent-gated path must not be advertised otherwise. The fix is to make the merged branch state necessary but not sufficient.

*Generated in-session (single-pass discovery)*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — make merged-ness necessary but not sufficient; add the dirtiness and in-flight evidence and route both to present-and-non-fixable
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modify) — real-Git fixtures for the REQ-412 and REQ-436 cases plus the still-fixable clean finished case
- `skills/do-work/actions/forensics.md` (modify) — keep the rendered category and remedy text consistent with the corrected classifier

**Files I will NOT touch:** `skills/do-work/actions/cleanup.md` (Pass 5's consent-gated behavior is already correct — this REQ stops verify from contradicting it), and anything adding a lock, heartbeat, PID probe, or mtime heuristic.

**Acceptance criteria (restated from REQ):**
- [x] No heartbeat, lock, PID probe, mtime heuristic, or claim registry is introduced
- [x] Ordinary worktree dirtiness and unfinished pipeline state are distinguished from clean merged residue using evidence the repository already records
- [x] A dirty or unfinished-run builder worktree is reported as present and non-fixable, with no `cleanup can fix` marker and no "nothing is lost" claim
- [x] Genuinely finished, clean, merged residue is still reported as a mechanically fixable leftover
- [x] `verify` stays read-only and still protects developer-owned worktrees outside the `worktree-agent-*` convention
- [x] Rendered category/remedy text and user-facing verify documentation match the corrected classifier

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified)

**What was done:** `classifyWorktreeMergeState` is unchanged and still answers only the ancestry question; a new `worktreeLeftoverDisposition` layer sits over it and makes the merged branch state necessary but not sufficient for `Fixable`. Two repository reads feed it: `worktreeHasUncommittedWork` (`git -C <worktreePath> status --porcelain --untracked-files=all`, returning an error rather than a silent "clean") and `classifyRequestPipelineState` (the leftover name's `REQ-NNN` id resolved through the board into `InFlight` / `Settled` / `AbsentFromBoard`). A merged leftover that is dirty, still in flight, or unreadable now reports as one of three new non-fixable categories — `worktree-present-uncommitted-work`, `worktree-present-run-in-flight`, `worktree-present-state-unknown` — none of which uses the words "leftover", "cleanup can fix", or "nothing is lost". Only a merged, clean, no-longer-in-`working/` leftover keeps `Fixable: true`, and its remedy now states all three conditions rather than the ancestry one alone. The `unmerged` and `undetermined` routing and remedies are byte-unchanged, as is the developer-owned-worktree exclusion.

## Decisions

- **D-01** — Reuse `board.RequestsById` + `TreeSection` for the in-flight read rather than a second filesystem scan, mirroring `appendStrandedFinishedFindings`. **DECIDE & STATE.** Cost is one new `board` parameter on `appendWorktreeFindings`, whose only caller already had the board built.
- **D-02** — Precedence among merged sub-states is in-flight → dirty → unknown → residue. **DECIDE & STATE.** A builder's worktree being dirty during its own run is expected, so "belongs to an unfinished run" is the fact that decides what to do. One finding per leftover name is preserved, so existing per-category counts and `FixableCount` assertions keep their meaning.
- **D-03** — `worktreeHasUncommittedWork` probes the whole worktree, not the `do-work/` subset `worktreeDirtyQueueState` asks about. **DECIDE & STATE.** They answer different questions: that one asks whether "state stays home" was broken, this one asks whether Pass 5's non-forced `git worktree remove` would refuse — and it refuses for dirt anywhere. Both are kept.
- **D-04** — Added a fourth disposition, `worktreeLeftoverStateUnknown`, covering a failed `git status` probe and a name carrying no `REQ-NNN` id. **DECIDE & STATE.** Both report present, non-fixable, plus a `SkippedProbes` line naming the failed read. This makes the fail-safe constraint structural rather than an inline nil-means-clean.
- **D-05** — `skills/do-work/actions/forensics.md` was declared in Scope and deliberately left unchanged. **ESCALATE.** Check 14 has no per-category description to make stale: line 80 states the board-output mapping is "keyed on the output class, not a hand-maintained category list, so a new verify category inherits it immediately", and line 78's only `[fixable]` prose is pass-through ("Preserve each emitted `[fixable]` … classification exactly"). Verified against the file, not taken on the builder's word. **Value:** the three new categories and the new skipped-probe line inherit their rows automatically, and the file keeps the property it advertises. **Risk:** an unused Scope declaration is scope drift, so the acceptance criterion about user-facing verify documentation is met by an unchanged file rather than an edit; if a future reader expects a per-category row there, they will not find one. Fully reversible — adding a row later costs nothing but reintroduces the hand-maintained list the file says it does not keep.
- **D-06** — The still-fixable residue remedy keeps "nothing is lost" but now names all three conditions: the branch is contained in HEAD, the worktree is clean, and its REQ has left `do-work/working/`. **DECIDE & STATE.** The claim is now true of everything it is asserted about.

## Discovered Tasks

- `worktreeDirtyQueueState` (`verify.go`) still returns `nil` when its `git status` fails, silently reporting "no queue-state writes found" for a probe that never ran — the same shape this REQ just fixed one level up, and the shape `VerifyReport`'s own doc comment argues against. `worktreeHasUncommittedWork` is the pattern to follow (error return plus a `SkippedProbes` line).
- `appendWorktreeFindings` was already flagged in the REQ-084 lesson as "the file's longest function — the next change here should extract, not append". This REQ added a parameter and a branch rather than extracting, because extraction is a refactor outside this write set. It is now roughly 90 lines running four sub-probes per leftover.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md:28` (REQ-083's entry) describes the superseded three-state model. The satellite is append-only, so the correction belongs in this REQ's own lesson bullet rather than an edit to that line.

## Qualification

**Passed with one declared-but-untouched scope drift** — 2 files verified, 6 requirements traced, P-A-U confirmed.

Mechanical: `tools/checks/qualify.sh` → `OK: mechanical qualification passed`.

`tools/checks/scope-drift.sh` → exit 1, `declared in ## Scope but never touched: skills/do-work/actions/forensics.md`. **Minor, and correct.** I checked the file rather than accepting the builder's reasoning: `forensics.md:80` states the board-output mapping is "keyed on the output class, not a hand-maintained category list, so a new verify category inherits it immediately", and `:78`'s only `[fixable]` prose is pass-through. There is no per-category row to make stale, so editing it would *add* the hand-maintained list the file says it does not keep. The Scope declaration was written before that was established; it stands as declared rather than being rewritten after the fact, with the reason recorded here and in D-05. Nothing outside the two remaining declared files was touched.

Independent (orchestrator-run, not the builder's report):
- `go build ./... && go vet ./...` clean; `gofmt -l .` printed nothing.
- Ran the three new tests plus `TestVerifyClassifiesWorktreeLeftoversByMergeState` (the REQ-083 fixture this REQ refines) and `TestVerifyWritesNothing` (the read-only invariant): all five pass.
- Read the added declarations and functions in the diff. The only subprocess introduced is `git status --porcelain --untracked-files=all`, read-only, so `verify` gains no write surface and CLAUDE.md's three-write-surface rule is untouched. No heartbeat, lock, PID probe, mtime heuristic, claim registry, or time threshold appears anywhere in the diff — the REQ-073 constraint holds.
- Requirement trace: dirtiness and in-flight state are both ordinary repository reads; the three new categories are all `Fixable: false` and none carries "leftover", "cleanup can fix", or "nothing is lost"; the clean-finished case keeps `Fixable: true`; the `worktree-agent-*` prefix guard and the `unmerged`/`undetermined` remedies are byte-unchanged.

## Testing

**Tests run:** `go build ./... && go vet ./... && gofmt -l .`; `go test ./...` for `queue-kanban`; canonical repository gate `bash _dev/tests/maintainer-verify.sh`; and the built binary run against a purpose-made five-leftover git fixture.
**Result:** ✓ `queue-kanban` suite green (`ok … 422s`). Gate exits 1 on the recorded baseline failure only — `internal/toolboxcommands` → `TestRemediationCancellationReachesMediaGitCommitAndRollback`, tracked as REQ-524 and unrelated to the board tool.

**Red-green validation:** traces the REQ's Red-Green Proof, which asked for the two mislabelled cases plus the preserved good case.
- `TestVerifyDoesNotAdvertiseADirtyMergedWorktreeAsFixable` (REQ-412's shape): ✗ `got 0 worktree-present-uncommitted-work findings, want 1`, with the report showing `! merged-worktree-leftover [fixable]` and `1 fixable: run do-work cleanup` → ✓
- `TestVerifyDoesNotAdvertiseAMergedWorktreeOfAnInFlightRequestAsFixable` (REQ-436's shape): ✗ `got 0 worktree-present-run-in-flight findings, want 1`, same `[fixable]` and "nothing is lost" remedy → ✓
- `TestVerifyStillAdvertisesCleanFinishedMergedResidueAsFixable`: green throughout — it is the preserve-the-good-case guard and must never have gone red. Proven to pin "in `working/`" rather than "the REQ file exists" by widening the read to `return exists`, which reddens it alone.

**New tests added:** `TestVerifyDoesNotAdvertiseADirtyMergedWorktreeAsFixable`, `TestVerifyDoesNotAdvertiseAMergedWorktreeOfAnInFlightRequestAsFixable`, `TestVerifyStillAdvertisesCleanFinishedMergedResidueAsFixable`, plus the three remediation tests below and the helpers `mergeFixtureBranchIntoIntegration`, `assertNotAdvertisedAsMechanicallyFixable`, `assertRemovabilityProbeSkipped`.

**Existing tests updated (cross-REQ impact):** `TestVerifyClassifiesWorktreeLeftoversByMergeState` (from REQ-083) plants archived `REQ-003`/`REQ-004` files. It previously asserted `FixableCount() == 2` with no REQ file anywhere, which held only because the boolean read treated absent-from-board as finished. Its subject is the merge-state split; finishedness is the second axis this REQ adds, and the fixture was silent on it. Same categories, same count, now discriminating on merge state as intended. `TestVerifyDoesNotAdvertiseADirtyMergedWorktreeAsFixable` likewise plants an archived REQ-412 so dirt is the only unestablished fact it tests.

**CLI evidence** (built binary, five-leftover fixture, not unit tests): dirty+archived → `worktree-present-uncommitted-work`; clean+in-flight → `worktree-present-run-in-flight`; clean+archived → `merged-worktree-leftover [fixable]`; clean+absent-from-board → `worktree-present-state-unknown` plus a skip line naming the REQ; name with no id → `worktree-present-state-unknown` plus `carries no REQ-NNN id`. Footer reads `1 fixable`. On `873b513` the fourth and fifth both printed `merged-worktree-leftover [fixable]`.

*Verified by work action*

## Review

**Overall: 73%** | 2026-09-03T02:00:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 85% |
| Code Quality | 82% |
| Test Adequacy | 75% |
| Scope | 90% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 — a REQ id absent from `board.RequestsById` was read as "finished", so a live run's clean merged worktree still printed `merged-worktree-leftover [fixable]` with a remedy asserting its REQ had left `do-work/working/` (`verify.go:947-976`) — `impact-user-visible` → fixed in remediation (D-07)
- F2 — `actions/cleanup.md:140` and `docs/cleanup-guide.md:29,56` still treat merged-plus-clean as automatically removable, so Pass 5 would mechanically remove the in-flight worktree verify now protects — `impact-user-visible` → REQ-527 created
- F3 — `worktreeLeftoverStateUnknown` had no lock-in test: both entry paths could be replaced with `worktreeLeftoverFinishedResidue` and the suite stayed green — `impact-rule-change` → fixed in remediation (three tests plus a shared skipped-probe assertion)

**Minor findings:** 7 (report only) — M1 the new `git status` refreshed the builder worktree's git index, against `prime-do-kanban.md`'s "`verify` write[s] nothing at all"; M2 `serve.go`'s stale ~40ms probe-set measurement; M3 a worktree both dirty and in-flight loses its dirtiness; M4 the state-unknown remedy pointed "beside this finding" where skipped probes never render; M5 double id derivation; M6 the declared-but-untouched `forensics.md` drift; M7 version/changelog (integrator-owned, Step 9). **M1, M4 and M5 were fixed in remediation.** M2, M3 and M6 are recorded as accepted.

**Acceptance:** Partial at review time — both named regression cases closed end-to-end through the CLI, but a merged clean worktree whose REQ was invisible to the board still reported fixable. That gap is closed by the remediation and re-verified at the CLI.
**Suggested testing:** 3 items — a live board-serve contention check, a badge-width render check for the longer category strings, and a Pass 5 dry run against an in-flight builder (now covered by REQ-527's own Red-Green Proof).
**Follow-ups created:** REQ-527; **sweeps appended to:** None

*Reviewed by review-work action*

## Remediation

The review's verdict was Approve with follow-ups at 73%, which requires a follow-up REQ for every Important finding. Two of the three were fixed here instead, because both are this REQ's own unfinished work: F1 is the very defect the REQ exists to prevent, still reachable through a third input class, and F3 is a shipped mechanism with no test on a `tdd: true` REQ. Only F2, which lives in two files this REQ declared out of scope, went to REQ-527.

**Remediation commit:** see `commit:` below.

- **D-07** — `isRequestStillInFlight(board, id) bool` replaced by `classifyRequestPipelineState` returning `InFlight` / `Settled` / `AbsentFromBoard`, the same tri-state shape as `worktreeMergeState`. `AbsentFromBoard` routes to `worktreeLeftoverStateUnknown` with an error, so it also emits a skipped-probe line. **DECIDE & STATE.** This covers both reviewer repros, including the stray REQ file parked under `do-work/user-requests/`: the walk files that as `StrayRequestFiles`, so it never enters `RequestsById` and its `status: claimed` is unreadable either way.
- **D-08** — Kept the pipeline read ahead of the cleanliness probe rather than reordering, matching the existing precedence comment. **DECIDE & STATE.** Consequence recorded honestly: the dirty-worktree test had no REQ-412 file at all, so it began answering the earlier question; its fixture now plants an archived REQ-412 so dirt is the only unestablished fact it tests.
- **D-09** — `--no-optional-locks` lifted into a named constant used by both `worktreeHasUncommittedWork` and `worktreeDirtyQueueState`, with the "top-level option, must precede `-C`" trap in its doc comment. **DECIDE & STATE.** Fixing only the new probe would have knowingly shipped the same violation next door.
- **D-10** — REQ-083's fixture gains archived `REQ-003`/`REQ-004` files. **DECIDE & STATE.** Rationale in `## Testing` above.
- **D-11** — The three state-unknown tests share `assertRemovabilityProbeSkipped`, which requires exactly one matching skip line and checks the evidence fragment inside it. **DECIDE & STATE.** Asserting the category alone would leave the `SkippedProbes` half unpinned, which was half of F3's complaint.
- **D-12** — The absent-from-board test is table-driven over the reviewer's two repro shapes rather than two near-duplicate functions. **DECIDE & STATE.**

**Lock-in confirmed by neutering, one guard at a time** — each replaced with `return worktreeLeftoverFinishedResidue, nil`; in every case exactly one test went red:
- `requestId == ""` → `TestVerifyDoesNotAdvertiseAMergedWorktreeWithoutARequestIdAsFixable`: `got 0 worktree-present-state-unknown findings, want 1`
- `case requestPipelineStateAbsentFromBoard` (F1's new path) → `TestVerifyDoesNotAdvertiseAMergedWorktreeOfARequestTheBoardNeverSawAsFixable`, **both** subtests
- `statusError != nil` → `TestVerifyDoesNotAdvertiseAMergedWorktreeWithAnUnreadableStatusAsFixable`. The status path is reached without stubbing: the fixture merges, archives, then deletes the worktree directory, so `git worktree list` still reports it while `git status` in a vanished path fails.

**M1 measured before and after**, all worktree files `touch`ed first to invalidate git's stat cache: before, all five `.git/worktrees/*/index` files changed hash (e.g. `ca98cf3cd74a → f508052409fa`); after, zero of five.

## Lessons Learned

**What worked:**
- Reaching for the tri-state the file already used next door. `worktreeMergeState` had exactly the right shape — merged / unmerged / *git declined to answer* — and the fix for F1 was to stop answering "is this finished?" with a boolean and copy that shape. REQ-457 landed on the same answer independently in the same week, which is what made it worth a family slug rather than a one-off note.
- Running the built binary against a purpose-made fixture instead of trusting the unit tests. Both F1 repros and the M1 index write were found that way; `TestVerifyWritesNothing` passes and cannot catch M1, because its fixture is not a git repo so the probe never runs there.

**What didn't:**
- The first pass fixed the two cases the user reported and stopped. A third input class — a REQ id the board never saw — still reached `Fixable: true` while its remedy asserted the REQ had left `working/`, which nothing had established. Closing the reported instances is not closing the defect.
- Declaring `forensics.md` in Scope before reading it. It turned out to have nothing to update by design, so the declaration became drift that had to be explained rather than an edit that had to be made.
- Adding a probe without asking whether the probe itself mutates. `git status` refreshes the index it reads; five of five worktree indices changed hash before `--no-optional-locks` went in, against the tool's own written promise that verify writes nothing at all.

**Worth knowing:**
- `git status` is not read-only, and the flag that makes it so (`--no-optional-locks`) is a top-level option — it must precede `-C`, not follow the subcommand.
- A REQ file parked outside the scanned sections is recorded as `StrayRequestFiles` and never enters `board.RequestsById`, so its `status:` is unreadable through the board no matter what it says. Any board lookup keyed on presence has to treat absence as unknown, not as a negative.
- A fixture that plants no REQ file is not testing finishedness — it is riding on whatever the absent case happens to do. REQ-083's fixture asserted two fixable leftovers exactly that way, which is why it had to change here.

## Orientation

`do-work verify` no longer tells you that an active builder's worktree is safe to delete. A `worktree-agent-*` leftover is advertised as mechanically fixable only when its branch is merged, its worktree is clean, and its REQ has left `do-work/working/`; every other state, including every state the probes could not establish, reports as present and non-fixable. Lives in the queue-kanban board tool's verify probe set.

`prime_files`: `_dev/primes/prime-kanban-board.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` — both spot-checked, referenced paths all still exist. The second had two stale claims this change corrects, both amended in the same commit: its Traps now state that merged never proves a worktree disposable, and that `git status` writes unless told not to.

**[MAP CHANGED]** — not a new module, but the meaning of `verify`'s `Fixable` flag for worktree leftovers changed, and `cleanup` Pass 5 has not caught up. Until REQ-527 lands, verify protects an in-flight worktree that cleanup would still remove.
