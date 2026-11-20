---
id: REQ-756
title: 'Addendum: Detect committed and ignored outside-bundle changes'
status: claimed
route: A
heavy_verified_at: 2026-11-20T16:27:06Z
heavy_verified_revision: 3295931f360da8157dbeb799674c6d27dfa86519
review_at: 2026-11-20T16:19:58Z
kb_status: pending
commit: ef938a96d404282c2ce42e6600127c06ba29c76a
integration_at: 2026-11-20T16:18:02Z
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-11-20T16:14:38Z
  basis: ["Route A", "6-file write set", "6 acceptance criteria"]
created_at: 2026-11-20T01:30:31Z
user_request: UR-188
addendum_to: REQ-754
related: [REQ-755]
domain: backend
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-releases.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-11-20T16:14:14Z
---
# Addendum: Detect committed and ignored outside-bundle changes

## What
Repair subprocess second-opinion checks so an outside edit cannot disappear from evidence merely because the backend committed it or Git ignores its path. Existing git status checks report no outside_paths in both cases and accept the run. These are inherited gaps that the worktree implementation retained.

## Requirements
- Detect and name tracked outside-bundle edits even after the backend commits them inside the detached worktree.
- Detect and name new ignored paths outside the isolated bundle.
- Compare against evidence captured before probes/generation, not mutable post-backend HEAD or ignore rules alone.
- Apply the same condition after model probes and generation, while preserving allowed bundle mutations and unrelated real-checkout dirt behavior.
- Keep existing-output preservation, bounded cancellation/process cleanup, and refusal-page recovery.
- Cover unchanged ignored baseline inputs and both staged/unstaged ordinary outside edits as controls where relevant to the chosen implementation.

## Red-Green Proof
**RED prompt/case:** `TestSecondOpinionWorktreeOutsideSnapshot` in the toolboxcommands Go harness, starting with the preserved `TestREQ754ReviewOutsideChanges` cases for committed README edits and ignored scratch-output.txt creation.
**Why RED now:** Both cases return state=ran, judge=pass, banner=inserted and outside_paths=[] against 4e48c2ef and the pre-change source.
**GREEN when:** Both fail with the offending relative path, leave real outside files unchanged, retain permitted finished output, and remove the owned worktree. Positive normal/dirty-real-checkout cases still succeed.
**Validation:** User authorized both reproduced fixes and execution after the review and retrospective.

## Prior Implementation
REQ-754 implemented detached temporary worktrees for subprocess second opinions, current-bundle copying, exclusive import and cleanup after process reaping. It removed real-checkout outside-status checks. Commits 60186a38 and e7bd242b were merged as d1156bc6, released as 0.331.2. The helper is `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go`; callers and checks are in `report_second_opinion.go`.

## Constraints
Keep the current uncommitted bundle contents and unrelated real-checkout changes intact. Keep backend permission flags and session-mode behavior. Retain permitted finished refusal output in the real bundle. A worktree is not a security boundary. No new dependencies or paid live-model calls are required.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 25337 tokens, slugged partial, exceeds the 2000-token capture budget. Relevant families include canonical-worktree-cleanup and final-boundary-identity. Read touch-required sections during implementation.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read CLI, release and action primes and relevant crew/lessons; recorded the six-file focused approach.
- [x] **[APPLY]:** Replaced Git status with a fixed pre-probe outside filesystem snapshot; added regression/positive controls and updated action/guide semantics.
- [x] **[UNIFY]:** Read the six-file diff, formatted Go, checked whitespace, reproduced RED then GREEN, ran race suite and Go 1.24/Windows builds; independent review approved. Integrated b681e515 as ef938a96.

## Context
No pending or pending-answers request in any UR shares this root cause. This is explicit user promotion of a report-only review finding, not an automatic critical follow-up.

## Assets
Original probes and observed results are preserved under `do-work/user-requests/UR-188/assets/`.

## Triage

**Route: A** — Reproduced detection bug with known call sites and test seams; no open design question.

## Plan

Planning not required (Route A). Replace mutable Git-status evidence with a fixed filesystem baseline of the isolated checkout outside the selected bundle and root Git control entry, captured before any backend probe. Compare file digests, types, permissions and symlink targets without following links. Check after every probe and generation; retain existing refusal recovery and cleanup. Cover committed/ignored edits at both phases, ordinary staged/unstaged edits, unchanged ignored inputs, and rename deletion evidence. Expected write set: shared runner, new snapshot helper and tests, existing rename test, action and guide (six files). No new dependencies.

## Decisions

- D-01 (DECIDE & STATE): Baseline filesystem evidence also handles untracked/ignored files and mutable ignore rules without combining several Git commands. Ignore only the checkout's root .git control entry; temporary Git internals are expected to change when a backend commits. This remains a worktree behavior check, not a security boundary.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified) — capture outside baseline before any probe and reuse it after each probe/generation; remove Git-status evidence.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_outside.go` (new) — fingerprint outside bytes, types, modes and link targets without following links; cancellable reads and sorted path differences.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_outside_test.go` (new) — probe/generation regressions for commits/deletions, ignored files, changed ignore rules, staged/unstaged changes and unchanged inputs; real checkout preservation and refusal cleanup assertions.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified) — preserve the existing rename-source assertion using the fixed baseline interface.
- `skills/do-work-toolbox/actions/ai-report.md` (modified) — describe fixed-snapshot evidence for subprocess backends.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified) — align user guidance with committed/ignored change detection.

