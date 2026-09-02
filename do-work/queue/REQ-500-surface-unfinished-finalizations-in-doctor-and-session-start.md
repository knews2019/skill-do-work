---
id: REQ-500
title: 'Surface unfinished finalizations in doctor and the session-start banner'
status: pending
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
---

# Surface Unfinished Finalizations in Doctor and the Session-Start Banner

## What
Give the archived-but-uncommitted state a name. `doctor` gains two findings, and the SessionStart hook prints one read-only pointer line when either fires.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-097/input.md` for complete verbatim input.

---
*Source: capture of the run-with-recovery request (UR-097).*
