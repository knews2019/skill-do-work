---
id: REQ-500
title: 'Surface unfinished finalizations in doctor and the session-start banner'
status: claimed
created_at: 2026-09-02T13:31:12Z
user_request: UR-097
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-498]
related: [REQ-499, REQ-501]
batch: run-with-recovery
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set: [skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan.go, skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan_test.go, skills/do-work/tools/do-work-cli/internal/hookcommands/session_start.go, skills/do-work/tools/do-work-cli/internal/hookcommands/session_start_test.go, _dev/tests/session-start-hook-behavior.sh]
claimed_at: 2026-09-02T16:41:44Z
planning_at: 2026-09-02T16:43:00Z
dispatch_at: 2026-09-02T16:44:00Z
builder_handback_at: 2026-09-02T16:59:18Z
integration_at: 2026-09-02T17:00:04Z
review_at: 2026-09-02T17:15:12Z
remediation_at: 2026-09-02T17:37:20Z
---

# Surface Unfinished Finalizations in Doctor and the Session-Start Banner

## What
Give the archived-but-uncommitted state a name. `doctor` gains two findings, and the SessionStart hook prints one read-only pointer line when either fires.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the CLI prime, touch-conditional lesson satellite, builder rules, and existing doctor/SessionStart seams. Add doctor-owned read-only detection for canonical finalization journals and Git-dirty archived terminal-success REQs with blank commits; project the first deterministic finding through SessionStart; prove valid tails, committed-archive false-positive control, exact output, and byte identity.
- [x] **[APPLY]:** Added `UNFINISHED-FINALIZATION` and `ARCHIVED-WITHOUT-COMMIT` safely-refused findings, canonical recovery argv/evidence, and a one-line SessionStart projection sourced from the same typed doctor findings. Added Go and retained-launcher regression coverage.
- [x] **[UNIFY]:** Reviewed `git diff --stat`, every changed file, `git diff --check`, and added-line debug-artifact scans. Verified doctor detection and Git false-positive controls, hook aggregation and exact banner behavior, retained launcher byte identity, and shell portability. The full canonical maintainer gate passed.

## Why
Today nothing reports the REQ-494 shape: `doctor` skips archived REQs with a blank `commit:` (`internal/doctor/doctor_scan.go:302-305`), `PREFLIGHT-DIRTY` filters `do-work/` out, `cleanup` refuses untracked archived REQs, and the session-start hook never reads `working/`, the checkpoint, or git status. The first sign of trouble was a claim refusal on an unrelated REQ.

## Context
The hook (`internal/hookcommands/session_start.go`) already holds a repository snapshot and already prints a banner, so one extra line costs nothing. Detection lives in `doctor` because `hookcommands` already imports it; `doctor` must not import `finalization` (cycle risk), so it reads the journal directory by its on-disk contract (`git rev-parse --git-path do-work-finalization`, one `*.json` per REQ with `phase` and `manifest.request_id`) using a minimal local struct.

## Detailed Requirements
- `UNFINISHED-FINALIZATION`: a finalization journal exists for a REQ. Evidence names the REQ and its `phase`. Fixability refused; next argv is REQ-498's `recover-finalization`.
- `ARCHIVED-WITHOUT-COMMIT`: an archived terminal-success REQ has a blank `commit:` **and** its file is untracked or modified relative to HEAD. The git condition is the false-positive control: many old archives legitimately carry a blank `commit:` and are committed. Same fixability and next argv; evidence also names `do-work run-with-recovery` for the sole-writer case.
- The existing divergence skip at `doctor_scan.go:302-305` stays; this is a separate scan placed beside it.
- session-start prints exactly one line, only when a tail is detected: `do-work: unfinished finalization for REQ-NNN — 'do-work run' resumes it; 'do-work run-with-recovery' if this checkout is the only writer.` No mutation. Existing exact-banner tests stay unchanged when no tail exists.

## Constraints
- Reading is not authority: the hook reports and never repairs (REQ-467 history on SessionStart authority restatements).
- Disjoint write set from REQ-499: this REQ never writes `internal/finalization/`.

