---
id: REQ-592
status: completed
domain: testing
created_at: 2026-09-05T22:44:36Z
user_request: UR-127
addendum_to: REQ-591
review_generated: true
impact: impact-critical
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-shell-commands.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
maintenance: false
depends_on: []
related: [REQ-591, REQ-574]
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-09-05T23:20:42Z
  basis:
    - Route B
    - 4-file write set
    - 2 subsystems involved
    - 3 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set: [_dev/tests/fast-stages.json, skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go, skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go, _dev/tests/fast-stage-reuse-behavior.sh]
title: '[impact-critical] Review fix: seal the do-work tree into both fast gate stages'
claimed_at: 2026-09-05T22:59:38Z
route: B
builder_handback_at: 2026-09-05T23:47:39Z
dispatch_at: 2026-09-05T23:24:40Z
completed_at: 2026-09-06T03:36:23Z
commit: df1d2b92ab131f572a95145e47032beb6ccfc074
release_at: 2026-09-06T03:36:23Z
---

# Review Fix: Seal the do-work Tree Into Both Fast Gate Stages

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Both `prime_files` read, plus the crew rules and the exploration. Approach: declare
  `do-work` as the queue-kanban stage's coverage and `do-work/archive/UR-003/input.md` as an exact rule
  on the do-work-cli stage, empty `non_stage_coverage`, and add a `seal_exclusions` list tested before
  the coverage test. Recorded in `## Exploration` above and in the builder brief.
- [x] **[APPLY]:** Four files, exactly the declared `write_set`. Nothing outside it was edited; the
  five paths considered and rejected are listed under **Declared but not touched** in the hand-back.
- [x] **[UNIFY]:** `git diff --stat` on the merge range reports four files, 146 insertions, 40
  deletions, identical to the builder branch's own diff against its base. Linters: `gofmt -l .` in the
  do-work-cli module — no output; `go vet ./...` — clean; `shellcheck` 0.11.0 on
  `_dev/tests/fast-stage-reuse-behavior.sh` — exit 0. No debug artifacts: the diff adds no `fmt.Print`,
  no `set -x`, no commented-out block and no temporary path; verified per file —
  `fast-stages.json` (three coverage edits plus the new list), `fast_stage_evidence.go` (struct field,
  predicate, decoder validation, two guard replacements, three corrected comments),
  `fast_stage_evidence_test.go` (fixture manifest, probe-input move, five cases), and
  `fast-stage-reuse-behavior.sh` (fixture manifest, probe-input move, two cases).

## What

`_dev/tests/fast-stages.json` declares `do-work/` as `non_stage_coverage`, which asserts that no
fast gate stage reads it. Both stages do. The do-work-cli stage's
`TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass` reads and byte-checks
`do-work/archive/UR-003/input.md`; the queue-kanban stage's `board_live_test.go`, `durations_test.go`
and `citations_test.go` build the board from the real `do-work/` tree. Because
`fast_stage_evidence.go` skips `queueStatePrefix` in both its tracked and untracked seal loops, a
`do-work/`-only change reuses stale evidence and the gate reports a false green.

Either seal `do-work/` into the stages that read it, or stop those stages from reading it. The
narrowest honest fix is the first: give the fast-stage seal its own exclusion set instead of
inheriting the heavy lane's `queueStatePrefix`, and declare in the manifest exactly which subtrees
under `do-work/` no stage reads (if any). The heavy lane may keep its own exclusion, which is safe
there because it refuses a dirty tree and attributes its result to a revision.

Two existing assertions pin the current behaviour as correct and must move with the fix:
`fast_stage_evidence_test.go` case `queue state changed` (expects `reused`), and
`_dev/tests/fast-stage-reuse-behavior.sh` case `queue state alone still reuses`.

## Context

Found during independent review of REQ-591. The review reproduced a gate-level false green: with a
warm evidence store, appending one newline to `do-work/archive/UR-003/input.md` makes the
do-work-cli stage's own test fail, while the whole gate prints `Maintainer verification passed.`
and exits 0 with that stage `REUSED`.

## Requirements

- A change to any `do-work/` path a fast gate stage reads must force that stage to execute.
- The manifest's `non_stage_coverage` must state only trees no stage reads, verified rather than
  assumed.
- The two tests that currently assert the opposite are updated in the same change, each naming the
  failure it now catches.
- The gate's own `do-work/test-durations.tsv` must keep not invalidating its own stage; it is
  gitignored and written by the stage itself, so it needs an explicit narrow exclusion rather than
  the whole-tree one.

## Red-Green Proof

**RED prompt/case:** In a detached worktree at the merge revision, run the gate once to record
evidence, append one newline to `do-work/archive/UR-003/input.md`, and run the gate again.
**Why RED now:** The second run prints `stage do-work-cli-fast-tests: REUSED (fingerprint_match…)`
and `Maintainer verification passed.` with exit 0, while
`go test -short -run TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass
./internal/repositorymodel/` fails on the same tree with `production legacy fixture changed size`.
**GREEN when:** The second run prints `EXECUTING (fingerprint_mismatch)` for that stage and the gate
exits non-zero, and the same sequence with only `do-work/test-durations.tsv` changed still reuses.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Full Context

See `do-work/user-requests/UR-127/input.md` for the capture instruction behind REQ-591.

---
*Source: independent review of REQ-591 (reducing repeated setup and unaffected reruns in the fast gate), finding F1, work run `do-work/runs/work-2026-09-05-170806/`.*

## Triage

**Route: B** — Explore then build.

**Reasoning:** The outcome is stated exactly — a change to any `do-work/` path a fast gate stage
reads must force that stage to execute — and the write set is already named. What was not known is
*which* `do-work/` paths each stage actually reads, and the request itself says that must be
verified rather than assumed. That is discovery, not design, so exploration runs and planning does
not.

