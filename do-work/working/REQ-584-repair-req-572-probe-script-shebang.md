---
id: REQ-584
title: 'Repository gate repair: add the missing shebang to the REQ-572 probe script'
status: claimed
route: A
created_at: 2026-09-05T09:50:18Z
user_request: UR-098
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-substantive
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: do-work-runs-probe-script-missing-shebang
depends_on: []
related: [REQ-507]
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-05T09:52:43Z
  basis:
    - Route A
    - 1 acceptance criteria
claimed_at: 2026-09-05T09:51:15Z
review_at: 2026-09-05T10:04:30Z
---

# Repository gate repair: add the missing shebang to the REQ-572 probe script

## What

Repair the repository-gate failure recorded below so dependency-gated requests can resume.

## Instances

- [x] repository gate: shellcheck:sc2148:do-work-runs-req-572-probe:missing-shebang affecting REQ-507 (found by REQ-507 / UR-098)

## Repository Gate Repair Intake

- **Parent:** REQ-507
- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** shellcheck:sc2148:do-work-runs-req-572-probe:missing-shebang
- **Repair dependency:** REQ-584
- **Diagnostic evidence:** "ShellCheck error SC2148 in do-work/runs/work-2026-09-04-232225/REQ-572-probe.sh line 1: Tips depend on target shell and yours is unknown. Add a shebang or a 'shell' directive."
- **Diagnostic evidence:** "Both direct canonical gate runs at 12d264c2 (detached worktree, clean tree) exited 1 on this one lint finding before any Go test ran. The probe file was committed by another session in 7ba3148a as REQ-572 run evidence and is outside REQ-507's implementation range 8e3dbf01..ad8bceb7, which touches no do-work/runs path."
- **Implementation base:** 8e3dbf01e0660424965d79acb2e386b6604e4780
- **Implementation merge:** ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9

## Triage

**Route: A** - Simple

**Reasoning:** The intake names one file and one lint finding (ShellCheck SC2148, a missing shebang on a tracked probe script under do-work/runs), so the change is a single known line. The defer-gate preset of Route C is overridden because nothing here needs a plan or exploration.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Repository Gate Repair No-Op

- **Expected diagnostic fingerprint:** shellcheck:sc2148:do-work-runs-req-572-probe:missing-shebang
- **Gate command:** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 0
- **Recorded green revision:** f9659f0f324ff5295610b464241186a38c2e16bb
- **Observed result:** green before implementation; repair already satisfied
- **Verified at:** 2026-09-05T10:00:29Z

## Implementation Summary

**Files changed:** None — verified repository-gate repair no-op.

**What was done:** Re-ran the repair's recorded canonical repository gate before source edits and confirmed it is already green; no implementation changes were necessary.

## Qualification

Passed — repository-gate repair no-op; durable gate evidence verified and project diff empty.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` (direct, unpiped, from the project root at `f9659f0f324ff5295610b464241186a38c2e16bb`, `GOMAXPROCS=2`)
**Result:** ✓ Maintainer verification passed (384 board tests, 758 CLI tests; slowest CLI file `internal/finalization/finalization_req499_test.go` 18.14s, under the 30s budget). Exit 0 on the first run, so the one-retry rule did not fire.

**Focused probe:** `shellcheck -S warning do-work/runs/work-2026-09-04-232225/REQ-572-probe.sh` — exit 0 (the exact lint that failed at intake), recorded through `advance` as the request-bound `run-blocked-check` record; the green-gate record was recorded with `recorded_revision` `f9659f0f324ff5295610b464241186a38c2e16bb`.

**Red-green validation:** none required — `validate-already-green-repair` returned `tdd_allowed: true` at `2026-09-05T10:01:03Z`; the intake fingerprint `shellcheck:sc2148:do-work-runs-req-572-probe:missing-shebang` matches and the project diff is empty.

**Gate history behind this no-op (detached diagnostic worktree, `GOMAXPROCS=2`):**
- `12d264c2`: exit 1 twice — ShellCheck SC2148 on the REQ-572 probe script (the intake failure).
- `c5ee4a8c`: exit 1 — the probe already carried a shebang (commit `4d47c821`, another session), but 23 archive links in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` failed the shipped-package reference contract; repaired by another session in `9f914188`.
- `9f914188`: exit 0 (3m04s wall).
- `f9659f0f`: exit 0, the recorded direct run above.

