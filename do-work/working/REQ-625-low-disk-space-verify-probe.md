---
id: REQ-625
title: 'Add a low-disk-space probe to the board VERIFY band'
status: claimed
created_at: 2026-10-02T13:58:27Z
user_request: UR-131
domain: backend
prime_files: ["_dev/primes/prime-kanban-board.md", "_dev/primes/prime-releases.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set: ["skills/do-work-board/tools/queue-kanban/verify.go", "skills/do-work-board/tools/queue-kanban/verify_test.go", "skills/do-work-board/tools/queue-kanban/disk_space_*.go", "skills/do-work-board/tools/queue-kanban/generate_test.go", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md", "CHANGELOG.md", "skills/do-work/CHANGELOG.md", "VERSION", "skills/do-work/VERSION", "skills/do-work/actions/version.md"]
claimed_at: 2026-10-02T14:03:40Z
---

# Add a Low-Disk-Space Probe to the Board VERIFY Band

## What
Add one read-only verify probe to the board tool that measures free disk space on the filesystem holding the repo root (and each enumerated `worktree-agent-*` worktree on a different device) and emits a `low-disk-space` finding below fixed thresholds. The finding reaches the CLI `verify` report and the board's VERIFY band through the existing `attachVerifyFindings` path, with no new UI.

## Why
During a long `do-work run --fan-out 2` on a game repo (October 2026), browser-QA runs wrote hundreds of megabytes of screenshots per comparison. The repo grew to about 20 GB in a few hours and the disk nearly filled, and nothing in the suite noticed. An earlier run (REQ-504, commit 1a726b6f) also lost its heavy test lanes to a full disk. Builders, gates and git fail in confusing ways at zero bytes free, and the VERIFY band is where the run operator already looks.

## Detailed Requirements
1. **Helper, build-tagged.** A `diskSpaceFor(directory) (free, total uint64, device uint64, err)` style helper split like `atomic_replace_*.go`: `syscall.Statfs` on unix, `GetDiskFreeSpaceExW` via `syscall.NewLazyDLL("kernel32.dll")` on Windows, and an unsupported file returning a typed "unsupported on GOOS" error. **The unix tag list must be narrower than `atomic_replace_unix.go`**: aix, illumos and solaris have no `syscall.Statfs` and belong to the unsupported file. All `Statfs_t` arithmetic stays inside the platform file (field types differ per OS); use `Bavail`, not `Bfree`. No new module dependency; `go.mod` stays as it is.
2. **Probe.** `appendDiskSpaceFindings(&report, repoRoot)` in `verify.go`, called from `collectVerifyFindings` after `appendWorktreeFindings`. It measures the repo root and every path returned by `listWorktreeAgentWorktrees`, dedupes by device id (`os.Stat(dir).Sys().(*syscall.Stat_t).Dev` cast to uint64 on unix; volume path or serial on Windows, inside the platform file), and emits at most one finding per device. On the unsupported error it appends to `SkippedProbes` (the invariant is real and went unverified) and never fails the run. A `Statfs` error on one directory is a skipped probe for that directory, not a crash.
3. **Thresholds.** Two named constants: `lowDiskSpaceWarnBytes = 10 GiB` and `lowDiskSpaceCriticalBytes = 3 GiB`. No percentage rule (dropped at triage: on a 20 GB disk a 10 % warn is 2 GiB, so critical would fire before warn).
4. **Finding shape.** `Category: verifyCategoryLowDiskSpace = "low-disk-space"`, `Subject: <measured directory>` (the repo root, which reduces to `.`, or the worktree path, which reduces to its name), `Fixable: false`. Detail: `"<free> free of <total> (<pct>%) — below the <warning|critical> threshold (<threshold>)"` with human-readable GiB values. Remedy, generic and project-agnostic: `"free space: clear regenerable QA output, finished builder worktrees (do-work cleanup), browser caches; `du -sh * | sort -h` at the repo root shows the largest directories"`. The suite must not name any project's output folders. **No largest-directories walk** in the probe (pushed back at triage: sizing direct children means walking them fully, so the 200 ms cap would make every result "(partial)" on exactly the repos that need it, on every serve request).
5. **Paths.** Nothing new: `attachVerifyFindings` already reduces Subject, Detail and Remedy. Confirm by test that the board payload carries no absolute path for this finding.
6. **Read-only.** The probe measures and reports. It deletes nothing. Cleanup stays human-consented through `actions/cleanup.md`; `Fixable` stays false because cleanup cannot resolve disk space.
7. **Docs.** Add the probe to `prime-do-kanban.md` where `verify.go`'s probes are named. Add one lesson to `lessons-do-kanban.md` under `[family: disk-space-blind-spot]` with the file's existing entry shape: a long fan-out run with browser QA can fill a disk in hours, and the board is where the operator looks. Refresh the `do-work/lessons-index.md` row for that satellite (families and token count) in the same commit.
8. **Release.** Shipped files change, so this is a release: changelog entry with a descriptive title and a version bump per `_dev/primes/prime-releases.md`. Forensics Check 14 counts every verify finding as a warning and needs no change.

