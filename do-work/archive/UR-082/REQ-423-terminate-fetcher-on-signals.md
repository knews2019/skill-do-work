---
id: REQ-423
title: 'Terminate archive fetching on interruption signals'
status: completed
claimed_at: 2026-08-29T20:31:49Z
completed_at: 2026-08-29T20:46:08Z
commit: f532dcf2
route: B
created_at: 2026-08-29T20:26:10Z
user_request: UR-082
domain: backend
prime_files: ['_dev/primes/prime-shell-commands.md']
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-421, REQ-422, REQ-424]
batch: accepted-review-fixes
write_set: ['tools/fetch-upstream-archive.sh', 'skills/do-work/tools/fetch-upstream-archive.sh', '_dev/tests/update-script-behavior.sh', '_dev/primes/prime-shell-commands.md', '_dev/primes/lessons-shell-commands.md']
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
---

# Terminate Archive Fetching on Interruption Signals

## What
Make the archive fetcher stop after HUP, INT, or TERM instead of cleaning up and continuing into fallback success.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Separate EXIT cleanup from signal termination and add one fixture loop that injects each signal before Git fallback can publish.
- [x] **[APPLY]:** Reserved `EXIT` for cleanup and made HUP, INT, and TERM exit 129, 130, and 143 in both fetcher copies; added interruption regressions with a valid fallback available.
- [x] **[UNIFY]:** Reviewed the complete diff, verified shell syntax/lint and mirror equality, ran the focused update behavior suite and the canonical maintainer verifier; all required lanes passed.

## Detailed Requirements
- Reserve `EXIT` for cleanup in both byte-identical fetcher copies.
- Make HUP, INT, and TERM terminate with statuses 129, 130, and 143 respectively.
- Test every signal while a valid Git fallback exists.
- Assert no target archive is published and no success report is emitted after interruption.

## Constraints
- Keep root and shipped fetcher scripts byte-identical.
- Preserve existing HTTP and Git fallback behavior when no signal occurs.
- Add concise signal-preservation guidance to the shell prime and a detailed linked lesson entry.

## Dependencies
None. It shares fetcher and shell-test files with REQ-424 and is implemented as one shell slice in this batch.

## Builder Guidance
Certainty: Firm. Conventional `128 + signal` statuses are explicitly required.

## Context
No pending or unassigned queue candidate shares this root cause. Provenance: accepted review finding `[P2] Exit after handling fetch interruption signals` against `skills/do-work/tools/fetch-upstream-archive.sh:34`. The review states that a cleanup-only trap returns and permits fallback or false success after cancellation.

## Red-Green Proof
**RED prompt/case:** Inject HUP, INT, and TERM during HTTP fetching while a working Git fallback is configured.
**Why RED now:** Each cleanup trap returns, allowing execution to continue and possibly publish/report success.
**GREEN when:** Each case exits 129/130/143, publishes no target, and reports no success.
**Validation:** User accepted the review finding and supplied the implementation plan.

## Full Context
See `do-work/user-requests/UR-082/input.md` for the approved plan and batch constraints.

---
*Source: accepted review finding [P2] on fetch interruption, followed by the user-approved plan.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The trap fix is small, while faithful signal injection and publication assertions require the existing shell integration harness.

**Planning:** Not required; the user supplied an implementation plan.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `tools/fetch-upstream-archive.sh` (modify) — terminating signal traps
- `skills/do-work/tools/fetch-upstream-archive.sh` (modify) — byte-identical shipped mirror
- `_dev/tests/update-script-behavior.sh` (modify) — HUP/INT/TERM regressions
- `_dev/primes/prime-shell-commands.md` (modify) — concise signal/ref guidance shared with REQ-424
- `_dev/primes/lessons-shell-commands.md` (modify) — linked detailed lessons shared with REQ-424

**Files I will NOT touch:** installer/updater caller interfaces.

**Acceptance criteria (restated from REQ):**
- [x] HUP/INT/TERM exit 129/130/143.
- [x] Interrupted fetches publish no archive and print no success.
- [x] Fetcher mirrors remain byte-identical.

## Implementation Summary

**Files changed:**
- `tools/fetch-upstream-archive.sh` (modified)
- `skills/do-work/tools/fetch-upstream-archive.sh` (modified)
- `_dev/tests/update-script-behavior.sh` (modified)
- `_dev/primes/prime-shell-commands.md` (modified)

**What was done:** Cleanup remains centralized in the EXIT trap, while each interruption trap terminates with its conventional status. Integration tests inject every signal while Git fallback could otherwise succeed and prove no archive or success message escapes.

## Testing

**Tests run:** `bash _dev/tests/update-script-behavior.sh`; Bash syntax and ShellCheck; fetcher `cmp`; `bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/maintainer-verify.sh`.

**Result:** All passed; both mirrors compare byte-identical and the canonical verifier exited 0.

**Red-green validation:**
- HUP/INT/TERM cases: RED with status 0, a published target, and success output for every signal → GREEN with statuses 129/130/143, no target, and no success output.

**New tests added:**
- Three signal interruption cases with a valid Git fallback present.

## Lessons Learned

An `EXIT` trap is a cleanup hook, while a signal trap is a control-flow decision. Reusing a cleanup function directly for HUP/INT/TERM lets the shell return to the interrupted workflow; the regression must keep a valid success fallback available so continuation cannot hide behind an already failing route.
