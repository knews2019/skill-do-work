---
id: REQ-542
title: 'SessionStart launches the background gate runner'
status: claimed
priority: now
created_at: 2026-09-03T14:49:02Z
user_request: UR-104
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-537, REQ-538, REQ-539, REQ-540, REQ-541, REQ-542]
batch: two-tier-gate
route: A
estimate:
  p50_active_minutes: 5
  rationale: "Update session-start.sh to launch gate-runner.sh with fail-soft pid guarding, and update gate-runner.sh to manage its pid file"
write_set:
  - skills/do-work/hooks/session-start.sh
  - _dev/tests/gate-runner.sh
  - skills/do-work/actions/version.md
  - CHANGELOG.md
  - skills/do-work/CHANGELOG.md
status_changed_at: 2026-09-04T13:08:46Z
claimed_at: 2026-09-04T13:13:09Z
---

# SessionStart Launches the Background Gate Runner

## What

The SessionStart hook launches `_dev/tests/gate-runner.sh` in the background when it is not already running for this checkout, so every session has the gate as evidence attached to revisions instead of a step. The runner never passes `--heavy`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** In `skills/do-work/hooks/session-start.sh`, check if running in a maintainer repository with `_dev/tests/gate-runner.sh`. If present, check PID file under `$TMPDIR/do-work-gate-runs/runner-<repo-id>.pid` to see if a runner process is already active. If not running, print log directory path, launch runner in background with nohup, and write PID file. In `_dev/tests/gate-runner.sh`, manage PID lifecycle with trap cleanup on exit. Add version bump and changelog entry.
- [x] **[APPLY]:** Code implemented in `skills/do-work/hooks/session-start.sh` and `_dev/tests/gate-runner.sh`. Bumped version to `0.274.1` in `version.md` and recorded in `CHANGELOG.md` and mirrored in `skills/do-work/CHANGELOG.md`.
- [x] **[UNIFY]:** Ran `session-start-hook-behavior.sh` and verified unchanged output for consumers. Verified runner launches in background and prevents second runner on subsequent session-start. Ran `maintainer-verify.sh`.

## Triage
- **Route:** A (fail-soft shell hook script update, mechanical integration)
- **Primary files:** `skills/do-work/hooks/session-start.sh`, `_dev/tests/gate-runner.sh`

## Plan
1. In `_dev/tests/gate-runner.sh`:
   - Compute `pid_file="$log_root/runner-$repo_id.pid"` where `repo_id` hashes `$repo_root`.
   - When not running with `--once`, verify if an existing runner is active via `kill -0 "$existing_pid"`. If active, exit 0.
   - Record `$$` into `$pid_file` and set `trap 'rm -f "$pid_file"' EXIT INT TERM`.
2. In `skills/do-work/hooks/session-start.sh`:
   - Determine repo root via `git -C "${CLAUDE_PROJECT_DIR:-.}" rev-parse --show-toplevel 2>/dev/null || echo "${CLAUDE_PROJECT_DIR:-.}"`.
   - If `"$REPO_ROOT/_dev/tests/gate-runner.sh"` exists:
     - Check `$LOG_ROOT/runner-$REPO_ID.pid`.
     - If not running (`! kill -0 "$RUNNER_PID"`):
       - Print `gate runner logging to $LOG_ROOT`.
       - Launch `nohup bash "$REPO_ROOT/_dev/tests/gate-runner.sh" >/dev/null 2>&1 &` and record `$!` to `$PID_FILE`.
   - Proceed with standard `exec bash "$CLI_LAUNCHER" ...`.
3. In `skills/do-work/actions/version.md` and `CHANGELOG.md` (and mirror `skills/do-work/CHANGELOG.md`):
   - Bump version to `0.274.1` and record release notes.
4. Verify:
   - Prove with `bash _dev/tests/gate-runner.sh --once`.
   - Verify `session-start-hook-behavior.sh`.
   - Verify second run starts no second runner.

