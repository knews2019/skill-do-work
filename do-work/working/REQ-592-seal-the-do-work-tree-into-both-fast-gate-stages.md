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
write_set: [_dev/tests/fast-stages.json, skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go, skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go, _dev/tests/fast-stage-reuse-behavior.sh]
title: '[impact-critical] Review fix: seal the do-work tree into both fast gate stages'
claimed_at: 2026-09-05T22:59:38Z
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