## Out of scope (pushed back at triage, do not build)
- An always-visible free-space chip in the board header.
- A top-3 largest-directories walk inside the finding.
- Any sentence in `work.md` or `work-reference.md` that pauses dispatch on this finding. If a dispatch pause is wanted later, it belongs in queue-mode `advance` as a typed exclusion in the core CLI, never as prose keyed on a board-package finding.

## Constraints
- Read-only probe; no deletion.
- No new Go module dependency.
- Thresholds are constants, not configuration.
- The probe runs on every serve request (findings are computed outside the mtime cache), so it must stay at one `Statfs` call per measured directory and nothing heavier.

## Dependencies
None. `listWorktreeAgentWorktrees`, `SkippedProbes`, `attachVerifyFindings` and `reduceAbsolutePaths` all exist today.

## Builder Guidance
High certainty on scope: the triage report fixed the shape and the exclusions above. Builder latitude on helper naming, on the exact Detail wording, and on where the Windows device identity comes from. Tests use an injected measurement function (a package-level `var` the test swaps), never the real filesystem, because real free space is nondeterministic.

## Red-Green Proof
**RED prompt/case:** With the measurement function faked to report 2 GiB free of 500 GiB for the repo root, `collectVerifyFindings` over a clean synthetic board returns no `low-disk-space` finding. With the function faked to return the unsupported error, `SkippedProbes` has no disk-space entry.
**Why RED now:** No probe exists; `collectVerifyFindings` ends at `appendWorktreeFindings` (`verify.go:195`).
**GREEN when:** The same inputs yield exactly one finding with `Category "low-disk-space"`, `Subject "."` after reduction, `Fixable false`, Detail naming the critical threshold; 20 GiB free yields no finding; 9 GiB yields a warning-level Detail; two worktree paths on one device yield one finding; the unsupported error yields one `SkippedProbes` entry and zero findings; the board payload test `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` still passes with the new finding present.
**Validation:** Inferred during capture

## Tests (each names the failure it pins)
- Fake measurement above and below the thresholds: pins "probe silently clean when disk is low".
- Boundary at exactly 3 GiB and exactly 10 GiB: pins an off-by-one at the constants.
- Two worktrees on the same device produce one finding: pins duplicate findings per device.
- Unsupported platform lands in `SkippedProbes`, never silently clean: pins "unknown reads as clean".
- Extend `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` (`generate_test.go:3125`): pins a machine path leaking into a shareable snapshot.

## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (6908 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matched: changing queue-kanban model and testing. Read it anyway if time allows; the `unknown-reads-as-clean` family is the closest prior.
- `_dev/primes/lessons-kanban-board.md` (5912 tokens, over budget; `slugged: partial`). Matched: finding-detail text.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Full Context
See `do-work/user-requests/UR-131/input.md` for the complete verbatim input and the original proposal.

*Source: capture the ones that are accepted*
