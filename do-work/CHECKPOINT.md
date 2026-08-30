---
session_updated: '2026-08-30T17:46:00Z'
session_ended: '2026-08-30T17:46:00Z'
last_completed: REQ-425
queue_state: [14 pending, 1 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress]
reqs_processed_this_session: 4
session_depth: moderate
---

# Session Checkpoint

## Completed This Session

- [REQ-406](archive/REQ-406-create-do-work-cli-foundation.md): Shared `do-work-cli` runtime and Git transaction foundation; implementation `2ca25d7`; released 0.245.0.
- [REQ-390](archive/REQ-390-timeline-trailing-window-periods.md): Timeline trailing windows replacing calendar periods; implementation `59105df`; released 0.246.0.
- [REQ-407](archive/REQ-407-migrate-install-update-bootstrap-to-go.md): Bootstrap, install, update, reconciliation, validation and fetching migrated into Go; implementation `f45cdca`; released 0.247.0.
- [REQ-425](archive/REQ-425-trailing-window-bound-assumptions.md): Trailing-window bound assumptions closed; implementation `04b8120`; released 0.248.0.

## In Progress (interrupted)

None. `do-work/working/` holds only `baseline.json`; no worktrees and no `worktree-agent-*` branches remain.

## Still Queued

REQ-408 through REQ-420 (the serial Go-platform chain, each gated on its predecessor), plus REQ-426 (mode-bit regression, ungated) and REQ-427 (`pending-answers` — the Go version floor).

## Session Notes

Every REQ ran the full pipeline with a worktree builder and an adversarial review: five dimensions, then three refutation lenses per finding. That review found real defects each time — a 10% flaky test already inside the canonical gate (REQ-406), two edge-of-range navigation defects (REQ-390), a goroutine race in the install signal handler and a BOM install-blocker (REQ-407) — while refuting the large majority of findings by measurement rather than argument.

Two process defects were fixed along the way. The review workflow counted an errored verifier as a refutation, which would have silently dismissed 24 unadjudicated findings; it now treats fewer than two live votes as unverified. And `tools/checks/scope-drift.sh` reads every backticked token in a "Files I will touch" bullet as a declared path, so a code span in the rationale produces phantom drift — folded onto REQ-414, which ports those checks to Go.

Usage limits interrupted the session three times. The recovery that mattered: a builder died with eight modified files and five new packages uncommitted, and only the worktree saved it. Every builder brief since instructs committing each unit as its tests pass rather than batching to the end.

## Environment

This container did not ship the maintainer toolchain. Installed during the session, outside the repo: ShellCheck 0.11.0, `just` 1.43.0, and Go 1.26.1 via `go env -w GOTOOLCHAIN=go1.26.1+auto`. Chrome for Testing 151.0.7922.174 was fetched to the scratchpad because the container's own Chromium is 141 — the build REQ-375 explicitly deprecated, and it reproduces that REQ's recorded vacuity-guard failure verbatim.

The canonical gate runs with its optional browser lane in the DEFAULT SKIPPED state. Do not set `QUEUE_KANBAN_BROWSER` for the gate: `TestBrowserBehaviorCompletionCompanionsKeepReadableContrast` fails in this Linux container at HEAD on both Chrome 141 and 151, with no source changes, where REQ-375 recorded a whole-lane pass on macOS. Timeline probes were run individually under Chrome 151 instead, where all 21 pass.

## Cleanup

All four builder worktrees removed and their branches deleted with `git branch -d` from the integration branch, which is the assertion that each merge actually landed. `do-work/working/baseline.json` is the pre-flight record and is left in place.