## Scope
- `skills/do-work/hooks/session-start.sh`
- `_dev/tests/gate-runner.sh`
- `skills/do-work/actions/version.md`
- `CHANGELOG.md`
- `skills/do-work/CHANGELOG.md`

## Implementation Summary
- Updated `skills/do-work/hooks/session-start.sh` to check for `_dev/tests/gate-runner.sh`. When present and not already active, prints `gate runner logging to $LOG_ROOT`, launches the runner in the background using `nohup`, and writes the PID file.
- Updated `_dev/tests/gate-runner.sh` to track and clean up its own repository-specific PID file on exit (`trap ... EXIT INT TERM`), exiting safely if another instance is already active.
- Bumped version to `0.274.1` in `version.md` and documented the change in `CHANGELOG.md` and `skills/do-work/CHANGELOG.md`.

## Testing
- `bash _dev/tests/session-start-hook-behavior.sh`: passed (consumer checkouts retain unchanged banner).
- `bash skills/do-work/hooks/session-start.sh`: verified first invocation starts background runner and prints log path.
- Subsequent invocation of `session-start.sh`: verified no duplicate runner is launched and no extra line is printed.
- `bash _dev/tests/maintainer-verify.sh`: passed.

## Review
- **Independent Verification:** Completed against all requirements.
  - Requirement 1 (detect existing runner via pid file, avoid duplicate): Verified.
  - Requirement 2 (only start in maintainer repo, never in consumer): Verified.
  - Requirement 3 (print one line naming log directory when starting): Verified.
- **Defects:** 0 critical, 0 non-critical.
- **Score:** 100%.

## Lessons Learned
- When background runners are launched by hooks, associating PID files with unique repository hashes ensures concurrent checkouts don't interfere with each other.

## Context

- The runner exists since 0.266.9 and records green through `record-green-gate`; a pipeline claim that finds HEAD proven green skips its baseline. Nothing starts it today.
- The approved draft left the choice open: SessionStart hook (Recommended) or a `just gate-watch` recipe. Capture took the recommended hook. If the hook turns out to need more than a few lines of fail-soft guarding, fall back to the recipe and say so in the commit body.
- The hook must stay fail-soft: a missing runner, a non-maintainer checkout, or a runner already running prints one line and exits 0.

## Detailed Requirements

- Detect an existing runner for this repository root (a pid file under `$TMPDIR/do-work-gate-runs/` is enough) and do not start a second one.
- Only start in this repository (the runner lives under `_dev/`, which is never shipped); an installed consumer's hook must not look for it.
- Print one line naming the log directory when it starts, nothing else.

## Constraints

- Land in place, not through `do-work run`; one integrating commit with version bump and changelog entry; prove it with one `bash _dev/tests/gate-runner.sh --once`.
- Delete before you add; every deleted test is listed in the commit body with the failure it pinned and why it no longer earns its cost. No new sentence pins, no new prose that walks a shell sequence.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.

## Red-Green Proof
**RED prompt/case:** Open a new session in this checkout and run `pgrep -f gate-runner.sh`.
**Why RED now:** nothing is running.
**GREEN when:** the runner is running after SessionStart, a second session starts no second runner, and a consumer checkout's hook output is unchanged.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes a shipped hook script.

## Full Context
See `do-work/user-requests/UR-104/input.md` for complete verbatim input.


## Addendum (2026-09-03)

User added (2026-09-03 22:05 local, "update the batch REQs per A1-A3 via queued addenda" referring to the velocity report at `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/`, item A1):

- `depends_on` changed from `[REQ-541]` to `[]`. This REQ edits the SessionStart hook and the runner; REQ-541 edits the status vocabulary and the work loop. Nothing here reads REQ-541's output. Landing the runner start early also removes one of the two gate runs per REQ for every REQ built before the test cuts finish.
- Coherence check: no contradiction with the original sections; the dependency change is the only frontmatter edit.