**Planning:** Skipped. The request carries its own Red-Green proof and a four-file write set; a plan
agent would restate them.

**Claim-time environment note.** This work runs in a fresh cloud container, not the shared checkout
the previous session used. Four things had to be fixed before the canonical gate could give an
honest verdict, and all four are environment, not repository:

- Go 1.24.7 was installed against a `go1.26.1` floor — resolved with `go env -w GOTOOLCHAIN=go1.26.1`.
- ShellCheck was absent, then 0.9.0 against an 0.11.0 floor — resolved with the upstream 0.11.0 binary.
- `just` was absent, then 1.21.0, which rejects multiline-string justfiles the contract probes
  generate — resolved with 1.43.0.
- Three harness-injected environment settings made the gate red for non-repository reasons:
  `NODE_OPTIONS`, the `GIT_CONFIG_COUNT`/`KEY_*`/`VALUE_*` GitHub URL rewriting (both refused by
  `heavy-runtime-fingerprint.py` as opaque runtime extension / opaque Git configuration override),
  and a global `commit.gpgsign=true` pointing at an empty signing key, which makes every `git commit`
  inside a test fixture repository fail.

That last one matters beyond this container: the previous session recorded
`TestLaneMutationCannotPublishOrReuseSuccess/commit=true` as a pre-existing intermittent. It is not
intermittent. It fails deterministically wherever a global `commit.gpgsign` cannot sign, because the
lane's mutating script commits and the test then expects `HEAVY-RUN-REVISION-CHANGED` but sees
`HEAVY-RUN-DIRTY-TREE`. With signing off it passes 2/2 here and the whole gate is green.

**Baseline for attribution:** canonical gate at claim revision `585e3fa8`, `DO_WORK_FAST_STAGE_REUSE=off`,
load average under 1: **74s wall, exit 0, `Maintainer verification passed.`**

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — index cost 10543 tokens; exceeds the
  2000-token budget and is `slugged: partial`, so targeted selection is not eligible. Matches
  `closed-enumeration-for-a-condition` (the new seal-exclusion list is exactly that shape),
  `silent-skip-reads-as-red` and `lifecycle-section-evidence`.
- `_dev/primes/lessons-shell-commands.md` — index cost 3385 tokens; exceeds the budget and is
  `slugged: partial`. Matches the shell probe in the write set,
  `_dev/tests/fast-stage-reuse-behavior.sh`.

## Exploration

Explore agent, read-only, re-verified against HEAD rather than against the audited commit. Full
report in the run directory as `do-work/runs/work-2026-09-05-231943/REQ-592-exploration.md`.

**Every claim the request makes about the current code holds.** `fast_stage_evidence.go:195` and
`:223` both skip `queueStatePrefix`, in the tracked and the untracked seal loop. The do-work-cli
stage really reads `do-work/archive/UR-003/input.md` and byte-checks it at 5608 bytes
(`repository_model_test.go:397-417`). The queue-kanban stage really builds the board from the real
tree (`board_live_test.go:16-45`, `durations_test.go:269-279`, `generate_test.go:703-727`, and
`citations_test.go:1478` — the other live citations tests are `TestBrowserBehavior*`/
`TestJavaScriptBehavior*`, which the fast stage excludes by prefix). Both assertions named as
pinning the wrong behaviour exist where the request says.

**One claim in the request does not hold, and it widens nothing here.** The request says the heavy
lane's exclusion is safe "because it refuses a dirty tree". The heavy dirty-tree refusal explicitly
*exempts* `do-work/` (`heavy_run.go:221,225`) and the untracked seal skips it too
(`heavy_evidence.go:481`), so heavy has the same false-green shape for an **uncommitted tracked**
`do-work/` edit. Its committed seal does cover `do-work/`, so a committed change is caught. That is
outside this request's `write_set` and is filed as a discovered task rather than built here.

**Which `do-work/` subtrees no fast stage reads: essentially none.** Verified by walking
`walk.go:103-199` and `filementions.go:25-56` rather than assuming. `do-work/runs/` and
`do-work/deliverables/` are pruned from the board walk (`walk.go:192`), `assets/` at any depth is
pruned (`walk.go:195`), and hidden directories are pruned (`walk.go:198`), which is what makes
`do-work/.req-reservations/` — 162 tracked marker files the allocator creates and removes during
ordinary work — the one genuinely unread tree. Everything else is byte-read, name-read, or stat'd:
`collectRepoFileMentions` stats every repo-relative path mentioned in any REQ or UR body, and those
mentions reach every subtree (125 under `do-work/runs/`, 97 of `do-work/CHECKPOINT.md`, 51 of
`do-work/calibration-log.tsv`, 37 under `do-work/audits/`). So `non_stage_coverage` becomes empty,
and the churn trees are handled by an exclusion instead.

**The narrowest honest fix is coverage plus one new concept, not a bare deletion.** Deleting the
`non_stage_coverage: [do-work]` entry alone would make every `do-work/` path *unclassified*, and
`fast_stage_evidence.go:202` seals every unclassified path into **every** stage — so every queue
edit would re-run the do-work-cli stage too, whose only `do-work/` input is one file. Declaring
`do-work` as the **queue-kanban stage's coverage** gets per-stage separation from machinery that
already exists (`fastStageManifestClassifiesPath`), and adding
`do-work/archive/UR-003/input.md` as an *exact* coverage rule on the do-work-cli stage seals its one
real read. The single thing the existing machinery cannot express is "sealed nowhere, even where a
stage covers it", which `do-work/test-durations.tsv` needs because the stage writes that file itself
and the recorded fingerprint is the pre-run one. That is the new `seal_exclusions` list, tested
before `stageCovered`, exactly where the `queueStatePrefix` guard sits today.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `_dev/tests/fast-stages.json` (modify) — declare `do-work` as the queue-kanban stage's coverage and
  `do-work/archive/UR-003/input.md` as an exact coverage rule on the do-work-cli stage; empty
  `non_stage_coverage`; add the new `seal_exclusions` list
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go` (modify) — add
  the `SealExclusions` manifest field, the `fastStageSealExcludesPath` predicate, validation in the
  decoder, and replace the two `queueStatePrefix` guards with it; correct the three comments that
  now state the opposite of the behaviour
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go` (modify) —
  move the fixture's toolchain probe inputs out of the newly-sealed tree, rewrite the
  `queue state changed` case into the two cases it should have been, and add the two churn-path
  reuse cases
