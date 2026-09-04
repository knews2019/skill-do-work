---
id: REQ-550
title: '[impact-negligible] Collapse four exported one-line Go delegates into their targets'
status: claimed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: backend
prime_files: []
tdd: true
suggested_spec:
depends_on: []
related: [REQ-549, REQ-551, REQ-552, REQ-553, REQ-554, REQ-555, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go, skills/do-work/tools/do-work-cli/internal/doctor/doctor_repair.go, skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch.go, skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-04T13:33:39Z
---

# Collapse four exported one-line Go delegates into their targets

## What
Four exported functions in do-work-cli have a single-expression body delegating to a private same-package function and no production caller (only tests call them). Collapse each delegate into its target under the exported name so the tests compile unchanged; one function per job.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
A wrapper whose only job is to call one other thing is surface without behaviour; all four live under `internal/`, so nothing outside the module can depend on them.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 10, sweep_key `exported-delegate-no-production-caller`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag SIMPLE; expected net line delta -12. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `internal/corehelpers/inventory.go` — `AssociateProjectPaths` delegates to `associatePaths` with an identical argument list.
- `internal/doctor/doctor_repair.go` — `ApplyTimestampPlan` delegates to `applyTimestampPlan(...).result`.
- `internal/archivefetch/archive_fetch.go` — `DownloadAtomic` delegates to `downloadAtomicValidated(..., nil)`.
- `internal/gateevidence/gate_evidence.go` — `CheckGreenGate` delegates to `checkGreenGateAtRevision(..., "")`.
- For each: keep the exported name, move the private body under it (or make the private function the exported one), delete the other; `go build ./... && go test ./...` in the module stays green with no test edits.
- Reproduce at dc8a64e3 (prints the four names): the Finding 10 Reproduce command in `do-work/audits/audit-2026-09-03.md`.

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (create it on first use, executable, invoked from `_dev/tests/contract-regressions.sh` in the fast tier the way `_dev/tests/defensive-surface-audit.sh` is, with the same missing-or-not-executable FAIL line), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- Tests unchanged: if a test breaks, the collapse was done in the wrong direction.
- Lock-in limit: exported one-line delegates with no production caller: 0 after this REQ (today 4).

## Dependencies
No dependency. REQ-552 (replace two coreutils exec sites) and REQ-557 (deduplicate Go helpers) run after it so the `corehelpers` export surface is settled first.

## Builder Guidance
Firm on the shape (one function per job), latitude on which name survives as long as the exported one does.

## Red-Green Proof
**RED prompt/case:** Run the Finding 10 Reproduce command from the audit report.
**Why RED now:** It prints four names (`AssociateProjectPaths`, `ApplyTimestampPlan`, `DownloadAtomic`, `CheckGreenGate`).
**GREEN when:** It prints nothing; `go test ./...` in `skills/do-work/tools/do-work-cli` is green with no test file changed; the lock-in pins exported single-expression delegates with no production caller at 0.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5660 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "do-work-cli internals".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for exported-delegate-no-production-caller.*
