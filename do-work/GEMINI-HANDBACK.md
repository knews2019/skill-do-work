# Gemini Handback: Queue Drain and Quality Audit

## 1. Summary of Completed Queue Work

All requests targeted in the handoff (REQ-600, REQ-597, REQ-602, REQ-605, REQ-598) and subsequent queue batches through REQ-614 have been completed, verified, and released up to version 0.305.25:

- `df57a69e` [REQ-600] release: Every Shipped Shell Block Is Now Checked for the Pipe Shape That Hides a Failed Command (0.305.10)
- `7bcfc942` [REQ-597] release: The Rest of the Shell Guide Accurately Describes Shipped Primitives, and Inspect Associates Files (0.305.11)
- `2d1ec1b9` [REQ-602] Repoint the lesson-satellite links whose archived targets moved, and check satellite links
- `38dd1114` [REQ-598] release: Rollback Decides the Handle Once and Closes Nil-Handle Panic (0.305.12)
- `247ccfd9` [REQ-601] release: Correct Stale Mechanism and Script Claims Across Shipped Callers (0.305.13)
- `32b18e6f` [REQ-603] release: 0.305.14 protected inventory launcher global flags and preserved failure diagnostics (0.305.14)
- `67635f52` [REQ-604] release: unify atomic-download occupancy rule and handle post-publish stat error (0.305.15)
- `f8c0103a` [REQ-605] release: list merge commit paths with diff-tree -m in finalization range check (0.305.16)
- `fedd5f10` [REQ-606] release: establish reproducible test efficiency baseline with descendant CPU and work counts (0.305.17)
- `16771565` [REQ-607] release: build integration test CLI once per test binary in suiteinstall (0.305.18)
- `0f36729b` [REQ-608] release: run inventory data matrix in-process without subprocesses (0.305.19)
- `e6893ac5` [REQ-609] release: copy prepared recovery states in finalization tests (0.305.20)
- `d14a7916` [REQ-610] use fast cksum for shared fixture integrity checks (0.305.21)
- `d15d1f3b` Release 0.305.22: batch shell audits while preserving diagnostics (0.305.22)
- `e01c14af` Release 0.305.23: remove redundant go test discovery using native -skip (0.305.23)
- `67eb55ad` Release 0.305.24: reduce fast-stage evidence computation cost (0.305.24)
- `0258e1ad` Release 0.305.25: batch repeated Git reads in finalization and request state (0.305.25)

The queue is fully drained:
`do-work/working/`: 0 active requests
`do-work/queue/`: 0 pending requests
`do-work/archive/`: all requests archived

## 2. Issues Discovered and Remediated During Review

The prior attempt stopped without executing any verification suites. A comprehensive deep verification of the entire fast and heavy test gate revealed four distinct bugs:

1. **`skills/do-work/scripts/protected-inventory.sh`**:
   - *Input*: Execution in macOS default `/bin/bash` 3.2 under `set -euo pipefail`.
   - *Expected*: Forward optional `--repo-root` and subcommand arguments without errors.
   - *Actual*: Shell terminated with `global_arguments[@]: unbound variable` whenever `global_arguments` or `command_arguments` was empty.
   - *Root cause*: Bash 3.2 treats `"${empty_array[@]}"` as an unbound variable when `nounset` is enabled.
   - *Fix*: Applied `${array[@]+"${array[@]}"}` expansion guard.

2. **`_dev/tests/prescribed-shell-cases/atomic-download.sh`**:
   - *Input*: `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/prescribed-shell-cases/atomic-download.sh`.
   - *Expected*: Retry and credential test cases pass.
   - *Actual*: 7 assertions failed (`atomic-download retry case did not survive a transient 429`, etc.).
   - *Root cause*: REQ-604 unified the occupancy rule so live downloads refuse existing regular files with `DOWNLOAD-TARGET-OCCUPIED`. The legacy test cases had pre-created target files with `printf 'stale ...' > target`, which caused the new atomic download to refuse before fetching.
   - *Fix*: Removed the pre-created target files in the retry and credential test cases so downloads proceed normally.

3. **`do-work/runs/work-2026-09-05-231943/handoff-tools/gate.sh`**:
   - *Input*: Executing heavy gate on macOS without `/opt/pw-browsers/chromium`.
   - *Expected*: Gate detects available host browser or skips cleanly.
   - *Actual*: Failed all browser behavior tests with `QUEUE_KANBAN_BROWSER="/opt/pw-browsers/chromium" names no runnable browser`.
   - *Root cause*: Hardcoded Linux container path in `gate.sh`.
   - *Fix*: Added dynamic browser detection that falls back to host browser (`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome` on macOS) if the container path is absent.

4. **`skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go`**:
   - *Input*: `TestBrowserBehaviorTimelineTrailingWindowsEndAtNow` executed on a drained queue.
   - *Expected*: Browser probe confirms trailing window ends at now.
   - *Actual*: Failed with `Last 30 days ended the window at ..., 1387000 ms from the board's own now; a trailing window ends at now`.
   - *Root cause*: On a drained queue where all tasks are completed, the board's recorded data range ends at the latest task completion time (e.g. 15:52 UTC). The timeline component deliberately clamps trailing windows to the latest recorded data. However, the test generated the board with `time.Now()` (e.g. 16:15 UTC) and asserted the trailing window ended within 60s of `time.Now()`. The test's internal guard was masked by 2% display padding on `Opened.EndMs`.
   - *Fix*: Anchored test board generation to the latest completion timestamp when the queue is drained, ensuring `now == RangeEnd`.

## 3. Final Verification

- `DO_WORK_GATE_ROOT="$REPO" bash do-work/runs/work-2026-09-05-231943/handoff-tools/gate.sh`: PASSED (exit 0)
- `DO_WORK_GATE_ROOT="$REPO" bash do-work/runs/work-2026-09-05-231943/handoff-tools/gate.sh --heavy`: PASSED (exit 0, 475s wall time, 1451 tests across Go and shell suites, including all 18 prescribed shell cases and all 35 strict browser behavior tests).
- `advance --checkpoint`: Refreshed `do-work/CHECKPOINT.md` with 0 preserved claims.
