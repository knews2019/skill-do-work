---
id: REQ-761
title: 'Match board executable names exactly before restart'
status: claimed
dispatch_at: 2026-11-25T23:27:36Z
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-11-25T23:27:30Z
  basis:
  - trivial short-circuit
route: A
created_at: 2026-11-25T23:09:19Z
user_request: UR-191
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
impact: impact-user-visible
effort_estimate: effort-mechanical
maintenance: false
claimed_at: 2026-11-25T23:27:19Z
builder_handback_at: 2026-11-26T00:18:52Z
integration_at: 2026-11-26T00:19:23Z
review_at: 2026-11-26T00:22:42Z
commit: bbfea807459e39c6e4fa19852e2c58ba566637f7
---
# Match board executable names exactly before restart
## What
Use an exact executable basename match in `skills/do-work-board/tools/run-kanban.sh`. Only `queue-kanban` and the legacy `go-bin-queue-kanban` qualify for an automatic stop. Commit `0a3d12e7` documents the legacy cross-repository restart.

## Red-Green Proof
**RED prompt/case:** In `_dev/tests/contracts/queue-kanban.sh`, run the foreign-listener fixture with `queue-kanban-helper`. The current wildcard at `skills/do-work-board/tools/run-kanban.sh:21` accepts the lookalike. The validation replay returned exit 1 but recorded a signal to stub PID 3131.
**GREEN when:** Run `bash _dev/tests/contracts/queue-kanban.sh`. For each documented board name, assert that the fixture signals its listener once and starts the replacement successfully. For `queue-kanban-helper`, assert refusal (exit 1) and an empty signal log. Keep the existing foreign-listener case.
**Validation:** User confirmed.

## Context
Accepted finding F1 from [the proposal report](../../ai-reports/2026-11-26_0041_proposal-repository-fix-priorities/index.html#f1), rated high impact. The fix narrows the existing signal decision.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` (9,341 tokens) covers shipped shell and exported-function fixtures.
- `_dev/primes/lessons-kanban-board.md` (7,493 tokens) covers exact basename matching.
Both index entries have partial family coverage and exceed the 2,000-token budget as whole files.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: [Original request](../user-requests/UR-191/input.md)*

## Triage

Route A. The executable check, both supported basenames and the fixture assertions are specified.

## Plan

**Planning not required** — Route A, direct implementation.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/run-kanban.sh` (modified)
- `_dev/tests/contracts/queue-kanban.sh` (modified)

The launcher now compares the exact basename reported by `lsof` with the two documented board names. The shell contract checks both successful restarts and both refusals with stubbed process discovery, signal delivery, build, and serve.

Integration range `a6a8d147a3bfb21032a9f54cf7abece0d52e80e1..bbfea807459e39c6e4fa19852e2c58ba566637f7`. The merge's second parent is `7cfb25b00ce99f4e75a233e86c070b180d8601f8`. The handback predates takeover, so no builder-work duration was inferred.

## Qualification

Passed. Typed `advance` returned a satisfied request-bound `qualify` record with `merged_range` provenance for `a6a8d147a3bfb21032a9f54cf7abece0d52e80e1..bbfea807459e39c6e4fa19852e2c58ba566637f7`. The two changed paths match the implementation summary and `git diff --check` passed. The existing stop, wait, and serve path now receives only the two exact executable basenames.

## Testing

- `bash _dev/tests/contracts/queue-kanban.sh` passed directly in under one second and through typed `advance` with a 29-second timeout. Both exact names recorded one signal to stub PID 4242 and reached the replacement with the expected arguments. `queue-kanban-helper` and `python3` each exited 1 with zero signals. The stuck-listener and static-board checks also passed.
- `bash _dev/tests/maintainer-verify.sh` exited 0 on its first direct run in 138 seconds. Typed `advance` recorded a satisfied green gate at `bbfea807459e39c6e4fa19852e2c58ba566637f7`. Go vet, fast tests, shell lint, and aggregate contracts passed.
- `bash -n` passed for both scripts. ShellCheck passed for the launcher and for the fixture with `SC2329` excluded for the exported child-shell stubs. The quiet-grep pipeline audit passed.

**Red-green validation:** The expanded contract exited 1 against the pre-merge launcher because `queue-kanban-helper` recorded a signal to PID 3131. It exited 0 against the merged launcher. Both executions stubbed signal delivery.

**Existing tests updated:** The launcher contract from REQ-717 now covers both authorized names and the lookalike. Its unrelated-listener, bounded-wait, static output, and build-failure checks remain in place.

**Heavy verification plan:** Range `a6a8d147a3bfb21032a9f54cf7abece0d52e80e1..bbfea807459e39c6e4fa19852e2c58ba566637f7`. See the exact selected lanes and reasons below.

*Verified by work action*

## Decisions

- **D-01.** Reuse the success and refusal fixtures with one iteration per executable and isolated logs for each case.

## Discovered Tasks

None.

## Review

Approved for local integration. Compared the two-file merge with the original UR, request, builder commit, and independent review. The launcher compares the executable basename returned by `lsof`. The full `ps` command remains in the diagnostic. The root and installed recipes both invoke the launcher. Independent verification also exercised the installed recipe in a temporary project path containing spaces.

**Overall: 100%** | 2026-11-26T00:22:40Z

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.

**Minor findings:** None.

**Acceptance:** Pass for implementation and local integration. The fixtures use stubs for process discovery, signal delivery, and replacement launch.

**Restatement sweep:** Redefined automatic stop eligibility for the board launcher. Checked the root and installed recipes, board instructions, prime pointers, and historical changelog entries. Earlier wave member REQ-760 is archived and recorded no redefinition. No stale current restatements found.

**Suggested testing:** 0 items.

**Follow-ups created:** None.

*Reviewed by review-work action*

## Orientation

The board restart command leaves unrelated listeners alone and can restart either documented board executable.

## Heavy Verification Plan

- Base revision `a6a8d147a3bfb21032a9f54cf7abece0d52e80e1`
- Target revision `bbfea807459e39c6e4fa19852e2c58ba566637f7`
- `queue-kanban-javascript` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`. _dev/tests/contracts/queue-kanban.sh matched subtree _dev/tests.
- `queue-kanban-browser` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`. _dev/tests/contracts/queue-kanban.sh matched subtree _dev/tests.
- `do-work-cli-integrations` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`. _dev/tests/contracts/queue-kanban.sh matched subtree _dev/tests.
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. _dev/tests/contracts/queue-kanban.sh matched subtree _dev/tests; skills/do-work-board/tools/run-kanban.sh matched subtree skills.
- `updater` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`. _dev/tests/contracts/queue-kanban.sh matched subtree _dev/tests.
- `installer` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`. _dev/tests/contracts/queue-kanban.sh matched subtree _dev/tests.
