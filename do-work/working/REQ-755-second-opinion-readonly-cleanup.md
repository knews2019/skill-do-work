---
id: REQ-755
title: 'Addendum: Remove second-opinion worktrees containing read-only directories'
status: claimed
route: A
review_at: 2026-11-20T16:14:06Z
kb_status: pending
commit: 10bd6e71103de2da6282d402ef890da25f68da2e
integration_at: 2026-11-20T01:34:19Z
estimate:
  p50_active_minutes: 10
  confidence: medium
  calculated_at: 2026-11-20T01:31:37Z
  basis:
    - Route A
    - 2-file write set
    - 4 acceptance criteria
created_at: 2026-11-20T01:30:31Z
user_request: UR-188
addendum_to: REQ-754
related: [REQ-756]
domain: backend
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-releases.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-11-20T01:30:48Z
---
# Addendum: Remove second-opinion worktrees containing read-only directories

## What
Repair cleanup when a valid report includes a readable but non-writable asset directory. The current helper reproduces that mode in the isolated copy, then both Git removal and os.RemoveAll fail. The generated page may already be published when the run reports cleanup failure and leaves scratch files.

## Requirements
- Remove owned temporary checkout files and registrations when a report contains a 0555 directory with regular assets.
- Cover success, refusal/backend failure, and cancellation using readable non-writable directories.
- Change deletion permissions only within the owned scratch tree, without following links or modifying original report modes or contents.
- Keep honest error reporting when cleanup truly cannot complete.

## Red-Green Proof
**RED prompt/case:** `TestSecondOpinionWorktreeReadOnlyCleanup` in the toolboxcommands Go harness. Begin with the preserved `TestREQ754ReviewCleanup` probe, adding meaningful failure/cancellation cases.
**Why RED now:** The current probe publishes and links the page, then reports state=failed with Permission denied and leaves its temporary checkout directory. The same input passed before REQ-754.
**GREEN when:** The same cases have their expected run result, preserve original bytes/modes, and leave neither owned disk files nor Git worktree registrations. Existing second-opinion tests remain green.
**Validation:** User authorized both reproduced fixes and execution after the review and retrospective.

## Prior Implementation
REQ-754 implemented detached temporary worktrees for subprocess second opinions, current-bundle copying, exclusive import and cleanup after process reaping. It removed real-checkout outside-status checks. Commits 60186a38 and e7bd242b were merged as d1156bc6, released as 0.331.2. The helper is `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go`; callers and checks are in `report_second_opinion.go`.

## Constraints
Keep the current uncommitted bundle contents and unrelated real-checkout changes intact. Keep backend permission flags and session-mode behavior. Retain permitted finished refusal output in the real bundle. A worktree is not a security boundary. No new dependencies or paid live-model calls are required.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 25337 tokens, slugged partial, exceeds the 2000-token capture budget. Relevant families include canonical-worktree-cleanup and final-boundary-identity. Read touch-required sections during implementation.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reproduce read-only-directory cleanup in success, backend-failure and cancellation cases. Restore owner access through verified directory handles only inside owned scratch before teardown; never follow links. Read the listed release lessons and relevant CLI cleanup/identity/process lessons.
- [x] **[APPLY]:** Added scratch-only directory permission restoration and four cleanup regressions on the isolated REQ branch.
- [x] **[UNIFY]:** Read both changed files and git diff, ran gofmt, diff-check, focused GREEN and race tests. Git/worktree registration and original bytes/modes assertions pass.

## Context
No pending or pending-answers request in any UR shares this root cause. This is explicit user promotion of a report-only review finding, not an automatic critical follow-up.

## Assets
Original probes and observed results are preserved under `do-work/user-requests/UR-188/assets/`.

## Triage

**Route: A** — Focused bug fix with reproduced input and known helper.

## Plan

Planning not required (Route A). The failure and code path are known.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go` (modified) — restore owner directory access on verified open handles in owned scratch after process reaping, before deletion; do not follow links.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_cleanup_test.go` (new) — success, backend failure, cancellation and outside-symlink refusal with read-only assets, original-byte/mode and residue assertions.

**What was done:** Valid read-only report directories no longer prevent worktree cleanup. Existing genuine cleanup errors remain reported. Implementation 24877ad3 integrated in 10bd6e71; merge range f2414954..10bd6e71.

## Decisions

- D-01 (DECIDE & STATE). Preserve copied permissions during generation and validation. Add owner access only for teardown after processes stop. WalkDir skips links, and chmod uses a verified directory handle.
- D-02 (DECIDE & STATE). Implement the small review repair in the main agent using an owned worktree. Main-tree lifecycle state stays with the orchestrator role.

## Qualification

Request-bound advance qualification passed for f2414954..10bd6e71. The new-file static-reference warning is Go test discovery, confirmed by the executed test. Both files and four acceptance criteria were checked; no scope drift.

## Testing

RED: four cases of TestSecondOpinionWorktreeReadOnlyCleanup reproduced permission-denied cleanup/residue before the fix (2.534s). GREEN: the same test passed (2.007s). The adjacent worktree/session tests passed with race detection (17.590s), and the resumed focused probe passed (17.467s). Fresh direct `bash _dev/tests/maintainer-verify.sh` exited 0 in 169s: 428 board and 1012 CLI fast tests, vet, formatting and contracts passed. The interrupted gate had no final result and was not counted. Heavy lanes below await the queue drain.

## Review

**Overall: 97%** | 2026-11-20T16:14:06Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.
**Minor findings:** None.
**Acceptance:** Pass — implementation and integration exercised using the real public handler with fake backend/pass-judge seams, plus cancellation. All four requirements met. WalkDir visits directories before reading their children, does not follow links, and chmod is bound to the observed directory handle; genuine errors remain joined. Original modes/bytes and worktree removal are asserted. Parent reviewed f2414954..10bd6e71; the interrupted independent review's preliminary comments are not counted as approval.
**Restatement sweep:** nothing redefined; cleanup behavior retains the existing documented contract.
**Suggested testing:** Live model/render acceptance unassessed; this repair changes teardown, covered by deterministic backend fixtures. Packaging lanes await drain.
**Follow-ups created:** None.

*Reviewed by review-work action in the current session*

## Lessons Learned

**What worked:** Reproduce read-only directories with nonempty assets and verify both filesystem and Git registration cleanup, including failure/cancellation.
**What didn't:** A force flag cannot bypass directory permissions; the previous cleanup assumed it could.
**Worth knowing:** Preserve original modes during generation and verification, then restore owner access only on owned scratch directories after writers stop. This is already covered by the CLI prime's canonical-worktree-cleanup/final-boundary-identity lessons; no duplicate satellite entry needed.

## Orientation

Second-opinion teardown in the do-work CLI now handles readable read-only report assets without altering original permissions. CLI/release prime paths remain valid and need no map update.

## Heavy Verification Plan

Base: `f2414954795c63dcb208a0410505e6869dfa1780`; target: `10bd6e71103de2da6282d402ef890da25f68da2e`.

- `do-work-cli-integrations`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`. Reasons: skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_cleanup_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go matched subtree skills/do-work/tools/do-work-cli.
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. Reasons: skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_cleanup_test.go matched subtree skills; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go matched subtree skills.
- `updater`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`. Reasons: skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_cleanup_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go matched subtree skills/do-work/tools/do-work-cli.
- `installer`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`. Reasons: skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_cleanup_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go matched subtree skills/do-work/tools/do-work-cli.