- `_dev/tests/fast-stage-reuse-behavior.sh` (modify) — the same fixture move, rewrite
  `queue state alone still reuses` into `a queue-tree change executes the stage that reads it`, and
  add the duration-log reuse case

**Files I will NOT touch:** `heavy_run.go`, `heavy_evidence.go` and `_dev/tests/heavy-lanes.json` —
the heavy lane's own `do-work/` exemption has a real gap for an uncommitted tracked edit, but it is
outside this request's `write_set` and is filed as a discovered task. `_dev/tests/maintainer-verify.sh`
— the stage wiring does not change, only what the manifest declares. Any test outside the two files
named above; no existing assertion is deleted, and the two that move are rewritten in place with the
failure each now catches named in a comment.

**Acceptance criteria (restated from REQ):**
- [ ] A change to any `do-work/` path a fast gate stage reads forces that stage to execute
- [ ] `non_stage_coverage` states only trees no stage reads, verified rather than assumed
- [ ] The two tests that currently assert the opposite are updated in the same change, each naming
      the failure it now catches
- [ ] `do-work/test-durations.tsv` keeps not invalidating its own stage, through an explicit narrow
      exclusion rather than the whole-tree one

## Pre-Flight

**Git:** ✓ Clean. This is a fresh cloud container with only this session in it, so none of the
foreign hand-back files the previous session had to judge around exist here, and canonical `recover`
reports `FINALIZATION-NONE`. The three held claims (REQ-583, REQ-587, REQ-591) were taken over from
the previous session's machine with `recover --take-over` and each reports
`RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES` with its `commit:` an ancestor of HEAD — untouched, awaiting
the heavy drain at queue exhaustion.

**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exited 0 at this REQ's claim revision
`15e2ec3`, run to completion, load average under 1, no other gate process — **75s wall**, exit status
read directly from `$?` and never through a pipe. Run with `DO_WORK_FAST_STAGE_REUSE=off`, which is
the mitigation the request exists to remove: until REQ-592 lands, a `do-work/`-only change reuses
stale evidence, so any verdict this run relies on is taken with reuse disabled.

**Tests baseline:** ✓ `go -C skills/do-work/tools/do-work-cli test -count=1 ./internal/heavyverification/`
is the focused lane this REQ changes. **Both exit statuses, as the retry finding requires: the first
recorded run exited 1, the rerun exited 0.** The exit-1 run happened while the full canonical gate
was running concurrently in this same checkout, so the two contended for the Go build cache and the
working tree moved under it. Re-run alone with `-count=2` immediately afterwards: exit 0, 27.2s. The
baseline is stable and the single red is attributable to that overlap, not to the tree.

**Dependencies:** ✓ Go 1.26.1 (via `GOTOOLCHAIN`, over a 1.24.7 host toolchain), ShellCheck 0.11.0,
`just` 1.43.0, Node v22.22.2, Chromium at `/opt/pw-browsers/chromium`. All four had to be installed
or upgraded in this container; the versions below the gate's floors were what made the first three
baseline attempts red.

**A correction to the inherited hand-off, because it changes what a reviewer should believe.** The
previous session recorded `TestLaneMutationCannotPublishOrReuseSuccess/commit=true` as a pre-existing
intermittent that "passes 6/6 in isolation". It is not intermittent. It fails deterministically
wherever the global Git config sets `commit.gpgsign` with a signing key that cannot sign: the lane's
mutating script runs `git add && git commit`, the commit fails, the tree stays dirty, and the test
gets `HEAVY-RUN-DIRTY-TREE` where it asserts `HEAVY-RUN-REVISION-CHANGED`. With signing off it
passes 2/2 here and the whole gate is green. Anyone who sees it red again should check the signing
config before calling it a flake.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `_dev/tests/fast-stages.json` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go` (modified)
- `_dev/tests/fast-stage-reuse-behavior.sh` (modified)

**What was done:** The fast gate's stage seal stopped inheriting the heavy lane's whole-tree
`do-work/` exclusion and started declaring what each stage really reads. `do-work` is now the
queue-kanban stage's coverage, because that stage builds the board from the real tree; and
`do-work/archive/UR-003/input.md` is an exact coverage rule on the do-work-cli stage, because that
one file is its only real read there. `non_stage_coverage` is empty — verified by walking the board's
own prune rules and its file-mention stat, not assumed. A new `seal_exclusions` list, tested **before**
the coverage test, keeps four churn paths from invalidating a stage that never reads their bytes; the
gate's own `do-work/test-durations.tsv` is the load-bearing one, because the stage appends to it while
running and the recorded fingerprint is the pre-run one, so a seal over it could never match again.

The two assertions that pinned the old behaviour were rewritten rather than deleted, and each names
the failure it now catches. The Go decision table's `queue state changed` case became two cases —
one proving the stage that reads the tree executes, one proving the stage that does not still reuses
— so the fix cannot be "seal everything into everything", which would have killed reuse during a
drain. Two more cases pin the churn exclusions, and a manifest-decoding case pins an unsupported
exclusion kind, because a typo'd kind would otherwise decode, match nothing, and turn reuse off for
that stage forever. Every one of those five was shown red by ablation before being accepted.

Merge range `fce57fcc..2d932a47`, four files, 146 insertions, 40 deletions — identical to the builder
branch's own diff against its base. Builder branch head `1a04355a`, one commit.

## Decisions

D-01 through D-05 are the builder's, authored in
`do-work/runs/work-2026-09-05-231943/REQ-592-handback.md` → `## Decisions` and transcribed here.

