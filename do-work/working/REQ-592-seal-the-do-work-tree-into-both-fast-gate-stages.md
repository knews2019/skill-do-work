---
id: REQ-592
status: claimed
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
---

# Review Fix: Seal the do-work Tree Into Both Fast Gate Stages

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