## Dependencies
Depends on REQ-498 (Make orchestrator finalization resumable) for the journal directory contract. Independent of REQ-499; runs in parallel with it.

## Builder Guidance
Firm. Follow the adjective-noun finding naming already in doctor (`STUCK-WORK`, `HOLLOW-COMPLETION`, `STRANDED-TERMINAL-REQUEST`).

## Red-Green Proof
**RED prompt/case:** A doctor test fixture with an archived `completed` REQ carrying a blank `commit:` whose file is untracked; and a second fixture with a journal file present. Run `doctor` and `session-start`.
**Why RED now:** `doctor` returns no finding for either fixture and the session-start banner is unchanged.
**GREEN when:** The first fixture yields `ARCHIVED-WITHOUT-COMMIT`, the second yields `UNFINISHED-FINALIZATION` with its phase, session-start prints the pointer line in both, and a control fixture with the same REQ committed at HEAD yields no finding and no line.
**Validation:** Inferred during capture from the audited incident.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2299 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on family `opaque-evidence-projection` (typed findings and evidence shape).

## Implementation Summary

- Added doctor-owned read-only detection for unfinished finalization journals and dirty archived terminal-success requests with blank provenance.
- Projected the first deterministic tail finding through SessionStart as the exact one-line recovery pointer.
- Added package and retained-launcher regressions for both tail shapes, committed-archive false positives, exact banner output, and byte identity.

## Decisions

None.

## Testing

- RED: `go test ./internal/doctor ./internal/hookcommands` failed because both findings and the SessionStart pointer were absent.
- GREEN: focused doctor and hookcommands tests passed.
- `bash _dev/tests/session-start-hook-behavior.sh` passed.
- `go vet ./...`, `go test -count=1 ./...`, and `go test -race ./internal/doctor ./internal/hookcommands` passed.
- Go 1.25 compatibility and `bash _dev/tests/maintainer-verify.sh` passed; the optional browser lane skipped because no browser was available.

## Review — Attempt 1

**Overall: 78%** | **Acceptance: Partial** | **Risk: Low**

Independent review approved the valid-tail behavior with follow-ups: both typed doctor findings, the exact one-line SessionStart projection, byte identity, and the committed blank-provenance false-positive control work. One bounded remediation is required for two Important findings:

- Preserve unreadable/malformed journal and Git-inventory observation failures as typed evidence instead of silently treating them as no unfinished finalization; SessionStart must still project the recovery pointer.
- Replace the per-archive `git status` subprocess loop with one repository-level inventory/probe so the common legacy blank-provenance population does not make SessionStart materially slow.
- Persist negative tests for malformed/unreadable journals, Git failure, scale behavior, and the committed-archive/no-banner control at the hook seam.

The reviewer measured 200 committed blank-provenance archives at 1.62 seconds versus 0.03 seconds with nonblank provenance. The stale README SessionStart summary was routed to `do-work/prose-backlog.md` under the fold-first prose rule.

## Remediation

The single remediation pass closed both Important findings. Doctor now reads one repository-wide NUL-delimited archive inventory, keeps committed blank-provenance archives quiet, and emits typed `FINALIZATION-TAIL-INSPECTION-FAILED` evidence when inventory observation fails. Malformed, unreadable, nonregular, or identity/phase-invalid canonical journal files retain `UNFINISHED-FINALIZATION` identity with `phase=unknown`, and SessionStart projects the same exact recovery pointer from the first affected request.

Persisted regressions cover malformed/unreadable journals, corrupt-index Git failure, 200 archives with one inventory probe, committed-archive/no-banner at the hook seam, recovery argv, byte identity, and read-only behavior. Focused, race, vet, full-module, Go 1.25, retained hook, and canonical maintainer gates passed. Remediation builder `cd72002a`; merge `ad260b5d`.

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-097/input.md` for complete verbatim input.

---
*Source: capture of the run-with-recovery request (UR-097).*