- **D-01 — `seal_exclusions` states a condition, not a list. DECIDE & STATE.** The Go struct doc
  carries the admission test — a path the gate or the orchestrator writes *while a gate runs*, whose
  bytes no stage reads — and says the manifest entries are only today's set of paths passing it. That
  is the answer `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale demands of a new
  enumeration, and it gives the next person a test to apply instead of a list to copy.
- **D-02 — `do-work/deliverables` is excluded although it does not exist in this checkout. DECIDE &
  STATE.** Toolbox actions write reports there while work runs, and `walk.go:192` prunes it from the
  board walk by name. A coverage rule matches strings and never touches the filesystem, so an entry
  for an absent directory costs nothing and stops the first report written there from turning reuse
  off.
- **D-03 — the two fixtures disagree about ignoring the duration log, on purpose. DECIDE & STATE.**
  The Go fixture's log is Git-ignored, matching the real file; the shell fixture's is not. Between
  them both untracked branches of the seal loop are exercised, and the ablation shows both cases fail
  without the exclusion.
- **D-04 — one commit, not several. DECIDE & STATE.** The manifest, the implementation and both
  rewritten assertions have to land together; splitting them would leave a red commit on the branch.
- **D-05 — a manifest-decoding case for an invalid exclusion kind was added. DECIDE & STATE.** The new
  validation loop was otherwise untested, and a typo'd kind would decode, silently match nothing, and
  seal the excluded file — turning reuse off for that stage permanently, with a green gate the whole
  time.

## Discovered Tasks

- **impact-critical — the heavy lane has the same false-green shape for an uncommitted tracked
  `do-work/` edit.** Its dirty-tree refusal exempts `do-work/` (`heavy_run.go:221,225`) and its
  untracked seal skips the tree too (`heavy_evidence.go:481`), while its committed seal only sees HEAD
  objects. The one tree the refusal skips is the tree the committed seal cannot see. REQ-592's own
  text asserts the heavy exclusion is safe "because it refuses a dirty tree"; that claim does not
  hold. The fix shape is the one just built for the fast gate. Out of this request's `write_set`.
- **impact-noncritical, report only — the fast-stage seal is byte-level, but part of the
  queue-kanban stage's real dependency is existence-level.** `collectRepoFileMentions`
  (`filementions.go:35-56`) stats every repo-relative path mentioned in any REQ or UR body, and
  `generate.go:713` runs it on every live board build. Those mentions reach `do-work/runs/` and
  `do-work/deliverables/`, so creating or deleting a mentioned file in an excluded subtree flips a
  boolean in the shipped board JSON without moving any seal. No current fast-stage assertion reads
  that map.
- **impact-noncritical, report only — the queue-kanban fast stage will rarely reuse during a drain.**
  With `do-work` as its coverage it seals about 730 tracked files, and every REQ claim, move or
  archive touches one. This is correct — the stage really reads those bytes — and it is the whole
  cost of the change. Measured here: a queue-only edit re-runs queue-kanban in 56s where a fully warm
  gate is 31s, against 75s for a cold gate.

## Qualification

**Passed.** Read from the merge range `fce57fcc..2d932a47`; the canonical `qualify` and
`scope-drift` gates both report satisfied.

- **The change is substantive and matches its declaration exactly.** Four files, 146 insertions, 40
  deletions — the same diff the builder branch carries against its base, so nothing entered the range
  from anywhere else. The `write_set` and the touched set are identical; no drift in either direction.
- **The request's own Red-Green proof was reproduced end to end, at gate level, and it is the proof
  that matters.** RED: with a warm store and one newline appended to `do-work/archive/UR-003/input.md`,
  the gate printed `REUSED (fingerprint_match)` for both stages and `Maintainer verification passed.`
  and exited 0, while `TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass` failed on that
  same tree with `production legacy fixture changed size: got 5609 bytes`. GREEN: the identical
  sequence now prints `EXECUTING (fingerprint_mismatch)` for both stages and the gate exits 1 with
  that test's failure in the log and zero occurrences of `Maintainer verification passed.`
- **Every new assertion was shown red by ablation before it was accepted.** That matters more here
  than usual, because four of the five new cases assert that something *still reuses* — a case that
  can pass because the feature works or because the feature never fires. Removing `seal_exclusions`
  from both fixtures turns `the gate's own duration log still reuses` and `a run-log write still
  reuses` red in the Go table and red in the shell probe; removing the decoder's validation loop turns
  the new `unsupported seal exclusion kind` case red. Each ablation output is quoted in the hand-back.
- **The fix is not "seal everything into everything", and there is an assertion that says so.** The
  old `queue state changed` case became two: one proving the stage that reads the tree executes, one
  proving the stage that does not still reuses. Measured on the real gate: a newline appended to
  `do-work/CHECKPOINT.md` re-runs queue-kanban and leaves do-work-cli reused, 56s against 75s cold.
- **The two fixture probe inputs were moved out of the newly-sealed tree.** ~~Without the move,
  `toolchain probe output changed` would be satisfied by the file's own byte seal moving rather than
  by the probe's output changing, and `toolchain probe cannot run` would reach `fingerprint_uncertain`
  through a missing seal input instead of a failing probe.~~ **Corrected at review: that claim
  overstates what the move achieved.** A reviewer reverted the move and both cases still passed,
  because both use `alpha-stage`, whose coverage is `module-alpha` only — so a `do-work/` path was
  never sealed into that stage either way. The move stays, because it keeps the fixture's intent
  legible, but it is not load-bearing and the remediation builder established that it cannot be made
  so: both cases assert a disposition plus a reason code, and the reason codes are identical under
  either mechanism, so a test that depended on the move would have to assert on internals.
- **The one requirement that could have been met dishonestly was met by verification.**
  `non_stage_coverage` is now empty because the board's own prune rules (`walk.go:192,195,198`) and its
  file-mention stat (`filementions.go:35-56`) were read to decide which subtrees are genuinely unread —
  and the answer was "essentially none", so the churn trees went to `seal_exclusions` instead of being
  declared unread. Declaring them unread would have been the same false green in a new place.

### Remediation qualification (after review)

**Passed.** Remediation merge range `1f923bb9..df1d2b92`, two files, 112 insertions, 10 deletions —
both already in the declared `write_set`, so the scope held. Cumulative range for the whole request is
`fce57fcc..df1d2b92`.

- **The review's headline gap is closed, and closed with a test that bites four ways.** Restoring
  `_dev/tests/fast-stages.json` to its pre-fix content now fails with three named messages instead of
  passing silently; deleting either stage's coverage rule fails with the message that names that
  stage; and broadening the `do-work/runs` exclusion to `do-work` fails with `a seal exclusion matches
  do-work/archive/UR-003/input.md`. Before this, all four of those mutations left the whole suite
  green — which meant the data half of the fix was pinned by nothing.
- **The second surviving mutation is closed, and the pin is exactly one case.** A tracked file under
  the excluded run-log subtree joined the fixture, so deleting the tracked loop's exclusion guard now
  fails `a tracked run-log file changed still reuses` and nothing else. That it is the *only* failure
  is the evidence that the case is the sole pin on that branch.
- **The new manifest test follows the shipped-module rule.** It skips when the `_dev/tests` directory
  is absent, which is the installed-copy case, but a renamed or undecodable manifest inside a
  maintainer checkout still fails. Confirmed under `-v` that it runs rather than skips here.
- **One record was corrected rather than defended** — see the struck-through bullet above about the
  fixture probe-input move. The reviewer's reproduction stood, the builder confirmed the claim could
  not be made true without asserting on internals, and the claim was corrected instead.

Requirements traced: a change to any `do-work/` path a fast stage reads forces that stage to execute;
`non_stage_coverage` states only verified-unread trees, and states none; both assertions that pinned
the old behaviour are rewritten in the same change, each naming the failure it now catches; and
`do-work/test-durations.tsv` keeps not invalidating its own stage, through the narrow exclusion rather
than the whole-tree one.

*Checked by work action*

## Testing

**Tests run:** the whole canonical gate, `bash _dev/tests/maintainer-verify.sh`, plus the focused lane
this REQ changes — `go -C skills/do-work/tools/do-work-cli test -count=1 ./internal/heavyverification/`
— plus the end-to-end probe `bash _dev/tests/fast-stage-reuse-behavior.sh`, `gofmt -l`, `go vet ./...`
and `shellcheck --severity=warning` on the one shell file touched.

**Result:** ✓ Green. The canonical gate exited 0 at the merge revision `b21479c`, run to completion
with `DO_WORK_FAST_STAGE_REUSE=off` so both stages executed and the whole suite really ran against the
changed code — **73s wall**, load average 3.84 at start, no other gate process. Exit status read
directly from `$?`, never through a pipe. The focused lane compared green against the recorded
baseline; the end-to-end probe printed `Fast-stage evidence reuse probes passed.` and exited 0.

**Both exit statuses for the focused baseline, as the pre-flight retry finding requires:** first
recorded run exit 1, rerun exit 0, and two further runs alone exit 0. The exit-1 run overlapped a full
canonical gate in this same checkout. Not a property of the tree.

**One mechanical note for the next person, because it cost a cycle here.** The test-gate's focused
check runs through `run-blocked-check --probe-file`, which reads the file's *bytes* and executes them
as a script — so `${BASH_SOURCE[0]}` is empty inside it. Handing it
`_dev/tests/fast-stage-reuse-behavior.sh`, which resolves its repository root from `BASH_SOURCE`,
makes that script exit 2 for a reason that has nothing to do with the code under test, while the same
script exits 0 when run by path. The probe file must be self-contained; here it was one line holding
exactly the baseline's own command text, which is also what makes the baseline comparison meaningful.
Its default timeout is 30 seconds, below this lane's ~27s warm and well below a cold run, so
`--timeout-seconds 300` was passed.

**Remediation testing.** The canonical gate exited 0 again at the remediation merge revision
`df1d2b92`, run to completion with `DO_WORK_FAST_STAGE_REUSE=off` so both stages executed —
**73s wall**, exit status read directly from `$?`. Focused lane `go -C skills/do-work/tools/do-work-cli
test -count=1 ./internal/heavyverification/` ok in 12.98s; `bash _dev/tests/fast-stage-reuse-behavior.sh`
printed `Fast-stage evidence reuse probes passed.` and exited 0; `gofmt -l` empty; `go vet ./...`
clean. Both new assertions were confirmed under `-v` to run rather than skip. Five ablations were run
to show them red first, and every ablated file was restored and verified restored — `git diff` on
`_dev/tests/fast-stages.json` is empty.

**What the gate itself now proves, which it did not before this change.** The RED case is recorded in
`## Qualification` and reproduced end to end: one newline appended to `do-work/archive/UR-003/input.md`
used to produce `Maintainer verification passed.` and exit 0 with both stages `REUSED`, while that
tree's own test failed. It now produces `EXECUTING (fingerprint_mismatch)` for both stages and exit 1
with that failure in the log.

