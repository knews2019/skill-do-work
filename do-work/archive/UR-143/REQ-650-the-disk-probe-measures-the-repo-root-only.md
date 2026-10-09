---
id: REQ-650
title: 'The disk probe measures the repo root only'
status: completed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-09T16:09:30Z
created_at: 2026-10-09T16:06:43Z
user_request: UR-143
domain: backend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-651, REQ-652, REQ-653]
batch: october-review-triage
write_set: ["skills/do-work-board/tools/queue-kanban/verify.go", "skills/do-work-board/tools/queue-kanban/disk_space_unix.go", "skills/do-work-board/tools/queue-kanban/disk_space_windows.go", "skills/do-work-board/tools/queue-kanban/disk_space_test.go", "skills/do-work-board/docs/board-guide.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "skills/do-work-board/tools/queue-kanban/generate_test.go"]
claimed_at: 2026-10-09T16:08:09Z
dispatch_at: 2026-10-09T16:11:57Z
builder_handback_at: 2026-10-09T16:35:42Z
integration_at: 2026-10-09T16:36:20Z
review_at: 2026-10-09T16:41:48Z
kb_status: pending
heavy_verified_at: 2026-10-09T16:45:01Z
heavy_verified_revision: 674fa943ed0ce7ce37421ebbd73ffc7f03168963
commit: 674fa943ed0ce7ce37421ebbd73ffc7f03168963
completed_at: 2026-10-09T16:45:27Z
release_at: 2026-10-09T16:45:27Z
---
# The Disk Probe Measures the Repo Root Only
## What
The `low-disk-space` verify probe measures the filesystem holding the repo root and then every `worktree-agent-*` worktree, deduplicating by device id so two directories on one disk report once. Reduce it to one measurement of the repo root. Delete the worktree loop, the device-id dedupe and their tests; keep the thresholds, the single finding, and the Testing-page line.
## Why
Consumer review (validate-feedback 2026-10-09, F8, accepted): "the disk space check is over-engineered; it treats the computer like a complex network of different volumes when a simple check of the project folder would have sufficed." The earning incident in REQ-625 (the low-disk-space probe) §Why was growth inside the repo: browser-QA screenshots grew the repo to about 20 GB and the disk nearly filled. A repo-root-only check catches that. No archived document names an incident where a worktree on another volume filled up; the multi-device branch rests only on "worktrees can sit on another disk than the repo root" (`verify.go:216-217`). The maintainer chose "Repo root only".
## Verified Facts (from triage)
- `skills/do-work-board/tools/queue-kanban/verify.go:279-290` measures the repo root first, then each worktree path in name order; the worktree paths come from the shared `git worktree list --porcelain` read (`verify.go:217-223`, `1491-1510`).
- `verify.go:313-316` skips a device already reported through `measuredDevices[measurement.deviceIdentity]`; the comments at `:238-241` and `:265-269` say "so two directories on one disk are reported once" and "One finding per device".
- Each path costs one `Statfs` plus one `os.Stat` for the device id (`disk_space_unix.go:19-44`).
- The Testing page shows only the repo root's reading (`verify.go:293-302`, `generate.go:709-720`), so the Testing line does not change.
- REQ-625's review noted there was no test for "a worktree on the repo root's device, the common layout" (`REQ-625.md:241`).
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md:23` records the incident under family `disk-space-blind-spot`; it stays.
## Detailed Requirements
1. `verify.go`: the probe measures exactly one path, the repo root. Remove the worktree-path loop, the `measuredDevices` map, the device-identity field on the measurement type if nothing else reads it, and the comments that describe per-device reporting. If the worktree read at `:216-223` feeds only this probe, remove that call too; if the VERIFY band's other probes share it, leave the shared read alone.
2. `disk_space_unix.go` and `disk_space_windows.go`: drop the `os.Stat` device-identity read when nothing else uses it. Keep `syscall.Statfs` and `GetDiskFreeSpaceExW`, the build tags, and `disk_space_unsupported.go`.
3. `disk_space_test.go`: remove the dedupe and multi-worktree cases. Keep the threshold cases, the unsupported-platform case and the root-reading case. Add the RED test below.
4. Thresholds (warn below 10 GiB, critical below 3 GiB), `diskSpaceLevelFor`, the finding's Detail and Remedy text, the `.` path reduction and the Testing-page payload are unchanged.
5. Docs: `skills/do-work-board/docs/board-guide.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` no longer say "one finding per device" or "each builder worktree"; they say the probe measures the repo root.
6. Changelog entry (written by the integrator at release) says why: the earning incident was growth inside the repo, and the multi-device branch had no incident behind it.
## Constraints
- Deletion only. No new flag, no new finding, no change to thresholds or to what the run loop reads (it reads nothing from this probe).
- Read `_dev/primes/prime-kanban-board.md` and `prime-do-kanban.md` before editing; the board's versioning and build rules apply.
- Browser lane: set `QUEUE_KANBAN_BROWSER` when running the heavy lane; a skipped lane is not a pass.
## Dependencies
None. May run beside REQ-651, REQ-652 and REQ-653.
## Builder Guidance
Certainty is high: the triage located every line. Latitude: whether the measurement type keeps a device field for the Windows helper's shape, and how the root-only test is phrased.
## Red-Green Proof
**RED prompt/case:** A queue-kanban test builds a repo with one `worktree-agent-REQ-1-*` worktree on a second device and runs the disk probe with a fake measurer that records the paths it was asked to measure.
**Why RED now:** The measurer is called twice (repo root and the worktree) and the probe emits one finding per device; `grep -n measuredDevices verify.go` finds the dedupe map.
**GREEN when:** The measurer is called exactly once, with the repo root; `grep -n measuredDevices skills/do-work-board/tools/queue-kanban/verify.go` prints nothing; `go test ./...` in `skills/do-work-board/tools/queue-kanban/` is green; `queue-kanban verify` on this repo prints at most one disk finding.
**Validation:** User confirmed. The maintainer chose "Repo root only (Recommended)" in the triage questions and approved the plan.
## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8110 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: family `disk-space-blind-spot` is this probe's own incident, and the satellite governs the queue-kanban model this REQ edits.
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over budget; `slugged: partial`). Matching reason: its owning prime governs the board tool this REQ edits.
## Full Context
See `do-work/user-requests/UR-143/input.md` for complete verbatim input (both pastes, the maintainer's answers, and the full triage). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (from the builder hand-back) Read the brief, the REQ, UR-143 F8, the crew rules, both board primes, prime-releases and the `disk-space-blind-spot` lesson. `gitState.worktreePathsByName` is also read by `appendWorktreeFindings`, so the shared read stays and only the disk probe's local use goes. The `deviceIdentity` field has no reader once the dedupe goes, so it is removed on both platforms. The signature drops the worktree map. The RED test goes through `runVerifyProbes` with a real git repo and one `worktree-agent-REQ-1-fixture` worktree, so it compiles against both the old and new code and fails on the assertion. No lesson contradicts the plan.
- [x] **[APPLY]:** (from the builder hand-back) RED test written and run first. Then verify.go, the two platform files, the test files and the prime, exactly as planned. No new flag, finding, threshold or run-loop read.
- [x] **[UNIFY]:** (from the builder hand-back) `git diff --stat 2eb14357..b504b4a3`: 6 files changed, 108 insertions(+), 200 deletions(-); `git diff --check` clean. `go vet ./` exit 0, cross-OS vet for windows, linux and netbsd exit 0. `go test -count=1 ./` exit 0 (78 s, browser tests skipped without `QUEUE_KANBAN_BROWSER`). `gofmt -l .` prints nothing. `grep -n measuredDevices verify.go` prints nothing; the per-device phrases are gone from board-guide.md and prime-do-kanban.md. Live `verify` on the main tree printed zero `low-disk-space` findings. Files checked: verify.go, disk_space_unix.go, disk_space_windows.go, disk_space_test.go, generate_test.go, prime-do-kanban.md; no debug artifacts, no build outputs.
*Source: consumer review "the disk space check is over-engineered; it treats the computer like a complex network of different volumes when a simple check of the project folder would have sufficed", accepted by `do-work-toolbox validate-feedback` on 2026-10-09 as F8; maintainer answer "Repo root only (Recommended)".*

---

## Triage

**Route: A** - Simple

**Reasoning:** Names every file and line to delete (verify.go worktree loop and device dedupe, the disk_space helpers' device read, the dedupe tests) and keeps everything else; a focused deletion with a runnable RED test in the existing Go harness.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/disk_space_unix.go` (modified)
- `skills/do-work-board/tools/queue-kanban/disk_space_windows.go` (modified)
- `skills/do-work-board/tools/queue-kanban/disk_space_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified)

**What was done:** `appendDiskSpaceFindings` now measures the repo root once; the worktree-path loop, the device-id dedupe map, the `deviceIdentity` field and both platform files' device reads are deleted, while thresholds, finding text, the `.` reduction and the Testing-page reading are unchanged. The dedupe and multi-worktree tests are gone, `TestDiskSpaceProbeMeasuresOnlyTheRepoRoot` is new, `generate_test.go` follows the removed field and signature, and the prime's disk trap now says repo root only (merge 674fa943, builder commit b504b4a3).

## Decisions

(from the builder hand-back)

- **D-01 DECIDE & STATE:** `generate_test.go` was edited, outside the write_set. The removed `deviceIdentity` field and the old four-argument signature forced it. The payload test still requires two low-disk-space findings reduced to `.`, so it still proves the reduction.
- **D-02 DECIDE & STATE:** `board-guide.md` was left unchanged. It never said "one finding per device" or "each builder worktree", and its only disk line already says the toolbar shows "the repo root's free disk space". Adding text would go against the deletion-only constraint.
- **D-03 DECIDE & STATE:** the measurement type keeps only free and total bytes. The REQ allowed keeping a device field for the Windows helper's shape, but nothing reads it. Removing it also removes the FNV hash and its known limit for folder-mounted volumes.
- **D-04 DECIDE & STATE:** the "disk-space probe for worktrees: <error>" skip line is gone along with the loop. A worktree-list failure is still reported by the worktree probes' own "worktree probes: ..." skip.
- **D-05 DECIDE & STATE:** no version bump or comment change in the Go tool. Under prime-kanban-board.md the tool rides the suite version and the root changelog, so the integrator owns that.

(integrator)

- **D-06 DECIDE & STATE:** D-01's write-set contradiction is accepted as the required-class case of work.md's "Write only inside the declared scope", not a scope expansion. The REQ's own GREEN requires `go test` green, and `generate_test.go` stops compiling once the field and the four-argument signature are gone, so the test class was required. `write_set` is extended with `skills/do-work-board/tools/queue-kanban/generate_test.go` so qualification judges it as declared. Coordinator ruling recorded in the integrator brief.

## Discovered Tasks

(from the builder hand-back)

- **impact-negligible** `skills/do-work/CHANGELOG.md` entry for 0.305.60 still says the probe measures "the repo root and each builder worktree, one finding per device". It is a historical entry and should stay as written. (The builder wrote `impact-cosmetic`, which is not an impact token; the integrator reclassified it as `impact-negligible`: no reader acts on a past release note.) → report only

## Qualification

**Gate record:** `advance --diff-range 54056e00..674fa943` returned the `qualify` gate `satisfied` with provenance `merged_range` and no findings (no debug artifact, no output primitive, P-A-U boxes ticked, no scope drift against the six-file `write_set`).

**Requirement trace against the merged files:**
1. `verify.go`: `appendDiskSpaceFindings(report, repoRoot, measure)` calls the measurer once with the repo root. The worktree-path loop, the `measuredDevices` map, the `deviceIdentity` field and the per-device comments are gone (`grep -rn "measuredDevices\|deviceIdentity"` over the tool prints nothing). The shared `gitState.worktreePathsByName` read stays because `appendWorktreeFindings` reads it (`verify.go:1246`), as requirement 1 asked.
2. `disk_space_unix.go` drops the `os.Stat` device read and the `os` import; `disk_space_windows.go` drops the FNV volume hash and the `hash/fnv` and `strings` imports. `syscall.Statfs`, `GetDiskFreeSpaceExW`, the build tags and `disk_space_unsupported.go` are unchanged.
3. `disk_space_test.go`: the per-device and per-directory-error tests are removed, the threshold, unsupported-platform and root-reading tests stay with their worktree halves cut, and `TestDiskSpaceProbeMeasuresOnlyTheRepoRoot` is new.
4. Thresholds, `diskSpaceLevelFor`, the Detail and Remedy strings, the root skip-reason strings and the `RepoRootDiskSpace` reading are byte-identical in the new body; the `.` reduction is still proven by `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths`, which now requires two `.` subjects and no other.
5. `prime-do-kanban.md`'s disk trap now says "the filesystem holding the repo root, the only path it measures". `board-guide.md` needed no edit: its only disk line already says "the repo root's free disk space" (D-02).
6. The changelog entry is the integrator's, written at finalization.

**Scope:** touched files = `write_set` minus `board-guide.md` (declared, correctly untouched) — six of seven. `generate_test.go` is declared by D-06's extension; its edit is the minimum the removed field and signature force.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge 674fa943 (passed to `advance` as `bash _dev/tests/maintainer-verify.sh`, exit 0); the REQ probe `do-work/runs/work-2026-10-09-160815/REQ-650-probe.sh` run by `advance` (`grep measuredDevices` absent plus `go test -run 'DiskSpace|LowDisk|FreeDisk'`, 3 s); `go test -count=1 -run 'DiskSpace|TestGeneratedVerifyPayloadCarriesNoAbsolutePaths' -v .` in `skills/do-work-board/tools/queue-kanban/` at the merge.
**Result:** ✓ All passing. Gate exit 0 in 127 s with every stage executed (`reuse_disabled`): queue-kanban 422 tests in 42 s, do-work-cli 879 tests in 59 s, contract regressions passed. `advance` recorded `green-gate` and `run-blocked-check` satisfied (`BLOCKED-PROBE-SUCCEEDED`). Uncached disk tests: 10 pass; `TestJavaScriptBehaviorTestingDiskSpaceLineReadsEveryCase` skips here and runs in the heavy JavaScript lane.

**Red-green validation:**
- `TestDiskSpaceProbeMeasuresOnlyTheRepoRoot` (`disk_space_test.go`): ✗ at base 2eb14357 with only the test added (`measured [<repo root> <worktree-agent-REQ-1-fixture path>], want exactly the repo root`) → ✓ at builder commit b504b4a3 and again uncached at merge 674fa943. It drives `runVerifyProbes` against a real git repo with one `worktree-agent-REQ-1-fixture` worktree and a recording fake measurer, which is the captured RED case. The captured case put the worktree "on a second device"; the test does not need a second device, because the GREEN condition is "the measurer is called exactly once, with the repo root", and a recording fake proves that on any layout.

**New tests added:**
- `TestDiskSpaceProbeMeasuresOnlyTheRepoRoot` (`skills/do-work-board/tools/queue-kanban/disk_space_test.go`)

**Existing tests updated (cross-REQ impact):**
- `disk_space_test.go` (from REQ-625, the low-disk-space probe): `TestDiskSpaceProbeReportsOneFindingPerDevice` and `TestDiskSpaceProbeMeasurementErrorSkipsOnlyThatDirectory` removed with the branch they covered; the root skip-reason and unsupported-platform tests lost their worktree halves. Intentional behaviour change.
- `generate_test.go` (from REQ-625): the fake measurer drops `deviceIdentity`; the payload test seeds a second root finding instead of a worktree finding and requires two `.` subjects and no other. Intentional behaviour change (D-06).

**Heavy verification plan:**
- Range: 54056e0073f8380133f6c1e9ce4a4660c26247de..674fa943ed0ce7ce37421ebbd73ffc7f03168963
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — all six changed paths match subtree `skills/do-work-board/tools/queue-kanban`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — all six changed paths match subtree `skills/do-work-board/tools/queue-kanban`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all six changed paths match subtree `skills`

*Verified by work action*

## Review

**Overall: 98%** | 2026-10-09T16:41:48Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None. Nits: F1 `verify.go:258` builds `repoRootReading` before the error switch that never uses it — impact-negligible → report only; F2 `generate_test.go:3165` seeded second repo-root finding repeats the live probe's — impact-negligible → report only
**Acceptance:** Pass — the measurer is called once with the repo root; disk tests and vet are green at 674fa943, and RED at base 2eb14357 matches the test code.
**Restatement sweep:** redefined the low-disk-space probe's scope ("repo root and each worktree, one finding per device" became "repo root only"); grepped skills/, _dev/, docs, README and both CHANGELOGs for "per device", "each builder worktree", "deviceIdentity", "measuredDevices", "disk-space probe for worktrees", and disk near worktree-agent. Only hits are the historical 0.305.60 entry in skills/do-work/CHANGELOG.md:222 and its root mirror CHANGELOG.md:222 (expected to stay). prime-do-kanban.md:36, board-guide.md:78 and actions/board.md:98 all say repo root. No stale restatements.
**Suggested testing:** 0 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** A RED test that drives `runVerifyProbes` against a real git repo with a recording fake measurer. It compiles against both the old and the new signature, so RED failed on the assertion instead of on the build, and it proves "measured once, with the repo root" without needing a second physical device.
**What didn't:** The REQ's `write_set` missed `generate_test.go`, which seeded a worktree disk finding through the old four-argument call and set the removed field. A deleted field or signature reaches every test that constructs it, so a deletion REQ's write set has to include those tests (D-01, D-06).
**Worth knowing:** A probe should watch where its earning incident happened. The per-worktree, per-device branch was built for a layout ("worktrees can sit on another disk") that no incident ever produced, and it cost a device read per path, an FNV volume hash on Windows and two tests. The 0.305.60 changelog entry still describes the old scope and stays as history.

## Orientation

Now `queue-kanban verify`'s `low-disk-space` probe takes one reading of the filesystem that holds the repo root, the same reading the Testing page shows; it lives in the board tool's VERIFY band (`skills/do-work-board/tools/queue-kanban/verify.go`, primes `prime-do-kanban.md` and `_dev/primes/prime-kanban-board.md`). Prime spot-check: every path the disk trap names still exists; no prime went stale.

## Heavy Verification Plan

- Base: 54056e0073f8380133f6c1e9ce4a4660c26247de
- Target: 674fa943ed0ce7ce37421ebbd73ffc7f03168963
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — reasons: the six changed paths under `skills/do-work-board/tools/queue-kanban/` (disk_space_test.go, disk_space_unix.go, disk_space_windows.go, generate_test.go, prime-do-kanban.md, verify.go) matched subtree `skills/do-work-board/tools/queue-kanban`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — reasons: the same six paths matched subtree `skills/do-work-board/tools/queue-kanban`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — reasons: the same six paths matched subtree `skills`

## Heavy Verification Result

- Target: 674fa943ed0ce7ce37421ebbd73ffc7f03168963
- Execution revision: 674fa943ed0ce7ce37421ebbd73ffc7f03168963 (detached checkout, `QUEUE_KANBAN_BROWSER` set to Google Chrome)
- queue-kanban-javascript: exit 0, executed, 8 s
- queue-kanban-browser: exit 0, executed, 89 s (no `HEAVY-RUN-LANE-SKIPPED`)
- staged-skills: exit 0, executed, 45 s

## Timing

Observed 2026-10-09T16:36:05Z to 2026-10-09T16:45:02Z: 8m 57s total, 6m 30s attributed across 4 events, 2m 27s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 4m 55s | 2 |
| review | 1m 20s | 1 |
| handback-merge | 15s | 1 |

Slowest stage: review / review agent, 1m 20s, outcome success.
Slowest command: verification-gate / heavy drain, 2m 40s, exit 0, .