**Heavy verification plan:** none — an already-green repair changes no project path, so no heavy lane is selected.

*Verified by work action*

## Review

**Overall: 99%** | 2026-09-05T10:02:58Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

Already-green repository-gate repair no-op accepted on the no-diff path. `validate-already-green-repair` was invoked freshly by this review and returned `outcome: success` with `already_green_repair.review_allowed: true` and `tdd_allowed: true`, empty `reason_codes` and `offending_paths`. Typed fingerprint: `intake_fingerprint` and `expected_fingerprint` are both `shellcheck:sc2148:do-work-runs-req-572-probe:missing-shebang`. Gate-evidence match: `record_provenance: persisted_green_run`, `gate_exit_status: 0`, `recorded_revision`, `head_revision` and `target_revision` all `f9659f0f324ff5295610b464241186a38c2e16bb`, `state: exact_revision_match`, `match_basis: exact_revision`, `matches: true`. `canonical_completion_paths`: `do-work/CHECKPOINT.md`, `do-work/archive/REQ-584-repair-req-572-probe-script-shebang.md`, `do-work/calibration-log.tsv`, `do-work/working/REQ-584-repair-req-572-probe-script-shebang.md`. `staged_paths`: empty. `project_changed_paths`: empty. Independently verified: the tracked probe script at `HEAD` begins `#!/usr/bin/env bash` (added by another session in commit `4d47c821`), `shellcheck -S warning do-work/runs/work-2026-09-04-232225/REQ-572-probe.sh` exits 0, and `git diff HEAD --stat -- . ':!do-work'` is empty. The predicate was not reconstructed and the maintainer gate was not re-run. Restatement Sweep not applicable — nothing redefined, no implementation diff.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- The sweep's `## Instances` line is still unticked (`- [ ] repository gate: shellcheck:sc2148:… affecting REQ-507`) although the instance is resolved at `HEAD`, so the archived sweep record will read as carrying one open instance — impact-negligible → report only
- Frontmatter carries `effort_estimate: effort-substantive` while the same file records Route A and a 5-minute `p50_active_minutes`; the token was set at minting, not by this build — impact-negligible → report only

**Acceptance:** Pass — validator green, the intake lint now exits 0, the shebang is present at `HEAD`, and the project diff and staging area are both empty.
**Suggested testing:** 1 item
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Deferring the red gate through `defer-gate` instead of patching another session's committed artifact by hand; the already-green no-op path closed the repair with one direct gate run and no source change.
**What didn't:** The first pre-build gate for this repair (at `c5ee4a8c`) was red for a second, unrelated reason (stale archive links in `lessons-do-work-cli.md`) that another live session repaired minutes later; with two sessions committing to `main`, a gate result is only meaningful at the exact revision it ran, so the recorded direct run had to be repeated at the newest `HEAD`.
**Worth knowing:** The Triage template's leading `---` line lands inside the preceding section, and `validate-already-green-repair` compares the durable `## Repository Gate Repair Intake` section byte-for-byte against the working copy, so that separator produced `REPAIR-INTAKE-NOT-DURABLE` until it was removed. Tick the sweep's `## Instances` line once the instance is resolved; the reviewer flagged it left open.

## Orientation

Gate repair, no code: the tracked REQ-572 probe script already carries its shebang (another session's commit `4d47c821`), so the maintainer gate is green again and REQ-507 (handing the archive and commit tails to `finalize`) can resume behind this repair. Lives in the do-work-cli repair-validation subsystem. No prime went stale.