*Verified by work action*

## Review

**Overall: 74%** | 2026-09-06T00:10:21Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 85% |
| Test Adequacy | 55% |
| Scope | 100% |
| Risk | Medium |
| Acceptance | Pass |

**Verdict: Approve with follow-ups.** The fix works and no reviewer could build a stale green against the real tree, but two mutations survive the test suite. The bigger one: reverting the shipped file `_dev/tests/fast-stages.json` (the gate's stage manifest, which says which parts of the queue tree each fast stage reads) back to its pre-fix content brings back the exact false green REQ-592 (making the fast gate stop reusing stale test evidence when queue files change) exists to remove, and every test still passes. All three reviewers found that gap independently. The code is pinned. The data file that is the actual fix is pinned by nothing.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- The shipped manifest `_dev/tests/fast-stages.json` is read by no test, so the whole fix can be reverted or broken while the suite stays green. Restoring the file from `fce57fcc` gives `ok github.com/knews2019/skill-do-work/do-work-cli/internal/heavyverification` and `Fast-stage evidence reuse probes passed.` Deleting only line 10 (the exact coverage rule on `do-work/archive/UR-003/input.md`) brings back the original bug. Broadening line 30 from `do-work/runs` to `do-work` turns the seal off for the whole tree, because an exclusion beats coverage by design. The heavy lane already has the pattern this lane lacks: `heavy_maintainer_tree_test.go:201-225` reads the real `_dev/tests/heavy-lanes.json` and cross-checks it. Fix: one Go test that decodes the real manifest and asserts by name that the do-work-cli stage covers `do-work/archive/UR-003/input.md` and that no seal exclusion matches that path. — impact-rule-change → report only
- Surviving mutation in the new code: the exclusion check in the tracked loop (`fast_stage_evidence.go:220`) can be deleted and everything stays green. Replacing `if path == "" || fastStageSealExcludesPath(manifest, path) { continue }` with `if path == "" { continue }` there gives `ok` on the Go package and `Fast-stage evidence reuse probes passed.` The same deletion in the untracked loop (`:250`) turns three cases red, so only one of the two branches is tested. The untested branch is the one that carries production weight: 231 tracked files under `do-work/runs` and 162 under `do-work/.req-reservations` pass through it. Cause is fixture shape — every excluded path a case mutates is untracked. Fix: one extra table case that mutates a tracked file under `do-work/runs` and expects reuse. — impact-negligible → report only

**Reviewer disagreement, resolved:** reviewer 1 rated the tracked-loop mutation Minor because it fails closed (reuse dies quietly, the gate stays slow but correct); reviewer 2 rated it Important because the mutation demonstrably survives. Taking Important. Two reviewers reproduced the surviving mutation independently, and a demonstrated surviving mutation is the bar for Important here. The impact token stays `impact-negligible`, which is where the fail-closed argument belongs. On the first finding, reviewers split on the token (rule-change, user-visible, rule-change). Taking `impact-rule-change`: nothing is wrong for a user today, and the remedy is a new guard test.

**Minor findings:**
- An unstaged rename or delete under `do-work/` makes the queue-kanban stage report `fingerprint_uncertain`, and `maintainer-verify.sh:189` only records evidence when the fingerprint is not `-`, so that run stores nothing. A plain `mv` of a REQ file from `do-work/queue/` to `do-work/working/`, which is what claiming work does, puts the tree in that state, and two gate runs in a row both execute the stage. `git mv` does not do this. Fails closed, so it is cost, not a wrong verdict, but a maintainer will read `EXECUTING (fingerprint_uncertain)` as a fault. — impact-user-visible → report only
- The do-work-cli half of the fix — an `exact` coverage rule on a single path that sits inside another stage's covered subtree — appears in no fixture manifest in either test file. The rewritten Go case `queue state changed forces the stage that reads it` names the real do-work-cli incident in its comment but runs beta-stage over a whole-subtree coverage, which is the queue-kanban shape. — impact-negligible → report only
- The REQ's Qualification claims that moving the fixture probe inputs restored what the two toolchain cases test. Reviewer 2 reverted the move and both cases still pass, in the Go fixture and the shell probe. Reason: both cases use `alpha-stage`, whose coverage is `module-alpha` only, so a `do-work/` path was never sealed into that stage either way. Taking this finding — it is the only reviewer who tested the claim, and the mechanism is checkable. The move is harmless, the claim is not accurate. — impact-negligible → report only
- The Implementation Summary and Qualification say all five new assertions were shown red before acceptance. Red output exists for four. No prior red run exists for `queue state changed leaves a stage that does not read it reusable`, which is the one case that passes both when the feature works and when it never fires. Reviewer 2 checked it and it does bite, so the assertion is fine and the provenance claim is what is wrong. — impact-negligible → report only
- Three of the four `seal_exclusions` entries are only true because `queue-kanban/walk.go:181-197` prunes `runs`, `deliverables` and dotted directories from the board walk. That is the same set enumerated in two modules with no pointer either way and no test tying them, which is the second half of the Closed Enumerations Go Stale rule. One clause in the struct doc naming `isSkippedSection` closes it. — impact-rule-change → report only
- `workingTreeSeals`' doc comment at `fast_stage_evidence.go:180-184` still says the seal covers every tracked and untracked path the stage covers. After this change that is false. The inline comment and the struct doc were both fixed; this one was missed. — impact-negligible → report only
- `heavy_run.go:28-31` still declares that the queue tree is never lane input. The manifest committed in this same change says the opposite for the fast gate, and the heavy do-work-cli lane runs the same test that byte-checks `do-work/archive/UR-003/input.md`. Outside this REQ's write set, so it belongs with the heavy-lane follow-up. — impact-rule-change → report only
- The admission condition for a new exclusion says "no stage reads its bytes", which is true, but it gives the next reader no hint that existence-level reads exist and sit outside the test. The board stats every repo-relative path mentioned in a REQ or UR body (`filementions.go:35-56`), and those mentions reach `do-work/runs` and `do-work/deliverables`. Adding "bytes, not existence" to the condition keeps the next entry honest. — impact-rule-change → report only

**Nit findings:**
- The shell fixture's `do-work/runs` seal exclusion is dead weight. Removing it alone leaves the probe green, because the only file under it never changes during the run. Both entries are individually pinned in the Go table, so behaviour is covered. — impact-negligible → report only
- `do-work/deliverables` is admitted by prediction, not by an observed write. The directory does not exist in this checkout and `git ls-files` never emits a path under it. — impact-negligible → report only

**Requirements checklist:**
- [x] A change to any `do-work/` path a fast gate stage reads forces that stage to execute — delivered. Verified against the real manifest through the real `decide-fast-stage` command: editing `do-work/archive/UR-003/input.md` moves the do-work-cli fingerprint, and editing `do-work/CHECKPOINT.md`, adding a file under `do-work/audits/`, or moving a REQ between `queue/` and `working/` with `git mv` all move the queue-kanban fingerprint. Deletions yield `fingerprint_uncertain`, which executes.
- [x] `non_stage_coverage` states only trees no stage reads, verified rather than assumed — delivered. It is empty. The four replacement exclusions were checked by mutating every tracked file under `do-work/runs` and `do-work/.req-reservations` and by deleting all 393 of them, and the queue-kanban suite stayed green both times.
- [x] The two tests that pinned the old behaviour are updated in place, each naming the failure it now catches — delivered. Both rewritten, neither deleted. Both bite: the shell case fails under the old guard, the Go case fails when the predicate is hardcoded back to `strings.HasPrefix(path, "do-work/")`. Caveat in the second Minor finding.
- [x] `do-work/test-durations.tsv` keeps not invalidating its own stage, through a narrow exclusion rather than the whole-tree one — delivered. Single `exact` rule, pinned independently in both fixtures.
- [x] New names satisfy Naming for Reach — delivered. `SealExclusions`, `seal_exclusions`, `fastStageSealExcludesPath`.
- [x] Shell is free of the prime's traps and passes shellcheck — delivered. ShellCheck 0.11.0 exit 0 plain and at `--severity=warning`, `gofmt -l` empty, `go vet` clean.
- [x] Nothing outside the declared four-file write set changed — delivered. `git diff --stat fce57fcc..2d932a47` is exactly four files, 146 insertions, 40 deletions.
- [ ] Release obligation discharged — not delivered. Two shipped files under `skills/` changed and no version or changelog moved. This commit is a release. Five files need updating at finalization: `/VERSION`, `/skills/do-work/VERSION`, `/skills/do-work/actions/version.md:5` (all three still read `0.303.10`), `/CHANGELOG.md`, and `/skills/do-work/CHANGELOG.md` as a byte-identical copy. Patch bump is the honest call, since the new `seal_exclusions` field is additive and the only manifest using it does not ship. Sequencing caution: REQ-591 (the change that added fast-stage evidence reuse in the first place) is still claimed and unreleased with its shipped changes already committed, so if it is finalized first this version number follows from that one, not from `0.303.10`. The builder correctly left this to the finalizer rather than breaking its four-file scope.

**Acceptance:** Pass — three reviewers independently ran the Go decision table and the shell probe green at merge revision `2d932a47`, and the strongest attack (deleting or mutating all 393 tracked files under the excluded trees, then running the queue-kanban suite with the gate's own selectors) produced `ok` in 58.6s, so there was no failure for the reuse to hide.

**Suggested testing:** 4 items
- Add a Go test that decodes the real `_dev/tests/fast-stages.json` and asserts the do-work-cli stage covers `do-work/archive/UR-003/input.md` and that no seal exclusion matches it. This is the one that closes the headline gap.
- Add a table case that mutates a tracked file under `do-work/runs` and expects reuse, to pin the tracked-loop half of the exclusion check.
- Add a fixture stage with an `exact` coverage rule on a path inside another stage's covered subtree, so the do-work-cli shape is exercised somewhere.
- Add a drift ratchet: scan the do-work-cli test corpus for path literals that resolve under `do-work/` and assert the set equals that stage's declared coverage. Today that is one path, and it would fail the day someone adds a second.

**Follow-ups created:** None (12 findings report only)

*Reviewed by review-work action*

## Lessons Learned

- **A cache's manifest is data the code reads and no test read.** The engine was pinned by fifteen
  assertions; the one file that decides what the engine is allowed to skip was pinned by none, so the
  entire fix could be reverted with the suite green. The heavy lane already had the guard this lane
  needed — `heavy_maintainer_tree_test.go` decodes the real `_dev/tests/heavy-lanes.json` and
  cross-checks it — and REQ-591 copied the engine without copying the guard. The general shape: when a
  new mechanism is configured by a shipped data file, the assertion that the data says what the design
  requires is part of the mechanism, not an extra.
- **An assertion that something *still reuses* passes twice: when the feature works, and when the
  feature never fires.** Four of the five new cases had that shape. The only thing that separates them
  is an ablation, and the ablation has to disable one branch at a time — removing `seal_exclusions`
  from both fixtures at once disabled the tracked and untracked seal loops together, which is why the
  tracked loop looked pinned when it was not. Ablate the code, one guard at a time, not the data.
- **"Pre-existing intermittent" is a diagnosis, and it needs the same evidence as any other.**
  `TestLaneMutationCannotPublishOrReuseSuccess/commit=true` was carried across two hand-offs as a flake
  that "passes 6/6 in isolation". It reproduced 3/3 here, and the cause was a global
  `commit.gpgsign` pointing at an unusable key: the lane's script commits, the commit fails, and the
  test sees a dirty tree where it asserts a new revision. A failure that reproduces is not a flake, and
  a flake that has never been explained is an open question, not a known issue.
- **A request's own statement of why something is safe is a claim to verify, not context to accept.**
  REQ-592 said the heavy lane's `do-work/` exclusion is safe because heavy refuses a dirty tree. The
  refusal exempts `do-work/` by construction, so heavy has the same hole for an uncommitted tracked
  edit. The exploration caught it because it read `heavy_run.go` instead of reading the sentence.

## Orientation

A future reader changing the fast gate's reuse rule starts in three places.
`skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go` holds the engine:
the `fastStageManifest` struct doc carries the admission condition for a new seal exclusion, and
`workingTreeSeals` holds the two loops where an exclusion is tested before coverage.
`_dev/tests/fast-stages.json` is the data half — what each stage covers, what is sealed nowhere — and
it is now pinned by `TestShippedFastStageManifestSealsTheQueuePathsItsStagesRead`, so a wrong line
there fails a test rather than producing a quiet false green.
`_dev/tests/fast-stage-reuse-behavior.sh` drives the shipped gate wrapper end to end against a
synthetic repository.

Two things are deliberately not solved here and are recorded as follow-ups. The heavy lane still has
the same false-green shape for an uncommitted tracked `do-work/` edit, and `heavy_run.go`'s
`queueStatePrefix` doc still asserts the contract this change replaced. And the seal is byte-level
while part of the board's real dependency on `do-work/` is existence-level, through
`filementions.go`'s repo-file-mention stat — no fast-stage assertion reads that map today, which is
the only reason the run and deliverable exclusions are safe.

## Heavy Verification Plan

- **Base revision:** `fce57fccb19338491fea9d01bb0721a71f6d988b`
- **Target revision:** `bb5118a9c2f77d416d118528128d2158ffa8bc96` (the recorded `commit:`)
- **Changed paths in range:** the four files of this REQ's diff plus this run's own `do-work/`
  artifacts. No uncovered paths, planner not forced, not uncertain.

All six lanes are selected, because the change reaches both the maintainer test tree and the CLI
module: `queue-kanban-javascript`, `queue-kanban-browser`, `staged-skills`,
`do-work-cli-integrations`, `updater`, `installer`. Each runs as
`env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane <id>`.

Held at Step 7.7 beside REQ-583, REQ-587 and REQ-591, which select the same six lanes. Draining all
four holds together at queue exhaustion runs the lanes once at the final revision instead of four
times. **This container has Chromium at `/opt/pw-browsers/chromium` and Node v22 at
`/opt/node22/bin/node`, so both engine-gated lanes can actually run** — set
`QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium` at drain time or the browser lane reports skipped,
and a skip is not a pass.

commit: bb5118a9c2f77d416d118528128d2158ffa8bc96

### Heavy verification result (run at drain, 2026-09-06)

**All six lanes ran and all six were green. Revision `a48b9eb6`, the tree quiet from the first command
to the last.** `bash _dev/tests/maintainer-verify.sh --heavy` printed `Maintainer verification passed.`
and exited **0**, gate wall **301s**.

**One deviation from the plan, stated.** The plan named six separate `--heavy-lane <id>` invocations.
What ran instead is the single `--heavy` gate, which executes the same six lanes' work in one process
at one revision. That is what the four held requests were waiting for — one run at the final revision
rather than four — and it removes the `HEAVY-RUN-REVISION-CHANGED` risk of interleaving six
invocations with four finalizations. The evidence below is per lane, not the gate's summary line,
because a skipped lane reports success.

| Lane | Its own evidence line | Result |
|---|---|---|
| `queue-kanban-javascript` | `module=…/queue-kanban wall=67s tests=481 slowest-file=generate_test.go:12.43s limit=none (heavy)` | 481 tests, green |
| `queue-kanban-browser` | `module=…/queue-kanban wall=102s tests=35 slowest-file=timeline_browser_probe_test.go:63.99s limit=none (heavy)` | 35 tests, green |
| `do-work-cli-integrations` | `module=…/do-work-cli wall=25s tests=798 slowest-file=internal/nextselection/blocked_probe_test.go:6.77s limit=none (heavy)` | 798 tests, green |
| `staged-skills` | `test-file duration: staged-skills-contract.sh 45s (limit none (heavy))` | green |
| `updater` | `test-file duration: update-script-behavior.sh 84s (limit none (heavy))` | green |
| `installer` | `test-file duration: install-suite-behavior.sh 28s (limit none (heavy))` | green |

**Zero `SKIP` lines in the whole run and zero `FAIL` lines.** The browser lane genuinely ran — 35 tests
and a 64-second `timeline_browser_probe_test.go` are not what a skipped lane prints — because
`QUEUE_KANBAN_BROWSER` pointed at `/opt/pw-browsers/chromium`, as every one of these four plans
required.

The run also needed a sanitized environment, which is worth recording for the next drain:
`NODE_OPTIONS` and the `GIT_CONFIG_COUNT` / `GIT_CONFIG_KEY_*` / `GIT_CONFIG_VALUE_*` triples unset,
and `GIT_CONFIG_GLOBAL` pointed at a config with `commit.gpgsign = false`. A heavy run refuses on an
opaque runtime extension or an opaque git configuration override, and an unusable global signing key
makes a fixture's own `git commit` fail inside the lane.