**What was done:** The isolated checkout's outside files are compared with evidence captured before probes, independent of mutable HEAD/index/ignore rules. The report bundle remains governed by its existing stricter preservation rules. Implementation b681e515 integrated as ef938a96; qualification range 435f7fa9..ef938a96. No new dependencies or live-model calls.

## Qualification

Request-bound qualification passed for 435f7fa9..ef938a96. The two new-file static-reference warnings are Go compilation/test discovery, confirmed by the public runner calling the new helper and the executed tests. All six acceptance criteria traced; six-file scope matches. The existing rename-source test was adapted to the new evidence interface while preserving its assertion (REQ-754 adjacent behavior).

## Testing

RED: TestSecondOpinionWorktreeOutsideSnapshot reproduced committed edits/deletions, ignored writes and changed ignore rules escaping both phases (10.469s). GREEN: the same test plus baseline/rename controls passed (13.224s). All second-opinion tests passed with race detection (56.971s). Go 1.24.0 build and Windows/amd64 cross-build passed. Independent targeted review tests passed (15.975s). The merged request-bound focused probe passed (11.058s; green-record run 8.907s). Fresh direct `bash _dev/tests/maintainer-verify.sh` exited 0 in 204s: contracts, vet, formatting, 427 board and 1014 CLI fast tests passed. Both typed advance test-gate records are satisfied. Heavy lanes below await the queue drain.

**Existing tests updated (cross-REQ impact):** The earlier second-opinion rename check now supplies a fixed before-snapshot, and still asserts that moving an outside tracked file into the bundle names the outside source. No assertion was weakened.

## Review

**Overall: 97.5%** | 2026-11-20T16:19:58Z

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
**Acceptance:** Pass — independent review of the six-file candidate exercised the public handler with deterministic backend/judge fixtures; snapshot, outside-change and cancellation/deadline tests passed in 15.975s. RED-before-GREEN evidence confirms closure of both captured detection gaps. Parent race tests cover every second-opinion test. Integration gate and packaging results are recorded separately below.
**Restatement sweep:** Outside-change evidence now means a fixed pre-probe filesystem snapshot. Active action and guide agree; historical changelog/lesson entries remain historical. The existing rename test retains its original behavioral assertion.
**Suggested testing:** Paid live-model/browser acceptance unassessed; no backend protocol or rendering logic changed. Snapshot cost grows with checkout bytes; file reads are streamed and cancellation-aware.
**Follow-ups created:** None.

*Reviewed independently by review-work action; parent owns integration and finalization.*

## Lessons Learned

**What worked:** Pair committed/ignored negative cases with unchanged ignored inputs and ordinary edits at both backend phases. The original regression probes remain preserved in UR-188.
**What didn't:** Post-backend Git status cannot prove preservation because the backend can change HEAD and ignore rules. Worktree isolation limits copy-back but does not supply evidence about what changed.
**Worth knowing:** Keep the baseline fixed before the first probe; compare bytes/types/modes/link targets and never traverse repository symlinks. The root .git entry is intentionally excluded. This is filesystem change detection, not a security sandbox. The earlier REQ-754 retrospective records the review gap; these local lessons elaborate it without adding duplicate prime instructions.

## Orientation

[MAP CHANGED] Subprocess second opinions in the do-work CLI now validate outside-bundle files against a fixed filesystem baseline before accepting outputs. Backend commits and ignore rules cannot erase the evidence. The CLI/release prime's referenced paths remain valid; no routing change is needed.

## Heavy Verification Plan

Base: `435f7fa97f054eab40725f8f0f1fa3cd200dcca9`; target: `ef938a96d404282c2ce42e6600127c06ba29c76a`.

- `do-work-cli-integrations`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`. Reasons: skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_outside.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_outside_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go matched subtree skills/do-work/tools/do-work-cli.
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. Reasons: skills/do-work-toolbox/actions/ai-report.md matched subtree skills; skills/do-work-toolbox/docs/ai-report-guide.md matched subtree skills; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go matched subtree skills; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_outside.go matched subtree skills; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_outside_test.go matched subtree skills; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go matched subtree skills.
- `updater`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`. Reasons: skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_outside.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_outside_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go matched subtree skills/do-work/tools/do-work-cli.
- `installer`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`. Reasons: skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_outside.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_outside_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go matched subtree skills/do-work/tools/do-work-cli.

## Heavy Verification Result

Target: `ef938a96d404282c2ce42e6600127c06ba29c76a`; execution revision: `3295931f360da8157dbeb799674c6d27dfa86519`. Recomputed plans matched their stored ranges; no pending request remained. All selected lanes passed, none skipped.

- `do-work-cli-integrations`: exit 0, 98s, executed (fingerprint_uncertain).
- `staged-skills`: exit 0, 40s, executed (fingerprint_uncertain).
- `updater`: exit 0, 95s, executed (fingerprint_uncertain).
- `installer`: exit 0, 63s, executed (fingerprint_uncertain).
